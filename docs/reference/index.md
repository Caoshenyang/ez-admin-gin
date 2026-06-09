---
title: 参考手册
description: "集中记录 EZ Admin Gin 的配置、接口、目录约定、模块规范、数据权限和部署参数。"
---

# 参考手册

参考手册用于快速查阅，不承担完整教程职责。这里保留交付、部署、扩展和长期维护时需要反复回看的稳定说明。

## 已有参考

| 参考 | 说明 |
| --- | --- |
| [GORM 快速入门](https://gorm.io/docs/) | GORM 官方文档 |
| [Casbin 快速入门](https://casbin.org/docs/overview) | Casbin 官方文档 |
| [接口风格决策](./api-style-decision) | RESTful 接口设计决策和统一响应格式 |
| [数据权限模型](./data-scope-model) | 五档范围、Actor、并集规则和资源接法 |
| [环境变量参考](./environment-variables-reference) | `EZ_*` 变量、默认值和部署优先检查项 |
| [权限码约定](./permission-code-conventions) | `policy.go`、菜单 code、component 和按钮权限码关系 |
| [错误码参考](./error-code-reference) | 统一响应体里的业务错误码、HTTP 映射和常见返回场景 |
| [目录约定](./directory-conventions) | 顶层目录、`server/internal` 主骨架和历史兼容区边界 |
| [模块规范](./module-conventions) | `routes / handler / service / repository / dto / policy / datascope` 固定分工 |
| [初始化数据参考](./init-data-reference) | 首批角色、菜单、Casbin 策略与管理员初始化接口的真实链路 |
| [动态菜单组件白名单](./dynamic-menu-component-reference) | `component`、`icon`、占位页回退和按钮权限来源 |
| [按钮权限消费示例](./button-permission-consumption-example) | `collectButtonCodes(...)`、`buttonPermissionCodes` 和页面 `canUse(code)` 的最小稳定写法 |
| [轻量验证与人工测试](./lightweight-verification) | 当前质量策略、轻量验证命令和发布前人工测试清单 |
| [上传与公开路径参考](./upload-public-path-reference) | `upload.dir`、`upload.public_path`、`sys_file.path/url` 和 `/uploads/` 代理关系 |
| [模块初始化模板](./module-init-template) | 新模块从目录、权限、菜单到前端页面接入的最小开工模板 |
| [查询与分页约定](./query-and-pagination-conventions) | `page/page_size`、`keyword`、`status` 与模块扩展筛选边界 |
| [数据库迁移工具选型](./migration-tool-selection) | Goose、golang-migrate、Atlas 对比和选型理由 |
| [数据库建表语句](./database-ddl) | MySQL / PostgreSQL 完整版 SQL 的权威入口 |
| [逻辑删除与唯一索引冲突](./logical-delete-and-unique-index) | 逻辑删除场景下唯一索引的处理方案 |

## 当前收口约定

- 数据库交付以 `server/migrations/{mysql,postgres}/full_schema_and_seed.sql` 为准
- 首个管理员账号通过 `POST /api/v1/setup/init` 创建，不内置在 SQL 种子里
- 更新日志以仓库根目录 `CHANGELOG.md` 为准，不再维护第二份文档站版本
