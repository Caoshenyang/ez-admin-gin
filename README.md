# EZ Admin Gin

面向个人项目快速上线的通用后台管理系统底座。登录、权限、菜单、配置、日志、文件上传等后台基础能力一次性沉淀，新项目直接复用。

## 适用场景

- 个人开发者需要一个"能直接用"的后台管理系统
- 想从零学习 Go + Vue 3 全栈后台搭建
- 已有业务 idea，不想重复写登录、权限、用户管理

## 技术栈

| 层 | 技术 |
| --- | --- |
| 后端 | Go 1.26、Gin、GORM、PostgreSQL、Redis、Casbin |
| 前端 | Vue 3.5、TypeScript、Naive UI、TailwindCSS 4、Vite 8 |
| 文档 | VitePress 2.0 |
| 部署 | Docker Compose、Nginx |

## 快速启动

### 环境要求

- Go >= 1.26
- Node.js >= 20.19
- pnpm
- Docker & Docker Compose（本地数据库和 Redis）

### 1. 启动基础服务

```bash
# macOS / Linux
docker compose -f deploy/compose.local.yml up -d

# Windows
docker compose -f deploy/compose.local.win.yml up -d
```

PostgreSQL 和 Redis 会自动启动，数据持久化到本机目录。

### 2. 启动后端

```bash
cd server
cp configs/config.yaml.example configs/config.yaml  # 按需修改
go run main.go
```

首次启动会自动创建数据库表和默认管理员账号。

### 3. 启动前端

```bash
cd admin
pnpm install
pnpm dev
```

### 4. 启动文档站

```bash
cd docs
pnpm install
pnpm docs:dev
```

## 默认账号

| 项目 | 值 |
| --- | --- |
| 用户名 | `admin` |
| 密码 | `EzAdmin@123456` |

::: warning
生产环境请务必修改默认密码和 JWT 密钥。
:::

## 功能清单

### 后端

- JWT 登录认证
- RBAC 角色权限（Casbin）
- 动态菜单（目录 / 菜单 / 按钮三级）
- 用户管理、角色管理、菜单管理
- 系统配置（键值对，支持启用 / 禁用）
- 文件上传（本地存储，白名单后缀）
- 操作日志、登录日志
- 公告管理
- 统一响应格式与错误处理
- 请求级操作日志中间件

### 前端

- 登录页
- 侧边栏 + 顶栏后台布局
- 动态菜单渲染
- 用户、角色、菜单、配置、文件、日志管理页面
- Dashboard 数据概览

### 部署

- 后端 Dockerfile（多阶段构建）
- 前端 Dockerfile（Nginx 托管）
- 生产环境 Docker Compose 编排
- Nginx 反向代理配置
- 环境变量配置（`.env.example`）
- 自动初始化数据和权限种子

## 项目结构

```
ez-admin-gin/
├── server/          # Go 后端
│   ├── configs/     # 配置文件
│   └── internal/    # 业务代码
├── admin/           # Vue 3 前端
│   └── src/
├── docs/            # VitePress 文档站
├── deploy/          # Docker Compose 和 Nginx 配置
│   └── nginx/
└── .agents/         # 开发辅助 skill
```

## 文档

- [使用指南](https://caoshenyang.github.io/ez-admin-gin/guide/)
- [从零搭建教程](https://caoshenyang.github.io/ez-admin-gin/tutorial/)（7 章，从空仓库到可部署）
- [参考手册](https://caoshenyang.github.io/ez-admin-gin/reference/)
- [路线图](https://caoshenyang.github.io/ez-admin-gin/roadmap)

## 部署

```bash
cd deploy
cp .env.example .env   # 修改生产环境配置
docker compose -f compose.prod.yml up -d
```

详见教程 [第 7 章：部署与复用](https://caoshenyang.github.io/ez-admin-gin/tutorial/chapter-7/)。

## 版本

当前稳定版本：**v1.1.0**

最后验证日期：2026-04-30

## License

MIT
