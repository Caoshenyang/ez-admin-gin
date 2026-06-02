---
title: 数据库建表语句
description: "当前系统表结构与完整版 SQL 交付物入口。"
---

# 数据库建表语句

这页不再把所有 DDL 逐段手工展开维护，而是把权威入口收口到两份完整版 SQL。

::: tip 权威入口
当前数据库交付物以这两份文件为准：

- `server/migrations/mysql/full_schema_and_seed.sql`
- `server/migrations/postgres/full_schema_and_seed.sql`

它们由 `./scripts/build-full-migrations.sh` 生成。生成结果是稳定版初始化文件：表结构直接呈现最终形态，内置种子集中整理，并写入当前迁移基线。
:::

## 为什么改成完整版 SQL 入口

这样做主要是为了避免三类漂移：

- 迁移文件已经更新，但文档里的 DDL 还停留在旧版本
- MySQL 和 PostgreSQL 两套说明页更新不同步
- 文档里的一份“示意 SQL”和仓库里真实可执行 SQL 不是同一份东西

现在的约定是：

- 要执行建库和内置种子初始化，直接使用完整版 SQL
- 要了解结构演进顺序，再回看保留的增量迁移链

## 当前主线表清单

| 表名 | 说明 |
| --- | --- |
| `sys_user` | 后台用户表 |
| `sys_role` | 后台角色表 |
| `sys_user_role` | 用户角色关系表 |
| `sys_menu` | 菜单与按钮表 |
| `sys_role_menu` | 角色菜单关系表 |
| `sys_config` | 系统配置表 |
| `sys_file` | 文件上传记录表 |
| `sys_attachment` | 附件中心表 |
| `sys_notice` | 公告表 |
| `sys_notification` | 站内通知表 |
| `sys_operation_log` | 操作日志表 |
| `sys_login_log` | 登录日志表 |
| `sys_department` | 部门表 |
| `sys_post` | 岗位表 |
| `sys_user_post` | 用户岗位关系表 |
| `sys_role_data_scope` | 角色自定义部门数据范围关系表 |
| `sys_dict_type` | 字典类型表 |
| `sys_dict_item` | 字典项表 |
| `casbin_rule` | Casbin 权限策略表 |

## 使用建议

### 新环境初始化

- MySQL 环境导入 `server/migrations/mysql/full_schema_and_seed.sql`
- PostgreSQL 环境导入 `server/migrations/postgres/full_schema_and_seed.sql`

### 首个管理员创建

完整版 SQL 不会写死管理员账号。导入完成后，需要再调用：

```bash
curl -X POST http://localhost:8080/api/v1/setup/init
```

### 后续维护

如果迁移链发生变化：

1. 先更新对应 `.up.sql`
2. 运行 `./scripts/build-full-migrations.sh`
3. 再同步检查这页、初始化说明页和 README 是否仍然一致

## 相关页面

- [迁移与种子数据](/backend/migration)
- [初始化数据参考](/reference/init-data-reference)
