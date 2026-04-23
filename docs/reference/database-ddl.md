---
title: 数据库建表语句
description: "集中记录后台底座中系统表的 PostgreSQL 与 MySQL 建表语句。"
---

# 数据库建表语句

这里集中记录教程中出现过的系统表建表语句。每张表都会给出 PostgreSQL 和 MySQL 两个版本，并补充表注释、字段注释和索引。

::: tip 怎么使用这页
教程中的表结构以这里的 SQL 为准。新增表时，先在这里确认 PostgreSQL / MySQL 两个版本，再到对应数据库中执行。
:::

::: warning 时间字段由代码维护
`created_at` 和 `updated_at` 不依赖数据库默认函数或触发器，统一由应用代码维护。使用 GORM 创建和更新数据时，`CreatedAt` / `UpdatedAt` 会自动写入。

如果绕过 GORM 直接执行 SQL，就需要在 SQL 中显式写入这两个时间字段。
:::

::: tip 主键生成策略
系统表默认使用数据库自增 BIGINT 主键。PostgreSQL 使用 `BIGSERIAL`，MySQL 使用 `BIGINT UNSIGNED AUTO_INCREMENT`。

应用代码创建数据时不手动生成 `id`，由数据库生成后回填。业务可读标识单独使用 `username`、`code` 这类字段。
:::

## 当前表清单

| 表名 | 说明 |
| --- | --- |
| `sys_user` | 后台用户表 |
| `sys_role` | 后台角色表 |
| `sys_user_role` | 用户角色关系表 |
| `sys_menu` | 后台菜单和按钮表 |
| `sys_role_menu` | 角色菜单关系表 |
| `sys_config` | 系统配置表 |
| `sys_file` | 文件上传记录表 |
| `sys_operation_log` | 操作日志表 |
| `casbin_rule` | Casbin 权限策略表 |

<a id="sys-user"></a>

## `sys_user` 后台用户表

`sys_user` 保存后台登录用户。用户名默认不允许在逻辑删除后复用，所以 `username` 使用普通唯一索引。

字段含义：

| 字段 | 说明 |
| --- | --- |
| `id` | 用户记录主键，数据库自增生成 |
| `username` | 登录用户名 |
| `password_hash` | 密码哈希 |
| `nickname` | 管理台展示名称 |
| `status` | 用户状态：`1` 启用，`2` 禁用 |
| `created_at` | 创建时间 |
| `updated_at` | 更新时间 |
| `deleted_at` | 逻辑删除时间，`NULL` 表示未删除 |

::: warning 唯一索引与逻辑删除
`username` 使用普通唯一索引后，即使一条用户记录被逻辑删除，相同用户名也不能再次创建。这样可以避免历史日志和操作者身份发生歧义。

如果后续某张表明确需要“删除后允许重新创建同名数据”，再使用 `delete_marker` + 联合唯一索引方案。更多背景可以看：[逻辑删除与唯一索引冲突](./logical-delete-and-unique-index)。
:::

### 建表语句

::: code-group

```sql [PostgreSQL]
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

CREATE UNIQUE INDEX uk_sys_user_username
ON sys_user (username);

CREATE INDEX idx_sys_user_deleted_at
ON sys_user (deleted_at);

COMMENT ON TABLE sys_user IS '后台用户表';
COMMENT ON COLUMN sys_user.id IS '用户记录主键，数据库自增生成';
COMMENT ON COLUMN sys_user.username IS '登录用户名';
COMMENT ON COLUMN sys_user.password_hash IS '密码哈希';
COMMENT ON COLUMN sys_user.nickname IS '管理台展示名称';
COMMENT ON COLUMN sys_user.status IS '用户状态：1 启用，2 禁用';
COMMENT ON COLUMN sys_user.created_at IS '创建时间';
COMMENT ON COLUMN sys_user.updated_at IS '更新时间';
COMMENT ON COLUMN sys_user.deleted_at IS '逻辑删除时间，NULL 表示未删除';
```

```sql [MySQL]
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
) ENGINE=InnoDB
  DEFAULT CHARSET=utf8mb4
  COLLATE=utf8mb4_unicode_ci
  COMMENT='后台用户表';
```

:::

