---
title: 初始化数据参考
description: "集中说明 EZ Admin Gin 当前首批角色、菜单、Casbin 策略、管理员初始化接口，以及首次部署时初始化链路的真实执行方式。"
---

# 初始化数据参考

这页专门回答一个高频问题：

> 当前仓库第一次跑起来后，系统里最初那批角色、菜单、权限和管理员账号到底是怎么来的？

## 先记住当前初始化分成两段

当前主线不是“一个 SQL 包打天下”，而是分成两段：

1. 迁移脚本先落系统表和内置种子数据
2. `/api/v1/setup/init` 再创建第一个管理员用户

也就是说：

- 角色、菜单、按钮、Casbin 策略来自迁移种子
- 第一个真正能登录的后台账号来自初始化接口

## 当前初始化链路

当前真实顺序是：

```text
数据库迁移
  ↓
000001_init_schema
  ↓
000002_seed_data
  ↓
系统具备 super_admin 角色、菜单树、Casbin 策略
  ↓
POST /api/v1/setup/init
  ↓
创建第一个管理员用户，并绑定 super_admin
```

相关文件入口：

- `server/migrations/postgres/000001_init_schema.up.sql`
- `server/migrations/postgres/000002_seed_data.up.sql`
- `server/migrations/mysql/000001_init_schema.up.sql`
- `server/migrations/mysql/000002_seed_data.up.sql`
- `server/internal/module/setup/handler.go`
- `scripts/setup-server.sh`

## 第一段：迁移种子到底写了什么

当前 `000002_seed_data.*` 主要写了四类内容：

| 内容 | 落表 | 作用 |
| --- | --- | --- |
| 超级管理员角色 | `sys_role` | 提供固定的 `super_admin` 主体 |
| 系统菜单与按钮 | `sys_menu` | 提供前端菜单树和按钮权限节点 |
| Casbin 接口策略 | `casbin_rule` | 提供 `super_admin` 对各系统接口的访问能力 |
| 角色菜单关系 | `sys_role_menu` | 让 `super_admin` 能拿到全部系统菜单 |

### 当前内置角色

当前种子只内置了一条系统角色：

| 字段 | 值 |
| --- | --- |
| `id` | `1` |
| `code` | `super_admin` |
| `name` | `超级管理员` |
| `remark` | `系统内置角色` |

这意味着后续很多链路都默认依赖它：

- 初始化管理员默认绑定到角色 `1`
- `LoadActor` 会按 `role.code == super_admin` 判断是否超级管理员
- Casbin 默认主体也是 `super_admin`

::: warning 超级管理员的角色编码和 ID 都不适合随意改
当前主线里，`super_admin` 既是角色语义，也是接口权限和初始化流程的一条稳定锚点。

如果你改了这条记录的编码、ID 或初始化假设，记得连带检查：

- `setup/init`
- `LoadActor`
- Casbin 种子
- 角色菜单种子
:::

### 当前内置菜单树

当前种子已经给出了一整套系统管理目录、菜单和按钮。

大体结构是：

| 层级 | 说明 |
| --- | --- |
| 目录 | 例如 `system` |
| 菜单 | 例如 `system:user`、`system:role` |
| 按钮 | 例如 `system:user:create`、`system:notice:update` |

当前 PostgreSQL 脚本里还明确写了 ID 规划规则：

- 目录从 `100` 起
- 菜单从 `200` 起
- 按钮从 `1000` 起

这不是为了“数字好看”，而是为了：

- 让跨环境数据保持稳定
- 让后续扩展有固定间距
- 让初始化数据和调试日志更容易对照

### 当前内置 Casbin 策略

当前种子把 `super_admin` 需要的系统接口权限一次性写进：

- `casbin_rule`

典型记录形态是：

```sql
INSERT INTO casbin_rule (ptype, v0, v1, v2)
VALUES ('p', 'super_admin', '/api/v1/system/users', 'GET');
```

可以直接理解为：

