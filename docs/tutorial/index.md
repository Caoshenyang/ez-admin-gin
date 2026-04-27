---
title: 从零搭建教程
description: "按阶段记录如何从空仓库搭建一个可复用的通用后台管理系统底座。"
---

# 从零搭建教程

这条主线会从空仓库开始，一步步搭出 `EZ Admin Gin`：一个面向个人项目快速上线的通用后台管理系统底座。

> 当前稳定版本：**v1.0.0** | 最后验证日期：2026-04-27

::: warning 验证优先
每一节都应该能被手动验证。只写“执行命令”不够，还要写清楚执行后应该看到什么。
:::

::: warning ⚠️ 你可能会碰到复制 import 不成功，被自动删除
在跟着教程修改 Go 文件时，如果你只先粘贴 `import` 里的新增依赖，GoLand 可能会因为“暂时未使用”把 import 自动删掉，看起来像是刚粘贴进去就没了。

更稳的方式是直接替换完整文件，或者先粘贴真正使用这些依赖的代码，再让 IDE 自动整理 import。
:::

::: details Q：怎么关闭 GoLand 自动删除未使用 import？
可以关，主要看你是哪个触发方式。

关闭保存时自动删 import：

1. 打开 `File` -> `Settings` -> `Tools` -> `Actions on Save`
2. 找到 `Optimize imports`
3. 取消勾选

这是最常见的“保存后 import 被删”。

关闭输入过程中自动优化 import：

1. 打开 `File` -> `Settings` -> `Go` -> `Imports`
2. 找到 `Optimize imports on the fly`
3. 如果已经勾选，就取消勾选

JetBrains 官方文档也提到，`Optimize Imports` 会移除未使用 import；自动保存时优化在 `Tools | Actions on Save` 中配置。参考：[GoLand Auto import 文档](https://www.jetbrains.com/help/go/creating-and-optimizing-imports.html)。

不过不用急着彻底关。Go 本身不允许未使用 import，IDE 帮你整理通常是好事。跟着教程复制代码时，更稳的方式还是直接替换完整文件，或者先粘贴下面真正使用依赖的代码，再让 GoLand 自动整理 import。
:::

## 怎么推进

当前阶段先完成 `tutorial/` 主线。等主线全部跑通后，再回头整理“开始这里”、参考手册和路线图。

## 章节入口

| 章节 | 目标 | 状态 |
| --- | --- | --- |
| [第 1 章：项目初始化](./chapter-1/) | 建立单仓库结构和三个子项目入口 | 待完善 |
| [第 2 章：后端基础设施](./chapter-2/) | 补齐配置、日志、数据库、Redis、响应和路由基础 | 待完善 |
| [第 3 章：认证与权限](./chapter-3/) | 完成登录、JWT、RBAC、Casbin 和菜单权限设计 | 待完善 |
| [第 4 章：通用系统模块](./chapter-4/) | 实现用户、角色、菜单、配置、文件和日志能力 | 待完善 |
| [第 5 章：前端管理台](./chapter-5/) | 搭建 Vue 3 后台页面、布局、菜单和管理页 | 待完善 |
| [第 6 章：业务模块接入规范](./chapter-6/) | 定义业务模块如何接入底座并实现一个示例模块 | 待完善 |
| [第 7 章：部署与复用](./chapter-7/) | 完成 Docker、Nginx、环境变量、初始化数据和复用说明 | 待完善 |

## 总大纲

完整小节划分见 [教程大纲](./curriculum)。
