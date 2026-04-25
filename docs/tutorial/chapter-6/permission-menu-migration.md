---
title: 权限、菜单与迁移接入
description: "说明业务模块如何接入接口权限、菜单权限、按钮权限和数据库种子数据。"
---

# 权限、菜单与迁移接入

写完 model、repository、service、handler 和 router 之后，接口能跑通了，但登录后台你会发现：侧边栏看不到新菜单、新接口返回 403、按钮全部隐藏。这些"看不见的水管"就是权限、菜单和种子数据——它们不参与业务逻辑，却决定了一个模块能不能真正用起来。

::: tip 🎯 本节目标
为一个业务模块同时补齐三件事：

1. **接口权限**：角色能访问哪些后端接口（Casbin 策略）。
2. **菜单权限**：侧边栏出现哪些目录、页面和按钮（菜单树）。
3. **数据库结构**：新表怎么建、种子数据怎么写。

验证标准：用 `super_admin` 登录后，侧边栏能看到新菜单，页面内按钮正常显示，接口请求返回 200 而不是 403。
:::

## 接口权限（Casbin）

### 权限是怎么判断的

后端所有需要权限的接口都挂在 `/api/v1/system` 路由分组下，这个分组在注册时挂了两层中间件：

- `middleware.Auth`：从 Token 中解析出当前用户 ID。
- `middleware.Permission`：根据用户角色和请求路径，查 Casbin 策略判断是否放行。

判断逻辑很直接：取当前用户的启用角色编码，对每个角色执行一次 `enforcer.Enforce(roleCode, fullPath, method)`。只要有一个角色命中策略就放行，否则返回 403。

Casbin 模型定义在 `server/configs/rbac_model.conf`：

<<< ../../../server/configs/rbac_model.conf{ini}

策略匹配规则是 `sub == p.sub && keyMatch2(obj, p.obj) && (act == p.act || p.act == "*")`，其中：

- `sub`：角色编码（如 `super_admin`）。
- `obj`：请求路径模板（如 `/api/v1/system/users/:id/update`），支持 `keyMatch2` 路径参数匹配。
- `act`：HTTP 方法（如 `GET`、`POST`），`*` 表示允许所有方法。

所有策略存储在 `casbin_rule` 表，字段 `ptype="p"` 表示基础策略，`v0` 是角色编码，`v1` 是路径，`v2` 是方法。

### 如何为新模块添加权限种子

权限种子定义在 `server/internal/bootstrap/bootstrap.go` 的 `defaultPermissionSeeds` 切片中。应用启动时，`seedDefaultPermissions` 会逐条检查 `casbin_rule` 是否已存在，不存在才插入。

假设新模块的接口路径是 `/api/v1/blog/posts`，需要增加以下种子：

```go
// server/internal/bootstrap/bootstrap.go

var defaultPermissionSeeds = []defaultPermissionSeed{
    // ... 已有的系统权限 ...
    {Path: "/api/v1/blog/posts", Method: "GET"},          // [!code ++]
    {Path: "/api/v1/blog/posts", Method: "POST"},         // [!code ++]
    {Path: "/api/v1/blog/posts/:id/update", Method: "POST"}, // [!code ++]
    {Path: "/api/v1/blog/posts/:id/status", Method: "POST"}, // [!code ++]
}
```

启动服务后，这些策略会被写入 `casbin_rule` 表。`super_admin` 角色会自动拥有这些接口的访问权限。

::: warning ⚠️ 路径必须和路由注册一致
`defaultPermissionSeeds` 中的 `Path` 必须和 `router.go` 中 `system.GET(...)` / `system.POST(...)` 注册的路径完全一致，包括 `:id` 等参数占位符。如果不一致，中间件在 `c.FullPath()` 拿到的路径模板就和策略对不上，即使角色有权限也会返回 403。
:::

::: details 为什么不用迁移文件来管理权限数据
权限策略是运行时数据，不是表结构。Casbin 的 gorm-adapter 负责读写 `casbin_rule`，策略的增删改查通过代码中的种子逻辑来保证幂等（已存在就跳过）。如果把策略放进独立的 SQL 迁移文件，反而需要额外处理"重复执行"和"与 Casbin 内部缓存同步"的问题。
:::

