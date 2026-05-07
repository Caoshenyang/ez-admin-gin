---
title: 项目结构
description: "EZ Admin Gin 的技术栈组成、当前目录结构，以及企业级底座主线下已经稳定下来的后端骨架说明。"
---

# 项目结构

::: tip 🎯 这页解决什么
帮你快速了解项目用了哪些技术、各目录负责什么，方便后续定位文件。
:::

## 技术栈

| 层 | 技术 |
| --- | --- |
| 后端 | Go 1.26、Gin、GORM、PostgreSQL、Redis、Casbin |
| 前端 | Vue 3.5、TypeScript、Naive UI、TailwindCSS 4、Vite 8 |
| 文档 | VitePress 2.0 |
| 部署 | Docker Compose、Nginx |

## 目录结构

```
ez-admin-gin/
├── server/          # Go 后端
│   ├── configs/     # 配置文件（config.yaml）
│   ├── internal/    # 启动装配、平台能力和业务模块
│   └── migrations/  # 数据库迁移
├── admin/           # Vue 3 前端
│   └── src/         # 页面、组件、路由、状态管理
├── docs/            # VitePress 文档站
├── deploy/          # Docker Compose 和 Nginx 配置
│   └── nginx/
└── .agents/         # 开发辅助 skill
```

各目录职责：

- **server/** — Go 后端，入口 `main.go` 使用 embed 嵌入迁移文件
- **admin/** — Vue 3 前端管理台，页面、组件、路由和状态管理都在 `src/` 下
- **docs/** — VitePress 文档站，就是你现在在读的站点
- **deploy/** — Docker Compose 文件和 Nginx 反向代理配置，分为本地开发环境和生产环境
- **.agents/** — 开发辅助工具配置，正常使用不需要关注

## 当前最终后端骨架

当前后端主线已经收敛到下面这套骨架，不再是“准备迁过去”的草图，而是后续继续扩模块时应该直接遵守的结构：

```text
server/
├── cmd/
│   └── server/
├── internal/
│   ├── bootstrap/
│   ├── platform/
│   │   ├── config/
│   │   ├── database/
│   │   ├── logger/
│   │   ├── redis/
│   │   ├── authn/
│   │   ├── authz/
│   │   └── datascope/
│   └── module/
│       ├── auth/
│       ├── setup/
│       ├── iam/
│       ├── system/
```

这样调整的核心原因只有一个：让认证、初始化、IAM、系统能力、组织体系、数据权限和后续业务模块扩展都有明确落点。

## 当前已经落地了什么

当前这轮升级已经不是“第一阶段刚开始”，而是已经落地了下面这些稳定结构：

- 新增 `bootstrap`，把启动装配从 `main.go` 中抽出来
- 新增 `platform` 命名空间，承接配置、数据库、日志、Redis、认证、鉴权和数据权限基础设施
- 新增 `module/auth`、`module/setup`、`module/iam/*`、`module/system/*` 等稳定模块边界
- 新增部门、岗位、用户岗位、角色数据范围的模型与迁移
- 核心系统模块已经按 `dto / repository / service / handler / routes / policy / datascope` 主线收敛

这也意味着：后续新增常用模块时，不应该再发明新的目录层级，而是优先沿这套骨架继续扩。

## 这套骨架已经证明过什么

当前这套 `module/*` 结构已经不只是“能放系统模块”的理论骨架，而是已经被几类真实能力验证过：

- `module/system/dict` 证明字典类型与字典项这种双表资源可以稳定落在标准模块结构里
- `module/auth` 下的账户中心证明“当前登录人自助能力”不需要再复制一套后台管理员模块
- `module/system/attachment` 证明业务化资源层可以复用既有底层上传能力，而不是重新发明一套文件链路
- 历史上的业务示例模块曾证明非 `system` 分组资源也能沿同一条主线接入菜单、按钮、前端页面和数据权限；当前主线保留的是这套接入规范，而不是继续内置 CRM 目录

也就是说，这套目录骨架现在已经覆盖了：

- 平台入口
- IAM 与系统基础模块
- 当前登录人维度的自助能力
- 非 `system` 分组的业务模块接入模式

后续继续扩模块时，应该直接沿这些已成立的路径延伸，而不是再开并行结构。

## 下一批模块应该怎么接

下一批可复用模块如果继续补，也应该保持同一套判断顺序：

1. 先判断它属于 `auth`、`iam`、`system` 还是独立业务分组
2. 再按 `dto / repository / service / handler / routes / policy / datascope` 落后端
3. 菜单、按钮、Casbin 与初始化种子一起补，不把权限接入推迟到最后
4. 前端同步补 `types / api / pages / dynamic-menu`

做到这里，一个新模块才算真正进入当前后台底座，而不是只多了几张表和一个接口。

## 怎么继续读

- 想理解为什么一定要做这次结构升级：看 [企业级架构升级](/guide/enterprise-architecture)
- 想从 Java 视角理解这套结构：看 [Go vs Java 工程结构](/guide/java-to-go-structure)
