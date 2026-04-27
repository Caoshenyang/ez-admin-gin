-- 系统用户表
CREATE TABLE `sys_user` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '用户记录主键，数据库自增生成',
  `username` VARCHAR(64) NOT NULL COMMENT '登录用户名',
  `password_hash` VARCHAR(255) NOT NULL COMMENT '密码哈希',
  `nickname` VARCHAR(64) NOT NULL DEFAULT '' COMMENT '管理台展示名称',
  `status` SMALLINT NOT NULL DEFAULT 1 COMMENT '用户状态：1 启用，2 禁用',
  `created_at` DATETIME(3) NOT NULL COMMENT '创建时间',
  `updated_at` DATETIME(3) NOT NULL COMMENT '更新时间',
  `deleted_at` DATETIME(3) NULL DEFAULT NULL COMMENT '逻辑删除时间，NULL 表示未删除',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_sys_user_username` (`username`),
  KEY `idx_sys_user_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='后台用户表';

-- 系统角色表
CREATE TABLE `sys_role` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '角色记录主键，数据库自增生成',
  `code` VARCHAR(64) NOT NULL COMMENT '角色编码',
  `name` VARCHAR(64) NOT NULL COMMENT '角色名称',
  `sort` INT NOT NULL DEFAULT 0 COMMENT '排序值，数字越小越靠前',
  `status` SMALLINT NOT NULL DEFAULT 1 COMMENT '角色状态：1 启用，2 禁用',
  `remark` VARCHAR(255) NOT NULL DEFAULT '' COMMENT '备注',
  `created_at` DATETIME(3) NOT NULL COMMENT '创建时间',
  `updated_at` DATETIME(3) NOT NULL COMMENT '更新时间',
  `deleted_at` DATETIME(3) NULL DEFAULT NULL COMMENT '逻辑删除时间，NULL 表示未删除',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_sys_role_code` (`code`),
  KEY `idx_sys_role_status` (`status`),
  KEY `idx_sys_role_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='后台角色表';

-- 用户角色关系表
CREATE TABLE `sys_user_role` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '关系记录主键，数据库自增生成',
  `user_id` BIGINT UNSIGNED NOT NULL COMMENT '用户 ID，对应 sys_user.id',
  `role_id` BIGINT UNSIGNED NOT NULL COMMENT '角色 ID，对应 sys_role.id',
  `created_at` DATETIME(3) NOT NULL COMMENT '绑定时间',
  `updated_at` DATETIME(3) NOT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_sys_user_role_user_role` (`user_id`, `role_id`),
  KEY `idx_sys_user_role_user_id` (`user_id`),
  KEY `idx_sys_user_role_role_id` (`role_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户角色关系表';

-- 系统菜单和按钮表
CREATE TABLE `sys_menu` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '菜单记录主键，数据库自增生成',
  `parent_id` BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '父级菜单 ID，根节点为 0',
  `type` SMALLINT NOT NULL COMMENT '节点类型：1 目录，2 菜单，3 按钮',
  `code` VARCHAR(128) NOT NULL COMMENT '菜单或按钮编码，系统内唯一',
  `title` VARCHAR(64) NOT NULL COMMENT '展示名称',
  `path` VARCHAR(255) NOT NULL DEFAULT '' COMMENT '前端路由路径',
  `component` VARCHAR(255) NOT NULL DEFAULT '' COMMENT '前端组件路径',
  `icon` VARCHAR(64) NOT NULL DEFAULT '' COMMENT '图标标识',
  `sort` INT NOT NULL DEFAULT 0 COMMENT '排序值，数字越小越靠前',
  `status` SMALLINT NOT NULL DEFAULT 1 COMMENT '菜单状态：1 启用，2 禁用',
  `remark` VARCHAR(255) NOT NULL DEFAULT '' COMMENT '备注',
  `created_at` DATETIME(3) NOT NULL COMMENT '创建时间',
  `updated_at` DATETIME(3) NOT NULL COMMENT '更新时间',
  `deleted_at` DATETIME(3) NULL DEFAULT NULL COMMENT '逻辑删除时间，NULL 表示未删除',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_sys_menu_code` (`code`),
  KEY `idx_sys_menu_parent_id` (`parent_id`),
  KEY `idx_sys_menu_type` (`type`),
  KEY `idx_sys_menu_status` (`status`),
  KEY `idx_sys_menu_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='后台菜单和按钮表';

