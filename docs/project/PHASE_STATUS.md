# Phase Status — 当前执行状态

## 评分

- 当前评分：约 8.8 / 10（Phase 1 全部完成，0 t.Skip，34 个真实测试全通过）
- 目标评分：9.0+

## 当前 Phase

**Phase 2：安全基线升级 — 已完成**

## 当前重点

- Phase 2 所有子任务已完成
- 下一步：Phase 3（可观测性和运维能力）

## 已完成

- [x] Makefile 统一开发入口
- [x] GitHub Actions CI（backend + integration + frontend + docker + security）
- [x] 集中式 server/tests 目录结构
- [x] server/tests/testutil/app.go — TestApp 完整基础设施
- [x] server/tests/testutil/app.go — SeedAdmin, SeedRestrictedUser, SeedUserWithPermissions
- [x] server/tests/testutil/app.go — SeedScopedUser, SeedDepartment
- [x] server/tests/testutil/app.go — CleanupTestData, ReloadPolicies, DecodeResponse
- [x] server/tests/testutil/app.go — CleanupTestData 支持 sys_menu 清理
- [x] server/tests/testutil/app.go — SeedCustomDeptUser（custom_department_ids via API）
- [x] server/tests/api/auth_api_test.go — 4 个真实测试
- [x] server/tests/api/user_api_test.go — 6 个真实测试
  - TestCreateUser, TestCreateUserDuplicateUsername, TestListUsers
  - TestUpdateUser, TestUpdateUserStatus, TestCreateUserMissingFields
- [x] server/tests/api/role_api_test.go — 6 个真实测试
  - TestCreateRole, TestCreateRoleDuplicateCode, TestListRoles
  - TestUpdateRole, TestUpdateRolePermissions, TestUpdateRoleStatus
- [x] server/tests/api/menu_api_test.go — 5 个真实测试
  - TestCreateMenu, TestListMenus, TestUpdateMenu, TestDeleteMenu, TestCreateMenuMissingFields
- [x] server/tests/contract/openapi_contract_test.go — 7 个真实测试
  - TestSwaggerFileExists, TestSwaggerParsable, TestKeyEndpointsDeclared, TestSwaggerInfo
  - TestResponseSchemaEnvelope, TestDefinitionsReachable, TestKeyEndpointDataSchemas
- [x] server/tests/rbac/permission_flow_test.go — 13 个真实测试
  - API 权限测试 (5): 401/403/200/200/method differentiation
  - Datascope 测试 (7): all/self/dept/dept_and_children/custom_dept/default_deny
  - 多角色权限合并 (1): TestMultiRolePermissionUnion
- [x] **datascope.go GORM scope 泄漏 bug 已修复** — newCleanSession 隔离子查询
- [x] server/internal/platform/authz/authz.go — ReloadPolicy() 方法
- [x] CI integration job 已包含 RBAC 测试
- [x] **test-integration Makefile 修复** — 添加 -p 1 防止并行包执行导致 DB 冲突

## Phase 2 进度

- [x] 安全响应头中间件（X-Content-Type-Options, X-Frame-Options, CSP, Referrer-Policy, Permissions-Policy）
- [x] 生产配置强校验增强（CORS origins、Swagger、上传大小限制）
- [x] CORS 生产校验增强（禁止通配符 * + credentials 组合）
- [x] 上传安全增强（MIME 类型嗅探校验，交叉验证文件内容与扩展名）
- [x] Access Token + Refresh Token 双 token
- [x] HttpOnly Secure Cookie
- [x] Refresh Token rotation
- [x] 服务端会话撤销
- [x] 登录限流增强（账号锁定 + IP 限流响应标准化）

## 未完成

- [ ] E2E 测试（Phase 4）

## 已修复的业务代码 Bug

1. **expandDepartmentTree GORM scope 泄漏** — 已通过 newCleanSession 修复
   - 根因：UserQueryScope/DepartmentQueryScope 的闭包捕获的 db 参数携带 GORM 链条件
   - 修复：在 scope 函数外层创建 cleanDB = newCleanSession(db)，所有子查询用 cleanDB
   - 影响：修复前 data_scope=dept_and_children 的用户查询会触发 500 错误

2. **Casbin 策略缓存不自动刷新** — 已通过 ReloadPolicy() 规避

3. **test-integration 并行包执行导致 TestDataScopeCustomDept 偶发失败** — 已通过 -p 1 修复
   - 根因：Go 默认并行运行不同 test package，API 和 RBAC 测试共享同一数据库
   - 修复：Makefile test-integration 目标添加 -p 1 强制串行执行

## 当前下一步

1. Phase 3：可观测性和运维能力

## 阻塞点

- 无

## 最近一次执行记录

- **日期：** 2026-05-11
- **修改内容（Phase 2 登录限流增强）：**
  - 修改 `server/internal/pkg/errorsx/errors.go` — 新增 CodeTooManyRequests + TooManyRequests() 构造函数
  - 修改 `server/internal/platform/middleware/ratelimit.go` — IP 限流 429 响应标准化 + 新增 IsUsernameLocked/RecordLoginFailure/ClearLoginFailures
  - 修改 `server/internal/platform/config/config.go` — RateLimitConfig 增加 LoginLockoutThreshold/LoginLockoutSec
  - 修改 `server/internal/modules/auth/api/handlers.go` — LoginHandler 增加账号锁定检查、失败计数、成功清除
  - 修改 `server/internal/modules/auth/api/routes.go` — RouteOptions 增加 LockoutMaxFailures/LockoutSec
  - 修改 `server/internal/modules/auth/routes.go` — 传递锁定配置
  - 修改 `server/tests/testutil/testdata/config.yaml` — 增加 lockout 配置
  - 修改 `server/configs/config.yaml` — 增加 lockout 默认值
- **测试结果：**
  - make test-contract: 7/7 PASS
  - make test-integration: 34/34 PASS（0 SKIP）
  - go vet ./...: 无错误