::: details 为什么不在数据库里写 `CURRENT_TIMESTAMP`
本项目把创建时间和更新时间作为应用层审计字段，由 GORM 统一写入。这样 PostgreSQL 和 MySQL 的行为更一致，也不会因为数据库默认函数、触发器或 `ON UPDATE` 规则不同导致时间维护逻辑分散。

如果后续有批处理脚本或初始化 SQL，需要显式写入 `created_at` 和 `updated_at`。
:::

::: details 为什么 MySQL 这里没有外键
`sys_user` 当前没有关联字段。后续出现 `user_id`、`role_id`、`menu_id` 这类字段时，也只建立普通索引，不创建数据库级外键约束。

关联是否存在、是否允许删除、是否允许绑定，由 service 层业务逻辑维护。
:::

<a id="sys-role"></a>

## `sys_role` 后台角色表

`sys_role` 保存后台角色。角色编码默认不允许在逻辑删除后复用，所以 `code` 使用普通唯一索引。

字段含义：

| 字段 | 说明 |
| --- | --- |
| `id` | 角色记录主键，数据库自增生成 |
| `code` | 角色编码 |
| `name` | 角色名称 |
| `sort` | 排序值，数字越小越靠前 |
| `status` | 角色状态：`1` 启用，`2` 禁用 |
| `remark` | 备注 |
| `created_at` | 创建时间 |
| `updated_at` | 更新时间 |
| `deleted_at` | 逻辑删除时间，`NULL` 表示未删除 |

### 建表语句

::: code-group

```sql [PostgreSQL]
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

CREATE UNIQUE INDEX uk_sys_role_code
ON sys_role (code);

CREATE INDEX idx_sys_role_status
ON sys_role (status);

CREATE INDEX idx_sys_role_deleted_at
ON sys_role (deleted_at);

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
```

```sql [MySQL]
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
) ENGINE=InnoDB
  DEFAULT CHARSET=utf8mb4
  COLLATE=utf8mb4_unicode_ci
  COMMENT='后台角色表';
```

:::

<a id="sys-user-role"></a>

## `sys_user_role` 用户角色关系表

`sys_user_role` 保存用户和角色之间的绑定关系。这里不创建数据库级外键，只通过字段、索引和业务逻辑维护关系。

字段含义：

| 字段 | 说明 |
| --- | --- |
| `id` | 关系记录主键，数据库自增生成 |
| `user_id` | 用户 ID，对应 `sys_user.id` |
| `role_id` | 角色 ID，对应 `sys_role.id` |
| `created_at` | 绑定时间 |
| `updated_at` | 更新时间 |

### 建表语句

::: code-group

```sql [PostgreSQL]
CREATE TABLE sys_user_role (
  id BIGSERIAL PRIMARY KEY,
  user_id BIGINT NOT NULL,
  role_id BIGINT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL,
  updated_at TIMESTAMPTZ NOT NULL
);

CREATE UNIQUE INDEX uk_sys_user_role_user_role
ON sys_user_role (user_id, role_id);

CREATE INDEX idx_sys_user_role_user_id
ON sys_user_role (user_id);

CREATE INDEX idx_sys_user_role_role_id
ON sys_user_role (role_id);

COMMENT ON TABLE sys_user_role IS '用户角色关系表';
COMMENT ON COLUMN sys_user_role.id IS '关系记录主键，数据库自增生成';
COMMENT ON COLUMN sys_user_role.user_id IS '用户 ID，对应 sys_user.id';
COMMENT ON COLUMN sys_user_role.role_id IS '角色 ID，对应 sys_role.id';
COMMENT ON COLUMN sys_user_role.created_at IS '绑定时间';
COMMENT ON COLUMN sys_user_role.updated_at IS '更新时间';
```

```sql [MySQL]
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
) ENGINE=InnoDB
  DEFAULT CHARSET=utf8mb4
  COLLATE=utf8mb4_unicode_ci
  COMMENT='用户角色关系表';
```

:::

::: details 为什么关系表没有数据库外键
本项目不使用数据库级外键约束。`user_id`、`role_id` 只表达关系，是否允许绑定、解绑、删除，由 service 层逻辑判断。

这种方式更适合逻辑删除、初始化脚本、批量导入和后续模块复用。
:::

<a id="sys-menu"></a>

## `sys_menu` 后台菜单和按钮表

