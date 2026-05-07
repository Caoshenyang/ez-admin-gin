---
title: 部署变体说明
description: "围绕 compose.server.yml、compose.deploy.yml、compose.prod.yml 和四份 Nginx 配置，讲清当前底座支持的几种部署形态以及选型方法。"
---

# 部署变体说明

本章前面几页已经把当前主线部署讲清楚了，但仓库里其实不只一套部署形态。

除了现在默认使用的：

- `pack.sh + setup-server.sh + compose.server.yml`

仓库里还保留了：

- 从 Docker Hub 拉镜像的全容器部署
- 从源码直接构建容器的全容器部署

这一页专门回答一个很实际的问题：

> 面对不同交付边界，当前底座到底该选哪一种部署形态。

::: tip 🎯 本节目标
读完后，你应该能判断：

1. 当前默认主线为什么是“后端跑宿主机，基础设施跑容器”
2. `compose.deploy.yml` 和 `compose.prod.yml` 分别适合什么场景
3. `nginx-native*.conf` 和 `nginx*.conf` 应该怎么配套选择
:::

## 先看仓库里现有的部署文件

当前和部署变体直接相关的文件主要有：

```text
deploy/compose.server.yml
deploy/compose.deploy.yml
deploy/compose.prod.yml
deploy/nginx/nginx-native.conf
deploy/nginx/nginx-native-ssl.conf
deploy/nginx/nginx.conf
deploy/nginx/nginx-ssl.conf
scripts/pack.sh
scripts/setup-server.sh
scripts/update-server.sh
```

可以把它们分成三种部署形态来理解。

## 一张选型表先看结论

| 部署形态 | 关键文件 | 适合什么场景 |
| --- | --- | --- |
| 当前默认主线：宿主机后端 + 容器基础设施 | `compose.server.yml`、`setup-server.sh`、`update-server.sh`、`nginx-native*.conf` | 个人项目、小团队交付、自己管服务器、希望更新和排障更直接 |
| 镜像分发型全容器部署 | `compose.deploy.yml`、`nginx*.conf` | 已经有镜像仓库，希望服务器只负责拉镜像和启动 |
| 源码构建型全容器部署 | `compose.prod.yml`、`nginx*.conf` | 服务器或 CI 可以直接访问源码，希望用容器一次构建和运行全部服务 |

## 形态一：当前默认主线

这一条就是本章前面几页一直在讲的主线：

```text
本地 pack.sh
  ↓
上传部署包
  ↓
setup-server.sh / update-server.sh
  ↓
compose.server.yml
  ↓
nginx-native*.conf
```

它的核心结构是：

- Go 后端二进制跑在宿主机 `systemd`
- PostgreSQL、Redis、Nginx 跑在 Docker Compose
- Nginx 通过 `host` 网络直接反代 `127.0.0.1:8080`

### 这一条为什么是当前默认主线

因为它最符合当前底座的交付目标：

- 打包物清晰
- 更新动作简单
- `journalctl` 和 `docker compose` 的职责边界非常明确
- 需要排障时，不会同时陷进“容器构建 + 容器编排 + 容器内日志”三层复杂度

### 这一条最适合谁

优先适合：

- 单体后台
- 个人项目
- 小团队首版交付
- 希望先把上线和运维成本压低的场景

## 形态二：镜像分发型全容器部署

这一条对应：

- `deploy/compose.deploy.yml`

它的核心思路是：

- PostgreSQL、Redis、Server、Nginx 都在容器里
- 后端和前端镜像从镜像仓库直接拉取
- Nginx 通过 Docker 网络里的 `server:8080` 访问后端

### 这一条最关键的前提

你需要先准备好：

- 镜像仓库
- `DOCKERHUB_USERNAME`
- 已经推送好的 `ez-admin-server:latest`
- 已经推送好的 `ez-admin-web:latest`

也就是说，这条线更适合：

- 已经有 CI/CD 推镜像流程
- 服务器不保留源码
- 想把发版动作收敛成“推镜像 + 拉镜像 + 重启”

### 这一条和当前默认主线最大的差别

