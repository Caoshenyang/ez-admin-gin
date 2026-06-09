---
title: 参考手册
description: "按使用场景整理 EZ Admin Gin 的系统地图、架构约定、接口、数据、前端、部署和质量参考。"
---

# 参考手册

参考手册用于快速查阅稳定事实，不承担完整教程职责。需要照着做时先看教程页，需要确认约定、路径、参数或边界时回到这里。

::: tip 先从系统地图开始
如果你不确定某个能力现在是否存在，先看 [当前系统地图](./current-system-map)。它会告诉你应该回到哪些代码文件确认事实。
:::

## 系统地图

| 参考 | 适合什么时候看 |
| --- | --- |
| [当前系统地图](./current-system-map) | 校准当前模块、页面、路由、配置和文档事实来源 |
| [目录约定](./directory-conventions) | 判断文件应该放到哪里 |
| [模块规范](./module-conventions) | 新增或维护后端模块时确认分层职责 |
| [模块初始化模板](./module-init-template) | 开新模块前确认最小接入清单 |

## 接口与权限

| 参考 | 适合什么时候看 |
| --- | --- |
| [接口风格决策](./api-style-decision) | 查 REST 风格、统一响应和非标准动作接口 |
| [错误码参考](./error-code-reference) | 查业务错误码和 HTTP 状态映射 |
| [权限码约定](./permission-code-conventions) | 对齐菜单 code、按钮权限码和接口权限 |
| [数据权限模型](./data-scope-model) | 理解五档数据范围、Actor 和并集规则 |
| [按钮权限消费示例](./button-permission-consumption-example) | 前端页面消费按钮权限时查最小写法 |

## 数据与初始化

| 参考 | 适合什么时候看 |
| --- | --- |
| [数据库建表语句](./database-ddl) | 找 MySQL / PostgreSQL 完整 SQL 的权威入口 |
| [初始化数据参考](./init-data-reference) | 查角色、菜单、接口权限、字典和管理员初始化链路 |
| [数据库迁移工具选型](./migration-tool-selection) | 理解当前迁移工具选择 |
| [逻辑删除与唯一索引冲突](./logical-delete-and-unique-index) | 处理逻辑删除资源的唯一索引问题 |
| [查询与分页约定](./query-and-pagination-conventions) | 查 `page/page_size`、`keyword`、`status` 等通用筛选 |

## 前端与文件

| 参考 | 适合什么时候看 |
| --- | --- |
| [动态菜单组件白名单](./dynamic-menu-component-reference) | 查菜单 `component`、图标和占位页回退 |
| [上传与公开路径参考](./upload-public-path-reference) | 查上传目录、公开 URL 和 Nginx 代理关系 |

## 部署与运维

| 参考 | 适合什么时候看 |
| --- | --- |
| [环境变量参考](./environment-variables-reference) | 查 `EZ_*` 环境变量、默认值和生产必改项 |
| [Docker 部署文件参考](./deploy-artifacts-reference) | 查 Compose 文件、服务和卷 |
| [Nginx 配置参考](./nginx-config-reference) | 查 HTTP / HTTPS 代理配置 |
| [SSH 隧道连接服务器数据库](./ssh-tunnel-database) | 临时连接远端数据库排查问题 |
| [VitePress 部署到 GitHub Pages](./vitepress-github-pages) | 发布文档站到 GitHub Pages |

## 质量与验证

| 参考 | 适合什么时候看 |
| --- | --- |
| [轻量验证与人工测试](./lightweight-verification) | 改代码、发布或复用前确认该跑哪些检查 |

## 外部资料

| 资料 | 用途 |
| --- | --- |
| [GORM 官方文档](https://gorm.io/docs/) | ORM、模型、查询和事务 |
| [Casbin 官方文档](https://casbin.org/docs/overview) | RBAC/ABAC 策略和模型 |

## 当前收口约定

- 数据库交付以 `server/migrations/{mysql,postgres}/full_schema_and_seed.sql` 为准。
- 首个管理员账号通过 `POST /api/v1/setup/init` 创建，不内置在 SQL 种子里。
- 接口类型以 Swagger 生成结果和 `admin/src/api/generated.ts` 的同步检查为准。
- 更新日志以仓库根目录 `CHANGELOG.md` 为准，不在文档站维护第二份版本记录。
