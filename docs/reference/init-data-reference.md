---
title: 初始化数据参考
description: "说明内置系统种子与首个管理员初始化的真实边界。"
---

# 初始化数据参考

这页只回答一个高频问题：

> 第一次把系统跑起来时，哪些数据来自数据库脚本，哪些数据来自初始化接口？

## 当前初始化边界

当前主线固定分成两段：

1. 完整版 SQL 创建系统表和内置系统种子
2. `POST /api/v1/setup/init` 创建第一个管理员用户

也就是说：

- 内置角色、菜单、按钮、Casbin 策略、角色菜单、字典、附件、通知权限种子来自数据库脚本
- 首个可登录管理员账号不内置在种子里，而是运行时创建

## 第一段：完整版 SQL 会写入什么

完整版 SQL 会落下面这些稳定锚点：

- `super_admin` 角色
- 系统管理菜单树和按钮权限节点
- `super_admin` 对系统接口的 Casbin 策略
- 角色菜单关系
- 字典、附件、部门、岗位、通知等主线模块的结构和内置系统数据

相关文件：

- `server/migrations/mysql/full_schema_and_seed.sql`
- `server/migrations/postgres/full_schema_and_seed.sql`

::: warning 不要把管理员账号写死到种子里
当前项目明确把“系统骨架”和“首个运营账号”分开处理。

这样做的好处是：

- 避免仓库内置固定密码
- 避免不同环境共享同一套管理员账号
- 保持首次初始化更接近真实交付场景
:::

## 第二段：`setup/init` 会做什么

`POST /api/v1/setup/init` 负责：

1. 检查系统里是否已经存在用户
2. 校验 `username / password / nickname`
3. 创建第一条管理员用户
4. 自动绑定 `super_admin` 角色

当前实现依赖一个稳定前提：

- 数据库里必须已经存在角色编码 `super_admin`

因此正确顺序一定是：

```text
导入完整版 SQL
  ↓
系统具备 super_admin、菜单树、Casbin 策略
  ↓
POST /api/v1/setup/init
  ↓
创建首个管理员用户并绑定 super_admin
```

## 首次初始化命令

```bash
curl -X POST http://localhost:8080/api/v1/setup/init \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"YourPassword123","nickname":"管理员"}'
```

执行成功后，你应该看到：

- `sys_user` 中出现首条管理员用户
- `sys_user_role` 中出现该用户与 `super_admin` 的绑定
- 再次执行会返回“系统已初始化”

## 当前最值得记住的约定

- 数据库脚本提供系统骨架，不提供真实管理员账号
- `super_admin` 是当前初始化链路的稳定角色锚点
- 菜单显隐和接口放行是两套配合关系：
  - `sys_menu` 负责前端菜单和按钮节点
  - `casbin_rule` 负责后端接口授权
- 完整版 SQL 是交付入口，增量迁移链是兼容和演进记录

## 相关页面

- [迁移与种子数据](/backend/migration)
- [数据库建表语句](/reference/database-ddl)