## 菜单树

### 三层结构：目录 → 菜单 → 按钮

菜单权限是一棵树，存储在 `sys_menu` 表中，通过 `parent_id` 形成层级关系：

| 层级 | type 值 | 作用 | 关键字段 |
| --- | --- | --- | --- |
| 目录（Directory） | `1` | 侧边栏分组标题 | `path`、`icon`、`sort` |
| 菜单（Menu） | `2` | 可访问的页面 | `path`、`component`、`icon`、`sort` |
| 按钮（Button） | `3` | 页面内的操作权限 | `code`、`sort` |

以"系统管理"为例，完整的树形结构是：

```text
系统管理 (directory, code=system)
├── 用户管理 (menu, code=system:user, component=system/UserView)
│   ├── 查看用户 (button, code=system:user:list)
│   ├── 创建用户 (button, code=system:user:create)
│   ├── 编辑用户 (button, code=system:user:update)
│   ├── 修改用户状态 (button, code=system:user:status)
│   └── 分配用户角色 (button, code=system:user:assign-role)
├── 角色管理 (menu, code=system:role, component=system/RoleView)
│   └── ...
└── ...
```

### Component 字段与前端路由的映射

菜单的 `component` 字段（如 `system/UserView`）决定了前端加载哪个页面组件。映射关系定义在 `admin/src/router/dynamic-menu.ts` 的 `routeComponentMap` 中：

```ts
// admin/src/router/dynamic-menu.ts
const routeComponentMap: Record<string, RouteComponent> = {
  'system/HealthView': placeholderPage,
  'system/UserView': () => import('../pages/system/UserView.vue'),
  'system/RoleView': () => import('../pages/system/RoleView.vue'),
  // ...
}
```

新增模块页面时，需要在这个 Map 中注册对应的组件路径。如果 `component` 值在 Map 中找不到，页面会回退到 `placeholderPage` 占位组件。

### 按钮权限与 canUse

按钮类型的菜单节点不会出现在侧边栏中，它的 `code` 字段用于前端按钮级权限控制。登录后，前端会从 `/api/v1/auth/menus` 接口拿到当前用户被授权的完整菜单树（包括按钮），`dynamic-menu.ts` 中的 `buttonPermissionCodes` 会递归收集所有按钮编码：

```ts
// admin/src/router/dynamic-menu.ts
export const buttonPermissionCodes = computed(() => {
  return collectButtonCodes(authMenus.value)
})
```

在页面中，通过 `canUse(code)` 判断某个按钮是否应该显示：

```vue
<!-- 只有拥有 system:user:create 权限时才显示"创建用户"按钮 -->
<NButton v-if="canUse('system:user:create')" type="primary" @click="openCreate">
  创建用户
</NButton>
```

每个页面组件中的 `canUse` 函数实现基本一致：

```ts
function canUse(code: string) {
  return buttonPermissionCodes.value.includes(code)
}
```

### 如何为新模块添加菜单种子

菜单种子同样定义在 `server/internal/bootstrap/bootstrap.go` 的 `seedDefaultMenus` 函数中。每个菜单通过 `seedMenu()` 创建，该函数会按 `code` 查重，已存在则直接跳过。

假设新模块是"博客管理"，添加步骤如下：

**第一步**：在文件顶部常量区定义编码。

```go
// server/internal/bootstrap/bootstrap.go

const (
    // ... 已有常量 ...
    defaultBlogMenuCode       = "blog"              // [!code ++]
    defaultBlogPostMenuCode   = "blog:post"          // [!code ++]
    defaultBlogPostListCode   = "blog:post:list"     // [!code ++]
    defaultBlogPostCreateCode = "blog:post:create"   // [!code ++]
    defaultBlogPostUpdateCode = "blog:post:update"   // [!code ++]
    defaultBlogPostStatusCode = "blog:post:status"   // [!code ++]
)
```

**第二步**：在 `seedDefaultMenus` 中按 目录 → 菜单 → 按钮的顺序创建节点。

