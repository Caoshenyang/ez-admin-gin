---
title: 逻辑删除与唯一索引冲突
description: "系统解释逻辑删除字段与唯一索引之间的冲突，并对比 PostgreSQL、MySQL、Java 项目的常见处理方式。"
---

# 逻辑删除与唯一索引冲突

设计后台系统表时，经常会同时遇到两个需求：数据不物理删除、某些字段又必须唯一。比如用户表既希望支持逻辑删除，又希望 `username` 唯一。这个组合看起来简单，但如果处理不好，会在“删除后能不能重新创建同名数据”上踩坑。

::: tip 这页先看结论
默认情况下，强身份字段删除后不允许复用，继续使用普通唯一索引；如果明确需要“删除后允许重新创建同名数据”，本项目优先使用兼容 MySQL 的 `delete_marker` 方案。PostgreSQL 部分唯一索引只作为 PostgreSQL 专用方案理解。
:::

## 1. 唯一索引和逻辑删除为什么会冲突

逻辑删除不会真正删除数据，只是给记录打上删除标记。

例如用户表里有唯一字段 `username`：

```text
id | username | deleted_at
1  | admin    | 2026-04-22 10:00:00
```

从业务角度看，这条 `admin` 已经被删除了；但从数据库角度看，这行数据仍然存在。普通唯一索引仍然会检查它，所以再次创建 `username = admin` 时，仍然可能发生唯一冲突。

也就是说，问题不在逻辑删除本身，而在于：**普通唯一索引并不知道你想忽略已删除数据**。

## 2. 默认应该怎么解决

后台底座默认采用这条规则：

```text
强身份字段：删除后不允许复用，使用普通唯一索引。
明确需要复用：再做额外索引设计。
```

适合“不允许删除后复用”的字段包括：

- `username`
- `email`
- `phone`
- 角色编码
- 权限标识

原因是这些字段通常会进入登录日志、操作日志、权限变更记录。如果一个历史用户叫 `admin`，删除后又创建一个新的 `admin`，后续排查日志时会更容易混淆。

所以对于 `sys_user.username`，普通唯一索引反而是更稳的默认策略。

## 3. 为什么推荐 `deleted_at`，而不是删除状态字段

Java 项目里常见：

```text
deleted = 0 / 1
is_deleted = 0 / 1
del_flag = 0 / 1
```

这些字段都能表达“是否删除”。但在当前项目里，更推荐使用：

```go
DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
```

它不是为了显得复杂，而是多解决了几个问题。

| 角度 | `deleted_at` 的好处 |
| --- | --- |
| 删除语义 | 不只知道是否删除，还知道什么时候删除 |
| 状态拆分 | `status` 表示启用、禁用；`deleted_at` 表示生命周期删除 |
| GORM 支持 | `gorm.DeletedAt` 会自动参与默认查询和删除 |
| 审计排查 | 删除时间本身就是排查线索 |
| 后续扩展 | 后面加 `deleted_by` 时语义自然 |
| PostgreSQL 索引 | 可以直接配合 `WHERE deleted_at IS NULL` |

### 和 `status` 的边界

`status` 和 `deleted_at` 不解决同一件事。

```text
status：这条数据当前是否启用、禁用、冻结
deleted_at：这条数据是否已经从业务视图中删除
```

例如一个用户可以是：

```text
status = disabled
deleted_at = NULL
```

这表示账号被禁用，但用户仍然存在。

如果是：

```text
deleted_at = 2026-04-22 10:00:00
```

这表示这条用户记录已经被逻辑删除，正常业务查询不应该再看到它。

## 4. `deleted_at` 会不会影响查询效率

会影响查询计划，但通常不是负面影响，关键看索引和查询条件是否稳定。

正常业务查询应该带上：

```sql
SELECT *
FROM sys_user
WHERE deleted_at IS NULL;
```

如果这个条件很常用，可以建立索引：

