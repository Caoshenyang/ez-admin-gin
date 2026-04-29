---
title: 部署验证与复用说明
description: "本地编译后端二进制、构建前端静态文件，打包成压缩包上传到服务器，通过 Docker 一键启动。"
---

# 部署验证与复用说明

这一节把后台底座部署到公网。思路很简单：本地编译好二进制和静态文件，打成压缩包上传到服务器解压运行，Docker 负责数据库、缓存和 Nginx。

::: tip 🎯 本节目标
完成后你能通过 `https://你的域名` 访问后台、登录并执行 CRUD。
:::

## 前置条件

- 本机有 Go 1.22+ 和 Node.js 22+
- 一个域名
- 腾讯云轻量应用服务器（Ubuntu 22.04）

---

## 🚀 第一步：本地打包

在项目根目录执行：

::: code-group

```powershell [Windows PowerShell]
.\scripts\pack.ps1
```

```bash [macOS / Linux]
bash scripts/pack.sh
```

:::

脚本会自动完成：

1. 编译后端（Linux amd64，已去掉调试符号）
2. 构建前端静态文件
3. 收集部署配置文件
4. 打包成压缩包（Windows 生成 `.zip`，macOS/Linux 生成 `.tar.gz`）

完成后你会看到：

```text
✅ 打包完成！上传 deploy-package.tar.gz 到服务器即可。
```

产物是项目根目录下的一个压缩包文件：`deploy-package.tar.gz` 或 `deploy-package.zip`。

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

::: warning ⚠️ Ubuntu 默认用户名
Ubuntu 默认用户名不是 `root`，而是 `ubuntu`。执行系统命令时需要用 `sudo` 提升权限。
:::

连接成功后，在终端中执行以下命令安装 Docker：

```bash
curl -fsSL https://get.docker.com | sudo sh
```

验证安装：

```bash
docker --version && docker compose version
```

两条命令都返回版本号即可。

