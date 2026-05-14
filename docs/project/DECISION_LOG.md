# Decision Log — 架构决策记录

## ADR-001：不采用大面积 co-located unit tests

**状态：** 已采纳

**原因：** 用户明确要求业务代码目录保持干净，不接受散落大量 *_test.go。项目更重视集中式测试治理。

**决策：** 后端测试集中放在 server/tests，业务代码目录不新增测试文件。

---

## ADR-002：采用集中式 API / RBAC / Contract / E2E 测试

**状态：** 已采纳

**原因：** 后台系统的核心价值在于权限和数据权限的正确性，集中式黑盒测试更接近真实行为，更适合企业级后台底座。

**决策：**
- `server/tests/api/` — API 黑盒测试
- `server/tests/rbac/` — 权限和数据权限流程测试
- `server/tests/contract/` — OpenAPI 契约测试
- `server/tests/testutil/` — 测试辅助工具

---

## ADR-003：当前阶段不优先堆功能

**状态：** 已采纳

**原因：** 项目目标是从 8 分提升到 9 分以上的生产级质量。当前最大短板是测试可信度、安全基线、可观测性。

**决策：** Phase 1 只做测试体系做实。暂不做主题、国际化、WebSocket、审批流。

---

## ADR-004：测试数据通过 API 层创建

**状态：** 已采纳

**原因：** 测试应通过公开 API 行为验证，而不是直接操作数据库。这样既测试了 API 本身，又避免了测试与内部实现耦合。但对于 Casbin 策略，如果 API 不提供直接写入能力，可以回退到 DB 层操作。

**决策：**
- 测试用户、角色通过 API 创建（POST /system/users, POST /system/roles）
- Casbin 策略通过 API 赋权（POST /system/roles/:id/permissions）
- 如果 API 不满足需求，可在 testutil 中直接操作 DB

---

## ADR-005：Casbin 策略缓存需要显式刷新

**状态：** 已采纳

**原因：** Casbin enforcer 在启动时从 DB 加载策略到内存。运行时通过 API 更新权限后，写入 DB 但内存中的策略未刷新，导致新权限不立即生效。在测试中发现此问题。

**决策：**
- 在 `authz.Enforcer` 上新增 `ReloadPolicy()` 方法
- 测试中通过 `testutil.ReloadPolicies()` 在权限变更后手动刷新
- 生产环境可能需要考虑事件驱动的自动刷新机制（Phase 3 或后续）

---

## ADR-006：expandDepartmentTree GORM scope 泄漏（已修复）

**状态：** 已修复

**原因：** `UserQueryScope` 和 `DepartmentQueryScope` 的闭包捕获了外部传入的 `db *gorm.DB`，该实例可能已携带 GORM 链条件（如 `.Model(&model.User{})`、`.Where(...)` 等）。scope 闭包内的 `expandDepartmentTree(db, ...)` 在这个有状态的 `db` 上执行 `.Table("sys_department")` 查询时，原有条件被泄漏到 `sys_department` 表上，导致 `column "department_id" does not exist` 错误。

**影响范围：** 所有使用 `data_scope=dept_and_children` 的用户在查询用户列表时触发 500 错误。`data_scope=dept` 在 DepartmentQueryScope 下的子部门查询同样受影响。

**修复：** 在 `UserQueryScope` 和 `DepartmentQueryScope` 入口处，用 `newCleanSession(db)` 创建一个 `Session(&gorm.Session{NewDB: true})` 的干净 session。所有子查询（`expandDepartmentTree`、`accessibleDepartmentIDs`）使用此 cleanDB 执行，彻底隔离 GORM 链条件。修复后通过 `TestDataScopeDeptAndChildren` 真实测试验证。

---

## ADR-007：Refresh Token 存储在 Redis 而非数据库

**状态：** 已采纳

**原因：** Refresh token 需要高频读写和自动过期，Redis 的 TTL 原生支持比数据库定时清理更高效。项目已有 go-redis 集成。

**决策：**
- Refresh token 存储在 Redis（key: `refresh_token:{sha256(token)}`）
- 辅助索引 `user_sessions:{user_id}` 支持全量会话撤销
- Access token 黑名单也存储在 Redis（短 TTL = access token 剩余寿命）
- 不增加数据库表

---

## ADR-008：Access Token 保持在 Authorization Header，Refresh Token 在 HttpOnly Cookie

**状态：** 已采纳

**原因：**
- Access token 保持无状态（Authorization header）兼容现有 API 测试和客户端
- Refresh token 放在 HttpOnly Secure Cookie 中，前端 JS 无法读取，防止 XSS 窃取
- 前端通过 `withCredentials: true` 自动发送 cookie

