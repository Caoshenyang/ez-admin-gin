# EZ Admin Gin

![EZ Admin Gin Badge](brand-assets/svg/readme-badge.svg)

EZ Admin Gin 是一个维护者自用优先的 Go + Gin + Vue 3 后台管理系统基座。

本仓库公开源码，主要用于个人项目、中小型后台系统和 SaaS/MVP 原型快速开发。项目不以社区协作为主要目标，也不追求完整自动化测试覆盖率。当前重点是保持基座稳定、结构清晰、部署简单，方便维护者基于它快速创建业务项目。

## 适合 / 不适合

| 适合 | 不适合 |
|------|--------|
| 个人项目后台 | 大型企业 IAM / 统一身份认证平台 |
| SaaS 原型 / MVP | 微服务架构的服务治理平台 |
| 中小型管理系统（ERP、CRM、CMS 底座） | 低代码 / 无代码平台底座 |
| 需要权限、数据权限、动态菜单的后台底座 | 无需二次开发的商业成品系统 |
| 想学习 Go + Vue 全栈工程化的开发者 | 高并发（万级 QPS+）独立场景 |

## 核心特性

- **权限体系** — JWT 认证、Casbin 接口授权、五级数据权限、动态菜单、按钮权限
- **组织体系** — 部门树、岗位管理、用户-角色-岗位多对多关联
- **模块化扩展** — 统一四层结构（handler/service/repository/domain），标准接入流程
- **部署简单** — Docker Compose、Nginx HTTP/HTTPS、服务器二进制部署
- **系统工具齐全** — 数据字典、系统配置、文件上传、操作日志、登录日志、通知公告

## 技术栈

| 层 | 技术 |
|----|------|
| 后端框架 | Go 1.26 + Gin |
| ORM | GORM（PostgreSQL / MySQL） |
| 权限引擎 | Casbin |
| 缓存 | Redis |
| 前端框架 | Vue 3 + TypeScript |
| UI 组件库 | Naive UI |
| 状态管理 | Pinia |
| CSS | Tailwind CSS 4 |
| 构建工具 | Vite |
| 文档 | VitePress |

## 项目结构

```text
ez-admin-gin/
├── server/                 Go 后端
│   ├── configs/            配置文件
│   ├── internal/           业务代码
│   │   ├── bootstrap/      启动引导
│   │   ├── modules/        业务模块（auth / iam / system / setup）
│   │   ├── platform/       平台层（authn / authz / datascope / middleware / ...）
│   │   └── pkg/            公共工具包
│   ├── migrations/         数据库迁移（MySQL + PostgreSQL）
│   └── docs/               Swagger / OpenAPI 规范
├── admin/                  Vue 3 管理台
│   └── src/
│       ├── modules/        业务模块（auth / iam / system）
│       ├── layouts/        布局组件
│       ├── router/         路由 + 动态菜单注册
│       └── stores/         Pinia 状态管理
├── docs/                   VitePress 文档站
├── deploy/                 Docker Compose、Nginx、部署配置
├── scripts/                部署与打包脚本
└── MANUAL_TEST.md          发布前人工测试清单
```

## 快速启动

### 环境要求

| 工具 | 最低版本 | 用途 |
|------|---------|------|
| Go | 1.26+ | 后端 |
| Node.js | 20.19+ 或 22.12+ | 前端 & 文档 |
| pnpm | 9+ | 包管理器 |
| Docker | 20+ | 本地 PostgreSQL + Redis |
| make | GNU Make 4+ | 构建自动化（Windows 可选） |

### 使用 Makefile（推荐）

```bash
# 查看所有可用命令
make help

# 1. 启动 PostgreSQL 和 Redis
make docker-up

# 2. 如需更新完整版初始化 SQL（可选）
make db-full-sql

# 3. 启动后端（另一个终端）
make server-dev

# 4. 初始化管理员账号
curl -X POST http://localhost:8080/api/v1/setup/init \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"YourPassword123","nickname":"管理员"}'

# 5. 启动前端（另一个终端）
make install && make admin-dev
```

初始化时通过 `setup/init` 接口自行指定用户名和密码。生产环境首次登录后请立即修改密码。

