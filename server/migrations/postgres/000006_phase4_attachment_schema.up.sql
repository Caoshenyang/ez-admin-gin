CREATE TABLE sys_attachment (
  id BIGSERIAL PRIMARY KEY,
  file_id BIGINT NOT NULL,
  display_name VARCHAR(255) NOT NULL,
  category VARCHAR(64) NOT NULL DEFAULT '',
  biz_type VARCHAR(64) NOT NULL DEFAULT '',
  uploader_id BIGINT NOT NULL DEFAULT 0,
  status SMALLINT NOT NULL DEFAULT 1,
  remark VARCHAR(255) NOT NULL DEFAULT '',
  created_at TIMESTAMPTZ NOT NULL,
  updated_at TIMESTAMPTZ NOT NULL,
  deleted_at TIMESTAMPTZ NULL
);

CREATE UNIQUE INDEX uk_sys_attachment_file_id ON sys_attachment (file_id);
CREATE INDEX idx_sys_attachment_category ON sys_attachment (category);
CREATE INDEX idx_sys_attachment_biz_type ON sys_attachment (biz_type);
CREATE INDEX idx_sys_attachment_uploader_id ON sys_attachment (uploader_id);
CREATE INDEX idx_sys_attachment_status ON sys_attachment (status);
CREATE INDEX idx_sys_attachment_deleted_at ON sys_attachment (deleted_at);

COMMENT ON TABLE sys_attachment IS '附件中心表';
COMMENT ON COLUMN sys_attachment.id IS '附件记录主键，数据库自增生成';
COMMENT ON COLUMN sys_attachment.file_id IS '底层文件记录 ID，对应 sys_file.id';
COMMENT ON COLUMN sys_attachment.display_name IS '附件展示名，默认使用原始文件名';
COMMENT ON COLUMN sys_attachment.category IS '附件分类，例如 invoice、contract、avatar';
COMMENT ON COLUMN sys_attachment.biz_type IS '业务类型，用于区分模块接入来源';
COMMENT ON COLUMN sys_attachment.uploader_id IS '上传用户 ID，对应 sys_user.id';
COMMENT ON COLUMN sys_attachment.status IS '附件状态：1 启用，2 停用';
COMMENT ON COLUMN sys_attachment.remark IS '备注';
COMMENT ON COLUMN sys_attachment.created_at IS '创建时间';
COMMENT ON COLUMN sys_attachment.updated_at IS '更新时间';
COMMENT ON COLUMN sys_attachment.deleted_at IS '逻辑删除时间，NULL 表示未删除';
