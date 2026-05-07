---
title: Nginx 与 HTTPS 入口层
description: "围绕当前 nginx-native.conf、nginx-native-ssl.conf、nginx.conf 和 nginx-ssl.conf，讲清静态资源、反代、上传和 HTTPS。"
---

# Nginx 与 HTTPS 入口层

本章里最容易被当成”配完就忘”的一层，其实是 Nginx。

但当前仓库里，Nginx 不只是一个可有可无的转发器，它承担了四件很核心的事：

- 前端静态资源入口
- API 反向代理
- 上传文件代理
- HTTPS 入口

这一页就对照当前真实配置，把这层讲清楚。

::: tip 🎯 本节目标
读完后，你应该能看懂：

1. 为什么前端页面能直接刷新而不 404
2. 为什么 `/api/` 和 `/uploads/` 都要反代
3. 当前仓库为什么同时保留 native 与非 native 两套 Nginx 配置
:::

## 先看当前真实配置文件

现在仓库里和 Nginx 直接相关的配置主要有四个：

```text
deploy/nginx/nginx-native.conf
deploy/nginx/nginx-native-ssl.conf
deploy/nginx/nginx.conf
deploy/nginx/nginx-ssl.conf
```

它们可以分成两组理解：

| 文件 | 面向哪种结构 |
| --- | --- |
| `nginx-native*.conf` | 当前服务器真实主线：后端跑宿主机 |
| `nginx*.conf` | 保留给“后端也进容器”或更传统容器结构参考 |

所以如果你跟着当前本章主线部署，最应该优先关注的是：

- `nginx-native.conf`
- `nginx-native-ssl.conf`

## 为什么当前主线是 `nginx-native.conf`

因为当前服务器部署的核心前提是：

- 后端服务直接监听宿主机 `127.0.0.1:8080`

所以 Nginx 配置里的 `proxy_pass` 当前也明确指向：

- `http://127.0.0.1:8080`

这和前一页讲的 `host` 网络 Nginx 容器正好对上。

## 当前 Nginx 配置最重要的四个 `location`

如果只记最关键的部分，可以记这四段：

| 路径 | 当前作用 |
| --- | --- |
| `/` | 前端 SPA 入口与历史路由回退 |
| `/assets/` | 前端构建资源长缓存 |
| `/api/` | 后端接口反向代理 |
| `/uploads/` | 上传文件访问代理 |

这四段加起来，基本就定义了当前后台“用户访问一个页面”时会经历的全部入口层行为。

## 为什么 `/` 一定要 `try_files ... /index.html`

当前前端使用的是单页应用路由。

这意味着当用户直接刷新：

- `/system/users`
- `/system/roles`

这类地址时，Nginx 不能把它当成服务器上的真实静态目录去找，否则就会 404。

所以当前配置里最关键的一段之一就是：

- `try_files $uri $uri/ /index.html;`

它的意义是：

> 所有不是实际文件的前端路由，都回退给前端应用自己接管。

## 为什么 `/assets/` 要单独长缓存

当前前端构建产物文件名带哈希，因此：

- 可以安全长缓存

这就是为什么 Nginx 里会单独给 `/assets/` 配：

- `expires 1y`
- `Cache-Control: public, immutable`

这样用户访问后台时，静态资源加载会更稳，也更省带宽。

## 为什么 `/api/` 不是直接开放后端端口给公网

当前主线有一个非常明确的取舍：

- 后端不直接对公网暴露
- 由 Nginx 统一承接公网入口

所以用户访问接口时，其实是：

```text
浏览器
  ↓
Nginx
  ↓
127.0.0.1:8080
```

这带来的好处包括：

- 后端服务可以继续只监听内网或回环地址
- 统一在 Nginx 层做 HTTPS、Header 和请求体限制
- API 和前端页面共享同一域名入口

## `/uploads/` 为什么也要走反代

当前上传文件访问并不是由 Nginx 直接去读某个静态目录，而是继续代理到后端。

这意味着当前系统把“上传资源访问”仍然视为后端资源入口的一部分，而不是完全独立的静态资源站。

这样做的好处是结构简单，但也意味着：

- 上传路径规则要和后端保持一致

所以你后面如果要拆独立对象存储或 CDN，这里就是很关键的演进点。

## HTTPS 当前是怎样切进去的

仓库已经给了 SSL 版本配置：

- `nginx-native-ssl.conf`

它做的事情主要有两层：

1. `80 -> 443` 跳转
2. `443 ssl` 主服务

并通过：

- `/opt/ez-admin/ssl/cert.pem`
- `/opt/ez-admin/ssl/key.pem`

读取证书。

这说明当前本章已经不是”以后再讲 HTTPS”，而是部署结构里已经把它预留好了。

## 为什么文档里总提 Cloudflare

因为当前 SSL 配置和部署说明，默认就是围绕 Cloudflare 源站证书写的。

这不是强制要求你一定用 Cloudflare，而是当前仓库已经给出了一条经过整理的、可直接复用的 HTTPS 主线。

如果你换别的证书来源，也主要还是替换：

- 证书文件路径
- 证书内容

而不是整个 Nginx 入口逻辑都要推倒重来。

## 安全头为什么也值得保留

当前 Nginx 配置里默认已经带了几组基础安全头，例如：

- `X-Frame-Options`
- `X-Content-Type-Options`
- `X-XSS-Protection`

SSL 版本里还多了：

- `Strict-Transport-Security`

它们不是完整安全方案，但对于后台系统来说，这层基础保护值得默认保留。

## `nginx.conf` / `nginx-ssl.conf` 为什么还没删

因为仓库并没有把“后端进容器”这条结构完全禁止掉，而是保留了一份参考配置。

这类文件现在更适合被理解成：

- 结构演进备用
- 其他部署模式参考

而不是当前本章默认主线。

## 一份入口层最小排障清单

如果用户反馈“页面打不开”或“接口不通”，可以先按下面顺序看：

1. `Nginx` 配置是否使用了当前主线对应的 native 版本
2. 前端静态文件是否真的在 `/opt/ez-admin/web`
3. `127.0.0.1:8080` 上的后端是否正常
4. `/api/`、`/uploads/`、`/health` 的反代是否都指向了对的目标
5. SSL 证书路径是否和配置一致

## 下一步

- 想继续看整套服务器运行结构，再回到 [Compose 与服务运行结构](./compose-and-service-layout)
- 想继续看部署执行顺序和更新流程，读 [部署验证与复用说明](./deployment-and-reuse)
