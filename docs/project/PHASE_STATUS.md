# Phase Status — 当前执行状态

## 评分

- 当前评分：约 8.5 / 10（custom_dept datascope 测试完成，OpenAPI 契约增强，t.Skip 从 2 减至 1）
- 目标评分：9.0+

## 当前 Phase

**Phase 1：集中式测试做实**

## 当前重点

- TestMultiRolePermissionUnion — 唯一剩余 t.Skip
- Phase 1 收尾检查

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
- [x] server/tests/rbac/permission_flow_test.go — 12 个真实测试
  - API 权限测试 (5): 401/403/200/200/method differentiation
  - Datascope 测试 (7): all/self/dept/dept_and_children/custom_dept/default_deny + dept_and_children 已修复
- [x] **datascope.go GORM scope 泄漏 bug 已修复** — newCleanSession 隔离子查询
- [x] server/internal/platform/authz/authz.go — ReloadPolicy() 方法
- [x] CI integration job 已包含 RBAC 测试

## 未完成

- [ ] TestMultiRolePermissionUnion — t.Skip（阻塞：API 不支持追加角色）
- [ ] E2E 测试（Phase 4）

## 已修复的业务代码 Bug

1. **expandDepartmentTree GORM scope 泄漏** — 已通过 newCleanSession 修复
   - 根因：UserQueryScope/DepartmentQueryScope 的闭包捕获的 db 参数携带 GORM 链条件
   - 修复：在 scope 函数外层创建 cleanDB = newCleanSession(db)，所有子查询用 cleanDB
   - 影响：修复前 data_scope=dept_and_children 的用户查询会触发 500 错误

2. **Casbin 策略缓存不自动刷新** — 已通过 ReloadPolicy() 规避

## 当前下一步

1. TestMultiRolePermissionUnion（唯一剩余 t.Skip，需要追加角色 API 或直接 DB 写入）
2. Phase 1 收尾检查，准备进入 Phase 2

## 阻塞点

- multi-role 需要一个追加角色而非重置角色的 API（或直接写 DB sys_user_role）

## 最近一次执行记录

- **日期：** 2026-05-11
- **修改内容：**
  - 新增 server/tests/testutil/app.go — SeedCustomDeptUser 方法
  - 实现 TestDataScopeCustomDept 真实测试（替代 t.Skip）
  - 新增 server/tests/contract/openapi_contract_test.go — 3 个契约增强测试
    - TestResponseSchemaEnvelope: 验证所有 200 响应引用 httpx.Body 统一信封
    - TestDefinitionsReachable: 验证所有 $ref 引用指向已定义的 definition
    - TestKeyEndpointDataSchemas: 验证关键端点有类型化的 data schema
- **测试结果：**
  - make test-contract: 7/7 PASS
  - make test-api: 21/21 PASS（auth 4 + user 6 + role 6 + menu 5）
  - make test-rbac: 12/12 PASS, 1 SKIP（multi-role）
  - make test-integration: 33/33 PASS, 1 SKIP
  - go vet ./...: 无错误
- **剩余风险：**
  - multi-role 测试仍阻塞（1 个 t.Skip）
