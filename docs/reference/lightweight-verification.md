---
title: 轻量验证与人工测试
description: "说明 EZ Admin Gin 当前采用的轻量质量策略、验证命令和发布前人工测试清单。"
---

# 轻量验证与人工测试

EZ Admin Gin 当前是维护者自用优先的后台系统基座，不维护复杂自动化测试体系，也不追求测试覆盖率。质量策略收敛为轻量验证命令 + 发布前人工冒烟检查。

::: tip 🎯 这页解决什么
当你修改代码、准备发布，或复制成本地 MVP 项目时，用这页确认应该跑哪些检查、哪些功能需要手动点一遍。
:::

## 轻量验证命令

在仓库根目录执行：

```bash
make verify
```

`make verify` 会执行：

| 检查 | 命令 |
| --- | --- |
| 后端静态检查 | `go vet ./...` |
| 后端构建 | `go build ./...` |
| 前端类型检查 | `pnpm type-check` |
| 前端 lint | `pnpm lint` |
| 前端生产构建 | `pnpm build` |
| Docker Compose 配置校验 | `docker compose ... config --quiet` |

## 拆开执行

如果你想定位是哪一段失败，可以拆开跑：

```bash
cd server
go mod tidy
go vet ./...
go build ./...
```

```bash
cd admin
pnpm install --frozen-lockfile
pnpm type-check
pnpm lint
pnpm build
```

```bash
docker compose -f deploy/compose.local.yml config --quiet
EZ_AUTH_JWT_SECRET=placeholder docker compose -f deploy/compose.prod.yml config --quiet
docker compose -f deploy/compose.server.yml config --quiet
```

## 人工测试清单

仓库根目录的 `MANUAL_TEST.md` 是发布前人工测试清单，覆盖：

- 本地启动、管理员初始化、登录/登出
- 未登录接口保护、动态菜单、用户和角色管理
- 权限分配、无权限拦截、部门/岗位、数据权限
- 字典、配置、日志、文件上传等系统模块
- 后端构建、前端构建、Docker Compose 配置校验
- MVP 项目复用前检查

::: warning ⚠️ 发布前不要只看构建结果
轻量验证只能证明代码能通过构建和基础静态检查，不能替代真实页面和权限链路的人工确认。涉及权限、菜单、数据范围或部署配置的改动，发布前至少按 `MANUAL_TEST.md` 做一轮冒烟验证。
:::
