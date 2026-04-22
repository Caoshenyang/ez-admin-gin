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

## 当前表清单

| 表名 | 说明 |
| --- | --- |
| `sys_user` | 后台用户表 |

<a id="sys-user"></a>

## `sys_user` 后台用户表

`sys_user` 保存后台登录用户。用户名默认不允许在逻辑删除后复用，所以 `username` 使用普通唯一索引。

字段含义：

| 字段 | 说明 |
| --- | --- |
| `id` | 用户记录主键 |
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
COMMENT ON COLUMN sys_user.id IS '用户记录主键';
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
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '用户记录主键',
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
