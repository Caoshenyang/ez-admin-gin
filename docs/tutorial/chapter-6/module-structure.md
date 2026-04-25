---
title: 模块固定结构
description: "定义每个业务模块在后端、前端、权限和菜单中的固定位置和命名规则。"
---

# 模块固定结构

当要接入第一个业务模块时，最常见的问题不是接口怎么写，而是文件该放哪里。后端 Model 放哪个包？Handler 跟 Router 要不要拆目录？前端页面、API、类型各放哪？权限编码用什么规则？

这一页把这些问题的答案固定下来。后面所有模块接入都按这个结构走，不需要每次重新判断。

::: tip 🎯 这一页的目标
记住一组目录和命名规则，让你在新增任何业务模块时，不需要猜文件该放哪里、权限该叫什么、菜单该映射到哪个组件。这页本身不写代码，只定义约定。
:::

## 后端目录结构

后端按领域分层，每新增一个业务资源，通常涉及以下位置：

| 层 | 文件位置 | 作用 |
| --- | --- | --- |
| Model | `server/internal/model/{module}.go` | 定义数据库表结构和 GORM 模型 |
| Handler | `server/internal/handler/{group}/{module}.go` | 处理 HTTP 请求、参数绑定和响应 |
| Router | `server/internal/router/router.go` | 注册路由分组和中间件 |
| Bootstrap | `server/internal/bootstrap/bootstrap.go` | 初始化权限种子和菜单种子 |

现有模块的映射：

| 模块 | Model 文件 | Handler 文件 |
| --- | --- | --- |
| 用户 | `model/user.go` | `handler/system/users.go` |
| 角色 | `model/role.go` | `handler/system/roles.go` |
| 菜单 | `model/menu.go` | `handler/system/menus.go` |
| 系统配置 | `model/system_config.go` | `handler/system/configs.go` |

可以看到，系统管理模块统一放在 `handler/system/` 目录下。未来新增业务领域（如博客、商品、订单）时，应该创建新的分组目录，比如 `handler/blog/`、`handler/product/`，而不是全部塞进 `system`。

路由注册在 `router.go` 中通过独立的函数组织，例如 `registerSystemRoutes` 和 `registerAuthRoutes`。新业务模块同样应该有独立的注册函数，比如 `registerBlogRoutes`。

## 前端目录结构

前端每新增一个业务资源，涉及以下四个位置：

| 层 | 文件位置 | 作用 |
| --- | --- | --- |
| 类型 | `admin/src/types/{module}.ts` | 定义接口请求和响应的 TypeScript 类型 |
| API | `admin/src/api/{module}.ts` | 封装接口调用函数 |
| 页面 | `admin/src/pages/{group}/{Module}View.vue` | 页面组件 |
| 路由映射 | `admin/src/router/dynamic-menu.ts` | 将菜单 `component` 字段映射到实际页面组件 |

现有模块的映射：

| 模块 | 类型文件 | API 文件 | 页面文件 |
| --- | --- | --- | --- |
| 用户 | `types/user.ts` | `api/user.ts` | `pages/system/UserView.vue` |
| 角色 | `types/role.ts` | `api/role.ts` | `pages/system/RoleView.vue` |
| 菜单 | `types/menu.ts` | `api/menu.ts` | `pages/system/MenuView.vue` |
| 系统配置 | `types/config.ts` | `api/config.ts` | `pages/system/ConfigView.vue` |

新业务领域同样使用新的分组目录，比如 `pages/blog/PostView.vue`。

路由映射在 `dynamic-menu.ts` 的 `routeComponentMap` 中完成。数据库菜单表的 `component` 字段填 `system/UserView` 这样的值，前端会在 `routeComponentMap` 中找到对应的懒加载组件。新增模块只需要在这个映射表里补一行。

## 命名约定

### 权限编码

权限编码采用三段式格式：`{module}:{resource}:{action}`

| 段 | 含义 | 示例 |
| --- | --- | --- |
| `module` | 所属分组 | `system`、`blog`、`product` |
| `resource` | 业务资源 | `user`、`role`、`config`、`post` |
| `action` | 操作类型 | `list`、`create`、`update`、`status`、`delete` |

实际示例：

- `system:user:list` — 用户列表
- `system:user:create` — 创建用户
- `system:role:update` — 更新角色
- `system:menu:delete` — 删除菜单
- `system:config:value` — 读取配置值

菜单和按钮的 `code` 字段也遵循这个规则。目录类型菜单用两段（如 `system:user`），按钮类型用三段（如 `system:user:create`）。

### 菜单路径与组件映射

菜单的 `path` 字段对应浏览器地址栏路径，`component` 字段对应前端组件：

