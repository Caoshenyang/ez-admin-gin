---
title: 组织体系
description: 部门树、岗位管理、用户归属关系
---

# 组织体系

EZ Admin 的组织体系由部门、岗位、用户三个实体构成，支撑数据权限的部门过滤能力。

## 组织模型

```
部门（sys_department）     岗位（sys_post）
  ┌──────────┐             ┌──────────┐
  │ 树形结构  │             │ 扁平列表  │
  │ ancestors │             │ sort     │
  └────┬─────┘             └────┬─────┘
       │                        │
       │    用户（sys_user）     │
       │   ┌──────────────┐    │
       └──►│ department_id │    │
           │              │◄───┘
           │ user_post (M:N)
           └──────────────┘
```

## 部门（Department）

部门采用树形结构，支持无限层级。

### 数据模型

```
sys_department
├── id              主键
├── parent_id       父部门 ID（0 为顶级）
├── ancestors       祖先路径（逗号分隔，如 "0,1,2"）
├── name            部门名称
├── code            部门编码
├── leader_user_id  负责人用户 ID
├── sort            排序值
├── status          状态（启用/禁用）
├── created_at / updated_at / deleted_at
```

### 祖先路径（ancestors）

`ancestors` 字段存储从根到当前节点的完整路径，用于高效子树查询：

```
总公司        ancestors: "0"
技术部        ancestors: "0,1"
  前端组      ancestors: "0,1,2"
  后端组      ancestors: "0,1,3"
市场部        ancestors: "0,4"
```

查询"技术部及下级所有部门"：`WHERE ancestors LIKE '0,1%'`。

### 部门管理 API

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/v1/iam/departments/tree` | 获取部门树 |
| POST | `/api/v1/iam/departments` | 创建部门 |
| PUT | `/api/v1/iam/departments/:id` | 更新部门 |
| DELETE | `/api/v1/iam/departments/:id` | 删除部门 |

## 岗位（Post）

岗位是扁平结构，用于标识用户的职能角色（如"经理"、"开发"、"测试"）。

### 数据模型

```
sys_post
├── id          主键
├── name        岗位名称
├── code        岗位编码
├── sort        排序值
├── status      状态
├── remark      备注
├── created_at / updated_at / deleted_at
```

### 用户-岗位关联

通过 `sys_user_post` 关联表实现多对多关系：一个用户可以有多个岗位，一个岗位可以属于多个用户。

## 用户归属关系

```
sys_user
├── department_id   → sys_department.id   （多对一：用户属于一个部门）
├── user_role       → sys_user_role        （多对多：用户拥有多个角色）
└── user_post       → sys_user_post        （多对多：用户拥有多个岗位）
```

## 与数据权限的关系

组织体系是数据权限的基础：

- **部门**直接参与 `dept`、`dept_and_children`、`custom_dept` 三种数据范围的过滤
- 用户的 `department_id` 决定了"本部门"的边界
- 部门的 `ancestors` 决定了"下级部门"的范围
- 角色的 `custom_dept` 通过 `sys_role_data_scope` 关联具体部门
