# Phase Status — 当前执行状态

## 评分

- 当前评分：约 9.0 / 10（Phase 1 + Phase 2 + Phase 3 完成，50 个后端测试 + 5 个 E2E 测试全部通过）
- 目标评分：9.0+

## 当前 Phase

**Phase 4：前端质量和 E2E — 进行中**

## 当前重点

- Phase 4 基础设施搭建完毕（Playwright + E2E 目录结构）
- 登录流程 E2E 测试已验证通过（5/5 PASS）
- 下一步：菜单权限 / 按钮权限 / 用户管理 E2E 测试

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

## 未完成

- [x] Playwright 基础设施（playwright.config.ts + e2e 目录 + fixtures）
- [x] 登录流程 E2E 测试（5 个用例）
- [ ] 菜单权限 E2E 测试
- [ ] 按钮权限 E2E 测试
- [ ] 用户管理 E2E 测试
- [ ] 角色授权 E2E 测试
- [ ] 无权限页面 E2E 测试
- [ ] Token 过期处理 E2E 测试
- [ ] OpenAPI 生成前端类型
- [ ] CI 检查前后端契约一致

## 当前下一步

1. 编写菜单权限 E2E 测试
2. 编写用户管理 E2E 测试

## 阻塞点

- 无

## 最近一次执行记录

- **日期：** 2026-05-12
- **修改内容：**
  - Phase 4 启动：Playwright E2E 基础设施搭建 + 登录流程测试
  - 新增: admin/playwright.config.ts (Playwright 配置，使用系统 Chrome channel)
  - 新增: admin/e2e/fixtures.ts (共享测试工具：loginViaApi, clearAuth, authedPage fixture)
  - 新增: admin/e2e/auth/login.spec.ts (登录流程 5 个 E2E 测试用例)
  - 修改: Makefile test-e2e target (调用 pnpm exec playwright test)
  - 新增依赖: @playwright/test 1.60.0
- **测试结果：**
  - make test-contract: 7/7 PASS（未改动后端）
  - make test-integration: 50/50 PASS, 0 SKIP
  - pnpm type-check: PASS
  - pnpm lint: PASS
  - playwright test: 5/5 PASS（登录流程 E2E）
- **剩余风险：**
  - E2E 测试依赖后端 + 前端同时运行
  - Chromium 下载因网络问题失败，已改用系统 Chrome channel
