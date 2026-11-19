package database

import (
	"database/sql"
	"errors"
	"fmt"
	"path/filepath"
	"sync"
	"testing"
	"time"
)

func testDatabase(t *testing.T) (AppDatabase, *sql.DB) {
	t.Helper()
	c, err := Open(filepath.Join(t.TempDir(), "test ?& database.db"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { c.Close() })
	db, err := New(c)
	if err != nil {
		t.Fatal(err)
	}
	return db, c
}
func testUsers(t *testing.T, db AppDatabase) []int {
	t.Helper()
	ids := make([]int, 0, 3)
	for _, name := range []string{"alice", "bobby", "carol"} {
		id, err := db.GetOrCreateUser(name)
		if err != nil {
			t.Fatal(err)
		}
		ids = append(ids, id)
	}
	return ids
}
func countRows(t *testing.T, c *sql.DB, table string) int {
	t.Helper()
	var count int
	if err := c.QueryRow("SELECT count(*) FROM " + table).Scan(&count); err != nil {
		t.Fatal(err)
	}
	return count
}

func TestNormalizedSchemaConstraints(t *testing.T) {
	db, c := testDatabase(t)
	users := testUsers(t, db)
	chat, err := db.CreateConversation(users[0], users[1])
	if err != nil {
		t.Fatal(err)
	}
	if _, err = db.CreateConversation(users[1], users[0]); err == nil {
		t.Fatal("duplicate private conversation allowed")
	}
	if _, err = c.Exec("INSERT INTO users_table(name) VALUES('alice')"); err == nil {
		t.Fatal("duplicate username allowed")
	}
	if _, err = c.Exec("INSERT INTO users_table(name) VALUES(NULL)"); err == nil {
		t.Fatal("NULL username allowed")
	}
	if _, err = c.Exec("INSERT INTO conversation_members VALUES(?,?)", chat, users[0]); err == nil {
		t.Fatal("duplicate membership allowed")
	}
	if _, err = c.Exec("INSERT INTO conversation_members VALUES(?,?)", chat, users[2]); err == nil {
		t.Fatal("third participant in private chat allowed")
	}
	if _, err = c.Exec("INSERT INTO conversation_members VALUES(99999,?)", users[0]); err == nil {
		t.Fatal("orphan membership allowed")
	}
	if err = db.AddUserToGroup(users[2], chat); !errors.Is(err, sql.ErrNoRows) {
		t.Fatalf("private conversation accepted group operation: %v", err)
	}
	if _, _, err = db.CreateMessage("not a member", chat, users[2]); err == nil {
		t.Fatal("outsider message allowed")
	}
	for table, removed := range map[string][]string{"users_table": {"conversation_ids"}, "conversations_table": {"members"}, "messages_table": {"comments", "hours", "author"}} {
		rows, err := c.Query("PRAGMA table_info(" + table + ")")
		if err != nil {
			t.Fatal(err)
		}
		for rows.Next() {
			var cid, notnull, pk int
			var name, kind string
			var defaultValue interface{}
			if err = rows.Scan(&cid, &name, &kind, &notnull, &defaultValue, &pk); err != nil {
				t.Fatal(err)
			}
			for _, old := range removed {
				if name == old {
					t.Errorf("%s still has %s", table, old)
				}
			}
		}
		if err = rows.Err(); err != nil {
			t.Fatal(err)
		}
		rows.Close()
	}
}

func TestForeignKeysSurviveConnectionReplacement(t *testing.T) {
	db, c := testDatabase(t)
	users := testUsers(t, db)
	group, err := db.CreateGroup("Team", users[1:], users[0])
	if err != nil {
		t.Fatal(err)
	}
	c.SetMaxIdleConns(0)
	for i := 0; i < 3; i++ {
		var enabled int
		if err = c.QueryRow("PRAGMA foreign_keys").Scan(&enabled); err != nil || enabled != 1 {
			t.Fatalf("foreign keys disabled: %d, %v", enabled, err)
		}
		if err = db.AddUserToGroup(999999, group); err == nil {
			t.Fatal("orphan membership allowed on replacement connection")
		}
	}
}

func TestConversationCreationRollsBackEveryMember(t *testing.T) {
	db, c := testDatabase(t)
	users := testUsers(t, db)
	for _, members := range [][]int{{users[1], 999999}, {users[1], users[1]}, {users[0], users[1]}} {
		if _, err := db.CreateGroup("Broken", members, users[0]); err == nil {
			t.Fatal("invalid group accepted")
		}
		if countRows(t, c, "conversations_table") != 0 || countRows(t, c, "conversation_members") != 0 {
			t.Fatal("failed group creation left partial data")
		}
	}
	_, err := c.Exec(fmt.Sprintf(`CREATE TRIGGER fail_second_member BEFORE INSERT ON conversation_members
 WHEN NEW.user_id=%d BEGIN SELECT RAISE(ABORT,'injected member failure'); END`, users[1]))
	if err != nil {
		t.Fatal(err)
	}
	if _, err = db.CreateConversation(users[0], users[1]); err == nil {
		t.Fatal("injected failure not returned")
	}
	if countRows(t, c, "conversations_table") != 0 || countRows(t, c, "conversation_members") != 0 {
		t.Fatal("private creation did not roll back")
	}
}

func TestTimestampsReactionsAndForwardingPersist(t *testing.T) {
	db, c := testDatabase(t)
	users := testUsers(t, db)
	source, err := db.CreateConversation(users[0], users[1])
	if err != nil {
		t.Fatal(err)
	}
	target, err := db.CreateConversation(users[0], users[2])
	if err != nil {
		t.Fatal(err)
	}
	original, originalID, err := db.CreateMessage("snapshot", source, users[0])
	if err != nil {
		t.Fatal(err)
	}
	if original.CreatedAt.IsZero() || time.Since(original.CreatedAt) > time.Minute || original.Forwarded || original.OriginalMessageID != nil || original.MessageId != originalID {
		t.Fatalf("invalid metadata: %#v", original)
	}
	first, created, err := db.SetReaction(users[1], "👍", originalID)
	if err != nil {
		t.Fatal(err)
	}
	updated, updatedCreated, err := db.SetReaction(users[1], "✨", originalID)
	if err != nil {
		t.Fatal(err)
	}
	if !created || updatedCreated || updated.CommentId != first.CommentId || countRows(t, c, "comments_table") != 1 {
		t.Fatal("reaction replacement duplicated a reaction")
	}
	if _, _, err = db.SetReaction(users[2], "✨", originalID); err == nil {
		t.Fatal("outsider reaction allowed")
	}
	if _, _, err = db.ForwardMessage(originalID, target, users[1]); err == nil {
		t.Fatal("forward into inaccessible target allowed")
	}
	if countRows(t, c, "messages_table") != 1 {
		t.Fatal("failed forward left a message")
	}
	forwarded, forwardID, err := db.ForwardMessage(originalID, target, users[0])
	if err != nil {
		t.Fatal(err)
	}
	if !forwarded.Forwarded || forwarded.OriginalMessageID == nil || *forwarded.OriginalMessageID != originalID || forwarded.Text != "snapshot" {
		t.Fatalf("invalid forward: %#v", forwarded)
	}
	// A new database wrapper applies no migration twice and reads persisted metadata.
	reopened, err := New(c)
	if err != nil {
		t.Fatal(err)
	}
	read, err := reopened.GetMessage(forwardID)
	if err != nil {
		t.Fatal(err)
	}
	if !read.Forwarded || !read.CreatedAt.Equal(forwarded.CreatedAt) || read.OriginalMessageID == nil {
		t.Fatal("forward metadata was not stored")
	}
	if err = db.SetName("renamed", users[0]); err != nil {
		t.Fatal(err)
	}
	read, err = db.GetMessage(originalID)
	if err != nil {
		t.Fatal(err)
	}
	if read.Author != "renamed" || len(read.Comments) != 1 || read.Comments[0].Emoji != "✨" {
		t.Fatalf("relations not reflected: %#v", read)
	}
	if err = db.RemoveMessage(originalID); err != nil {
		t.Fatal(err)
	}
	if countRows(t, c, "comments_table") != 0 {
		t.Fatal("source reactions not cascaded")
	}
	read, err = db.GetMessage(forwardID)
	if err != nil {
		t.Fatal(err)
	}
	if !read.Forwarded || read.OriginalMessageID != nil || read.Text != "snapshot" {
		t.Fatal("source deletion lost forwarded content or flag")
	}
}

func TestMessageOrderingUsesFullTimestamp(t *testing.T) {
	db, c := testDatabase(t)
	users := testUsers(t, db)
	chat, err := db.CreateConversation(users[0], users[1])
	if err != nil {
		t.Fatal(err)
	}
	_, first, err := db.CreateMessage("first inserted", chat, users[0])
	if err != nil {
		t.Fatal(err)
	}
	_, second, err := db.CreateMessage("earlier day", chat, users[0])
	if err != nil {
		t.Fatal(err)
	}
	if _, err = c.Exec("UPDATE messages_table SET created_at=? WHERE message_id=?", "2026-11-08T00:01:00.000000000Z", first); err != nil {
		t.Fatal(err)
	}
	if _, err = c.Exec("UPDATE messages_table SET created_at=? WHERE message_id=?", "2026-11-07T23:59:00.000000000Z", second); err != nil {
		t.Fatal(err)
	}
	messages, err := db.GetConversationMessages(chat)
	if err != nil {
		t.Fatal(err)
	}
	if len(messages) != 2 || messages[0].MessageId != second || messages[1].MessageId != first {
		t.Fatalf("timestamp ordering failed: %#v", messages)
	}
}

func TestLeavingLastGroupMemberIsAtomic(t *testing.T) {
	db, c := testDatabase(t)
	users := testUsers(t, db)
	group, err := db.CreateGroup("Team", users[1:], users[0])
	if err != nil {
		t.Fatal(err)
	}
	_, messageID, err := db.CreateMessage("group message", group, users[0])
	if err != nil {
		t.Fatal(err)
	}
	if _, _, err = db.SetReaction(users[1], "👍", messageID); err != nil {
		t.Fatal(err)
	}
	for _, id := range users[1:] {
		if err = db.LeaveGroup(id, group); err != nil {
			t.Fatal(err)
		}
	}
	if _, err = c.Exec("CREATE TRIGGER fail_group_delete BEFORE DELETE ON conversations_table BEGIN SELECT RAISE(ABORT,'injected delete failure'); END"); err != nil {
		t.Fatal(err)
	}
	if err = db.LeaveGroup(users[0], group); err == nil {
		t.Fatal("injected failure not returned")
	}
	member, err := db.CheckUser_Conversation(group, users[0])
	if err != nil || !member {
		t.Fatal("failed deletion did not restore membership")
	}
	if countRows(t, c, "messages_table") != 1 || countRows(t, c, "comments_table") != 1 {
		t.Fatal("failed deletion lost chat history")
	}
	if _, err = c.Exec("DROP TRIGGER fail_group_delete"); err != nil {
		t.Fatal(err)
	}
	if err = db.LeaveGroup(users[0], group); err != nil {
		t.Fatal(err)
	}
	for _, table := range []string{"conversations_table", "conversation_members", "messages_table", "comments_table"} {
		if countRows(t, c, table) != 0 {
			t.Fatalf("last member left orphan rows in %s", table)
		}
	}
}

func TestConcurrentLoginsAndPrivateConversationUniqueness(t *testing.T) {
	db, _ := testDatabase(t)
	var wg sync.WaitGroup
	ids := make(chan int, 12)
	errs := make(chan error, 12)
	for i := 0; i < 12; i++ {
		wg.Add(1)
		go func() { defer wg.Done(); id, err := db.GetOrCreateUser("alice"); ids <- id; errs <- err }()
	}
	wg.Wait()
	close(ids)
	close(errs)
	for err := range errs {
		if err != nil {
			t.Fatal(err)
		}
	}
	var first int
	for id := range ids {
		if first == 0 {
			first = id
		}
		if id != first {
			t.Fatal("concurrent login produced multiple users")
		}
	}
	bob, err := db.GetOrCreateUser("bobby")
	if err != nil {
		t.Fatal(err)
	}
	results := make(chan error, 12)
	for i := 0; i < 12; i++ {
		wg.Add(1)
		go func() { defer wg.Done(); _, err := db.CreateConversation(first, bob); results <- err }()
	}
	wg.Wait()
	close(results)
	successes := 0
	for err := range results {
		if err == nil {
			successes++
		}
	}
	if successes != 1 {
		t.Fatalf("created %d duplicate private conversations", successes)
	}
}
