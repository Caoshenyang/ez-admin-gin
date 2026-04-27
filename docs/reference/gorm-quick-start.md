---
title: GORM 快速入门
description: "面向当前后台底座开发的 GORM 扫盲参考，覆盖模型、连接、查询、软删除、事务和常见约定。"
---

# GORM 快速入门

GORM 是 Go 生态里常用的 ORM。它把数据库表映射成 Go 结构体，让你用链式 API 完成增删改查，同时保留在必要时写 SQL 的能力。

::: tip 这页怎么读
第一次读先抓住一个主线：`model` 定义表结构，`*gorm.DB` 执行查询，`Error` 判断结果。后续写业务模块时，回到这里查常用写法即可。
:::

## 当前项目怎么使用 GORM

当前后端使用：

| 依赖 | 用途 |
| --- | --- |
| `gorm.io/gorm` | GORM 核心库 |
| `gorm.io/driver/postgres` | PostgreSQL 驱动 |
| `github.com/casbin/gorm-adapter/v3` | 让 Casbin 从数据库读取权限策略 |

项目中的主要落点：

| 文件 | 职责 |
| --- | --- |
| `server/internal/database/database.go` | 创建数据库连接、设置连接池、健康检查 |
| `server/internal/model/*.go` | 定义数据库模型和表名 |
| `server/migrations/{postgres,mysql}/000002_seed_data.up.sql` | 初始化默认账号、角色和权限策略 |
| `server/internal/middleware/permission.go` | 查询用户角色并交给 Casbin 判断权限 |

::: info 版本以项目为准
本文示例服务于当前仓库。具体版本可以看 `server/go.mod`，不要把别的项目里的 GORM 写法直接照搬进来。
:::

## 先认识 `*gorm.DB`

在业务代码里最常见的对象是 `*gorm.DB`。它不是一条数据库连接，而是 GORM 的数据库操作入口，底层会使用 `database/sql` 的连接池。

当前项目启动时会创建一次数据库对象：

```go
db, err := gorm.Open(postgres.Open(dsn(cfg)), &gorm.Config{
	// Warn 级别会记录慢查询和潜在问题，同时避免开发日志过多。
	Logger: gormLogger.Default.LogMode(gormLogger.Warn),
})
```

拿到 `db` 后，常见使用方式是：

```go
var user model.User
err := db.Where("username = ?", username).First(&user).Error
```

这段代码可以拆成三步理解：

| 片段 | 含义 |
| --- | --- |
| `Where("username = ?", username)` | 拼接查询条件，参数会被绑定，避免手动拼接字符串 |
| `First(&user)` | 查询一条记录并写入 `user` |
| `.Error` | 取出本次数据库操作的错误 |

::: warning 不要手动拼接用户输入
查询条件里的用户输入要通过 `?` 传参，不要把字符串直接拼进 SQL。后台系统里登录名、关键词、路径参数都可能来自请求，默认按不可信输入处理。
:::

## 模型就是表结构的 Go 表达

GORM 通过结构体描述表字段、索引、默认值和 JSON 输出。

当前用户模型的核心写法如下：

```go
type User struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	Username     string         `gorm:"size:64;not null;uniqueIndex" json:"username"`
	PasswordHash string         `gorm:"size:255;not null" json:"-"`
	Nickname     string         `gorm:"size:64;not null;default:''" json:"nickname"`
	Status       UserStatus     `gorm:"type:smallint;not null;default:1" json:"status"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

func (User) TableName() string {
	return "sys_user"
}
```

需要重点记住几类标签：

| 写法 | 作用 |
| --- | --- |
| `primaryKey` | 主键 |
| `size:64` | 字符串长度 |
| `not null` | 数据库非空 |
| `default:''` | 数据库默认值 |
| `uniqueIndex` | 唯一索引 |
| `index` | 普通索引 |
| `type:smallint` | 指定数据库字段类型 |

### 为什么要写 `TableName`

GORM 默认会根据结构体名推导表名。当前项目显式写 `TableName()`，是为了固定系统表名，例如：

```go
func (Role) TableName() string {
	return "sys_role"
}
```

这样后续就算调整 GORM 命名策略，也不会影响已经上线的表。

## 创建数据

创建数据使用 `Create`：

```go
role := model.Role{
	Code:   "super_admin",
	Name:   "超级管理员",
	Status: model.RoleStatusEnabled,
	Remark: "系统内置角色",
}

