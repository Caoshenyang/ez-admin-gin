---
title: 数据库迁移工具选型
description: "对比 Goose、golang-migrate 和 Atlas 三款 Go 数据库迁移工具，为本项目选择合适的方案。"
---

# 数据库迁移工具选型

EZ Admin Gin 之前一直采用手动执行 SQL 建表，没有引入迁移工具。这种方式在项目初期足够简单，但随着系统表增多、种子数据变复杂，几个实际问题开始显现：

- **没有迁移历史** — 谁改了哪个字段、什么时候改的，无迹可查
- **无法回滚** — 上线出问题没法回退到上一个版本
- **种子数据和启动代码耦合** — 角色、菜单、权限规则全部写在迁移文件里，每次启动都要查重判断，影响启动性能

::: tip 项目决策
引入 **golang-migrate** 管理建表和种子数据，把 DDL 和静态初始化数据全部写成 SQL 迁移文件，编号管理、可审查、可回滚。管理员账号（需要 bcrypt）改为一次性初始化接口，不再在启动流程中硬编码。
:::

## golang-migrate 简介

**golang-migrate** 是一个轻量级的数据库迁移工具，专注于纯 SQL 文件的版本化管理。它的核心特点是：

- **纯 SQL 迁移**：使用 `.up.sql` 和 `.down.sql` 文件分别定义升级和回滚操作
- **版本化管理**：通过文件编号（如 `000001_`）确保迁移顺序
- **广泛的数据库支持**：原生支持 20+ 种数据库，包括 PostgreSQL、MySQL、SQLite 等
- **Go 库集成**：可以通过 `embed.FS` 嵌入迁移文件，在应用启动时自动执行
- **命令行工具**：提供 `migrate` CLI 工具，支持手动执行迁移操作

### 基本使用方式

1. **创建迁移文件**：在 `migrations/{postgres,mysql}/` 目录下创建编号命名的 SQL 文件
   - `000001_init_schema.up.sql`：创建表结构
   - `000001_init_schema.down.sql`：回滚表结构
   - `000002_seed_data.up.sql`：插入种子数据
   - `000002_seed_data.down.sql`：清除种子数据

2. **执行迁移**：在应用启动时通过 `golang-migrate` 库自动执行

3. **回滚迁移**：需要时执行 `migrate down` 命令或调用 `m.Down()` 方法

### 项目中的集成

本项目通过 `internal/migrate/migrate.go` 封装了迁移逻辑：

- 使用 `embed.FS` 嵌入 `migrations/` 目录下的所有 SQL 文件
- 根据数据库驱动（`postgres` 或 `mysql`）加载对应子目录的迁移文件
- 在应用启动时自动执行 `migrate.Up()`，确保数据库结构和种子数据与代码版本一致

这种方式既保证了数据库变更的可追踪性，又简化了部署流程，是中小型项目的理想选择。

## golang-migrate 整合步骤

下面是在 Go 项目中整合 golang-migrate 的详细步骤，适合第一次接触的开发者：

### 步骤 1：安装依赖

在项目的 `server/` 目录下执行：

```bash
# 安装 golang-migrate 库（会自动包含所有子包）
go get github.com/golang-migrate/migrate/v4@latest
```

**说明**：虽然代码中会使用 `database/mysql`、`database/postgres` 和 `source/iofs` 等子包，但 Go 的依赖管理机制会自动处理。安装主库后，所有子包都可以直接使用，无需单独安装。

### 步骤 2：创建迁移文件目录结构

在项目中创建如下目录结构：

```text
server/
└─ migrations/
   ├── postgres/          # PostgreSQL 迁移文件
   │   ├── 000001_init_schema.up.sql
   │   ├── 000001_init_schema.down.sql
   │   ├── 000002_seed_data.up.sql
   │   └── 000002_seed_data.down.sql
   └── mysql/          # MySQL 迁移文件
       ├── 000001_init_schema.up.sql
       ├── 000001_init_schema.down.sql
       ├── 000002_seed_data.up.sql
       └── 000002_seed_data.down.sql
```

- **文件命名规则**：`{版本号}_{描述}.{up|down}.sql`
- **版本号**：使用 6 位数字，确保顺序正确
- **up.sql**：升级操作（创建表、插入数据等）
- **down.sql**：回滚操作（删除表、清除数据等）

### 步骤 3：编写迁移文件

#### 示例：创建表结构（000001_init_schema.up.sql）

```sql
-- PostgreSQL 版本
CREATE TABLE IF NOT EXISTS sys_user (
    id SERIAL PRIMARY KEY,
    username VARCHAR(64) NOT NULL UNIQUE,
    password VARCHAR(128) NOT NULL,
    nickname VARCHAR(64) NOT NULL,
    status INTEGER NOT NULL DEFAULT 1,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP
);

-- MySQL 版本
CREATE TABLE IF NOT EXISTS sys_user (
    id INT AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(64) NOT NULL UNIQUE,
    password VARCHAR(128) NOT NULL,
    nickname VARCHAR(64) NOT NULL,
    status INT NOT NULL DEFAULT 1,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at DATETIME
);
```

