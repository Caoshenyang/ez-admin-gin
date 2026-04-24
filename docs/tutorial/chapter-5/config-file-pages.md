---
title: 配置与文件页面
description: "实现系统配置管理页面和文件上传管理页面。"
---

# 配置与文件页面

上一节已经把角色和菜单接成了真实页面。现在继续补齐系统管理剩下的两个功能页：配置管理和文件管理。

完成这一节后，侧边栏里的"配置管理"和"文件管理"不再停留在占位页。配置页面负责维护系统键值配置，按分组归类管理；文件页面负责上传附件、查看文件列表和复制文件链接。

::: tip 🎯 本节目标
这一节会把 `system/ConfigView` 和 `system/FileView` 从占位页换成真实页面，并补齐配置和文件相关的类型和 API 封装。配置页面采用搜索 + 数据表 + 弹框表单布局；文件页面使用上传按钮 + 数据表布局，支持按文件名和类型筛选。
:::

## 先看接口边界

配置管理接口：

| 方法 | 路径 | 用途 |
| --- | --- | --- |
| `GET` | `/api/v1/system/configs` | 配置分页列表 |
| `POST` | `/api/v1/system/configs` | 创建配置 |
| `POST` | `/api/v1/system/configs/:id/update` | 编辑配置 |
| `POST` | `/api/v1/system/configs/:id/status` | 修改配置状态 |

文件管理接口：

| 方法 | 路径 | 用途 |
| --- | --- | --- |
| `GET` | `/api/v1/system/files` | 文件分页列表 |
| `POST` | `/api/v1/system/files` | 上传文件 |

::: warning ⚠️ 配置键创建后不可更改
配置键（`key`）会被后端缓存到 Redis，也被其他模块引用。允许随意修改键名，容易导致缓存失效或引用断裂。所以编辑模式下，配置键是只读字段。
:::

## 本节会改什么

本节会新增或修改下面这些文件：

```text
admin/
└─ src/
   ├─ api/
   │  ├─ config.ts
   │  └─ file.ts
   ├─ pages/
   │  └─ system/
   │     ├─ ConfigView.vue
   │     └─ FileView.vue
   ├─ router/
   │  └─ dynamic-menu.ts
   └─ types/
      ├─ config.ts
      └─ file.ts
```

## 开始前先确认

开始之前，先确认下面几件事：

- 已完成上一节 [角色与菜单页面](./role-menu-pages)。
- 登录后侧边栏能看到"配置管理"和"文件管理"。
- 当前账号拥有配置与文件相关按钮权限。
- 后端 `/api/v1/system/configs` 和 `/api/v1/system/files` 可以正常返回数据。

## 🛠️ 完整代码

下面直接引入本节对应的完整项目文件，默认折叠。需要复制或对照时点击展开即可。

::: details `admin/src/types/config.ts` — 配置类型

<<< ../../../admin/src/types/config.ts

:::

::: details `admin/src/types/file.ts` — 文件类型

<<< ../../../admin/src/types/file.ts

:::

::: details `admin/src/api/config.ts` — 配置接口

<<< ../../../admin/src/api/config.ts

:::

::: details `admin/src/api/file.ts` — 文件接口

<<< ../../../admin/src/api/file.ts

:::

::: details `admin/src/pages/system/ConfigView.vue` — 配置管理页面

<<< ../../../admin/src/pages/system/ConfigView.vue

:::

::: details `admin/src/pages/system/FileView.vue` — 文件管理页面

<<< ../../../admin/src/pages/system/FileView.vue

:::

::: details `admin/src/router/dynamic-menu.ts` — 动态路由映射

修改后，`system/ConfigView` 和 `system/FileView` 会从占位页切换为真实页面。

<<< ../../../admin/src/router/dynamic-menu.ts

:::

::: details 为什么配置编辑不允许改键
配置键会被后端缓存到 Redis，其他模块也可能通过 `GET /api/v1/system/configs/value/:key` 按键读取值。如果允许编辑键名，需要同步清理旧缓存、写入新缓存，还要更新所有引用方。后端当前的编辑接口也没有接收 `key` 字段，所以前端编辑表单会把键作为只读信息处理。
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
2. 点击"系统管理 / 配置管理"，确认配置列表能正常加载。
3. 点击"+ 新增配置"，填写分组、键、名称和值，保存后列表中出现新配置。
4. 点击编辑，确认键为只读，修改名称和值后保存成功。
5. 点击"禁用"按钮，确认状态切换成功。
6. 进入"系统管理 / 文件管理"，确认文件列表能正常加载。
7. 点击"上传文件"，选择一张图片或文档，确认上传成功后列表中出现新记录。
8. 点击复制链接按钮，确认剪贴板中包含文件 URL。

::: details 如果上传失败，先检查这几件事
- 文件大小是否超过后端配置的 `max_size_mb`（默认 10 MB）。
- 文件扩展名是否在 `allowed_exts` 列表中。
- `uploads` 目录是否有写入权限。
- 后端控制台是否有 `save file` 相关错误日志。
:::

## 本节小结

这一节把系统管理剩余的两个页面补齐了：

- 配置页面负责维护分组键值配置，支持搜索、筛选、新增、编辑和状态切换。
- 文件页面负责上传附件、查看文件列表和复制文件链接。
- 配置键创建后不可更改，避免缓存和引用断裂。
- 文件上传通过 `multipart/form-data` 提交，后端自动生成文件名并计算 SHA256 校验。

到这里，第 5 章前端管理台的所有基础页面都已完成。下一节继续补齐日志查询页面：[日志页面](./log-pages)。
