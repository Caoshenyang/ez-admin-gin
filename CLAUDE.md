# EZ Admin Gin — Project Rules

## Skills

本项目在 `.agents/skills/` 下维护了自定义 skill，遇到以下场景时必须先读取对应 SKILL.md 再工作：

- 数据库表结构设计、GORM 模型相关 → `.agents/skills/database-schema-design/SKILL.md`
- 编写或修改 docs/ 下的 Markdown/VitePress 文档 → `.agents/skills/vitepress-doc-writing/SKILL.md`
- 编写或修改 admin 前端代码 → `.agents/skills/ez-admin-frontend-guidelines/SKILL.md`
- 新增业务 CRUD 模块 → `.agents/skills/module-generator/SKILL.md`

## 项目目标

- 当前阶段：稳定收尾与交付面整理
- 核心方向：生产可交付的后台底座，不再继续扩展示例功能
- 判断标准：减少双套真相、降低维护成本、优先修复漂移

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
- `server/tests/` — 集中式测试目录
  - `server/tests/api/` — API 黑盒测试
  - `server/tests/rbac/` — 权限和数据权限流程测试
  - `server/tests/contract/` — OpenAPI 契约测试
  - `server/tests/testutil/` — 测试辅助工具

## 测试组织原则

- 不在业务代码目录中散落大量 *_test.go
- 后端测试集中在 server/tests
- server/tests/api：API 黑盒测试
- server/tests/rbac：权限和数据权限流程测试
- server/tests/contract：OpenAPI 契约测试
- server/tests/testutil：测试辅助工具
- tests/e2e 或 admin/e2e：Playwright E2E（Phase 4）

## 禁止事项

- 禁止大面积创建 co-located unit tests
- 禁止大规模重构（除非有明确收益且先说明风险）
- 禁止绕过权限逻辑
- 禁止为了测试通过而降低安全性
- 禁止连接生产数据库
- 禁止破坏现有 CI
- 禁止在没有说明的情况下扩大任务范围
- 禁止把 TODO 当成完成
- 禁止把 t.Skip 当成有效覆盖

## 每次任务前必须阅读

- CLAUDE.md（本文件）
- README.md
- docs/backend/migration.md
- docs/reference/database-ddl.md
- docs/reference/init-data-reference.md

## 数据库交付约定

- 数据库对外交付以 `server/migrations/mysql/full_schema_and_seed.sql` 和 `server/migrations/postgres/full_schema_and_seed.sql` 为准。
- 这两份文件由 `./scripts/build-full-migrations.sh` 从真实 `.up.sql` 迁移链归并生成。
- 不要手工维护第二套“示意 SQL”真相源。
- 首个管理员继续通过 `POST /api/v1/setup/init` 创建，不写死到种子里。
