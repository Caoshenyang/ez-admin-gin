---
title: 环境变量参考
description: "集中记录当前底座支持的 EZ_* 环境变量、默认值、用途和部署时最值得优先确认的项。"
---

# 环境变量参考

这页只做快速查阅，不替代第 9 章的部署步骤。  
如果你需要确认“某个变量有没有、默认值是什么、上线前哪些最该先改”，直接查这里。

::: tip 🎯 这页解决什么
快速回答下面这些问题：

- 当前支持哪些 `EZ_*` 环境变量
- 它们分别映射到哪段配置
- 哪些是必须改、哪些是建议改
:::

## 环境变量加载规则

当前服务端配置遵循这条规则：

```text
config.yaml 默认值
  ↓
bindEnvs(...)
  ↓
EZ_* 环境变量覆盖
```

关键实现位置：

- `server/internal/config/config.go`
- `deploy/.env.example`

也就是说：

- 本地开发可以主要看 `configs/config.yaml`
- 服务器部署时更应该围绕 `.env` 工作

## 部署前最值得优先确认的变量

当前最重要的几项可以先直接看这张表：

| 变量 | 优先级 | 原因 |
| --- | --- | --- |
| `EZ_AUTH_JWT_SECRET` | 必须 | 默认占位值不能进生产 |
| `EZ_DATABASE_PASSWORD` | 强烈建议 | 默认密码过弱 |
| `EZ_APP_ENV` | 建议 | 生产环境应设为 `prod` |
| `EZ_LOG_FORMAT` | 建议 | 生产更适合 `json` |
| `EZ_SERVER_ADDR` | 视环境决定 | 默认监听 `:8080` |

## 当前 `deploy/.env.example` 中实际出现的变量

| 变量 | 默认值 | 作用 |
| --- | --- | --- |
| `EZ_AUTH_JWT_SECRET` | `change-me-to-a-random-string-at-least-32-chars` | JWT 签名密钥 |
| `EZ_DATABASE_HOST` | `127.0.0.1` | 数据库地址 |
| `EZ_DATABASE_PORT` | `5432` | 数据库端口 |
| `EZ_DATABASE_USER` | `ez_admin` | 数据库用户名 |
| `EZ_DATABASE_PASSWORD` | `ez_admin_123456` | 数据库密码 |
| `EZ_DATABASE_NAME` | `ez_admin` | 数据库名 |
| `EZ_REDIS_HOST` | `127.0.0.1` | Redis 地址 |
| `EZ_REDIS_PORT` | `6379` | Redis 端口 |
| `EZ_REDIS_PASSWORD` | 空 | Redis 密码 |
| `EZ_APP_ENV` | `prod` | 当前运行环境 |
| `EZ_SERVER_ADDR` | `:8080` | 后端监听地址 |
| `EZ_LOG_LEVEL` | `info` | 日志级别 |
| `EZ_LOG_FORMAT` | `json` | 日志格式 |

## 配置结构支持的完整环境变量分组

除了部署模板里已经写出来的变量，当前配置系统还支持下面这些覆盖项。

### App / Server

| 环境变量 | 对应配置键 | 默认值 |
| --- | --- | --- |
| `EZ_APP_NAME` | `app.name` | `ez-admin` |
| `EZ_APP_ENV` | `app.env` | `dev` |
| `EZ_SERVER_ADDR` | `server.addr` | `:8080` |

### Database

| 环境变量 | 对应配置键 | 默认值 |
| --- | --- | --- |
| `EZ_DATABASE_DRIVER` | `database.driver` | `postgres` |
| `EZ_DATABASE_HOST` | `database.host` | `localhost` |
| `EZ_DATABASE_PORT` | `database.port` | `5432` |
| `EZ_DATABASE_USER` | `database.user` | `ez_admin` |
| `EZ_DATABASE_PASSWORD` | `database.password` | `ez_admin_123456` |
| `EZ_DATABASE_NAME` | `database.name` | `ez_admin` |
| `EZ_DATABASE_MAX_IDLE_CONNS` | `database.max_idle_conns` | `10` |
| `EZ_DATABASE_MAX_OPEN_CONNS` | `database.max_open_conns` | `50` |
| `EZ_DATABASE_CONN_MAX_LIFETIME` | `database.conn_max_lifetime` | `3600` |