if err := db.Create(&role).Error; err != nil {
	return err
}
```

GORM 会做几件事：

- 根据 `TableName()` 找到表。
- 把结构体字段写入对应列。
- 自动维护 `CreatedAt` 和 `UpdatedAt`。
- 数据库生成主键后，回填到 `role.ID`。

::: warning 时间字段由 GORM 维护
当前项目约定 `created_at` 和 `updated_at` 由应用代码维护。通过 GORM 创建和更新数据时会自动处理；如果直接写 SQL，需要显式写入这两个字段。
:::

## 查询一条数据

登录时查询用户就是典型的一条记录查询：

```go
var user model.User
err := db.Where("username = ?", req.Username).First(&user).Error
if err != nil {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		// 没查到记录，通常转换成业务错误。
		return apperror.Unauthorized("用户名或密码错误")
	}

	return apperror.Internal("登录失败", err)
}
```

常见查询方法：

| 方法 | 用途 |
| --- | --- |
| `First(&user)` | 查询第一条，查不到会返回 `gorm.ErrRecordNotFound` |
| `Find(&users)` | 查询多条，查不到通常不会返回 `ErrRecordNotFound` |
| `Take(&user)` | 查询一条，不附加排序 |
| `Where(...)` | 添加查询条件 |
| `Order(...)` | 排序 |
| `Limit(...)` / `Offset(...)` | 分页 |

::: tip 错误处理习惯
查询单条记录时，优先判断 `errors.Is(err, gorm.ErrRecordNotFound)`，再处理其它数据库错误。这样业务语义更清楚。
:::

## 查询列表和分页

业务列表通常会同时需要总数和当前页数据：

```go
var total int64
var users []model.User

query := db.Model(&model.User{}).
	Where("status = ?", model.UserStatusEnabled)

if err := query.Count(&total).Error; err != nil {
	return err
}

if err := query.
	Order("id DESC").
	Limit(pageSize).
	Offset((page - 1) * pageSize).
	Find(&users).Error; err != nil {
	return err
}
```

这个模式适合后台管理页面：

| 步骤 | 作用 |
| --- | --- |
| `Model(&model.User{})` | 指定查询的表 |
| `Count(&total)` | 查询总数 |
| `Limit` / `Offset` | 查询当前页 |
| `Order` | 固定排序，避免翻页结果漂移 |

## 更新数据

更新单个字段：

```go
err := db.Model(&model.User{}).
	Where("id = ?", userID).
	Update("status", model.UserStatusDisabled).Error
```

更新多个字段：

```go
err := db.Model(&model.User{}).
	Where("id = ?", userID).
	Updates(map[string]any{
		"nickname": nickname,
		"status":   model.UserStatusEnabled,
	}).Error
```

::: warning `Updates(struct)` 会跳过零值
如果用结构体更新，GORM 默认会跳过 `0`、`""`、`false` 这类零值。后台表单里经常需要把字段更新为空字符串或 `0`，这时优先使用 `map[string]any`，语义更直接。
:::

## 删除和软删除

模型里只要包含：

```go
DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
```

GORM 的普通删除就会变成软删除：

```go
err := db.Delete(&model.User{}, userID).Error
```

它不是物理删除，而是更新 `deleted_at`。之后普通查询会自动过滤掉已删除记录。

如果确实需要查历史记录，可以使用 `Unscoped()`：

```go
var user model.User
err := db.Unscoped().
	Where("username = ?", "admin").
	First(&user).Error
