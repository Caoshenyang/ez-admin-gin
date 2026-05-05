---
title: Compose 与服务运行结构
description: "围绕当前 compose.server.yml、ez-admin.service、setup-server.sh 和 update-server.sh，讲清服务器上的真实运行拓扑。"
---

# Compose 与服务运行结构

如果只看“怎么部署”，很容易把第 9 章理解成一串命令。但真正决定后续维护成本的，其实是：

> 这套系统在服务器上到底是怎样跑起来的。

这一页专门拆开当前真实部署结构：

- `deploy/compose.server.yml`
- `deploy/ez-admin.service`
- `scripts/setup-server.sh`
- `scripts/update-server.sh`

::: tip 🎯 本节目标
读完后，你应该能回答三个问题：

1. 为什么当前后端不是跑在 Docker 里
2. 为什么基础设施用 Compose，而后端用 systemd
3. 为什么首次部署和后续更新要分成两套脚本
:::

## 先看当前真实运行拓扑

现在服务器上的运行结构可以直接画成这样：

```text
宿主机
├─ /opt/ez-admin/server
│  └─ systemd: ez-admin.service
└─ docker compose -f compose.server.yml
   ├─ postgres
   ├─ redis
   └─ nginx
```

这意味着当前项目采用的是一种很务实的混合结构：

- 后端二进制直接跑在宿主机
- PostgreSQL、Redis、Nginx 跑在容器里

这不是“半成品”，而是当前仓库已经明确选择的一条部署路线。

## 为什么后端没有塞进容器

当前仓库把后端放在宿主机运行，有几个很直接的好处：

- systemd 管进程重启比自己写容器编排更轻
- 查看日志可以直接走 `journalctl`
- 更新后端时只需要替换二进制并重启服务
- 和当前打包脚本、部署脚本更自然地对齐

也就是说，这条路线更像“单体后台首版上线”的工程解，而不是一开始就强行追求全容器化。

## 为什么 PostgreSQL、Redis、Nginx 反而进了容器

因为这三类基础设施最适合标准化容器托管：

- PostgreSQL：卷挂载 + 健康检查
- Redis：卷挂载 + 启动参数
- Nginx：挂配置 + 挂静态文件 + 重启快

它们的共性是：

- 运行边界清晰
- 环境依赖固定
- 用 Compose 管理很省心

这也是为什么 `compose.server.yml` 当前只收这三类服务，而没有把 Go 后端也混进来。

## `compose.server.yml` 现在在表达什么

当前 Compose 文件的重点不是“服务很多”，而是把服务器侧最关键的三类运行资源固定下来：

| 服务 | 当前职责 |
| --- | --- |
| `postgres` | 主数据库 |
| `redis` | 缓存与会话相关能力 |
| `nginx` | 静态资源分发 + 反向代理 |

其中两个值得特别注意的设计是：

### 1. PostgreSQL 和 Redis 只监听 `127.0.0.1`

当前映射端口是：

- `127.0.0.1:5432:5432`
- `127.0.0.1:6379:6379`

这说明数据库和 Redis 默认不对公网直接暴露，而是只给宿主机上的后端服务访问。

### 2. Nginx 使用 `host` 网络

当前 Nginx 不是桥接网络，而是：

- `network_mode: host`

这样它就可以直接把 `/api/` 反代到：

- `127.0.0.1:8080`

也就是宿主机上的 Go 后端。

## `ez-admin.service` 为什么值得单独看

当前后端服务文件很短，但信息量很大：

- `WorkingDirectory=/opt/ez-admin`
- `EnvironmentFile=/opt/ez-admin/.env`
- `ExecStart=/opt/ez-admin/server`

这三行基本就定义了当前后端在服务器上的运行基线：

1. 工作目录固定在部署目录
2. 环境变量统一来自 `.env`
3. 启动目标就是打包出来的 `server` 二进制

所以如果你部署后遇到“为什么代码没生效”“为什么配置没吃到”，通常都要先回来看这三点。

## 为什么首次部署要走 `setup-server.sh`

因为首次部署不只是“把服务跑起来”，还包含了一次性的初始化动作：

- 目录创建
- `dist -> web` 整理
- `.env.example -> .env`
- JWT 密钥首次生成
- Compose 启动
- 后端 systemd 启动
- 管理员初始化

这些动作很多都不应该在后续每次更新时重复发生。

所以当前仓库把首次部署单独拆成：

- `setup-server.sh`

是非常合理的。

## 为什么更新要走 `update-server.sh`

后续更新时，真正需要变动的通常只有：

- 前端静态文件
- 后端二进制

而不应该再去反复碰：

- `.env`
- 数据卷
- 管理员初始化
- 基础容器重建

这就是为什么 `update-server.sh` 当前只做两件核心事：

1. 把新 `dist/` 变成 `web/`
2. 重启 `ez-admin` 后端服务

这种“首次部署脚本”和“日常更新脚本”分离的方式，会明显减少误操作。

## 这一套结构最适合什么场景

当前这套运行结构特别适合：

- 单体后台
- 一台云服务器
- 小团队交付
- 首版快速上线
- 想保留容器便利，但又不想过早把所有东西都容器化

换句话说，它的目标不是“云原生炫技”，而是：

> 用足够稳、足够清晰的方式，把一个真实后台系统交付出去。

## 一份最小排障思路

如果部署后出现问题，建议先按这条顺序判断：

1. `docker compose -f /opt/ez-admin/compose.server.yml ps`
2. `sudo systemctl status ez-admin`
3. `sudo journalctl -u ez-admin -f`
4. `curl http://localhost/health`

这样基本就能快速判断问题是在：

- 容器层
- 后端进程层
- 还是 Nginx 入口层

## 下一步

- 想继续看入口层怎样把静态资源、API、上传和 HTTPS 接起来，下一页读 [Nginx 与 HTTPS 入口层](./nginx-and-https)
- 想回到具体部署步骤，读 [部署验证与复用说明](./deployment-and-reuse)
