---
title: 日志页面
description: "实现操作日志和登录日志查询页面。"
---

# 日志页面

上一节已经把配置和文件管理接成了真实页面。现在补齐系统管理的最后一组功能：操作日志和登录日志。

完成这一节后，侧边栏里的"操作日志"和"登录日志"不再停留在占位页。操作日志页面展示每次 API 调用的方法、路径、状态码和耗时；登录日志页面展示每次登录尝试的用户、IP、状态和 User-Agent。

::: tip 🎯 本节目标
这一节会把 `system/OperationLogView` 和 `system/LoginLogView` 从占位页换成真实页面，并补齐日志相关的类型和 API 封装。两个页面都是只读查询，不需要新增、编辑或删除操作。
:::

## 先看接口边界

日志接口只有查询，不提供写入：

| 方法 | 路径 | 用途 |
| --- | --- | --- |
| `GET` | `/api/v1/system/operation-logs` | 操作日志分页列表 |
| `GET` | `/api/v1/system/login-logs` | 登录日志分页列表 |

操作日志支持按用户名、请求方法、路径和是否成功筛选。登录日志支持按用户名、IP 和登录状态筛选。

::: details 操作日志是怎么产生的
后端通过 `OperationLog` 中间件自动记录每次 API 调用。请求完成后，中间件会把用户、方法、路径、状态码、耗时和错误信息写入 `sys_operation_log` 表。登录日志则是在登录接口内部写入 `sys_login_log` 表，记录每次登录尝试的结果。
:::

## 本节会改什么

本节会新增或修改下面这些文件：

```text
admin/
└─ src/
   ├─ api/
   │  ├─ operation-log.ts
   │  └─ login-log.ts
   ├─ pages/
   │  └─ system/
   │     ├─ OperationLogView.vue
   │     └─ LoginLogView.vue
   ├─ router/
   │  └─ dynamic-menu.ts
   └─ types/
      ├─ operation-log.ts
      └─ login-log.ts
```

## 开始前先确认

开始之前，先确认下面几件事：

- 已完成上一节 [配置与文件页面](./config-file-pages)。
- 登录后侧边栏能看到"操作日志"和"登录日志"。
- 后端 `/api/v1/system/operation-logs` 和 `/api/v1/system/login-logs` 可以正常返回数据。
- 数据库中已经有若干操作日志和登录日志记录（登录和操作几次后自动产生）。

## 🛠️ 完整代码

下面直接引入本节对应的完整项目文件，默认折叠。需要复制或对照时点击展开即可。

::: details `admin/src/types/operation-log.ts` — 操作日志类型

<<< ../../../admin/src/types/operation-log.ts

:::

::: details `admin/src/types/login-log.ts` — 登录日志类型

<<< ../../../admin/src/types/login-log.ts

:::

::: details `admin/src/api/operation-log.ts` — 操作日志接口

<<< ../../../admin/src/api/operation-log.ts

:::

::: details `admin/src/api/login-log.ts` — 登录日志接口

<<< ../../../admin/src/api/login-log.ts

:::

::: details `admin/src/pages/system/OperationLogView.vue` — 操作日志页面

<<< ../../../admin/src/pages/system/OperationLogView.vue

:::

::: details `admin/src/pages/system/LoginLogView.vue` — 登录日志页面

<<< ../../../admin/src/pages/system/LoginLogView.vue

:::

::: details `admin/src/router/dynamic-menu.ts` — 动态路由映射

修改后，`system/OperationLogView` 和 `system/LoginLogView` 会从占位页切换为真实页面。

<<< ../../../admin/src/router/dynamic-menu.ts

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
2. 点击"系统管理 / 操作日志"，确认日志列表能正常加载。
3. 按用户名筛选，确认只显示匹配记录。
4. 按请求方法筛选（如 POST），确认过滤生效。
5. 按路径关键词筛选，确认模糊匹配正常。
6. 点击某条记录的"详情"，确认能看到错误信息、请求参数或 User-Agent。
7. 进入"系统管理 / 登录日志"，确认日志列表能正常加载。
8. 按用户名、IP 或状态筛选，确认过滤生效。

## 本节小结

这一节把系统管理的日志查询页面补齐了：

- 操作日志页面展示 API 调用记录，支持按用户、方法、路径筛选，点击详情可查看错误信息和请求参数。
- 登录日志页面展示登录尝试记录，支持按用户名、IP 和状态筛选。
- 两个页面都是只读查询，没有新增、编辑或删除操作。
- 日志数据由后端中间件和登录接口自动产生，不需要前端手动写入。

到这里，第 5 章前端管理台的所有页面都已完成。