| 字段 | 含义 |
| --- | --- |
| `ptype` | 策略类型，当前主要是 `p` |
| `v0` | 角色主体 |
| `v1` | 接口路径 |
| `v2` | HTTP 方法 |

## 第二段：为什么还需要 `/setup/init`

因为迁移种子只创建了：

- 角色
- 菜单
- Casbin 策略

但并没有创建真正的后台用户。

所以系统首次可登录，还差这最后一步：

- 创建管理员用户
- 绑定到 `super_admin`

这一步当前通过下面接口完成：

- `POST /api/v1/setup/init`

## `setup/init` 当前做了什么

当前实现位于：

- `server/internal/module/setup/handler.go`

它会做下面几件事：

1. 检查 `sys_user` 是否已经有记录
2. 如果已经有用户，则拒绝重复初始化
3. 校验 `username / password / nickname`
4. 用 `bcrypt` 生成密码哈希
5. 创建第一条管理员用户记录
6. 创建 `sys_user_role(user_id, role_id=1)` 绑定

也就是说，当前接口默认有一个非常明确的前提：

- 角色 `1` 必须就是 `super_admin`

## 当前首次部署脚本怎么调用它

当前服务器首次部署脚本：

- `scripts/setup-server.sh`

会在基础设施和后端服务启动后，自动执行：

```bash
curl -X POST http://localhost/api/v1/setup/init \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"Admin@123456","nickname":"管理员"}'
```

并按返回码做分支：

| HTTP 状态码 | 当前解释 |
| --- | --- |
| `200` | 管理员创建成功 |
| `409` | 系统已初始化，跳过 |
| 其他 | 需要手工排查和重试 |

## PostgreSQL 和 MySQL 初始化种子的差异

两套脚本的核心语义保持一致，但幂等写法不同：

| 数据库 | 幂等方式 |
| --- | --- |
| PostgreSQL | `ON CONFLICT DO NOTHING` |
| MySQL | `INSERT IGNORE` |

另外时间函数也略有差异：

| 数据库 | 时间函数 |
| --- | --- |
| PostgreSQL | `NOW()` |
| MySQL | `NOW(3)` |

## 当前最值得记住的几条初始化约定

### 1. 先有角色和菜单，再有管理员用户

当前管理员账号不是直接把权限写到用户身上，而是：

- 先绑定 `super_admin`
- 再通过角色继承菜单和 Casbin 能力

### 2. 菜单树和接口权限是两张表

初始化数据里同时写了：

- `sys_menu`
- `casbin_rule`

这说明：

- 前端菜单和按钮显隐依赖 `sys_menu`
- 后端接口放行依赖 `casbin_rule`

二者是配套关系，不是同一件事。

### 3. `setup/init` 只适合首个管理员

它的职责不是“通用建用户接口”，而是：

- 一次性创建第一位管理员

一旦系统里已经有用户，它就应该退出。

## 扩展初始化数据时最稳的做法

如果你后续要扩一批新模块的初始化数据，当前更稳的顺序是：

1. 先补迁移脚本里的 `sys_menu`
2. 再补对应 `casbin_rule`
3. 再补 `sys_role_menu`
4. 如果需要内置角色，再补 `sys_role`
5. 不要把真实业务用户直接写死进种子，优先通过后台接口或单独初始化脚本创建

## 最常见的初始化问题

### 菜单有了，但接口 403

通常说明：

- `sys_menu` 补了
- 但 `casbin_rule` 没同步补

### 接口能通，但前端没有入口

通常说明：

- Casbin 策略补了
- 但 `sys_menu` 或 `sys_role_menu` 没同步补

### `setup/init` 返回冲突

通常说明：

- `sys_user` 已经存在记录
- 当前系统已完成首次管理员初始化

## 相关教程与参考页

- [权限码约定](./permission-code-conventions)
- [动态菜单组件白名单](./dynamic-menu-component-reference)
- [数据库建表语句](./database-ddl)
- [第 8 章：部署与项目复用](../tutorial/chapter-8/)