## 常用命令速查

| 命令 | 说明 |
|------|------|
| `make help` | 显示所有可用命令 |
| `make install` | 安装前端依赖 |
| `make server-dev` | 启动后端 (`go run .`) |
| `make admin-dev` | 启动前端 (`pnpm dev`) |
| `make server-vet` | 后端代码检查 (`go vet ./...`) |
| `make admin-check` | 前端类型检查 + lint |
| `make lint` | 后端 vet + 前端检查 + API 类型同步检查 |
| `make build` | 构建后端二进制 + 前端产物 |
| `make verify` | 轻量验证：vet、前端检查、构建、Docker Compose 配置校验 |
| `make docker-up` | 启动 PostgreSQL + Redis |
| `make docker-down` | 停止 PostgreSQL + Redis |
| `make docker-config` | 验证所有 Docker Compose 配置 |
| `make generate-types` | 从 Swagger spec 生成前端 TypeScript 类型 |
| `make check-types` | 检查生成的 API 类型是否与 Swagger spec 同步 |

## 质量策略

本项目不维护复杂自动化测试体系，不追求测试覆盖率。当前采用轻量质量策略：后端 `go vet`、后端 `go build`、前端 TypeScript 类型检查、前端 lint、前端生产构建、Docker Compose 配置校验，以及发布前人工冒烟测试。

本项目不维护 API 自动化测试、RBAC 自动化测试、Contract 测试、E2E 测试和覆盖率报告。

本地轻量验证：

```bash
make verify
```

需要更细地拆开看时，可以分别执行：

```bash
cd server && go mod tidy && go vet ./... && go build ./...
cd admin && pnpm install --frozen-lockfile && pnpm type-check && pnpm lint && pnpm build
docker compose -f deploy/compose.local.yml config --quiet
EZ_AUTH_JWT_SECRET=placeholder docker compose -f deploy/compose.prod.yml config --quiet
docker compose -f deploy/compose.server.yml config --quiet
```

发布或复制到新 MVP 项目前，按 [MANUAL_TEST.md](MANUAL_TEST.md) 做人工验证。

## 文档

在线文档：[https://caoshenyang.github.io/ez-admin-gin/](https://caoshenyang.github.io/ez-admin-gin/)

- [快速开始](https://caoshenyang.github.io/ez-admin-gin/getting-started/)
- [系统架构](https://caoshenyang.github.io/ez-admin-gin/architecture/overview)
- [权限体系](https://caoshenyang.github.io/ez-admin-gin/architecture/rbac)
- [后端开发](https://caoshenyang.github.io/ez-admin-gin/backend/overview)
- [前端开发](https://caoshenyang.github.io/ez-admin-gin/frontend/overview)
- [部署概览](https://caoshenyang.github.io/ez-admin-gin/deployment/overview)
- [生产环境检查清单](https://caoshenyang.github.io/ez-admin-gin/deployment/production-checklist)
- [参考手册](https://caoshenyang.github.io/ez-admin-gin/reference/)

## Roadmap

已完成：

- [x] JWT 认证 + Casbin 权限控制
- [x] 动态菜单与按钮权限
- [x] 组织体系（部门/岗位）
- [x] 五级数据权限
- [x] 系统模块（用户/角色/菜单/配置/字典/文件/日志/公告）
- [x] 前端管理台（登录/壳子/动态菜单/管理页面）
- [x] 多场景部署方案（Docker Compose + Nginx + 部署脚本）
- [x] WebSocket 通知公告实时推送

暂不计划：

- [ ] 微服务拆分
- [ ] 低代码 / 无代码引擎
- [ ] 复杂多租户隔离
- [ ] 大型 IAM 平台能力
- [ ] 社区治理和长期 PR 协作流程

## Contributing

本项目主要是维护者自用的后台系统基座，公开源码供参考和复用。欢迎通过 Issue 反馈 Bug 或建议，但项目不以社区协作为主要目标，Pull Request 不保证接受或处理。当前优先级是保持基座稳定，并支撑维护者自己的 MVP 项目快速开发。

详见 [CONTRIBUTING.md](CONTRIBUTING.md)。

## License

[MIT](LICENSE)
