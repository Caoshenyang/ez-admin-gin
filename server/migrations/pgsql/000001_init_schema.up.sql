-- 系统用户表
CREATE TABLE sys_user (
  id BIGSERIAL PRIMARY KEY,
  username VARCHAR(64) NOT NULL,
  password_hash VARCHAR(255) NOT NULL,
  nickname VARCHAR(64) NOT NULL DEFAULT '',
  status SMALLINT NOT NULL DEFAULT 1,
  created_at TIMESTAMPTZ NOT NULL,
  updated_at TIMESTAMPTZ NOT NULL,
  deleted_at TIMESTAMPTZ NULL
);

CREATE UNIQUE INDEX uk_sys_user_username ON sys_user (username);
CREATE INDEX idx_sys_user_deleted_at ON sys_user (deleted_at);

COMMENT ON TABLE sys_user IS '后台用户表';
COMMENT ON COLUMN sys_user.id IS '用户记录主键，数据库自增生成';
COMMENT ON COLUMN sys_user.username IS '登录用户名';
COMMENT ON COLUMN sys_user.password_hash IS '密码哈希';
COMMENT ON COLUMN sys_user.nickname IS '管理台展示名称';
COMMENT ON COLUMN sys_user.status IS '用户状态：1 启用，2 禁用';
COMMENT ON COLUMN sys_user.created_at IS '创建时间';
COMMENT ON COLUMN sys_user.updated_at IS '更新时间';
COMMENT ON COLUMN sys_user.deleted_at IS '逻辑删除时间，NULL 表示未删除';

