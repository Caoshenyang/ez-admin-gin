---
title: 部署与运维
description: "把 EZ Admin Gin 从本地跑通推进到服务器上线、更新、回滚和长期维护的一组说明页面。"
---

# 部署与运维

如果你现在最关心的是：

- 这套项目怎么上线
- 服务器上到底怎么跑
- 环境变量和初始化数据怎么准备
- Nginx、HTTPS、更新、回滚怎么做

那这一组页面就是给你准备的。

::: tip 🎯 这组页面解决什么
它不再按“教程章节”来组织，而是按“上线任务”来组织。你可以把这里当成当前仓库的部署说明站入口。
:::

## 先明确这组页面适合谁

这组内容特别适合三类读者：

1. 你已经把项目代码拉下来，准备把它部署到一台真实服务器上
2. 你不是来从零学习实现，而是来确认现成仓库的上线方法
3. 你后面准备把这套底座复用到自己的项目里，想先看清部署边界

如果你现在想从空仓库一步一步跟着做，还是建议回到 [从零搭建教程](/tutorial/)；  
如果你现在主要是“让它上线”，直接从这组页面往下读会更顺。

## 推荐阅读顺序

部署主线建议按下面顺序读：

1. [部署验证与复用说明](/tutorial/chapter-8/deployment-and-reuse)
2. [Compose 与服务运行结构](/tutorial/chapter-8/compose-and-service-layout)
3. [环境变量与初始化数据](/tutorial/chapter-8/env-and-init-data)
4. [Nginx 与 HTTPS 入口层](/tutorial/chapter-8/nginx-and-https)
5. [更新与回滚策略](/tutorial/chapter-8/update-and-rollback)
6. [部署排障 FAQ](/tutorial/chapter-8/deployment-troubleshooting-faq)
7. [长期运维 FAQ](/tutorial/chapter-8/operations-maintenance-faq)
8. [新项目复用清单](/tutorial/chapter-8/project-reuse-checklist)

## 如果你只想先跑通一版上线

最短路径建议是：

1. 先读 [部署验证与复用说明](/tutorial/chapter-8/deployment-and-reuse)
2. 再读 [环境变量与初始化数据](/tutorial/chapter-8/env-and-init-data)
3. 接着读 [Compose 与服务运行结构](/tutorial/chapter-8/compose-and-service-layout)
4. 最后补 [Nginx 与 HTTPS 入口层](/tutorial/chapter-8/nginx-and-https)

这样你会先知道：

- 当前打包和上传链路是什么
- 服务器上的真实运行拓扑是什么
- 配置、迁移和管理员初始化怎么衔接
- 最终入口层怎样把页面、API 和 HTTPS 接起来

## 页面职责怎么分

为了避免“主线、FAQ、参考页全混在一起”，部署篇现在按下面方式分工：

### 主线页

负责：

- 讲顺序
- 讲动作
- 讲验证
- 讲失败排查入口

当前主线页包括：

- [部署验证与复用说明](/tutorial/chapter-8/deployment-and-reuse)
- [Compose 与服务运行结构](/tutorial/chapter-8/compose-and-service-layout)
- [环境变量与初始化数据](/tutorial/chapter-8/env-and-init-data)
- [Nginx 与 HTTPS 入口层](/tutorial/chapter-8/nginx-and-https)
- [更新与回滚策略](/tutorial/chapter-8/update-and-rollback)

### FAQ 页

负责：

- 回答上线后最常遇到的问题
- 提供故障判断顺序
- 尽量不重复主线正文

当前 FAQ 页包括：

- [部署排障 FAQ](/tutorial/chapter-8/deployment-troubleshooting-faq)
- [长期运维 FAQ](/tutorial/chapter-8/operations-maintenance-faq)

### 参考页

负责：

- 查参数
- 查模板
- 查文件清单

当前最常用的部署参考页包括：

- [环境变量参考](/reference/environment-variables-reference)
- [初始化数据参考](/reference/init-data-reference)
- [Docker 部署文件参考](/reference/deploy-artifacts-reference)
- [Nginx 配置参考](/reference/nginx-config-reference)

## 当前仓库默认部署路线

当前仓库默认走的是一条非常务实的部署主线：

```text
本地打包
  ↓
上传部署包
  ↓
服务器首次执行 setup-server.sh
  ↓
后续版本执行 update-server.sh
```

运行结构是：

```text
宿主机 systemd
  └─ Go 后端二进制

Docker Compose
  ├─ PostgreSQL
  ├─ Redis
  └─ Nginx
```

也就是说，这组页面默认讲的不是“全容器部署大全”，而是当前仓库最推荐、最稳妥的上线方式。

## 如果你还想继续跟教程章节走

部署篇现在是独立入口，但不影响教程本身继续保留：

- 想走教学主线：回 [第 8 章：部署、升级与复用](/tutorial/chapter-8/)
- 想回到完整章节路线：回 [教程首页](/tutorial/)

## 下一步

- 直接开始部署：读 [部署验证与复用说明](/tutorial/chapter-8/deployment-and-reuse)
- 想先看当前部署文件分工：读 [Docker 部署文件参考](/reference/deploy-artifacts-reference)
- 想知道这一轮部署篇后面怎么重构：读 [部署篇页面规划](/guide/deployment-document-plan)
