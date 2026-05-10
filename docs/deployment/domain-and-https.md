---
title: 域名与 HTTPS
description: 域名解析、Nginx HTTPS 配置、证书管理和反向代理规则
---

# 域名与 HTTPS

本文介绍如何为 EZ Admin 配置域名和 HTTPS。内容适用于服务器二进制部署方式；Docker 全容器化部署的 Nginx 配置逻辑相同，区别仅在于配置文件的放置路径。

::: tip 前置条件
- 已完成 [服务器二进制部署](/deployment/server-binary-deploy) 或 [Docker 部署](/deployment/docker-deploy)
- 拥有一个域名，且能管理其 DNS 记录
:::

## 请求转发关系

```
Browser
  │
  ▼ :80 / :443
Nginx
  ├── /            → 前端静态文件（SPA，history 模式 fallback 到 index.html）
  ├── /assets/*    → 带哈希的构建文件（缓存 1 年）
  ├── /api/*       → 反向代理到 Go 后端 :8080
  ├── /uploads/*   → 反向代理到 Go 后端 :8080
  └── /health      → 反向代理到 Go 后端 :8080
```

所有 Nginx 配置变体的路由规则一致，区别仅在于是否包含 SSL 终止。

## 域名解析

在域名服务商的 DNS 管理中添加 A 记录：

| 记录类型 | 主机记录 | 记录值 | TTL |
|---------|---------|--------|-----|
| A | `admin`（或 `@`） | 你的服务器公网 IP | 600 |

例如，将 `admin.example.com` 指向服务器 IP：

```
admin.example.com  →  A  →  1.2.3.4
```

解析生效通常需要几分钟。验证：

```bash
ping admin.example.com
# 应返回你的服务器 IP
```

::: tip
本文不绑定特定云厂商。无论你使用哪家的 DNS 服务，操作方式相同：添加 A 记录指向服务器 IP。
:::

## HTTP 配置

服务器二进制部署默认使用 `deploy/nginx/nginx-native.conf`，核心配置如下：

### 前端静态资源

```nginx
server {
    listen 80;
    server_name _;

    # 前端静态资源目录
    root /opt/ez-admin/web;
    index index.html;

    # SPA fallback：所有非文件请求回退到 index.html
    location / {
        try_files $uri $uri/ /index.html;
    }

    # 带哈希的构建文件，缓存 1 年
    location /assets/ {
        expires 1y;
        add_header Cache-Control "public, immutable";
    }
}
```

### API 反向代理

```nginx
    # API 请求代理到后端
    location /api/ {
        proxy_pass http://127.0.0.1:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        client_max_body_size 20m;
    }
```

### 上传文件和健康检查

```nginx
    # 上传文件代理
    location /uploads/ {
        proxy_pass http://127.0.0.1:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }

    # 健康检查
    location /health {
        proxy_pass http://127.0.0.1:8080;
    }
```

### 安全响应头

```nginx
    add_header X-Frame-Options "SAMEORIGIN" always;
    add_header X-Content-Type-Options "nosniff" always;
    add_header X-XSS-Protection "1; mode=block" always;
```

## HTTPS 配置

使用 `deploy/nginx/nginx-native-ssl.conf` 配置 HTTPS。

### HTTPS 证书方案

| 方案 | 适合场景 | 有效期 |
|------|---------|--------|
| Let's Encrypt（免费，自动续期） | 大多数场景 | 90 天，Certbot 自动续期 |
| Cloudflare 源站证书（免费） | 使用 Cloudflare DNS 的场景 | 最长 15 年 |
| 自有证书 | 已有证书（如从域名服务商获取） | 取决于证书 |

::: tip
项目内置的 SSL 配置示例路径为 `/opt/ez-admin/ssl/cert.pem` 和 `/opt/ez-admin/ssl/key.pem`，可以适配 Cloudflare Origin Certificate。如果使用 Let's Encrypt，只需保持证书路径一致，或修改 Nginx 配置中的 `ssl_certificate` / `ssl_certificate_key`。
:::

### 1. 获取并放置证书

```bash
mkdir -p /opt/ez-admin/ssl
cp cert.pem /opt/ez-admin/ssl/cert.pem
cp key.pem /opt/ez-admin/ssl/key.pem
```

::: details 各方案的证书获取方式
- **Let's Encrypt**：`sudo certbot certonly --standalone -d admin.example.com`，证书位于 `/etc/letsencrypt/live/admin.example.com/`
- **Cloudflare 源站证书**：Cloudflare 控制台 → SSL/TLS → 源服务器 → 创建证书，下载 `cert.pem` 和 `key.pem`
- **自有证书**：直接使用服务商提供的证书文件
:::

### 2. 切换 Nginx 配置

