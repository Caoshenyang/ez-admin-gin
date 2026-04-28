---
title: 部署验证与复用说明
description: "本地编译后端二进制、构建前端静态文件，部署到腾讯云轻量服务器，通过 Cloudflare 配置 HTTPS 域名访问。"
---

# 部署验证与复用说明

这一节把后台底座部署到公网。思路很简单：本地编译好二进制和静态文件，上传到服务器直接运行，Docker 负责数据库、缓存和 Nginx。

::: tip 🎯 本节目标
完成后你能通过 `https://你的域名` 访问后台、登录并执行 CRUD。
:::

## 前置条件

- 本机有 Go 1.22+ 和 Node.js 22+
- 一个域名
- 腾讯云轻量应用服务器（Ubuntu 22.04）

---

## 🚀 第一步：本地构建

在项目根目录执行：

::: code-group

```powershell [Windows PowerShell]
cd server
$env:GOOS="linux"; $env:GOARCH="amd64"; go build -o server .
cd ..

cd admin
pnpm install; pnpm build
cd ..
```

```bash [macOS / Linux]
cd server
GOOS=linux GOARCH=amd64 go build -o server .
cd ..

cd admin
pnpm install && pnpm build
cd ..
```

:::

构建完成后你得到两个东西：
- `server/server` — Linux 可执行文件
- `admin/dist/` — 前端静态文件目录

---

## 🛠️ 第二步：服务器准备

### 购买服务器

