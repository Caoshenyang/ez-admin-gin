# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/).

## [Unreleased]

## [1.2.0] - 2026-06-09

### Changed

- 移除误提交的 server/server 二进制文件（38 MB ELF）
- 修复 SECURITY.md 限流默认值：Per-IP 从 200 req/60s 修正为 10 req/60s（与代码一致）
- 修复 docs/architecture/rbac.md 代码块语言标记（conf → ini）
- .gitignore 补充 *.out、*.test、tmp/ 等构建产物忽略规则

### Documentation

- README 升级：增加适合/不适合表格、环境要求、默认账号说明、License 链接
- README Roadmap 同步：WebSocket 通知推送移至已完成，暂不计划项明确列出
- 新增 .github/workflows/ci.yml（后端 + 集成测试 + 前端 + Docker + 安全扫描）
- 新增 .github/workflows/deploy-docs.yml（VitePress 部署到 GitHub Pages）

### Added

- 企业级组织体系：部门树、岗位管理、用户-岗位多对多关联
- 五级数据权限（所有数据 / 本部门 / 本部门及下级 / 仅本人 / 自定义部门）
- 角色数据范围配置表（`sys_role_data_scope`）
- PostgreSQL / MySQL 迁移文件（`000003_enterprise_foundation`）
- 文档页：架构设计、权限体系、数据权限、组织体系、模块扩展机制
- 后端开发指南、前端开发指南、部署指南（二进制部署、Docker 部署、域名 HTTPS、更新回滚）
- 参考手册：环境变量、错误码、权限码约定、目录约定、模块规范、Nginx 配置等
- CI 流水线：后端测试 + 前端检查 + Docker 校验 + 安全扫描 + 集成测试
- 安全：JWT secret 生产环境强制校验、HttpOnly cookie refresh token、安全 headers、上传白名单、登录限流、生产配置校验
- WebSocket 通知公告实时推送
- 前端品牌 Logo 与品牌色统一

### Changed

- 项目定位从"功能型后台模板"升级为"面向个人项目和小型团队的通用后台底座"
- 前端样式从硬编码颜色迁移到 CSS 变量体系
- 默认管理员密码统一为 `EzAdmin@123456`
- 动态菜单图标白名单和前端组件注册机制同步更新
- 双 Token 机制：access token（Bearer header）+ refresh token（HttpOnly cookie + Redis 轮换）

## [1.1.0] - 2026-04-30

### Added

- 基于 `golang-migrate` 的数据库迁移流程，补充 MySQL 支持
- 文档站 GitHub Pages 部署工作流（`deploy-docs.yml`）
- Linux 一键部署脚本 `scripts/deploy.sh`，以及 `setup-server.sh` / `update-server.sh`
- Windows 打包脚本 `scripts/pack.ps1`
- 前端品牌 Logo 资源与复用组件

### Changed

- 默认管理员密码统一为 `EzAdmin@123456`
- 侧栏品牌区调整为图形 Logo + 品牌文字
- 动态菜单图标白名单同步更新
- 文档部署环境切换为 pnpm + Node 22

### Fixed

- 迁移目录命名和 DSN 处理逻辑
- 文档部署环境兼容性
- favicon 路径等文档站细节问题

## [1.0.0] - 2026-04-27

### Added

- Go 后端骨架：Gin 路由、Viper 配置、Zap 日志、GORM 数据库、Redis 连接
- JWT 登录认证 + Token 签发、解析、中间件校验
- RBAC 权限：用户-角色-菜单三级关联 + Casbin 接口权限控制
- 动态菜单系统：目录/菜单/按钮三种类型
- 用户管理：CRUD、状态切换、角色分配
- 角色管理：CRUD、接口权限分配、菜单权限分配
- 菜单管理：CRUD、树形结构
- 系统配置、文件上传（白名单）、操作日志、登录日志、公告管理
- Dashboard 统计接口
- 环境变量覆盖（`EZ_` 前缀）
- 统一响应格式和错误处理
- Vue 3 前端管理台（Naive UI + TailwindCSS + TypeScript）
- 登录页、后台布局、动态菜单渲染、管理页面
- VitePress 文档站
- Docker Compose 本地环境（PostgreSQL + Redis）
- 生产部署配置：Dockerfile、Nginx 反向代理
- 环境变量模板（`.env.example`）
