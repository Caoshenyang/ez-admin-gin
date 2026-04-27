---
title: 快速启动
description: "五步在本地跑起 EZ Admin Gin，附带初始化细节和常见问题。"
---

# 快速启动

EZ Admin Gin 是一个面向个人项目快速上线的通用后台管理系统底座——登录、权限、菜单、配置、日志、文件上传，这些每个后台都要写的东西，这里已经沉淀好了。

::: tip 🎯 这页解决什么
用最短的时间帮你在本地把项目跑起来，确认它是不是你想要的。
:::

## 📋 先确认环境

| 依赖 | 版本要求 | 用途 |
| --- | --- | --- |
| Go | >= 1.26 | 后端运行 |
| Node.js | >= 20.19 | 前端运行 |
| pnpm | 最新稳定版 | 前端包管理 |
| Docker & Docker Compose | 最新稳定版 | 本地 PostgreSQL 和 Redis |

::: info 没有安装 Docker？
如果你本地已经有 PostgreSQL 18 和 Redis 8，可以跳过第 1 步，直接修改 `server/configs/config.yaml` 里的数据库和 Redis 连接信息即可。
:::

## 🚀 五步跑起来

### 第 1 步：启动 PostgreSQL 和 Redis

::: code-group

```bash [macOS / Linux]
docker compose -f deploy/compose.local.yml up -d
```

```bash [Windows]
docker compose -f deploy/compose.local.win.yml up -d
```

:::

这一步会启动两个容器：

| 服务 | 本机端口 | 默认账号 |
| --- | --- | --- |
| PostgreSQL 18 | 5432 | 用户 `ez_admin` / 密码 `ez_admin_123456` / 数据库 `ez_admin` |
| Redis 8 | 6379 | 无密码 |

数据持久化到本机 `~/ez-admin-gin-data/` 目录，删掉容器数据不会丢。

::: details 验证服务是否启动成功
```bash
# PostgreSQL
docker exec ez-admin-postgres pg_isready -U ez_admin
# 应输出：accepting connections

# Redis
docker exec ez-admin-redis redis-cli ping
# 应输出：PONG
```
:::

### 第 2 步：启动后端

```bash
cd server
go run main.go
```

后端监听地址：`http://localhost:8080`

::: tip 🔑 这一步会自动完成数据库迁移
首次启动时，程序会通过 golang-migrate 自动执行 `server/migrations/pgsql/` 下的 SQL 迁移文件，**不需要手动建表或导入 SQL**：

1. **自动建表** — 执行 `000001_init_schema.up.sql`，创建所有系统表（用户、角色、菜单、配置、文件、日志、公告、权限策略等）
2. **初始化种子数据** — 执行 `000002_seed_data.up.sql`，自动写入：
   - 超级管理员角色（`super_admin`，ID 固定为 1）
   - 系统管理菜单（目录、菜单、按钮）
   - 全量接口权限规则（Casbin）
   - 超级管理员与所有菜单的绑定关系
3. **幂等安全** — golang-migrate 通过 `schema_migrations` 表追踪版本，已执行的迁移不会重复执行

如果你想重新初始化，删掉 PostgreSQL 容器和数据卷重来即可。
:::

::: details 支持哪些数据库
后端同时支持 **PostgreSQL** 和 **MySQL**。默认使用 PostgreSQL，切换方式：

1. 修改 `server/configs/config.yaml` 中 `database.driver` 为 `mysql`
2. 修改数据库连接信息（host、port、user、password、name）指向你的 MySQL 实例
3. 程序启动时会自动加载 `server/migrations/mysql/` 下的迁移文件
:::

::: details 配置文件说明
后端配置在 `server/configs/config.yaml`，已内置本地开发默认值（数据库连接 `localhost:5432`、Redis `localhost:6379` 等），开箱即用无需修改。

如果需要改端口或数据库密码，直接编辑该文件。生产环境建议通过环境变量覆盖，参考 [部署章节](/tutorial/chapter-7/)。
:::

### 第 3 步：创建管理员账号

首次启动后，数据库中还没有管理员用户（管理员密码需要 bcrypt 加密，不适合写在 SQL 中）。调用初始化接口创建：

::: code-group

```powershell [Windows PowerShell]
Invoke-RestMethod -Method Post -Uri http://localhost:8080/api/v1/setup/init -ContentType "application/json" -Body '{"username":"admin","password":"YourPassword123","nickname":"管理员"}'
```

```bash [macOS / Linux]
curl -X POST http://localhost:8080/api/v1/setup/init \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"YourPassword123","nickname":"管理员"}'
```

:::

成功后会返回：

```json
{
  "message": "管理员账号创建成功",
  "user_id": 1,
  "username": "admin"
}
```

::: warning ⚠️ 这个接口只能调用一次
如果数据库中已有用户记录，接口会返回 `409`，提示"系统已初始化，不能重复执行"。
:::

### 第 4 步：启动前端

```bash
cd admin
pnpm install
pnpm dev
```

启动后终端会打印前端访问地址。

### 第 5 步：登录系统

打开浏览器访问前端地址，使用第 3 步中设置的用户名和密码登录。

::: danger ⚠️ 上线前务必修改
- 替换 `auth.jwt_secret` 为随机字符串（至少 32 位）
- 确保管理员密码足够强壮
:::

## 🧭 接下来去哪

- **了解项目用了什么技术** → [项目结构](/guide/project-structure)
- **想学每一步怎么搭出来的** → [从零搭建教程](/tutorial/)
- **查配置、接口、建表语句** → [参考手册](/reference/)
