---
title: Nginx 配置参考
description: "逐块解析项目中 Nginx 配置文件的作用：SPA 路由回退、API 反向代理、静态资源缓存、gzip 压缩、安全头和 SSL 配置。"
---

# Nginx 配置参考

这一页解析项目中两份 Nginx 配置文件的每个部分：`deploy/nginx/nginx.conf`（HTTP 模式）和 `deploy/nginx/nginx-ssl.conf`（HTTPS 模式）。

::: tip 这页怎么读
部署时不需要改这些配置文件，照着上传就行。遇到问题需要排查、或者想自定义配置时，再回来查对应的章节。
:::

## 完整配置文件

::: details `deploy/nginx/nginx.conf` — HTTP 模式
```nginx
server {
    listen 80;
    server_name _;

    # 前端静态资源
    root /usr/share/nginx/html;
    index index.html;

    # gzip 压缩
    gzip on;
    gzip_types text/plain text/css application/json application/javascript text/xml application/xml application/xml+rss text/javascript image/svg+xml;
    gzip_min_length 1024;
    gzip_vary on;

    # 前端路由：所有非文件请求回退到 index.html（SPA 历史模式）
    location / {
        try_files $uri $uri/ /index.html;
    }

    # 静态资源长缓存（前端构建文件带 hash）
    location /assets/ {
        expires 1y;
        add_header Cache-Control "public, immutable";
    }

    # API 请求反向代理到后端
    location /api/ {
        proxy_pass http://server:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;

        # 上传文件大小限制
        client_max_body_size 20m;
    }

    # 后端上传文件的静态资源代理
    location /uploads/ {
        proxy_pass http://server:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }

    # 健康检查代理
    location /health {
        proxy_pass http://server:8080;
    }

    # 安全头
    add_header X-Frame-Options "SAMEORIGIN" always;
    add_header X-Content-Type-Options "nosniff" always;
    add_header X-XSS-Protection "1; mode=block" always;
}
```
:::

::: details `deploy/nginx/nginx-ssl.conf` — HTTPS 模式（Cloudflare 源站证书）
```nginx
# HTTP → HTTPS 301 跳转
server {
    listen 80;
    server_name _;
    return 301 https://$host$request_uri;
}

# HTTPS 主服务
server {
    listen 443 ssl;
    server_name _;

    # Cloudflare 源站证书
    ssl_certificate     /etc/nginx/ssl/cert.pem;
    ssl_certificate_key /etc/nginx/ssl/key.pem;
    ssl_protocols       TLSv1.2 TLSv1.3;

    # 前端静态资源
    root /usr/share/nginx/html;
    index index.html;

    # gzip 压缩
    gzip on;
    gzip_types text/plain text/css application/json application/javascript text/xml application/xml application/xml+rss text/javascript image/svg+xml;
    gzip_min_length 1024;
    gzip_vary on;

    # 前端路由：所有非文件请求回退到 index.html（SPA 历史模式）
    location / {
        try_files $uri $uri/ /index.html;
    }

    # 静态资源长缓存（前端构建文件带 hash）
    location /assets/ {
        expires 1y;
        add_header Cache-Control "public, immutable";
    }

    # API 请求反向代理到后端
    location /api/ {
        proxy_pass http://server:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;

        # 上传文件大小限制
        client_max_body_size 20m;
    }

    # 后端上传文件的静态资源代理
    location /uploads/ {
        proxy_pass http://server:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }

    # 健康检查代理
    location /health {
        proxy_pass http://server:8080;
    }

    # 安全头 + HSTS
    add_header X-Frame-Options "SAMEORIGIN" always;
    add_header X-Content-Type-Options "nosniff" always;
    add_header X-XSS-Protection "1; mode=block" always;
    add_header Strict-Transport-Security "max-age=31536000; includeSubDomains" always;
}
```
:::

## 配置块解析

### SPA 路由回退

```nginx
location / {
    try_files $uri $uri/ /index.html;
}
```

前端使用 Vue Router 的 `history` 模式。用户访问 `/dashboard` 时：

1. Nginx 先查找 `/dashboard` 文件 → 不存在
2. 再查找 `/dashboard/` 目录 → 不存在
3. 返回 `index.html` → Vue Router 接管路由，渲染 `/dashboard` 页面

::: warning 去掉 try_files 会怎样
没有这个回退，用户直接访问 `/dashboard` 或刷新页面时，Nginx 返回 404——因为它真的去找 `/dashboard` 文件了，找不到就报错。这是 SPA 部署最常见的踩坑点。
:::

### 静态资源长缓存

```nginx
location /assets/ {
    expires 1y;
    add_header Cache-Control "public, immutable";
}
```

