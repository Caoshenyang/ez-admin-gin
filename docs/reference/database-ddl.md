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
