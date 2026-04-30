---
title: 数据权限模型
description: "说明 EZ Admin Gin v2 首版数据权限的五档模型、角色字段、关系表和合并规则。"
---

# 数据权限模型

::: tip 🎯 这页解决什么
快速说明 `EZ Admin Gin` 在 `v2` 第一阶段引入的数据权限模型长什么样，以及它为什么要依赖组织体系一起设计。
:::

## 首版固定支持 5 档

`EZ Admin Gin` 的首版数据权限直接采用企业后台里最常见的一组范围：

| 值 | 含义 |
| --- | --- |
| `all` | 全部数据 |
| `dept` | 本部门数据 |
| `dept_and_children` | 本部门及子部门数据 |
| `self` | 仅本人数据 |
| `custom_dept` | 自定义授权部门数据 |

这组枚举已经在后端 `platform/datascope` 中固化，后续模块接入时不需要再自己发明一套命名。

## 为什么要和组织体系一起设计

数据权限不是单独加一个字段就能完成的能力。  
只要支持 `dept` 和 `dept_and_children`，系统就必须同时具备：

- 部门表
- 用户归属部门
- 部门树查询能力

只要支持 `custom_dept`，系统还必须同时具备：

- 角色到部门的授权关系表

所以 `v2` 第一阶段先补下面这些结构：

- `sys_department`
- `sys_post`
- `sys_user.department_id`
- `sys_user_post`
- `sys_role.data_scope`
- `sys_role_data_scope`

## 角色如何表达数据范围

角色表新增字段：

| 字段 | 说明 |
| --- | --- |
| `data_scope` | 当前角色的数据权限范围 |

对于 `custom_dept`，角色还会通过关系表绑定一组可见部门：

| 表名 | 说明 |
| --- | --- |
| `sys_role_data_scope` | 角色和自定义部门范围的绑定关系 |

## 多角色怎么合并

同一个用户可能拥有多个角色。  
首版规则固定为：**按并集合并**。

也就是说：

- 只要任一角色是 `all`，结果就是全部数据
- 一个角色给本部门，另一个角色给自定义部门，最终结果就是两者并集
- `super_admin` 永远绕过数据权限限制

这套规则的好处是简单、稳定，也更贴近企业后台的常见预期。

## 部门树为什么用 `parent_id + ancestors`

部门表首版采用：

- `parent_id`
- `ancestors`

例如一条部门记录的 `ancestors` 可能是：

```text
0,1,3
```

这表示它的祖先链路是根节点 → 1 → 3。

这样设计的原因是：

- PostgreSQL / MySQL 都容易实现
- 查询“本部门及子部门”更直接
- 对教程和排查也更友好

## 当前阶段的数据权限接入边界

`v2` 第一阶段先把模型和基础设施打底，不会一口气让所有系统模块都启用数据过滤。  
优先接入的资源会是：

- 用户管理
- 部门管理
- 岗位管理
- 后续真实业务示例模块

而像下面这些系统级资源，通常不会在第一阶段就做数据权限限制：

- 角色管理
- 菜单管理
- 权限策略管理
- 系统配置
- 登录日志
- 操作日志

## 相关代码与迁移

第一阶段已经落地的关键位置：

- 枚举定义：`server/internal/platform/datascope/`
- 模型定义：`server/internal/model/department.go`、`post.go`、`role_data_scope.go`、`user_post.go`
- 迁移文件：`server/migrations/postgres/000003_enterprise_foundation.up.sql`
- MySQL 迁移：`server/migrations/mysql/000003_enterprise_foundation.up.sql`

## 下一步通常会接什么

有了这套模型后，后续的数据权限真正落地会继续补：

- 当前登录人 `Actor` 上下文
- 角色数据范围加载
- `gorm.Scopes(...)` 查询过滤
- 模块内 `datascope.go` 规则声明

想理解这一步为什么要先做结构升级，可以继续看 [企业级架构升级](/guide/enterprise-architecture)。