-- 角色菜单关系表
CREATE TABLE `sys_role_menu` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '关系记录主键，数据库自增生成',
  `role_id` BIGINT UNSIGNED NOT NULL COMMENT '角色 ID，对应 sys_role.id',
  `menu_id` BIGINT UNSIGNED NOT NULL COMMENT '菜单 ID，对应 sys_menu.id',
  `created_at` DATETIME(3) NOT NULL COMMENT '绑定时间',
  `updated_at` DATETIME(3) NOT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_sys_role_menu_role_menu` (`role_id`, `menu_id`),
  KEY `idx_sys_role_menu_role_id` (`role_id`),
  KEY `idx_sys_role_menu_menu_id` (`menu_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='角色菜单关系表';

-- 系统配置表
CREATE TABLE `sys_config` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '配置记录主键，数据库自增生成',
  `group_code` VARCHAR(64) NOT NULL COMMENT '配置分组，例如 site、upload',
  `config_key` VARCHAR(128) NOT NULL COMMENT '配置键，系统内唯一，例如 site:title',
  `name` VARCHAR(64) NOT NULL COMMENT '配置名称',
  `value` TEXT NOT NULL COMMENT '配置值，统一按字符串存储',
  `sort` INT NOT NULL DEFAULT 0 COMMENT '排序值，数字越小越靠前',
  `status` SMALLINT NOT NULL DEFAULT 1 COMMENT '配置状态：1 启用，2 禁用',
  `remark` VARCHAR(255) NOT NULL DEFAULT '' COMMENT '备注',
  `created_at` DATETIME(3) NOT NULL COMMENT '创建时间',
  `updated_at` DATETIME(3) NOT NULL COMMENT '更新时间',
  `deleted_at` DATETIME(3) NULL DEFAULT NULL COMMENT '逻辑删除时间，NULL 表示未删除',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_sys_config_key` (`config_key`),
  KEY `idx_sys_config_group_code` (`group_code`),
  KEY `idx_sys_config_status` (`status`),
  KEY `idx_sys_config_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='系统配置表';

-- 文件上传记录表
CREATE TABLE `sys_file` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '文件记录主键，数据库自增生成',
  `storage` VARCHAR(32) NOT NULL DEFAULT 'local' COMMENT '存储类型，本节使用 local',
  `original_name` VARCHAR(255) NOT NULL COMMENT '用户上传时的原始文件名',
  `file_name` VARCHAR(255) NOT NULL COMMENT '后端生成的保存文件名',
  `ext` VARCHAR(32) NOT NULL DEFAULT '' COMMENT '文件后缀，例如 .png、.pdf',
  `mime_type` VARCHAR(128) NOT NULL DEFAULT '' COMMENT '上传请求中的文件 MIME 类型',
  `size` BIGINT NOT NULL DEFAULT 0 COMMENT '文件大小，单位字节',
  `sha256` VARCHAR(64) NOT NULL DEFAULT '' COMMENT '文件内容 SHA-256 哈希',
  `path` VARCHAR(500) NOT NULL COMMENT '服务端保存路径',
  `url` VARCHAR(500) NOT NULL COMMENT '前端可访问地址',
  `uploader_id` BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '上传用户 ID，对应 sys_user.id',
  `status` SMALLINT NOT NULL DEFAULT 1 COMMENT '文件状态：1 启用，2 停用',
  `remark` VARCHAR(255) NOT NULL DEFAULT '' COMMENT '备注',
  `created_at` DATETIME(3) NOT NULL COMMENT '创建时间',
  `updated_at` DATETIME(3) NOT NULL COMMENT '更新时间',
  `deleted_at` DATETIME(3) NULL DEFAULT NULL COMMENT '逻辑删除时间，NULL 表示未删除',
  PRIMARY KEY (`id`),
  KEY `idx_sys_file_ext` (`ext`),
  KEY `idx_sys_file_sha256` (`sha256`),
  KEY `idx_sys_file_uploader_id` (`uploader_id`),
  KEY `idx_sys_file_status` (`status`),
  KEY `idx_sys_file_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='文件上传记录表';

-- 操作日志表
CREATE TABLE `sys_operation_log` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '操作日志主键，数据库自增生成',
  `user_id` BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '操作人 ID，对应 sys_user.id',
  `username` VARCHAR(64) NOT NULL DEFAULT '' COMMENT '操作人用户名',
  `method` VARCHAR(10) NOT NULL COMMENT 'HTTP 请求方法',
  `path` VARCHAR(255) NOT NULL COMMENT '实际请求路径',
  `route_path` VARCHAR(255) NOT NULL DEFAULT '' COMMENT 'Gin 路由模板',
  `query` VARCHAR(1000) NOT NULL DEFAULT '' COMMENT '查询参数',
  `ip` VARCHAR(64) NOT NULL DEFAULT '' COMMENT '客户端 IP',
  `user_agent` VARCHAR(500) NOT NULL DEFAULT '' COMMENT '浏览器或客户端标识',
  `status_code` INT NOT NULL DEFAULT 0 COMMENT 'HTTP 状态码',
  `latency_ms` BIGINT NOT NULL DEFAULT 0 COMMENT '请求耗时，单位毫秒',
  `success` BOOLEAN NOT NULL DEFAULT TRUE COMMENT '是否成功',
  `error_message` VARCHAR(500) NOT NULL DEFAULT '' COMMENT '错误摘要',
  `created_at` DATETIME(3) NOT NULL COMMENT '创建时间',
  PRIMARY KEY (`id`),
  KEY `idx_sys_operation_log_user_id` (`user_id`),
  KEY `idx_sys_operation_log_username` (`username`),
  KEY `idx_sys_operation_log_method` (`method`),
  KEY `idx_sys_operation_log_path` (`path`),
  KEY `idx_sys_operation_log_route_path` (`route_path`),
  KEY `idx_sys_operation_log_status_code` (`status_code`),
  KEY `idx_sys_operation_log_success` (`success`),
  KEY `idx_sys_operation_log_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='操作日志表';

