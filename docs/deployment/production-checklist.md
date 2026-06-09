---
title: 生产环境检查清单
description: 上线前必须逐项检查的安全、数据、服务和应用验证清单
---

# 生产环境检查清单

正式上线前，逐项确认以下配置。遗漏任何一项都可能带来安全风险或线上故障。

## 安全

### 密钥和密码

- [ ] `EZ_AUTH_JWT_SECRET` 已替换为随机字符串（至少 32 字符）
  ```bash
  # 生成密钥
  openssl rand -hex 32
  # 更新 .env
  vim /opt/ez-admin/.env
  ```
- [ ] `EZ_DATABASE_PASSWORD` 已修改为强密码
  ```bash
  # 生成随机密码
  openssl rand -base64 24
  ```
  ::: warning
  PostgreSQL 的 `POSTGRES_PASSWORD` 仅在首次创建数据库时生效。如果已用旧密码创建了数据库，需要进入容器修改：
  ```bash
  docker compose -f compose.server.yml exec postgres psql -U ez_admin -c "ALTER USER ez_admin PASSWORD '新密码';"
  ```
  :::
- [ ] `EZ_REDIS_PASSWORD` 已设置（建议）
  ```bash
  # 生成 Redis 密码
  openssl rand -hex 16
  # 修改后重建 Redis 容器
  cd /opt/ez-admin
  docker compose -f compose.server.yml up -d redis
  sudo systemctl restart ez-admin
  ```
- [ ] 默认管理员密码已修改（首次登录后立即修改）

### 网络和访问

- [ ] `EZ_CORS_ALLOWED_ORIGINS` 已设为实际域名
  ```bash
  EZ_CORS_ALLOWED_ORIGINS=https://admin.example.com
  ```
  留空或使用 `*` 会导致任意网站都能调用 API。
- [ ] HTTPS 已配置，参见 [域名与 HTTPS](/deployment/domain-and-https)
- [ ] 防火墙仅开放必要端口
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
- [ ] 数据库和 Redis 未对外暴露
  ```bash
  # 只应该看到 127.0.0.1 的监听
  sudo ss -tlnp | grep -E "5432|6379"
  # 不应该有 0.0.0.0 的监听
  ```
  `compose.server.yml` 已将数据库和 Redis 绑定到 `127.0.0.1`，确认没有意外暴露。
- [ ] 上传目录权限正确
  ```bash
  ls -la /opt/ez-admin/uploads/
  # 确保后端进程有写入权限
  ```

### 应用安全

- [ ] Swagger / API 文档已关闭（如不需要对外暴露）
  ```bash
  EZ_SWAGGER_ENABLED=false
  ```
- [ ] 日志级别为 `info` 或 `warn`，不是 `debug`
  ```bash
  EZ_LOG_LEVEL=info
  EZ_LOG_FORMAT=json
  ```

## 数据

- [ ] 数据库备份策略已配置（建议每天备份，保留 7 天）
  ```bash
  # 手动备份
  docker compose -f compose.server.yml exec -T postgres \
    pg_dump -U ez_admin ez_admin > backup-$(date +%Y%m%d).sql
  ```
  配置定时备份参见 [更新与回滚 → 数据库备份](/deployment/update-and-rollback#数据库备份)。
- [ ] 升级前会先备份数据库
- [ ] migration 执行前会确认 SQL 内容
- [ ] 种子数据不会重复污染生产环境（`setup-server.sh` 在管理员已存在时会跳过初始化）

## 服务

### systemd 服务

- [ ] 后端服务状态正常
  ```bash
  sudo systemctl status ez-admin
  # 期望：Active: active (running)
  ```
- [ ] 后端服务已设置为开机自启
  ```bash
  sudo systemctl is-enabled ez-admin
  # 期望：enabled
  ```

### Docker 容器

- [ ] 所有容器正常运行
  ```bash
  cd /opt/ez-admin
  docker compose -f compose.server.yml ps
  # 期望：postgres (healthy)、redis (healthy)、nginx (running)
  ```

### Nginx

- [ ] Nginx 配置测试通过
  ```bash
  docker exec ez-admin-nginx nginx -t
  ```
- [ ] 前端静态资源可访问
- [ ] API 代理正常工作

### 磁盘和日志

- [ ] 日志目录权限正确
  ```bash
  ls -la /opt/ez-admin/logs/ 2>/dev/null || echo "目录不存在"
  ```
- [ ] 上传目录权限正确
- [ ] 磁盘空间充足
  ```bash
  df -h /opt
  ```
- [ ] 时区配置正确
  ```bash
  timedatectl
  ```

## 应用验证

### 登录和权限

- [ ] 能正常登录后台
- [ ] 登录后菜单正常显示
- [ ] 角色权限验证正常（不同角色看到不同菜单）
- [ ] 按钮权限验证正常（无权限按钮不显示）
- [ ] 数据权限验证正常（不同角色看到不同数据范围）

### 功能验证

- [ ] 前端页面刷新后不出现 404（SPA fallback 正常）
- [ ] 文件上传功能正常
- [ ] API 请求正常代理（`/api` 路径）
- [ ] 上传文件可正常访问（`/uploads` 路径）
- [ ] 健康检查端点正常
  ```bash
  curl -sf http://localhost/health && echo "OK"
  ```

## 配置文件检查清单

确认 `/opt/ez-admin/.env` 中以下变量已正确配置：

```bash
# 必须修改
EZ_AUTH_JWT_SECRET=<随机字符串>
EZ_DATABASE_PASSWORD=<强密码>

# 建议修改
EZ_REDIS_PASSWORD=<密码>
EZ_CORS_ALLOWED_ORIGINS=https://your-domain.com
EZ_LOG_LEVEL=info

# 确认环境
EZ_APP_ENV=prod
```
