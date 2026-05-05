CREATE TABLE sys_customer (
  id BIGSERIAL PRIMARY KEY,
  name VARCHAR(128) NOT NULL,
  contact_name VARCHAR(64) NOT NULL DEFAULT '',
  phone VARCHAR(32) NOT NULL DEFAULT '',
  level VARCHAR(32) NOT NULL DEFAULT '',
  source VARCHAR(32) NOT NULL DEFAULT '',
  department_id BIGINT NOT NULL DEFAULT 0,
  owner_user_id BIGINT NOT NULL DEFAULT 0,
  status SMALLINT NOT NULL DEFAULT 1,
  remark VARCHAR(255) NOT NULL DEFAULT '',
  created_at TIMESTAMPTZ NOT NULL,
  updated_at TIMESTAMPTZ NOT NULL,
  deleted_at TIMESTAMPTZ NULL
);

CREATE INDEX idx_sys_customer_level ON sys_customer (level);
CREATE INDEX idx_sys_customer_source ON sys_customer (source);
CREATE INDEX idx_sys_customer_department_id ON sys_customer (department_id);
CREATE INDEX idx_sys_customer_owner_user_id ON sys_customer (owner_user_id);
CREATE INDEX idx_sys_customer_status ON sys_customer (status);
CREATE INDEX idx_sys_customer_deleted_at ON sys_customer (deleted_at);

COMMENT ON TABLE sys_customer IS 'CRM 客户档案表';
COMMENT ON COLUMN sys_customer.id IS '客户记录主键，数据库自增生成';
COMMENT ON COLUMN sys_customer.name IS '客户名称';
COMMENT ON COLUMN sys_customer.contact_name IS '联系人姓名';
COMMENT ON COLUMN sys_customer.phone IS '联系电话';
COMMENT ON COLUMN sys_customer.level IS '客户等级，例如 a、b、vip';
COMMENT ON COLUMN sys_customer.source IS '客户来源，例如 referral、ads、offline';
COMMENT ON COLUMN sys_customer.department_id IS '归属部门 ID，对应 sys_department.id';
COMMENT ON COLUMN sys_customer.owner_user_id IS '负责人用户 ID，对应 sys_user.id';
COMMENT ON COLUMN sys_customer.status IS '客户状态：1 启用，2 停用';
COMMENT ON COLUMN sys_customer.remark IS '备注';
COMMENT ON COLUMN sys_customer.created_at IS '创建时间';
COMMENT ON COLUMN sys_customer.updated_at IS '更新时间';
COMMENT ON COLUMN sys_customer.deleted_at IS '逻辑删除时间，NULL 表示未删除';