#### 示例：回滚表结构（000001_init_schema.down.sql）

```sql
-- PostgreSQL 版本
DROP TABLE IF EXISTS sys_user;

-- MySQL 版本
DROP TABLE IF EXISTS sys_user;
```

#### 示例：插入种子数据（000002_seed_data.up.sql）

```sql
-- 插入超级管理员角色
INSERT INTO sys_role (id, code, name, status, remark, created_at, updated_at)
VALUES (1, 'super_admin', '超级管理员', 1, '系统内置角色', NOW(), NOW())
ON CONFLICT (code) DO NOTHING;

-- 插入系统管理菜单
INSERT INTO sys_menu (id, parent_id, type, code, title, path, component, icon, sort, status, remark, created_at, updated_at)
VALUES (1, 0, 1, 'system', '系统管理', '/system', '', 'setting', 10, 1, '系统内置目录', NOW(), NOW())
ON CONFLICT (code) DO NOTHING;
```

### 步骤 4：在 Go 代码中集成

创建 `internal/migrate/migrate.go` 文件：

```go
package migrate

import (
	"embed"
	"fmt"
	"io/fs"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"go.uber.org/zap"
)

// Run 执行数据库迁移
func Run(driver, dsn string, migrationsFS embed.FS, log *zap.Logger) error {
	// 根据驱动加载对应子目录
	sub, err := fs.Sub(migrationsFS, "migrations/"+driver)
	if err != nil {
		return fmt.Errorf("open migrations/%s: %w", driver, err)
	}

	// 创建 iofs 源
	source, err := iofs.New(sub, ".")
	if err != nil {
		return fmt.Errorf("create migration source: %w", err)
	}

	// 创建 migrate 实例
	m, err := migrate.NewWithSourceInstance("iofs", source, dsn)
	if err != nil {
		return fmt.Errorf("create migrate instance: %w", err)
	}
	defer m.Close()

	// 执行迁移
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("run migrations: %w", err)
	}

	// 输出日志
	if err == migrate.ErrNoChange {
		log.Info("database migrations up to date", zap.String("driver", driver))
	} else {
		log.Info("database migrations applied", zap.String("driver", driver))
	}

	return nil
}
```

### 步骤 5：在主程序中使用

修改 `main.go`，在启动时执行迁移：

```go
package main

import (
	"embed"

	"ez-admin-gin/server/internal/config"
	"ez-admin-gin/server/internal/database"
	"ez-admin-gin/server/internal/migrate"
	appLogger "ez-admin-gin/server/internal/logger"

	"go.uber.org/zap"
)

//go:embed migrations/postgres migrations/mysql
var migrationsFS embed.FS

func main() {
	// 加载配置
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}

	// 创建日志
	log, err := appLogger.New(cfg.Log)
	if err != nil {
		panic(err)
	}

	// 连接数据库
	db, err := database.New(cfg.Database, log)
	if err != nil {
		log.Fatal("connect database", zap.Error(err))
	}

	// 执行迁移
	migrateDSN, err := database.MigrateDSN(cfg.Database)
	if err != nil {
		log.Fatal("build migration dsn", zap.Error(err))
	}
	if err := migrate.Run(cfg.Database.Driver, migrateDSN, migrationsFS, log); err != nil {
		log.Fatal("run database migrations", zap.Error(err))
	}

	// 启动服务...
}
```

### 步骤 6：验证迁移

启动服务后，查看日志：

```text
INFO database migrations applied {"driver": "postgres"}
INFO server started {"addr": ":8080", "env": "dev"}
```

如果看到 `database migrations applied`，说明迁移执行成功。

### 常见问题处理

1. **迁移文件未执行**
   - 检查 `schema_migrations` 表，确认迁移版本是否已记录
   - 检查迁移文件路径和命名是否正确

2. **迁移失败**
   - 查看错误信息，通常是 SQL 语法错误或依赖表不存在
   - 修复后重新启动服务

3. **脏状态处理**
   - 如果迁移中途失败，可能会留下脏状态
   - 使用 `migrate force <version>` 命令恢复到指定版本

4. **回滚迁移**
   - 如需回滚，调用 `m.Down()` 方法或使用 `migrate down` 命令

### 最佳实践

- **版本控制**：将迁移文件纳入版本控制，确保团队协作时的一致性
- **幂等性**：使用 `IF NOT EXISTS` 和 `ON CONFLICT DO NOTHING` 确保迁移可以重复执行
- **分阶段**：将建表和种子数据分开，便于管理和回滚
- **测试**：在开发环境充分测试迁移，避免生产环境出错
- **备份**：执行迁移前备份数据库，尤其是生产环境

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

1. **数据库支持最广** — 原生支持 20+ 数据库，本项目同时需要 PostgreSQL 和 MySQL，两套迁移文件可以按 `migrations/postgres/` 和 `migrations/mysql/` 分别管理
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
