---
title: 部署排障 FAQ
description: "围绕当前 setup-server.sh、update-server.sh、compose.server.yml、systemd 和 Nginx 配置，整理最常见的部署故障排查顺序。"
---

# 部署排障 FAQ

部署真正让人卡住的，通常不是“不会执行命令”，而是：

> 命令已经跑完了，但页面打不开、服务没起来、登录失败，或者更新后看起来完全没生效。

这一页把当前部署链路里最常见的几类问题集中整理成 FAQ，默认按现在这套真实结构排查：

```text
Nginx
  ↓
systemd 后端服务
  ↓
PostgreSQL / Redis 容器
  ↓
.env / 初始化数据 / 部署目录
```

::: tip 🎯 本节目标
读完后，你应该能先判断问题落在哪一层，再决定是去看 `docker compose`、`systemctl`、Nginx 配置，还是初始化数据。
:::

## 先记住一条最小排障顺序

不管遇到什么现象，建议先按下面顺序看：

1. `docker compose -f /opt/ez-admin/compose.server.yml ps`
2. `sudo systemctl status ez-admin`
3. `curl -I http://localhost/health`
4. `curl -I http://localhost/api/v1/auth/captcha`
5. 浏览器打开首页，再看登录和系统页

这样做的好处是，你会很快知道：

- 是基础容器没起来
- 还是后端没起来
- 还是 Nginx 没把 API 正确代理过去
- 还是前端页面可开但业务链路没通

## Q1：浏览器打不开首页，应该先看哪里？

先不要急着怀疑前端构建，优先检查入口层和静态文件：

1. 确认 `nginx` 容器是否在运行：

```bash
docker compose -f /opt/ez-admin/compose.server.yml ps
```

2. 确认前端静态资源是否已经落到 `/opt/ez-admin/web`
3. 确认 `deploy/nginx/nginx-native.conf` 对应的 `root` 仍然是 `/opt/ez-admin/web`

如果容器在跑，但首页仍然打不开，最常见的原因通常是：

- `dist/` 没有被整理成 `web/`
- `nginx-native.conf` 没被放进 `nginx/` 目录
- 80 或 443 端口被别的服务占用

## Q2：首页能打开，但接口全部 404 或 502，应该怎么判断？

这通常说明：

- 前端静态文件正常
- Nginx 到后端这一层有问题

优先看两件事：

1. 后端服务是否真的起来了：

```bash
sudo systemctl status ez-admin
```

2. 本机直连后端是否通：

```bash
curl -I http://127.0.0.1:8080/health
```

如果这里不通，问题在后端或其配置，不在 Nginx。

如果这里能通，但浏览器访问 `/api/*` 还是报错，就去看 `nginx-native.conf` 或 `nginx-native-ssl.conf` 里的：

- `location /api/`
- `proxy_pass http://127.0.0.1:8080;`

## Q3：`setup-server.sh` 跑完了，但管理员没有创建成功，怎么办？

当前脚本最后会调用：

- `POST /api/v1/setup/init`

如果没有成功，先看后端健康检查是否已经通过：

```bash
curl -I http://localhost/health
```

如果健康检查正常，再手动执行一次初始化：

```bash
curl -X POST http://localhost/api/v1/setup/init \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"Admin@123456","nickname":"管理员"}'
```

几种常见返回可以这样理解：

- `200`：管理员已创建
- `409`：管理员已经存在
- 其他错误：继续看 `journalctl -u ez-admin -f`

## Q4：数据库或 Redis 容器没起来，先看什么？

当前 `compose.server.yml` 最关键的几个点是：

- PostgreSQL 只监听 `127.0.0.1:5432`
- Redis 只监听 `127.0.0.1:6379`
- 数据目录在 `/opt/ez-admin/data/*`

所以先看：

```bash
docker compose -f /opt/ez-admin/compose.server.yml ps
docker compose -f /opt/ez-admin/compose.server.yml logs postgres
docker compose -f /opt/ez-admin/compose.server.yml logs redis
```

如果容器反复重启，最常见的方向通常是：

