---
title: 中间件
description: HTTP 中间件栈、执行顺序和各中间件的职责
---

# 中间件

## 中间件执行顺序

```
Request → CORS → RequestID → Logger → Recovery → [Auth → LoadActor → Permission] → [OperationLog] → Handler
```

方括号内的中间件仅在需要认证的路由上执行。

## 全局中间件

### CORS

跨域资源共享，配置在 `cors.allowed_origins` 中。

- 开发模式下 `localhost:*` 自动放行
- 生产环境必须配置实际前端域名

### RequestID

为每个请求生成唯一 ID，写入响应头 `X-Request-ID`，用于日志追踪。

### Logger

结构化请求日志，记录请求方法、路径、状态码、耗时。

- 开发模式：console 格式
- 生产模式：JSON 格式

### Recovery

Panic 恢复，防止未处理的 panic 导致服务崩溃。捕获后返回 500 错误并记录日志。

## 认证路由中间件

### Auth

JWT Token 验证：

1. 从 `Authorization: Bearer <token>` 提取 Token
2. 解析 JWT Claims
3. 提取 UserID 和 Username
4. 注入 gin.Context

未携带 Token 或 Token 无效时返回 401。

### LoadActor

加载用户完整上下文：

1. 从 Context 获取 UserID
2. 查询用户信息（部门、状态）
3. 查询用户角色列表
4. 查询角色关联的菜单 ID
5. 查询角色关联的按钮权限码
6. 构建 Actor 对象注入 Context

### Permission

Casbin 策略匹配：

1. 从 Actor 获取角色列表
2. 构造 Casbin 请求 `(角色, URL, HTTP方法)`
3. 遍历角色逐一匹配
4. 任一角色匹配即放行
5. 全部不匹配返回 403

## 特定路由中间件

### RateLimit

登录接口限流，防止暴力破解。

- 默认：每 IP 60 秒内最多 10 次登录请求
- 超限返回 429 Too Many Requests
- 配置项：`rate_limit.login_max_requests`、`rate_limit.login_window_sec`

### OperationLog

操作日志记录，记录 API 请求的详细信息：

- 操作人、请求方法、请求路径、请求参数
- 响应状态码、IP 地址、User-Agent
- 记录到 `sys_operation_log` 表

## 中间件源码位置

```
server/internal/platform/middleware/
├── actor.go           Actor 上下文
├── auth.go            JWT 认证
├── cors.go            跨域
├── operation_log.go   操作日志
├── permission.go      权限校验
├── ratelimit.go       限流
└── requestid.go       请求 ID
```
