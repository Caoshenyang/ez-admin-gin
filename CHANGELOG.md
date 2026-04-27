# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/).

## [1.0.0] - 2026-04-27

### Added

- Go 后端项目骨架：Gin 路由、Viper 配置、Zap 日志、GORM 数据库、Redis 连接
- JWT 登录认证，支持 Token 签发、解析和中间件校验
- RBAC 权限模型：用户、角色、菜单三级关联，Casbin 接口权限控制
- 动态菜单系统：目录 / 菜单 / 按钮三种类型，支持启用 / 禁用
- 用户管理：增删改查、状态切换、角色分配
- 角色管理：增删改查、状态切换、接口权限分配、菜单权限分配
- 菜单管理：增删改查、树形结构、状态切换
- 系统配置：键值对存储、启用 / 禁用、按 key 读取
- 文件上传：本地存储、后缀白名单、大小限制
- 操作日志中间件：自动记录请求级操作日志
- 登录日志：记录登录时间、IP、User-Agent
- 公告管理：增删改查、状态切换
- Dashboard 接口：用户数、角色数、菜单数等统计数据
- 服务启动时自动初始化默认管理员、角色、菜单和权限种子
- 环境变量覆盖配置（`EZ_` 前缀）
- 统一响应格式和错误处理

- Vue 3 前端管理台：Naive UI + TailwindCSS 4 + TypeScript
- 登录页
- 后台布局：侧边栏 + 顶栏 + 内容区
- 动态菜单渲染
- 用户、角色、菜单、配置、文件、日志管理页面
- Axios 封装和 Pinia 状态管理

- VitePress 文档站
- 从零搭建教程（7 章 40+ 节，从空仓库到可部署）
- 参考手册（GORM、Casbin 快速入门、接口风格决策、DDL、逻辑删除）
- 路线图

- Docker Compose 本地环境（PostgreSQL + Redis）
- 生产部署配置：Dockerfile、Docker Compose、Nginx 反向代理
- 环境变量模板（`.env.example`）
- macOS 和 Windows 平台适配

### Docs

- 完整教程大纲和 7 章教学内容
- GORM、Casbin 参考手册
- 接口风格决策记录
- 数据库建表语句参考
- 逻辑删除与唯一索引冲突分析
