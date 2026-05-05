---
title: RBAC 角色权限模型
description: "把用户、角色、接口权限、菜单权限和数据权限串成一条企业级后台可复用的授权主线。"
---

# RBAC 角色权限模型

前一章已经把登录、Token 和登录态接进来了。接下来要补的，不再只是一个简单的 `role_id` 字段，而是一套能支撑企业级后台继续演进的角色模型。

::: tip 🎯 本节目标
读完这一节，你应该能说清三件事：

- 为什么用户和角色要用关系表做多对多绑定
- 为什么角色既承接接口权限，也承接菜单权限和数据权限
- 为什么当前项目把角色模型放进 `iam/role`，而不是继续留在零散的认证代码里
:::

## 先看最终版授权关系

当前主线里的角色模型，已经不是“用户绑定一个角色名”这么简单，而是下面这条完整链路：

```text
当前登录用户
  ↓
sys_user_role
  ↓
sys_role
  ├─ casbin_rule         → 决定接口能不能访问
  ├─ sys_role_menu       → 决定菜单和按钮能不能看到
  └─ sys_role_data_scope → 决定数据范围怎么过滤
```

这也是为什么当前仓库在接口权限体系里，不能只停留在 `sys_role` 和 `sys_user_role` 两张表。

## 当前代码落点

这一节现在主要对应下面这些位置：

```text
server/
├─ internal/
│  ├─ middleware/
│  │  └─ actor.go
│  ├─ model/
│  │  ├─ role.go
│  │  ├─ user_role.go
│  │  └─ role_data_scope.go
│  ├─ module/
│  │  ├─ auth/
│  │  │  └─ dto.go
│  │  └─ iam/
│  │     └─ role/
│  │        ├─ dto.go
│  │        ├─ repository.go
│  │        ├─ service.go
│  │        ├─ handler.go
│  │        └─ routes.go
│  └─ platform/
│     └─ datascope/
│        └─ datascope.go
└─ migrations/
   ├─ mysql/
   └─ postgres/
```

| 位置 | 职责 |
| --- | --- |
| `model/role.go` | 定义角色基础字段，包括 `data_scope` |
| `model/user_role.go` | 定义用户与角色的多对多关系 |
| `model/role_data_scope.go` | 定义角色与“自定义部门范围”的绑定关系 |
| `middleware/actor.go` | 在请求期加载当前登录人的角色编码和数据范围摘要 |
| `module/iam/role/*` | 提供角色管理、接口权限维护、菜单权限维护 |
| `platform/datascope/datascope.go` | 把多角色数据范围合并成可复用的查询规则 |

## 角色模型为什么已经不是一张简单字典表

当前 `sys_role` 最关键的字段，不只是在 UI 上展示角色名称，而是要稳定承接后续三类授权能力：

```go
type Role struct {
	ID        uint
	Code      string
	Name      string
	Sort      int
	DataScope datascope.Scope
	Status    RoleStatus
	Remark    string
}
```

每个字段真正承担的职责是：

| 字段 | 作用 |
| --- | --- |
| `code` | 角色稳定标识，Casbin 接口策略直接引用它 |
| `name` | 给管理台和运营人员看的展示名 |
| `sort` | 控制角色列表展示顺序 |
| `data_scope` | 声明这个角色的数据范围类型 |
| `status` | 决定角色是否仍然参与授权 |

这里最值得注意的是 `data_scope`。

这意味着当前项目里的角色，不只是“接口权限的容器”，也是“数据权限的入口”。后面第 5 章做部门、岗位和数据权限时，会继续沿用这套模型，而不是额外再造一套授权体系。

## 为什么用户和角色一定要用关系表

当前项目没有把 `role_id` 直接塞进 `sys_user`，而是使用 `sys_user_role`：

```text
sys_user
  ↕
sys_user_role
  ↕
sys_role
```

原因很直接：

- 一个后台用户可能同时拥有多个角色
- 一个角色也会被多个用户复用
- 多角色并集，是后续菜单权限和数据权限成立的前提

比如某个用户既是“部门管理员”，又是“内容运营”，那他的最终能力不该由单一角色决定，而应该是：

- 接口权限取并集
- 菜单权限取并集
- 数据范围按规则合并

这也是 `middleware.LoadActor` 在请求期要一次性加载 `RoleCodes` 和 `Grants` 的原因。

## 当前角色到底承接了哪三类权限

### 1. 接口权限：角色编码 → `casbin_rule`

接口权限是这样挂到角色上的：

```text
角色 code
  ↓
casbin_rule.v0
  ↓
路径 + 方法
```

也就是说，Casbin 的主体不是用户 ID，而是角色编码。这样用户与角色关系调整时，只要重新绑定角色，就能复用既有接口策略。

### 2. 菜单权限：角色 ID → `sys_role_menu`

菜单权限使用角色 ID 与菜单 ID 建立绑定：

```text
sys_role.id
  ↓
sys_role_menu
  ↓
sys_menu.id
```

这条链路最终服务于 `/api/v1/auth/menus`，由认证模块返回“当前登录用户可见的完整菜单树”。

### 3. 数据权限：角色范围 → `sys_role_data_scope`

当角色的 `data_scope = custom_dept` 时，还会额外挂接自定义部门列表：

```text
sys_role.data_scope = custom_dept
  ↓
sys_role_data_scope
  ↓
可访问的部门 ID 列表
```

后续请求里，`LoadActor` 会把这些规则组装成 `datascope.Actor`，再由 `datascope.Merge(...)` 合成当前登录人的最终数据范围摘要。

## 当前角色管理模块暴露了哪些真实能力

现在角色相关能力已经收进 `server/internal/module/iam/role/`，并且不再只是“查个角色列表”：

| 接口 | 用途 |
| --- | --- |
| `GET /api/v1/system/roles` | 查询角色列表 |
| `POST /api/v1/system/roles` | 创建角色 |
| `POST /api/v1/system/roles/:id/update` | 更新角色基础信息和数据范围 |
| `POST /api/v1/system/roles/:id/status` | 启停角色 |
| `POST /api/v1/system/roles/:id/permissions` | 更新角色接口权限 |
| `POST /api/v1/system/roles/:id/menus` | 更新角色菜单权限 |

这几个接口一起说明了一件事：

> 当前项目已经把“角色”当成统一授权入口，而不是把接口权限、菜单权限、数据权限拆成三个彼此孤立的模块。

## 多角色在请求期是怎么合并的

当请求进入受保护接口后，`middleware.LoadActor` 会做两件事：

1. 查出当前用户拥有的启用角色编码。
2. 查出这些角色对应的数据范围规则，并压成 `Actor`。

最终 `/api/v1/auth/me` 会返回这份结果的一部分，例如：

```json
{
  "user_id": 1,
  "username": "admin",
  "department_id": 1,
  "role_codes": ["super_admin"],
  "is_super_admin": true,
  "data_scope": {
    "allow_all": true,
    "require_self": false,
    "include_department": false,
    "include_dept_tree": false,
    "custom_department_ids": []
  }
}
```

这说明角色模型真正服务的，不只是“登录以后做一次权限判断”，而是：

- 请求上下文
- 接口访问控制
- 菜单返回
- 数据范围过滤

## 本节最关键的结论

这一节真正要建立的判断是：

> 在当前最终版结构里，角色不是一个简单字典表，而是整个授权系统的统一承接点。

接口权限、菜单权限和数据权限虽然最终落在不同链路上，但起点都是角色。

下一节继续把“角色如何承接接口权限”单独讲透：[接口级权限控制](./casbin-permission)。
