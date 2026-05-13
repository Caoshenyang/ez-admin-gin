# Testing Strategy — 测试策略

## 测试目标

- 防止权限回归（RBAC）
- 防止数据权限越权（datascope）
- 防止登录认证回归（auth）
- 防止 OpenAPI 文档和真实接口脱节（contract）
- 防止核心管理流程损坏（user/role/menu CRUD）

## 测试类型

| 类型 | 目录 | 运行方式 | 需要 DB/Redis |
|------|------|----------|--------------|
| API 黑盒测试 | server/tests/api/ | `make test-api` | 是 |
| RBAC 流程测试 | server/tests/rbac/ | `make test-rbac` | 是 |
| datascope 流程测试 | server/tests/rbac/ | `make test-rbac` | 是 |
| OpenAPI 契约测试 | server/tests/contract/ | `make test-contract` | 否 |
| E2E 测试 | tests/e2e 或 admin/e2e | `make test-e2e` | Phase 4 |

## 目录规则

```
server/tests/
├── api/              # API 黑盒测试（auth, user, role, menu 等）
├── rbac/             # 权限和数据权限流程测试
├── contract/         # OpenAPI 契约测试
└── testutil/         # 测试辅助工具（TestApp, seed, reset, cleanup）
    └── testdata/     # 测试配置文件
```

## 不采用

- 大面积业务目录旁边 *_test.go
- 为每个 service 创建机械测试
- 为覆盖率测试无意义 getter/setter
- mock 数据库（集成测试用真实 DB）

## 测试数据原则

- 不连接生产数据库
- 不依赖开发者本机旧数据
- 测试数据可重复初始化（通过 testutil seed）
- 测试之间不依赖执行顺序
- 测试用户、角色、部门、权限必须明显是测试数据
- 每个测试调用 CleanupTestData() 清理上一次测试残留

## testutil 能力清单

| 方法 | 说明 |
|------|------|
| NewTestApp | 启动完整测试应用（DB/Redis/httptest） |
| SeedAdmin | 创建初始 admin 用户 |
| SeedRestrictedUser | 创建无 API 权限的受限用户 |
| SeedUserWithPermissions | 创建有特定 API 权限的用户 |
| SeedScopedUser | 创建带指定 data_scope 和部门的用户 |
| SeedCustomDeptUser | 创建 data_scope=custom_dept 的用户（指定自定义部门列表） |
| SeedDepartment | 通过 API 创建部门 |
| CleanupTestData | 清理测试产生的角色/用户/部门/菜单/Casbin 策略 |
| ReloadPolicies | 刷新 Casbin 内存策略 |
| LoginAs | 以指定用户登录并返回 token |
| LoginWithCookies | 以指定用户登录并返回 token + cookies（含 refresh token） |
| AuthRequest | 构造带 Bearer token 的请求 |
| DecodeResponse | JSON 响应解码 |

## CI 原则

- contract tests 必须跑（不需要 DB）
- api tests 必须跑（CI 中 PostgreSQL + Redis service container）
- rbac tests 必须跑（已加入 CI integration job）
- e2e smoke 准备好后再加入
- 不成熟测试不得伪装成完成

## 当前测试清单

### 真实可跑（非 t.Skip）— 共 50 个后端 + 43 个 E2E（全部通过）

**admin/e2e/auth/login.spec.ts (5) — 已验证:**
- Login Flow: redirects to login page when not authenticated
- Login Flow: shows error on wrong password
- Login Flow: logs in successfully and redirects to dashboard
- Login Flow: already logged in user is redirected to dashboard from /login
- Login Flow: form validation shows errors for empty fields

**admin/e2e/iam/menu.spec.ts (8) — 已验证:**
- Menu Permission: displays menu management page with correct header
- Menu Permission: shows create root directory button for admin
- Menu Permission: displays menu table with seed data columns
- Menu Permission: shows seed menu items in table
- Menu Permission: creates a new root directory menu
- Menu Permission: opens edit modal with correct data
- Menu Permission: deletes a menu item
- Menu Permission: action buttons are visible for admin user

**admin/e2e/iam/user.spec.ts (8) — 已验证:**
- User Management: displays user management page with correct header
- User Management: shows create user button for admin
- User Management: displays user table with columns
- User Management: shows admin user in the list
- User Management: creates a new user
- User Management: opens edit modal with existing data
- User Management: toggles user status via API
- User Management: action buttons visible for admin

