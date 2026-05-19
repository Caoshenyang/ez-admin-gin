# EZ Admin Gin

基于 Go + Gin + Vue 3 的全栈后台管理系统底座，适合个人项目快速上线、中小型后台系统和 SaaS 原型二次开发。

当前仓库已进入稳定收尾阶段，主目标是保持交付完整、文档一致和数据库初始化链路清晰，而不是继续扩展示例功能。

## 适合 / 不适合

| 适合 | 不适合 |
|------|--------|
| 个人项目后台 | 大型企业 IAM / 统一身份认证平台 |
| SaaS 原型 / MVP | 微服务架构的服务治理平台 |
| 中小型管理系统（ERP、CRM、CMS 底座） | 低代码 / 无代码平台底座 |
| 需要 RBAC / 数据权限 / 动态菜单的后台底座 | 无需二次开发的商业成品系统 |
| 想学习 Go + Vue 全栈工程化的开发者 | 高并发（万级 QPS+）独立场景 |

## 核心特性

- **完整 RBAC 权限体系** — JWT 认证 + Casbin 接口授权 + 五级数据权限 + 动态菜单 + 按钮权限
- **企业级组织体系** — 部门树、岗位管理、用户-角色-岗位多对多关联
- **模块化扩展** — 统一四层结构（handler/service/repository/domain），标准接入流程
- **全场景部署** — 5 种 Docker Compose 变体 + Nginx HTTP/HTTPS + 一键部署脚本
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

## 系统架构

```
┌─────────────────────────────────────────────────┐
│                    Nginx                         │
│            (反向代理 + 静态资源)                   │
├────────────────────┬────────────────────────────┤
│    Go + Gin 后端    │     Vue 3 前端              │
│  ┌──────────────┐  │  ┌──────────────────────┐  │
│  │  Middleware   │  │  │  Router (动态菜单)    │  │
│  ├──────────────┤  │  ├──────────────────────┤  │
│  │  Handler     │  │  │  Composables         │  │
│  ├──────────────┤  │  ├──────────────────────┤  │
│  │  Service     │  │  │  Components          │  │
│  ├──────────────┤  │  └──────────────────────┘  │
│  │  Repository  │  │                             │
│  │  (datascope) │  │                             │
│  └──────┬───────┘  │                             │
│         │          │                             │
├─────────┼──────────┼─────────────────────────────┤
│  PostgreSQL / MySQL  │  Redis                     │
└─────────────────────────────────────────────────┘
```

## 项目结构

```
ez-admin-gin/
├── server/                 Go 后端
│   ├── configs/            配置文件
│   ├── internal/           业务代码
│   │   ├── bootstrap/      启动引导
│   │   ├── modules/        业务模块（auth / iam / system / setup）
│   │   ├── platform/       平台层（authn / authz / datascope / middleware / ...）
│   │   └── pkg/            公共工具包
│   ├── migrations/         数据库迁移（MySQL + PostgreSQL）
│   ├── tests/              测试（API / RBAC / Contract / Testutil）
│   └── docs/               Swagger / OpenAPI 规范
├── admin/                  Vue 3 管理台
│   └── src/
│       ├── modules/        业务模块（auth / iam / system）
│       ├── layouts/        布局组件
│       ├── router/         路由 + 动态菜单注册
│       └── stores/         Pinia 状态管理
├── docs/                   VitePress 文档站
├── deploy/                 Docker Compose、Nginx、部署配置
└── scripts/                部署与打包脚本
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

### 不使用 make 的等效命令

```bash
# 1. 启动 PostgreSQL 和 Redis
docker compose -f deploy/compose.local.yml up -d

# 2. 如需更新完整版初始化 SQL（可选）
./scripts/build-full-migrations.sh

# 3. 启动后端
cd server && go run .

# 4. 初始化管理员账号
curl -X POST http://localhost:8080/api/v1/setup/init \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"YourPassword123","nickname":"管理员"}'

