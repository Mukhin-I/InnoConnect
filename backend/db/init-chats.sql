CREATE TABLE chats (
    id BIGSERIAL PRIMARY KEY,

    -- REQUEST or MEETING
    type VARCHAR(20) NOT NULL,

    -- ID of request or meeting from another service
    related_id BIGINT NOT NULL,

    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE chat_participants (
    chat_id BIGINT NOT NULL,
    user_id BIGINT NOT NULL,
    user_name TEXT NOT NULL,

    PRIMARY KEY (chat_id, user_id),

    FOREIGN KEY (chat_id)
        REFERENCES chats(id)
        ON DELETE CASCADE
);

CREATE TABLE messages (
    id BIGSERIAL PRIMARY KEY,

    chat_id BIGINT NOT NULL,

    -- User ID from User Service
    sender_id BIGINT NOT NULL,

    sender_name TEXT NOT NULL,

    text TEXT NOT NULL,

    sent_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (chat_id)
        REFERENCES chats(id)
        ON DELETE CASCADE
);

CREATE INDEX idx_messages_chat_sent
    ON messages(chat_id, sent_at);

CREATE INDEX idx_chat_participants_user
    ON chat_participants(user_id);