**admin/e2e/iam/button-permission.spec.ts (7) — 已验证:**
- Button Permission / Admin: admin sees create role button
- Button Permission / Admin: admin sees edit button on role cards
- Button Permission / Admin: admin sees status toggle on non-super-admin roles
- Button Permission / Admin: admin sees save permission button
- Button Permission / Restricted: user without create permission does not see create button
- Button Permission / Restricted: user without update permission does not see edit button
- Button Permission / Restricted: user without status permission does not see status toggle

**admin/e2e/iam/role.spec.ts (9) — 已验证:**
- Role Authorization: displays role page with correct header
- Role Authorization: shows super_admin role in role list
- Role Authorization: super admin role shows protected tag
- Role Authorization: creates a new role
- Role Authorization: opens edit modal with existing data
- Role Authorization: toggles role status via UI
- Role Authorization: permission panel shows menu tree for selected role
- Role Authorization: assigns menu permissions to a role
- Role Authorization: adds API permission row

**admin/e2e/iam/no-permission.spec.ts (3) — 已验证:**
- No Permission Page: restricted user does not see unauthorized menu in sidebar
- No Permission Page: restricted user navigating to unauthorized route is redirected to dashboard
- No Permission Page: API request without permission shows error message

**admin/e2e/iam/token-expired.spec.ts (3) — 已验证:**
- Token Expiration: expired access token triggers redirect to login page
- Token Expiration: valid token allows normal page access
- Token Expiration: removing token redirects to login on next navigation

**server/tests/api/auth_api_test.go (4):**
- TestLoginSuccess
- TestLoginWrongPassword
- TestUnauthorizedAccessWithoutToken
- TestUnauthorizedAccessWithInvalidToken

**server/tests/api/auth_refresh_test.go (6):**
- TestRefreshSuccess
- TestRefreshWithInvalidToken
- TestRefreshWithoutCookie
- TestRefreshRotation
- TestLogoutSuccess
- TestLogoutRevokesAccessToken

**server/tests/api/user_api_test.go (6):**
- TestCreateUser
- TestCreateUserDuplicateUsername
- TestListUsers
- TestUpdateUser
- TestUpdateUserStatus
- TestCreateUserMissingFields

**server/tests/api/role_api_test.go (6):**
- TestCreateRole
- TestCreateRoleDuplicateCode
- TestListRoles
- TestUpdateRole
- TestUpdateRolePermissions
- TestUpdateRoleStatus

**server/tests/api/menu_api_test.go (5):**
- TestCreateMenu
- TestListMenus
- TestUpdateMenu
- TestDeleteMenu
- TestCreateMenuMissingFields

**server/tests/api/health_api_test.go (3):**
- TestHealthzReturnsOK
- TestReadyzReturnsOK
- TestMetricsEndpoint

**server/tests/contract/openapi_contract_test.go (7):**
- TestSwaggerFileExists
- TestSwaggerParsable
- TestKeyEndpointsDeclared
- TestSwaggerInfo
- TestResponseSchemaEnvelope
- TestDefinitionsReachable
- TestKeyEndpointDataSchemas

**server/tests/rbac/permission_flow_test.go (13):**
- API 权限测试 (5):
  - TestUnauthenticatedAccessToSystemEndpoint
  - TestPermissionDeniedWithoutRole
  - TestPermissionGrantedWithRole
  - TestAdminCanAccessAllEndpoints
  - TestHTTPMethodPermissionDifferentiation
- Datascope 测试 (7):
  - TestDataScopeAll
  - TestDataScopeSelf
  - TestDataScopeDept
  - TestDataScopeDeptAndChildren
  - TestDataScopeCustomDept
  - TestDataScopeDefaultDeny
- 多角色权限联合测试 (1):
  - TestMultiRolePermissionUnion

### t.Skip / TODO — 0 个

（Phase 1 全部完成，无剩余 t.Skip）

### 已修复的业务代码 Bug

1. **expandDepartmentTree GORM scope 泄漏** — 已通过 newCleanSession 修复（ADR-006）
2. **Casbin 策略缓存不自动刷新** — 已通过 ReloadPolicy() 规避（ADR-005）