`sys_menu` 保存管理台的目录、菜单和按钮权限点。目录、菜单、按钮都放在同一张表中，通过 `type` 区分。

字段含义：

| 字段 | 说明 |
| --- | --- |
| `id` | 菜单记录主键，数据库自增生成 |
| `parent_id` | 父级菜单 ID，根节点为 `0` |
| `type` | 节点类型：`1` 目录，`2` 菜单，`3` 按钮 |
| `code` | 菜单或按钮编码，系统内唯一 |
| `title` | 展示名称 |
| `path` | 前端路由路径 |
| `component` | 前端组件路径 |
| `icon` | 图标标识 |
| `sort` | 排序值，数字越小越靠前 |
| `status` | 菜单状态：`1` 启用，`2` 禁用 |
| `remark` | 备注 |
| `created_at` | 创建时间 |
| `updated_at` | 更新时间 |
| `deleted_at` | 逻辑删除时间，`NULL` 表示未删除 |

### 建表语句

::: code-group

```sql [PostgreSQL]
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

CREATE UNIQUE INDEX uk_sys_menu_code
ON sys_menu (code);

CREATE INDEX idx_sys_menu_parent_id
ON sys_menu (parent_id);

CREATE INDEX idx_sys_menu_type
ON sys_menu (type);

CREATE INDEX idx_sys_menu_status
ON sys_menu (status);

CREATE INDEX idx_sys_menu_deleted_at
ON sys_menu (deleted_at);

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
```

```sql [MySQL]
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
) ENGINE=InnoDB
  DEFAULT CHARSET=utf8mb4
  COLLATE=utf8mb4_unicode_ci
  COMMENT='后台菜单和按钮表';
```

:::

::: details 为什么 `code` 使用普通唯一索引
菜单编码是稳定权限标识，例如 `system:health:view`。即使菜单被逻辑删除，也默认不复用相同编码，避免历史授权、操作日志和前端判断产生歧义。
:::

<a id="sys-role-menu"></a>

## `sys_role_menu` 角色菜单关系表

`sys_role_menu` 保存角色和菜单之间的绑定关系。一个角色可以拥有多个菜单，一个菜单也可以授权给多个角色。

字段含义：

| 字段 | 说明 |
| --- | --- |
| `id` | 关系记录主键，数据库自增生成 |
| `role_id` | 角色 ID，对应 `sys_role.id` |
| `menu_id` | 菜单 ID，对应 `sys_menu.id` |
| `created_at` | 绑定时间 |
| `updated_at` | 更新时间 |

### 建表语句

::: code-group

```sql [PostgreSQL]
CREATE TABLE sys_role_menu (
  id BIGSERIAL PRIMARY KEY,
  role_id BIGINT NOT NULL,
  menu_id BIGINT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL,
  updated_at TIMESTAMPTZ NOT NULL
);

CREATE UNIQUE INDEX uk_sys_role_menu_role_menu
ON sys_role_menu (role_id, menu_id);

CREATE INDEX idx_sys_role_menu_role_id
ON sys_role_menu (role_id);

CREATE INDEX idx_sys_role_menu_menu_id
ON sys_role_menu (menu_id);

COMMENT ON TABLE sys_role_menu IS '角色菜单关系表';
COMMENT ON COLUMN sys_role_menu.id IS '关系记录主键，数据库自增生成';
COMMENT ON COLUMN sys_role_menu.role_id IS '角色 ID，对应 sys_role.id';
COMMENT ON COLUMN sys_role_menu.menu_id IS '菜单 ID，对应 sys_menu.id';
COMMENT ON COLUMN sys_role_menu.created_at IS '绑定时间';
COMMENT ON COLUMN sys_role_menu.updated_at IS '更新时间';
```

```sql [MySQL]
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
) ENGINE=InnoDB
  DEFAULT CHARSET=utf8mb4
  COLLATE=utf8mb4_unicode_ci
  COMMENT='角色菜单关系表';
```

:::

::: details 为什么关系表不建数据库外键
本项目通过普通索引和业务逻辑维护关系，不使用数据库级外键约束。绑定菜单前应由业务逻辑确认角色和菜单存在，删除角色或菜单前也应检查依赖关系。
:::

<a id="sys-file"></a>

## `sys_file` 文件上传记录表

