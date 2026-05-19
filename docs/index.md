title: EZ Admin Gin
description: "EZ Admin Gin 是面向个人项目快速上线、也适合中小型后台系统复用的全栈后台底座。"
layout: home

hero:
  name: EZ Admin Gin
  text: 全栈后台管理系统底座
  tagline: 高效、易用、企业级后台框架，开箱即用的权限体系与部署方案
  image:
    src: /images/logo-stacked.svg
    alt: EZ Admin Gin 品牌 Logo
  actions:
    - theme: brand
      text: 快速开始
      link: /getting-started/
    - theme: alt
      text: 系统架构
      link: /architecture/overview
    - theme: alt
      text: GitHub
      link: https://github.com/caoshenyang/ez-admin-gin

features:
  - title: 完整 RBAC 权限体系
    details: 基于 Casbin 的接口级权限控制，五级数据权限作用域，动态菜单与按钮权限，开箱即用。
  - title: 企业级后台能力
    details: 用户、角色、部门、岗位、菜单、字典、配置、文件、日志、公告——后台底座的标准能力全部就绪。
  - title: 前后端分离架构
    details: Go + Gin 后端，Vue 3 + Naive UI 前端，模块化分层设计，清晰的扩展边界。
  - title: 多场景部署方案
    details: Docker Compose 一键编排，Nginx 反向代理，支持本地开发、服务器部署、云端部署和生产环境。
---

## 适合谁

- 个人项目需要快速上线一个管理后台
- SaaS 原型或 MVP 需要权限体系和用户管理
- 中小型内部管理系统（ERP、CRM、CMS 底座）
- 二次开发：基于现有模块扩展业务功能

## 不适合谁

- 直接当大型企业 IAM / 统一身份认证平台
- 微服务架构的服务治理平台
- 低代码 / 无代码平台底座
- 高并发（万级 QPS+）独立场景 — 需要额外优化
