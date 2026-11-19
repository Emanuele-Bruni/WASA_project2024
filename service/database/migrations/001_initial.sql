CREATE TABLE users_table (
    id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL UNIQUE CHECK(length(name) BETWEEN 3 AND 16),
    photo TEXT NOT NULL DEFAULT ''
);

CREATE TABLE conversations_table (
    conversation_id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    title TEXT NOT NULL DEFAULT '',
    photo TEXT NOT NULL DEFAULT '',
    is_group INTEGER NOT NULL CHECK(is_group IN (0, 1)),
    direct_user_low INTEGER REFERENCES users_table(id),
    direct_user_high INTEGER REFERENCES users_table(id),
    UNIQUE(direct_user_low, direct_user_high),
    CHECK((is_group = 1 AND direct_user_low IS NULL AND direct_user_high IS NULL AND length(title) BETWEEN 1 AND 25)
       OR (is_group = 0 AND direct_user_low IS NOT NULL AND direct_user_high IS NOT NULL AND direct_user_low < direct_user_high))
);

CREATE INDEX conversations_direct_high ON conversations_table(direct_user_high);

CREATE TABLE conversation_members (
    conversation_id INTEGER NOT NULL REFERENCES conversations_table(conversation_id) ON DELETE CASCADE,
    user_id INTEGER NOT NULL REFERENCES users_table(id),
    PRIMARY KEY(conversation_id, user_id)
);
CREATE INDEX conversation_members_user ON conversation_members(user_id, conversation_id);
CREATE TRIGGER private_members_only BEFORE INSERT ON conversation_members
WHEN EXISTS (
    SELECT 1 FROM conversations_table c WHERE c.conversation_id = NEW.conversation_id
    AND c.is_group = 0 AND NEW.user_id NOT IN (c.direct_user_low, c.direct_user_high)
)
BEGIN
    SELECT RAISE(ABORT, 'private conversation only accepts its two participants');
END;

CREATE TABLE messages_table (
    message_id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    conversation_id INTEGER NOT NULL REFERENCES conversations_table(conversation_id) ON DELETE CASCADE,
    user_id INTEGER NOT NULL REFERENCES users_table(id),
    text TEXT NOT NULL CHECK(length(text) > 0),
    created_at TEXT NOT NULL CHECK(length(created_at) = 30 AND strftime('%Y-%m-%dT%H:%M:%fZ', created_at) IS NOT NULL),
    forwarded INTEGER NOT NULL DEFAULT 0 CHECK(forwarded IN (0, 1)),
    original_message_id INTEGER REFERENCES messages_table(message_id) ON DELETE SET NULL,
    CHECK(forwarded = 1 OR original_message_id IS NULL)
);
CREATE INDEX messages_conversation_time ON messages_table(conversation_id, created_at, message_id);
CREATE INDEX messages_author ON messages_table(user_id);
CREATE INDEX messages_original ON messages_table(original_message_id);

CREATE TABLE comments_table (
    comment_id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    message_id INTEGER NOT NULL REFERENCES messages_table(message_id) ON DELETE CASCADE,
    author_id INTEGER NOT NULL REFERENCES users_table(id),
    emoji TEXT NOT NULL CHECK(length(emoji) > 0),
    UNIQUE(message_id, author_id)
);
CREATE INDEX comments_author ON comments_table(author_id);
