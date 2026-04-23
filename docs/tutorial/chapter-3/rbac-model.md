---
title: 角色与权限模型
description: "设计角色表和用户角色关系表，为后续接口权限与菜单权限打基础。"
---

# 角色与权限模型

前面已经能识别“当前登录用户是谁”。这一节开始补权限模型的基础：用户可以绑定角色，角色再承接后续的接口权限和菜单权限。

::: tip 🎯 本节目标
完成后，数据库中会新增 `sys_role` 和 `sys_user_role` 两张表；启动服务时会初始化 `super_admin` 角色，并把默认管理员绑定到这个角色。
:::

## 本节会改什么

本节会新增或修改下面这些文件：

```text
server/
├─ internal/
│  ├─ bootstrap/
│  │  └─ bootstrap.go
│  └─ model/
│     ├─ role.go
│     └─ user_role.go
```

| 位置 | 用途 |
| --- | --- |
| `internal/model/role.go` | 定义角色表结构 |
| `internal/model/user_role.go` | 定义用户与角色的绑定关系 |
| `internal/bootstrap/bootstrap.go` | 初始化超级管理员角色，并绑定默认管理员 |

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
| Casbin 权限控制 | 让角色拥有接口访问权限 |
| 菜单权限设计 | 让角色拥有菜单和按钮权限 |

::: warning ⚠️ 本项目不使用数据库外键约束
`sys_user_role.user_id` 和 `sys_user_role.role_id` 只表达关联关系，并建立普通索引和联合唯一索引，不创建数据库级外键。

用户是否存在、角色是否存在、删除前能不能删除，都由后续 service 层逻辑维护。
:::

## 先创建数据表

本节新增 `sys_role` 和 `sys_user_role`，分别用于保存后台角色和用户角色绑定关系。

::: tip 建表 SQL
字段说明、索引设计、关系表约定和 PostgreSQL / MySQL 建表语句统一放在参考手册：

