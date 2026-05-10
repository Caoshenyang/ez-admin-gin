---
title: 生产环境检查清单
description: 上线前必须逐项检查的安全和配置清单
---

# 生产环境检查清单

正式上线前，请逐项确认以下配置。遗漏任何一项都可能带来安全风险。

## 必须修改

### JWT 密钥

`.env` 中的 `EZ_AUTH_JWT_SECRET` 必须替换为随机字符串。

```bash
# 生成密钥
openssl rand -hex 32

# 更新 .env
vim /opt/ez-admin/.env
```

::: danger
`EZ_APP_ENV=prod` 且密钥包含 `change-me` 时，后端会拒绝启动。即使你绕过了这个检查，使用默认密钥意味着任何人都能伪造 JWT token。
:::

### 数据库密码

将 `EZ_DATABASE_PASSWORD` 从默认的 `ez_admin_123456` 改为强密码：

```bash
# 生成随机密码
openssl rand -base64 24
```

修改后需要同步更新 PostgreSQL 容器：

```bash
# 1. 更新 .env 中的 EZ_DATABASE_PASSWORD
# 2. 重建容器
cd /opt/ez-admin
docker compose -f compose.server.yml down
docker compose -f compose.server.yml up -d
# 3. 重启后端
sudo systemctl restart ez-admin
```

::: warning
PostgreSQL 的密码通过 `POSTGRES_PASSWORD` 环境变量设置。这个变量仅在**首次创建数据库**时生效。如果你已经用旧密码创建了数据库，需要进入容器修改用户密码：

```bash
docker compose -f compose.server.yml exec postgres psql -U ez_admin -c "ALTER USER ez_admin PASSWORD '新密码';"
```
:::

### 默认管理员密码

部署脚本会创建 `admin / Admin@123456` 管理员账号。首次登录后必须修改密码。

### CORS 白名单

`EZ_CORS_ALLOWED_ORIGINS` 必须设为实际域名：

```bash
EZ_CORS_ALLOWED_ORIGINS=https://admin.example.com
```

留空或使用 `*` 会导致任意网站都能调用你的 API。

## 建议配置

### HTTPS

配置 HTTPS 证书，确保传输安全。参见 [域名与 HTTPS](/deployment/domain-and-https)。

### 关闭 Swagger

生产环境不应暴露 API 文档：

```bash
EZ_SWAGGER_ENABLED=false
```

### 日志级别

生产环境建议使用 `info` 或 `warn`，避免 `debug` 级别输出敏感信息：

```bash
EZ_LOG_LEVEL=info
EZ_LOG_FORMAT=json
```

### Redis 密码

设置 Redis 密码防止未授权访问：

```bash
EZ_REDIS_PASSWORD=$(openssl rand -hex 16)
```

修改后重建 Redis 容器：

```bash
cd /opt/ez-admin
docker compose -f compose.server.yml up -d redis
sudo systemctl restart ez-admin
```

## 网络安全

### 端口暴露检查

在服务器二进制部署中，`compose.server.yml` 已将数据库和 Redis 绑定到 `127.0.0.1`：

```yaml
ports:
  - "127.0.0.1:5432:5432"   # PostgreSQL 仅本机
  - "127.0.0.1:6379:6379"   # Redis 仅本机
```

确认没有意外暴露：

```bash
# 只应该看到 80 / 443 对外开放
sudo ss -tlnp | grep -E "5432|6379"
# 不应该有 0.0.0.0 的监听
```

### 防火墙

只开放必要端口：

```bash
# Ubuntu UFW 示例
sudo ufw allow 22/tcp    # SSH
sudo ufw allow 80/tcp    # HTTP
sudo ufw allow 443/tcp   # HTTPS
sudo ufw enable
```

::: warning
确保 SSH 端口（22）已放行后再启用防火墙，否则会锁死自己。
:::

## 数据安全

### 文件上传目录权限

上传目录应只允许后端进程写入：

```bash
# 检查上传目录权限
ls -la /opt/ez-admin/uploads/ 2>/dev/null || echo "目录不存在"
```

### 数据库备份

配置定时备份，参见 [更新与回滚 → 数据库备份](/deployment/update-and-rollback#数据库备份)。

最低要求：每天备份一次，保留最近 7 天。

## 完整检查清单

- [ ] `EZ_AUTH_JWT_SECRET` 已替换为随机字符串
- [ ] `EZ_DATABASE_PASSWORD` 已修改为强密码
- [ ] `EZ_CORS_ALLOWED_ORIGINS` 已设为实际域名
- [ ] 默认管理员密码已修改
- [ ] HTTPS 已配置
- [ ] Swagger 已关闭（`EZ_SWAGGER_ENABLED=false`）
- [ ] 日志级别为 `info` 或 `warn`
- [ ] Redis 密码已设置
- [ ] 数据库和 Redis 未对外暴露
- [ ] 防火墙仅开放 22 / 80 / 443
- [ ] 文件上传目录权限正确
- [ ] 数据库定时备份已配置