```sql
CREATE INDEX idx_sys_user_deleted_at
ON sys_user (deleted_at);
```

如果查询经常按用户名查未删除用户，则可以进一步考虑：

```sql
CREATE INDEX idx_sys_user_username_deleted_at
ON sys_user (username, deleted_at);
```

::: info 查询效率的关键
`deleted_at IS NULL` 本身不是硬伤。真正重要的是：业务查询要稳定过滤未删除数据，索引要服务真实查询条件。
:::

在 GORM 中，使用 `gorm.DeletedAt` 后，普通查询默认会自动带上未删除条件；这比每次手写 `deleted = 0` 更不容易漏。

## 5. PostgreSQL 可以怎么解决复用问题

如果项目明确使用 PostgreSQL，并且业务要求“逻辑删除后允许重新创建同名记录”，可以使用部分唯一索引。

```sql
CREATE UNIQUE INDEX idx_sys_user_username_alive
ON sys_user (username)
WHERE deleted_at IS NULL;
```

含义是：

```text
只约束未删除记录的 username 唯一。
已逻辑删除记录不参与唯一判断。
```

它能表达这种规则：

```text
admin + deleted_at IS NULL      只能有一条
admin + deleted_at 有删除时间   可以保留多条历史记录
```

这就是 PostgreSQL 在这个问题上的优势：唯一约束可以只作用在满足条件的数据上。不过本项目需要兼容 MySQL，所以它不是默认落地方案。

## 6. PostgreSQL 部分索引会不会影响效率

部分索引通常是正向影响，前提是查询条件能命中它。

比如这个部分唯一索引：

```sql
CREATE UNIQUE INDEX idx_sys_user_username_alive
ON sys_user (username)
WHERE deleted_at IS NULL;
```

它只索引未删除用户。好处是：

- 索引更小。
- 写入和维护成本更低。
- 查询未删除数据更精准。
- 唯一约束只作用于有效数据。

典型查询能用上它：

```sql
SELECT *
FROM sys_user
WHERE username = 'admin'
  AND deleted_at IS NULL;
```

但如果查询没有带上 `deleted_at IS NULL`，就不一定能用上：

```sql
SELECT *
FROM sys_user
WHERE username = 'admin';
```

::: warning 注意查询条件
部分索引不是“建了就一定用”。查询条件需要包含或能推出 `deleted_at IS NULL`，PostgreSQL 才更容易使用它。
:::

## 7. PostgreSQL 部分索引和 MySQL 联合索引一样吗

不一样。

```text
PostgreSQL 部分索引：只给满足 WHERE 条件的行建索引。
MySQL 联合索引：把多个字段组合起来一起判断唯一。
```

PostgreSQL 部分索引写的是：

```sql
CREATE UNIQUE INDEX idx_sys_user_username_alive
ON sys_user (username)
WHERE deleted_at IS NULL;
```

它的重点是 `WHERE deleted_at IS NULL`：只索引未删除记录。

MySQL 联合索引通常是：

```sql
UNIQUE KEY uk_user_username_deleted (username, deleted_at)
```

它的重点是 `(username, deleted_at)` 两个值组合起来唯一。

这两个思路完全不同，不要混在一起理解。

## 8. 为什么 MySQL 不能直接用 `(username, deleted_at)`

MySQL 没有 PostgreSQL 这种原生部分索引。

很多人第一反应会写：

```sql
UNIQUE KEY uk_user_username_deleted_at (username, deleted_at)
```

但这个写法不稳，因为 MySQL 唯一索引允许多个 `NULL`。如果未删除记录的 `deleted_at` 都是 `NULL`，它可能挡不住多个未删除的同名记录。

所以对于 MySQL，不建议直接用 `(username, deleted_at)` 来模拟 PostgreSQL 部分唯一索引。

## 9. 如果要兼容 MySQL，怎么办

先定默认策略：

