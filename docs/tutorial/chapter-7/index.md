---
title: 第 7 章：部署与复用
description: "完成 Docker 镜像构建、Compose 编排、Nginx 配置、环境变量管理和初始化数据，让后台底座可以被快速部署和复用。"
---

# 第 7 章：部署与复用

前六章已经完成了后台底座的核心功能：后端认证、权限、系统模块，前端登录、布局、菜单和页面，以及业务模块的接入规范。现在要把这些能力打包成可部署的完整产物，让它能在 Docker 环境下稳定运行，也能被下一个个人项目快速复用。

::: tip 🎯 本章怎么验证
本章完成后，你应该能通过一条 `docker compose up` 命令把整个后台跑起来，用浏览器访问、登录、看到完整菜单并执行 CRUD 操作。之后如果需要开一个新项目，也知道该改哪些地方。
:::

## 为什么需要这一章

一个后台底座如果只能在本机 `go run` 跑起来，还不能算"可用"。真正让它变成可交付产物的几个条件：

- 可以用 Docker 一键构建和启动，不依赖本机环境。
- 环境变量和配置分离，本地开发和生产部署互不干扰。
- 服务启动时自动初始化必要的数据，不需要手动导入 SQL。
- 结构清晰到你可以直接复制仓库、改名、加入新业务，不需要重新搭建。

这一章会把上面这些能力补齐。

## 开始前准备

进入本章前，建议先确认：

- 后端服务可以正常启动，系统模块接口可用。
- 前端管理台可以登录，菜单和页面功能正常。
- 第 6 章的业务模块接入规范已经理解。
- 本机已安装 Docker 和 Docker Compose（V2）。

## 本章会完成什么

本章会按下面的顺序推进：

<table>
  <colgroup>
    <col style="width: 11rem;">
    <col>
    <col>
  </colgroup>
  <thead>
    <tr>
      <th>小节</th>
      <th>完成什么</th>
      <th>验证重点</th>
    </tr>
  </thead>
  <tbody>
    <tr>
      <td>后端与前端 Dockerfile</td>
      <td>为后端服务和前端管理台编写多阶段构建的 Dockerfile</td>
      <td>镜像可以成功构建并运行</td>
    </tr>
    <tr>
      <td>Docker Compose 编排</td>
      <td>编排 PostgreSQL、Redis、后端和 Nginx 四个服务</td>
      <td>所有容器正常启动，服务间网络互通</td>
    </tr>
    <tr>
      <td>Nginx 配置</td>
      <td>配置静态资源托管和 API 反向代理</td>
      <td>前端页面和后端接口都能通过 Nginx 访问</td>
    </tr>
    <tr>
      <td>环境变量与初始化数据</td>
      <td>整理部署环境变量，了解启动时自动初始化的数据</td>
      <td>配置覆盖机制正确，初始化数据幂等</td>
    </tr>
    <tr>
      <td>部署验证与复用说明</td>
      <td>完成端到端部署验证，说明如何复用到新项目</td>
      <td>浏览器可以登录并操作，复用步骤可执行</td>
    </tr>
  </tbody>
</table>

## 本章小节

- [后端与前端 Dockerfile](./dockerfile)
- [Docker Compose 编排](./docker-compose)
- [Nginx 配置](./nginx-config)
- [环境变量与初始化数据](./env-and-init-data)
- [部署验证与复用说明](./deployment-and-reuse)

## 本章完成后的状态

完成本章后，后台底座会从"本地可开发"推进到"线上可部署、新项目可复用"：

```text
编写 Dockerfile（后端 + 前端）
  ↓
编排 Docker Compose（数据库 + 缓存 + 后端 + Nginx）
  ↓
配置 Nginx（静态托管 + API 代理）
  ↓
整理环境变量，理解初始化数据
  ↓
端到端部署验证 + 复用指南
```

下一节从 Dockerfile 开始：[后端与前端 Dockerfile](./dockerfile)。