::: details `seedDefaultMenus` 中新增博客管理模块的代码
```go
// seedDefaultMenus 函数末尾，return menus 之前

// 1. 创建目录
blogMenu, err := seedMenu(db, model.Menu{
    ParentID: 0,
    Type:     model.MenuTypeDirectory,
    Code:     defaultBlogMenuCode,
    Title:    "博客管理",
    Path:     "/blog",
    Icon:     "edit",
    Sort:     20,
    Status:   model.MenuStatusEnabled,
    Remark:   "博客业务目录",
}, log)
if err != nil {
    return nil, err
}

// 2. 创建菜单页面
blogPostMenu, err := seedMenu(db, model.Menu{
    ParentID:  blogMenu.ID,
    Type:      model.MenuTypeMenu,
    Code:      defaultBlogPostMenuCode,
    Title:     "文章管理",
    Path:      "/blog/posts",
    Component: "blog/PostView",
    Icon:      "document",
    Sort:      10,
    Status:    model.MenuStatusEnabled,
    Remark:    "博客文章管理菜单",
}, log)
if err != nil {
    return nil, err
}

// 3. 创建按钮权限
blogPostButtons := []model.Menu{
    {ParentID: blogPostMenu.ID, Type: model.MenuTypeButton, Code: defaultBlogPostListCode, Title: "查看文章", Sort: 10, Status: model.MenuStatusEnabled, Remark: "博客文章按钮"},
    {ParentID: blogPostMenu.ID, Type: model.MenuTypeButton, Code: defaultBlogPostCreateCode, Title: "创建文章", Sort: 20, Status: model.MenuStatusEnabled, Remark: "博客文章按钮"},
    {ParentID: blogPostMenu.ID, Type: model.MenuTypeButton, Code: defaultBlogPostUpdateCode, Title: "编辑文章", Sort: 30, Status: model.MenuStatusEnabled, Remark: "博客文章按钮"},
    {ParentID: blogPostMenu.ID, Type: model.MenuTypeButton, Code: defaultBlogPostStatusCode, Title: "修改文章状态", Sort: 40, Status: model.MenuStatusEnabled, Remark: "博客文章按钮"},
}

menus = append(menus, *blogMenu, *blogPostMenu)
for _, button := range blogPostButtons {
    createdButton, err := seedMenu(db, button, log)
    if err != nil {
        return nil, err
    }
    menus = append(menus, *createdButton)
}
```
:::

**第三步**：前端注册组件映射。

```ts
// admin/src/router/dynamic-menu.ts
const routeComponentMap: Record<string, RouteComponent> = {
  // ... 已有映射 ...
  'blog/PostView': () => import('../pages/blog/PostView.vue'), // [!code ++]
}
```

### 角色菜单绑定

`seedRoleMenus` 函数会把 `seedDefaultMenus` 返回的所有菜单绑定到 `super_admin` 角色。这意味着只要新模块的菜单被加入 `menus` 切片，`super_admin` 就会自动拥有这些菜单和按钮权限——不需要手动在角色管理页面勾选。

对于非超管角色，需要通过"角色管理"页面的"分配菜单权限"功能手动授权。

## 数据库迁移

### 自动建表，无需迁移文件

本项目不使用独立的迁移文件或迁移工具。GORM 在连接数据库时通过 `gorm.Open` 自动完成表结构创建——只要模型结构体定义了 `TableName()` 方法，GORM 就会在首次访问时创建对应的表。

新增业务模块时，只需要在 `server/internal/model/` 下定义模型文件，确保：

1. 结构体有 `gorm` 标签标注字段约束。
2. 实现了 `TableName()` 方法指定表名。
3. GORM 能扫描到该模型（项目通过 `gorm.Open` 时自动检测所有注册的模型）。

::: warning ⚠️ 已有表的结构变更
GORM 的 `AutoMigrate` 只会新增字段，不会删除已有字段，也不会修改字段类型。如果需要变更已有列的定义（比如从 `varchar(64)` 改成 `varchar(128)`），需要手动执行 SQL。新增字段可以直接在结构体中添加，GORM 会自动补列。
:::

