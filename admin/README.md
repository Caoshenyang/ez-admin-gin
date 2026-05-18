# EZ Admin — 管理台前端

基于 Vue 3 + TypeScript 的后台管理前端，配合 [ez-admin-gin](https://github.com/Caoshenyang/ez-admin-gin) Go 后端使用。

## 技术栈

| 技术 | 版本 | 用途 |
|------|------|------|
| Vue 3 | 3.5+ | 框架（Composition API + `<script setup>`） |
| TypeScript | 6.0+ | 类型安全 |
| Vite | 8+ | 构建工具 |
| Naive UI | 2.44+ | UI 组件库 |
| Pinia | 3+ | 状态管理 |
| Tailwind CSS | 4+ | 原子化 CSS |
| Vue Router | 5+ | 路由 + 动态菜单 |
| Axios | 1.15+ | HTTP 客户端 |

## 本地启动

```bash
# 1. 安装依赖
pnpm install

# 2. 启动开发服务器（默认 http://localhost:5173）
pnpm dev
```

前端通过 Vite proxy 将 `/api` 请求转发到后端 `http://localhost:8080`，需要后端和数据库先运行。

完整启动流程参见项目根目录 [README](../README.md)。

## 环境变量

| 变量 | 默认值 | 说明 |
|------|--------|------|
| `VITE_API_BASE_URL` | `/api` | 后端 API 地址 |
| `VITE_APP_TITLE` | `EZ Admin` | 浏览器标签页标题 |

通过 `.env.local` 覆盖，不要提交真实密钥。

## 目录结构

```
admin/src/
├── api/              # HTTP 客户端封装 + 生成的 TypeScript 类型
├── modules/          # 业务模块
│   ├── auth/         #   登录认证
│   ├── iam/          #   用户、角色、部门、岗位管理
│   └── system/       #   配置、字典、文件、日志、公告管理
├── layouts/          # 后台布局（侧栏 + 顶栏 + 标签页）
├── router/           # 路由定义 + 动态菜单注册 + 路由守卫
├── stores/           # Pinia 状态（认证、用户、应用）
├── styles/           # 全局样式 + Tailwind 配置
└── utils/            # 工具函数
```

## 权限、菜单、路由的前端接入

1. **登录** → 后端返回 access token + 通过 HttpOnly cookie 设置 refresh token
2. **路由守卫** → `router/guard.ts` 拦截未登录请求，跳转登录页
3. **动态菜单** → 登录后从 `/api/v1/menus/tree` 获取菜单树，根据 `component` 字段动态注册路由
4. **按钮权限** → 通过 `v-permission` 指令或 `usePermission()` composable 控制按钮显隐

## 常用命令

```bash
pnpm dev              # 启动开发服务器
pnpm build            # 类型检查 + 生产构建
pnpm type-check       # TypeScript 类型检查
pnpm lint             # oxlint + ESLint 检查并修复
pnpm generate:api     # 从 Swagger spec 生成 TypeScript 类型
```

## 开发规范

- 使用 `<script setup>` + Composition API
- TypeScript strict 模式
- 页面放在 `src/modules/{module}/pages/{Name}View.vue`
- API 类型由 `pnpm generate:api` 自动生成，不要手动编辑 `src/api/generated.ts`