# 5. 启动前端
cd admin && pnpm install && pnpm dev
```

::: warning Windows 用户
安装 make：`choco install make` 或 `scoop install make`。如果暂时不想安装，也可以直接查看 Makefile 中的对应命令手动执行。
:::

### 默认账号说明

初始化时通过 `setup/init` 接口自行指定用户名和密码。**生产环境首次登录后请立即修改密码。**

### 数据库初始化说明

数据库对外交付以两份完整版 SQL 为准：

- `server/migrations/mysql/full_schema_and_seed.sql`
- `server/migrations/postgres/full_schema_and_seed.sql`

它们只负责系统表和内置种子，不会写死真实管理员账号。首个管理员仍通过 `setup/init` 创建。

## 常用命令速查

| 命令 | 说明 |
|------|------|
| `make help` | 显示所有可用命令 |
| `make server-dev` | 启动后端 (`go run .`) |
| `make admin-dev` | 启动前端 (`pnpm dev`) |
| `make docs-dev` | 启动文档站 |
| `make test` | 运行后端测试 (`go test ./...` + 契约测试) |
| `make test-contract` | OpenAPI 契约测试（不需要 DB） |
| `make test-integration` | 集成测试（需要 DB + Redis） |
| `make server-vet` | 后端代码检查 (`go vet ./...`) |
| `make admin-check` | 前端类型检查 + lint |
| `make lint` | 运行所有 lint (后端 vet + 前端检查 + 契约一致性) |
| `make build` | 构建后端二进制 + 前端产物 |
| `make docs-build` | 构建文档站 |
| `make db-full-sql` | 生成 MySQL / PostgreSQL 完整版初始化 SQL |
| `make docker-up` | 启动 PostgreSQL + Redis |
| `make docker-down` | 停止 PostgreSQL + Redis |
| `make docker-config` | 验证所有 Docker Compose 配置 |
| `make generate-types` | 从 Swagger spec 生成前端 TypeScript 类型 |

## 文档

在线文档：[https://caoshenyang.github.io/ez-admin-gin/](https://caoshenyang.github.io/ez-admin-gin/)

- [快速开始](https://caoshenyang.github.io/ez-admin-gin/getting-started/)
- [系统架构](https://caoshenyang.github.io/ez-admin-gin/architecture/overview)
- [权限体系](https://caoshenyang.github.io/ez-admin-gin/architecture/rbac)
- [后端开发](https://caoshenyang.github.io/ez-admin-gin/backend/overview)
- [前端开发](https://caoshenyang.github.io/ez-admin-gin/frontend/overview)
- [部署概览](https://caoshenyang.github.io/ez-admin-gin/deployment/overview)
- [服务器二进制部署](https://caoshenyang.github.io/ez-admin-gin/deployment/server-binary-deploy)
- [Docker 部署](https://caoshenyang.github.io/ez-admin-gin/deployment/docker-deploy)
- [生产环境检查清单](https://caoshenyang.github.io/ez-admin-gin/deployment/production-checklist)
- [参考手册](https://caoshenyang.github.io/ez-admin-gin/reference/)

## 权限体系

```
登录 → JWT Token → Auth 中间件 → LoadActor（角色+菜单+按钮权限）
  → Permission 中间件（Casbin: 角色 × URL × HTTP方法）
  → Repository 层（datascope: 五级数据权限过滤）
```

五级数据作用域：所有数据 / 本部门 / 本部门及下级 / 仅本人 / 自定义部门

## CI / 质量门禁

每次 push 和 pull request 都会自动运行以下检查：

| Job | 检查内容 | 阻塞 |
|-----|---------|------|
| **Backend** | `go mod tidy` 一致性、`go vet`、`go test` | 是 |
| **Integration** | API 黑盒测试 + RBAC 权限流程测试（PostgreSQL + Redis） | 是 |
| **Frontend** | API 类型同步、TypeScript 类型检查、ESLint + oxlint、生产构建 | 是 |
| **Docker** | compose.local / compose.prod / compose.server 配置校验 | 是 |
| **Security** | Gitleaks 密钥扫描、govulncheck Go 漏洞检查 | 密钥扫描阻塞，漏洞扫描仅告警 |

> **govulncheck 为非阻塞**：当前项目依赖链较复杂，部分已知漏洞可能需要依赖上游修复，因此仅作为告警。后续稳定后可改为阻塞。

本地模拟 CI：

```bash
# 后端
cd server && go mod tidy && git diff --exit-code go.mod go.sum
cd server && go vet ./... && go test ./...
cd server && go test -v -timeout 60s ./tests/contract/...

# 前端
cd admin && pnpm install --frozen-lockfile && pnpm type-check && pnpm lint && pnpm build

# Docker
docker compose -f deploy/compose.local.yml config --quiet
EZ_AUTH_JWT_SECRET=placeholder docker compose -f deploy/compose.prod.yml config --quiet
docker compose -f deploy/compose.server.yml config --quiet
```

## Roadmap

已完成：

- [x] JWT 认证 + Casbin RBAC
- [x] 动态菜单与按钮权限
- [x] 组织体系（部门/岗位）
- [x] 五级数据权限
- [x] 系统模块（用户/角色/菜单/配置/字典/文件/日志/公告）
- [x] 前端管理台（登录/壳子/动态菜单/管理页面）
- [x] 多场景部署方案（5 种 Docker Compose 变体 + 一键部署脚本）
- [x] WebSocket 通知公告实时推送

计划中：

- [ ] 前端主题切换（暗色模式）
- [ ] 国际化支持
- [ ] 审批工作流
- [ ] 更多业务模板

暂不计划：

- [ ] 微服务拆分
- [ ] 低代码引擎
- [ ] 多租户隔离

## Contributing

欢迎通过 Issue 反馈 Bug、提出建议或分享使用心得。

当前以仓库维护者主导的稳定化和收尾为主。若要协作实现，建议先通过 Issue 对齐范围和方案，详见 [CONTRIBUTING.md](CONTRIBUTING.md)。

## License

[MIT](LICENSE)
