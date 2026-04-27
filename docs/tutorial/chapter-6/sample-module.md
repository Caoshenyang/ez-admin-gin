---
title: 示例业务模块
description: "用公告管理模块走完一整条接入链路，证明前面定义的规范可以落地。"
---

# 示例业务模块

前四页已经把模块接入的每一步拆开了：目录放哪、后端怎么接、权限菜单怎么挂、前端页面怎么写。但拆开看和串起来跑是两件事。这一页用一个完整的公告管理模块，从 Model 到页面，把前面所有约定串成一条能跑通的链路。

::: tip 这一页做完你能得到什么
一个完整的公告管理模块，包含后端接口、数据库迁移、权限菜单种子和前端 CRUD 页面。更重要的是，你会看到前面定义的每一条规范在真实代码里是怎么落地的，以后照着这个模式接新模块即可。
:::

## 为什么选公告管理

公告管理是后台系统里常见的轻量模块，数据结构简单、操作清晰，但同时又覆盖了分页查询、关键字搜索、状态切换、新建编辑这些后台页面最常见的交互。用它做示例，既能讲清楚接入流程，又不会因为业务本身太复杂而分散注意力。

## 后端：Model

公告表需要记录标题、正文、排序、状态和备注，同时支持软删除。下面是完整的 Model 定义：

<<< ../../../server/internal/model/notice.go

几个要点：

- `TableName()` 固定表名为 `sys_notice`，与系统表保持 `sys_` 前缀一致。
- `DeletedAt` 使用 `gorm.DeletedAt`，GORM 会自动处理软删除，`json:"-"` 表示不返回给前端。
- `Status` 使用自定义类型 `NoticeStatus`，配合常量 `Enabled = 1` / `Disabled = 2`，让代码语义更清晰。
- 排序字段 `Sort` 默认为 `0`，列表查询时按 `sort ASC, id DESC` 排序。

## 后端：Handler

公告 Handler 包含四个方法：`List`、`Create`、`Update`、`UpdateStatus`，对应分页查询、新建、编辑和状态变更。文件较长，折叠查看：

::: details `server/internal/handler/system/notices.go` — 公告 Handler 完整实现
<<< ../../../server/internal/handler/system/notices.go
:::

Handler 的写法和前面系统模块完全一致，值得关注的几个设计：

- **请求 / 响应结构体定义在文件内部**。`noticeListQuery`、`createNoticeRequest` 等结构体只在 Handler 里使用，不需要对外暴露，所以不放到 Model 包。
- **分页参数归一化**。`normalizeNoticePage` 把非法的 `page` 和 `page_size` 修正为合理值，上限 100，避免一次查太多数据。
- **关键字搜索用 `LIKE`**。公告量通常不大，`LIKE` 足够；如果后续数据量变大，可以换全文检索。
- **`buildNoticeResponse` 统一响应格式**。从 Model 转成响应结构体时集中在一个函数里处理，后续加字段只需要改一处。
- **`Update` 用 `map[string]any` 做批量更新**。GORM 的 `Updates` 方法传入 struct 时会忽略零值字段，用 map 可以避免这个问题。

::: warning 为什么 UpdateStatus 单独拆一个方法
状态变更是高频操作，而且只需要传一个字段。如果复用 Update 方法，前端每次切换状态都要把公告的全部字段回传，既浪费带宽又容易出错。拆出来后，状态切换只需传 `status` 一个值，接口更轻。
:::

## 后端：Router

路由注册只需要在 `registerSystemRoutes` 里新增两行：创建 Handler 实例，注册路由。下面用 diff 标记标出新增的部分：

```go
func registerSystemRoutes(r *gin.Engine, opts Options) {
    health := systemHandler.NewHealthHandler(opts.Config, opts.DB, opts.Redis, opts.Log)
    users := systemHandler.NewUserHandler(opts.DB, opts.Log)
    roles := systemHandler.NewRoleHandler(opts.DB, opts.Log)
    menus := systemHandler.NewMenuAdminHandler(opts.DB, opts.Log)
    configs := systemHandler.NewSystemConfigHandler(opts.DB, opts.Redis, opts.Log)
    files := systemHandler.NewFileHandler(opts.DB, opts.Config.Upload, opts.Log)
    operationLogs := systemHandler.NewOperationLogHandler(opts.DB, opts.Log)
    loginLogs := systemHandler.NewLoginLogHandler(opts.DB, opts.Log)
    notices := systemHandler.NewNoticeHandler(opts.DB, opts.Log) // [!code ++]

    // ... 省略中间代码 ...

    system.GET("/login-logs", loginLogs.List)
    system.GET("/notices", notices.List)              // [!code ++]
    system.POST("/notices", notices.Create)            // [!code ++]
    system.POST("/notices/:id/update", notices.Update) // [!code ++]
    system.POST("/notices/:id/status", notices.UpdateStatus) // [!code ++]
}
```

