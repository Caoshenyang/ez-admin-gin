# EZ Admin Gin

![EZ Admin Gin Badge](brand-assets/svg/readme-badge.svg)

EZ Admin Gin 是一个维护者自用优先的后台系统基座，面向个人项目、中小型管理系统和 SaaS/MVP 原型快速开发。

> 当前稳定版本：**v1.2.0**  
> 当前定位：公开源码供参考和复用，但不以社区协作型框架为主要目标。

它提供一套已经能落地使用的后台基础能力：认证、权限、组织、动态菜单、系统配置、附件、日志、通知、邮件、消息提醒、健康检查和部署文档。项目更重视结构清晰、可复制、可上线，而不是把自己包装成大而全的通用平台。

## 适合 / 不适合

| 适合 | 不适合 |
| --- | --- |
| 个人项目后台 | 大型企业 IAM / 统一身份认证平台 |
| SaaS 原型 / MVP | 微服务治理平台 |
| 中小型管理系统（ERP、CRM、CMS 底座） | 低代码 / 无代码平台底座 |
| 需要权限、数据权限、动态菜单的后台底座 | 无需二次开发的商业成品系统 |
| 想学习 Go + Vue 全栈工程化的开发者 | 万级 QPS+ 的独立高并发场景 |

## 当前能力

### 后台基础

- **认证会话**：JWT access token + HttpOnly refresh token，支持 Redis 轮换和账号中心。
- **权限体系**：Casbin 接口授权、动态菜单、按钮权限、API 资源管理。
- **组织体系**：部门树、岗位管理、用户-角色-岗位关联。
- **数据权限**：支持所有数据、本部门、本部门及下级、仅本人、自定义部门五类范围。

### 系统模块

- **配置与字典**：系统配置、运行时配置读取、数据字典类型和字典项。
- **文件与附件**：文件上传、附件管理、上传白名单和公开访问路径。
- **日志审计**：登录日志、操作日志、操作详情查看。
- **通知消息**：通知公告、WebSocket 实时推送、消息模板、消息提醒。
- **邮件能力**：邮件账号、邮件模板、邮件日志、邮件预览、测试发送和手动发送。
- **运行观测**：健康检查页面、基础指标展示、运行状态反馈。

### 前端体验

- Vue 3 + TypeScript + Naive UI 管理台。
- 登录页、后台壳子、侧边栏、顶部栏、工作标签和动态菜单。
- 用户、角色、菜单、部门、岗位、系统配置、字典、附件、日志、通知、邮件、消息等管理页面。
- 统一的页面头部、帮助抽屉、表格容器、搜索区、状态标签和操作按钮组件。

## 技术栈

| 层 | 技术 |
| --- | --- |
| 后端框架 | Go 1.26 + Gin |
| ORM | GORM |
| 数据库 | PostgreSQL / MySQL |
| 权限引擎 | Casbin |
| 缓存与会话 | Redis |
| 前端框架 | Vue 3 + TypeScript |
| UI 组件库 | Naive UI |
| 状态管理 | Pinia |
| CSS | Tailwind CSS 4 |
| 构建工具 | Vite |
| 文档 | VitePress |
| 部署 | Docker Compose / Nginx / 服务器二进制 |

## 项目结构

```text
ez-admin-gin/
├── server/                 Go 后端
│   ├── configs/            配置文件
│   ├── docs/               Swagger / OpenAPI 规范
│   ├── internal/
│   │   ├── bootstrap/      启动引导、路由注册、Swagger 接入
│   │   ├── modules/        业务模块（auth / iam / system / setup）
│   │   ├── platform/       平台能力（authn / authz / config / middleware / model / ...）
│   │   └── pkg/            公共工具包
│   └── migrations/         MySQL / PostgreSQL 完整初始化 SQL
├── admin/                  Vue 3 管理台
│   └── src/
│       ├── components/     通用 UI 与后台壳子组件
│       ├── modules/        前端业务模块（auth / iam / system）
│       ├── router/         静态路由、动态菜单和路由守卫
│       ├── stores/         Pinia 状态
│       └── styles/         主题变量和全局样式
├── docs/                   VitePress 文档站
├── deploy/                 Docker Compose、Nginx、部署配置
├── scripts/                Swagger、初始化 SQL、部署脚本
└── MANUAL_TEST.md          发布前人工测试清单
```

## 快速启动

### 环境要求

| 工具 | 最低版本 | 用途 |
| --- | --- | --- |
| Go | 1.26.1+ | 后端 |
| Node.js | 20.19+ 或 22.12+ | 前端和文档 |
| pnpm | 9+ | 前端和文档包管理 |
| Docker | 20+ | 本地 PostgreSQL + Redis |
| make | GNU Make 4+ | 统一开发入口（Windows 可选） |

### 本地运行

```bash
# 查看所有可用命令
make help

# 1. 启动 PostgreSQL 和 Redis
make docker-up

# 2. 安装前端依赖
make install

# 3. 启动后端（另一个终端）
make server-dev

# 4. 初始化管理员账号
curl -X POST http://localhost:8080/api/v1/setup/init

# 5. 启动前端（另一个终端）
make admin-dev
```

启动后：

- 后端接口：`http://localhost:8080`
- 前端管理台：`http://localhost:5173`
- 默认管理员：`admin / EzAdmin@123456`

生产环境首次登录后请立即修改默认密码，并配置强随机 `EZ_AUTH_JWT_SECRET`。

## 常用命令

