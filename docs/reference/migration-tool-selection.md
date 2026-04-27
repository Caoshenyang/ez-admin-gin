---
title: 数据库迁移工具选型
description: "对比 Goose、golang-migrate 和 Atlas 三款 Go 数据库迁移工具，为本项目选择合适的方案。"
---

# 数据库迁移工具选型

EZ Admin Gin 之前一直采用手动执行 SQL 建表，没有引入迁移工具。这种方式在项目初期足够简单，但随着系统表增多、种子数据变复杂，几个实际问题开始显现：

- **没有迁移历史** — 谁改了哪个字段、什么时候改的，无迹可查
- **无法回滚** — 上线出问题没法回退到上一个版本
- **种子数据和启动代码耦合** — 角色、菜单、权限规则全部写在 `bootstrap.go` 里，每次启动都要查重判断，影响启动性能

::: tip 项目决策
引入 **golang-migrate** 管理建表和种子数据，把 DDL 和静态初始化数据全部写成 SQL 迁移文件，编号管理、可审查、可回滚。管理员账号（需要 bcrypt）改为一次性初始化接口，不再在启动流程中硬编码。
:::

## 候选工具

社区主流的 Go 数据库迁移工具有三款：

| 工具 | GitHub | 迁移格式 | 核心思路 |
| --- | --- | --- | --- |
| [Goose](https://github.com/pressly/goose) | ~8k stars | SQL + Go | 简洁实用，同时支持纯 SQL 和 Go 代码迁移 |
| [golang-migrate](https://github.com/golang-migrate/migrate) | ~16k stars | SQL | 最多人用，纯 SQL 文件，数据库支持最广 |
| [Atlas](https://atlasgo.io) | ~6k stars | SQL + HCL | 声明式 schema 管理，可自动生成迁移文件 |

## 对比

| 维度 | Goose | golang-migrate | Atlas |
| --- | --- | --- | --- |
| **迁移格式** | SQL + Go | 仅 SQL | SQL + HCL + Go |
| **Up / Down 回滚** | ✅ | ✅ | ✅ |
| **Go 代码迁移** | ✅ 原生支持 | ❌ | ✅ |
| **CLI + Go 库** | ✅ 两者都有 | ✅ 两者都有 | ✅ 两者都有 |
| **自动生成迁移** | ❌ 需手写 | ❌ 需手写 | ✅ 从 schema diff 自动生成 |
| **声明式 schema** | ❌ | ❌ | ✅ |
| **数据库支持** | PostgreSQL、MySQL、SQLite、ClickHouse 等 | 20+ 种（最广） | PostgreSQL、MySQL、SQLite、SQL Server 等 |
| **ORM 集成** | 基础 | 基础 | GORM、Ent、Prisma 深度集成 |
| **CI/CD 集成** | 手动配置 | 手动配置 | 内置 lint 和 CI 支持 |
| **脏状态处理** | 干净，版本锁简单 | 需要手动 `force` | 自动修复 |
| **学习曲线** | 低 | 低 | 中高 |
| **复杂度** | 轻量 | 轻量 | 较重，概念多 |

## 各工具特点

### Goose

Goose 的核心优势是**同时支持纯 SQL 和 Go 代码迁移**。

```sql
-- +goose Up
CREATE TABLE users (
  id         BIGSERIAL PRIMARY KEY,
  username   VARCHAR(64) NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- +goose Down
DROP TABLE users;
```

当你需要做复杂的数据迁移（比如拆分字段、批量转换数据格式）时，可以用 Go 迁移文件代替 SQL，直接操作 GORM 或原生 `sql.DB`。

适合：想要简洁工具、偶尔需要 Go 迁移、不想引入复杂概念的项目。

### golang-migrate

社区使用量最大，数据库支持最广，纯 SQL 文件。迁移文件分 `_up.sql` 和 `_down.sql` 两个文件：

```sql
-- 000001_create_users_table.up.sql
CREATE TABLE users (...);
```

```sql
-- 000001_create_users_table.down.sql
DROP TABLE users;
```

遇到脏状态（迁移中途失败）时需要手动 `migrate force <version>` 恢复，这是被吐槽最多的点。

适合：数据库种类多、团队习惯纯 SQL、追求最简单方案的项目。

### Atlas

最现代的方案，核心思路是"声明式 schema 管理"——你定义目标状态，Atlas 自动计算差异并生成迁移文件。

```hcl
table "users" {
  schema = schema.public
  column "id" {
    type = bigserial
  }
  column "username" {
    type = varchar(64)
  }
  primary_key {
    columns = [column.id]
  }
}
```

支持从 GORM 模型自动生成迁移（`atlas migrate diff`），也支持版本化迁移模式（和 Goose / golang-migrate 类似的手写 SQL 方式）。

功能最强大，但概念多、学习曲线较陡，对小型后台底座来说可能过重。

适合：团队较大、schema 变更频繁、重视 CI/CD 迁移检查的项目。

## 本项目选择

EZ Admin Gin 的特点：

- 同时支持 PostgreSQL 和 MySQL，需要跨数据库兼容
- 已有 GORM 模型，迁移工具需要和 GORM 共存
- 项目规模不大，复杂的数据迁移场景较少
- 种子数据（角色、菜单、权限规则、角色-菜单绑定）全部可以用纯 SQL 表达
- 优先选择简洁、学习成本低、维护活跃的方案

基于以上特点，**golang-migrate** 是比较合适的选择：

1. **数据库支持最广** — 原生支持 20+ 数据库，本项目同时需要 PostgreSQL 和 MySQL，两套迁移文件可以按 `migrations/pgsql/` 和 `migrations/mysql/` 分别管理
2. **纯 SQL 迁移** — 建表和种子数据都用 SQL 写，直观、可审查、不依赖 Go 运行时
3. **轻量** — 不引入额外概念，`Up()` / `Down()` 即可
4. **Go 库集成** — 可以通过 `embed.FS` 嵌入迁移文件，在应用启动时自动执行
5. **社区最大** — GitHub 16k stars，问题排查资料丰富

::: details Goose 也不错，为什么不选？
Goose 和 golang-migrate 在纯 SQL 迁移上几乎等价。Goose 额外支持 Go 迁移文件，但本项目的种子数据（角色、菜单、权限规则、角色-菜单绑定）都可以用纯 SQL 的 `INSERT` 搞定，不需要 Go 运行时。唯一需要 Go 代码的是管理员用户的 bcrypt 密码哈希，这部分改为通过一次性初始化接口处理，不走迁移文件。

选 golang-migrate 的核心理由是数据库支持更广——本项目需要同时维护 PostgreSQL 和 MySQL 两套 DDL。
:::

::: details Atlas 为什么不选？
Atlas 功能最强大，但对本项目来说概念偏多（HCL schema、声明式模式、Atlas Cloud），引入后需要团队额外学习一套工具链。如果后续项目规模增长、需要自动生成迁移或 CI/CD 检查，可以再考虑迁移到 Atlas。
:::
