---
title: 更新日志
description: "EZ Admin Gin 版本更新记录。"
---

# 更新日志

本页记录每个版本的变更内容。完整文件见仓库根目录 [CHANGELOG.md](https://github.com/Caoshenyang/ez-admin-gin/blob/main/CHANGELOG.md)。

## v1.1.0 (2026-04-30)

围绕“更稳地部署、更顺地阅读文档、前端品牌感更统一”完成了一轮迭代。

### 后端与部署

- 引入 `golang-migrate` 迁移流程，并补充 MySQL 支持
- 统一默认管理员密码为 `EzAdmin@123456`
- 拆分服务端 `setup` / `update` 脚本，补充 Linux 一键部署脚本

### 前端

- 后台侧栏品牌区升级为 logo + 文字组合
- 登录页与后台壳子统一品牌 Logo 展示
- 动态菜单图标白名单和菜单示例同步整理

### 文档

- 补充 VitePress GitHub Pages 部署方案
- 第 5-7 章教程改为更完整的内联代码块，减少引用跳转
- 修正文档中的示例语言标记、favicon 路径和部署细节

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
