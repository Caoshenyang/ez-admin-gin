---
title: 更新日志
description: "EZ Admin Gin 版本更新记录。"
---

# 更新日志

本页记录每个版本的变更内容。完整文件见仓库根目录 [CHANGELOG.md](https://github.com/Caoshenyang/ez-admin-gin/blob/main/CHANGELOG.md)。

## v1.0.0 (2026-04-27)

首个正式版本，包含完整的后台管理系统底座能力。

### 后端

- Go + Gin 路由、Viper 配置、Zap 日志、GORM 数据库、Redis 连接
- JWT 登录认证和中间件校验
- RBAC 权限：用户、角色、菜单三级关联 + Casbin 接口控制
- 动态菜单（目录 / 菜单 / 按钮三级）
- 用户、角色、菜单、系统配置、文件上传、操作日志、登录日志、公告管理
- Dashboard 统计接口
- 启动时自动初始化默认数据和权限种子

### 前端

- Vue 3 + Naive UI + TailwindCSS 4 + TypeScript 管理台
- 登录页、后台布局、动态菜单
- 用户、角色、菜单、配置、文件、日志管理页面

### 文档

- VitePress 文档站
- 7 章 40+ 节从零搭建教程
- 参考手册（GORM、Casbin、接口风格、DDL、逻辑删除）

### 部署

- Docker Compose 本地环境（PostgreSQL + Redis）
- 生产环境 Dockerfile、Docker Compose、Nginx 配置
- macOS 和 Windows 平台适配
