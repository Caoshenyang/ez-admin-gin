---
title: 接口级权限控制
description: "用 Casbin 把角色编码、接口路径和请求方法连接起来，并对齐当前最终版路由结构。"
---

# 接口级权限控制

::: warning 历史页说明
这页保留的是旧章节位置下的 Casbin 权限说明，用来兼容早期引用。

当前教程主线里的 canonical 位置已经切到 [第 4 章：接口级权限控制](../chapter-4/casbin-permission)。
:::

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

请求体结构很直接：

```json
{
  "permissions": [
    { "path": "/api/v1/system/users", "method": "GET" },
    { "path": "/api/v1/system/users", "method": "POST" }
  ]
}
```

服务端会做两件事：

1. 先把这组 `path + method` 规范化去重。
2. 再替换 `casbin_rule` 里该角色编码对应的全部 `p` 策略。

也就是说，当前接口不是“增量补一条权限”，而是“用一份完整权限集覆盖当前角色接口权限”。

## 一个需要明确告诉读者的当前现状

::: warning ⚠️ 当前实现还没有做 Casbin 策略热刷新
`Enforcer` 现在只会在服务启动时执行一次 `LoadPolicy()`。

这意味着：

- 修改 `sys_user_role`、`sys_role_menu`、`sys_role_data_scope` 这类直接走数据库读取的关系，通常可以在后续请求里立即体现
- 修改 `casbin_rule` 后，当前进程里的 Casbin 内存策略不会自动刷新

如果你通过 `/api/v1/system/roles/:id/permissions` 改了接口权限，想让结果立刻生效，目前最稳妥的做法仍然是重启服务；后续如果要做在线权限管理，这里还需要补“策略重载”能力
:::

这一点在文档里一定要讲清楚。否则读者会误以为“接口返回成功了，权限就已经实时切换”，然后卡在一个很难定位的问题上。

## 新增一个受保护接口时，应该怎么接

如果你后面继续扩模块，建议按下面顺序接入接口权限：

1. 把新接口挂到需要保护的路由组下，或者显式挂上 `middleware.Permission(...)`。
2. 确认策略里的路径使用 Gin 路由模式，而不是具体请求路径。
3. 给初始化角色补默认策略，或者通过角色权限接口写入策略。
4. 如果你刚改的是 `casbin_rule`，记得当前版本需要重启服务才能稳定验证。

一个最小判断标准是：

> 接口是否受保护，不是看“你有没有写 Handler”，而是看“这个请求有没有进入统一的认证与授权链路”。

## 怎么验证这一节已经做成

### 1. 管理员访问受保护系统接口应当成功

登录拿到 Token 后，请求：

```text
GET /api/v1/system/health
```

默认管理员拥有 `super_admin` 角色，应该可以正常访问。

### 2. 没有 Token 时应当直接被拒绝

同样的接口如果不带 Token，请求会先被 `middleware.Auth` 拦住，而不会进入后续 Casbin 判断。

### 3. 权限不足时应当返回统一拒绝

如果某个角色没有对应 `casbin_rule` 策略，请求应当返回“没有权限访问”，而不是 404，也不应该绕过中间件直接执行 Handler。

### 4. 修改角色接口权限后，要按当前实现方式验证

如果你刚调用过：

```text
POST /api/v1/system/roles/:id/permissions
```

那验证时要额外确认一件事：当前服务进程是否已经重新加载策略。就现阶段实现来说，最直接的验证方式仍然是重启服务后再测一次。

## 本节最关键的结论

这一节真正要记住的是：

> 当前项目的接口权限，不是“用户直接拥有接口访问能力”，而是“用户通过角色编码命中 Casbin 策略后，才拥有接口访问能力”。

只要这条主线保持稳定，后续再扩系统模块、业务模块或者角色管理页，权限体系就不会散掉。

下一节继续看角色如何承接菜单和按钮权限：[角色菜单权限](./menu-permission)。