| 命令 | 说明 |
| --- | --- |
| `make help` | 显示所有可用命令 |
| `make docker-up` | 启动本地 PostgreSQL + Redis |
| `make docker-down` | 停止本地 PostgreSQL + Redis |
| `make install` | 安装管理台依赖 |
| `make server-dev` | 启动后端（`go run .`） |
| `make admin-dev` | 启动前端（`pnpm dev`） |
| `make docs-dev` | 启动文档站 |
| `make server-vet` | 后端 `go vet ./...` |
| `make admin-check` | 前端类型检查 + lint |
| `make build` | 构建后端二进制和前端产物 |
| `make docker-config` | 校验 Docker Compose 配置 |
| `make db-full-sql` | 生成 MySQL / PostgreSQL 完整初始化 SQL |
| `make swagger` | 生成 Swagger / OpenAPI 文档 |
| `make generate-types` | 基于 OpenAPI 生成前端类型 |
| `make check-types` | 检查前端生成类型是否同步 |

## 质量策略

本项目采用轻量质量策略，不追求复杂测试覆盖率。当前 GitHub Actions 的 `Lightweight Verify` 包含：

- 后端：`go mod download`
- 后端：`go mod tidy` 后检查 `go.mod` / `go.sum` 是否仍有差异
- 后端：`go vet ./...`
- 后端：`go build ./...`
- 前端：`pnpm install --frozen-lockfile`
- 前端：`pnpm type-check`
- 前端：`pnpm lint`
- 前端：`pnpm build`
- Docker：校验 `compose.local.yml`、`compose.prod.yml`、`compose.server.yml`

本地发布前建议执行：

```bash
cd server
go mod tidy
git diff --exit-code go.mod go.sum
go vet ./...
go build ./...
go test ./...

cd ../admin
pnpm install --frozen-lockfile
pnpm type-check
pnpm lint
pnpm build

cd ../docs
pnpm install --frozen-lockfile
pnpm docs:build

cd ..
docker compose -f deploy/compose.local.yml config --quiet
EZ_AUTH_JWT_SECRET=placeholder docker compose -f deploy/compose.prod.yml config --quiet
docker compose -f deploy/compose.server.yml config --quiet
```

复制到新项目或正式发布前，再按 [MANUAL_TEST.md](MANUAL_TEST.md) 做人工冒烟验证。

## 文档

在线文档：[https://caoshenyang.github.io/ez-admin-gin/](https://caoshenyang.github.io/ez-admin-gin/)

- [快速开始](https://caoshenyang.github.io/ez-admin-gin/getting-started/)
- [路线图](https://caoshenyang.github.io/ez-admin-gin/getting-started/roadmap)
- [系统架构](https://caoshenyang.github.io/ez-admin-gin/architecture/overview)
- [权限体系](https://caoshenyang.github.io/ez-admin-gin/architecture/rbac)
- [动态菜单](https://caoshenyang.github.io/ez-admin-gin/architecture/dynamic-menu)
- [后端开发](https://caoshenyang.github.io/ez-admin-gin/backend/overview)
- [前端开发](https://caoshenyang.github.io/ez-admin-gin/frontend/overview)
- [部署概览](https://caoshenyang.github.io/ez-admin-gin/deployment/overview)
- [生产环境检查清单](https://caoshenyang.github.io/ez-admin-gin/deployment/production-checklist)
- [参考手册](https://caoshenyang.github.io/ez-admin-gin/reference/)

## 版本状态

### v1.2.0

- 新增系统邮件管理：账号、模板、日志、预览、测试发送和手动发送。
- 新增系统消息与提醒管理：消息模板、提醒记录和管理页面。
- 强化后台体验：顶部栏、侧边栏、仪表盘、登录页、页面帮助抽屉和通用列表组件。
- 增强健康检查与系统配置页面。
- 同步动态菜单、权限文档、数据库 DDL 和完整初始化 SQL。
- 发布验证已对齐 GitHub Actions：后端 tidy/vet/build、前端 type/lint/build、Docker Compose 校验。

更多历史记录见 [CHANGELOG.md](CHANGELOG.md)。

## Roadmap

已完成：

- [x] JWT 认证 + Refresh Token 轮换
- [x] Casbin 接口权限、动态菜单和按钮权限
- [x] 组织体系：部门、岗位、用户-角色-岗位关联
- [x] 五级数据权限
- [x] 系统模块：配置、字典、附件、文件、日志、通知、邮件、消息、健康检查
- [x] 前端管理台：登录、后台壳子、动态菜单、账号中心、系统管理页面
- [x] PostgreSQL / MySQL 完整初始化 SQL
- [x] 多场景部署方案：Docker Compose、Nginx、服务器二进制部署
- [x] VitePress 文档站和 GitHub Pages 部署 workflow

轻量维护中：

- [ ] 按真实项目需要补充少量高复用业务模块
- [ ] 持续完善人工测试清单和复用前检查
- [ ] 继续收敛页面交互和模块接入规范

暂不计划：

- [ ] 微服务拆分
- [ ] 低代码 / 无代码引擎
- [ ] 复杂多租户隔离
- [ ] 大型 IAM 平台能力
- [ ] 社区治理和长期 PR 协作流程

## Contributing

本项目主要服务维护者自己的项目复用，公开源码供参考和学习。欢迎通过 Issue 反馈 Bug 或建议，但 Pull Request 不保证接受或处理。

详见 [CONTRIBUTING.md](CONTRIBUTING.md)。

## License

[MIT](LICENSE)
