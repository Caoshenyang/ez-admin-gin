---
title: 部署验证与复用说明
description: "完成完整的部署验证流程，并把后台底座复用到下一个个人项目。"
---

# 部署验证与复用说明

这一节会完成两件事：先把整个部署流程从头到尾跑一遍，验证后台底座在 Docker 环境下能正常工作；再说明如何把这套代码复用到你自己的新项目里。

::: tip 🎯 本节目标
跟着步骤走完之后，你能通过浏览器访问后台、登录系统、看到完整菜单，并确认业务接口可用。之后如果需要开一个新项目，也知道该怎么改。
:::

## 前置条件

在开始部署之前，确认本机已经安装了：

- Docker（20.10 及以上）
- Docker Compose（V2，支持 `docker compose` 命令）

如果不确定版本，可以执行：

```bash
docker --version
docker compose version
```

两者都能正常返回版本号即可。

## 🚀 部署步骤

### 第一步：准备环境变量

进入部署目录，复制模板并修改 JWT 密钥：

```bash
cd deploy
cp .env.example .env
```

打开 `.env` 文件，找到 `EZ_AUTH_JWT_SECRET`，替换为一个随机字符串：

```bash
# 生成一个随机密钥（macOS / Linux）
openssl rand -hex 32
```

把生成的字符串填入 `.env`：

```bash
EZ_AUTH_JWT_SECRET=你刚才生成的随机字符串
```

::: warning ⚠️ 不改密钥会启动失败
`compose.prod.yml` 中 `EZ_AUTH_JWT_SECRET` 使用了 `${EZ_AUTH_JWT_SECRET:?JWT_SECRET is required}` 语法，意味着如果这个变量为空或不存在，Docker Compose 会直接报错退出。这是故意的——防止你带着开发密钥上线。
:::

### 第二步：构建并启动所有服务

在 `deploy/` 目录下执行：

```bash
docker compose -f compose.prod.yml up -d --build
```

这条命令会：

- 构建后端和前端镜像
- 启动 PostgreSQL、Redis、后端服务和 Nginx 四个容器
- 等待数据库和 Redis 的健康检查通过后，再启动后端服务

首次构建可能需要几分钟，取决于网络速度。

### 第三步：等待服务就绪

使用以下命令查看容器状态：

```bash
docker compose -f compose.prod.yml ps
```

所有容器的状态应该显示为 `running`（或 `healthy`）：

```text
NAME                 STATUS
ez-admin-postgres    running (healthy)
ez-admin-redis       running (healthy)
ez-admin-server      running
ez-admin-nginx       running
```

如果后端服务刚启动，可以查看初始化日志：

```bash
docker compose -f compose.prod.yml logs server
```

日志中应该能看到类似这些信息：

```text
default admin user created   username=admin
default admin role created   role_code=super_admin
default menu created         menu_code=system
...
```

::: details 日志怎么看？
- `docker compose -f compose.prod.yml logs` 查看所有服务的日志。
- 加上服务名（如 `logs server`）只看某一个服务。
- 加 `-f` 参数可以持续跟踪：`logs -f server`。
:::

## ✅ 部署验证清单

服务全部启动后，按顺序验证以下四项：

### 1. 健康检查

```bash
curl http://localhost/health
```

期望返回类似：

```json
{
  "code": 0,
  "message": "ok",
  "data": {
    "env": "prod",
    "database": "ok",
    "redis": "ok"
  }
}
```

### 2. 登录

用浏览器打开 `http://localhost`，进入登录页面。输入默认管理员账号：

| 项目 | 值 |
| --- | --- |
| 用户名 | `admin` |
| 密码 | `EzAdmin@123456` |

登录成功后，页面应该跳转到后台首页，侧边栏显示系统管理菜单。

### 3. 菜单加载

登录后检查侧边栏，应该能看到以下菜单项：

```text
系统管理
  ├── 系统状态
  ├── 用户管理
  ├── 角色管理
  ├── 菜单管理
  ├── 系统配置
  ├── 文件管理
  ├── 操作日志
  ├── 登录日志
  └── 公告管理
```

点击“系统状态”后，页面里应该能看到当前环境（例如 `prod`）以及 `database = ok`、`redis = ok` 的检查结果。这一步可以顺手验证后台菜单、登录态和依赖状态页都已经接通。

### 4. CRUD 操作验证

进入任意一个管理页面（比如用户管理），确认以下操作可用：

- 列表数据正常加载
- 新增记录可以提交
- 编辑记录可以保存
- 状态切换有响应

::: tip 💡 快速排查思路
如果菜单可见但接口报 403，检查角色是否绑定了对应菜单权限。如果接口返回数据库错误，查看后端日志确认数据库连接是否正常。
:::

## 🛠️ 复用：开始一个新项目

