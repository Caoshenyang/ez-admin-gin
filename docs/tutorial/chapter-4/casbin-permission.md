---
title: 接口级权限控制
description: "用 Casbin 把角色编码、接口路径和请求方法连接起来，并对齐当前最终版路由结构。"
---

# 接口级权限控制

上一节已经把“用户拥有多个角色，角色再承接授权能力”这件事建立起来了。现在要把其中一条链路单独做实：角色到底怎么决定一个接口能不能访问。

::: tip 🎯 本节目标
读完这一节，你应该能说清下面这条判断链路：

- 当前请求先经过认证中间件
- `LoadActor` 会把角色编码装进上下文
- `Permission` 中间件把 `角色编码 + 路由路径 + HTTP 方法` 交给 Casbin
- 只要任意一个角色命中允许策略，请求就会放行
:::

## 当前主线里的权限校验链路

接口级权限控制现在已经收敛到下面这条路径：

```text
请求进入 /api/v1/system/*
  ↓
middleware.Auth
  ↓
middleware.LoadActor
  ↓
middleware.OperationLog
  ↓
middleware.Permission
  ↓
Casbin Enforcer
  ↓
允许 / 拒绝
```

它对应的真实代码入口，是 `server/internal/module/system/routes.go`：

```go
system := api.Group("/system")
system.Use(middleware.Auth(opts.Token, opts.Log))
system.Use(middleware.LoadActor(opts.DB, opts.Log))
system.Use(middleware.OperationLog(opts.DB, opts.Log))
system.Use(middleware.Permission(opts.DB, opts.Permission, opts.Log))
```

::: warning ⚠️ 当前 Casbin 主要保护的是 `/api/v1/system/*`
认证模块下的 `/api/v1/auth/me`、`/api/v1/auth/menus`、`/api/v1/auth/dashboard` 目前要求“先登录”，但没有再额外挂 Casbin。

这不是遗漏，而是职责划分：`auth/*` 更偏“登录后的身份消费接口”，`system/*` 才是后台管理操作入口。
:::

## 当前代码落点

这一节现在主要对应下面这些文件：

```text
server/
├─ configs/
│  └─ rbac_model.conf
├─ internal/
│  ├─ bootstrap/
│  │  └─ run.go
│  ├─ middleware/
│  │  └─ permission.go
│  ├─ model/
│  │  └─ casbin_rule.go
│  ├─ permission/
│  │  └─ enforcer.go
│  └─ platform/
│     └─ authz/
│        └─ authz.go
└─ migrations/
   ├─ mysql/
   │  └─ 000002_seed_data.up.sql
   └─ postgres/
      └─ 000002_seed_data.up.sql
```

| 位置 | 职责 |
| --- | --- |
| `configs/rbac_model.conf` | 定义 Casbin 模型 |
| `internal/model/casbin_rule.go` | 映射策略表 `casbin_rule` |
| `internal/permission/enforcer.go` | 创建 Enforcer 并加载数据库策略 |
| `internal/platform/authz/authz.go` | 对外暴露统一的授权命名空间 |
| `internal/middleware/permission.go` | 在请求期执行权限判断 |
| `migrations/*/000002_seed_data.up.sql` | 初始化 `super_admin` 的默认接口权限策略 |

## Casbin 在当前项目里到底判断什么

当前模型只关注三个维度：

```text
sub = 角色编码
obj = 接口路径
act = HTTP 方法
```

一条典型策略长这样：

```text
p, super_admin, /api/v1/system/health, GET
```

含义是：

```text
角色 super_admin 可以用 GET 访问 /api/v1/system/health
```

这里最关键的判断有两个：

- 主体不是用户 ID，而是角色编码
- 接口权限不是菜单权限，它只回答“这个请求能不能进”

## 当前 Casbin 模型文件

`server/configs/rbac_model.conf` 的核心内容如下：

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

你可以把最后一行拆成三条规则理解：

| 条件 | 作用 |
| --- | --- |
| `r.sub == p.sub` | 当前角色必须与策略主体一致 |
| `keyMatch2(r.obj, p.obj)` | 路径支持参数匹配 |
| `r.act == p.act || p.act == "*"` | 方法相同，或者策略允许全部方法 |

::: info 为什么这里强调 `keyMatch2`
因为当前系统路由里有大量 `:id` 形式的路径，例如：

- `/api/v1/system/users/:id/update`
- `/api/v1/system/roles/:id/permissions`

如果不用模式匹配，而是拿某次请求里的真实路径去做字符串全等判断，权限策略会很快失控。
:::

