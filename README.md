# EZ Admin Gin

面向个人项目快速上线的通用后台管理系统底座。功能完整、可直接部署、可二次扩展。

## 项目简介

EZ Admin 不是单独的后端模板，也不是只有页面壳子的前端 Demo，而是一套放在同一个单仓库里的后台底座：

- **`server/`** — Go + Gin 后端（认证、RBAC、数据权限、系统模块）
- **`admin/`** — Vue 3 管理台（Naive UI + Tailwind CSS）
- **`docs/`** — VitePress 文档站
- **`deploy/`** — Docker Compose、Nginx、部署脚本

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

## 快速启动

```bash
# 1. 启动 PostgreSQL 和 Redis
docker compose -f deploy/compose.local.yml up -d

# 2. 启动后端
cd server && go run .

# 3. 初始化管理员账号
curl -X POST http://localhost:8080/api/v1/setup/init \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"YourPassword123","nickname":"管理员"}'

# 4. 启动前端
cd admin && pnpm install && pnpm dev
```

## 项目结构

```
ez-admin-gin/
├── server/
│   ├── configs/            配置文件
│   ├── internal/           业务代码
│   │   ├── bootstrap/      启动引导
│   │   ├── modules/        业务模块（auth / iam / system / setup）
│   │   ├── platform/       平台层（authn / authz / datascope / middleware / ...）
│   │   └── pkg/            公共工具包
│   └── migrations/         数据库迁移（MySQL + PostgreSQL）
├── admin/
│   └── src/
│       ├── modules/        业务模块（auth / iam / system）
│       ├── layouts/        布局组件
│       ├── router/         路由 + 动态菜单注册
│       └── stores/         Pinia 状态管理
├── docs/                   VitePress 文档站
├── deploy/                 Docker Compose、Nginx、部署配置
└── scripts/                部署与打包脚本
```

## 文档

- [快速开始](https://caoshenyang.github.io/ez-admin-gin/getting-started/)
- [系统架构](https://caoshenyang.github.io/ez-admin-gin/architecture/overview)
- [权限体系](https://caoshenyang.github.io/ez-admin-gin/architecture/rbac)
- [后端开发](https://caoshenyang.github.io/ez-admin-gin/backend/overview)
- [前端开发](https://caoshenyang.github.io/ez-admin-gin/frontend/overview)
- [部署方案](https://caoshenyang.github.io/ez-admin-gin/deployment/overview)
- [参考手册](https://caoshenyang.github.io/ez-admin-gin/reference/)

## 权限体系

```
登录 → JWT Token → Auth 中间件 → LoadActor（角色+菜单+按钮权限）
  → Permission 中间件（Casbin: 角色 × URL × HTTP方法）
  → Repository 层（datascope: 五级数据权限过滤）
```

五级数据作用域：所有数据 / 本部门 / 本部门及下级 / 仅本人 / 自定义部门

## Roadmap

- [x] JWT 认证 + Casbin RBAC
- [x] 动态菜单与按钮权限
- [x] 组织体系（部门/岗位）
- [x] 五级数据权限
- [x] 系统模块（用户/角色/菜单/配置/字典/文件/日志/公告）
- [x] 前端管理台（登录/壳子/动态菜单/管理页面）
- [x] 多场景部署方案
- [ ] 前端主题切换（暗色模式）
- [ ] 国际化支持
- [ ] WebSocket 消息推送
- [ ] 审批工作流

## Contributing

欢迎提交 Issue 和 Pull Request。

## License

MIT