后台底座的设计初衷就是让你快速复用到不同的个人项目中。下面是完整的复用步骤。

### 第一步：复制仓库

把仓库复制到你自己的项目目录：

```bash
cp -r ez-admin-gin my-new-project
cd my-new-project
```

或者直接 Fork 仓库后 Clone 到本地。

### 第二步：修改模块名称

打开 `server/go.mod`，把模块名改成你自己的项目名：

```go
module my-new-project/server
```

然后在 `server/` 目录下执行一次替换，把所有 import 路径中的旧模块名更新为新模块名：

```bash
# macOS / Linux
cd server
find . -name "*.go" -exec sed -i '' 's|ez-admin-gin/server|my-new-project/server|g' {} +
```

::: warning ⚠️ 替换后需要重新整理依赖
模块名改完之后，执行 `go mod tidy` 确保依赖和 import 路径一致。
:::

### 第三步：更新配置

修改以下位置，把默认的 `ez-admin` 相关名称换成你自己的：

| 文件 | 修改项 |
| --- | --- |
| `server/configs/config.yaml` | 应用名称、数据库名等 |
| `deploy/.env.example` | 复制为 `.env`，修改密钥和数据库信息 |
| `deploy/compose.prod.yml` | 项目名称、容器名称（可选） |

### 第四步：添加业务模块

按照第 6 章的模块接入规范，在 `server/internal/` 下创建新的业务模块目录。新增模块的结构和注册方式与系统模块一致。

### 第五步：更新前端

在 `admin/src/` 下创建业务模块对应的页面、API 和路由，复用已有的布局、请求封装和类型定义。

### 第六步：部署到生产环境

更新 `.env` 中的生产配置：

```bash
EZ_APP_ENV=prod
EZ_AUTH_JWT_SECRET=生产环境随机密钥
EZ_DATABASE_PASSWORD=生产环境数据库密码
EZ_LOG_FORMAT=json
EZ_NGINX_PORT=80
```

然后按前面讲过的部署流程启动即可。

::: details 复用后还需要注意什么？
- **修改默认管理员密码**：首次登录后立即在用户管理页面修改。
- **调整上传目录**：如果业务需要不同的文件存储路径，修改 `EZ_UPLOAD_DIR`。
- **关闭调试信息**：确认 `EZ_APP_ENV=prod`，日志级别不低于 `info`。
- **数据库连接池**：小型项目用默认值即可，并发量上来后再调整 `EZ_DATABASE_MAX_OPEN_CONNS` 等参数。
:::

## 常见问题排查

部署过程中如果遇到问题，按下面的清单逐项排查。

::: details 容器启动失败：JWT_SECRET is required
`.env` 文件中没有填写 `EZ_AUTH_JWT_SECRET`，或者 `.env` 文件不在 `deploy/` 目录下。

确认方法：在 `deploy/` 目录下执行 `cat .env | grep JWT_SECRET`，应该能看到你填写的密钥。
:::

::: details 后端服务不断重启，日志报数据库连接失败
可能原因：数据库容器还没有通过健康检查，后端就尝试连接了。

排查方法：`docker compose -f compose.prod.yml ps` 确认 postgres 容器状态为 `healthy`。如果一直是 `starting`，查看数据库日志：`docker compose -f compose.prod.yml logs postgres`。
:::

::: details 前端页面白屏，浏览器控制台报 502
后端服务还没有完全启动，Nginx 反向代理找不到上游。

排查方法：确认 server 容器正常运行，检查 `docker compose -f compose.prod.yml logs server`，等待看到路由注册完成的日志。
:::

::: details 登录后侧边栏没有菜单
可能原因：`bootstrap.go` 的初始化数据没有写入成功，或者角色菜单绑定出了问题。

排查方法：查看后端启动日志中是否有 `default menu created` 和 `default role menu bound` 相关记录。
:::

::: details 修改 .env 后配置没有生效
修改 `.env` 后需要重启服务：

```bash
docker compose -f compose.prod.yml down
docker compose -f compose.prod.yml up -d
```

只是重启单个服务不会重新读取 `.env`。
:::

## 小结

到这里，第 7 章的所有内容就完成了。回顾一下本章做了什么：

- 为后端和前端分别编写了 Dockerfile，实现了多阶段构建。
- 用 Docker Compose 把 PostgreSQL、Redis、后端和 Nginx 四个服务编排在一起。
- 配置了 Nginx 静态资源托管和 API 反向代理。
- 整理了所有环境变量，理解了 `EZ_` 前缀的覆盖机制。
- 了解了 `bootstrap.go` 如何自动创建管理员、角色、菜单和权限数据。
- 完成了从构建到验证的完整部署流程。
- 说明了如何把这套底座复用到新的个人项目。

回到本章总览：[第 7 章：部署与复用](./)。
