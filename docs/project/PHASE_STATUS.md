# Phase Status — 当前执行状态

## 评分

- 当前评分：约 9.2 / 10（Phase 1–4 核心完成，50 个后端测试 + 43 个 E2E 测试全部通过）
- 目标评分：9.0+

## 当前 Phase

**Phase 4：前端质量和 E2E — 进行中**

## 当前重点

- Phase 4 核心测试用例全部完成
- 43 个 E2E 测试全部通过
- 下一步：OpenAPI 生成前端类型 + CI 检查前后端契约一致

## 已完成

### Phase 1: 集中式测试做实

- [x] Makefile 统一开发入口
- [x] GitHub Actions CI
- [x] 集中式 server/tests 目录结构
- [x] 34 个测试全部通过（详见 git history）

### Phase 2: 安全基线升级

- [x] **Refresh Token 存储层** — server/internal/platform/authn/refresh.go
  - Redis-backed RefreshTokenStore: Create/Verify/Revoke/RevokeAllForUser
  - Token 黑名单: BlacklistAccessToken/IsBlacklisted
  - SHA-256 hash key, crypto/rand token 生成
- [x] **双 Token 签发** — Login 返回 access token + HttpOnly Secure refresh cookie
  - server/internal/platform/authn/authn.go — SetRefreshStore, GenerateRefreshToken
  - server/internal/modules/auth/application/login_service.go — 返回双 token
  - server/internal/modules/auth/api/handlers.go — LoginWithRefresh + setRefreshTokenCookie
- [x] **POST /auth/refresh** — Refresh token rotation
  - server/internal/modules/auth/application/refresh_service.go
  - server/internal/modules/auth/api/refresh_handler.go
- [x] **POST /auth/logout** — 服务端会话撤销
  - server/internal/modules/auth/application/logout_service.go
  - server/internal/modules/auth/api/logout_handler.go
  - 撤销 refresh token + 黑名单 access token
- [x] **Auth middleware 黑名单检查** — 可选 TokenBlacklistChecker
- [x] **安全响应头** — server/internal/platform/middleware/security_headers.go
  - X-Content-Type-Options, X-Frame-Options, Referrer-Policy, X-XSS-Protection
  - 生产环境 Strict-Transport-Security
- [x] **CORS Vary: Origin** — 防止缓存问题
- [x] **生产配置强校验** — ValidateProduction 扩展
  - JWT secret, CORS origins, Swagger, 数据库密码
- [x] **上传安全增强** — ValidateFilename 拒绝双重扩展名/路径穿越
- [x] **登录限流增强** — AccountLockChecker 账户级锁定
- [x] **前端双 Token 适配**
  - admin/src/api/http.ts — withCredentials + 静默刷新
  - admin/src/modules/auth/api/auth.ts — logout API
  - admin/src/layouts/AdminLayout.vue — 后端 logout 调用

### Phase 3: 可观测性和运维能力

- [x] **X-Request-ID** — 已有（Phase 1 中实现）
  - server/internal/platform/middleware/requestid.go — UUID 生成或传递
  - 已注入 context + 响应头 + GinLogger 日志字段
- [x] **结构化日志** — 已有（Phase 1 中实现）
  - server/internal/platform/logger/logger.go — Zap + JSON 格式支持
  - GinLogger/GinRecovery 自动记录 request_id, method, path, status, latency
- [x] **错误码规范** — 已有（Phase 1 中实现）
  - server/internal/pkg/errorsx/errors.go — 统一错误码 + HTTP 状态码映射
- [x] **Kubernetes 健康探针**
  - GET /healthz — liveness probe（进程存活即返回 200，无依赖检查）
  - GET /readyz — readiness probe（检查 DB + Redis 连通性）
  - GET /health — 向后兼容（等同 readyz）
  - server/internal/modules/system/health_handler.go — Liveness/Readiness/Check
- [x] **Prometheus 指标**
  - server/internal/platform/middleware/metrics.go — Gin 中间件
  - http_requests_total（counter: method, path, status）
  - http_request_duration_seconds（histogram: method, path）
  - GET /metrics — Prometheus text format 指标暴露
  - bootstrap/router.go — 中间件注册 + promhttp.Handler

### Phase 4: 前端质量和 E2E

- [x] Playwright 基础设施（playwright.config.ts + e2e 目录 + fixtures）
- [x] 登录流程 E2E 测试（5 个用例）
- [x] 菜单权限 E2E 测试（8 个用例）
- [x] 用户管理 E2E 测试（8 个用例）
- [x] 按钮权限 E2E 测试（7 个用例）
- [x] 角色授权 E2E 测试（9 个用例）
- [x] 无权限页面 E2E 测试（3 个用例）
- [x] Token 过期处理 E2E 测试（3 个用例）
- [x] **Casbin 策略自动刷新** — UpdatePermissions 后调用 ReloadPolicy()

## 未完成

- [ ] OpenAPI 生成前端类型
- [ ] CI 检查前后端契约一致

## 当前下一步

1. OpenAPI 生成前端类型
2. CI 检查前后端契约一致

## 阻塞点

- 无

## 最近一次执行记录

- **日期：** 2026-05-13
- **修改内容：**
  - **修复 Bug：** Casbin 策略更新后未自动 ReloadPolicy，导致运行时权限变更不生效
    - server/internal/modules/iam/role/application/service.go — 注入 PolicyReloader，UpdatePermissions 后调用 ReloadPolicy()
    - server/internal/modules/iam/role/application/ports.go — 新增 PolicyReloader 接口
    - server/internal/modules/iam/role/services.go — ServiceOptions 增加 Enforcer 字段
    - server/internal/modules/iam/role/routes.go — RouteOptions 增加 Enforcer 字段
    - server/internal/modules/iam/routes.go — 向 role 模块传递 Enforcer
  - **修复 E2E 测试：** 角色状态切换断言从不存在的 "成功" 改为匹配实际消息
    - admin/e2e/iam/role.spec.ts — 修复 getByText(/成功/) → getByText(/已禁用|已启用/)
  - **新增 E2E 测试：**
    - admin/e2e/iam/no-permission.spec.ts（3 个用例：侧边栏无权限菜单、URL 重定向、API 403）
    - admin/e2e/iam/token-expired.spec.ts（3 个用例：过期 token 跳转、正常访问、清除 token）
  - **Playwright 配置：** 移除 channel: 'chrome' 使用 Playwright 内置 Chromium
- **测试结果：**
  - 后端 contract tests 通过
  - E2E 测试 43/43 全部通过（37.9s）
- **剩余风险：**
  - E2E 测试依赖后端 + 前端同时运行
  - 残留测试数据需定期清理（E2E 前缀的用户、角色和菜单）
