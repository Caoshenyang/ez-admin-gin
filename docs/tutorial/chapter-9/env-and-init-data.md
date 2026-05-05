---
title: 环境变量与初始化数据
description: "围绕当前 config 覆盖规则、.env 模板、迁移脚本和 setup/init 接口，讲清部署前必须确认的运行基线。"
---

# 环境变量与初始化数据

第 9 章先从最容易在真实部署中出问题的一层开始：

- 环境变量到底怎么覆盖配置
- 第一次启动后系统会自动准备什么
- 管理员账号到底什么时候创建

这一页对齐当前真实链路：

- `deploy/.env.example`
- `server/internal/platform/config`
- `server/migrations/*`
- `/api/v1/setup/init`

::: tip 🎯 本节目标
读完后，你应该能清楚判断：

1. 哪些变量必须改
2. 哪些初始化是迁移自动完成的
3. 哪些初始化需要你部署后手动或脚本触发
:::

## 先看当前真实部署基线

当前服务器部署依赖的一套核心文件是：

```text
deploy/.env.example
deploy/compose.server.yml
deploy/ez-admin.service
server/configs/config.yaml
server/migrations/postgres/*
server/migrations/mysql/*
```

它们分别负责：

| 文件 | 当前作用 |
| --- | --- |
| `deploy/.env.example` | 生产环境变量模板 |
| `compose.server.yml` | PostgreSQL、Redis、Nginx 基础容器 |
| `ez-admin.service` | 宿主机上的后端 systemd 服务 |
| `config.yaml` | 配置默认值基线 |
| `migrations/*` | 数据库结构和种子数据 |

## 环境变量是怎样覆盖配置的

当前后端配置使用的是：

- `Viper`

并统一支持 `EZ_` 前缀环境变量覆盖。

规则非常简单：

- `database.host` → `EZ_DATABASE_HOST`
- `auth.jwt_secret` → `EZ_AUTH_JWT_SECRET`

也就是说：

- 点号改成下划线
- 统一转大写
- 统一加 `EZ_` 前缀

这让部署时可以主要围绕一份 `.env` 工作，而不是直接改仓库里的配置文件。

## 当前最值得先改的变量有哪些

真实部署时，最应该优先确认的是下面几项：

| 变量 | 是否必须关注 | 原因 |
| --- | --- | --- |
| `EZ_AUTH_JWT_SECRET` | 必须 | 默认占位值不能进生产 |
| `EZ_DATABASE_PASSWORD` | 强烈建议 | 默认值过弱 |
| `EZ_APP_ENV` | 建议 | 生产环境应设为 `prod` |
| `EZ_LOG_FORMAT` | 建议 | 生产更适合 `json` |
| `EZ_SERVER_ADDR` | 视环境决定 | 默认监听 `:8080` |

::: warning ⚠️ 生产环境第一优先级就是换 JWT 密钥
当前 `deploy/.env.example` 里 `EZ_AUTH_JWT_SECRET` 只是占位值。`setup-server.sh` 首次部署时虽然会自动替换它，但你仍然应该理解这一步的重要性，因为它直接关系到 Token 是否可伪造。
:::

## 为什么 `compose.server.yml` 里很多地址是 `127.0.0.1`

当前服务器部署不是把后端也塞进 Docker，而是：

- PostgreSQL / Redis / Nginx 跑容器
- 后端二进制跑宿主机

所以 `.env` 里的数据库和 Redis 地址默认会用：

- `127.0.0.1`

这是因为容器端口已经绑定回宿主机本地地址，后端服务直接访问宿主机回环地址即可。

这也是当前部署结构里一个很重要的判断：

> 这是“容器托底基础服务 + 宿主机运行后端”的混合模式，不是全容器化单体。

## 第一次启动时，迁移到底会做什么

当前后端启动后，会通过迁移工具自动执行：

- `000001_init_schema`
- `000002_seed_data`

这两层分别负责：

| 迁移层 | 主要内容 |
| --- | --- |
| `init_schema` | 创建系统表结构 |
| `seed_data` | 初始化角色、菜单、按钮、Casbin 规则、角色菜单绑定等种子数据 |

这意味着很多“后台最初就该有的结构”不是靠手工点界面创建，而是在启动时就自动建好。

## 超级管理员角色是怎样来的

当前种子数据会自动准备：

- `super_admin` 角色
- 固定的系统菜单树
- 对应的 Casbin 接口权限规则

但它不会直接创建管理员用户。

原因是：

- 用户密码需要 bcrypt 加密

而这一步不适合纯 SQL 种子完成。

## 管理员账号为什么要靠 `/setup/init`

当前管理员账号通过：

- `POST /api/v1/setup/init`

创建。

这个接口会做的事情是：

1. 先检查系统里是否已经存在用户
2. 如果没有，创建 `admin`
3. 对密码做 bcrypt 加密
4. 把这个用户绑定到 `super_admin`

这也解释了一个很常见的现象：

> 数据库里虽然已经有角色、菜单和规则，但你第一次部署后仍然还不能直接登录。

因为管理员用户本身还没有被创建出来。

## `setup-server.sh` 在这里帮你补了什么

当前服务器首次部署脚本里，已经把管理员初始化纳入了流程：

```text
启动 compose.server.yml
  ↓
等待数据库就绪
  ↓
启动 ez-admin systemd 服务
  ↓
调用 /api/v1/setup/init
```

这意味着如果你按仓库当前的服务器脚本部署，管理员创建通常不需要你再额外手敲一次 `curl`。

不过理解这条链路仍然很重要，因为后面排障时你会知道问题到底出在：

- 迁移没跑
- 后端没起来
- 还是初始化接口没执行成功

## 一份最小部署前检查清单

正式部署前，建议先快速核对这几项：

1. `deploy/.env.example` 是否已经变成服务器上的 `.env`
2. `EZ_AUTH_JWT_SECRET` 是否已经替换为真实随机值
3. 数据库和 Redis 地址是否符合当前部署方式
4. 你这次部署到底是 PostgreSQL 还是 MySQL
5. 是否清楚管理员账号是靠 `/setup/init` 创建，而不是靠迁移直接插入

## 读完这页后，应该知道什么

如果把这一页压缩成一句话，那就是：

> 当前系统的“初始化”其实分成了三层：配置覆盖、数据库迁移、管理员创建。

只有把这三层都跑通，部署后的后台才算真正可登录、可使用。

## 下一步

- 想继续看真实的打包、上传、启动和更新流程，下一页读 [部署验证与复用说明](./deployment-and-reuse)
- 想先查看部署产物清单，可以读 [部署产物参考](../../reference/deploy-artifacts-reference)
