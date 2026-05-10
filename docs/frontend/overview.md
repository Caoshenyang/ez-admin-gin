---
title: 前端概览
description: Vue 3 前端的架构分层、技术栈、模块结构
---

# 前端概览

## 技术栈

| 技术 | 版本 | 用途 |
|------|------|------|
| Vue | 3.5 | UI 框架，Composition API + `<script setup>` |
| Naive UI | 2.44 | 组件库（表单、表格、菜单、布局等） |
| Pinia | 3.0 | 状态管理 |
| Vite | 8.0 | 构建工具 |
| Tailwind CSS | 4.2 | CSS 工具类（尺寸、间距、布局补充） |
| Axios | 1.5 | HTTP 客户端 |
| TypeScript | 全量 | 类型安全 |

## 目录结构

```
admin/src/
├── api/                HTTP 客户端封装
│   └── http.ts         Axios 实例、拦截器
├── assets/             静态资源
├── components/         全局共享组件
│   ├── app-shell/      布局壳子（AppSidebar, AppHeader, WorkTabs）
│   └── brand/          品牌 Logo
├── composables/        全局组合式函数
│   └── usePermission.ts 按钮权限
├── constants/          全局常量
├── layouts/            布局组件
│   └── AdminLayout.vue 主布局
├── modules/            业务模块
│   ├── auth/           认证（登录、Dashboard、账户中心）
│   ├── iam/            身份管理（用户、角色、菜单、部门、岗位）
│   └── system/         系统管理（配置、字典、文件、日志、公告）
├── router/             路由
│   ├── index.ts        路由配置
│   ├── guard.ts        路由守卫
│   └── dynamic-menu.ts 动态菜单注册
├── stores/             Pinia 状态
│   └── admin-shell.ts  管理台 Shell 状态
├── styles/             全局样式
├── types/              全局类型
└── utils/              工具函数
    └── auth.ts         Token 管理
```

## 模块结构约定

每个业务模块遵循三层结构：

```
modules/{module}/
├── api/            接口调用，只做 HTTP 请求和类型转换
├── types/          TypeScript 类型定义
├── composables/    状态管理 + 副作用逻辑（useXxxPage）
├── components/     展示组件（通过 props/events 通信）
└── pages/          编排层（拼装 composable + component）
```

**职责边界：**

- **pages/** 只做拼装，不包含业务逻辑
- **components/** 通过 props 接收数据，通过 events 上报交互，不导入 api
- **composables/** 封装全部状态和副作用，返回页面所需的响应式数据和方法
- 全局共享组件放在 `components/`，模块级组件放在各自 `components/` 目录

## API 代理

开发模式下通过 Vite 代理转发 API 请求：

```typescript
// vite.config.ts
server: {
  proxy: {
    '/api': {
      target: 'http://localhost:8080',
      changeOrigin: true,
    }
  }
}
```

## HTTP 客户端

`api/http.ts` 封装了 Axios 实例：

- **请求拦截器**：自动注入 `Authorization: Bearer <token>`
- **响应拦截器**：401/403 自动清除认证状态并跳转登录页
- **基础路径**：`/api/v1`

## 页面组件模式

以用户管理为例，模块内的文件组织：

```
modules/iam/
├── api/
│   └── user.ts                  getUserList, createUser, updateUser, deleteUser
├── types/
│   ├── user.ts                  User, UserStatus
│   └── user-page.ts             ListQuery, ListResponse
├── composables/
│   ├── useUserPage.ts           列表加载、搜索、分页、CRUD 操作
│   └── user-page.utils.ts      辅助函数
├── components/
│   ├── UserFilterBar.vue        搜索栏（关键词、状态筛选）
│   ├── UserTable.vue            数据表格
│   ├── UserFormModal.vue        新增/编辑弹窗
│   └── UserRoleModal.vue        角色分配弹窗
└── pages/
    └── UserView.vue             页面编排（引入 composable + 组件）
```