公告路由注册在 `system` 分组下，自动继承了三条中间件：

1. **Auth** — 验证登录状态，未登录返回 401。
2. **OperationLog** — 记录操作日志，方便审计。
3. **Permission** — 校验角色是否有对应接口的访问权限。

::: details 路由路径为什么要统一用复数
`/notices` 而不是 `/notice`，与已有的 `/users`、`/roles`、`/menus` 保持一致。REST 风格里资源名用复数是常见约定，团队统一一种写法比争论哪一种更正确更有价值。
:::

## 后端：数据库迁移

公告模块需要通过数据库迁移文件来初始化权限种子和菜单种子。创建新的迁移文件来添加公告管理的权限和菜单：

### 权限种子

在 `server/migrations/{postgres,mysql}/` 目录下创建新的迁移文件，添加公告的接口权限：

::: code-group

```sql [PostgreSQL — 000003_notice_seed_data.up.sql]
-- 公告管理接口权限
INSERT INTO casbin_rule (ptype, v0, v1, v2) VALUES ('p', 'super_admin', '/api/v1/system/notices', 'GET');
INSERT INTO casbin_rule (ptype, v0, v1, v2) VALUES ('p', 'super_admin', '/api/v1/system/notices', 'POST');
INSERT INTO casbin_rule (ptype, v0, v1, v2) VALUES ('p', 'super_admin', '/api/v1/system/notices/:id/update', 'POST');
INSERT INTO casbin_rule (ptype, v0, v1, v2) VALUES ('p', 'super_admin', '/api/v1/system/notices/:id/status', 'POST');
```

```sql [MySQL — 000003_notice_seed_data.up.sql]
-- 公告管理接口权限
INSERT INTO `casbin_rule` (`ptype`, `v0`, `v1`, `v2`) VALUES ('p', 'super_admin', '/api/v1/system/notices', 'GET');
INSERT INTO `casbin_rule` (`ptype`, `v0`, `v1`, `v2`) VALUES ('p', 'super_admin', '/api/v1/system/notices', 'POST');
INSERT INTO `casbin_rule` (`ptype`, `v0`, `v1`, `v2`) VALUES ('p', 'super_admin', '/api/v1/system/notices/:id/update', 'POST');
INSERT INTO `casbin_rule` (`ptype`, `v0`, `v1`, `v2`) VALUES ('p', 'super_admin', '/api/v1/system/notices/:id/status', 'POST');
```

:::

### 菜单种子

在同一个迁移文件中，添加公告菜单和按钮：

::: code-group

```sql [PostgreSQL — 000003_notice_seed_data.up.sql]
-- 公告管理菜单
INSERT INTO sys_menu (id, parent_id, type, code, title, path, component, icon, sort, status, remark, created_at, updated_at)
VALUES (20, 1, 2, 'system:notice', '公告管理', '/system/notices', 'system/NoticeView', 'notification', 90, 1, '系统内置菜单', NOW(), NOW());

-- 公告管理按钮
INSERT INTO sys_menu (id, parent_id, type, code, title, path, component, icon, sort, status, remark, created_at, updated_at)
VALUES (1000, 20, 3, 'system:notice:list', '查看公告', '', '', '', 10, 1, '系统内置按钮', NOW(), NOW());
INSERT INTO sys_menu (id, parent_id, type, code, title, path, component, icon, sort, status, remark, created_at, updated_at)
VALUES (1001, 20, 3, 'system:notice:create', '创建公告', '', '', '', 20, 1, '系统内置按钮', NOW(), NOW());
INSERT INTO sys_menu (id, parent_id, type, code, title, path, component, icon, sort, status, remark, created_at, updated_at)
VALUES (1002, 20, 3, 'system:notice:update', '编辑公告', '', '', '', 30, 1, '系统内置按钮', NOW(), NOW());
INSERT INTO sys_menu (id, parent_id, type, code, title, path, component, icon, sort, status, remark, created_at, updated_at)
VALUES (1003, 20, 3, 'system:notice:status', '修改公告状态', '', '', '', 40, 1, '系统内置按钮', NOW(), NOW());

-- 绑定到 super_admin 角色
INSERT INTO sys_role_menu (role_id, menu_id, created_at, updated_at) VALUES (1, 20, NOW(), NOW());
INSERT INTO sys_role_menu (role_id, menu_id, created_at, updated_at) VALUES (1, 1000, NOW(), NOW());
INSERT INTO sys_role_menu (role_id, menu_id, created_at, updated_at) VALUES (1, 1001, NOW(), NOW());
INSERT INTO sys_role_menu (role_id, menu_id, created_at, updated_at) VALUES (1, 1002, NOW(), NOW());
INSERT INTO sys_role_menu (role_id, menu_id, created_at, updated_at) VALUES (1, 1003, NOW(), NOW());
```