-- 系统角色表
CREATE TABLE sys_role (
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

CREATE UNIQUE INDEX uk_sys_role_code ON sys_role (code);
CREATE INDEX idx_sys_role_status ON sys_role (status);
CREATE INDEX idx_sys_role_deleted_at ON sys_role (deleted_at);

COMMENT ON TABLE sys_role IS '后台角色表';
COMMENT ON COLUMN sys_role.id IS '角色记录主键，数据库自增生成';
COMMENT ON COLUMN sys_role.code IS '角色编码';
COMMENT ON COLUMN sys_role.name IS '角色名称';
COMMENT ON COLUMN sys_role.sort IS '排序值，数字越小越靠前';
COMMENT ON COLUMN sys_role.status IS '角色状态：1 启用，2 禁用';
COMMENT ON COLUMN sys_role.remark IS '备注';
COMMENT ON COLUMN sys_role.created_at IS '创建时间';
COMMENT ON COLUMN sys_role.updated_at IS '更新时间';
COMMENT ON COLUMN sys_role.deleted_at IS '逻辑删除时间，NULL 表示未删除';

-- 用户角色关系表
CREATE TABLE sys_user_role (
  id BIGSERIAL PRIMARY KEY,
  user_id BIGINT NOT NULL,
  role_id BIGINT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL,
  updated_at TIMESTAMPTZ NOT NULL
);

CREATE UNIQUE INDEX uk_sys_user_role_user_role ON sys_user_role (user_id, role_id);
CREATE INDEX idx_sys_user_role_user_id ON sys_user_role (user_id);
CREATE INDEX idx_sys_user_role_role_id ON sys_user_role (role_id);

COMMENT ON TABLE sys_user_role IS '用户角色关系表';
COMMENT ON COLUMN sys_user_role.id IS '关系记录主键，数据库自增生成';
COMMENT ON COLUMN sys_user_role.user_id IS '用户 ID，对应 sys_user.id';
COMMENT ON COLUMN sys_user_role.role_id IS '角色 ID，对应 sys_role.id';
COMMENT ON COLUMN sys_user_role.created_at IS '绑定时间';
COMMENT ON COLUMN sys_user_role.updated_at IS '更新时间';

-- 系统菜单和按钮表
CREATE TABLE sys_menu (
  id BIGSERIAL PRIMARY KEY,
  parent_id BIGINT NOT NULL DEFAULT 0,
  type SMALLINT NOT NULL,
  code VARCHAR(128) NOT NULL,
  title VARCHAR(64) NOT NULL,
  path VARCHAR(255) NOT NULL DEFAULT '',
  component VARCHAR(255) NOT NULL DEFAULT '',
  icon VARCHAR(64) NOT NULL DEFAULT '',
  sort INTEGER NOT NULL DEFAULT 0,
  status SMALLINT NOT NULL DEFAULT 1,
  remark VARCHAR(255) NOT NULL DEFAULT '',
  created_at TIMESTAMPTZ NOT NULL,
  updated_at TIMESTAMPTZ NOT NULL,
  deleted_at TIMESTAMPTZ NULL
);

CREATE UNIQUE INDEX uk_sys_menu_code ON sys_menu (code);
CREATE INDEX idx_sys_menu_parent_id ON sys_menu (parent_id);
CREATE INDEX idx_sys_menu_type ON sys_menu (type);
CREATE INDEX idx_sys_menu_status ON sys_menu (status);
CREATE INDEX idx_sys_menu_deleted_at ON sys_menu (deleted_at);

COMMENT ON TABLE sys_menu IS '后台菜单和按钮表';
COMMENT ON COLUMN sys_menu.id IS '菜单记录主键，数据库自增生成';
COMMENT ON COLUMN sys_menu.parent_id IS '父级菜单 ID，根节点为 0';
COMMENT ON COLUMN sys_menu.type IS '节点类型：1 目录，2 菜单，3 按钮';
COMMENT ON COLUMN sys_menu.code IS '菜单或按钮编码，系统内唯一';
COMMENT ON COLUMN sys_menu.title IS '展示名称';
COMMENT ON COLUMN sys_menu.path IS '前端路由路径';
COMMENT ON COLUMN sys_menu.component IS '前端组件路径';
COMMENT ON COLUMN sys_menu.icon IS '图标标识';
COMMENT ON COLUMN sys_menu.sort IS '排序值，数字越小越靠前';
COMMENT ON COLUMN sys_menu.status IS '菜单状态：1 启用，2 禁用';
COMMENT ON COLUMN sys_menu.remark IS '备注';
COMMENT ON COLUMN sys_menu.created_at IS '创建时间';
COMMENT ON COLUMN sys_menu.updated_at IS '更新时间';
COMMENT ON COLUMN sys_menu.deleted_at IS '逻辑删除时间，NULL 表示未删除';

-- 角色菜单关系表
CREATE TABLE sys_role_menu (
  id BIGSERIAL PRIMARY KEY,
  role_id BIGINT NOT NULL,
  menu_id BIGINT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL,
  updated_at TIMESTAMPTZ NOT NULL
);

CREATE UNIQUE INDEX uk_sys_role_menu_role_menu ON sys_role_menu (role_id, menu_id);
CREATE INDEX idx_sys_role_menu_role_id ON sys_role_menu (role_id);
CREATE INDEX idx_sys_role_menu_menu_id ON sys_role_menu (menu_id);

COMMENT ON TABLE sys_role_menu IS '角色菜单关系表';
COMMENT ON COLUMN sys_role_menu.id IS '关系记录主键，数据库自增生成';
COMMENT ON COLUMN sys_role_menu.role_id IS '角色 ID，对应 sys_role.id';
COMMENT ON COLUMN sys_role_menu.menu_id IS '菜单 ID，对应 sys_menu.id';
COMMENT ON COLUMN sys_role_menu.created_at IS '绑定时间';
COMMENT ON COLUMN sys_role_menu.updated_at IS '更新时间';

-- 系统配置表
CREATE TABLE sys_config (
  id BIGSERIAL PRIMARY KEY,
  group_code VARCHAR(64) NOT NULL,
  config_key VARCHAR(128) NOT NULL,
  name VARCHAR(64) NOT NULL,
  value TEXT NOT NULL,
  sort INTEGER NOT NULL DEFAULT 0,
  status SMALLINT NOT NULL DEFAULT 1,
  remark VARCHAR(255) NOT NULL DEFAULT '',
  created_at TIMESTAMPTZ NOT NULL,
  updated_at TIMESTAMPTZ NOT NULL,
  deleted_at TIMESTAMPTZ NULL
);

CREATE UNIQUE INDEX uk_sys_config_key ON sys_config (config_key);
CREATE INDEX idx_sys_config_group_code ON sys_config (group_code);
CREATE INDEX idx_sys_config_status ON sys_config (status);
CREATE INDEX idx_sys_config_deleted_at ON sys_config (deleted_at);

COMMENT ON TABLE sys_config IS '系统配置表';
COMMENT ON COLUMN sys_config.id IS '配置记录主键，数据库自增生成';
COMMENT ON COLUMN sys_config.group_code IS '配置分组，例如 site、upload';
COMMENT ON COLUMN sys_config.config_key IS '配置键，系统内唯一，例如 site:title';
COMMENT ON COLUMN sys_config.name IS '配置名称';
COMMENT ON COLUMN sys_config.value IS '配置值，统一按字符串存储';
COMMENT ON COLUMN sys_config.sort IS '排序值，数字越小越靠前';
COMMENT ON COLUMN sys_config.status IS '配置状态：1 启用，2 禁用';
COMMENT ON COLUMN sys_config.remark IS '备注';
COMMENT ON COLUMN sys_config.created_at IS '创建时间';
COMMENT ON COLUMN sys_config.updated_at IS '更新时间';
COMMENT ON COLUMN sys_config.deleted_at IS '逻辑删除时间，NULL 表示未删除';

-- 文件上传记录表
CREATE TABLE sys_file (
  id BIGSERIAL PRIMARY KEY,
  storage VARCHAR(32) NOT NULL DEFAULT 'local',
  original_name VARCHAR(255) NOT NULL,
  file_name VARCHAR(255) NOT NULL,
  ext VARCHAR(32) NOT NULL DEFAULT '',
  mime_type VARCHAR(128) NOT NULL DEFAULT '',
  size BIGINT NOT NULL DEFAULT 0,
  sha256 VARCHAR(64) NOT NULL DEFAULT '',
  path VARCHAR(500) NOT NULL,
  url VARCHAR(500) NOT NULL,
  uploader_id BIGINT NOT NULL DEFAULT 0,
  status SMALLINT NOT NULL DEFAULT 1,
  remark VARCHAR(255) NOT NULL DEFAULT '',
  created_at TIMESTAMPTZ NOT NULL,
  updated_at TIMESTAMPTZ NOT NULL,
  deleted_at TIMESTAMPTZ NULL
);

CREATE INDEX idx_sys_file_ext ON sys_file (ext);
CREATE INDEX idx_sys_file_sha256 ON sys_file (sha256);
CREATE INDEX idx_sys_file_uploader_id ON sys_file (uploader_id);
CREATE INDEX idx_sys_file_status ON sys_file (status);
CREATE INDEX idx_sys_file_deleted_at ON sys_file (deleted_at);

COMMENT ON TABLE sys_file IS '文件上传记录表';
COMMENT ON COLUMN sys_file.id IS '文件记录主键，数据库自增生成';
COMMENT ON COLUMN sys_file.storage IS '存储类型，本节使用 local';
COMMENT ON COLUMN sys_file.original_name IS '用户上传时的原始文件名';
COMMENT ON COLUMN sys_file.file_name IS '后端生成的保存文件名';
COMMENT ON COLUMN sys_file.ext IS '文件后缀，例如 .png、.pdf';
COMMENT ON COLUMN sys_file.mime_type IS '上传请求中的文件 MIME 类型';
COMMENT ON COLUMN sys_file.size IS '文件大小，单位字节';
COMMENT ON COLUMN sys_file.sha256 IS '文件内容 SHA-256 哈希';
COMMENT ON COLUMN sys_file.path IS '服务端保存路径';
COMMENT ON COLUMN sys_file.url IS '前端可访问地址';
COMMENT ON COLUMN sys_file.uploader_id IS '上传用户 ID，对应 sys_user.id';
COMMENT ON COLUMN sys_file.status IS '文件状态：1 启用，2 停用';
COMMENT ON COLUMN sys_file.remark IS '备注';
COMMENT ON COLUMN sys_file.created_at IS '创建时间';
COMMENT ON COLUMN sys_file.updated_at IS '更新时间';
COMMENT ON COLUMN sys_file.deleted_at IS '逻辑删除时间，NULL 表示未删除';

-- 操作日志表
CREATE TABLE sys_operation_log (
  id BIGSERIAL PRIMARY KEY,
  user_id BIGINT NOT NULL DEFAULT 0,
  username VARCHAR(64) NOT NULL DEFAULT '',
  method VARCHAR(10) NOT NULL,
  path VARCHAR(255) NOT NULL,
  route_path VARCHAR(255) NOT NULL DEFAULT '',
  query VARCHAR(1000) NOT NULL DEFAULT '',
  ip VARCHAR(64) NOT NULL DEFAULT '',
  user_agent VARCHAR(500) NOT NULL DEFAULT '',
  status_code INTEGER NOT NULL DEFAULT 0,
  latency_ms BIGINT NOT NULL DEFAULT 0,
  success BOOLEAN NOT NULL DEFAULT TRUE,
  error_message VARCHAR(500) NOT NULL DEFAULT '',
  created_at TIMESTAMPTZ NOT NULL
);

CREATE INDEX idx_sys_operation_log_user_id ON sys_operation_log (user_id);
CREATE INDEX idx_sys_operation_log_username ON sys_operation_log (username);
CREATE INDEX idx_sys_operation_log_method ON sys_operation_log (method);
CREATE INDEX idx_sys_operation_log_path ON sys_operation_log (path);
CREATE INDEX idx_sys_operation_log_route_path ON sys_operation_log (route_path);
CREATE INDEX idx_sys_operation_log_status_code ON sys_operation_log (status_code);
CREATE INDEX idx_sys_operation_log_success ON sys_operation_log (success);
CREATE INDEX idx_sys_operation_log_created_at ON sys_operation_log (created_at);

COMMENT ON TABLE sys_operation_log IS '操作日志表';
COMMENT ON COLUMN sys_operation_log.id IS '操作日志主键，数据库自增生成';
COMMENT ON COLUMN sys_operation_log.user_id IS '操作人 ID，对应 sys_user.id';
COMMENT ON COLUMN sys_operation_log.username IS '操作人用户名';
COMMENT ON COLUMN sys_operation_log.method IS 'HTTP 请求方法';
COMMENT ON COLUMN sys_operation_log.path IS '实际请求路径';
COMMENT ON COLUMN sys_operation_log.route_path IS 'Gin 路由模板';
COMMENT ON COLUMN sys_operation_log.query IS '查询参数';
COMMENT ON COLUMN sys_operation_log.ip IS '客户端 IP';
COMMENT ON COLUMN sys_operation_log.user_agent IS '浏览器或客户端标识';
COMMENT ON COLUMN sys_operation_log.status_code IS 'HTTP 状态码';
COMMENT ON COLUMN sys_operation_log.latency_ms IS '请求耗时，单位毫秒';
COMMENT ON COLUMN sys_operation_log.success IS '是否成功';
COMMENT ON COLUMN sys_operation_log.error_message IS '错误摘要';
COMMENT ON COLUMN sys_operation_log.created_at IS '创建时间';

-- 登录日志表
CREATE TABLE sys_login_log (
  id BIGSERIAL PRIMARY KEY,
  user_id BIGINT NOT NULL DEFAULT 0,
  username VARCHAR(64) NOT NULL DEFAULT '',
  status SMALLINT NOT NULL,
  message VARCHAR(255) NOT NULL DEFAULT '',
  ip VARCHAR(64) NOT NULL DEFAULT '',
  user_agent VARCHAR(500) NOT NULL DEFAULT '',
  created_at TIMESTAMPTZ NOT NULL
);

CREATE INDEX idx_sys_login_log_user_id ON sys_login_log (user_id);
CREATE INDEX idx_sys_login_log_username ON sys_login_log (username);
CREATE INDEX idx_sys_login_log_status ON sys_login_log (status);
CREATE INDEX idx_sys_login_log_ip ON sys_login_log (ip);
CREATE INDEX idx_sys_login_log_created_at ON sys_login_log (created_at);

COMMENT ON TABLE sys_login_log IS '登录日志表';
COMMENT ON COLUMN sys_login_log.id IS '登录日志主键，数据库自增生成';
COMMENT ON COLUMN sys_login_log.user_id IS '用户 ID，对应 sys_user.id；用户名不存在时为 0';
COMMENT ON COLUMN sys_login_log.username IS '登录用户名';
COMMENT ON COLUMN sys_login_log.status IS '登录状态：1 成功，2 失败';
COMMENT ON COLUMN sys_login_log.message IS '登录结果说明';
COMMENT ON COLUMN sys_login_log.ip IS '客户端 IP';
COMMENT ON COLUMN sys_login_log.user_agent IS '浏览器或客户端标识';
COMMENT ON COLUMN sys_login_log.created_at IS '创建时间';

-- 公告表
CREATE TABLE sys_notice (
  id BIGSERIAL PRIMARY KEY,
  title VARCHAR(128) NOT NULL,
  content TEXT NOT NULL,
  sort INTEGER NOT NULL DEFAULT 0,
  status SMALLINT NOT NULL DEFAULT 1,
  remark VARCHAR(255) NOT NULL DEFAULT '',
  created_at TIMESTAMPTZ NOT NULL,
  updated_at TIMESTAMPTZ NOT NULL,
  deleted_at TIMESTAMPTZ NULL
);

CREATE INDEX idx_sys_notice_deleted_at ON sys_notice (deleted_at);

COMMENT ON TABLE sys_notice IS '公告表';
COMMENT ON COLUMN sys_notice.id IS '公告记录主键，数据库自增生成';
COMMENT ON COLUMN sys_notice.title IS '公告标题';
COMMENT ON COLUMN sys_notice.content IS '公告内容';
COMMENT ON COLUMN sys_notice.sort IS '排序值，数字越小越靠前';
COMMENT ON COLUMN sys_notice.status IS '公告状态：1 启用，2 禁用';
COMMENT ON COLUMN sys_notice.remark IS '备注';
COMMENT ON COLUMN sys_notice.created_at IS '创建时间';
COMMENT ON COLUMN sys_notice.updated_at IS '更新时间';
COMMENT ON COLUMN sys_notice.deleted_at IS '逻辑删除时间，NULL 表示未删除';

-- Casbin 权限策略表
CREATE TABLE casbin_rule (
  id BIGSERIAL PRIMARY KEY,
  ptype VARCHAR(100) NOT NULL DEFAULT '',
  v0 VARCHAR(100) NOT NULL DEFAULT '',
  v1 VARCHAR(100) NOT NULL DEFAULT '',
  v2 VARCHAR(100) NOT NULL DEFAULT '',
  v3 VARCHAR(100) NOT NULL DEFAULT '',
  v4 VARCHAR(100) NOT NULL DEFAULT '',
  v5 VARCHAR(100) NOT NULL DEFAULT ''
);

CREATE UNIQUE INDEX uk_casbin_rule_policy ON casbin_rule (ptype, v0, v1, v2, v3, v4, v5);
CREATE INDEX idx_casbin_rule_ptype ON casbin_rule (ptype);
CREATE INDEX idx_casbin_rule_subject ON casbin_rule (v0);

COMMENT ON TABLE casbin_rule IS 'Casbin 权限策略表';
COMMENT ON COLUMN casbin_rule.id IS '策略记录主键，数据库自增生成';
COMMENT ON COLUMN casbin_rule.ptype IS '策略类型，例如 p';
COMMENT ON COLUMN casbin_rule.v0 IS '策略主体，本项目存角色编码';
COMMENT ON COLUMN casbin_rule.v1 IS '资源路径';
COMMENT ON COLUMN casbin_rule.v2 IS '请求方法';
COMMENT ON COLUMN casbin_rule.v3 IS '预留字段';
COMMENT ON COLUMN casbin_rule.v4 IS '预留字段';
COMMENT ON COLUMN casbin_rule.v5 IS '预留字段';

-- 迁移版本追踪表（golang-migrate 自动管理）
CREATE TABLE IF NOT EXISTS schema_migrations (
 version bigint NOT NULL PRIMARY KEY,
 dirty boolean NOT NULL
);