### Redis

| 环境变量 | 对应配置键 | 默认值 |
| --- | --- | --- |
| `EZ_REDIS_HOST` | `redis.host` | `localhost` |
| `EZ_REDIS_PORT` | `redis.port` | `6379` |
| `EZ_REDIS_PASSWORD` | `redis.password` | 空 |
| `EZ_REDIS_DB` | `redis.db` | `0` |
| `EZ_REDIS_MAX_RETRIES` | `redis.max_retries` | `3` |
| `EZ_REDIS_MIN_IDLE_CONNS` | `redis.min_idle_conns` | `5` |
| `EZ_REDIS_POOL_SIZE` | `redis.pool_size` | `10` |

### Auth

| 环境变量 | 对应配置键 | 默认值 |
| --- | --- | --- |
| `EZ_AUTH_JWT_SECRET` | `auth.jwt_secret` | 开发占位密钥 |
| `EZ_AUTH_ACCESS_TOKEN_TTL` | `auth.access_token_ttl` | `7200` |
| `EZ_AUTH_ISSUER` | `auth.issuer` | `ez-admin` |

### Log

| 环境变量 | 对应配置键 | 默认值 |
| --- | --- | --- |
| `EZ_LOG_LEVEL` | `log.level` | `info` |
| `EZ_LOG_FORMAT` | `log.format` | `console` |
| `EZ_LOG_FILENAME` | `log.filename` | `logs/app.log` |
| `EZ_LOG_MAX_SIZE` | `log.max_size` | `100` |
| `EZ_LOG_MAX_BACKUPS` | `log.max_backups` | `7` |
| `EZ_LOG_MAX_AGE` | `log.max_age` | `30` |
| `EZ_LOG_COMPRESS` | `log.compress` | `false` |

### Upload

| 环境变量 | 对应配置键 | 默认值 |
| --- | --- | --- |
| `EZ_UPLOAD_DIR` | `upload.dir` | `uploads` |
| `EZ_UPLOAD_PUBLIC_PATH` | `upload.public_path` | `/uploads` |
| `EZ_UPLOAD_MAX_SIZE_MB` | `upload.max_size_mb` | `10` |
| `EZ_UPLOAD_ALLOWED_EXTS` | `upload.allowed_exts` | 图片、文档等白名单 |

## 当前默认部署主线为什么用 `127.0.0.1`

在第 9 章默认部署主线里：

- PostgreSQL 容器映射到 `127.0.0.1:5432`
- Redis 容器映射到 `127.0.0.1:6379`
- 后端直接跑在宿主机

所以 `.env` 里当前默认写的是：

- `EZ_DATABASE_HOST=127.0.0.1`
- `EZ_REDIS_HOST=127.0.0.1`

如果你切到全容器部署变体，这两个值就不一定还能继续沿用。

## 哪些值通常不建议频繁改

虽然所有环境变量都可以改，但下面这些值更应该谨慎：

- `EZ_AUTH_JWT_SECRET`
- `EZ_DATABASE_DRIVER`
- `EZ_DATABASE_NAME`
- `EZ_UPLOAD_PUBLIC_PATH`

原因不是“不能改”，而是它们更容易牵一发而动全身：

- Token 全量失效
- 数据连接切换
- 上传访问路径变化

## 快速排查建议

如果你怀疑线上问题和环境变量有关，优先先确认这三件事：

1. `.env` 是否真的被 systemd 服务加载
2. 新版本是否新增了必填环境变量
3. 当前部署形态下，`HOST` 类变量是不是还符合真实网络结构

## 相关教程页

- [第 9 章：环境变量与初始化数据](/tutorial/chapter-9/env-and-init-data)
- [第 9 章：部署验证与复用说明](/tutorial/chapter-9/deployment-and-reuse)
- [第 9 章：长期运维 FAQ](/tutorial/chapter-9/operations-maintenance-faq)
