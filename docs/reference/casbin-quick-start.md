---
title: Casbin 快速入门
description: "面向当前后台底座开发的 Casbin 扫盲参考，讲清模型、策略表、权限校验链路和新增接口权限的基本做法。"
---

# Casbin 快速入门

Casbin 是一个权限控制库。它不负责登录，也不负责识别当前用户是谁；它只回答一个问题：**某个主体能不能对某个资源执行某个动作**。

::: tip 这页先看结论
当前项目把角色编码当成主体，把接口路径当成资源，把 HTTP 方法当成动作。也就是：`角色 code + API 路径 + 请求方法` 决定接口能不能访问。
:::

## Casbin 解决什么问题

后台系统里，认证和授权是两件不同的事：

| 概念 | 解决的问题 | 当前项目例子 |
| --- | --- | --- |
| 认证 | 你是谁 | JWT 解析出当前用户 ID |
| 授权 | 你能做什么 | Casbin 判断角色能否访问接口 |

登录成功只说明用户身份可信，不代表他能访问所有接口。进入受保护接口后，权限中间件会：

1. 从上下文拿当前用户 ID。
2. 查询用户拥有的启用角色编码。
3. 把角色编码、接口路径、请求方法交给 Casbin。
4. 只要任意一个角色允许访问，就放行。

## 当前项目用到哪些依赖

| 依赖 | 用途 |
| --- | --- |
| `github.com/casbin/casbin/v3` | Casbin 核心权限判断能力 |
| `github.com/casbin/gorm-adapter/v3` | 从数据库表 `casbin_rule` 读取策略 |
| `gorm.io/gorm` | 为 gorm-adapter 提供数据库连接 |

项目中的主要落点：

| 文件 | 职责 |
| --- | --- |
| `server/configs/rbac_model.conf` | 定义 Casbin 模型 |
| `server/internal/permission/enforcer.go` | 创建 Enforcer 并加载策略 |
| `server/internal/middleware/permission.go` | 在接口请求中执行权限判断 |
| `server/internal/model/casbin_rule.go` | 对应 `casbin_rule` 策略表 |
| `server/migrations/{pgsql,mysql}/000002_seed_data.up.sql` | 初始化默认权限策略 |

## 先理解 `sub / obj / act`

当前模型使用三个维度表达权限：

```text
sub = subject，主体，这里是角色编码
obj = object，资源，这里是接口路径
act = action，动作，这里是 HTTP 方法
```

一条策略可以这样读：

```text
p, super_admin, /api/v1/system/health, GET
```

含义是：

```text
角色 super_admin 可以用 GET 访问 /api/v1/system/health
```

对应到 `casbin_rule` 表：

| 字段 | 值 | 含义 |
| --- | --- | --- |
| `ptype` | `p` | 这是一条权限策略 |
| `v0` | `super_admin` | 主体，角色编码 |
| `v1` | `/api/v1/system/health` | 资源，接口路径 |
| `v2` | `GET` | 动作，请求方法 |

::: info 为什么主体用角色编码
用户和角色的绑定会变，但角色编码更稳定。策略绑定角色编码后，一个用户只要拥有这个角色，就自然拥有对应接口权限。
:::

## Casbin 模型文件

当前模型文件是 `server/configs/rbac_model.conf`：

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

逐段理解：

| 配置 | 含义 |
| --- | --- |
| `request_definition` | 调用 `Enforce` 时传入哪些参数 |
| `policy_definition` | 策略表里的每条策略有哪些字段 |
| `policy_effect` | 只要有一条允许策略命中，就允许访问 |
| `matchers` | 请求和策略如何匹配 |

重点是最后一行：

```ini
m = r.sub == p.sub && keyMatch2(r.obj, p.obj) && (r.act == p.act || p.act == "*")
```

它表达三条规则：

| 条件 | 说明 |
| --- | --- |
| `r.sub == p.sub` | 当前角色必须等于策略里的角色 |
| `keyMatch2(r.obj, p.obj)` | 路径支持参数匹配，例如 `/api/v1/users/:id` |
| `r.act == p.act || p.act == "*"` | 方法相同才允许，或者策略方法写 `*` 表示全部方法 |

