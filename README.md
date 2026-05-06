# EZ Admin Gin

面向个人项目快速上线的通用后台管理系统底座，也是一套面向 Java 转 Go 工程师的完整实战样板。当前稳定版已经收敛到登录、权限、菜单、组织体系、数据权限、前端管理台、部署复用这一条完整主线。

当前稳定版本：**v1.1.0**

## 项目是什么

这个仓库不是单独的后端模板，也不是只有页面壳子的前端 Demo，而是一套放在同一个单仓库里的后台底座：

- `server/`：Go + Gin 后端
- `admin/`：Vue 3 管理台
- `docs/`：VitePress 文档站
- `deploy/`：Docker Compose、Nginx、部署脚本

适合的场景主要有三类：

- 想把熟悉的 Java 后台经验迁到 Go
- 想要一套可长期扩展、可部署、可复用的后台底座
- 想系统理解认证、权限、菜单、组织体系、数据权限和模块化接入

## 五步跑起来

### 1. 启动 PostgreSQL 和 Redis

```bash
# macOS / Linux
docker compose -f deploy/compose.local.yml up -d

# Windows
docker compose -f deploy/compose.local.win.yml up -d
```

### 2. 启动后端

```bash
cd server
go run ./cmd/server
```

后端默认监听 `http://localhost:8080`，首次启动会自动执行迁移。

### 3. 初始化管理员账号

```bash
curl -X POST http://localhost:8080/api/v1/setup/init \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"YourPassword123","nickname":"管理员"}'
```

### 4. 启动前端

```bash
cd admin
pnpm install
pnpm dev
```

### 5. 启动文档站

```bash
cd docs
pnpm install
pnpm docs:dev
```

## 从哪里开始读

新读者建议固定按这条路径进入：

1. [快速启动与项目入口](https://caoshenyang.github.io/ez-admin-gin/guide/)
2. [从零搭建教程](https://caoshenyang.github.io/ez-admin-gin/tutorial/)
3. [参考手册](https://caoshenyang.github.io/ez-admin-gin/reference/)

当前教程主线已经稳定为 **9 章**：

- 第 1-2 章：项目初始化、后端基础设施
- 第 3-4 章：认证、RBAC、Casbin、菜单权限
- 第 5 章：组织体系与数据权限
- 第 6 章：核心系统模块
- 第 7 章：前端企业级管理台
- 第 8 章：模块化接入规范
- 第 9 章：部署、升级与复用

## 当前已经具备的能力

### 后端

- JWT 登录认证
- RBAC 角色权限与 Casbin 接口授权
- 动态菜单与按钮权限
- 组织体系基础模型：部门、岗位、用户岗位关系
- 数据权限基础模型：角色数据范围、自定义部门范围
- 用户、角色、菜单、配置、文件、日志、公告、字典、附件等系统模块
- 统一响应格式、错误处理和请求级操作日志

### 前端

- 登录页
- 后台壳子、动态菜单、工作标签
- 用户、角色、菜单、配置、文件、日志、字典、附件、CRM 示例页面
- Dashboard 和账户中心

### 部署

- 后端打包脚本与 Linux 二进制部署
- 前端静态资源构建
- Docker Compose 基础设施编排
- Nginx 反向代理配置
- 初始化数据与升级回滚说明

## 项目结构

```text
ez-admin-gin/
├── server/
│   ├── cmd/
│   ├── configs/
│   ├── internal/
│   └── migrations/
├── admin/
├── docs/
├── deploy/
└── scripts/
```

## 更多文档

- [快速启动](https://caoshenyang.github.io/ez-admin-gin/guide/)
- [项目结构](https://caoshenyang.github.io/ez-admin-gin/guide/project-structure)
- [企业级架构升级](https://caoshenyang.github.io/ez-admin-gin/guide/enterprise-architecture)
- [Go vs Java 工程结构](https://caoshenyang.github.io/ez-admin-gin/guide/java-to-go-structure)
- [路线图](https://caoshenyang.github.io/ez-admin-gin/guide/roadmap)

## License

MIT