```

当前初始化默认管理员时就使用了 `Unscoped()`，目的是避免默认账号被逻辑删除后又重复创建同名账号。

## 关联查询怎么写

当前项目还没有大规模使用 GORM 的自动关联，更多是用显式查询表达关系。例如权限中间件查询当前用户拥有的角色编码：

```go
var roleCodes []string
err := db.
	Table("sys_role AS r").
	Select("r.code").
	Joins("JOIN sys_user_role AS ur ON ur.role_id = r.id").
	Where("ur.user_id = ?", userID).
	Where("r.status = ?", model.RoleStatusEnabled).
	Where("r.deleted_at IS NULL").
	Pluck("r.code", &roleCodes).Error
```

这类写法的好处是：

- SQL 结构清晰。
- 不依赖隐式外键。
- 适合后台权限、菜单、日志这类需要明确控制查询条件的场景。

::: info 什么时候用 `Preload`
如果只是加载普通的一对多、一对一关系，`Preload` 很方便；但权限、菜单、复杂筛选和跨表统计，更推荐写清楚 `Joins`、`Select` 和 `Where`。
:::

## 事务

当一个业务动作必须同时写多张表时，使用事务。例如创建角色并初始化权限：

```go
err := db.Transaction(func(tx *gorm.DB) error {
	if err := tx.Create(&role).Error; err != nil {
		return err
	}

	if err := tx.Create(&rule).Error; err != nil {
		return err
	}

	return nil
})
```

事务函数里：

| 返回值 | 结果 |
| --- | --- |
| `return nil` | 提交事务 |
| `return err` | 回滚事务 |

::: warning 事务里要一直使用 `tx`
进入 `db.Transaction` 后，事务内部的读写都要用 `tx`，不要又写回外层 `db`。否则有些操作会跑到事务外面。
:::

## 当前项目的 GORM 约定

| 场景 | 推荐做法 |
| --- | --- |
| 表名 | 模型上显式实现 `TableName()` |
| 主键 | 使用 `uint`，数据库自增生成 |
| 时间字段 | `CreatedAt` / `UpdatedAt` 由 GORM 自动维护 |
| 软删除 | 使用 `gorm.DeletedAt` |
| 强身份唯一字段 | 普通唯一索引，删除后不复用 |
| 复杂查询 | 优先写清 `Table` / `Select` / `Joins` / `Where` |
| 单条查询错误 | 使用 `errors.Is(err, gorm.ErrRecordNotFound)` 判断 |
| 批量写多表 | 使用 `db.Transaction` |
| 建表 | 以参考手册 SQL 为准，不依赖 GORM 自动迁移 |

## 常见坑速查

| 问题 | 原因 | 建议 |
| --- | --- | --- |
| 查不到已删除数据 | `gorm.DeletedAt` 默认过滤软删除记录 | 需要查历史时使用 `Unscoped()` |
| 更新空字符串不生效 | `Updates(struct)` 默认跳过零值 | 使用 `Updates(map[string]any{...})` |
| 同名用户删除后仍不能创建 | 普通唯一索引仍约束历史记录 | 这是当前项目的默认设计 |
| 事务里部分数据提前写入 | 事务内部混用了外层 `db` | 事务函数内统一使用 `tx` |
| 查询拼接字符串有风险 | 用户输入直接进入 SQL | 使用 `?` 参数绑定 |

## 继续查资料

- [GORM 官方文档](https://gorm.io/)
- [GORM Go 包文档](https://pkg.go.dev/gorm.io/gorm)
- [PostgreSQL Driver 文档](https://pkg.go.dev/gorm.io/driver/postgres)
- [数据库建表语句](./database-ddl)
- [逻辑删除与唯一索引冲突](./logical-delete-and-unique-index)

## 小结

在当前后台底座里，GORM 不需要先学到很深。先熟悉这几件事就足够支撑日常开发：模型标签、`Where` 查询、`Create` 创建、`Updates` 更新、`gorm.DeletedAt` 软删除、`Transaction` 事务，以及永远从 `.Error` 判断数据库操作结果。