```bash
cd /opt/ez-admin/nginx

# 备份 HTTP 配置
cp nginx-native.conf nginx-native.conf.bak

# 替换为 HTTPS 配置
cp ../nginx-native-ssl.conf nginx-native.conf
```

### 3. 修改 server_name

编辑 `nginx-native.conf`，将 `server_name _;` 改为你的域名：

```nginx
# 修改前
server_name _;

# 修改后
server_name admin.example.com;
```

HTTP 配置中的 `server_name` 也建议同步修改。

### 4. 重启 Nginx

```bash
cd /opt/ez-admin
docker compose -f compose.server.yml restart nginx
```

### 5. 验证

```bash
curl -I https://admin.example.com
# 应返回 200，且包含 Strict-Transport-Security 头
```

### HTTPS 额外配置

SSL 版本相比 HTTP 版本增加了以下内容：

- **HTTP → HTTPS 跳转**：80 端口自动 301 重定向到 443
- **SSL 配置**：`ssl_certificate`、`ssl_certificate_key`、`ssl_protocols TLSv1.2 TLSv1.3`
- **HSTS 头**：`Strict-Transport-Security: max-age=31536000; includeSubDomains`

## Docker 全容器化部署：配置 HTTPS

Docker 部署使用 `nginx-ssl.conf`：

### 1. 放置证书和配置

```bash
mkdir -p /opt/ez-admin/nginx/ssl

# 放置证书
cp cert.pem /opt/ez-admin/nginx/ssl/cert.pem
cp key.pem /opt/ez-admin/nginx/ssl/key.pem

# 复制 HTTPS 配置
cp nginx-ssl.conf /opt/ez-admin/nginx/nginx-ssl.conf
```

### 2. 修改 server_name

编辑 `nginx-ssl.conf`，将 `server_name _;` 改为你的域名。

### 3. 切换配置

编辑 `/opt/ez-admin/.env`：

```bash
EZ_NGINX_CONF=nginx-ssl.conf
```

### 4. 重启 Nginx

```bash
docker compose -f compose.deploy.yml restart nginx
```

## Cloudflare 配置参考

如果你使用 Cloudflare 作为 DNS 和 CDN：

1. **SSL/TLS 模式**：设为 `Full (Strict)`，要求有效源站证书
2. **创建源站证书**：SSL/TLS → 源服务器 → 创建证书 → 下载 `cert.pem` 和 `key.pem`
3. **Always Use HTTPS**：开启，Cloudflare 会自动将 HTTP 重定向到 HTTPS
4. **Minimum TLS Version**：建议 TLS 1.2

::: details Cloudflare Full vs Full (Strict) 的区别
- **Full**：Cloudflare 到源站使用 HTTPS，但不验证源站证书。源站甚至可以使用自签证书。
- **Full (Strict)**：Cloudflare 到源站使用 HTTPS，并验证源站证书有效性。需要使用 Cloudflare 源站证书或受信任 CA 签发的证书。

生产环境建议使用 `Full (Strict)`。
:::

::: details 源站证书有效期
Cloudflare 源站证书最长有效期 15 年。到期前需要重新创建并替换服务器上的证书文件。

Let's Encrypt 证书有效期 90 天，Certbot 会自动续期。两种方案各有利弊，按需选择。
:::

## 常见问题

### 域名能访问但 HTTPS 不行

- 确认证书文件路径正确：`ls -la /opt/ez-admin/ssl/cert.pem /opt/ez-admin/ssl/key.pem`
- 确认 Nginx 已使用 SSL 配置：`docker exec ez-admin-nginx cat /etc/nginx/conf.d/default.conf | grep ssl`
- 确认 443 端口已开放：`sudo ss -tlnp | grep 443`

### 证书路径错误

Nginx SSL 配置中的证书路径必须与实际放置位置一致：

```nginx
ssl_certificate     /opt/ez-admin/ssl/cert.pem;    # 确认文件存在
ssl_certificate_key /opt/ez-admin/ssl/key.pem;      # 确认文件存在
```

如果使用 Let's Encrypt，证书路径通常是 `/etc/letsencrypt/live/your-domain/fullchain.pem`。

### Mixed Content 报错

确保前端 API 请求地址使用 `https://` 而非 `http://`。检查浏览器开发者工具 Network 面板中的请求 URL。

如果使用 Cloudflare，开启 "Always Use HTTPS" 可以自动解决这个问题。

### /api 返回 502

后端未启动或端口不匹配：

```bash
# 二进制部署
sudo systemctl status ez-admin
curl http://localhost:8080/health

# Docker 部署
docker compose -f compose.deploy.yml ps server
docker compose -f compose.deploy.yml logs server --tail 20
```

### 前端刷新 404

确认 Nginx 配置中包含 SPA fallback：

```nginx
location / {
    try_files $uri $uri/ /index.html;
}
```

检查 Nginx 配置是否正确挂载到容器中。