-- 登录日志表
CREATE TABLE `sys_login_log` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '登录日志主键，数据库自增生成',
  `user_id` BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '用户 ID，对应 sys_user.id；用户名不存在时为 0',
  `username` VARCHAR(64) NOT NULL DEFAULT '' COMMENT '登录用户名',
  `status` SMALLINT NOT NULL COMMENT '登录状态：1 成功，2 失败',
  `message` VARCHAR(255) NOT NULL DEFAULT '' COMMENT '登录结果说明',
  `ip` VARCHAR(64) NOT NULL DEFAULT '' COMMENT '客户端 IP',
  `user_agent` VARCHAR(500) NOT NULL DEFAULT '' COMMENT '浏览器或客户端标识',
  `created_at` DATETIME(3) NOT NULL COMMENT '创建时间',
  PRIMARY KEY (`id`),
  KEY `idx_sys_login_log_user_id` (`user_id`),
  KEY `idx_sys_login_log_username` (`username`),
  KEY `idx_sys_login_log_status` (`status`),
  KEY `idx_sys_login_log_ip` (`ip`),
  KEY `idx_sys_login_log_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='登录日志表';

-- 公告表
CREATE TABLE `sys_notice` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '公告记录主键，数据库自增生成',
  `title` VARCHAR(128) NOT NULL COMMENT '公告标题',
  `content` TEXT NOT NULL COMMENT '公告内容',
  `sort` INT NOT NULL DEFAULT 0 COMMENT '排序值，数字越小越靠前',
  `status` SMALLINT NOT NULL DEFAULT 1 COMMENT '公告状态：1 启用，2 禁用',
  `remark` VARCHAR(255) NOT NULL DEFAULT '' COMMENT '备注',
  `created_at` DATETIME(3) NOT NULL COMMENT '创建时间',
  `updated_at` DATETIME(3) NOT NULL COMMENT '更新时间',
  `deleted_at` DATETIME(3) NULL DEFAULT NULL COMMENT '逻辑删除时间，NULL 表示未删除',
  PRIMARY KEY (`id`),
  KEY `idx_sys_notice_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='公告表';

-- Casbin 权限策略表
CREATE TABLE `casbin_rule` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '策略记录主键，数据库自增生成',
  `ptype` VARCHAR(100) NOT NULL DEFAULT '' COMMENT '策略类型，例如 p',
  `v0` VARCHAR(100) NOT NULL DEFAULT '' COMMENT '策略主体，本项目存角色编码',
  `v1` VARCHAR(100) NOT NULL DEFAULT '' COMMENT '资源路径',
  `v2` VARCHAR(100) NOT NULL DEFAULT '' COMMENT '请求方法',
  `v3` VARCHAR(100) NOT NULL DEFAULT '' COMMENT '预留字段',
  `v4` VARCHAR(100) NOT NULL DEFAULT '' COMMENT '预留字段',
  `v5` VARCHAR(100) NOT NULL DEFAULT '' COMMENT '预留字段',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_casbin_rule_policy` (`ptype`, `v0`, `v1`, `v2`, `v3`, `v4`, `v5`),
  KEY `idx_casbin_rule_ptype` (`ptype`),
  KEY `idx_casbin_rule_subject` (`v0`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='Casbin 权限策略表';

-- 迁移版本追踪表（golang-migrate 自动管理）
CREATE TABLE IF NOT EXISTS `schema_migrations` (
 `version` bigint NOT NULL PRIMARY KEY,
 `dirty` boolean NOT NULL
);
