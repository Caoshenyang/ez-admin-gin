---
title: 权限、菜单与迁移接入
description: "对齐当前真实实现，讲清模块如何同时接入 Casbin、菜单树、超级管理员默认授权和前端动态路由。"
---

# 权限、菜单与迁移接入

把模块代码写完后，最常见的错觉是：接口能调通，模块就算接完了。

但对后台系统来说，真正决定一个模块能不能“被使用起来”的，往往是下面这些看不见的连接层：

- Casbin 接口权限
- 菜单树和按钮权限
- 超级管理员默认授权
- 前端动态路由与按钮消费
- 迁移文件里的初始化种子

如果这些地方没接上，就会出现一组很典型的问题：

- 接口返回 `403`
- 侧边栏看不到页面
- 页面能进，但按钮全没了
- 数据库里有表和接口，系统里却像这个模块不存在

::: tip 🎯 本节目标
读完这一节，你应该能顺着当前真实实现走通一条完整接入链路：

- 接口权限写进 `casbin_rule`
- 菜单和按钮写进 `sys_menu`
- 默认角色绑定写进 `sys_role_menu`
- 前端通过 `dynamic-menu.ts` 把页面和按钮权限真正消费起来
:::

## 先看当前主线里的完整权限链路

一个系统模块真正被接进来，当前主线一般会经过下面这条路径：

```text
module/*/policy.go
  ↓
migrations/*/000002_seed_data.up.sql
  ├─ casbin_rule
  ├─ sys_menu
  └─ sys_role_menu
  ↓
/api/v1/auth/menus
  ↓
admin/src/router/dynamic-menu.ts
  ↓
前端页面 / 按钮权限消费
```

这意味着模块接入不是只补后端接口，而是至少要把“接口能访问”和“前端知道它存在”这两层同时打通。

## 第一层：接口权限是怎么接进来的

当前项目的接口权限由 Casbin 负责，核心模型是：

```text
sub = 角色编码
obj = 接口路径
act = HTTP 方法
```

一条典型策略像这样：

```text
p, super_admin, /api/v1/system/notices, GET
```

也就是说：

- 主体是角色编码，例如 `super_admin`
- 资源是 Gin 路由模板路径
- 动作是 `GET` / `POST`

当前这些策略都落在 `casbin_rule` 表中，并通过迁移种子写入。

## 为什么 `policy.go` 很适合先把权限点列清楚

虽然当前 `casbin_rule` 的初始化最终写在 SQL 里，但模块本身最好先在 `policy.go` 里把权限点常量收住。

以公告模块为例：

```go
const (
	PermissionList         = "system:notice:list"
	PermissionCreate       = "system:notice:create"
	PermissionUpdate       = "system:notice:update"
	PermissionUpdateStatus = "system:notice:status"
)
```

注意这里有一个很重要的边界：

- `policy.go` 里的权限码常量更像“按钮 / 菜单语义权限点”
- `casbin_rule` 里的 `path + method` 更像“接口访问权限”

两者不是同一字段，但通常应该保持稳定对应关系。这样前后端和迁移种子才能讲同一种语言。

## 第二层：Casbin 种子到底写到哪里

当前内置系统模块的默认权限，已经直接写在：

- `server/migrations/postgres/000002_seed_data.up.sql`
- `server/migrations/mysql/000002_seed_data.up.sql`

例如公告模块已经有：

```text
/api/v1/system/notices                     GET
/api/v1/system/notices                     POST
/api/v1/system/notices/:id/update          POST
/api/v1/system/notices/:id/status          POST
```

这说明当前主线并不是靠“启动后手动点一遍角色权限”来让系统可用，而是通过迁移种子把最小可管理状态直接固定下来。

::: warning ⚠️ Casbin 路径要和 Gin 路由模板保持一致
种子里的路径应该写成：

- `/api/v1/system/notices/:id/update`

而不是某次请求里的真实路径：

- `/api/v1/system/notices/3/update`

因为权限中间件优先拿的是 `c.FullPath()`，它返回的是 Gin 注册路由时的模板路径。
:::

## 第三层：菜单和按钮为什么也要进迁移种子

只有接口权限还不够。  
因为后台页面并不是通过硬编码侧边栏出现的，而是由菜单树驱动。

当前项目里：

- 目录、页面、按钮都在 `sys_menu`
- 角色和菜单的绑定关系在 `sys_role_menu`
- `/api/v1/auth/menus` 会按当前登录用户返回完整权限树

所以一个模块如果想在系统里“真正出现”，至少还要补两类菜单数据：

1. 页面菜单节点
2. 按钮权限节点

## 当前菜单树里一共分几类节点

当前菜单模型已经固定为三类：

| `type` | 含义 | 作用 |
| --- | --- | --- |
| `1` | 目录 | 侧边栏分组 |
| `2` | 页面菜单 | 动态路由和实际页面入口 |
| `3` | 按钮 | 页面内权限点 |

以公告模块为例，当前种子里已经有：

```text
system:notice                -> 页面菜单
system:notice:list           -> 按钮
system:notice:create         -> 按钮
system:notice:update         -> 按钮
system:notice:status         -> 按钮
```

这说明：

- 菜单不是只管“侧边栏显示”
- 按钮权限也统一挂在同一棵树里

## 第四层：一个页面菜单节点至少要配哪些字段

当前页面级菜单节点最关键的几个字段是：

| 字段 | 作用 |
| --- | --- |
| `code` | 菜单稳定编码 |
| `title` | 菜单标题 |
| `path` | 前端页面路径 |
| `component` | 前端组件映射键 |
| `icon` | 图标标识 |
| `sort` | 排序 |
| `status` | 启用状态 |

