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

## ADR-007：双 Token 架构（Access Token + Refresh Token）

**状态：** 已采纳

**原因：** 单一 JWT Access Token 无法撤销，泄漏后攻击窗口长。后台系统需要会话管理能力（踢用户下线、密码修改后强制重新登录）。

**决策：**
- Access Token：短寿命 JWT（2h），通过 Authorization header 传递，保持不变
- Refresh Token：长寿命 opaque token（7d），SHA-256 哈希后存 Redis，通过 HttpOnly Secure Cookie 传递
- 登录时同时签发两种 token；Refresh 时旋转（吊销旧 token，签发新 token）
- 前端 401 拦截器自动调 /auth/refresh，成功后重试原请求，失败则清除会话跳登录页
- Cookie 配置：HttpOnly=true, Secure=true, SameSite=Lax, Path=/api/v1/auth

**不修改的部分：**
- Auth middleware 仍从 Bearer header 读 access token
- 不修改数据库 schema
- 现有测试无需改动（LoginAs 仍返回 access token）

---

## ADR-008：登录限流双层防护（IP + 账号锁定）

**状态：** 已采纳

**原因：** 纯 IP 限流无法防御针对特定用户名的暴力破解（攻击者可分散 IP）。同时旧实现返回英文裸 JSON，与项目 httpx envelope 规范不一致。

**决策：**
- 保留 IP 滑动窗口限流（第一道防线，防止高频请求）
- 新增 per-username 账号锁定（第二道防线）：连续 N 次失败后临时锁定用户名
- 锁定计数和锁定状态存 Redis，键名 `lockout:fail:{username}` / `lockout:blocked:{username}`
- 登录成功时清除失败计数
- 429 响应统一使用 httpx.Error + errorsx.TooManyRequests，中文错误消息
- maxFailures ≤ 0 时自动禁用锁定功能（防御配置缺失）