[腾讯云轻量应用服务器](https://console.cloud.tencent.com/lighthouse)，推荐 2 核 2G + Ubuntu 22.04。购买后记下公网 IP。

### 配置防火墙

实例详情 → 防火墙，确认开放 22（SSH）、80（HTTP）、443（HTTPS）。

### 连接服务器并安装 Docker

推荐使用图形化 SSH 工具连接服务器，操作更直观方便：

- **FinalShell**（免费）：[https://www.hostbuf.com/t/988.html](https://www.hostbuf.com/t/988.html)
- **Xshell**（个人免费）：[https://www.netsarang.com/zh/xshell/](https://www.netsarang.com/zh/xshell/)

连接信息：
- 主机：你的服务器公网 IP
- 端口：22
- 用户名：`ubuntu`
- 密码：创建实例时设置的密码

连接成功后，在终端中执行以下命令安装 Docker：

```bash
curl -fsSL https://get.docker.com | sudo sh
```

验证安装：

```bash
docker --version && docker compose version
```

两条命令都返回版本号即可。

---

## 📦 第三步：上传文件

### 查看服务器系统信息（可选）

连接服务器后，可查看系统版本信息：

```bash
# 查看系统版本
cat /etc/os-release

# 查看内核版本
uname -a

# 查看当前用户
whoami
```

::: warning ⚠️ Ubuntu 默认用户名

如果使用 **Ubuntu 24.04.4 LTS**（推荐），默认用户名不是 `root`，而是 `ubuntu`。执行系统命令时需要用 `sudo` 提升权限。
:::

### 上传文件（使用图形化工具）

使用 FinalShell 或 Xshell 连接服务器后：

1. **创建目录**：在终端中执行
   ```bash
   mkdir -p /opt/ez-admin/nginx /opt/ez-admin/web /opt/ez-admin/ssl && sudo mkdir -p /etc/systemd/system
   ```

2. **上传文件**：通过工具的文件传输功能（SFTP）拖放以下文件：
   - 本地 `server/server` → 服务器 `/opt/ez-admin/`
   - 本地 `admin/dist/` 目录下所有文件 → 服务器 `/opt/ez-admin/web/`
   - 本地 `deploy/compose.server.yml` → 服务器 `/opt/ez-admin/`
   - 本地 `deploy/nginx/nginx-native.conf` → 服务器 `/opt/ez-admin/nginx/`
   - 本地 `deploy/.env.example` → 服务器 `/opt/ez-admin/.env`
   - 本地 `deploy/ez-admin.service` → 服务器 `/tmp/`，然后执行：
     ```bash
     sudo mv /tmp/ez-admin.service /etc/systemd/system/
     ```

---

## ⚙️ 第四步：配置环境变量

在服务器上编辑 `.env`：

```bash
nano /opt/ez-admin/.env
```

重点修改：

```bash {hl_lines="2 4 7"}
# 后端连接本地 Docker 中的数据库和缓存
EZ_DATABASE_HOST=127.0.0.1
EZ_DATABASE_PORT=5432
EZ_REDIS_HOST=127.0.0.1
EZ_REDIS_PORT=6379

# JWT 密钥（必须改，用 openssl rand -hex 32 生成）
EZ_AUTH_JWT_SECRET=你生成的随机字符串

# 数据库密码（建议改掉默认值）
EZ_DATABASE_PASSWORD=你的数据库密码
```

::: warning ⚠️ 注意 HOST 地址
后端直接运行在服务器上，不是在 Docker 容器里，所以数据库和 Redis 的 HOST 是 `127.0.0.1`，不是容器名。
:::

---

## 🚀 第五步：启动所有服务

在服务器上执行（Ubuntu 用户需在 `systemctl` 命令前加 `sudo`）：

```bash
# 1. 启动 PostgreSQL + Redis + Nginx
cd /opt/ez-admin && docker compose -f compose.server.yml up -d

# 2. 启动后端（通过 systemd，需要 sudo）
chmod +x /opt/ez-admin/server
sudo systemctl daemon-reload
sudo systemctl enable --now ez-admin
```

验证：

```bash
# 检查容器状态（应该看到 postgres、redis、nginx 三个容器）
docker compose -f /opt/ez-admin/compose.server.yml ps

# 检查后端是否运行（需要 sudo）
sudo systemctl status ez-admin

# 检查健康接口
curl http://localhost/health
```

三个都返回正常，服务就起来了。

### 初始化管理员

```bash
curl -X POST http://localhost/api/v1/setup/init \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"Admin@123456","nickname":"管理员"}'
```

### 浏览器验证

打开 `http://服务器IP`，登录并确认菜单和 CRUD 正常。

---

## 🌐 第六步：域名与 HTTPS

### Cloudflare 托管域名

1. 登录 [Cloudflare](https://dash.cloudflare.com)，添加域名，选择 Free 计划
2. 按提示到域名注册商修改 NS 为 Cloudflare 提供的地址，等待生效
3. 添加 A 记录：Type `A`，Name `@`，IPv4 填服务器 IP，Proxy status 暂选 **DNS only**

**验证**：`ping 你的域名` 解析到服务器 IP，`http://域名` 能访问。

### 配置 SSL 证书

1. Cloudflare → 你的域名 → SSL/TLS → Overview，设为 **Full（完全）**
2. SSL/TLS → Origin Server → Create Certificate（RSA 2048，15 年）
3. 保存证书和私钥到服务器：

```bash
# 粘贴 Origin Certificate
nano /opt/ez-admin/ssl/cert.pem
# 粘贴 Private Key
nano /opt/ez-admin/ssl/key.pem
```

4. 切换 Nginx 为 SSL 配置：
   - 通过 SFTP 将本地 `deploy/nginx/nginx-native-ssl.conf` 上传到服务器 `/opt/ez-admin/nginx/nginx-native.conf`

然后在服务器上重启 Nginx 容器：

```bash
cd /opt/ez-admin && docker compose -f compose.server.yml restart nginx
```

5. DNS 记录开启代理（橙色云朵）

**验证**：`https://你的域名` 正常访问，浏览器显示锁头图标。

---

## 🔄 更新发布

改完代码后，按以下步骤更新：

### 1. 本地编译构建
::: code-group

```powershell [Windows PowerShell]
# 编译后端
cd server
$env:GOOS="linux"; $env:GOARCH="amd64"; go build -o server .
cd ..

# 构建前端
cd admin; pnpm build; cd ..
```

```bash [macOS / Linux]
# 编译后端
cd server && GOOS=linux GOARCH=amd64 go build -o server . && cd ..

# 构建前端
cd admin && pnpm build && cd ..
```

:::

### 2. 上传文件（使用图形化工具）
通过 FinalShell 或 Xshell 的 SFTP 功能上传：
- 本地 `server/server` → 服务器 `/opt/ez-admin/server`（覆盖原有文件）
- 本地 `admin/dist/` 目录下所有文件 → 服务器 `/opt/ez-admin/web/`（覆盖原有文件）

### 3. 重启后端
在服务器终端执行：
```bash
sudo systemctl restart ez-admin
```

如果只改了后端，只需要编译+上传后端+重启。如果只改了前端，只需要构建+上传前端。

---

## ✅ 部署验证清单

| 验证项 | 期望结果 |
| --- | --- |
| 容器状态 | PostgreSQL、Redis、Nginx 均 running/healthy |
| 后端服务 | `systemctl status ez-admin` 显示 active |
| 健康接口 | `curl http://localhost/health` 返回 ok |
| 管理员初始化 | `/api/v1/setup/init` 返回成功 |
| IP 访问 | `http://服务器IP` 能登录 |
| HTTPS 域名 | `https://域名` 正常访问 |
| CDN 代理 | `ping 域名` 不显示真实 IP |

---

## 常见问题排查

::: details 后端启动失败，报数据库连接拒绝
确认 Docker 容器在运行：`docker compose -f /opt/ez-admin/compose.server.yml ps`。

确认 `.env` 中 `EZ_DATABASE_HOST=127.0.0.1`（不是 `postgres`）。
:::

::: details Nginx 报 502 Bad Gateway
后端还没启动或已崩溃。检查：`systemctl status ez-admin`，查看日志：`journalctl -u ez-admin -f`。
:::

::: details 前端白屏
确认前端文件已上传到 `/opt/ez-admin/web/` 目录，且 Nginx 配置中有 `try_files $uri $uri/ /index.html;`。
:::

::: details Cloudflare ERR_TOO_MANY_REDIRECTS
SSL 加密模式设成了 Flexible。改为 **Full（完全）**。
:::

::: details 更新后端后接口没变化
确认上传了新二进制且执行了 `systemctl restart ez-admin`。
:::

---

## 🛠️ 复用：开始一个新项目

### 1. 复制仓库

```bash
cp -r ez-admin-gin my-new-project && cd my-new-project
```

### 2. 改模块名

`server/go.mod` 中把 `ez-admin-gin/server` 替换为你的项目名，然后 `find . -name "*.go" -exec sed -i '' 's|ez-admin-gin/server|my-new-project/server|g' {} +`，最后 `go mod tidy`。

### 3. 加业务模块

按第 6 章的规范在 `server/internal/` 下新增模块，在 `admin/src/` 下新增页面。

### 4. 部署

按本节的步骤操作即可。

---

## 小结

- Docker Compose 负责 PostgreSQL、Redis 和 Nginx，后端二进制直接运行在宿主机上。
- 后端是一个 Linux 二进制，前端是一份静态文件，上传就能跑。
- 更新只需要重新编译、上传、重启后端。
- Cloudflare 提供免费的 CDN + HTTPS + IP 隐藏。

回到本章总览：[第 7 章：部署与复用](./)。
