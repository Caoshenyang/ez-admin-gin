---
title: 权限、菜单与迁移接入
description: "说明业务模块如何接入接口权限、菜单权限、按钮权限和数据库种子数据。"
---

# 权限、菜单与迁移接入

写完 model、handler 和 router 之后，接口能跑通了，但登录后台你会发现：侧边栏看不到新菜单、新接口返回 403、按钮全部隐藏。这些"看不见的水管"就是权限、菜单和种子数据——它们不参与业务逻辑，却决定了一个模块能不能真正用起来。

::: tip 🎯 本节目标
为一个业务模块同时补齐三件事：

1. **接口权限**：角色能访问哪些后端接口（Casbin 策略）。
2. **菜单权限**：侧边栏出现哪些目录、页面和按钮（菜单树）。
3. **数据库结构**：新表怎么建、种子数据怎么写。

验证标准：用 `super_admin` 登录后，侧边栏能看到新菜单，页面内按钮正常显示，接口请求返回 200 而不是 403。
:::

## 接口权限（Casbin）

### 权限是怎么判断的

后端所有需要权限的接口都挂在 `/api/v1/system` 路由分组下，这个分组在注册时挂了三层中间件：

- `middleware.Auth`：从 Token 中解析出当前用户 ID。
- `middleware.OperationLog`：记录请求信息，方便审计和排查。
- `middleware.Permission`：根据用户角色和请求路径，查 Casbin 策略判断是否放行。

判断逻辑很直接：取当前用户的启用角色编码，对每个角色执行一次 `enforcer.Enforce(roleCode, fullPath, method)`。只要有一个角色命中策略就放行，否则返回 403。

Casbin 模型定义在 `server/configs/rbac_model.conf`：

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

策略匹配规则是 `sub == p.sub && keyMatch2(obj, p.obj) && (act == p.act || p.act == "*")`，其中：

- `sub`：角色编码（如 `super_admin`）。
- `obj`：请求路径模板（如 `/api/v1/system/users/:id/update`），支持 `keyMatch2` 路径参数匹配。
- `act`：HTTP 方法（如 `GET`、`POST`），`*` 表示允许所有方法。

所有策略存储在 `casbin_rule` 表，字段 `ptype="p"` 表示基础策略，`v0` 是角色编码，`v1` 是路径，`v2` 是方法。

### 如何为新模块添加权限种子

权限种子通过 SQL 迁移文件管理，位于 `server/migrations/{postgres,mysql}/` 目录下。系统初始权限写在 `000002_seed_data.up.sql` 中，新模块的权限应该写在新的迁移文件里（例如 `000003_blog_seed_data.up.sql`）。

假设新模块的接口路径是 `/api/v1/blog/posts`，需要创建新的迁移文件来添加权限：

::: code-group

```sql [PostgreSQL — 000003_blog_seed_data.up.sql]
INSERT INTO casbin_rule (ptype, v0, v1, v2) VALUES ('p', 'super_admin', '/api/v1/blog/posts', 'GET')
ON CONFLICT (ptype, v0, v1, v2, v3, v4, v5) DO NOTHING;
INSERT INTO casbin_rule (ptype, v0, v1, v2) VALUES ('p', 'super_admin', '/api/v1/blog/posts', 'POST')
ON CONFLICT (ptype, v0, v1, v2, v3, v4, v5) DO NOTHING;
INSERT INTO casbin_rule (ptype, v0, v1, v2) VALUES ('p', 'super_admin', '/api/v1/blog/posts/:id/update', 'POST')
ON CONFLICT (ptype, v0, v1, v2, v3, v4, v5) DO NOTHING;
INSERT INTO casbin_rule (ptype, v0, v1, v2) VALUES ('p', 'super_admin', '/api/v1/blog/posts/:id/status', 'POST')
ON CONFLICT (ptype, v0, v1, v2, v3, v4, v5) DO NOTHING;
```

