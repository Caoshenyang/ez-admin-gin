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
