---
title: 更新与回滚
description: 版本更新流程、日志查看和回滚建议
---

# 更新与回滚

本文介绍如何更新已部署的 EZ Admin 服务，以及如何在出现问题时回滚到上一个版本。

::: tip 适用场景
- 代码更新后需要重新部署
- 部署后发现问题需要回滚
- 日常运维中检查服务状态和查看日志
:::

## 什么时候需要更新

- 后端代码有新提交，需要重新编译部署
- 前端页面有更新，需要重新构建部署
- 修复了线上 Bug，需要紧急发布

## 更新前检查

更新前请确认：

- [ ] 本地代码已拉取最新
- [ ] 如果有数据库表结构变更，已准备好备份方案
- [ ] 确认服务器磁盘空间充足（`df -h /opt`）
- [ ] 确认当前服务运行正常（`sudo systemctl status ez-admin`）

## 重新打包

::: code-group

```bash [macOS / Linux]
bash scripts/pack.sh
```

```powershell [Windows]
.\scripts\pack.ps1
```

:::

打包完成后会生成 `deploy-package.tar.gz`（或 Windows 下的 `.zip`）。

## 上传更新包

### 方式一：一键更新（推荐）

在本地直接执行：

```bash
bash scripts/deploy.sh user@your-server
```

这个脚本会自动完成打包、上传和远端初始化。

::: warning
`deploy.sh` 在远端执行的是 `setup-server.sh`（首次部署脚本），不是 `update-server.sh`。`setup-server.sh` 会检测已有配置和密钥，不会覆盖。但如果服务器禁用了远程 sudo，此方式不可用。
:::

### 方式二：手动上传

```bash
# 本地上传
scp deploy-package.tar.gz user@your-server:/tmp/

# SSH 到服务器
ssh user@your-server

# 备份当前版本（建议）
cd /opt/ez-admin
sudo tar czf /opt/ez-admin-backup-$(date +%Y%m%d%H%M).tar.gz \
  --exclude='data' \
  --exclude='ssl' \
  .

# 解压更新包
cd /opt/ez-admin
tar xzf /tmp/deploy-package.tar.gz
```

## 执行更新

```bash
sudo bash /opt/ez-admin/update-server.sh
```

### update-server.sh 做了什么

1. 将新的 `dist/` 重命名为 `web/`（替换前端文件）
2. 赋予新二进制可执行权限
3. 重启后端 systemd 服务
4. 等待健康检查通过（最多 15 秒）

### update-server.sh 不负责什么

::: warning update-server.sh 的边界
`update-server.sh` 只做文件替换和后端重启。它**不会**：

- 备份数据库
- 执行或回滚数据库迁移
- 更新 Docker 容器镜像
- 自动保留旧版本二进制
- 回滚 Redis 数据
- 合并 `.env` 配置变更
:::

## 查看服务状态

```bash
# 后端服务状态
sudo systemctl status ez-admin
# 期望：Active: active (running)

# Docker 容器状态
cd /opt/ez-admin
docker compose -f compose.server.yml ps
# 期望：三个容器均为 Up / running (healthy)

# 快速健康检查
curl -sf http://localhost/health && echo "OK" || echo "FAIL"
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

## 回滚建议

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

::: warning 回滚不包含数据库
回滚只恢复程序文件（后端二进制、前端静态文件、配置文件）。如果新版本已经执行了数据库迁移，回滚后旧版本可能不兼容新的表结构。

建议在每次更新前备份数据库，见下方"数据库备份"。
:::

## 推荐实践：版本化部署目录

当前脚本的回滚依赖手动备份。如果项目规模增长，建议采用版本化目录结构：

```
/opt/ez-admin/
├── releases/
│   ├── 20250101-120000/     # 按时间戳归档版本
│   │   ├── server
│   │   ├── web/
│   │   └── configs/
│   └── 20250102-150000/
├── current -> releases/20250102-150000/   # 软链接指向当前版本
├── data/                   # 数据（不跟随版本）
├── .env                    # 配置（不跟随版本）
└── ssl/                    # 证书（不跟随版本）
```

::: warning
这是推荐的未来优化方案，当前脚本没有实现此结构。如果你需要版本化部署，需要自行改造 `update-server.sh`。
:::

## 数据库迁移注意事项

如果本次更新包含数据库表结构变更：

1. **备份数据库**（见下方）
2. 启动新版本后端 — 后端会自动执行 GORM AutoMigrate
3. 观察日志确认迁移完成

```bash
sudo journalctl -u ez-admin -f
# 看到 "auto migrate success" 或服务正常启动即表示迁移完成
```

::: warning
GORM AutoMigrate 只会添加新列和新表，**不会**删除列或修改列类型。如果你的迁移涉及删除或重命名列，需要手动执行 SQL。
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
# 创建备份目录
sudo mkdir -p /opt/backups

# 编辑 crontab
sudo crontab -e

# 每天凌晨 3 点备份
0 3 * * * cd /opt/ez-admin && docker compose -f compose.server.yml exec -T postgres pg_dump -U ez_admin ez_admin | gzip > /opt/backups/postgres-$(date +\%Y\%m\%d).sql.gz
```

最低要求：每天备份一次，保留最近 7 天。