```sql [MySQL — 000003_blog_seed_data.up.sql]
INSERT IGNORE INTO `casbin_rule` (`ptype`, `v0`, `v1`, `v2`) VALUES ('p', 'super_admin', '/api/v1/blog/posts', 'GET');
INSERT IGNORE INTO `casbin_rule` (`ptype`, `v0`, `v1`, `v2`) VALUES ('p', 'super_admin', '/api/v1/blog/posts', 'POST');
INSERT IGNORE INTO `casbin_rule` (`ptype`, `v0`, `v1`, `v2`) VALUES ('p', 'super_admin', '/api/v1/blog/posts/:id/update', 'POST');
INSERT IGNORE INTO `casbin_rule` (`ptype`, `v0`, `v1`, `v2`) VALUES ('p', 'super_admin', '/api/v1/blog/posts/:id/status', 'POST');
```

:::

启动服务后，golang-migrate 会自动执行新的迁移文件，这些策略会被写入 `casbin_rule` 表。`super_admin` 角色会自动拥有这些接口的访问权限。

::: warning ⚠️ 路径必须和路由注册一致
迁移文件中的 `v1`（路径列）必须和 `router.go` 中 `system.GET(...)` / `system.POST(...)` 注册的路径完全一致，包括 `:id` 等参数占位符。如果不一致，中间件在 `c.FullPath()` 拿到的路径模板就和策略对不上，即使角色有权限也会返回 403。
:::

::: details 为什么权限数据用迁移文件管理
权限策略通过 SQL 迁移文件管理，和表结构变更保持一致的版本化追踪。golang-migrate 通过 `schema_migrations` 表保证幂等——已执行的迁移不会重复运行。如果后续需要通过管理界面动态调整权限，Casbin 的 gorm-adapter 会直接读写 `casbin_rule` 表，和迁移文件互不冲突。
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
  'system/HealthView': () => import('../pages/system/HealthView.vue'),
  'system/UserView': () => import('../pages/system/UserView.vue'),
  'system/RoleView': () => import('../pages/system/RoleView.vue'),
  // ...
}
```

新增模块页面时，需要在这个 Map 中注册对应的组件路径。如果 `component` 值在 Map 中找不到，页面会回退到 `placeholderPage` 占位组件。

### Icon 字段与前端图标白名单

菜单的 `icon` 字段不是前端组件名，而是一个稳定的图标标识，例如 `setting`、`notification`、`layout-dashboard`。`admin/src/router/dynamic-menu.ts` 会先把这个字符串归一化，再去命中前端维护的 `menuIconMap`：

```ts
function resolveMenuIcon(icon: string) {
  return renderMenuIcon(menuIconMap[normalizeMenuIcon(icon)] ?? defaultMenuIcon)
}
```

这条链路有两个好处：

- 数据库只保存业务可读的图标标识，不和具体前端组件实现耦合。
- 如果后端返回了空值或未知值，侧边栏会安全回退到默认图标，不会因为菜单配置错误把渲染打挂。

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

菜单种子通过 SQL 迁移文件管理。新模块的菜单数据应该写在新的迁移文件中（与权限种子可以放在同一个迁移文件里）。

假设新模块是"博客管理"，添加步骤如下：

**第一步**：设计菜单编码和固定 ID。

```text
博客管理目录  — code=blog,     ID=300
文章管理菜单  — code=blog:post, ID=301
  查看文章    — code=blog:post:list,   ID=1100
  创建文章    — code=blog:post:create, ID=1101
  编辑文章    — code=blog:post:update, ID=1102
  修改文章状态 — code=blog:post:status, ID=1103
