---
title: 迁移与种子数据
description: 数据库迁移机制、种子数据内容、双数据库支持
---

# 迁移与种子数据

## 迁移机制

使用 `golang-migrate` 管理数据库版本，支持 MySQL 和 PostgreSQL 双方言。

### 迁移文件结构

```
server/migrations/
├── mysql/
│   ├── 000001_init_schema.up.sql
│   ├── 000001_init_schema.down.sql
│   ├── 000002_seed_data.up.sql
│   ├── 000002_seed_data.down.sql
│   └── ...
└── postgres/
    ├── 000001_init_schema.up.sql
    ├── 000001_init_schema.down.sql
    └── ...
```

每个版本包含 `.up.sql`（执行）和 `.down.sql`（回滚）。

### 迁移版本

| 版本 | 内容 |
|------|------|
| 000001 | 基础表结构（用户、角色、菜单、配置、日志等） |
| 000002 | 种子数据（super_admin 角色、系统菜单、权限策略） |
| 000003 | 企业基础（部门、岗位、数据权限表结构） |
| 000004 | 字典表结构 |
| 000005 | 字典种子数据 |
| 000006 | 附件表结构 |
| 000007 | 附件种子数据 |
| 000008 | 组织菜单种子数据 |
| 000009 | 菜单图标对齐 |

### 迁移执行

迁移在服务启动时自动执行：

```
bootstrap.MustRun()
  → 连接数据库
  → migrate.MustRun(migrationsFS, driver)
  → 已执行的版本自动跳过
```

也提供了验证脚本 `scripts/verify-realdb-migrations.sh`，可针对真实数据库测试迁移。

## 数据库表清单

| 表 | 用途 | 逻辑删除 |
|----|------|---------|
| `sys_user` | 用户 | ✅ |
| `sys_role` | 角色 | ✅ |
| `sys_menu` | 菜单/权限 | ✅ |
| `sys_department` | 部门 | ✅ |
| `sys_post` | 岗位 | ✅ |
| `sys_user_role` | 用户-角色关联 | ❌ |
| `sys_role_menu` | 角色-菜单关联 | ❌ |
| `sys_user_post` | 用户-岗位关联 | ❌ |
| `sys_role_data_scope` | 角色自定义数据范围 | ❌ |
| `sys_config` | 系统配置 | ✅ |
| `sys_dict_type` | 字典类型 | ✅ |
| `sys_dict_item` | 字典项 | ✅ |
| `sys_login_log` | 登录日志 | ❌ |
| `sys_operation_log` | 操作日志 | ❌ |
| `sys_notice` | 通知公告 | ✅ |
| `sys_file` | 文件 | ✅ |
| `sys_attachment` | 附件 | ✅ |
| `casbin_rule` | Casbin 策略 | ❌ |

## 种子数据

### 超级管理员角色

- 角色编码：`super_admin`
- 数据范围：`all`（所有数据）
- 关联全部菜单和按钮权限

### 系统菜单结构

```
系统管理
├── 用户管理 (system:user:*)
├── 角色管理 (system:role:*)
├── 菜单管理 (system:menu:*)
├── 部门管理 (system:dept:*)
└── 岗位管理 (system:post:*)
系统工具
├── 系统配置 (system:config:*)
├── 数据字典 (system:dict:*)
├── 文件管理 (system:file:*)
├── 操作日志 (system:operationlog:*)
├── 登录日志 (system:loginlog:*)
└── 通知公告 (system:notice:*)
```

### Casbin 策略

所有系统管理接口的权限策略分配给 `super_admin` 角色。

## 系统初始化

管理员账号通过 Setup API 创建，不在迁移中硬编码：

```bash
curl -X POST http://localhost:8080/api/v1/setup/init \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"YourPassword123","nickname":"管理员"}'
```

初始化只能执行一次，重复调用会返回错误。

## 新增迁移

添加新迁移文件时：

1. 在 `migrations/mysql/` 和 `migrations/postgres/` 同时添加
2. 编号递增，格式 `000010_xxx.up.sql` / `000010_xxx.down.sql`
3. 同时提供 `.up.sql` 和 `.down.sql`
4. `.up.sql` 包含建表、索引、种子数据
5. `.down.sql` 包含反向操作（DROP TABLE、DELETE）
