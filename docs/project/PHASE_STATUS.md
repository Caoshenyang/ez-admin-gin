# Phase Status — 当前执行状态

## 评分

- 当前评分：约 9.2 / 10（Phase 1–5 完成，50 个后端测试 + 43 个 E2E 测试全部通过）
- 目标评分：9.0+

## 当前 Phase

**Phase 6：产品能力补强 — 进行中**

## 当前重点

- Phase 6 第一项「暗色模式 / 主题系统」已完成
- 下一步：Phase 6 其余项目（国际化、WebSocket 通知等）

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

### Phase 5: 发布治理和文档成熟

- [x] **SECURITY.md** — 安全策略和漏洞报告流程
- [x] **CONTRIBUTING.md** — 贡献指南（开发环境、代码风格、PR 流程、测试规范）
- [x] **Issue 模板** — `.github/ISSUE_TEMPLATE/bug_report.md` + `feature_request.md`
- [x] **PR 模板** — `.github/PULL_REQUEST_TEMPLATE.md`
- [x] **CHANGELOG.md** — 已有（Keep a Changelog 格式）
- [x] **生产检查清单** — 已有 `docs/deployment/production-checklist.md`
- [x] **测试策略文档** — 已有 `docs/project/TESTING_STRATEGY.md`
- [x] **架构决策记录** — 已有 `docs/project/DECISION_LOG.md`（ADR-001 ~ ADR-012）
- [x] **迁移指南** — 已有 `docs/backend/migration.md`
- [x] **Demo 环境** — 已有 `deploy/compose.local.yml`

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
- [x] **OpenAPI 生成前端类型** — openapi-typescript 从 swagger.json 生成 TypeScript 类型
  - admin/src/api/generated.ts — 自动生成的 paths + definitions 类型
  - admin/package.json — generate:api / check:api-types 脚本
  - Makefile — generate-types / check-types 目标
- [x] **CI 契约一致性检查** — frontend job 中检查 generated.ts 是否与 swagger.json 同步

### Phase 6: 产品能力补强

- [x] **暗色模式 / 主题系统**
  - admin/src/stores/theme.ts — Pinia 主题 store（light/dark/auto + localStorage + 系统偏好）
  - admin/src/styles/main.css — `.dark` 暗色 CSS 变量（surface/text/border/status 全覆盖）
  - admin/src/styles/theme.ts — darkThemeOverrides（Naive UI 暗色组件主题）
  - admin/src/App.vue — 条件切换 darkTheme + dark class
  - admin/src/components/app-shell/AppHeader.vue — 月亮/太阳图标主题切换（light → dark → auto 循环）
  - 全局硬编码颜色替换为 CSS 变量（~23 个 Vue 文件）

## 未完成

- Phase 6 国际化 — **基础设施已搭建，逐文件替换进行中**：
  - [x] vue-i18n 安装 + 插件注册
  - [x] admin/src/i18n/index.ts（createI18n 实例）
  - [x] admin/src/i18n/locales/zh-CN.ts（完整中文翻译）
  - [x] admin/src/i18n/locales/en-US.ts（完整英文翻译）
  - [x] admin/src/stores/locale.ts（Pinia locale store + Naive UI locale 映射）
  - [x] admin/src/main.ts（注册 i18n 插件）
  - [x] admin/src/App.vue（Naive UI locale 绑定到 locale store）
  - [ ] 逐文件替换硬编码中文 → t() 调用（~60 个文件待处理）
  - [ ] AppHeader 添加语言切换按钮
  - [ ] 验证（lint + type-check + build）
- Phase 6 其余项目：
  - WebSocket 通知
  - 审批工作流
  - 业务模板
  - 模块生成器

## 当前下一步

1. 继续国际化逐文件替换：先从布局组件（AppHeader、AppSidebar、WorkTabs、AdminLayout）开始
2. 然后替换 auth / iam / system 模块中的硬编码文本
3. 最后验证并完成

## 阻塞点

- 无

## 最近一次执行记录

- **日期：** 2026-05-14
- **修改内容：**
  - **Phase 6 国际化基础设施搭建：**
    - 安装 vue-i18n 11.4.2
    - 新建 admin/src/i18n/index.ts — createI18n 实例（legacy: false，默认 zh-CN）
    - 新建 admin/src/i18n/locales/zh-CN.ts — 完整中文翻译（common / layout / auth / iam / system / theme / error 全覆盖）
    - 新建 admin/src/i18n/locales/en-US.ts — 完整英文翻译
    - 新建 admin/src/stores/locale.ts — Pinia locale store（localStorage 持久化，Naive UI locale/dateLocale 映射）
    - 修改 admin/src/main.ts — 注册 i18n 插件
    - 修改 admin/src/App.vue — Naive UI locale 从 locale store 动态获取（替换硬编码 zhCN）
- **剩余工作：**
  - ~60 个 Vue/TS 文件中的硬编码中文需替换为 t() 调用
  - AppHeader 添加语言切换 UI
