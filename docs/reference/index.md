---
title: 参考手册
description: "集中记录 EZ Admin Gin 的配置、接口、目录约定、模块规范、数据权限和部署参数。"
---

# 参考手册

参考手册用于快速查阅，不承担完整教程职责。这里放固定约定、参数说明、接口说明和需要反复翻阅的内容。

## 已有参考

| 参考 | 说明 |
| --- | --- |
| [GORM 快速入门](./gorm-quick-start) | GORM 基础用法和本项目中的使用方式 |
| [Casbin 快速入门](./casbin-quick-start) | Casbin 权限模型和策略配置 |
| [接口风格决策](./api-style-decision) | RESTful 接口设计决策和统一响应格式 |
| [数据权限模型](./data-scope-model) | 五档数据范围、角色字段和组织体系约定 |
| [数据库迁移工具选型](./migration-tool-selection) | Goose、golang-migrate、Atlas 对比和选型理由 |
| [数据库建表语句](./database-ddl) | 完整建表 SQL 和字段说明 |
| [逻辑删除与唯一索引冲突](./logical-delete-and-unique-index) | 逻辑删除场景下唯一索引的处理方案 |

## 计划补充

- 配置参考（环境变量、配置文件、默认值）
- 目录约定（各子目录职责和文件命名规范）
- 模块规范（model / repository / service / handler / router 分层约定）
- 部署参数（Docker、Nginx、数据库、Redis 配置说明）