公告模块当前真实菜单节点就是：

```text
code      = system:notice
title     = 公告管理
path      = /system/notices
component = system/NoticeView
icon      = notification
```

这几项一旦定稳，后端菜单树、前端动态路由和图标映射才能真正对起来。

## 第五层：`component` 字段为什么必须和前端映射表一致

后端 `sys_menu.component` 并不是任意字符串，它必须能命中前端 `admin/src/router/dynamic-menu.ts` 里的 `routeComponentMap`。

当前公告模块已经有：

```ts
'system/NoticeView': () => import('../pages/system/NoticeView.vue')
```

而种子里写的是：

```text
component = system/NoticeView
```

只要这两边对得上，前端就能把菜单节点转换成真实页面路由。

如果对不上，前端虽然还能渲染菜单，但页面会回退到占位组件，读者就会看到“菜单能点开，但不是目标页面”。

## 第六层：图标字段为什么也要按白名单来

当前后端菜单种子里的 `icon` 值，也不是直接绑定前端组件名，而是命中前端维护的图标白名单。

例如公告模块当前用的是：

```text
icon = notification
```

前端再在 `dynamic-menu.ts` 里把它归一化后命中：

```ts
notification -> NotificationsOutline
```

这条链路的好处是：

- 数据库只存稳定标识
- 前端图标实现可以替换
- 未知图标会安全回退到默认图标

## 第七层：为什么还要补 `sys_role_menu`

把菜单节点写进 `sys_menu` 之后，还不能直接显示。  
因为当前登录用户最终能看到什么，取决于角色菜单绑定。

这就是 `sys_role_menu` 的作用：

```text
角色 ID
  ↓
sys_role_menu
  ↓
菜单 / 按钮节点 ID
```

当前主线里，初始化阶段会把这些系统内置节点直接绑定给 `super_admin`。

这意味着只要系统启动后初始化成功，默认管理员就应该天然能看到这些页面和按钮，而不需要再手动去角色管理页里补一次授权。

## 第八层：前端是怎么消费这棵权限树的

当前前端会通过 `/api/v1/auth/menus` 拿到当前登录用户的完整权限树，然后在 `dynamic-menu.ts` 做两件事：

1. 把页面菜单节点转成动态路由和侧边栏菜单
2. 递归收集按钮节点的 `code`

这就是为什么按钮权限控制现在可以在页面里写成：

```ts
buttonPermissionCodes.value.includes('system:notice:create')
```

而不用额外再请求一条“按钮权限列表”接口。

## 第九层：迁移文件里为什么推荐“一次补全一组模块种子”

当前这套主线最稳的做法，是在迁移文件里一次把同一模块相关的初始化数据补完整：

- `casbin_rule`
- `sys_menu`
- `sys_role_menu`

这样能避免一种常见混乱：

- 接口权限先补了，菜单忘了
- 菜单补了，角色菜单绑定忘了
- 页面能打开，按钮权限没补

对于系统内置模块来说，初始化阶段直接形成“接口可访问 + 页面可见 + 按钮可用”的最小闭环，会比让读者上线后手动补齐稳定得多。

## 当前还有一个必须告诉读者的现状

::: warning ⚠️ Casbin 策略当前仍然是启动时加载
即使你通过接口或手工改了 `casbin_rule`，当前进程里的 Casbin 内存策略也不会自动热刷新。

这意味着：

- 菜单相关数据改完后，前端下一次请求 `/auth/menus` 通常就能看到变化
- 但接口权限改完后，当前服务进程里的 Casbin 规则未必立刻生效

现阶段最稳妥的验证方式，仍然是重启服务后再验证接口权限变化。
:::

这个现状在“模块接入”场景里尤其容易误导人，所以这里必须单独讲清楚。

## 用公告模块回看一遍这条链路

当前公告模块已经把这几层全部串起来了：

### 接口权限

- `/api/v1/system/notices`
- `/api/v1/system/notices/:id/update`
- `/api/v1/system/notices/:id/status`

### 菜单节点

- 页面：`system:notice`
- 按钮：`system:notice:list`
- 按钮：`system:notice:create`
- 按钮：`system:notice:update`
- 按钮：`system:notice:status`

### 角色绑定

- `super_admin` 默认绑定这些菜单与按钮节点

### 前端映射

- `component = system/NoticeView`
- `dynamic-menu.ts` 里有对应页面映射
- `NoticeView.vue` 已真实消费按钮权限码

这说明公告模块不是只在“模块结构”层面成立，也已经把本页这一整套接入链路完整跑通了。

## 接一个新模块时，最小接入检查表

如果你后面要新增一个系统模块，推荐至少按下面顺序检查：

- [ ] `policy.go` 已经列出稳定权限点
- [ ] `casbin_rule` 种子已补接口路径和方法
- [ ] `sys_menu` 已补页面菜单节点
- [ ] `sys_menu` 已补按钮权限节点
- [ ] `sys_role_menu` 已把这些节点绑定给 `super_admin`
- [ ] 前端 `dynamic-menu.ts` 已补 `component` 映射
- [ ] 页面里已按按钮权限码控制关键操作显隐
- [ ] 已考虑 Casbin 当前不热刷新的验证方式

## 本节最关键的结论

这一节真正要建立的判断是：

> 对后台模块来说，权限、菜单、角色绑定和迁移种子不是附属配置，而是模块真正进入系统的一部分。

只要这一层没接完，哪怕模块代码写得再完整，它在后台里也还是“不存在”的。

下一章会继续把这套后端能力对接到前端管理台：[第 7 章：前端企业级管理台](../chapter-7/)。
