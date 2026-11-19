-- Binary media are immutable and deduplicated by the SHA-256 of the validated bytes.
CREATE TABLE images (
    image_id INTEGER PRIMARY KEY AUTOINCREMENT,
    mime_type TEXT NOT NULL CHECK(mime_type IN ('image/png', 'image/jpeg')),
    name TEXT NOT NULL,
    width INTEGER NOT NULL CHECK(width > 0 AND width <= 4096),
    height INTEGER NOT NULL CHECK(height > 0 AND height <= 4096),
    digest TEXT NOT NULL UNIQUE,
    data BLOB NOT NULL CHECK(length(data) BETWEEN 1 AND 1500000),
    CHECK(width * height <= 4194304)
);
ALTER TABLE users_table ADD COLUMN photo_image_id INTEGER REFERENCES images(image_id);
ALTER TABLE conversations_table ADD COLUMN photo_image_id INTEGER REFERENCES images(image_id);
ALTER TABLE messages_table ADD COLUMN content_type TEXT NOT NULL DEFAULT 'text' CHECK(content_type IN ('text', 'image'));
ALTER TABLE messages_table ADD COLUMN image_id INTEGER REFERENCES images(image_id)
    CHECK((content_type = 'text' AND image_id IS NULL) OR (content_type = 'image' AND image_id IS NOT NULL));
ALTER TABLE messages_table ADD COLUMN image_name TEXT NOT NULL DEFAULT '' CHECK(length(image_name) <= 128);
CREATE INDEX users_photo_image ON users_table(photo_image_id);
CREATE INDEX conversations_photo_image ON conversations_table(photo_image_id);
CREATE INDEX messages_image ON messages_table(image_id);

-- No message foreign key here: deleted message IDs must remain in the change stream.
CREATE TABLE message_changes (
    revision INTEGER PRIMARY KEY AUTOINCREMENT,
    conversation_id INTEGER NOT NULL REFERENCES conversations_table(conversation_id) ON DELETE CASCADE,
    message_id INTEGER NOT NULL
);
CREATE INDEX message_changes_conversation ON message_changes(conversation_id, revision);
CREATE TRIGGER message_insert_change AFTER INSERT ON messages_table BEGIN
    INSERT INTO message_changes(conversation_id, message_id) SELECT NEW.conversation_id, NEW.message_id WHERE EXISTS(SELECT 1 FROM conversations_table WHERE conversation_id=NEW.conversation_id);
END;
CREATE TRIGGER message_update_change AFTER UPDATE ON messages_table BEGIN
    INSERT INTO message_changes(conversation_id, message_id) SELECT NEW.conversation_id, NEW.message_id WHERE EXISTS(SELECT 1 FROM conversations_table WHERE conversation_id=NEW.conversation_id);
END;
CREATE TRIGGER message_delete_change AFTER DELETE ON messages_table BEGIN
    INSERT INTO message_changes(conversation_id, message_id) SELECT OLD.conversation_id, OLD.message_id WHERE EXISTS(SELECT 1 FROM conversations_table WHERE conversation_id=OLD.conversation_id);
END;
CREATE TRIGGER reaction_insert_change AFTER INSERT ON comments_table BEGIN
    INSERT INTO message_changes(conversation_id, message_id)
    SELECT conversation_id, message_id FROM messages_table WHERE message_id=NEW.message_id;
END;
CREATE TRIGGER reaction_update_change AFTER UPDATE ON comments_table BEGIN
    INSERT INTO message_changes(conversation_id, message_id)
    SELECT conversation_id, message_id FROM messages_table WHERE message_id=NEW.message_id;
END;
CREATE TRIGGER reaction_delete_change AFTER DELETE ON comments_table BEGIN
    INSERT INTO message_changes(conversation_id, message_id)
    SELECT conversation_id, message_id FROM messages_table WHERE message_id=OLD.message_id AND EXISTS(SELECT 1 FROM conversations_table c WHERE c.conversation_id=messages_table.conversation_id);
END;
CREATE TRIGGER author_name_change AFTER UPDATE OF name ON users_table WHEN OLD.name != NEW.name BEGIN
    INSERT INTO message_changes(conversation_id, message_id)
    SELECT conversation_id, message_id FROM messages_table WHERE user_id=NEW.id;
END;
INSERT INTO message_changes(conversation_id, message_id)
SELECT conversation_id, message_id FROM messages_table ORDER BY message_id;
-- migrateLegacyMedia performs base64 decoding and image validation in this transaction.