```sql [MySQL — 000003_notice_seed_data.up.sql]
-- 公告管理菜单
INSERT INTO `sys_menu` (`id`, `parent_id`, `type`, `code`, `title`, `path`, `component`, `icon`, `sort`, `status`, `remark`, `created_at`, `updated_at`)
VALUES (20, 1, 2, 'system:notice', '公告管理', '/system/notices', 'system/NoticeView', 'notification', 90, 1, '系统内置菜单', NOW(), NOW());

-- 公告管理按钮
INSERT INTO `sys_menu` (`id`, `parent_id`, `type`, `code`, `title`, `path`, `component`, `icon`, `sort`, `status`, `remark`, `created_at`, `updated_at`)
VALUES (1000, 20, 3, 'system:notice:list', '查看公告', '', '', '', 10, 1, '系统内置按钮', NOW(), NOW());
INSERT INTO `sys_menu` (`id`, `parent_id`, `type`, `code`, `title`, `path`, `component`, `icon`, `sort`, `status`, `remark`, `created_at`, `updated_at`)
VALUES (1001, 20, 3, 'system:notice:create', '创建公告', '', '', '', 20, 1, '系统内置按钮', NOW(), NOW());
INSERT INTO `sys_menu` (`id`, `parent_id`, `type`, `code`, `title`, `path`, `component`, `icon`, `sort`, `status`, `remark`, `created_at`, `updated_at`)
VALUES (1002, 20, 3, 'system:notice:update', '编辑公告', '', '', '', 30, 1, '系统内置按钮', NOW(), NOW());
INSERT INTO `sys_menu` (`id`, `parent_id`, `type`, `code`, `title`, `path`, `component`, `icon`, `sort`, `status`, `remark`, `created_at`, `updated_at`)
VALUES (1003, 20, 3, 'system:notice:status', '修改公告状态', '', '', '', 40, 1, '系统内置按钮', NOW(), NOW());

-- 绑定到 super_admin 角色
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`, `created_at`, `updated_at`) VALUES (1, 20, NOW(), NOW());
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`, `created_at`, `updated_at`) VALUES (1, 1000, NOW(), NOW());
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`, `created_at`, `updated_at`) VALUES (1, 1001, NOW(), NOW());
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`, `created_at`, `updated_at`) VALUES (1, 1002, NOW(), NOW());
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`, `created_at`, `updated_at`) VALUES (1, 1003, NOW(), NOW());
```

:::

::: warning 菜单种子的 `Component` 字段必须与前端路由映射一致
迁移文件里写的 `component: "system/NoticeView"` 必须和前端 `dynamic-menu.ts` 中 `routeComponentMap` 的 key 完全匹配。如果这里写 `Notice` 而前端写 `system/NoticeView`，菜单能查到但页面会加载占位组件，不会报错但也不会显示真实页面。这类问题排查起来很费时间，建议在接入新模块时把 `Component` 值直接从前端 `routeComponentMap` 里复制过来。
:::

## 前端：Types

类型定义是前端接入的起点。公告模块的类型文件包含状态枚举、列表项、查询参数、响应结构和请求载荷：

<<< ../../../admin/src/types/notice.ts

类型定义和后端 Model 一一对应，几个设计考虑：

- `NoticeStatus` 用 `as const` 定义常量对象，同时导出类型和值，在模板和逻辑中都能直接使用。
- `NoticeListQuery` 的 `status` 类型写成 `NoticeStatus | 0`，`0` 表示"查询全部"，不传给后端。
- `CreateNoticePayload` 和 `UpdateNoticePayload` 结构相同，但分开定义。如果后续创建和编辑的字段出现差异（比如编辑时多一个版本号），改动不会互相影响。

## 前端：API

API 层负责类型安全的请求封装，每个函数对应一个后端接口：

<<< ../../../admin/src/api/notice.ts

注意接口路径和后端路由的对应关系：

| 前端函数 | HTTP 方法 | 路径 |
| --- | --- | --- |
| `getNotices` | GET | `/system/notices` |
| `createNotice` | POST | `/system/notices` |
| `updateNotice` | POST | `/system/notices/:id/update` |
| `updateNoticeStatus` | POST | `/system/notices/:id/status` |

所有函数都通过 `http` 实例发送请求，自动带上 Token 和错误处理。返回值直接解包为业务数据，页面调用时不需要再处理 `response.data.data`。

## 前端：页面

公告管理页面包含搜索栏、数据表格、分页和弹窗表单，是一个典型的后台 CRUD 页面。文件较长，折叠查看：

::: details `admin/src/pages/system/NoticeView.vue` — 公告管理页面完整代码
<<< ../../../admin/src/pages/system/NoticeView.vue
:::

页面的核心结构可以拆成四个部分来理解：

1. **搜索区** — 关键字输入框 + 状态下拉 + 查询/重置按钮。查询时重置到第一页，重置时清空所有条件。
2. **表格区** — `NDataTable` 使用 `remote` 模式，分页、排序都由后端处理。列定义中用 `render` 函数自定义了标题加粗、状态标签、时间格式化和操作按钮。
3. **分页区** — `NPagination` 放在表格底部，支持切换页码和每页条数。
4. **弹窗表单** — `NModal` + `NForm`，支持新建和编辑两种模式。表单校验规则只要求标题必填。

::: details 按钮权限是怎么生效的
页面上每个操作按钮都用 `canUse('system:notice:create')` 这样的方式控制可见性。`canUse` 函数读取 `dynamic-menu.ts` 中导出的 `buttonPermissionCodes`，这个值是从后端 `/auth/menus` 接口返回的按钮权限列表中收集的。只有当前用户所属角色被授权了对应的按钮权限编码，按钮才会渲染出来。

如果没有看到某个按钮，排查顺序是：角色管理里是否勾选了该按钮权限 → 菜单管理里按钮是否启用 → Bootstrap 里菜单种子是否正确创建。
:::

## 前端：路由映射

最后一步，在 `dynamic-menu.ts` 的 `routeComponentMap` 中加一行，把后端菜单的 `Component` 值映射到实际的 Vue 组件：

```ts
const routeComponentMap: Record<string, RouteComponent> = {
  'system/HealthView': () => import('../pages/system/HealthView.vue'),
  'system/UserView': () => import('../pages/system/UserView.vue'),
  'system/RoleView': () => import('../pages/system/RoleView.vue'),
  'system/MenuView': () => import('../pages/system/MenuView.vue'),
  'system/ConfigView': () => import('../pages/system/ConfigView.vue'),
  'system/FileView': () => import('../pages/system/FileView.vue'),
  'system/OperationLogView': () => import('../pages/system/OperationLogView.vue'),
  'system/LoginLogView': () => import('../pages/system/LoginLogView.vue'),
  'system/NoticeView': () => import('../pages/system/NoticeView.vue'), // [!code ++]
}
```

这一行是菜单能加载到真实页面的关键。`dynamic-menu.ts` 中的 `resolveRouteComponent` 函数会拿后端返回的 `Component` 字段（这里是 `"system/NoticeView"`）去 `routeComponentMap` 里查找对应的懒加载函数。找到就加载真实组件，找不到就降级到占位页面。

## 验证

模块接入完成后，按下面的步骤逐一验证。

### 1. 数据库迁移

先在数据库里手动创建 `sys_notice` 表。当前项目不会在启动时自动帮你建业务表，所以这一步要在 `go run main.go` 之前完成。

下面给一个 PostgreSQL 示例：

```sql
CREATE TABLE sys_notice (
  id BIGSERIAL PRIMARY KEY,
  title VARCHAR(128) NOT NULL,
  content TEXT NOT NULL,
  sort INTEGER NOT NULL DEFAULT 0,
  status SMALLINT NOT NULL DEFAULT 1,
  remark VARCHAR(255) NOT NULL DEFAULT '',
  created_at TIMESTAMPTZ NOT NULL,
  updated_at TIMESTAMPTZ NOT NULL,
  deleted_at TIMESTAMPTZ NULL
);