::: details 为什么选择不引入迁移工具
对于个人项目，独立迁移文件带来的好处（版本化、回滚）有限，而引入额外工具的配置和维护成本相对更高。当前方式在模型定义和数据库结构之间保持了最短的路径：改模型 → 重启服务 → 表结构同步。如果后续项目规模增长，可以按需引入 `golang-migrate` 或 `goose` 等工具。
:::

### 新模块的模型示例

假设博客模块需要一个 `Post` 模型：

```go
// server/internal/model/post.go
package model

import (
    "time"

    "gorm.io/gorm"
)

type PostStatus int

const (
    PostStatusDraft     PostStatus = 1 // 草稿
    PostStatusPublished PostStatus = 2 // 已发布
)

// Post 是博客文章模型。
type Post struct {
    ID        uint           `gorm:"primaryKey" json:"id"`
    Title     string         `gorm:"size:128;not null" json:"title"`
    Content   string         `gorm:"type:text;not null" json:"content"`
    Status    PostStatus     `gorm:"type:smallint;not null;default:1" json:"status"`
    CreatedAt time.Time      `json:"created_at"`
    UpdatedAt time.Time      `json:"updated_at"`
    DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 固定文章表名。
func (Post) TableName() string {
    return "biz_post"
}
```

定义好模型后，重启服务，GORM 会自动创建 `biz_post` 表。

::: tip 📌 表名前缀约定
系统模块的表名以 `sys_` 前缀（如 `sys_menu`、`sys_role`）。业务模块建议使用 `biz_` 前缀（如 `biz_post`），方便在数据库层面区分系统表和业务表。
:::

## 验证清单

完成权限、菜单和模型的定义后，重启服务，按以下顺序逐一验证：

| # | 检查项 | 验证方式 | 期望结果 |
| --- | --- | --- | --- |
| 1 | 数据库有表 | 连接数据库执行 `\dt` 或查看表列表 | `biz_post` 表存在，字段与模型一致 |
| 2 | 后端路由已注册 | 启动服务，查看控制台日志或直接 curl | 路由路径和方法与 `router.go` 注册一致 |
| 3 | Casbin 策略已写入 | 查询 `casbin_rule` 表 | 新增的 `{role_code, path, method}` 记录存在 |
| 4 | 菜单已写入 | 查询 `sys_menu` 表 | 目录、菜单、按钮节点齐全，`parent_id` 层级正确 |
| 5 | 角色菜单已绑定 | 查询 `sys_role_menu` 表 | `super_admin` 角色绑定了所有新菜单 |
| 6 | 前端侧边栏可见 | 用 `super_admin` 登录后台 | 侧边栏出现新目录和菜单项 |
| 7 | 按钮权限生效 | 打开新模块页面 | 拥有权限的按钮正常显示 |
| 8 | 接口权限生效 | 通过前端操作或 curl 调用新接口 | 返回 200 而不是 403 |

::: warning ⚠️ 权限编码和按钮编码必须前后端一致
后端 `defaultPermissionSeeds` 中的路径、前端 `routeComponentMap` 中的组件键名、按钮 `canUse()` 中的编码字符串——这三者分别和 `casbin_rule.v1`、`sys_menu.component`、`sys_menu.code` 对应。只要有一处拼写不一致（比如 `blog:post:create` vs `blog:posts:create`），就会出现"菜单能点但接口报 403"或"按钮不显示"的问题。

建议在开发新模块时，先把编码设计写在纸上或注释里，统一确认后再写代码，而不是写到一半才发现前后端编码对不上。
:::

## 小结

这一节补齐了让业务模块"从能调用到能用起来"的三个关键环节：

- **接口权限**：在 `defaultPermissionSeeds` 中添加 `{Path, Method}` 种子，启动后自动写入 `casbin_rule`。
- **菜单树**：按 目录 → 菜单 → 按钮三层结构在 `seedDefaultMenus` 中添加节点，`component` 字段对应前端组件映射，按钮 `code` 对应前端 `canUse()`。
- **数据库迁移**：在 `model/` 下定义模型结构体，GORM 自动建表，无需迁移文件。

三件事全部完成后，用 `super_admin` 登录验证：侧边栏有菜单、页面有按钮、接口不报 403。

下一节会把权限和菜单落地到前端页面层面：[前端页面接入流程](./frontend-page-flow)。
