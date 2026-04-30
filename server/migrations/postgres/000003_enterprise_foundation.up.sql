ALTER TABLE sys_user
ADD COLUMN department_id BIGINT NOT NULL DEFAULT 0;

CREATE INDEX idx_sys_user_department_id
ON sys_user (department_id);

ALTER TABLE sys_role
ADD COLUMN data_scope VARCHAR(32) NOT NULL DEFAULT 'self';

CREATE TABLE sys_department (
  id BIGSERIAL PRIMARY KEY,
  parent_id BIGINT NOT NULL DEFAULT 0,
  ancestors VARCHAR(500) NOT NULL DEFAULT '',
  name VARCHAR(64) NOT NULL,
  code VARCHAR(64) NOT NULL,
  leader_user_id BIGINT NOT NULL DEFAULT 0,
  sort INTEGER NOT NULL DEFAULT 0,
  status SMALLINT NOT NULL DEFAULT 1,
  remark VARCHAR(255) NOT NULL DEFAULT '',
  created_at TIMESTAMPTZ NOT NULL,
  updated_at TIMESTAMPTZ NOT NULL,
  deleted_at TIMESTAMPTZ NULL
);

CREATE UNIQUE INDEX uk_sys_department_code
ON sys_department (code);

CREATE INDEX idx_sys_department_parent_id
ON sys_department (parent_id);

CREATE INDEX idx_sys_department_leader_user_id
ON sys_department (leader_user_id);

CREATE INDEX idx_sys_department_status
ON sys_department (status);

CREATE INDEX idx_sys_department_deleted_at
ON sys_department (deleted_at);

COMMENT ON TABLE sys_department IS '组织部门表';
COMMENT ON COLUMN sys_department.parent_id IS '父部门 ID，根节点为 0';
COMMENT ON COLUMN sys_department.ancestors IS '祖先路径，例如 0,1,3';
COMMENT ON COLUMN sys_department.name IS '部门名称';
COMMENT ON COLUMN sys_department.code IS '部门编码';
COMMENT ON COLUMN sys_department.leader_user_id IS '负责人用户 ID';
COMMENT ON COLUMN sys_department.sort IS '排序值，数字越小越靠前';
COMMENT ON COLUMN sys_department.status IS '部门状态：1 启用，2 禁用';
COMMENT ON COLUMN sys_department.remark IS '备注';

CREATE TABLE sys_post (
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

CREATE UNIQUE INDEX uk_sys_post_code
ON sys_post (code);

CREATE INDEX idx_sys_post_status
ON sys_post (status);

CREATE INDEX idx_sys_post_deleted_at
ON sys_post (deleted_at);

COMMENT ON TABLE sys_post IS '岗位表';
COMMENT ON COLUMN sys_post.code IS '岗位编码';
COMMENT ON COLUMN sys_post.name IS '岗位名称';
COMMENT ON COLUMN sys_post.sort IS '排序值，数字越小越靠前';
COMMENT ON COLUMN sys_post.status IS '岗位状态：1 启用，2 禁用';
COMMENT ON COLUMN sys_post.remark IS '备注';

CREATE TABLE sys_user_post (
  id BIGSERIAL PRIMARY KEY,
  user_id BIGINT NOT NULL,
  post_id BIGINT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL,
  updated_at TIMESTAMPTZ NOT NULL
);

CREATE UNIQUE INDEX uk_sys_user_post_user_post
ON sys_user_post (user_id, post_id);

CREATE INDEX idx_sys_user_post_user_id
ON sys_user_post (user_id);

CREATE INDEX idx_sys_user_post_post_id
ON sys_user_post (post_id);

COMMENT ON TABLE sys_user_post IS '用户岗位关系表';
COMMENT ON COLUMN sys_user_post.user_id IS '用户 ID，对应 sys_user.id';
COMMENT ON COLUMN sys_user_post.post_id IS '岗位 ID，对应 sys_post.id';

CREATE TABLE sys_role_data_scope (
  id BIGSERIAL PRIMARY KEY,
  role_id BIGINT NOT NULL,
  department_id BIGINT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL,
  updated_at TIMESTAMPTZ NOT NULL
);

CREATE UNIQUE INDEX uk_sys_role_data_scope_role_department
ON sys_role_data_scope (role_id, department_id);

CREATE INDEX idx_sys_role_data_scope_role_id
ON sys_role_data_scope (role_id);

CREATE INDEX idx_sys_role_data_scope_department_id
ON sys_role_data_scope (department_id);

COMMENT ON TABLE sys_role_data_scope IS '角色自定义部门数据范围关系表';
COMMENT ON COLUMN sys_role_data_scope.role_id IS '角色 ID，对应 sys_role.id';
COMMENT ON COLUMN sys_role_data_scope.department_id IS '部门 ID，对应 sys_department.id';
