# Contributing to EZ Admin Gin

感谢你对 EZ Admin Gin 的关注！

## 当前贡献状态

- **Issue**：欢迎通过 [GitHub Issues](https://github.com/Caoshenyang/ez-admin-gin/issues) 报告 Bug、提出建议或分享使用心得。
- **Pull Request**：**暂不接受**。项目仍处于架构快速迭代阶段，代码结构和接口设计调整较频繁，外部 PR 容易产生冲突。等架构稳定后会开放 PR。

开放 PR 的条件：
- 核心架构（platform 层、模块接入流程）稳定
- API 接口版本化管理就绪
- 贡献指南和 PR 模板完善

届时会在 README 和本文件中更新。

## 开发环境

| 工具 | 版本 | 用途 |
|------|------|------|
| Go | 1.26+ | 后端 |
| Node.js | 20.19+ 或 22.12+ | 前端 & 文档 |
| pnpm | 9+ | 包管理器 |
| Docker | 20+ | 本地 PostgreSQL + Redis |
| make | Any | 构建自动化（Windows 可选） |

## 本地启动

```bash
# 1. 启动 PostgreSQL 和 Redis
make docker-up

# 2. 启动后端（另一个终端）
make server-dev

# 3. 初始化管理员账号
curl -X POST http://localhost:8080/api/v1/setup/init \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"YourPassword123","nickname":"管理员"}'

# 4. 启动前端（另一个终端）
make install && make admin-dev
```

## 代码风格

### 后端（Go）

- 遵循标准 Go 规范（`gofmt`、`go vet`）
- 包命名：小写，无下划线
- 错误处理：使用 `errorsx` 包统一错误码
- 测试放在 `server/tests/`，不在业务代码目录中创建 `*_test.go`

### 前端（Vue 3 + TypeScript）

- 使用 `<script setup>` + Composition API
- TypeScript strict 模式
- 页面：`admin/src/modules/{module}/pages/{Name}View.vue`
- API 类型：通过 `make generate-types` 生成，不手动编辑 `admin/src/api/generated.ts`

### 文档

- 文档在 `docs/`（VitePress）
- 中文为主
- 提交前运行 `make docs-build` 验证构建

## 质量检查

```bash
make lint              # 后端 vet + 前端类型检查 + lint
make test              # 后端测试
make test-contract     # OpenAPI 契约测试（不需要 DB）
make test-integration  # 集成测试（需要 DB + Redis）
make build             # 后端二进制 + 前端构建
make docs-build        # 文档站构建
```

## Commit 规范

使用 [Conventional Commits](https://www.conventionalcommits.org/)：

```
feat: add user import from CSV
fix: resolve menu tree rendering for empty children
docs: update deployment guide with HTTPS setup
test: add role permission assignment test
refactor: extract pagination helper from list handlers
chore: update Go dependencies
```

## 关键文件

在做任何修改前，建议先阅读：

- `CLAUDE.md` — 项目规则和约束
- `docs/project/PHASE_STATUS.md` — 当前进度和状态
- `docs/project/QUALITY_ROADMAP.md` — 长期路线图
- `docs/project/TESTING_STRATEGY.md` — 测试规范
- `docs/project/DECISION_LOG.md` — 架构决策记录

## 问题？

欢迎在 [GitHub Issues](https://github.com/Caoshenyang/ez-admin-gin/issues) 提问，使用 `question` 标签。