::: warning 路径要和 Gin 路由模式对齐
权限中间件优先使用 `c.FullPath()`，它拿到的是 Gin 注册路由时的模式路径。比如路由是 `/api/v1/users/:id`，策略也应该写这个模式，而不是某一次请求里的 `/api/v1/users/1`。
:::

## 策略表 `casbin_rule`

Casbin GORM 适配器默认使用 `casbin_rule` 表。当前项目保留了这个表名，并用 SQL 手动建表。

模型结构如下：

```go
type CasbinRule struct {
	ID    uint   `gorm:"primaryKey" json:"id"`
	Ptype string `gorm:"size:100;not null;default:''" json:"ptype"`
	V0    string `gorm:"size:100;not null;default:''" json:"v0"`
	V1    string `gorm:"size:100;not null;default:''" json:"v1"`
	V2    string `gorm:"size:100;not null;default:''" json:"v2"`
	V3    string `gorm:"size:100;not null;default:''" json:"v3"`
	V4    string `gorm:"size:100;not null;default:''" json:"v4"`
	V5    string `gorm:"size:100;not null;default:''" json:"v5"`
}

func (CasbinRule) TableName() string {
	return "casbin_rule"
}
```

当前模型只使用 `ptype`、`v0`、`v1`、`v2`。`v3` 到 `v5` 是 Casbin 适配器预留字段，保留即可。

::: info 建表方式
当前项目统一使用参考手册里的 SQL 建表，不依赖 gorm-adapter 自动迁移。这样表结构、索引和注释都更可控。
:::

## Enforcer 是什么

`Enforcer` 是 Casbin 的权限判断器。当前项目做了一层轻量封装：

```go
func NewEnforcer(db *gorm.DB, modelPath string) (*Enforcer, error) {
	// 项目统一使用 SQL 建表，不让 gorm-adapter 自动迁移表结构。
	gormadapter.TurnOffAutoMigrate(db)

	adapter, err := gormadapter.NewAdapterByDB(db)
	if err != nil {
		return nil, fmt.Errorf("create casbin adapter: %w", err)
	}

	enforcer, err := casbin.NewEnforcer(modelPath, adapter)
	if err != nil {
		return nil, fmt.Errorf("create casbin enforcer: %w", err)
	}

	if err := enforcer.LoadPolicy(); err != nil {
		return nil, fmt.Errorf("load casbin policy: %w", err)
	}

	return &Enforcer{inner: enforcer}, nil
}
```

这里做了三件关键事情：

| 步骤 | 作用 |
| --- | --- |
| `TurnOffAutoMigrate` | 不让适配器自动改表结构 |
| `NewAdapterByDB` | 复用项目现有 GORM 连接 |
| `LoadPolicy` | 从 `casbin_rule` 表加载策略到内存 |

::: warning 策略变更后要重新加载
Casbin 判断依赖内存里的策略。当前项目启动时会执行 `LoadPolicy()`。如果后续做了“在线修改权限策略”的页面，写入数据库后还需要设计重新加载策略的机制。
:::

## 权限判断链路

权限中间件的核心调用是：

```go
allowed, err := enforcer.Enforce(roleCode, obj, act)
```

其中：

| 参数 | 来源 |
| --- | --- |
| `roleCode` | 当前用户拥有的角色编码 |
| `obj` | `c.FullPath()`，如果为空则使用 `c.Request.URL.Path` |
| `act` | `c.Request.Method` |

判断过程可以这样理解：

```text
请求 GET /api/v1/system/health
当前用户角色：super_admin
策略表存在：p, super_admin, /api/v1/system/health, GET
结果：允许访问
```

如果用户有多个角色，当前项目会逐个判断：

```go
for _, roleCode := range roleCodes {
	allowed, err := enforcer.Enforce(roleCode, obj, act)
	if err != nil {
		return err
	}

	if allowed {
		return nil
	}
}
```

只要一个角色命中允许策略，就会放行。

## 新增接口权限怎么做

假设新增了一个用户列表接口：

```go
users := api.Group("/users")
users.Use(middleware.Auth(opts.Token, opts.Log))
users.Use(middleware.Permission(opts.DB, opts.Permission, opts.Log))
users.GET("", userHandler.List)
```

它的权限资源通常是：

```text
GET /api/v1/users
```

需要给角色新增一条策略：

```text
p, super_admin, /api/v1/users, GET
```

如果直接写 SQL，可以插入：