`sys_file` 保存文件上传后的元数据。文件内容保存在本地目录或后续扩展的存储服务中，数据库只保存文件名、路径、访问地址、大小、哈希等信息。

字段含义：

| 字段 | 说明 |
| --- | --- |
| `id` | 文件记录主键，数据库自增生成 |
| `storage` | 存储类型，本节使用 `local` |
| `original_name` | 用户上传时的原始文件名 |
| `file_name` | 后端生成的保存文件名 |
| `ext` | 文件后缀，例如 `.png`、`.pdf` |
| `mime_type` | 上传请求中的文件 MIME 类型 |
| `size` | 文件大小，单位字节 |
| `sha256` | 文件内容 SHA-256 哈希 |
| `path` | 服务端保存路径 |
| `url` | 前端可访问地址 |
| `uploader_id` | 上传用户 ID，对应 `sys_user.id` |
| `status` | 文件状态：`1` 启用，`2` 停用 |
| `remark` | 备注 |
| `created_at` | 创建时间 |
| `updated_at` | 更新时间 |
| `deleted_at` | 逻辑删除时间，`NULL` 表示未删除 |

### 建表语句

::: code-group

```sql [PostgreSQL]
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

CREATE INDEX idx_sys_file_ext
ON sys_file (ext);

CREATE INDEX idx_sys_file_sha256
ON sys_file (sha256);

CREATE INDEX idx_sys_file_uploader_id
ON sys_file (uploader_id);

CREATE INDEX idx_sys_file_status
ON sys_file (status);

CREATE INDEX idx_sys_file_deleted_at
ON sys_file (deleted_at);

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
```

```sql [MySQL]
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
) ENGINE=InnoDB
  DEFAULT CHARSET=utf8mb4
  COLLATE=utf8mb4_unicode_ci
  COMMENT='文件上传记录表';
```

:::

::: details 为什么 `sys_file` 不保存文件内容
数据库保存文件元数据，文件内容保存在磁盘或对象存储中。这样数据库体积更可控，文件访问和迁移也更灵活。
:::

<a id="sys-operation-log"></a>

## `sys_operation_log` 操作日志表

`sys_operation_log` 保存后台用户的关键写操作。它是审计事实记录，不做逻辑删除，也不保存完整请求体。

字段含义：

| 字段 | 说明 |
| --- | --- |
| `id` | 操作日志主键，数据库自增生成 |
| `user_id` | 操作人 ID，对应 `sys_user.id` |
| `username` | 操作人用户名 |
| `method` | HTTP 请求方法 |
| `path` | 实际请求路径 |
| `route_path` | Gin 路由模板 |
| `query` | 查询参数 |
| `ip` | 客户端 IP |
| `user_agent` | 浏览器或客户端标识 |
| `status_code` | HTTP 状态码 |
| `latency_ms` | 请求耗时，单位毫秒 |
| `success` | 是否成功 |
| `error_message` | 错误摘要 |
| `created_at` | 创建时间 |

### 建表语句

::: code-group

```sql [PostgreSQL]
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

CREATE INDEX idx_sys_operation_log_user_id
ON sys_operation_log (user_id);

CREATE INDEX idx_sys_operation_log_username
ON sys_operation_log (username);

CREATE INDEX idx_sys_operation_log_method
ON sys_operation_log (method);

CREATE INDEX idx_sys_operation_log_path
ON sys_operation_log (path);

CREATE INDEX idx_sys_operation_log_route_path
ON sys_operation_log (route_path);

CREATE INDEX idx_sys_operation_log_status_code
ON sys_operation_log (status_code);

CREATE INDEX idx_sys_operation_log_success
ON sys_operation_log (success);

CREATE INDEX idx_sys_operation_log_created_at
ON sys_operation_log (created_at);

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
```

```sql [MySQL]
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
) ENGINE=InnoDB
  DEFAULT CHARSET=utf8mb4
  COLLATE=utf8mb4_unicode_ci
  COMMENT='操作日志表';
```

:::

::: details 为什么操作日志没有 `updated_at` 和 `deleted_at`
操作日志记录的是已经发生的事实，正常情况下不编辑、不逻辑删除。后续如果数据量变大，可以按时间做归档或清理策略。
:::

<a id="sys-config"></a>

## `sys_config` 系统配置表

