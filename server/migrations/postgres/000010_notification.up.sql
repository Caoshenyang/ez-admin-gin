-- +migrate Up
CREATE TABLE sys_notification (
    id         BIGSERIAL PRIMARY KEY,
    user_id    BIGINT NOT NULL,
    type       SMALLINT NOT NULL DEFAULT 1,
    title      VARCHAR(128) NOT NULL DEFAULT '',
    content    TEXT NOT NULL DEFAULT '',
    extra      JSONB,
    is_read    BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    read_at    TIMESTAMP
);

CREATE INDEX idx_notification_user_unread ON sys_notification(user_id, is_read, created_at DESC);