**决策：**
- 登录返回：access token 在 JSON body，refresh token 在 Set-Cookie header
- 刷新返回：新 access token 在 JSON body，新 refresh token 在 Set-Cookie header（rotation）
- 登出：清除 cookie + 黑名单 access token

---

## ADR-009：TokenBlacklistChecker 通过可选接口注入 Auth Middleware

**状态：** 已采纳

**原因：** Auth middleware 签名需要向后兼容——测试环境中可能没有 Redis（contract tests）。

**决策：** `middleware.Auth(tokenManager, blacklist, log)` — blacklist 为 nil 时跳过黑名单检查。所有调用点传入 `*RefreshTokenStore`（nil-safe）。

---

## ADR-010：Kubernetes 健康探针分离 liveness 和 readiness

**状态：** 已采纳

**原因：** Kubernetes 需要区分进程存活（liveness）和服务就绪（readiness）。如果只提供一个探针检查依赖，当 DB/Redis 短暂不可用时 K8s 会重启 pod，而不是暂时摘除流量。

**决策：**
- `/healthz`（liveness）：仅检查进程存活，始终返回 200
- `/readyz`（readiness）：检查 DB + Redis 连通性，任一失败返回 503
- `/health` 保留为向后兼容，行为等同 `/readyz`

---

## ADR-011：Prometheus 指标使用 promauto 自动注册

**状态：** 已采纳

**原因：** 项目使用标准 Prometheus client_golang 库。promauto 包在 init 时自动注册指标到默认 registry，避免手动注册的样板代码。

**决策：**
- `http_requests_total`：CounterVec，按 method/path/status 分标签
- `http_request_duration_seconds`：HistogramVec，按 method/path 分标签
- 使用 `promhttp.Handler()` 暴露 `/metrics` 端点
- Metrics 中间件跳过 `/metrics` 路径自身，避免递归计数

---

## ADR-012：Casbin 策略变更后自动 ReloadPolicy

**状态：** 已采纳

**原因：** `UpdatePermissions` 和 `UpdateMenus` 只将新策略写入数据库，但 Casbin Enforcer 在启动时一次性加载策略到内存。运行时通过 API 更新权限后，内存策略不刷新，导致新权限不立即生效。E2E 测试中受限用户因 403 无法访问已授权的 API。

**决策：**
- 在 `role/application` 层定义 `PolicyReloader` 接口（`ReloadPolicy() error`）
- `Service` 注入 `PolicyReloader`，在 `UpdatePermissions` 事务提交后调用 `ReloadPolicy()`
- 通过依赖链 `routes → services → application` 传递 `*authz.Enforcer`
- `UpdateMenus` 不需要 ReloadPolicy（菜单权限不经过 Casbin）

---

## ADR-013：暗色模式采用 CSS 变量 + Naive UI darkTheme 双层切换

**状态：** 已采纳

**原因：** 前端已有的 CSS 变量体系（`--ez-*`）和 Naive UI `themeOverrides` 是两套独立的样式通道。暗色模式需要两者同步切换：CSS 变量控制自定义组件和 Tailwind 样式，Naive UI `darkTheme` 控制组件库内置样式。

**决策：**
- 新建 Pinia theme store（`admin/src/stores/theme.ts`），管理 light/dark/auto 三种模式
- CSS 变量层：在 `.dark` class 下覆盖所有 `--ez-*` 变量
- Naive UI 层：`n-config-provider` 条件传入 `darkTheme` + `darkThemeOverrides`
- 主题偏好持久化到 localStorage，auto 模式监听 `prefers-color-scheme` 系统偏好
- Header 月亮/太阳图标点击循环切换 light → dark → auto
- 全局硬编码颜色统一替换为 CSS 变量（约 23 个 Vue 文件）
- LoginPage 保持浅色设计，不纳入暗色模式

---

## ADR-014：国际化采用 vue-i18n + Pinia locale store

**状态：** 已采纳

**原因：** 前端所有用户可见文本为硬编码中文。引入 i18n 需要同时处理 vue-i18n 翻译和 Naive UI 组件库 locale 切换。

**决策：**
- 使用 vue-i18n（Composition API 模式，legacy: false）
- Pinia locale store（`admin/src/stores/locale.ts`）管理当前语言和 Naive UI locale 映射
- 翻译文件集中放在 `admin/src/i18n/locales/`（zh-CN.ts / en-US.ts）
- App.vue 的 NConfigProvider locale 绑定到 locale store 的 computed 属性
- 后端返回的动态数据（菜单标题、角色名称等）暂不翻译
- 默认语言 zh-CN，localStorage 持久化用户选择