| 菜单 path | 菜单 component | 前端文件 |
| --- | --- | --- |
| `/system` | — (目录类型无组件) | — |
| `/system/users` | `system/UserView` | `pages/system/UserView.vue` |
| `/system/roles` | `system/RoleView` | `pages/system/RoleView.vue` |
| `/system/configs` | `system/ConfigView` | `pages/system/ConfigView.vue` |

规律：`/system/{resources}` 对应组件 `system/{Resource}View`，其中资源名用单数、组件名用 PascalCase + `View` 后缀。

### API 路径

后端接口路径统一放在 `/api/v1` 下，按分组组织：

| 路径 | 含义 |
| --- | --- |
| `/api/v1/system/users` | 系统管理 - 用户 |
| `/api/v1/system/roles` | 系统管理 - 角色 |
| `/api/v1/system/configs` | 系统管理 - 系统配置 |

新业务模块遵循同样模式：`/api/v1/{group}/{resources}`。分组名和 handler 目录名保持一致。

## 四层映射速查表

下面用现有模块展示完整映射关系，方便新增模块时对照：

<table>
  <colgroup>
    <col style="width: 8rem;">
    <col>
    <col>
    <col>
    <col>
  </colgroup>
  <thead>
    <tr>
      <th>模块</th>
      <th>后端</th>
      <th>前端</th>
      <th>权限编码</th>
      <th>API 路径</th>
    </tr>
  </thead>
  <tbody>
    <tr>
      <td>用户</td>
      <td><code>model/user.go</code><br><code>handler/system/users.go</code></td>
      <td><code>types/user.ts</code><br><code>api/user.ts</code><br><code>pages/system/UserView.vue</code></td>
      <td><code>system:user:list</code><br><code>system:user:create</code><br><code>system:user:update</code><br><code>system:user:status</code></td>
      <td><code>/api/v1/system/users</code></td>
    </tr>
    <tr>
      <td>角色</td>
      <td><code>model/role.go</code><br><code>handler/system/roles.go</code></td>
      <td><code>types/role.ts</code><br><code>api/role.ts</code><br><code>pages/system/RoleView.vue</code></td>
      <td><code>system:role:list</code><br><code>system:role:create</code><br><code>system:role:update</code><br><code>system:role:status</code></td>
      <td><code>/api/v1/system/roles</code></td>
    </tr>
    <tr>
      <td>菜单</td>
      <td><code>model/menu.go</code><br><code>handler/system/menus.go</code></td>
      <td><code>types/menu.ts</code><br><code>api/menu.ts</code><br><code>pages/system/MenuView.vue</code></td>
      <td><code>system:menu:list</code><br><code>system:menu:create</code><br><code>system:menu:update</code><br><code>system:menu:delete</code></td>
      <td><code>/api/v1/system/menus</code></td>
    </tr>
    <tr>
      <td>系统配置</td>
      <td><code>model/system_config.go</code><br><code>handler/system/configs.go</code></td>
      <td><code>types/config.ts</code><br><code>api/config.ts</code><br><code>pages/system/ConfigView.vue</code></td>
      <td><code>system:config:list</code><br><code>system:config:create</code><br><code>system:config:update</code><br><code>system:config:status</code></td>
      <td><code>/api/v1/system/configs</code></td>
    </tr>
  </tbody>
</table>

::: warning ⚠️ 四层命名必须一致
新增模块时，后端分组名、前端目录名、权限编码前缀、API 路径分组这四处的命名必须保持一致。如果后端用 `blog`，前端就不能用 `articles`；权限编码不能用 `article:list`，必须用 `blog:post:list`。命名一旦不统一，菜单能点但接口被拦、权限授权了但页面不显示这类问题就会频繁出现，而且很难排查。
:::

## 新增模块时的检查清单

每次新增业务模块时，按这个顺序确认文件是否到位：

- [ ] 后端 Model：`server/internal/model/{module}.go`
- [ ] 后端 Handler：`server/internal/handler/{group}/{module}.go`
- [ ] 后端路由注册：在 `router.go` 中新增独立的注册函数
- [ ] Bootstrap 种子：在 `bootstrap.go` 中补齐权限和菜单初始数据
- [ ] 前端类型：`admin/src/types/{module}.ts`
- [ ] 前端 API：`admin/src/api/{module}.ts`
- [ ] 前端页面：`admin/src/pages/{group}/{Module}View.vue`
- [ ] 路由映射：在 `dynamic-menu.ts` 的 `routeComponentMap` 中补一行

---

下一节进入后端模块接入的具体步骤：[后端模块接入流程](./backend-module-flow)。