```

**第二步**：在迁移文件中按 目录 → 菜单 → 按钮的顺序插入数据。

::: details `000003_blog_seed_data.up.sql` 中新增博客管理模块的 SQL（PostgreSQL 版）
```sql
-- 1. 创建目录（ON CONFLICT DO NOTHING 保证幂等）
INSERT INTO sys_menu (id, parent_id, type, code, title, path, component, icon, sort, status, remark, created_at, updated_at)
VALUES (300, 0, 1, 'blog', '博客管理', '/blog', '', 'edit', 20, 1, '博客业务目录', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;

-- 2. 创建菜单页面
INSERT INTO sys_menu (id, parent_id, type, code, title, path, component, icon, sort, status, remark, created_at, updated_at)
VALUES (301, 300, 2, 'blog:post', '文章管理', '/blog/posts', 'blog/PostView', 'document', 10, 1, '博客文章管理菜单', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;

-- 3. 创建按钮权限
INSERT INTO sys_menu (id, parent_id, type, code, title, path, component, icon, sort, status, remark, created_at, updated_at)
VALUES (1100, 301, 3, 'blog:post:list', '查看文章', '', '', '', 10, 1, '博客文章按钮', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;
INSERT INTO sys_menu (id, parent_id, type, code, title, path, component, icon, sort, status, remark, created_at, updated_at)
VALUES (1101, 301, 3, 'blog:post:create', '创建文章', '', '', '', 20, 1, '博客文章按钮', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;
INSERT INTO sys_menu (id, parent_id, type, code, title, path, component, icon, sort, status, remark, created_at, updated_at)
VALUES (1102, 301, 3, 'blog:post:update', '编辑文章', '', '', '', 30, 1, '博客文章按钮', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;
INSERT INTO sys_menu (id, parent_id, type, code, title, path, component, icon, sort, status, remark, created_at, updated_at)
VALUES (1103, 301, 3, 'blog:post:status', '修改文章状态', '', '', '', 40, 1, '博客文章按钮', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;

-- 4. 绑定到 super_admin 角色
INSERT INTO sys_role_menu (role_id, menu_id, created_at, updated_at) VALUES (1, 300, NOW(), NOW())
ON CONFLICT (role_id, menu_id) DO NOTHING;
INSERT INTO sys_role_menu (role_id, menu_id, created_at, updated_at) VALUES (1, 301, NOW(), NOW())
ON CONFLICT (role_id, menu_id) DO NOTHING;
INSERT INTO sys_role_menu (role_id, menu_id, created_at, updated_at) VALUES (1, 1100, NOW(), NOW())
ON CONFLICT (role_id, menu_id) DO NOTHING;
INSERT INTO sys_role_menu (role_id, menu_id, created_at, updated_at) VALUES (1, 1101, NOW(), NOW())
ON CONFLICT (role_id, menu_id) DO NOTHING;
INSERT INTO sys_role_menu (role_id, menu_id, created_at, updated_at) VALUES (1, 1102, NOW(), NOW())
ON CONFLICT (role_id, menu_id) DO NOTHING;
INSERT INTO sys_role_menu (role_id, menu_id, created_at, updated_at) VALUES (1, 1103, NOW(), NOW())
ON CONFLICT (role_id, menu_id) DO NOTHING;
```

MySQL 版本将 `ON CONFLICT (...) DO NOTHING` 替换为 `INSERT IGNORE INTO`，其余写法一致。
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

新模块的菜单需要在迁移文件中显式绑定到 `super_admin` 角色（ID=1），通过 `sys_role_menu` 表的 `INSERT` 语句实现。这意味着只要新模块的菜单绑定写入迁移文件，`super_admin` 就会自动拥有这些菜单和按钮权限——不需要手动在角色管理页面勾选。

对于非超管角色，需要通过"角色管理"页面的"分配菜单权限"功能手动授权。

## 数据库迁移

### golang-migrate 自动建表

本项目使用 golang-migrate 管理数据库迁移。启动时，程序通过 `embed.FS` 嵌入 `server/migrations/{postgres,mysql}/` 目录下的 SQL 文件，并自动执行未应用的迁移。

迁移文件按序号命名，格式为 `NNNNNN_name.up.sql` / `NNNNNN_name.down.sql`：

```text
server/migrations/postgres/
├── 000001_init_schema.up.sql        -- 建所有系统表
├── 000001_init_schema.down.sql      -- 反向 DROP
├── 000002_seed_data.up.sql          -- 角色、菜单、权限、绑定
├── 000002_seed_data.down.sql        -- 反向清空
└── 000003_blog_seed_data.up.sql     -- 新模块的种子数据（示例）
```

::: tip 📌 迁移文件命名规范
- 序号固定 6 位，新迁移递增（`000003`、`000004`...）
- 名称用小写 + 短横线（`add_biz_post`、`blog_seed_data`）
- PostgreSQL 和 MySQL 各维护一份，放在对应子目录
- `_up.sql` 是正向操作，`_down.sql` 是反向回滚
:::

### 新增业务模块时的迁移步骤

新增业务模块时，需要在 `server/internal/model/` 下定义模型，同时创建迁移文件：

1. **定义模型** — 在 `model/` 下新增模型文件，确保有 `gorm` 标签和 `TableName()` 方法。
2. **编写 DDL 迁移** — 创建 `00000X_add_xxx.up.sql`，包含 `CREATE TABLE` 和索引。
3. **编写种子数据迁移** — 如果模块需要初始权限、菜单，创建对应的 `00000X_xxx_seed.up.sql`。
4. **启动验证** — 重启服务，golang-migrate 会自动执行新迁移。

::: warning ⚠️ PostgreSQL 需要重置序列计数器
如果在迁移文件中使用了固定 ID 的 `INSERT`，PostgreSQL 版本需要在末尾添加 `SELECT setval(...)` 语句，确保后续 INSERT 的自增 ID 不会和固定 ID 冲突。
:::

### 新模块的模型示例

假设博客模块需要一个 `Post` 模型：

::: details `server/internal/model/post.go` — 文章模型

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

:::

定义好模型后，还需要在迁移文件中创建对应的表。在 `server/migrations/{postgres,mysql}/` 下新增迁移文件，写入 `biz_post` 的建表语句。重启服务后 golang-migrate 会自动执行，不需要手动建表。

::: tip 📌 表名前缀约定
系统模块的表名以 `sys_` 前缀（如 `sys_menu`、`sys_role`）。业务模块建议使用 `biz_` 前缀（如 `biz_post`），方便在数据库层面区分系统表和业务表。
:::

## 验证清单

完成权限、菜单和模型的定义后，重启服务，按以下顺序逐一验证：

| # | 检查项 | 验证方式 | 期望结果 |
| --- | --- | --- | --- |
| 1 | 数据库有表 | 重启服务执行迁移，再用 `\dt` 或查看表列表 | `biz_post` 表存在，字段与模型一致 |
| 2 | 后端路由已注册 | 启动服务，查看控制台日志或直接 curl | 路由路径和方法与 `router.go` 注册一致 |
| 3 | Casbin 策略已写入 | 查询 `casbin_rule` 表 | 新增的 `{role_code, path, method}` 记录存在 |
| 4 | 菜单已写入 | 查询 `sys_menu` 表 | 目录、菜单、按钮节点齐全，`parent_id` 层级正确 |
| 5 | 角色菜单已绑定 | 查询 `sys_role_menu` 表 | `super_admin` 角色绑定了所有新菜单 |
| 6 | 前端侧边栏可见 | 用 `super_admin` 登录后台 | 侧边栏出现新目录和菜单项 |
| 7 | 按钮权限生效 | 打开新模块页面 | 拥有权限的按钮正常显示 |
| 8 | 接口权限生效 | 通过前端操作或 curl 调用新接口 | 返回 200 而不是 403 |

::: warning ⚠️ 权限编码和按钮编码必须前后端一致
迁移文件中 `casbin_rule.v1` 的路径、前端 `routeComponentMap` 中的组件键名、按钮 `canUse()` 中的编码字符串——这三者分别和 `router.go` 注册的路径、`sys_menu.component`、`sys_menu.code` 对应。只要有一处拼写不一致（比如 `blog:post:create` vs `blog:posts:create`），就会出现"菜单能点但接口报 403"或"按钮不显示"的问题。

建议在开发新模块时，先把编码设计写在纸上或注释里，统一确认后再写代码，而不是写到一半才发现前后端编码对不上。
:::

## 小结

这一节补齐了让业务模块"从能调用到能用起来"的三个关键环节：

- **接口权限**：通过 SQL 迁移文件向 `casbin_rule` 表插入 `{ptype, v0, v1, v2}` 记录，启动时 golang-migrate 自动执行。
- **菜单树**：按 目录 → 菜单 → 按钮三层结构在迁移文件中 `INSERT sys_menu`，`component` 字段对应前端组件映射，按钮 `code` 对应前端 `canUse()`。
- **数据库结构**：在 `model/` 下定义模型结构体，并通过迁移文件创建对应的表和索引。

三件事全部完成后，用 `super_admin` 登录验证：侧边栏有菜单、页面有按钮、接口不报 403。

下一节会把权限和菜单落地到前端页面层面：[前端页面接入流程](./frontend-page-flow)。