::: warning ⚠️ 国内服务器拉镜像可能失败
如果后续执行 setup 脚本时报 `failed to resolve reference "docker.io"` 错误，说明无法连接 Docker Hub。请参考 [腾讯云镜像源加速 Docker](https://cloud.tencent.com/document/product/213/8623#.E4.BD.BF.E7.94.A8.E8.85.BE.E8.AE.AF.E4.BA.91.E9.95.9C.E5.83.8F.E6.BA.90.E5.8A.A0.E9.80.9F-docker) 配置镜像加速，修改后执行 `sudo systemctl restart docker` 重启 Docker，再重新运行 setup 脚本。
:::

---

## 📦 第三步：上传并启动

### 1. 上传压缩包

在 FinalShell / Xshell 的文件管理器中，把本地生成的压缩包（`deploy-package.tar.gz` 或 `deploy-package.zip`）上传到服务器的 `/opt/` 目录下。

### 2. 解压

在服务器终端执行：

::: code-group

```bash [tar.gz]
sudo mkdir -p /opt/ez-admin
sudo tar xzf /opt/deploy-package.tar.gz -C /opt/ez-admin/
```

```bash [zip]
sudo mkdir -p /opt/ez-admin
sudo unzip /opt/deploy-package.zip -d /opt/ez-admin/
```

:::

### 3. 首次部署：运行 setup 脚本

```bash
sudo bash /opt/ez-admin/setup-server.sh
```

::: warning ⚠️ 此脚本仅在首次部署时执行
`setup-server.sh` 负责 Docker 环境初始化、密钥生成和管理员创建。后续更新请使用 `update-server.sh`，不要重复执行 setup。
:::

脚本会自动完成：

- 整理文件结构（创建子目录、移动配置文件到正确位置）
- 自动生成 JWT 密钥（仅首次）
- 启动 PostgreSQL + Redis + Nginx 容器
- 等待数据库就绪后启动后端
- 自动初始化管理员账号（仅首次）
- 打印访问地址和默认账号

完成后会看到：

```text
=========================================
✅ 部署完成！

  访问地址：http://你的服务器IP
  默认账号：admin / Admin@123456

  查看后端日志：sudo journalctl -u ez-admin -f
  查看容器状态：docker compose -f /opt/ez-admin/compose.server.yml ps
=========================================
```

打开 `http://服务器IP`，登录并确认菜单和 CRUD 正常。**首次登录后请修改默认密码。**

---

## 🌐 第四步：域名与 HTTPS

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
# 粘贴 Origin Certificate 内容
sudo nano /opt/ez-admin/ssl/cert.pem
# 粘贴 Private Key 内容
sudo nano /opt/ez-admin/ssl/key.pem
```

4. 切换 Nginx 为 SSL 配置：将本地 `deploy/nginx/nginx-native-ssl.conf` 上传到 `/opt/ez-admin/nginx/nginx-native.conf`（覆盖原文件）

然后在服务器上重启 Nginx 容器：

```bash
cd /opt/ez-admin && sudo docker compose -f compose.server.yml restart nginx
```

5. DNS 记录开启代理（橙色云朵）

**验证**：`https://你的域名` 正常访问，浏览器显示锁头图标。

---

## 🔄 更新发布

改完代码后，按以下步骤更新：

**1. 本地重新打包**

::: code-group

```powershell [Windows PowerShell]
.\scripts\pack.ps1
```

```bash [macOS / Linux]
bash scripts/pack.sh
```

:::

**2. 上传压缩包到服务器**

通过 FinalShell / Xshell 上传新的压缩包到 `/opt/`。

**3. 解压并执行更新脚本**

::: code-group

```bash [tar.gz]
sudo tar xzf /opt/deploy-package.tar.gz -C /opt/ez-admin/
sudo bash /opt/ez-admin/update-server.sh
```

```bash [zip]
sudo unzip -o /opt/deploy-package.zip -d /opt/ez-admin/
sudo bash /opt/ez-admin/update-server.sh
```

:::

`update-server.sh` 只替换文件并重启后端，不会动 Docker 环境和已有配置。

---

## ✅ 验证清单

| 验证项 | 命令 | 期望结果 |
| --- | --- | --- |
| 容器状态 | `sudo docker compose -f /opt/ez-admin/compose.server.yml ps` | postgres、redis、nginx 均 running/healthy |
| 后端服务 | `sudo systemctl status ez-admin` | 显示 active (running) |
| IP 访问 | 浏览器打开 `http://服务器IP` | 能登录 |
| HTTPS 域名 | 浏览器打开 `https://域名` | 正常访问，显示锁头 |

---

## 常见问题排查

::: details 后端启动失败，报数据库连接拒绝
确认 Docker 容器在运行：`sudo docker compose -f /opt/ez-admin/compose.server.yml ps`。

确认 `.env` 中 `EZ_DATABASE_HOST=127.0.0.1`（不是 `postgres`）。

如果容器刚启动，可能健康检查还没通过，等 30 秒再试。
:::

::: details Nginx 报 502 Bad Gateway
后端还没启动或已崩溃。检查：`sudo systemctl status ez-admin`，查看日志：

```bash
sudo journalctl -u ez-admin -f
```
:::

::: details 前端白屏
确认前端文件在 `/opt/ez-admin/web/` 目录（不是 `dist/`），且目录下有 `index.html`。

确认 Nginx 配置中有 `try_files $uri $uri/ /index.html;`。
:::

::: details Cloudflare ERR_TOO_MANY_REDIRECTS
SSL 加密模式设成了 Flexible。改为 **Full（完全）**。
:::

::: details 更新后端后接口没变化
确认上传了新二进制且执行了 `sudo bash /opt/ez-admin/update-server.sh`（脚本会自动重启后端）。
:::

::: details 初始化管理员返回 409
说明管理员已经创建过了，不是错误。直接用已有账号登录即可。
:::

---

## 🛠️ 复用：开始一个新项目

### 1. 复制仓库

```bash
cp -r ez-admin-gin my-new-project && cd my-new-project
```

### 2. 改模块名

`server/go.mod` 中把 `ez-admin-gin/server` 替换为你的项目名，然后批量替换所有 Go 文件中的 import 路径：

::: code-group

```powershell [Windows PowerShell]
Get-ChildItem -Recurse -Filter *.go | ForEach-Object {
    (Get-Content $_.FullName) -replace 'ez-admin-gin/server', 'my-new-project/server' | Set-Content $_.FullName
}
go mod tidy
```

```bash [macOS / Linux]
find . -name "*.go" -exec sed -i 's|ez-admin-gin/server|my-new-project/server|g' {} +
go mod tidy
```

:::

### 3. 加业务模块

按第 6 章的规范在 `server/internal/` 下新增模块，在 `admin/src/` 下新增页面。

### 4. 部署

按本节的步骤操作即可。

---

## 小结

- Docker Compose 负责 PostgreSQL、Redis 和 Nginx，后端二进制直接运行在宿主机上。
- 后端是一个 Linux 二进制，前端是一份静态文件，打成一个压缩包上传即可部署。
- 更新只需要重新打包、上传、解压、执行 setup 脚本。
- Cloudflare 提供免费的 CDN + HTTPS + IP 隐藏。

回到本章总览：[第 7 章：部署与复用](./)。