Vite 构建时 `/assets/` 下的文件名带内容哈希（如 `dashboard-3a8f2b1c.js`）。内容变则文件名变，可以放心设置 1 年缓存：

- `expires 1y`：`Expires` 头设为一年后。
- `Cache-Control: public, immutable`：告诉浏览器和 CDN 这个资源不会变，不需要重新验证。
- 代码更新后文件名变了，浏览器自动请求新文件。

### API 反向代理

```nginx
location /api/ {
    proxy_pass http://server:8080;
    proxy_set_header Host $host;
    proxy_set_header X-Real-IP $remote_addr;
    proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    proxy_set_header X-Forwarded-Proto $scheme;
    client_max_body_size 20m;
}
```

| 指令 | 作用 |
| --- | --- |
| `proxy_pass http://server:8080` | 转发到后端容器。`server` 是 Docker Compose 服务名，不是 localhost |
| `X-Real-IP` / `X-Forwarded-For` | 透传客户端真实 IP，否则后端拿到的都是 Nginx 容器的内部 IP |
| `X-Forwarded-Proto` | 告诉后端原始请求协议（HTTP 或 HTTPS） |
| `client_max_body_size 20m` | 允许上传最大 20MB 文件，覆盖 Nginx 默认的 1MB 限制 |

::: details 为什么用 server:8080 而不是 localhost
Nginx 和后端跑在不同容器里。Nginx 容器内的 `localhost` 指的是 Nginx 自己，不是后端。Docker Compose 的服务名 `server` 会被解析为后端容器的内部 IP。
:::

### 上传文件代理

```nginx
location /uploads/ {
    proxy_pass http://server:8080;
    proxy_set_header Host $host;
    proxy_set_header X-Real-IP $remote_addr;
}
```

用户上传的文件由后端存储和管理，前端通过 `/uploads/` 路径访问。上传大小限制在 `/api/` 的 `client_max_body_size` 中已设置，这里不需要重复。

### 健康检查代理

```nginx
location /health {
    proxy_pass http://server:8080;
}
```

外部负载均衡器或监控系统通过 `/health` 检查服务是否存活。配置保持最简，不需要 header 转发。

### gzip 压缩

```nginx
gzip on;
gzip_types text/plain text/css application/json application/javascript
           text/xml application/xml application/xml+rss
           text/javascript image/svg+xml;
gzip_min_length 1024;
gzip_vary on;
```

| 指令 | 作用 |
| --- | --- |
| `gzip on` | 开启压缩 |
| `gzip_types` | 需要压缩的 MIME 类型。Nginx 默认只压 `text/html`，JS/CSS/JSON/SVG 等必须显式列出 |
| `gzip_min_length 1024` | 只压缩超过 1KB 的响应，小文件压缩收益不大反而增加 CPU 开销 |
| `gzip_vary on` | 加 `Vary: Accept-Encoding` 头，让 CDN 和浏览器正确处理缓存 |

> [!NOTE]
> 图片（PNG、JPG、WebP）本身已是压缩格式，再 gzip 不会变小，所以 `gzip_types` 里没有列图片类型。

### 安全头

```nginx
add_header X-Frame-Options "SAMEORIGIN" always;
add_header X-Content-Type-Options "nosniff" always;
add_header X-XSS-Protection "1; mode=block" always;
```

| 头 | 作用 |
| --- | --- |
| `X-Frame-Options: SAMEORIGIN` | 只允许同源页面嵌入 iframe，防止点击劫持 |
| `X-Content-Type-Options: nosniff` | 禁止浏览器猜测响应类型，防止 MIME 嗅探攻击 |
| `X-XSS-Protection: 1; mode=block` | 启用浏览器内置 XSS 过滤器 |

`always` 确保即使响应状态码是 4xx 或 5xx，这些头也会被添加。

## HTTPS 模式新增内容

`nginx-ssl.conf` 在 `nginx.conf` 基础上增加了三部分：

### HTTP → HTTPS 跳转

```nginx
server {
    listen 80;
    server_name _;
    return 301 https://$host$request_uri;
}
```

所有 HTTP 请求被 301 重定向到 HTTPS。

### SSL 配置

```nginx
listen 443 ssl;
ssl_certificate     /etc/nginx/ssl/cert.pem;
ssl_certificate_key /etc/nginx/ssl/key.pem;
ssl_protocols       TLSv1.2 TLSv1.3;
```

加载 Cloudflare 源站证书。证书文件通过 Docker 卷挂载到 `/etc/nginx/ssl/` 目录。

### HSTS

```nginx
add_header Strict-Transport-Security "max-age=31536000; includeSubDomains" always;
```

告诉浏览器一年内始终使用 HTTPS 访问这个域名，即使用户输入 `http://` 也会自动跳转。
