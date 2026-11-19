package database

import (
	"WASA_project2024/service/media"
	"encoding/base64"
	"io/fs"
	"path/filepath"
	"strings"
	"testing"
	"testing/fstest"
)

func TestVersionOneMediaUpgradePreservesHistory(t *testing.T) {
	c, err := Open(filepath.Join(t.TempDir(), "upgrade.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()
	script, _ := fs.ReadFile(migrationFiles, "migrations/001_initial.sql")
	if err = migrate(c, fstest.MapFS{"migrations/001_initial.sql": {Data: script}}); err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"alice", "bobby"} {
		if _, err = c.Exec("INSERT INTO users_table(name,photo) VALUES(?,?)", name, DefaultContactPhoto); err != nil {
			t.Fatal(err)
		}
	}
	if _, err = c.Exec("INSERT INTO conversations_table(title,photo,is_group,direct_user_low,direct_user_high) VALUES('',?,0,1,2)", DefaultContactPhoto); err != nil {
		t.Fatal(err)
	}
	if _, err = c.Exec("INSERT INTO conversation_members VALUES(1,1),(1,2)"); err != nil {
		t.Fatal(err)
	}
	if _, err = c.Exec("INSERT INTO messages_table(conversation_id,user_id,text,created_at) VALUES(1,1,?,'2026-11-08T12:00:00.000000000Z')", "[foto]"+DefaultContactPhoto); err != nil {
		t.Fatal(err)
	}
	if _, err = c.Exec("INSERT INTO comments_table(message_id,author_id,emoji) VALUES(1,2,'👍')"); err != nil {
		t.Fatal(err)
	}
	db, err := New(c)
	if err != nil {
		t.Fatal(err)
	}
	message, err := db.GetMessage(1)
	if err != nil {
		t.Fatal(err)
	}
	if message.Type != "image" || message.Image == nil || len(message.Comments) != 1 || message.Author != "alice" {
		t.Fatal("history lost during upgrade")
	}
	user, err := db.GetUser(1)
	if err != nil || !strings.HasPrefix(user.Photo, "/media/") {
		t.Fatal("profile was not converted")
	}
	img, err := db.GetImage(message.Image.ImageID)
	if err != nil {
		t.Fatal(err)
	}
	original, _ := base64.StdEncoding.DecodeString(DefaultContactPhoto)
	if string(img.Data) != string(original) || img.MimeType != "image/png" {
		t.Fatal("image bytes changed")
	}
	if _, err = New(c); err != nil {
		t.Fatal(err)
	}
	if countRows(t, c, "images") != 1 {
		t.Fatal("identical legacy images were not deduplicated")
	}
	if _, _, err = db.CreateImageMessage(media.Image{}, 1, 1); err == nil {
		t.Fatal("invalid image bypassed SQL constraints")
	}
}

func TestGroupCascadeWithForwardReferencesDoesNotRecreateChanges(t *testing.T) {
	db, c := testDatabase(t)
	users := testUsers(t, db)
	group, err := db.CreateGroup("Group", users[1:], users[0])
	if err != nil {
		t.Fatal(err)
	}
	_, source, err := db.CreateMessage("source", group, users[0])
	if err != nil {
		t.Fatal(err)
	}
	if _, _, err = db.ForwardMessage(source, group, users[1]); err != nil {
		t.Fatal(err)
	}
	for _, user := range users {
		if err = db.LeaveGroup(user, group); err != nil {
			t.Fatal(err)
		}
	}
	for _, table := range []string{"conversations_table", "messages_table", "comments_table", "message_changes", "conversation_members"} {
		if countRows(t, c, table) != 0 {
			t.Fatalf("orphan records in %s", table)
		}
	}
}
