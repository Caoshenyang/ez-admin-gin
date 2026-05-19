---
title: 迁移与种子数据
description: 完整版 SQL、历史迁移链与管理员初始化之间的关系说明
---

# 迁移与种子数据

这页只回答三个问题：

1. 生产环境到底以哪份 SQL 为准
2. 仓库里保留的增量迁移链还扮演什么角色
3. 第一个可登录管理员账号是怎么创建出来的

::: tip 先记结论
对外交付时，数据库以 `server/migrations/mysql/full_schema_and_seed.sql` 和 `server/migrations/postgres/full_schema_and_seed.sql` 为准。

现有 `000001` 到 `000011` 的增量迁移链继续保留，主要用于程序启动兼容和历史演进追踪；它不是新的对外交付入口。
:::

## 当前初始化分成两段

当前主线不是“一个 SQL 直接连管理员账号都帮你建好”，而是固定分成两段：

1. 导入完整版 SQL，创建全部系统表和内置系统种子
2. 调用 `POST /api/v1/setup/init` 创建第一个管理员用户

也就是说：

- 角色、菜单、按钮、Casbin 策略、字典、附件、通知等内置数据来自数据库脚本
- 第一个真正能登录后台的管理员账号来自初始化接口

## 推荐使用方式

### 生产环境 / 新环境初始化

按数据库类型选择下面其中一份：

- MySQL：`server/migrations/mysql/full_schema_and_seed.sql`
- PostgreSQL：`server/migrations/postgres/full_schema_and_seed.sql`

导入完成后，再调用：

```bash
curl -X POST http://localhost:8080/api/v1/setup/init \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"YourPassword123","nickname":"管理员"}'
```

### 程序启动时的增量迁移

服务启动时仍会执行 `server/migrations/{driver}/` 下的增量迁移链。

这条链路继续保留的原因是：

- 兼容已有开发环境和历史部署方式
- 让仓库保留清晰的结构演进记录
- 避免一次性推翻当前启动逻辑带来的额外风险

但从交付口径上，它现在是“内部兼容机制”，不是新的用户入口。

::: warning 不要把两种入口重复执行成两遍初始化
如果你已经从空库导入了 `full_schema_and_seed.sql`，就不要再把它理解成“还需要额外手动补一遍 000001~000011”。

程序启动后看到迁移链是 `no change` 或幂等跳过，属于正常现象。
:::

## 完整版 SQL 是怎么来的

两份完整版 SQL 不是手工拼出来的，它们由仓库脚本从真实的 `.up.sql` 迁移链顺序归并生成：

```bash
./scripts/build-full-migrations.sh
```

如果后续迁移链继续演进，更新完整版 SQL 时也应该走这个脚本，而不是直接手改产物文件。

## 当前迁移版本

当前保留的增量迁移版本为：

| 版本 | 内容 |
| --- | --- |
| `000001` | 基础系统表结构 |
| `000002` | 超级管理员角色、系统菜单、Casbin、角色菜单种子 |
| `000003` | 组织体系、岗位、数据权限结构 |
| `000004` | 字典表结构 |
| `000005` | 字典种子数据 |
| `000006` | 附件表结构 |
| `000007` | 附件菜单与权限种子 |
| `000008` | 部门 / 岗位菜单与权限种子 |
| `000009` | 菜单图标对齐 |
| `000010` | 通知表结构 |
| `000011` | 通知权限种子 |

## 验证你是否初始化成功

完成完整版 SQL 导入后，系统应至少满足下面这些条件：

- 能查到 `super_admin` 角色
- 能查到系统管理主菜单及其按钮节点
- `casbin_rule` 中存在 `super_admin` 的系统接口策略
- 字典、附件、通知等扩展模块表已经存在

完成 `setup/init` 后，还应看到：

- `sys_user` 出现首条管理员用户
- `sys_user_role` 出现该用户与 `super_admin` 的绑定
- 重复调用 `setup/init` 会返回冲突，而不是重复创建账号

## 相关入口

- 完整建表与种子说明：[/reference/database-ddl](/reference/database-ddl)
- 初始化链路说明：[/reference/init-data-reference](/reference/init-data-reference)
