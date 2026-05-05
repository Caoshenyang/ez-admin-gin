CREATE TABLE sys_customer_followup (
  id BIGSERIAL PRIMARY KEY,
  customer_id BIGINT NOT NULL,
  department_id BIGINT NOT NULL DEFAULT 0,
  owner_user_id BIGINT NOT NULL DEFAULT 0,
  follow_type VARCHAR(32) NOT NULL DEFAULT '',
  subject VARCHAR(128) NOT NULL,
  content VARCHAR(1000) NOT NULL,
  result VARCHAR(255) NOT NULL DEFAULT '',
  next_follow_at TIMESTAMPTZ NULL,
  status SMALLINT NOT NULL DEFAULT 1,
  created_at TIMESTAMPTZ NOT NULL,
  updated_at TIMESTAMPTZ NOT NULL,
  deleted_at TIMESTAMPTZ NULL
);

CREATE INDEX idx_sys_customer_followup_customer_id ON sys_customer_followup (customer_id);
CREATE INDEX idx_sys_customer_followup_department_id ON sys_customer_followup (department_id);
CREATE INDEX idx_sys_customer_followup_owner_user_id ON sys_customer_followup (owner_user_id);
CREATE INDEX idx_sys_customer_followup_follow_type ON sys_customer_followup (follow_type);
CREATE INDEX idx_sys_customer_followup_status ON sys_customer_followup (status);
CREATE INDEX idx_sys_customer_followup_next_follow_at ON sys_customer_followup (next_follow_at);
CREATE INDEX idx_sys_customer_followup_deleted_at ON sys_customer_followup (deleted_at);

COMMENT ON TABLE sys_customer_followup IS 'CRM 客户跟进记录表';
COMMENT ON COLUMN sys_customer_followup.id IS '客户跟进记录主键，数据库自增生成';
COMMENT ON COLUMN sys_customer_followup.customer_id IS '关联客户 ID，对应 sys_customer.id';
COMMENT ON COLUMN sys_customer_followup.department_id IS '继承客户归属部门 ID，对应 sys_department.id';
COMMENT ON COLUMN sys_customer_followup.owner_user_id IS '继承客户负责人 ID，对应 sys_user.id';
COMMENT ON COLUMN sys_customer_followup.follow_type IS '跟进方式，例如 phone、wechat、visit';
COMMENT ON COLUMN sys_customer_followup.subject IS '跟进主题';
COMMENT ON COLUMN sys_customer_followup.content IS '跟进内容';
COMMENT ON COLUMN sys_customer_followup.result IS '跟进结果摘要';
COMMENT ON COLUMN sys_customer_followup.next_follow_at IS '下次计划跟进时间';
COMMENT ON COLUMN sys_customer_followup.status IS '客户跟进状态：1 待跟进，2 已完成，3 已关闭';
COMMENT ON COLUMN sys_customer_followup.created_at IS '创建时间';
COMMENT ON COLUMN sys_customer_followup.updated_at IS '更新时间';
COMMENT ON COLUMN sys_customer_followup.deleted_at IS '逻辑删除时间，NULL 表示未删除';
