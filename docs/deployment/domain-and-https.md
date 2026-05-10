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

## HTTPS 证书方案

EZ Admin 的 Nginx 配置支持多种 HTTPS 方案：

### 方案一：Let's Encrypt（免费，自动续期）

适合大多数场景。使用 Certbot 自动获取和续期证书。

### 方案二：Cloudflare 源站证书（免费，需使用 Cloudflare DNS）

EZ Admin 已内置 Cloudflare 源站证书的 Nginx 配置模板（`nginx-native-ssl.conf` 和 `nginx-ssl.conf`）。

### 方案三：自有证书

如果你已有证书（如从域名服务商获取），直接放置即可。

## 服务器二进制部署：配置 HTTPS

服务器二进制部署使用 `nginx-native-ssl.conf`，步骤如下：

### 1. 获取并放置证书

证书和私钥放到 `/opt/ez-admin/ssl/` 目录：

```bash
mkdir -p /opt/ez-admin/ssl
cp cert.pem /opt/ez-admin/ssl/cert.pem
cp key.pem /opt/ez-admin/ssl/key.pem
```

::: details 各方案的证书获取方式
- **Let's Encrypt**：`sudo certbot certonly --standalone -d admin.example.com`，证书位于 `/etc/letsencrypt/live/admin.example.com/`
- **Cloudflare 源站证书**：在 Cloudflare 控制台 → SSL/TLS → 源服务器 → 创建证书，下载 cert.pem 和 key.pem
- **自有证书**：直接使用服务商提供的证书文件
:::

### 2. 切换 Nginx 配置

将 `/opt/ez-admin/nginx/` 下的配置替换为 SSL 版本：

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

### 4. 重启 Nginx

```bash
cd /opt/ez-admin
docker compose -f compose.server.yml restart nginx
```

### 5. 验证

```bash
curl -I https://admin.example.com
# 应返回 200，且包含严格传输安全头
```

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

## Nginx 反向代理规则

EZ Admin 的 Nginx 配置包含以下路由规则，所有变体一致：

| 路径 | 代理目标 | 说明 |
|------|---------|------|
| `/` | 前端静态文件 | SPA，`try_files $uri $uri/ /index.html` |
| `/assets/*` | 前端静态文件 | 带 hash 的构建文件，缓存 1 年 |
| `/api/*` | 后端 `:8080` | API 请求代理 |
| `/uploads/*` | 后端 `:8080` | 上传文件代理 |
| `/health` | 后端 `:8080` | 健康检查端点 |

### SPA Fallback

前端使用 Vue Router 的 history 模式，所有非文件请求都回退到 `index.html`：

```nginx
location / {
    try_files $uri $uri/ /index.html;
}
```

### /api 代理

API 请求转发到后端，并附加标准代理头：

```nginx
location /api/ {
    proxy_pass http://127.0.0.1:8080;  # 二进制部署
    # proxy_pass http://server:8080;   # Docker 部署
    proxy_set_header Host $host;
    proxy_set_header X-Real-IP $remote_addr;
    proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    proxy_set_header X-Forwarded-Proto $scheme;
    client_max_body_size 20m;
}
```

### 安全响应头

所有响应都包含安全头：

```nginx
add_header X-Frame-Options "SAMEORIGIN" always;
add_header X-Content-Type-Options "nosniff" always;
add_header X-XSS-Protection "1; mode=block" always;
```

HTTPS 模式额外添加 HSTS：

```nginx
add_header Strict-Transport-Security "max-age=31536000; includeSubDomains" always;
```

## Cloudflare 配置参考

如果你使用 Cloudflare 作为 DNS 和 CDN：

1. **SSL/TLS 模式**：设为 `Full (Strict)`，要求有效源站证书
2. **创建源站证书**：SSL/TLS → 源服务器 → 创建证书 → 下载 `cert.pem` 和 `key.pem`
3. **Always Use HTTPS**：开启，Cloudflare 会自动将 HTTP 重定向到 HTTPS
4. **Minimum TLS Version**：建议 TLS 1.2

::: details Cloudflare 源站证书有效期
源站证书最长有效期 15 年。到期前需要重新创建并替换服务器上的证书文件。

Let's Encrypt 证书有效期 90 天，Certbot 会自动续期。两种方案各有利弊，按需选择。
:::

## 常见问题

### 浏览器提示证书不信任

- Let's Encrypt：确认 Certbot 续期正常（`sudo certbot renew --dry-run`）
- Cloudflare：确认 SSL 模式为 `Full (Strict)`，源站证书未过期
- 自有证书：确认证书链完整

### HTTP 没有跳转到 HTTPS

Nginx SSL 配置中已包含 HTTP → HTTPS 301 跳转。如果没有跳转，检查是否正确替换了 Nginx 配置文件。

### /api 返回 502

后端未启动或端口不匹配：

```bash
# 检查后端状态（二进制部署）
sudo systemctl status ez-admin

# 检查后端状态（Docker 部署）
docker compose -f compose.deploy.yml ps server
```
