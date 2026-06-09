---
title: EZ Admin Gin
description: "EZ Admin Gin 是维护者自用优先、适合个人项目和 SaaS/MVP 快速复用的全栈后台底座。"
layout: home

hero:
  name: EZ Admin Gin
  text: 自用优先的后台系统基座
  tagline: 面向个人项目和 SaaS/MVP 快速开发，保持结构清晰、部署简单、维护成本低
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
  - title: 认证与权限
    details: 登录、刷新令牌、退出登录、当前用户、账户中心、Casbin 接口授权、动态菜单和按钮权限已经打通。
  - title: IAM 基础模块
    details: 用户、角色、菜单、部门、岗位和接口资源管理都有后端接口与管理台页面支撑。
  - title: 系统管理能力
    details: 配置、字典、文件、附件、操作日志、登录日志、公告、消息、邮件、站内通知和健康检查集中在 System 模块。
  - title: 前后端分离架构
    details: Go + Gin 后端，Vue 3 + Naive UI 前端，模块按职责分层，适合继续扩展业务模块。
  - title: 多场景部署方案
    details: Docker Compose、Nginx、二进制部署、全容器化部署、更新回滚和生产检查清单都已收口到 deploy 与 scripts。
---

## 适合谁

- 个人项目需要快速上线一个管理后台
- SaaS 原型或 MVP 需要用户、角色、菜单和常用系统模块
- 中小型内部管理系统需要一个可继续扩展的后台底座
- 想基于现有权限、组织、文件、消息和部署能力做二次开发

## 怎么读这套文档

如果你第一次打开项目，先走 [快速开始](/getting-started/)；如果你准备改代码，先看 [项目结构](/getting-started/project-structure) 和 [系统架构概览](/architecture/overview)；如果你在查配置、权限码、数据库或部署参数，直接进 [参考手册](/reference/)。

::: tip 文档维护原则
这些页面以当前代码目录、路由注册、配置结构和初始化 SQL 为准。遇到文档与代码不一致时，优先相信代码，并回到 [当前系统地图](/reference/current-system-map) 找事实来源。
:::

## 不适合谁

- 直接当大型企业 IAM / 统一身份认证平台
- 微服务架构的服务治理平台
- 低代码 / 无代码平台底座
- 万级 QPS 以上的独立高并发系统
