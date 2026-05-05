---
title: 部署验证与复用说明
description: "围绕当前 pack.sh、setup-server.sh、update-server.sh、compose.server.yml 和 ez-admin.service，讲清打包、部署、更新和复用。"
---

# 部署验证与复用说明

当前仓库的部署方式已经不是“只给几个零散命令”，而是围绕一条明确的服务器链路组织好的：

```text
本地打包
  ↓
上传部署包
  ↓
首次执行 setup-server.sh
  ↓
后续执行 update-server.sh
```

这一页就沿着这条真实链路来讲。

::: tip 🎯 本节目标
读完后，你应该能顺着当前仓库已有脚本，把项目从本地源码推进到“服务器可访问、可更新、可复用”的状态。
:::

## 先看当前真实部署文件

和部署最直接相关的文件主要有：

```text
scripts/pack.sh
scripts/pack.ps1
scripts/setup-server.sh
scripts/update-server.sh
deploy/compose.server.yml
deploy/ez-admin.service
deploy/nginx/nginx-native.conf
deploy/nginx/nginx-native-ssl.conf
```

它们的职责分工是：

| 文件 | 当前作用 |
| --- | --- |
| `pack.sh / pack.ps1` | 本地打包后端、前端和部署文件 |
| `setup-server.sh` | 首次部署初始化 |
| `update-server.sh` | 后续版本更新 |
| `compose.server.yml` | PostgreSQL、Redis、Nginx 容器 |
| `ez-admin.service` | 后端 systemd 服务 |
| `nginx-native*.conf` | 宿主机网络模式下的 Nginx 反代配置 |

## 第一步：本地打包到底打出了什么

当前 `scripts/pack.sh` 会按下面顺序做事：

1. 编译 Linux `amd64` 后端二进制
2. 构建前端静态资源
3. 复制部署配置与脚本
4. 生成 `deploy-package.tar.gz`

最终打包目录里主要会有：

- `server`
- `dist/`
- `compose.server.yml`
- `.env.example`
- `ez-admin.service`
- `setup-server.sh`
- `update-server.sh`
- `configs/`
- `nginx/`

这说明当前部署包不是“只打后端”或“只打前端”，而是一个已经可以直接传服务器的完整运行包。

## 第二步：首次部署时真正发生了什么

当前服务器首次部署的主入口是：

- `setup-server.sh`

它做的事情比“启动一下服务”更多，真实顺序是：

1. 创建部署目录和数据目录
2. 把 `dist/` 转成 `web/`
3. 把 `nginx-native.conf` 挪到正确位置
4. 第一次把 `.env.example` 变成 `.env`
5. 安装并加载 `ez-admin.service`
6. 如果还是占位值，就自动生成 JWT 密钥
7. 启动 `compose.server.yml`
8. 等数据库就绪
9. 启动后端 systemd 服务
10. 调用 `/api/v1/setup/init` 创建管理员

也就是说，当前首次部署脚本已经在帮你把：

- 文件整理
- 基础服务启动
- 后端启动
- 管理员初始化

串成一条完整闭环。

## 第三步：当前服务器运行结构是什么样

部署完成后，系统运行结构可以理解成：

```text
宿主机 systemd
  └─ /opt/ez-admin/server

Docker Compose
  ├─ postgres
  ├─ redis
  └─ nginx
```

这里最重要的判断是：

- 后端不是容器运行
- 数据库、缓存、Nginx 才是容器运行

这也解释了为什么：

- 后端日志看 `journalctl`
- 容器状态看 `docker compose ps`

而不是统一只看 Docker。

## 为什么 `compose.server.yml` 里的 Nginx 用 `host` 网络

当前 Nginx 选择的是：

- `network_mode: host`

这样它就能直接用：

- `127.0.0.1:8080`

去反代宿主机上的后端服务，而不用再额外折腾容器到宿主机的特殊网络映射。

这是一种非常务实的服务器部署写法，尤其适合当前这种“后端在宿主机，基础设施在容器”的混合结构。

## 第四步：后续更新为什么不能再跑 `setup-server.sh`

因为 `setup-server.sh` 的定位是：

- 首次部署脚本

它会处理：

- `.env` 初建
- JWT 密钥首次生成
- 管理员初始化

后续更新时，应该走的是：

- `update-server.sh`

当前更新脚本只做两件事：

1. 更新前端静态文件
2. 重启后端服务

这正是它安全的地方，因为它不会乱动你已经存在的：

- 数据卷
- 容器环境
- `.env`
- 管理员数据

## 现在怎样判断部署是否成功

当前部署完成后，最小验证建议按下面顺序做：

1. `docker compose -f /opt/ez-admin/compose.server.yml ps`
2. `sudo systemctl status ez-admin`
3. 浏览器打开服务器 IP 或域名
4. 用管理员账号登录
5. 随便进入一个系统页，确认页面和接口都可用

如果你想更稳一点，再补两步：

6. 打开一个需要按钮权限的页面，确认按钮显隐正常
7. 查看操作日志和登录日志，确认部署后链路真的在工作

## 当前 HTTPS 是怎样接进去的

仓库现在已经给了两类 Nginx 配置：

- 非 SSL：`deploy/nginx/nginx-native.conf`
- SSL：`deploy/nginx/nginx-native-ssl.conf`

这意味着第 9 章当前的真实做法不是临场手写配置，而是：

> 先用非 SSL 跑通，再按证书路径切到 SSL 配置文件。

如果你使用 Cloudflare，这条链路会更顺，因为当前文档和配置默认就是围绕它来写的。

## 复用到新项目时，最值得保留的是什么

第 9 章的复用价值不只是“再部署一次”，而是：

- 打包脚本可复用
- 部署目录结构可复用
- systemd 服务方式可复用
- Docker Compose 基础设施可复用

所以如果后面你基于这个仓库开新项目，最应该保留的往往不是某一条命令，而是这整条部署链路的分层方式。

## 当前这条链路特别适合哪类项目

它非常适合：

- 单体后台
- 小团队交付
- 自己掌控服务器
- 想快速上线但不想一开始就上复杂 K8s 体系

也就是说，这套部署方式的重点不是“炫技”，而是：

> 让一个真实后台底座能更稳地完成首版上线。

## 一份快速部署清单

你可以直接按下面这份清单执行：

1. 本地执行 `bash scripts/pack.sh` 或 `.\scripts\pack.ps1`
2. 上传部署包到服务器
3. 解压到 `/opt/ez-admin`
4. 首次运行 `sudo bash /opt/ez-admin/setup-server.sh`
5. 后续更新运行 `sudo bash /opt/ez-admin/update-server.sh`
6. 用浏览器和日志命令完成最小验收

## 下一步

- 想先补环境变量和迁移初始化这层背景，回到 [环境变量与初始化数据](./env-and-init-data)
- 想继续查更细的部署文件含义，可以读 [部署产物参考](../../reference/deploy-artifacts-reference)
