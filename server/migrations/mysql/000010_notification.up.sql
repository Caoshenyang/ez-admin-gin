-- +migrate Up
CREATE TABLE sys_notification (
    id         BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    user_id    BIGINT UNSIGNED NOT NULL,
    type       SMALLINT NOT NULL DEFAULT 1,
    title      VARCHAR(128) NOT NULL DEFAULT '',
    content    TEXT NOT NULL DEFAULT '',
    extra      JSON,
    is_read    TINYINT(1) NOT NULL DEFAULT 0,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    read_at    DATETIME NULL,
    INDEX idx_notification_user_unread (user_id, is_read, created_at DESC)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