最大的差别不在 PostgreSQL 或 Redis，而在：

- 后端不再交给宿主机 `systemd`
- Nginx 不再走 `127.0.0.1:8080`
- 入口层配置要改用 `nginx.conf` 或 `nginx-ssl.conf`

## 形态三：源码构建型全容器部署

这一条对应：

- `deploy/compose.prod.yml`

它和 `compose.deploy.yml` 的结构很像，但区别是：

- 后端镜像现场 `build`
- 前端镜像现场 `build`

也就是说，服务器或 CI 节点上要有：

- 完整源码
- Docker 构建能力

### 这一条更适合什么时候用

比较适合：

- 你正在做一条自建构建链路
- 暂时还没把镜像发布流程标准化
- 想先验证“全容器结构是不是适合自己”

### 这一条的代价是什么

它的代价主要是：

- 构建时间更长
- 服务器或 CI 机器压力更大
- 排障时要多看一层镜像构建过程

所以对当前这个底座来说，它更像：

> 一条可选的全容器运行方案，而不是首选交付方式。

## Nginx 配置要跟部署形态配套

当前仓库里的四份 Nginx 配置，并不是随便挑一份用，而是要和部署形态配套：

| 部署形态 | HTTP 配置 | HTTPS 配置 |
| --- | --- | --- |
| 宿主机后端 + 容器基础设施 | `nginx-native.conf` | `nginx-native-ssl.conf` |
| 全容器部署 | `nginx.conf` | `nginx-ssl.conf` |

它们的最大区别是反代目标不同：

- `nginx-native*.conf` 指向 `127.0.0.1:8080`
- `nginx*.conf` 指向 `server:8080`

所以如果部署形态和 Nginx 配置混用了，最容易出现的现象就是：

- 首页能打开
- 但接口全部 502 或 404

## 怎么做部署形态选择

如果你现在要落一个真实项目，可以直接按下面这套判断：

### 先选当前默认主线，如果你满足下面大多数条件

- 项目是单体后台
- 你自己能登录服务器
- 你更看重可维护和排障简单
- 还不想先铺镜像发布平台

### 选 `compose.deploy.yml`，如果你满足下面大多数条件

- 团队已经有镜像仓库
- 想把发版动作做成镜像推送
- 服务器只负责拉镜像和跑容器
- 后端也希望完全容器化

### 选 `compose.prod.yml`，如果你满足下面大多数条件

- 你要做一条从源码直接构建生产容器的链路
- 服务器或 CI 本身就有足够的 Docker 构建资源
- 你明确想验证“全容器构建 + 全容器运行”的生产形态

## 三种形态各自最小验证看什么

### 当前默认主线

先看：

```bash
docker compose -f /opt/ez-admin/compose.server.yml ps
sudo systemctl status ez-admin
curl -I http://localhost/health
```

### 镜像分发型全容器部署

先看：

```bash
docker compose -f deploy/compose.deploy.yml ps
docker compose -f deploy/compose.deploy.yml logs server
docker compose -f deploy/compose.deploy.yml logs nginx
```

### 源码构建型全容器部署

先看：

```bash
docker compose -f deploy/compose.prod.yml ps
docker compose -f deploy/compose.prod.yml logs server
docker compose -f deploy/compose.prod.yml logs nginx
```

## 当前底座的推荐顺序

如果没有非常明确的团队基础设施前提，当前更推荐按这个顺序选择：

1. 先用默认主线把项目稳定上线
2. 再根据团队交付方式，决定是否切到镜像分发型全容器部署
3. 最后才考虑源码构建型全容器部署

这个顺序的核心原因很简单：

> 先把“能稳定交付”做成，再去追求“更统一的容器化形态”。

## 下一步

- 想回到当前默认主线的执行步骤，读 [部署验证与复用说明](./deployment-and-reuse)
- 想继续看当前主线为什么选 `compose.server.yml`，读 [Compose 与服务运行结构](./compose-and-service-layout)
- 想细查这些部署文件本身，读 [Docker 部署文件参考](../../reference/deploy-artifacts-reference)