`sys_config` 保存后台可维护的普通业务配置，例如站点标题、上传目录、默认分页大小等。这里不存数据库密码、JWT 密钥这类敏感信息。

字段含义：

| 字段 | 说明 |
| --- | --- |
| `id` | 配置记录主键，数据库自增生成 |
| `group_code` | 配置分组，例如 `site`、`upload` |
| `config_key` | 配置键，系统内唯一，例如 `site:title` |
| `name` | 配置名称 |
| `value` | 配置值，统一按字符串存储 |
| `sort` | 排序值，数字越小越靠前 |
| `status` | 配置状态：`1` 启用，`2` 禁用 |
| `remark` | 备注 |
| `created_at` | 创建时间 |
| `updated_at` | 更新时间 |
| `deleted_at` | 逻辑删除时间，`NULL` 表示未删除 |

### 建表语句

::: code-group

```sql [PostgreSQL]
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

CREATE UNIQUE INDEX uk_sys_config_key
ON sys_config (config_key);

CREATE INDEX idx_sys_config_group_code
ON sys_config (group_code);

CREATE INDEX idx_sys_config_status
ON sys_config (status);

CREATE INDEX idx_sys_config_deleted_at
ON sys_config (deleted_at);

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
```

```sql [MySQL]
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
) ENGINE=InnoDB
  DEFAULT CHARSET=utf8mb4
  COLLATE=utf8mb4_unicode_ci
  COMMENT='系统配置表';
```

:::

::: warning 系统配置表不存敏感密钥
`sys_config` 适合放业务可调参数，不适合放数据库密码、JWT 密钥、第三方平台 Secret 这类敏感内容。敏感配置仍然应通过环境变量或配置文件维护。
:::

::: details 为什么 `value` 使用 `TEXT`
配置值统一按字符串存储时，最容易遇到的情况就是“短值为主，但偶尔会放一段 JSON 字符串或较长文本”。`TEXT` 更省心，也避免后续因为长度限制再改表。
:::

<a id="casbin-rule"></a>

## `casbin_rule` Casbin 权限策略表

`casbin_rule` 保存 Casbin 权限策略。本项目使用角色编码作为策略主体，默认策略形如：

```text
p, super_admin, /api/v1/system/health, GET
```

字段含义：

| 字段 | 说明 |
| --- | --- |
| `id` | 策略记录主键，数据库自增生成 |
| `ptype` | 策略类型，例如 `p` |
| `v0` | 策略主体，本项目存角色编码 |
| `v1` | 资源路径 |
| `v2` | 请求方法 |
| `v3` | 预留字段 |
| `v4` | 预留字段 |
| `v5` | 预留字段 |

### 建表语句

::: code-group

```sql [PostgreSQL]
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

CREATE UNIQUE INDEX uk_casbin_rule_policy
ON casbin_rule (ptype, v0, v1, v2, v3, v4, v5);

CREATE INDEX idx_casbin_rule_ptype
ON casbin_rule (ptype);

CREATE INDEX idx_casbin_rule_subject
ON casbin_rule (v0);

COMMENT ON TABLE casbin_rule IS 'Casbin 权限策略表';
COMMENT ON COLUMN casbin_rule.id IS '策略记录主键，数据库自增生成';
COMMENT ON COLUMN casbin_rule.ptype IS '策略类型，例如 p';
COMMENT ON COLUMN casbin_rule.v0 IS '策略主体，本项目存角色编码';
COMMENT ON COLUMN casbin_rule.v1 IS '资源路径';
COMMENT ON COLUMN casbin_rule.v2 IS '请求方法';
COMMENT ON COLUMN casbin_rule.v3 IS '预留字段';
COMMENT ON COLUMN casbin_rule.v4 IS '预留字段';
COMMENT ON COLUMN casbin_rule.v5 IS '预留字段';
```

```sql [MySQL]
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
) ENGINE=InnoDB
  DEFAULT CHARSET=utf8mb4
  COLLATE=utf8mb4_unicode_ci
  COMMENT='Casbin 权限策略表';
```

:::

::: details 为什么 `casbin_rule` 没有时间字段
`casbin_rule` 是 Casbin GORM 适配器使用的策略表，不是普通业务实体。这里优先保持适配器表结构简单稳定。

后续如果需要审计权限变更，可以单独设计权限操作日志表。
:::
