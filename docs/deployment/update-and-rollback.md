---
title: 更新与回滚
description: 版本更新流程、日志查看和服务回滚
---

# 更新与回滚

本文介绍如何更新已部署的 EZ Admin 服务，以及如何在出现问题时回滚到上一个版本。

::: tip 适用场景
- 代码更新后需要重新部署
- 部署后发现问题需要回滚
- 日常运维中检查服务状态和查看日志
:::

## 更新流程

### 方式一：一键更新（推荐）

在本地执行 `deploy.sh`，它会自动完成打包、上传和远端更新：

```bash
bash scripts/deploy.sh user@your-server
```

这个脚本会：
1. 编译后端、构建前端、打包配置文件
2. 上传到服务器 `/tmp/`
3. 远端解压并执行 `setup-server.sh`

### 方式二：手动更新

如果需要更精细的控制，可以分步操作：

#### 1. 本地打包

::: code-group
```bash [macOS / Linux]
bash scripts/pack.sh
```
```powershell [Windows]
.\scripts\pack.ps1
```
:::

#### 2. 备份当前部署

```bash
# SSH 到服务器
ssh user@your-server

# 备份当前版本
cd /opt/ez-admin
sudo tar czf /opt/ez-admin-backup-$(date +%Y%m%d%H%M).tar.gz \
  --exclude='data' \
  --exclude='ssl' \
  .
```

::: tip
备份时排除 `data/`（数据库数据）和 `ssl/`（证书），这些文件不需要跟随版本更新。
:::

#### 3. 上传并解压

```bash
# 本地执行
scp deploy-package.tar.gz user@your-server:/tmp/

# 服务器执行
cd /opt/ez-admin
tar xzf /tmp/deploy-package.tar.gz
```

#### 4. 执行更新

```bash
# 在服务器上执行更新脚本
sudo bash /opt/ez-admin/update-server.sh
```

`update-server.sh` 做了什么：

1. 将新的 `dist/` 替换为 `web/`（前端文件更新）
2. 赋予新二进制可执行权限
3. 重启后端 systemd 服务
4. 等待健康检查通过（最多 15 秒）

::: warning
`update-server.sh` **不会**修改 `.env`、数据库数据或 SSL 证书。只替换程序文件并重启后端。
:::

## 数据库迁移注意事项

如果本次更新包含数据库表结构变更：

1. **备份数据库**（见下方备份方法）
2. 启动新版本后端 — 后端会自动执行 GORM AutoMigrate
3. 观察日志确认迁移完成：

```bash
sudo journalctl -u ez-admin -f
# 看到 "auto migrate success" 或服务正常启动即表示迁移完成
```

::: warning
GORM AutoMigrate 只会添加新列和新表，**不会**删除列或修改列类型。如果你的迁移涉及删除或重命名列，需要手动执行 SQL。
:::

## 检查服务状态

```bash
# 后端服务状态
sudo systemctl status ez-admin

# Docker 容器状态
cd /opt/ez-admin
docker compose -f compose.server.yml ps

# 快速健康检查
curl -sf http://localhost/health && echo "OK" || echo "FAIL"
```

### 期望输出

`systemctl status` 显示：

```
Active: active (running)
```

`docker compose ps` 显示三个容器均为 `Up`：

```
ez-admin-postgres   running (healthy)
ez-admin-redis      running (healthy)
ez-admin-nginx      running
```

## 查看日志

### 后端日志

```bash
# 实时跟踪
sudo journalctl -u ez-admin -f

# 最近 100 行
sudo journalctl -u ez-admin -n 100

# 指定时间段
sudo journalctl -u ez-admin --since "2025-01-01" --until "2025-01-02"
```

### Docker 容器日志

```bash
# PostgreSQL 日志
docker compose -f compose.server.yml logs postgres

# Redis 日志
docker compose -f compose.server.yml logs redis

# Nginx 日志
docker compose -f compose.server.yml logs nginx
```

## 回滚

如果新版本出现问题，可以回滚到之前备份的版本。

### 1. 停止后端

```bash
sudo systemctl stop ez-admin
```

### 2. 恢复备份

```bash
cd /opt/ez-admin

# 查看可用备份
ls -lt /opt/ez-admin-backup-*.tar.gz

# 恢复（替换时间戳为实际备份文件）
sudo tar xzf /opt/ez-admin-backup-202501011200.tar.gz
```

### 3. 重启后端

```bash
sudo systemctl start ez-admin
```

### 4. 验证

```bash
sudo systemctl status ez-admin
curl -sf http://localhost/health && echo "OK"
```

::: warning
回滚不会回滚数据库。如果新版本已经执行了数据库迁移，回滚后旧版本可能不兼容新的表结构。建议在更新前备份数据库。
:::

## 数据库备份

### 手动备份

```bash
# 导出完整数据库
docker compose -f compose.server.yml exec -T postgres \
  pg_dump -U ez_admin ez_admin > backup-$(date +%Y%m%d).sql

# 恢复
cat backup-20250101.sql | docker compose -f compose.server.yml exec -T postgres \
  psql -U ez_admin ez_admin
```

### 定时备份（crontab）

```bash
# 编辑 crontab
sudo crontab -e

# 每天凌晨 3 点备份
0 3 * * * cd /opt/ez-admin && docker compose -f compose.server.yml exec -T postgres pg_dump -U ez_admin ez_admin | gzip > /opt/backups/postgres-$(date +\%Y\%m\%d).sql.gz
```

```bash
# 创建备份目录
sudo mkdir -p /opt/backups
```