```text
默认不允许唯一字段在逻辑删除后复用。
只有明确需要复用时，才做额外设计。
```

如果确实要跨 PostgreSQL / MySQL 支持“删除后允许复用唯一值”，优先额外增加一个删除标记字段参与联合唯一索引。

约定：

```text
未删除：delete_marker = 0
已删除：delete_marker = 当前记录 id 或删除时间戳
唯一索引：username + delete_marker
```

示例：

```go
Username     string         `gorm:"size:64;not null;index:uk_user_username_marker,unique" json:"username"`
DeleteMarker uint           `gorm:"not null;default:0;index:uk_user_username_marker,unique" json:"-"`
DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
```

数据会变成：

```text
admin, delete_marker = 0   当前有效记录
admin, delete_marker = 1   已删除历史记录
admin, delete_marker = 8   已删除历史记录
```

这样只有 `admin + 0` 会被唯一约束，已删除历史记录之间不会互相冲突，也不会阻止重新创建新的 `admin`。

::: details 为什么不只用 `deleted = 0 / 1`
如果唯一索引是 `(username, deleted)`，第一次删除后会有：

```text
admin, deleted = 1
```

再次创建并删除 `admin` 时，又会出现另一条：

```text
admin, deleted = 1
```

这两条历史记录会互相冲突。所以允许复用时，删除标记不能只有 `0 / 1`，已删除记录需要使用不同标记值。
:::

## 10. Java 项目通常怎么处理

Java 项目中，逻辑删除通常由框架帮忙自动补查询条件。

常见方式：

| 技术栈 | 常见写法 |
| --- | --- |
| MyBatis-Plus | `@TableLogic` + `deleted` / `del_flag` |
| JPA / Hibernate | `@SQLDelete` + `@Where` |

这些框架主要解决：

```text
查询时默认过滤已删除数据。
删除时改成 update 标记删除。
```

但它们不能自动解决唯一索引冲突。唯一字段是否允许删除后复用，最终仍然要靠数据库索引设计决定。

Java 项目常见策略也是这几类：

| 场景 | 常见处理 |
| --- | --- |
| 强身份字段 | 普通唯一索引，删除后不允许复用 |
| 允许复用字段 | `deleted = 0 / id / 时间戳` + 联合唯一索引 |
| PostgreSQL 项目 | 使用部分唯一索引 |

所以这个问题本质不是语言差异，而是数据库索引策略差异。Java 框架负责逻辑删除行为，数据库负责唯一约束。

## 11. 当前项目怎么落地

当前后台底座需要兼容 MySQL，默认采用：

```text
逻辑删除字段：deleted_at
业务状态字段：status
强身份唯一字段：删除后不允许复用
默认索引：普通唯一索引
PostgreSQL 专用复用方案：部分唯一索引
跨库复用方案：delete_marker + 联合唯一索引
```

用于 `sys_user.username` 时，默认建议：

```go
Username  string         `gorm:"size:64;not null;uniqueIndex" json:"username"`
DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
```

含义是：

```text
用户名一旦使用过，即使用户被逻辑删除，也不允许重新创建同名用户。
```

这能让登录日志、操作日志和权限变更记录更清晰。

## 12. 选型速查

| 需求 | 推荐方案 |
| --- | --- |
| 用户名、手机号、邮箱等强身份字段 | 普通唯一索引，删除后不复用 |
| 需要兼容 MySQL，并允许删除后复用 | `delete_marker` + 联合唯一索引 |
| 只绑定 PostgreSQL，并允许删除后复用 | `WHERE deleted_at IS NULL` 部分唯一索引 |
| 只想表达是否删除 | 仍推荐 `deleted_at`，因为能保留删除时间 |
| 操作日志、登录日志这类事实记录 | 通常不做逻辑删除 |

如果拿不准，就先使用默认策略：**强身份唯一字段删除后不复用，保留普通唯一索引**。
