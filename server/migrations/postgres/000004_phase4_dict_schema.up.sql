CREATE TABLE sys_dict_type (
  id BIGSERIAL PRIMARY KEY,
  code VARCHAR(64) NOT NULL,
  name VARCHAR(64) NOT NULL,
  sort INTEGER NOT NULL DEFAULT 0,
  status SMALLINT NOT NULL DEFAULT 1,
  remark VARCHAR(255) NOT NULL DEFAULT '',
  created_at TIMESTAMPTZ NOT NULL,
  updated_at TIMESTAMPTZ NOT NULL,
  deleted_at TIMESTAMPTZ NULL
);

CREATE UNIQUE INDEX uk_sys_dict_type_code ON sys_dict_type (code);
CREATE INDEX idx_sys_dict_type_status ON sys_dict_type (status);
CREATE INDEX idx_sys_dict_type_deleted_at ON sys_dict_type (deleted_at);

COMMENT ON TABLE sys_dict_type IS '系统字典类型表';
COMMENT ON COLUMN sys_dict_type.id IS '字典类型主键，数据库自增生成';
COMMENT ON COLUMN sys_dict_type.code IS '字典编码，系统内唯一';
COMMENT ON COLUMN sys_dict_type.name IS '字典名称';
COMMENT ON COLUMN sys_dict_type.sort IS '排序值，数字越小越靠前';
COMMENT ON COLUMN sys_dict_type.status IS '字典状态：1 启用，2 禁用';
COMMENT ON COLUMN sys_dict_type.remark IS '备注';
COMMENT ON COLUMN sys_dict_type.created_at IS '创建时间';
COMMENT ON COLUMN sys_dict_type.updated_at IS '更新时间';
COMMENT ON COLUMN sys_dict_type.deleted_at IS '逻辑删除时间，NULL 表示未删除';

CREATE TABLE sys_dict_item (
  id BIGSERIAL PRIMARY KEY,
  type_id BIGINT NOT NULL,
  item_key VARCHAR(64) NOT NULL,
  label VARCHAR(64) NOT NULL,
  value VARCHAR(255) NOT NULL,
  tag_type VARCHAR(32) NOT NULL DEFAULT '',
  sort INTEGER NOT NULL DEFAULT 0,
  status SMALLINT NOT NULL DEFAULT 1,
  remark VARCHAR(255) NOT NULL DEFAULT '',
  created_at TIMESTAMPTZ NOT NULL,
  updated_at TIMESTAMPTZ NOT NULL,
  deleted_at TIMESTAMPTZ NULL
);

CREATE UNIQUE INDEX uk_sys_dict_item_type_key ON sys_dict_item (type_id, item_key);
CREATE INDEX idx_sys_dict_item_type_id ON sys_dict_item (type_id);
CREATE INDEX idx_sys_dict_item_status ON sys_dict_item (status);
CREATE INDEX idx_sys_dict_item_deleted_at ON sys_dict_item (deleted_at);

COMMENT ON TABLE sys_dict_item IS '系统字典项表';
COMMENT ON COLUMN sys_dict_item.id IS '字典项主键，数据库自增生成';
COMMENT ON COLUMN sys_dict_item.type_id IS '字典类型 ID，对应 sys_dict_type.id';
COMMENT ON COLUMN sys_dict_item.item_key IS '字典项编码，同一类型内唯一';
COMMENT ON COLUMN sys_dict_item.label IS '字典项名称';
COMMENT ON COLUMN sys_dict_item.value IS '字典项值';
COMMENT ON COLUMN sys_dict_item.tag_type IS '前端标签样式提示';
COMMENT ON COLUMN sys_dict_item.sort IS '排序值，数字越小越靠前';
COMMENT ON COLUMN sys_dict_item.status IS '字典项状态：1 启用，2 禁用';
COMMENT ON COLUMN sys_dict_item.remark IS '备注';
COMMENT ON COLUMN sys_dict_item.created_at IS '创建时间';
COMMENT ON COLUMN sys_dict_item.updated_at IS '更新时间';
COMMENT ON COLUMN sys_dict_item.deleted_at IS '逻辑删除时间，NULL 表示未删除';
