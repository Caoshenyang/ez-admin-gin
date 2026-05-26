# Contributing to EZ Admin Gin

感谢你对 EZ Admin Gin 的关注。

本项目主要是维护者自用的后台系统基座，公开源码供参考和复用。欢迎通过 [GitHub Issues](https://github.com/Caoshenyang/ez-admin-gin/issues) 反馈 Bug 或建议，但项目不以社区协作为主要目标，Pull Request 不保证接受或处理。当前优先级是保持基座稳定，并支撑维护者自己的 MVP 项目快速开发。

## 当前协作方式

- Issue：欢迎报告 Bug、提出建议或分享使用心得。
- Pull Request：不保证接受或处理；如确需协作，建议先通过 Issue 对齐范围。
- 文档和代码以维护者自己的使用节奏为主，不承诺社区路线图。

## 当前不在范围内

- 完整自动化测试覆盖率
- 大型测试框架
- 微服务重构
- 低代码 / 无代码引擎
- 复杂多租户隔离
- 大型 IAM 平台能力
- 社区治理和长期 PR 协作

## 本地开发

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

## 轻量验证

```bash
make verify
```

`make verify` 只做轻量质量检查：后端 vet、后端 build、前端类型检查、前端 lint、前端 build、Docker Compose 配置校验。发布前请结合 [MANUAL_TEST.md](MANUAL_TEST.md) 做人工冒烟验证。

## 代码风格

- 后端遵循标准 Go 规范，保持 `gofmt`、`go vet` 通过。
- 前端使用 Vue 3、`<script setup>`、TypeScript 和项目已有 Naive UI / Tailwind 写法。
- API 类型通过 `make generate-types` 生成，不手动编辑 `admin/src/api/generated.ts`。
- 文档以中文为主，保持面向个人项目快速复用的定位。

## Commit 规范

推荐使用 [Conventional Commits](https://www.conventionalcommits.org/)：

```text
feat: add user import from CSV
fix: resolve menu tree rendering for empty children
docs: update deployment guide with HTTPS setup
refactor: extract pagination helper from list handlers
chore: update Go dependencies
```
