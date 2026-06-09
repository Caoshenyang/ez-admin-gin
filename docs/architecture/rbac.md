---
title: 权限体系
description: RBAC 角色权限、Casbin 接口授权、按钮权限、数据权限的完整链路
---

# 权限体系

EZ Admin 的权限体系由四个层级组成：认证、角色授权、接口权限、数据权限。

## 权限架构总览

```
用户登录 → JWT Token 签发
  ↓
请求携带 Token → Auth 中间件验证身份
  ↓
LoadActor 中间件 → 加载用户角色 + 菜单权限 + 按钮权限码
  ↓
Permission 中间件 → Casbin 策略匹配（角色 × URL × HTTP 方法）
  ↓
Repository 层 → datascope 注入数据过滤条件
```

## 认证（Authentication）

使用 JWT（HS256）进行身份认证。

**Token 结构：**

```go
type Claims struct {
    UserID   uint   `json:"user_id"`
    Username string `json:"username"`
    jwt.RegisteredClaims
}
```

**配置项：**

| 配置 | 默认值 | 说明 |
|------|-------|------|
| `auth.jwt_secret` | 开发用默认值 | 生产环境必须通过环境变量覆盖 |
| `auth.access_token_ttl` | 7200（2 小时） | Token 有效期 |
| `auth.issuer` | ez-admin | 签发方标识 |

**认证流程：**

1. 用户提交 username + password
2. 后端验证密码（bcrypt）
3. 签发 JWT Token
4. 前端存储 Token（localStorage / sessionStorage）
5. 后续请求通过 `Authorization: Bearer <token>` 携带

## 角色授权（RBAC）

采用经典的用户-角色-权限模型：

```
sys_user ──M:N── sys_role ──M:N── sys_menu
                       └──M:N── sys_api ──sync── casbin_rule
```

- 一个用户可以拥有多个角色
- 一个角色可以关联多个菜单/按钮权限
- 一个角色可以关联多个接口权限元数据
- 菜单和按钮权限通过 `sys_role_menu` 管理
- 接口权限通过 `sys_api` + `sys_role_api` 管理，再同步到 `casbin_rule` 给 Casbin 执行

**角色数据范围：** 每个角色可以配置数据权限作用域（详见[数据权限](#数据权限)）。

## Casbin 接口权限

使用接口元数据管理可选权限点，使用 Casbin 进行接口级（URL + HTTP 方法）权限控制。

**RBAC 模型（`configs/rbac_model.conf`）：**

```ini
[request_definition]
r = sub, obj, act

[policy_definition]
p = sub, obj, act

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = r.sub == p.sub && keyMatch2(r.obj, p.obj) && (r.act == p.act || p.act == "*")
```

- **sub**：角色编码（如 `super_admin`）
- **obj**：URL 路径模式（如 `/api/v1/system/users/:id/update`）
- **act**：HTTP 方法（GET/POST/PUT/DELETE），`*` 表示全部放行

**接口权限元数据存储在 `sys_api`，角色关联存储在 `sys_role_api`。** 保存角色接口权限时，后端会用角色关联重建该角色的 `casbin_rule`，并重新加载 Casbin 策略。

**权限匹配流程：**

1. Permission 中间件从 Actor 上下文获取用户角色列表
2. 构造 Casbin 请求 `(角色, URL路径, HTTP方法)`
3. 逐一匹配 Casbin 策略
4. 任一角色匹配即放行，全部不匹配返回 403

**管理流程：**

1. `sys_api` 维护接口名称、权限编码、路径和方法
2. 角色接口权限页从 `/api/v1/system/apis` 拉取接口元数据
3. 保存角色时写入 `sys_role_api`
4. 同步生成该角色的 `casbin_rule`

## 按钮权限

菜单表中的 `type=3` 记录为按钮级权限，每个按钮有唯一 `code`。

**权限码命名约定：** `模块:资源:操作`

```
system:user:list
system:user:create
system:user:update
system:user:delete
```

**前端消费方式：**

```typescript
import { usePermission } from '@/composables/usePermission'

const { canUse } = usePermission()

// 在模板中控制按钮显隐
<button v-if="canUse('system:user:create')">新建用户</button>
```

**后端保障：** 按钮权限对应的 API 接口同时在 Casbin 中注册，前端隐藏不代表绕过。

## 数据权限

数据权限在 Repository 层注入，控制用户能看到哪些数据。

### 五级数据作用域

| 作用域 | 值 | 含义 | 过滤条件 |
|--------|---|------|---------|
| all | 1 | 所有数据 | 不追加过滤 |
| dept | 2 | 本部门 | `department_id = 用户部门ID` |
| dept_and_children | 3 | 本部门及下级 | `department_id IN (本部门及子部门)` |
| self | 4 | 仅本人 | `creator_id = 用户ID` |
| custom_dept | 5 | 自定义部门 | `department_id IN (配置的部门列表)` |

### 数据权限注入

```go
// Repository 层调用
datascope.ApplyScopes(db, actor)
```

根据 Actor 的角色数据范围类型，自动注入对应的 GORM scope：

- **all** → 不追加条件
- **dept** → `WHERE department_id = actor.DepartmentID`
- **dept_and_children** → `WHERE department_id IN (部门子树 ID 列表)`
- **self** → `WHERE creator_id = actor.UserID`
- **custom_dept** → `WHERE department_id IN (sys_role_data_scope 配置的部门)`

### 相关数据表

| 表 | 用途 |
|----|------|
| `sys_role.data_scope` | 角色的数据作用域类型 |
| `sys_role_data_scope` | 自定义部门范围关联（多对多） |
| `sys_department.ancestors` | 部门祖先路径（逗号分隔 ID），用于子树查询 |

## Actor 上下文

`LoadActor` 中间件在每次请求时构建完整的用户上下文：

```go
type Actor struct {
    UserID         uint     // 用户 ID
    Username       string   // 用户名
    DepartmentID   uint     // 所属部门
    Roles          []string // 角色编码列表
    MenuIDs        []uint   // 授权菜单 ID
    ButtonCodes    []string // 按钮权限码
}
```

当用户拥有多个角色时，权限取并集：所有角色的菜单、按钮、数据范围合并计算。

## 中间件执行顺序

```
CORS → RequestID → Logger → Recovery → Auth → LoadActor → Permission → OperationLog → Handler
```

- **Auth**：验证 JWT，提取 userID + username
- **LoadActor**：查询用户完整信息，注入 Actor 到 gin.Context
- **Permission**：用 Actor 的角色做 Casbin 匹配
- **OperationLog**：记录 API 操作（需认证的路由）
