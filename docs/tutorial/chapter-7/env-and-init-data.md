---
title: 环境变量与初始化数据
description: "整理部署所需环境变量，理解配置覆盖机制，并了解启动时如何通过迁移自动初始化角色、权限和菜单。"
---

# 环境变量与初始化数据

这一节会完成两件事：把所有部署相关的环境变量讲清楚，让你知道哪些必须改、哪些可以保持默认；然后了解服务启动时如何通过 golang-migrate 自动完成数据库建表和种子数据准备。

::: tip 🎯 本节目标
读完这一节，你能回答三个问题：环境变量有哪些、哪些必须手动改、第一次启动后系统会自动准备好什么。
:::

## 环境变量机制

后台底座使用 [Viper](https://pkg.go.dev/github.com/spf13/viper) 读取配置，所有配置项都支持通过 `EZ_` 前缀的环境变量覆盖。覆盖规则很简单：

- 配置键 `database.host` 对应环境变量 `EZ_DATABASE_HOST`。
- 点号替换为下划线，统一加上 `EZ_` 前缀。
- 环境变量优先级高于配置文件中的默认值。

这意味着你可以只维护一份 `.env` 文件，不用手动改 `config.yaml`。

### 完整环境变量清单

下面是 `.env.example` 的完整内容，你可以在 `deploy/.env.example` 中找到它：

::: details deploy/.env.example — 部署环境变量模板
<<< ../../../deploy/.env.example
:::

### 环境变量参考表

下面的表格按功能分组，标注了每个变量的默认值、是否必填和用途。默认值来自 `config.go` 中的 `setDefaults` 函数。

<table>
  <colgroup>
    <col style="width: 15rem;">
    <col style="width: 8rem;">
    <col style="width: 5rem;">
    <col>
  </colgroup>
  <thead>
    <tr>
      <th>变量名</th>
      <th>默认值</th>
      <th>必填</th>
      <th>说明</th>
    </tr>
  </thead>
  <tbody>
    <tr>
      <td colspan="4"><strong>认证</strong></td>
    </tr>
    <tr>
      <td><code>EZ_AUTH_JWT_SECRET</code></td>
      <td><code>ez-admin-dev-secret-...</code></td>
      <td>✅ 是</td>
      <td>JWT 签名密钥。生产环境必须替换为 32 位以上的随机字符串</td>
    </tr>
    <tr>
      <td><code>EZ_AUTH_ACCESS_TOKEN_TTL</code></td>
      <td><code>7200</code>（2 小时）</td>
      <td>否</td>
      <td>Access Token 有效期，单位秒</td>
    </tr>
    <tr>
      <td><code>EZ_AUTH_ISSUER</code></td>
      <td><code>ez-admin</code></td>
      <td>否</td>
      <td>Token 签发方标识</td>
    </tr>
    <tr>
      <td colspan="4"><strong>数据库</strong></td>
    </tr>
    <tr>
      <td><code>EZ_DATABASE_DRIVER</code></td>
      <td><code>postgres</code></td>
      <td>否</td>
      <td>数据库驱动：<code>postgres</code> 或 <code>mysql</code></td>
    </tr>
    <tr>
      <td><code>EZ_DATABASE_HOST</code></td>
      <td><code>localhost</code></td>
      <td>否</td>
      <td>数据库主机地址。Docker 环境下通常是容器名 <code>postgres</code></td>
    </tr>
    <tr>
      <td><code>EZ_DATABASE_PORT</code></td>
      <td><code>5432</code></td>
      <td>否</td>
      <td>数据库端口</td>
    </tr>
    <tr>
      <td><code>EZ_DATABASE_USER</code></td>
      <td><code>ez_admin</code></td>
      <td>否</td>
      <td>数据库用户名</td>
    </tr>
    <tr>
      <td><code>EZ_DATABASE_PASSWORD</code></td>
      <td><code>ez_admin_123456</code></td>
      <td>否</td>
      <td>数据库密码。生产环境建议替换</td>
    </tr>
    <tr>
      <td><code>EZ_DATABASE_NAME</code></td>
      <td><code>ez_admin</code></td>
      <td>否</td>
      <td>数据库名称</td>
    </tr>
    <tr>
      <td colspan="4"><strong>Redis</strong></td>
    </tr>
    <tr>
      <td><code>EZ_REDIS_HOST</code></td>
      <td><code>localhost</code></td>
      <td>否</td>
      <td>Redis 主机地址。Docker 环境下通常是容器名 <code>redis</code></td>
    </tr>
    <tr>
      <td><code>EZ_REDIS_PORT</code></td>
      <td><code>6379</code></td>
      <td>否</td>
      <td>Redis 端口</td>
    </tr>
    <tr>
      <td><code>EZ_REDIS_PASSWORD</code></td>
      <td><em>空</em></td>
      <td>否</td>
      <td>Redis 密码。本地开发可以为空</td>
    </tr>
    <tr>
      <td colspan="4"><strong>应用与服务</strong></td>
    </tr>
    <tr>
      <td><code>EZ_APP_ENV</code></td>
      <td><code>dev</code></td>
      <td>否</td>
      <td>运行环境标识：<code>dev</code> 或 <code>prod</code></td>
    </tr>
    <tr>
      <td><code>EZ_SERVER_ADDR</code></td>
      <td><code>:8080</code></td>
      <td>否</td>
      <td>后端 HTTP 服务监听地址</td>
    </tr>
    <tr>
      <td colspan="4"><strong>日志</strong></td>
    </tr>
    <tr>
      <td><code>EZ_LOG_LEVEL</code></td>
      <td><code>info</code></td>
      <td>否</td>
      <td>日志级别：<code>debug</code> / <code>info</code> / <code>warn</code> / <code>error</code></td>
    </tr>
    <tr>
      <td><code>EZ_LOG_FORMAT</code></td>
      <td><code>console</code></td>
      <td>否</td>
      <td>日志格式：<code>console</code> 适合开发，<code>json</code> 适合生产采集</td>
    </tr>
    <tr>
      <td colspan="4"><strong>Nginx</strong></td>
    </tr>
    <tr>
      <td><code>EZ_NGINX_PORT</code></td>
      <td><code>80</code></td>
      <td>否</td>
      <td>Nginx 对外暴露端口</td>
    </tr>
  </tbody>
</table>

::: warning ⚠️ 生产环境至少改两项
`EZ_AUTH_JWT_SECRET` 是唯一标记为"必填"的环境变量，Docker Compose 会在它缺失时直接报错退出。除此之外，建议把 `EZ_DATABASE_PASSWORD` 也换成生产级别的密码，避免使用默认值。
:::

::: details 还有哪些可选变量？
除了上面列出的常用变量，`config.go` 还支持数据库连接池（`EZ_DATABASE_MAX_IDLE_CONNS`、`EZ_DATABASE_MAX_OPEN_CONNS`、`EZ_DATABASE_CONN_MAX_LIFETIME`）、Redis 连接池（`EZ_REDIS_POOL_SIZE`、`EZ_REDIS_MIN_IDLE_CONNS`、`EZ_REDIS_MAX_RETRIES`）和上传配置（`EZ_UPLOAD_DIR`、`EZ_UPLOAD_MAX_SIZE_MB` 等）。这些都有合理默认值，小型后台起步时不需要手动调整。
:::

## 启动时自动初始化

服务第一次启动时，会通过 golang-migrate 自动执行 `server/migrations/pgsql/`（或 `mysql/`）目录下的 SQL 迁移文件，完成数据库建表和种子数据写入。

### 迁移执行顺序

```text
000001_init_schema.up.sql
  → 创建所有系统表（用户、角色、菜单、配置、文件、日志、公告、权限策略等）
  ↓
000002_seed_data.up.sql
  → 插入超级管理员角色（super_admin，ID=1）
  → 插入系统菜单和按钮（目录 → 菜单 → 按钮，使用固定 ID）
  → 插入 Casbin 接口权限规则
  → 插入角色-菜单绑定关系
  → PostgreSQL 版本还会重置序列计数器
```

golang-migrate 通过 `schema_migrations` 表追踪已执行的迁移版本。如果所有迁移都已执行，启动时会跳过并输出 `database migrations up to date`。

::: tip 📌 迁移文件在哪里
迁移文件位于 `server/migrations/` 目录，按数据库驱动分目录：

```text
server/migrations/
├── pgsql/
│   ├── 000001_init_schema.up.sql
│   ├── 000001_init_schema.down.sql
│   ├── 000002_seed_data.up.sql
│   └── 000002_seed_data.down.sql
└── mysql/
    ├── 000001_init_schema.up.sql
    ├── 000001_init_schema.down.sql
    ├── 000002_seed_data.up.sql
    └── 000002_seed_data.down.sql
```

程序启动时通过 `embed.FS` 嵌入这些文件，不需要额外部署迁移文件到服务器。
:::

### 管理员账号初始化

管理员账号**不在迁移文件中创建**，因为密码需要 bcrypt 加密，纯 SQL 无法完成。管理员通过一次性初始化接口创建：

```bash
curl -X POST http://localhost:8080/api/v1/setup/init \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"YourPassword123","nickname":"管理员"}'
```

这个接口会：

1. 检查 `sys_user` 是否已有记录 → 有则返回 `409`，不允许重复初始化
2. bcrypt 加密密码
3. 插入 `sys_user`
4. 绑定到 `super_admin` 角色（ID=1）

::: warning ⚠️ 部署后立即调用初始化接口
服务启动后、管理员账号创建前，登录接口不可用。部署验证时，应把调用 `/api/v1/setup/init` 作为启动后的第一个操作。
:::

### 超级管理员角色

迁移文件自动创建的角色：

| 项目 | 值 |
| --- | --- |
| 角色编码 | `super_admin` |
| 角色名称 | 超级管理员 |
| 角色ID | `1`（固定） |
| 权限范围 | 所有系统接口 + 所有系统菜单 |

### 系统菜单结构

菜单按 目录 → 菜单 → 按钮 三级层次创建，所有 ID 固定分配。创建完成后，系统管理目录下会包含以下菜单和按钮权限：

| 菜单 | 类型 | 路由 | 按钮权限 |
| --- | --- | --- | --- |
| 系统管理 | 目录 | `/system` | — |
| 系统状态 | 菜单 | `/system/health` | 查看系统状态 |
| 用户管理 | 菜单 | `/system/users` | 查看、创建、编辑、修改状态、分配角色 |
| 角色管理 | 菜单 | `/system/roles` | 查看、创建、编辑、修改状态、分配接口权限、分配菜单权限 |
| 菜单管理 | 菜单 | `/system/menus` | 查看、创建、编辑、修改状态、删除 |
| 系统配置 | 菜单 | `/system/configs` | 查看、创建、编辑、修改状态、读取配置值 |
| 文件管理 | 菜单 | `/system/files` | 查看、上传 |
| 操作日志 | 菜单 | `/system/operation-logs` | 查看 |
| 登录日志 | 菜单 | `/system/login-logs` | 查看 |
| 公告管理 | 菜单 | `/system/notices` | 查看、创建、编辑、修改状态 |

其中"系统状态"页面会调用受保护的 `/api/v1/system/health`，把当前运行环境、数据库和 Redis 的检查结果展示在管理台里；而公开的 `/health` 仍然保留给部署探针和外部监控使用。

### 接口权限规则

迁移文件会在 Casbin 规则表中为 `super_admin` 角色写入所有系统接口的访问规则。每条规则包含角色编码、接口路径和 HTTP 方法三部分。

这意味着超级管理员角色可以直接访问所有 `/api/v1/system/*` 路径的接口，不需要额外配置。

::: details 幂等性是怎么保证的？
golang-migrate 通过 `schema_migrations` 表记录当前迁移版本。已执行的迁移不会重复运行。如果需要重新初始化整个数据库，删掉数据库或对应数据卷即可。
:::

## 小结

这一节把部署时的配置和初始数据讲清楚了。核心要点：

- 所有配置项都可以通过 `EZ_` 前缀的环境变量覆盖，优先级高于配置文件。
- 生产环境只需修改 `EZ_AUTH_JWT_SECRET`，建议同时修改数据库密码。
- 服务启动时通过 golang-migrate 自动执行建表和种子数据迁移，幂等安全。
- 管理员账号通过 `/api/v1/setup/init` 接口创建，部署后应立即调用。

下一节进入实际的部署操作和验收：[部署验证与复用说明](./deployment-and-reuse)。