CREATE INDEX idx_sys_notice_deleted_at
ON sys_notice (deleted_at);

CREATE INDEX idx_sys_notice_status
ON sys_notice (status);
```

执行完成后，可以先用 `\d sys_notice` 或数据库客户端确认字段存在；然后再启动后端服务：

```bash
cd server
go run main.go
```

启动日志中应该能看到类似输出：

```text
INFO	database migrations applied
INFO	server started	{"addr": ":8080", "env": "dev"}
```

### 2. 创建管理员账号

服务启动后，先通过初始化接口创建管理员账号：

```bash
# 创建管理员账号
curl -X POST http://localhost:8080/api/v1/setup/init \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"YourPassword123","nickname":"管理员"}'
```

### 3. 接口验证

使用 `curl` 验证接口是否正常工作。先登录获取 Token：

```bash
TOKEN=$(curl -s http://localhost:8080/api/v1/auth/login \
  -H 'Content-Type: application/json' \
  -d '{"username":"admin","password":"YourPassword123"}' \
  | jq -r '.data.access_token')
```

查询公告列表（应该返回空列表）：

```bash
curl -s http://localhost:8080/api/v1/system/notices \
  -H "Authorization: Bearer $TOKEN" | jq
```

期望输出：

```json
{
  "code": 0,
  "data": {
    "items": [],
    "total": 0,
    "page": 1,
    "page_size": 10
  }
}
```

创建一条公告：

```bash
curl -s -X POST http://localhost:8080/api/v1/system/notices \
  -H "Authorization: Bearer $TOKEN" \
  -H 'Content-Type: application/json' \
  -d '{
    "title": "系统上线公告",
    "content": "后台管理系统已正式上线，欢迎使用。",
    "sort": 0,
    "status": 1,
    "remark": "首条公告"
  }' | jq