- `.env` 里的数据库用户名、密码或库名有问题
- 数据目录权限不对
- 服务器磁盘空间不足
- Redis 密码改了，但应用配置没同步改

## Q5：后端服务起不来，最先看哪条日志？

直接看：

```bash
sudo journalctl -u ez-admin -f
```

当前后端不是跑在 Docker 里，而是跑在 `systemd` 里，所以看到接口报错时，不要先去找 Nginx 日志或 Compose 日志。

如果服务起不来，最常见的原因通常在：

- `.env` 配置缺失或格式错误
- 数据库/Redis 连不上
- 二进制没有执行权限
- 端口已被占用

## Q6：更新后页面还是旧的，像是完全没生效，怎么办？

先区分是前端没更新，还是后端没更新。

### 前端没更新时

重点看两处：

- `/opt/ez-admin/web`
- 浏览器缓存

当前 Nginx 对 `/assets/` 配了长缓存，如果你已经替换了新静态文件，但浏览器还是老页面，先强刷浏览器再判断。

### 后端没更新时

重点看两处：

- `/opt/ez-admin/server` 是否已经被新包覆盖
- `sudo systemctl restart ez-admin` 后服务是否真的重启成功

可以重新执行：

```bash
sudo bash /opt/ez-admin/update-server.sh
sudo systemctl status ez-admin
```

## Q7：刷新前端路由后直接 404，问题通常出在哪？

这通常不是前端页面坏了，而是 SPA 历史路由回退没配好。

当前 Nginx 配置里，首页路由依赖这一段：

```nginx
location / {
    try_files $uri $uri/ /index.html;
}
```

如果这段被改掉，用户直接访问 `/dashboard` 或刷新系统页时，Nginx 就会按真实文件路径去找，最后返回 404。

## Q8：HTTPS 配好了证书，站点还是打不开，先查什么？

优先看三件事：

1. 服务器当前是不是已经切到 `nginx-native-ssl.conf`
2. `/opt/ez-admin/ssl/cert.pem` 和 `/opt/ez-admin/ssl/key.pem` 是否真实存在
3. Nginx 重载前有没有先执行 `nginx -t`

当前 SSL 配置默认围绕 Cloudflare 源站证书写，所以如果你不是走这条证书链路，也要确认自己的证书路径和文件名是否一致。

## Q9：登录页能打开，但验证码或登录接口报错，先查哪一层？

先查接口链路，不要先改前端页面。

建议按这个顺序判断：

1. `curl -I http://localhost/api/v1/auth/captcha`
2. `sudo systemctl status ez-admin`
3. `sudo journalctl -u ez-admin -f`

如果验证码接口都不通，问题通常在后端服务、Nginx 代理或初始化阶段，不在登录页组件本身。

## Q10：上传文件失败，应该先怀疑哪一层？

文件上传最常见的故障点有三层：

1. Nginx `client_max_body_size`
2. 后端上传目录权限
3. 前后端接口链路

当前 `/api/` 代理已经设置：

- `client_max_body_size 20m`

如果你上传的文件明显超过这个值，Nginx 会先拦下来。否则就继续看后端日志，确认是不是文件落盘或元数据写入失败。

## 一份当前结构下的快速判断表

| 现象 | 先看哪里 |
| --- | --- |
| 首页打不开 | `docker compose ps`、`/opt/ez-admin/web`、Nginx `root` |
| 首页能开但 API 404/502 | `systemctl status ez-admin`、`proxy_pass` |
| 管理员没创建 | `/health`、`/api/v1/setup/init`、`journalctl` |
| 数据库/Redis 不通 | `docker compose logs`、`.env`、数据目录 |
| 更新后没生效 | `/opt/ez-admin/server`、`/opt/ez-admin/web`、浏览器缓存 |
| 刷新系统页 404 | `try_files ... /index.html` |
| HTTPS 不通 | 证书文件、SSL 配置、`nginx -t` |

## 下一步

- 想回到完整部署主线，读 [部署验证与复用说明](./deployment-and-reuse)
- 想继续看更新和回滚，读 [更新与回滚策略](./update-and-rollback)
- 想细查 Nginx 配置字段，读 [Nginx 配置参考](../../reference/nginx-config-reference)