## 请求期是怎么把角色交给 Casbin 的

`middleware.Permission` 的逻辑并不复杂，但非常关键：

1. 先从 `CurrentActor` 里取当前用户的 `RoleCodes`。
2. 如果上下文里没有 `Actor`，再退回到数据库查询当前用户启用角色。
3. 取当前请求的路径和方法。
4. 逐个角色调用 `enforcer.Enforce(...)`。
5. 只要任意一个角色通过，就放行。

核心判断大致可以理解成：

```go
for _, roleCode := range roleCodes {
	allowed, err := enforcer.Enforce(roleCode, obj, act)
	if err != nil {
		return err
	}

	if allowed {
		c.Next()
		return
	}
}
```

这意味着：

- 多角色接口权限按并集处理
- 角色被禁用后，不会再参与当前用户权限判断

## 为什么路径优先取 `c.FullPath()`

当前中间件里有一段容易被忽略，但很重要：

```go
obj := c.FullPath()
if obj == "" {
	obj = c.Request.URL.Path
}
```

这样做是为了优先拿到 Gin 的“路由模式路径”，例如：

| 实际请求 | `c.FullPath()` |
| --- | --- |
| `/api/v1/system/users/12/update` | `/api/v1/system/users/:id/update` |
| `/api/v1/system/roles/5/permissions` | `/api/v1/system/roles/:id/permissions` |

这可以确保 Casbin 策略写的是稳定路由模式，而不是某一次请求里的具体主键值。

## Enforcer 是在哪里初始化的

当前最终版启动入口已经收进 `server/internal/bootstrap/run.go`。启动时会直接初始化授权判断器：

```go
permissionEnforcer, err := authzPlatform.NewEnforcer(db, rbacModelPath)
if err != nil {
	log.Fatal("create permission enforcer", zap.Error(err))
}
```

`internal/platform/authz/authz.go` 又对旧的 `internal/permission` 做了一层平滑包装，所以现在对外应优先把它理解成平台级授权能力，而不是零散的工具函数。

## 默认策略从哪里来

当前项目的默认接口权限不是在代码里硬编码写入，而是通过迁移初始化：

- `server/migrations/mysql/000002_seed_data.up.sql`
- `server/migrations/postgres/000002_seed_data.up.sql`

里面会为 `super_admin` 插入一组 `casbin_rule` 记录，例如：

```text
p, super_admin, /api/v1/system/health, GET
p, super_admin, /api/v1/system/users, GET
p, super_admin, /api/v1/system/roles/:id/permissions, POST
```

这说明当前主线已经把“系统初始可管理状态”放进迁移，而不是要求你第一次启动后再手工点一堆权限。

## 后台是怎么维护角色接口权限的

角色接口权限现在走的是：

```text
POST /api/v1/system/roles/:id/permissions
```

也就是说，角色页不是直接改某个内存策略，而是：

1. 先更新数据库里的策略记录
2. 再由 Enforcer 加载这些策略

当前实现里还有一个非常值得说明的边界：

::: warning ⚠️ 当前角色接口权限改完后，策略不会自动热刷新
现在 Casbin 策略是在服务启动时加载的。角色接口权限修改后，数据库里会更新，但内存里的策略不会自动热刷新。

也就是说，如果你刚改完角色接口权限，想稳定看到效果，当前更可靠的方式仍然是重启后端服务。
:::

## 怎么验证这一节已经成立

### 1. 无权限角色访问系统接口会被拒绝

请求：

```text
GET /api/v1/system/users
```

如果当前角色没有对应权限策略，应返回统一的拒绝结果。

### 2. 超级管理员能访问默认受保护接口

初始化管理员登录后，请求：

```text
GET /api/v1/system/health
```

应当可以通过，这说明默认种子策略已经进入 `casbin_rule`。

### 3. 带路径参数的接口也能按模式命中

例如：

```text
POST /api/v1/system/users/1/update
```

应该能命中类似 `/api/v1/system/users/:id/update` 这类策略，而不是要求每个 ID 都单独写策略。

## 本节最关键的结论

这一节真正要建立的判断是：

> 当前接口权限控制，不是“用户直接对接口做判断”，而是“请求期拿到角色编码，再由 Casbin 按 角色 + 路由模式 + 方法 做统一授权”。

这条链路一旦稳定，后面新增模块时就只需要：

- 固定权限码
- 补默认策略
- 把接口挂到受保护路由组

下一节继续看角色如何承接菜单和按钮权限：[角色菜单权限](./menu-permission)。