```

期望输出中 `data.id` 大于 0，`data.title` 为 `"系统上线公告"`。

再次查询列表，`total` 应为 `1`，`items` 中包含刚创建的记录。

### 3. 前端页面验证

1. 打开浏览器，登录后台管理系统。
2. 侧边栏"系统管理"下应该出现"公告管理"菜单项（图标为 `notification`）。
3. 点击进入，页面顶部显示"公告管理"标题和"+ 新增公告"按钮。
4. 点击"新增公告"，弹窗中填写标题和内容，点击"保存"。
5. 表格中出现新建的公告，状态显示绿色"启用"标签。
6. 点击"禁用"按钮，确认后状态切换为红色"禁用"标签。

::: warning 菜单看不到的排查顺序
如果侧边栏没有出现"公告管理"，按这个顺序检查：

1. 后端是否正常启动，日志里有没有 `default menu created menu_code=system:notice`。
2. 角色管理中 `super_admin` 角色的菜单权限是否包含公告相关条目（Bootstrap 会自动绑定，但如果数据库里已有旧数据，可能需要手动勾选）。
3. 浏览器控制台 Network 面板，查看 `/auth/menus` 接口返回的菜单列表是否包含 `system:notice`。
4. 清除浏览器缓存后重新登录。
:::

## 小结

公告管理模块走完了一条完整的接入链路，涉及的所有文件和改动点可以汇总成一张表：

| 层 | 文件 | 改动类型 |
| --- | --- | --- |
| Model | `server/internal/model/notice.go` | 新增 |
| Handler | `server/internal/handler/system/notices.go` | 新增 |
| Router | `server/internal/router/router.go` | 追加 5 行 |
| Migration | `server/migrations/{postgres,mysql}/000003_notice_seed_data.up.sql` | 新增 |
| Types | `admin/src/types/notice.ts` | 新增 |
| API | `admin/src/api/notice.ts` | 新增 |
| Page | `admin/src/pages/system/NoticeView.vue` | 新增 |
| Route | `admin/src/router/dynamic-menu.ts` | 追加 1 行 |

这就是[模块固定结构](./module-structure)里定义的约定在真实代码里的落地方式。以后接入新模块，按同样的顺序和结构走一遍就行：先写 Model，再写 Handler，然后注册路由，创建数据库迁移文件来初始化权限和菜单，最后接前端。

回到本章目录：[第 6 章：业务模块接入规范](./index)。