```sql
INSERT INTO casbin_rule (ptype, v0, v1, v2, v3, v4, v5)
VALUES ('p', 'super_admin', '/api/v1/users', 'GET', '', '', '');
```

如果接口带路径参数，例如：

```go
users.GET("/:id", userHandler.Detail)
```

策略应写成：

```text
p, super_admin, /api/v1/users/:id, GET
```

因为 `keyMatch2` 能匹配 Gin 风格的 `:id` 参数。

## 常见策略写法

| 策略 | 含义 |
| --- | --- |
| `p, super_admin, /api/v1/system/health, GET` | 允许超级管理员访问系统健康检查接口 |
| `p, admin, /api/v1/users, GET` | 允许管理员查看用户列表 |
| `p, admin, /api/v1/users/:id, GET` | 允许管理员查看某个用户详情 |
| `p, admin, /api/v1/users/:id/update, POST` | 允许管理员更新某个用户 |
| `p, auditor, /api/v1/logs, GET` | 允许审计角色查看日志资源 |

::: warning 慎用 `*`
`*` 会放开同一个资源下的所有 HTTP 方法。本项目后台接口主要使用 `GET` 和 `POST`，除非非常确定，否则优先精确到具体方法。
:::

## Casbin 和菜单权限的边界

接口权限和菜单权限不要混成一件事。

| 类型 | 控制什么 | 推荐落点 |
| --- | --- | --- |
| 接口权限 | 后端接口能不能访问 | Casbin |
| 菜单权限 | 管理台显示哪些菜单和按钮 | 菜单表、角色菜单关系 |

前端菜单隐藏只能改善体验，不能当成安全边界。真正的安全判断必须在后端接口里做。

::: tip 后台权限的基本原则
前端负责“看不看得到”，后端负责“能不能真的执行”。只要是会改数据、查敏感数据的接口，都应该经过后端权限判断。
:::

## 开发时怎么排查权限问题

| 现象 | 优先检查 |
| --- | --- |
| 返回未登录 | JWT 认证中间件是否通过，`Authorization` 是否正确 |
| 返回没有权限 | 当前用户是否绑定启用角色 |
| 有角色仍然没权限 | `casbin_rule` 是否存在对应 `角色 + 路径 + 方法` |
| 路径看起来一样但不命中 | 策略路径是否使用 Gin 路由模式，例如 `:id` |
| 修改策略后没生效 | 当前进程是否重新加载了策略 |
| 所有人都能访问 | 路由是否真的挂了 `Permission` 中间件 |

一个最小检查顺序：

1. 看路由是否在受保护分组下。
2. 看请求是否带有效 JWT。
3. 查 `sys_user_role` 和 `sys_role`，确认用户有启用角色。
4. 查 `casbin_rule`，确认策略存在。
5. 确认策略里的路径和 HTTP 方法与实际路由一致。

## 当前项目的 Casbin 约定

| 场景 | 约定 |
| --- | --- |
| 主体 `sub` | 使用角色编码，例如 `super_admin` |
| 资源 `obj` | 使用 Gin 路由模式路径，例如 `/api/v1/users/:id` |
| 动作 `act` | 使用 HTTP 方法，本项目默认只使用 `GET` 和 `POST` |
| 策略表 | 使用 `casbin_rule` |
| 建表方式 | 参考手册 SQL 手动建表 |
| 策略加载 | 服务启动时 `LoadPolicy()` |
| 多角色判断 | 任意角色命中允许策略即可放行 |
| 菜单权限 | 不放在 Casbin 里，单独建模 |

## 继续查资料

- [Casbin 官方文档](https://casbin.org/docs/overview)
- [Casbin Go 包文档](https://pkg.go.dev/github.com/casbin/casbin/v3)
- [Casbin GORM Adapter 文档](https://pkg.go.dev/github.com/casbin/gorm-adapter/v3)
- [数据库建表语句](./database-ddl#casbin-rule)

## 小结

在当前后台底座里，Casbin 先记住四句话就够用：用户通过 JWT 完成认证；角色编码进入 Casbin 做授权；策略存放在 `casbin_rule`；接口权限按 `角色 + 路径 + 方法` 判断。后续新增业务接口时，只要同步补策略，并确保路由挂上权限中间件，就能支撑基础 RBAC 开发。
