# EZ Admin Gin — Project Rules

## Skills

本项目在 `.agents/skills/` 下维护了自定义 skill，遇到以下场景时必须先读取对应 SKILL.md 再工作：

- 数据库表结构设计、GORM 模型相关 → `.agents/skills/database-schema-design/SKILL.md`
- 编写或修改 docs/ 下的 Markdown/VitePress 文档 → `.agents/skills/vitepress-doc-writing/SKILL.md`
- 编写或修改 admin 前端代码 → `.agents/skills/ez-admin-frontend-guidelines/SKILL.md`
- 新增业务 CRUD 模块 → `.agents/skills/module-generator/SKILL.md`

## 项目目标

- 当前定位：维护者自用优先的后台系统基座
- 核心方向：快速支撑维护者自己的 SaaS/MVP 项目开发
- 公开源码目的：供他人参考和复用，不以社区协作为主要目标
- 判断标准：结构清晰、部署简单、维护成本低、方便复制到真实业务项目

## 技术栈

- Go + Gin + GORM + Casbin + Redis
- Vue 3 + TypeScript + Naive UI + Pinia
- Docker Compose
- VitePress docs
- GitHub Actions CI

## 项目结构

- `server/` — Go 后端
- `admin/` — Vue 3 前端
- `docs/` — 文档站
- `deploy/` — 部署配置
- `scripts/` — 部署、打包和迁移辅助脚本
- `MANUAL_TEST.md` — 发布前人工测试清单

## 质量策略

本项目不维护复杂自动化测试体系，不追求测试覆盖率。修改后优先运行轻量验证：

- 后端 vet
- 后端 build
- 前端 type-check
- 前端 lint
- 前端 build
- Docker Compose config 校验（修改部署文件时必须跑）

除非维护者明确要求，不要新增复杂自动化测试。

## 禁止事项

- 禁止大规模重构（除非有明确收益且先说明风险）
- 禁止绕过权限逻辑
- 禁止连接生产数据库
- 禁止破坏构建、启动和部署链路
- 禁止在没有说明的情况下扩大任务范围
- 禁止把 TODO 当成完成
- 禁止为了质量检查通过而降低安全性

## 每次任务前建议阅读

- `CLAUDE.md`（本文件）
- `README.md`
- `docs/backend/migration.md`
- `docs/reference/database-ddl.md`
- `docs/reference/init-data-reference.md`

## 数据库交付约定

- 数据库对外交付以 `server/migrations/mysql/full_schema_and_seed.sql` 和 `server/migrations/postgres/full_schema_and_seed.sql` 为准。
- 这两份文件由 `./scripts/build-full-migrations.sh` 从真实 `.up.sql` 迁移链归并生成。
- 不要手工维护第二套“示意 SQL”真相源。
- 首个管理员继续通过 `POST /api/v1/setup/init` 创建，不写死到种子里。
