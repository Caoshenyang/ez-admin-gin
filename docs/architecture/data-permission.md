---
title: 数据权限
description: 五级数据作用域、Actor 上下文、Repository 层注入机制
---

# 数据权限

数据权限控制用户能看到哪些数据，在 Repository 层通过 GORM scope 注入，对 Service 层透明。

## 设计思路

数据权限的核心问题是：同一个列表接口，不同用户看到的记录范围不同。

EZ Admin 的方案是：

1. **角色配置数据范围**：每个角色设置一种数据作用域类型
2. **Repository 层注入**：查询时根据当前用户的角色数据范围，自动追加过滤条件
3. **多角色取并集**：用户拥有多个角色时，数据范围取并集（最宽范围生效）

## 五级数据作用域

| 作用域 | 值 | 含义 | 过滤条件 |
|--------|---|------|---------|
| all | 1 | 所有数据 | 不追加过滤 |
| dept | 2 | 本部门 | `department_id = 用户部门ID` |
| dept_and_children | 3 | 本部门及下级 | `department_id IN (部门子树)` |
| self | 4 | 仅本人 | `creator_id = 用户ID` |
| custom_dept | 5 | 自定义部门 | `department_id IN (指定部门列表)` |

### 作用域选择指南

| 场景 | 推荐作用域 |
|------|-----------|
| 超级管理员 | all |
| 部门经理 | dept_and_children |
| 普通员工 | self |
| 跨部门协作 | custom_dept |

## 数据权限注入机制

### Actor 上下文

每次认证请求通过 `LoadActor` 中间件构建 Actor：

```go
type Actor struct {
    UserID         uint
    Username       string
    DepartmentID   uint
    Roles          []RoleInfo  // 角色列表（含 data_scope）
    MenuIDs        []uint
    ButtonCodes    []string
}
```

### Repository 层调用

```go
// 在 Repository 查询中使用
func (r *Repository) List(db *gorm.DB, actor *actorx.Actor, query ListQuery) ([]Item, int64, error) {
    db = datascope.ApplyScopes(db, actor)
    // ... 正常查询逻辑
}
```

`datascope.ApplyScopes` 根据 Actor 的角色数据范围，自动追加 GORM scope：

- **all** → 不追加条件
- **dept** → `WHERE department_id = actor.DepartmentID`
- **dept_and_children** → 通过 `sys_department.ancestors` 字段查子树，`WHERE department_id IN (...)`
- **self** → `WHERE creator_id = actor.UserID`
- **custom_dept** → 查询 `sys_role_data_scope` 表获取部门列表，`WHERE department_id IN (...)`

## 部门树与祖先路径

部门表通过 `ancestors` 字段实现高效的子树查询：

```
sys_department
├── id: 1, name: "总公司", ancestors: "0"
├── id: 2, name: "技术部", ancestors: "0,1"
├── id: 3, name: "前端组", ancestors: "0,1,2"
└── id: 4, name: "后端组", ancestors: "0,1,2"
```

查询"技术部及下级"时：`WHERE ancestors LIKE '0,1,2%'`。

## 自定义部门范围

当角色的数据范围为 `custom_dept` 时，通过 `sys_role_data_scope` 关联表指定可见部门：

```
sys_role_data_scope
├── role_id        角色ID
└── department_id  部门ID
```

一个角色可以关联多个部门，查询时取并集。

## 多角色并集

用户拥有多个角色时，数据范围取并集（最宽范围）：

1. 遍历用户所有角色的数据范围
2. 收集所有允许的部门 ID
3. 合并去重
4. 如果任一角色为 `all`，则直接放行全部数据

## 相关数据表

| 表 | 用途 |
|----|------|
| `sys_role.data_scope` | 角色数据作用域类型（1-5） |
| `sys_role_data_scope` | 自定义部门范围关联 |
| `sys_department.ancestors` | 部门祖先路径 |
| `sys_user.department_id` | 用户所属部门 |
