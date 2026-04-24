---
title: 角色与菜单页面
description: "实现角色管理和菜单管理前端页面。"
---

# 角色与菜单页面

上一节已经把“用户管理”接成了真实页面：可以查询用户、维护状态，并把用户绑定到角色。现在继续补齐权限体系的另一半：角色管理和菜单管理。

完成这一节后，侧边栏里的“角色管理”和“菜单管理”不再停留在占位页，而是可以进入真实管理页面。角色页面负责维护角色、接口权限和菜单权限；菜单页面负责维护目录、页面菜单和按钮权限。

::: tip 🎯 本节目标
这一节会把 `system/RoleView` 和 `system/MenuView` 从占位页换成真实页面，并补齐角色、菜单相关的类型和 API 封装。页面会参考原型里的“左侧列表 / 树表格 + 右侧权限或编辑面板”结构，保持后台页面紧凑、可扫描。
:::

## 先看接口边界

角色管理接口：

| 方法 | 路径 | 用途 |
| --- | --- | --- |
| `GET` | `/api/v1/system/roles` | 角色分页列表 |
| `POST` | `/api/v1/system/roles` | 创建角色 |
| `POST` | `/api/v1/system/roles/:id/update` | 编辑角色基础信息 |
| `POST` | `/api/v1/system/roles/:id/status` | 修改角色状态 |
| `POST` | `/api/v1/system/roles/:id/permissions` | 替换角色接口权限 |
| `POST` | `/api/v1/system/roles/:id/menus` | 替换角色菜单权限 |

菜单管理接口：

| 方法 | 路径 | 用途 |
| --- | --- | --- |
| `GET` | `/api/v1/system/menus` | 获取完整菜单树 |
| `POST` | `/api/v1/system/menus` | 创建目录、菜单或按钮 |
| `POST` | `/api/v1/system/menus/:id/update` | 编辑菜单 |
| `POST` | `/api/v1/system/menus/:id/status` | 修改菜单状态 |
| `POST` | `/api/v1/system/menus/:id/delete` | 删除菜单 |

::: warning ⚠️ 菜单权限和接口权限是两件事
菜单权限决定“看不看得到入口”，接口权限决定“能不能真的调用接口”。只给角色分配菜单但没有接口权限，页面可能能打开，但请求会被后端拦截；只给接口权限但没有菜单，用户可能有能力访问接口，却没有侧边栏入口。
:::

## 本节会改什么

本节会新增或修改下面这些文件：

```text
admin/
└─ src/
   ├─ api/
   │  ├─ menu.ts
   │  └─ role.ts
   ├─ pages/
   │  └─ system/
   │     ├─ RoleView.vue
   │     └─ MenuView.vue
   ├─ router/
   │  └─ dynamic-menu.ts
   └─ types/
      ├─ menu.ts
      └─ role.ts
```

## 开始前先确认

开始之前，先确认下面几件事：

- 已完成上一节 [用户管理页面](./user-pages)。
- 登录后侧边栏能看到“角色管理”和“菜单管理”。
- 当前账号拥有角色与菜单相关按钮权限。
- 后端 `/api/v1/system/roles` 和 `/api/v1/system/menus` 可以正常返回数据。

## 🛠️ 完整代码

下面直接引入本节对应的完整项目文件。这样阅读时看到的代码和实际项目文件保持一致，也方便后续继续维护。

### 角色类型

`admin/src/types/role.ts`

<<< ../../../admin/src/types/role.ts

### 菜单类型

`admin/src/types/menu.ts`

<<< ../../../admin/src/types/menu.ts

### 角色接口

`admin/src/api/role.ts`

<<< ../../../admin/src/api/role.ts

### 菜单接口

`admin/src/api/menu.ts`

<<< ../../../admin/src/api/menu.ts

### 角色权限页面

`admin/src/pages/system/RoleView.vue`

<<< ../../../admin/src/pages/system/RoleView.vue

### 菜单管理页面

`admin/src/pages/system/MenuView.vue`

<<< ../../../admin/src/pages/system/MenuView.vue

### 动态路由映射

修改 `admin/src/router/dynamic-menu.ts` 后，`system/RoleView` 和 `system/MenuView` 会从占位页切换为真实页面。

<<< ../../../admin/src/router/dynamic-menu.ts

::: warning ⚠️ 按钮权限的 `code` 要和页面判断一致
例如用户页里判断的是 `system:user:create`、`system:user:update` 这类编码。菜单管理页新增按钮节点时，`code` 必须和页面代码里的 `canUse(code)` 保持一致，否则按钮权限不会生效。
:::

::: details 为什么创建菜单需要 `code`，编辑菜单不允许改 `code`
菜单编码会被按钮权限、角色菜单权限和前端权限判断使用。允许随意修改编码，很容易出现“页面还在，按钮突然不显示”的问题。

后端当前的编辑接口也没有接收 `code` 字段，所以前端编辑表单会把编码作为只读信息处理。
:::

## ✅ 验证结果

先启动后端和前端：

::: code-group

```bash [后端]
cd server
go run .
```

```bash [前端]
cd admin
pnpm dev
```

:::

然后按下面顺序验证：

1. 使用 `admin / EzAdmin@123456` 登录。
2. 点击“系统管理 / 角色管理”，确认角色列表和右侧权限树能正常加载。
3. 新建一个测试角色，例如 `demo_operator`，保存后左侧角色列表中能看到它。
4. 给测试角色分配菜单权限和按钮权限，点击“保存权限”。
5. 切换到“接口权限”，给测试角色分配必要接口权限。
6. 进入“系统管理 / 菜单管理”，新增一个测试菜单或按钮，确认树表格会刷新。
7. 回到“用户管理”，把测试用户绑定到这个角色。
8. 退出登录，再用测试用户登录。
9. 确认侧边栏只显示被授权的菜单，页面按钮也按按钮权限显示。

::: details 如果菜单没有变化，先检查这几件事
- 是否重新登录了。当前菜单在登录后加载，修改角色菜单后建议重新登录验证。
- `sys_role_menu` 是否写入了新的菜单 ID。
- `sys_menu.code` 是否和前端 `canUse(code)` 判断一致。
- 角色是否处于启用状态。
- 用户是否已经绑定到刚刚修改的角色。
:::

## 本节小结

这一节把权限体系最关键的两个页面补齐了：

- 角色页面负责维护角色、接口权限和菜单权限。
- 菜单页面负责维护目录、菜单和按钮权限。
- 动态路由通过 `component` 字段加载真实 Vue 页面。
- 菜单权限控制入口，接口权限控制后端访问，按钮权限控制页面操作体验。

下一节继续补齐系统管理里的剩余页面：[配置与文件页面](./config-file-pages)。
