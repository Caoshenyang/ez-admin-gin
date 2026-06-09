---
title: 前端概览
description: "按当前代码说明 Vue 管理台的目录结构、动态菜单、页面组织、API 类型和本地验证方式。"
---

# 前端概览

前端位于 `admin/`，使用 Vue 3、TypeScript、Naive UI、Pinia、Vite 和 Tailwind CSS。页面按 `auth`、`iam`、`system` 三个模块分组，动态菜单根据后端返回的菜单树生成。

::: tip 🎯 这页怎么读
想新增页面，先看“页面如何被动态菜单发现”；想改接口，先看“API 与类型”；想改页面交互，先看“模块内部结构”。
:::

## 当前技术栈

| 技术 | 当前版本来源 | 用途 |
| --- | --- | --- |
| Vue | `admin/package.json` | UI 框架，Composition API + `<script setup>` |
| Vue Router | `admin/package.json` | 静态路由与动态菜单路由 |
| Pinia | `admin/package.json` | 管理台状态 |
| Naive UI | `admin/package.json` | 表单、表格、弹窗、菜单、布局组件 |
| Vite | `admin/package.json` | 开发服务器和构建 |
| TypeScript | `admin/package.json` | 类型约束 |
| Tailwind CSS | `admin/package.json` | 布局、间距和样式补充 |

## 目录结构

```text
admin/src/
├── api/
│   ├── http.ts              # Axios 实例、认证头、错误处理
│   ├── generated.ts         # Swagger 生成的 API 类型
│   └── contract-check.ts    # 编译期契约检查
├── components/
│   ├── app-shell/           # 侧栏、头部、工作标签、通知抽屉
│   ├── brand/
│   └── ez/                  # 页面、表格、工具栏、空状态等后台基础组件
├── composables/             # 全局组合式函数
├── layouts/AdminLayout.vue
├── modules/
│   ├── auth/
│   ├── iam/
│   └── system/
├── router/
│   ├── index.ts
│   ├── guard.ts
│   ├── dynamic-menu.ts
│   └── route-components.ts
├── stores/
├── styles/
├── types/
└── utils/
```

## 当前页面分组

| 模块 | 页面 |
| --- | --- |
| Auth | 登录页、Dashboard、账户中心 |
| IAM | 用户、角色、菜单、部门、岗位 |
| System | 健康检查、配置、字典、文件、附件、操作日志、登录日志、公告、消息、邮件、占位页 |

站内通知由全局 Shell 的 `NotificationDrawer.vue` 承载，不是独立路由页。

## 模块内部结构

资源型页面默认按这套结构组织：

```text
admin/src/modules/iam/
├── api/             # 请求函数，例如 user.ts / role.ts
├── components/      # 表格、筛选栏、表单弹窗、权限面板
├── composables/     # useUserPage.ts 这类页面状态和副作用
├── pages/           # UserView.vue / RoleView.vue
└── types/           # 资源类型和页面状态类型
```

职责边界：

- `pages/` 只做页面编排，尽量不直接写请求细节。
- `composables/` 管理列表加载、筛选、分页、提交、删除和消息提示。
- `components/` 通过 props 接收数据，通过 emits 上报操作。
- `api/` 只封装请求，不处理页面状态。

## 动态菜单如何发现页面

`admin/src/router/route-components.ts` 使用下面的 glob 自动收集页面：

```ts
const routeModules = import.meta.glob('../modules/**/pages/*View.vue')
```

因此，后端菜单的 `component` 字段通常应匹配：

| 页面文件 | 菜单 component |
| --- | --- |
| `admin/src/modules/iam/pages/UserView.vue` | `iam/UserView` |
| `admin/src/modules/system/pages/FileView.vue` | `system/FileView` |

为了兼容历史菜单数据，前端仍保留 `system/*` 兼容别名。新菜单建议使用真实模块名，减少后续迁移成本。

::: warning 组件不存在时不会白屏
如果菜单里的 `component` 找不到，前端会回退到 `admin/src/modules/system/pages/PlaceholderPage.vue`。这能保护运行时体验，但也说明菜单数据或页面文件需要同步修正。
:::

## 权限与菜单状态

`admin/src/router/dynamic-menu.ts` 负责三件事：

- 把后端授权菜单树转换为 Naive UI 侧栏菜单。
- 从按钮菜单中提取 `buttonPermissionCodes`。
- 为顶部搜索生成页面索引。

页面侧通常通过 `usePermission.ts` 判断按钮是否可见或可用。权限码约定见 [权限码约定](/reference/permission-code-conventions)。

## API 与类型

HTTP 客户端在 `admin/src/api/http.ts`：

- 基础路径是 `/api/v1`。
- 请求拦截器自动注入 `Authorization: Bearer <token>`。
- 认证失效时清理本地认证状态并回到登录页。

接口类型由 Swagger 生成：

```bash
make generate-types
make check-types
```

`make check-types` 会先重新生成后端 Swagger，再确认 `admin/src/api/generated.ts` 没有未提交差异。

## 本地验证

```bash
make admin-check
make admin-build
```

需要联调后端时：

```bash
make docker-up
make server-dev
make admin-dev
```

打开 `http://localhost:5173`，登录后检查 Dashboard、左侧动态菜单、列表页加载、按钮权限和通知抽屉。
