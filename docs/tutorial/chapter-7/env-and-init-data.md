---
title: 环境变量与初始化数据
description: "整理部署所需环境变量，理解配置覆盖机制，并了解启动时如何自动初始化管理员、角色、权限和菜单。"
---

# 环境变量与初始化数据

这一节会完成两件事：把所有部署相关的环境变量讲清楚，让你知道哪些必须改、哪些可以保持默认；然后了解服务启动时 `bootstrap.go` 如何自动完成初始数据准备。

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

服务第一次启动时，`bootstrap.go` 中的 `Run` 函数会自动执行一系列初始化动作。整个过程是幂等的：如果数据已经存在，不会重复创建。

### 初始化顺序

初始化按以下顺序依次执行，每一步都有明确的依赖关系：

```text
1. 创建默认管理员账号
   ↓
2. 创建超级管理员角色
   ↓
3. 创建系统菜单（目录 → 菜单 → 按钮）
   ↓
4. 创建 Casbin 接口权限规则
   ↓
5. 绑定管理员与角色
   ↓
6. 绑定角色与菜单
```

### 默认管理员账号

系统会自动创建一个管理员用户：

| 项目 | 值 |
| --- | --- |
| 用户名 | `admin` |
| 默认密码 | `EzAdmin@123456` |
| 昵称 | 系统管理员 |
| 状态 | 启用 |

::: warning ⚠️ 首次登录后请立即修改密码
默认密码仅用于首次部署验证，登录后请尽快在用户管理页面修改。
:::

### 超级管理员角色

自动创建的角色：

| 项目 | 值 |
| --- | --- |
| 角色编码 | `super_admin` |
| 角色名称 | 超级管理员 |
| 权限范围 | 所有系统接口 + 所有系统菜单 |

### 系统菜单结构

菜单按 目录 → 菜单 → 按钮 三级层次创建。创建完成后，系统管理目录下会包含以下菜单和按钮权限：

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

### 接口权限规则

除了菜单权限，`bootstrap.go` 还会在 Casbin 规则表中为 `super_admin` 角色写入所有系统接口的访问规则。每条规则包含角色编码、接口路径和 HTTP 方法三部分。

这意味着超级管理员角色可以直接访问所有 `/api/v1/system/*` 路径的接口，不需要额外配置。

::: details 幂等性是怎么保证的？
每条初始化数据在写入前都会先查询是否已存在。存在则跳过，不存在才创建。所以你可以放心地重启服务，不用担心数据被重复写入。

关键查询条件：
- 管理员账号：按 `username` 查询，并使用 `Unscoped` 避免逻辑删除记录导致重复创建。
- 角色：按 `code` 查询。
- 菜单：按 `code` 查询。
- Casbin 规则：按 `ptype` + `v0` + `v1` + `v2` 组合查询。
- 角色菜单绑定：按 `role_id` + `menu_id` 组合查询。
:::

## 小结

这一节把部署时的配置和初始数据讲清楚了。核心要点：

- 所有配置项都可以通过 `EZ_` 前缀的环境变量覆盖，优先级高于配置文件。
- 生产环境只需修改 `EZ_AUTH_JWT_SECRET`，建议同时修改数据库密码。
- 服务启动时会自动创建管理员、角色、菜单和权限规则，整个过程幂等安全。

下一节进入实际的部署操作和验收：[部署验证与复用说明](./deployment-and-reuse)。
