---
title: RBAC 角色权限模型
description: "设计角色表和用户角色关系表，为后续接口权限与菜单权限打基础。"
---

# RBAC 角色权限模型

前面已经能识别“当前登录用户是谁”。这一节开始补权限模型的基础：用户可以绑定角色，角色再承接后续的接口权限和菜单权限。

::: tip 🎯 本节目标
完成后，数据库中会新增 `sys_role` 和 `sys_user_role` 两张表；启动服务时会初始化 `super_admin` 角色，并把默认管理员绑定到这个角色。
:::

## 本节会改什么

本节会新增或修改下面这些文件：

```text
server/
└─ internal/
   └─ model/
      ├─ role.go
      └─ user_role.go
```

| 位置 | 用途 |
| --- | --- |
| `internal/model/role.go` | 定义角色表结构 |
| `internal/model/user_role.go` | 定义用户与角色的绑定关系 |

## 先看关系

本节先落地下面这条关系：

```text
用户 sys_user
  ↓
用户角色关系 sys_user_role
  ↓
角色 sys_role
```

一个用户可以绑定多个角色，一个角色也可以绑定多个用户。后续：

| 后续小节 | 继续完成什么 |
| --- | --- |
| 接口级权限控制 | 让角色拥有接口访问权限 |
| 角色菜单权限 | 让角色拥有菜单和按钮权限 |

::: warning ⚠️ 本项目不使用数据库外键约束
`sys_user_role.user_id` 和 `sys_user_role.role_id` 只表达关联关系，并建立普通索引和联合唯一索引，不创建数据库级外键。

用户是否存在、角色是否存在、删除前能不能删除，都由后续 service 层逻辑维护。
:::

## 先创建数据表

本节新增 `sys_role` 和 `sys_user_role`，分别用于保存后台角色和用户角色绑定关系。

`sys_role` 表保存后台角色编码、名称和状态；`sys_user_role` 表保存用户与角色的绑定关系。字段和索引详情见 [数据库建表语句 - `sys_role`](/reference/database-ddl#sys-role) 和 [数据库建表语句 - `sys_user_role`](/reference/database-ddl#sys-user-role)。

## 🛠️ 创建角色模型

创建 `server/internal/model/role.go`。这是新增文件，直接完整写入即可。

::: details `server/internal/model/role.go` — 角色模型

```go
package model

import (
	"time"

	"gorm.io/gorm"
)

// RoleStatus 表示角色状态。
type RoleStatus int

const (
	// RoleStatusEnabled 表示角色可以正常使用。
	RoleStatusEnabled RoleStatus = 1
	// RoleStatusDisabled 表示角色已被禁用。
	RoleStatusDisabled RoleStatus = 2
)

// Role 是后台角色表模型。
type Role struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Code      string         `gorm:"size:64;not null;uniqueIndex" json:"code"`
	Name      string         `gorm:"size:64;not null" json:"name"`
	Sort      int            `gorm:"not null;default:0" json:"sort"`
	Status    RoleStatus     `gorm:"type:smallint;not null;default:1" json:"status"`
	Remark    string         `gorm:"size:255;not null;default:''" json:"remark"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 固定角色表名，避免后续调整命名策略时影响已有表。
func (Role) TableName() string {
	return "sys_role"
}
```

:::

## 🛠️ 创建用户角色关系模型

创建 `server/internal/model/user_role.go`。这是新增文件，直接完整写入即可。

```go
package model

import "time"

// UserRole 是用户与角色的绑定关系。
type UserRole struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"not null;uniqueIndex:uk_sys_user_role_user_role;index:idx_sys_user_role_user_id" json:"user_id"`
	RoleID    uint      `gorm:"not null;uniqueIndex:uk_sys_user_role_user_role;index:idx_sys_user_role_role_id" json:"role_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TableName 固定用户角色关系表名。
func (UserRole) TableName() string {
	return "sys_user_role"
}
```

::: tip 📌 超级管理员角色初始化
超级管理员角色（`super_admin`）通过数据库迁移文件自动创建，不需要在代码中手动初始化。当服务启动时，会执行 `server/migrations/{postgres,mysql}/000002_seed_data.up.sql` 迁移文件，创建超级管理员角色、系统菜单和权限。

管理员账号需要通过 `/api/v1/setup/init` 接口创建，创建时会自动绑定到 `super_admin` 角色。
:::

::: details 为什么默认角色叫 `super_admin`
`admin` 通常更像用户名，而 `super_admin` 更像角色编码。这样可以避免“账号名”和“角色名”混在一起。
:::

## ✅ 整理依赖并启动

本节没有新增第三方依赖，整理一次依赖：

```bash
# 在 server/ 目录执行
go mod tidy
```

确认数据库和 Redis 正在运行：

```bash
# 在项目根目录执行，确认本地依赖服务处于运行状态
docker compose -f deploy/compose.local.yml ps
```

回到 `server/` 目录启动服务：

```bash
# 在 server/ 目录启动服务
go run .
```

第一次启动后，控制台应该能看到类似日志：

```text
INFO	database migrations applied
INFO	server started	{"addr": ":8080", "env": "dev"}
```

## ✅ 创建管理员账号并验证角色和绑定关系

服务启动后，先通过初始化接口创建管理员账号：

```bash
# 创建管理员账号
curl -X POST http://localhost:8080/api/v1/setup/init \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"YourPassword123","nickname":"管理员"}'
```

然后验证角色和绑定关系：

```bash
# 查看默认角色
docker compose -f deploy/compose.local.yml exec postgres psql -U ez_admin -d ez_admin -c "select id, code, name, status, deleted_at from sys_role;"
```

应该看到类似结果：

```text
 id |    code     |    name    | status | deleted_at
----+-------------+------------+--------+------------
  1 | super_admin | 超级管理员 |      1 |
```

继续查看管理员和角色的绑定关系：

```bash
# 查看管理员绑定了哪些角色
docker compose -f deploy/compose.local.yml exec postgres psql -U ez_admin -d ez_admin -c "select ur.id, u.username, r.code from sys_user_role ur join sys_user u on u.id = ur.user_id join sys_role r on r.id = ur.role_id;"
```

应该看到类似结果：

```text
 id | username |    code
----+----------+-------------  1 | admin    | super_admin
```

::: details 如果提示 `relation "sys_role" does not exist`
说明角色表还没有创建。服务启动时会自动执行数据库迁移，创建表结构。如果迁移失败，查看服务启动日志获取详细信息。
:::

::: details 如果绑定关系没有出现
确认管理员账号已经通过 `/api/v1/setup/init` 接口创建成功，该接口会自动绑定 `super_admin` 角色。
:::

## 常见问题

::: details 为什么不直接在 `sys_user` 表里放 `role_id`
一个后台用户后续可能同时拥有多个角色，例如既是“内容管理员”，又是“运营管理员”。如果直接在用户表里放一个 `role_id`，很快就会不够用。

单独使用 `sys_user_role` 关系表，可以自然支持多角色。
:::

::: details 角色和权限现在是什么关系
本节先让用户拥有角色。下一节会用 Casbin 表达“角色可以访问哪些接口”；再下一节会用菜单模型表达“角色能看到哪些菜单和按钮”。
:::

下一节会把角色和接口权限连接起来：[接口级权限控制](./casbin-permission)。
