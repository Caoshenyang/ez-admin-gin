---
title: Nginx 配置
description: "配置 Nginx 实现前端 SPA 托管、API 反向代理、静态资源缓存、gzip 压缩和安全头。"
---

# Nginx 配置

上一节在 Compose 里看到 Nginx 作为唯一对外入口，负责把前端页面和 API 请求分别转发到正确的地方。这一节拆开这份配置，看看每一块做了什么。

::: tip 🎯 本节目标
理解 `nginx.conf` 里每个 location 块的用途，知道前端路由、API 代理、缓存、gzip 和安全头分别是怎么配置的。
:::

## 完整配置

<<< ../../../deploy/nginx/nginx.conf

## 前端静态资源托管

```nginx
root /usr/share/nginx/html;
index index.html;
```

`root` 指向 Nginx 容器内前端构建产物的目录。`admin/Dockerfile` 在构建阶段执行 `pnpm build`，把打包结果复制到这个路径。`index` 指定默认首页文件。

## SPA 历史模式路由

```nginx
location / {
    try_files $uri $uri/ /index.html;
}
```

这一行是 SPA 部署的核心配置，解决的是 Vue Router 的 `history` 模式问题：

1. 用户访问 `/dashboard` 时，Nginx 先查找是否存在 `/dashboard` 这个文件。
2. 不存在，再查找 `/dashboard/` 这个目录。
3. 也不存在，最终返回 `index.html`。

浏览器拿到 `index.html` 后，Vue Router 接管路由，解析 `/dashboard` 并渲染对应页面。

::: warning ⚠️ 去掉 try_files 会怎样
如果没有 `try_files` 回退，用户直接访问 `/dashboard` 或刷新页面时，Nginx 会返回 404——因为它真的去找 `/dashboard` 这个文件了，找不到就报错。这是 SPA 部署最常见的踩坑点。
:::

## 静态资源长缓存

```nginx
location /assets/ {
    expires 1y;
    add_header Cache-Control "public, immutable";
}
```

Vite 构建前端时，`/assets/` 下的文件名会带上内容哈希（如 `dashboard-3a8f2b1c.js`）。文件内容不变，文件名就不变；内容一变，文件名也跟着变。

利用这个特性，可以放心设置 1 年长缓存：

- `expires 1y` 设置 `Expires` 头为一年后。
- `Cache-Control: public, immutable` 告诉浏览器和 CDN：这个资源永远不会变，不需要重新验证。
- 当代码更新后，文件名变了，浏览器自然请求新文件，不会读到旧缓存。

## API 反向代理

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

所有 `/api/` 开头的请求被转发到后端服务：

- `proxy_pass http://server:8080` — `server` 是 Compose 里的服务名，Nginx 容器通过 Docker 内部网络直接访问后端，不需要经过宿主机端口。
- `X-Real-IP` / `X-Forwarded-For` / `X-Forwarded-Proto` — 把客户端真实 IP 和协议传给后端，否则后端拿到的都是 Nginx 容器的内部 IP。
- `client_max_body_size 20m` — 允许上传最大 20MB 的文件，覆盖 Nginx 默认的 1MB 限制。

::: details 为什么用 server:8080 而不是 localhost
Nginx 和后端服务跑在不同的容器里。在 Nginx 容器内部，`localhost` 指的是 Nginx 自己，不是后端。Docker Compose 的服务名（`server`）会被解析为后端容器的内部 IP，所以用 `http://server:8080` 才能正确转发。
:::

## 上传文件代理

```nginx
location /uploads/ {
    proxy_pass http://server:8080;
    proxy_set_header Host $host;
    proxy_set_header X-Real-IP $remote_addr;
}
```

用户上传的文件由后端存储和管理，前端通过 `/uploads/` 路径访问。这里的配置把 `/uploads/` 请求转发到后端，和 API 代理类似，但不设 `client_max_body_size`——上传大小限制在 API 代理那个 location 里已经处理过了。

## 健康检查代理

```nginx
location /health {
    proxy_pass http://server:8080;
}
```

把 `/health` 请求直接转发给后端。外部负载均衡器或监控系统可以通过这个端点检查服务是否存活。配置保持最简——健康检查不需要复杂的 header 转发。

## gzip 压缩

```nginx
gzip on;
gzip_types text/plain text/css application/json application/javascript
           text/xml application/xml application/xml+rss
           text/javascript image/svg+xml;
gzip_min_length 1024;
gzip_vary on;
```

- `gzip on` 开启压缩。
- `gzip_types` 指定需要压缩的 MIME 类型。注意 Nginx 默认只压缩 `text/html`，所以 JS、CSS、JSON、SVG 等必须显式列出。
- `gzip_min_length 1024` 只压缩超过 1KB 的响应，小文件压缩收益不大反而增加 CPU 开销。
- `gzip_vary on` 在响应头里加 `Vary: Accept-Encoding`，让 CDN 和浏览器正确处理缓存。

::: tip 💡 为什么不压缩图片
图片文件（PNG、JPG、WebP）本身已经是压缩格式，再 gzip 压缩几乎不会变小，反而浪费 CPU。所以 `gzip_types` 里没有列图片类型。
:::

## 安全头

```nginx
add_header X-Frame-Options "SAMEORIGIN" always;
add_header X-Content-Type-Options "nosniff" always;
add_header X-XSS-Protection "1; mode=block" always;
```

三个安全头的作用：

| 头 | 作用 |
| --- | --- |
| `X-Frame-Options: SAMEORIGIN` | 只允许同源页面嵌入 iframe，防止点击劫持 |
| `X-Content-Type-Options: nosniff` | 禁止浏览器猜测响应类型，防止 MIME 嗅探攻击 |
| `X-XSS-Protection: 1; mode=block` | 启用浏览器内置 XSS 过滤器，检测到攻击时阻止页面渲染 |

`always` 关键字确保即使响应状态码是 4xx 或 5xx，这些头也会被添加。

::: details 生产环境还可以加什么
如果后续需要 HTTPS，可以补充：

- `Strict-Transport-Security`（HSTS）— 强制浏览器使用 HTTPS。
- `Content-Security-Policy`（CSP）— 限制页面可以加载哪些资源。
- `Referrer-Policy` — 控制请求头中 Referer 的暴露程度。

这些配置通常在 HTTPS 配置完成后统一添加。
:::

## 小结

`nginx.conf` 做了五件事：

- **SPA 路由回退** — `try_files` 确保前端路由不会 404。
- **API 反向代理** — `/api/` 转发到后端，附带真实客户端信息。
- **静态资源缓存** — `/assets/` 利用文件名哈希设置一年长缓存。
- **gzip 压缩** — 减少 JS、CSS、JSON 等文本资源的传输体积。
- **安全头** — 基础的点击劫持、MIME 嗅探和 XSS 防护。

接下来看环境变量和初始化数据的准备：[环境变量与初始化数据](./env-and-init-data)。