- [数据库建表语句 - `sys_role`](../../reference/database-ddl#sys-role)
- [数据库建表语句 - `sys_user_role`](../../reference/database-ddl#sys-user-role)
:::

## 🛠️ 创建角色模型

创建 `server/internal/model/role.go`。这是新增文件，直接完整写入即可。

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

## 🛠️ 修改启动初始化

修改 `server/internal/bootstrap/bootstrap.go`。这一节建议直接替换完整文件，避免漏掉函数返回值调整。

```go
package bootstrap

import (
	"errors"
	"fmt"

	"ez-admin-gin/server/internal/model"

	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

const (
	defaultAdminUsername = "admin"
	defaultAdminPassword = "EzAdmin@123456"
	defaultAdminRoleCode = "super_admin"
	defaultAdminRoleName = "超级管理员"
)

// Run 执行服务启动时必须完成的初始化动作。
func Run(db *gorm.DB, log *zap.Logger) error {
	admin, err := seedDefaultAdmin(db, log)
	if err != nil {
		return fmt.Errorf("seed default admin: %w", err)
	}

	role, err := seedSuperAdminRole(db, log)
	if err != nil {
		return fmt.Errorf("seed super admin role: %w", err)
	}

	if err := seedAdminRole(db, admin.ID, role.ID, log); err != nil {
		return fmt.Errorf("seed admin role: %w", err)
	}

	return nil
}

// seedDefaultAdmin 创建本地起步用的默认管理员。
func seedDefaultAdmin(db *gorm.DB, log *zap.Logger) (*model.User, error) {
	var user model.User
	// Unscoped 会把已逻辑删除记录也查出来，避免重复创建同名默认账号。
	err := db.Unscoped().Where("username = ?", defaultAdminUsername).First(&user).Error
	if err == nil {
		return &user, nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(defaultAdminPassword), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("hash default admin password: %w", err)
	}

	user = model.User{
		Username:     defaultAdminUsername,
		PasswordHash: string(passwordHash),
		Nickname:     "系统管理员",
		Status:       model.UserStatusEnabled,
	}

	if err := db.Create(&user).Error; err != nil {
		return nil, err
	}

	log.Info("default admin user created", zap.String("username", defaultAdminUsername))

	return &user, nil
}

// seedSuperAdminRole 创建超级管理员角色。
func seedSuperAdminRole(db *gorm.DB, log *zap.Logger) (*model.Role, error) {
	var role model.Role
	// 角色编码唯一，查询历史记录可以避免逻辑删除后重复创建同名编码。
	err := db.Unscoped().Where("code = ?", defaultAdminRoleCode).First(&role).Error
	if err == nil {
		return &role, nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	role = model.Role{
		Code:   defaultAdminRoleCode,
		Name:   defaultAdminRoleName,
		Sort:   0,
		Status: model.RoleStatusEnabled,
		Remark: "系统内置角色",
	}

	if err := db.Create(&role).Error; err != nil {
		return nil, err
	}

	log.Info("default admin role created", zap.String("role_code", defaultAdminRoleCode))

	return &role, nil
}

// seedAdminRole 绑定默认管理员和超级管理员角色。
func seedAdminRole(db *gorm.DB, userID uint, roleID uint, log *zap.Logger) error {
	var userRole model.UserRole
	err := db.Where("user_id = ? AND role_id = ?", userID, roleID).First(&userRole).Error
	if err == nil {
		return nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	userRole = model.UserRole{
		UserID: userID,
		RoleID: roleID,
	}

	if err := db.Create(&userRole).Error; err != nil {
		return err
	}

	log.Info(
		"default admin role bound",
		zap.Uint("user_id", userID),
		zap.Uint("role_id", roleID),
	)

	return nil
}
```

::: warning ⚠️ 启动初始化不是权限管理接口
这里初始化 `super_admin` 角色，只是为了让本地起步有一条完整的用户角色关系。后续真正的角色创建、授权、解绑，要放在系统管理接口中处理。
:::

::: details 为什么默认角色叫 `super_admin`
`admin` 通常更像用户名，而 `super_admin` 更像角色编码。这样可以避免“账号名”和“角色名”混在一起。
:::

## ✅ 整理依赖并启动

本节没有新增第三方依赖，但修改了模型和初始化逻辑，仍然可以整理一次：

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
INFO	default admin role created	{"role_code": "super_admin"}
INFO	default admin role bound	{"user_id": 1, "role_id": 1}
INFO	server started	{"addr": ":8080", "env": "dev"}
```

如果角色和绑定关系已经存在，后续启动不会重复创建。

## ✅ 验证角色和绑定关系

打开另一个终端，在项目根目录执行：

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

继续查看默认管理员和角色的绑定关系：

```bash
# 查看默认管理员绑定了哪些角色
docker compose -f deploy/compose.local.yml exec postgres psql -U ez_admin -d ez_admin -c "select ur.id, u.username, r.code from sys_user_role ur join sys_user u on u.id = ur.user_id join sys_role r on r.id = ur.role_id;"
```

应该看到类似结果：

```text
 id | username |    code
----+----------+-------------
  1 | admin    | super_admin
```

::: details 如果提示 `relation "sys_role" does not exist`
说明角色表还没有创建。先回到 [`sys_role` 建表语句](../../reference/database-ddl#sys-role)，执行对应数据库版本的 SQL，然后重新启动服务。
:::

::: details 如果绑定关系没有出现
先确认 `bootstrap.Run` 中已经调用了 `seedAdminRole`，再确认 `sys_user` 和 `sys_role` 中分别存在默认管理员和 `super_admin` 角色。
:::

## 常见问题

::: details 为什么不直接在 `sys_user` 表里放 `role_id`
一个后台用户后续可能同时拥有多个角色，例如既是“内容管理员”，又是“运营管理员”。如果直接在用户表里放一个 `role_id`，很快就会不够用。

单独使用 `sys_user_role` 关系表，可以自然支持多角色。
:::

::: details 角色和权限现在是什么关系
本节先让用户拥有角色。下一节会用 Casbin 表达“角色可以访问哪些接口”；再下一节会用菜单模型表达“角色能看到哪些菜单和按钮”。
:::

下一节会把角色和接口权限连接起来：[Casbin 权限控制](./casbin-permission)。
