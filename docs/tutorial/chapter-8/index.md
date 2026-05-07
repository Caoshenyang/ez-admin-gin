---
title: 第 8 章：部署、升级与复用
description: "围绕当前已有部署正文，收口环境变量、初始化数据、部署验证、升级发布和新项目复用。"
---

# 第 8 章：部署、升级与复用

::: tip 🎯 这一章会做成什么
这一章会把项目从“本地能跑”推进到“可以部署、可以升级、可以复制到新项目继续用”。
:::

## 本章聚焦什么

企业级完整版主线的最后一章，不只讲 Docker 和 Nginx，还会把部署产物、环境变量、升级迁移、回滚说明和新项目复用清单一起收口。

::: info 边界提醒
本章统一承担部署、升级、回滚、排障和复用相关内容。旧章节位置下的部署页只保留兼容跳转，不再承载主线正文。
:::

## 本章正文

当前已经可以直接阅读的部署正文有：

- [环境变量与初始化数据](./env-and-init-data)
- [部署验证与复用说明](./deployment-and-reuse)
- [Compose 与服务运行结构](./compose-and-service-layout)
- [Nginx 与 HTTPS 入口层](./nginx-and-https)
- [部署变体说明](./deployment-variants)
- [更新与回滚策略](./update-and-rollback)
- [回滚分级策略](./rollback-strategy-levels)
- [部署排障 FAQ](./deployment-troubleshooting-faq)
- [长期运维 FAQ](./operations-maintenance-faq)
- [新项目复用清单](./project-reuse-checklist)

## 本章完成后的判断标准

这一章收稳后，至少应该能回答下面几个问题：

1. 部署前哪些环境变量和初始化数据必须准备好
2. 基础设施为什么交给 `docker compose`，后端为什么交给 `systemd`
3. Nginx 和 HTTPS 应该放在什么位置，反向代理链路怎么验证
4. 面对不同交付边界，应该选哪一种部署形态
5. 日常发版如何更新，出现异常后如何回滚
6. 不同故障级别下，回滚动作应该做到什么程度
7. 常见故障出现后，应该先看哪一层、先跑哪些检查命令
8. 部署上线后，平时应该固定看哪些指标、日志和备份动作
9. 换一个新项目继续复用这套底座时，哪些内容可以直接复制，哪些内容必须重配

## 怎么继续读

- 如果你正在跟主线推进，先完成 [第 7 章：前端企业级管理台](../chapter-7/)
- 想先看总体大纲，可以回到 [教程大纲](../curriculum)

## 本章小节

- [环境变量与初始化数据](./env-and-init-data)
- [部署验证与复用说明](./deployment-and-reuse)
- [Compose 与服务运行结构](./compose-and-service-layout)
- [Nginx 与 HTTPS 入口层](./nginx-and-https)
- [部署变体说明](./deployment-variants)
- [更新与回滚策略](./update-and-rollback)
- [回滚分级策略](./rollback-strategy-levels)
- [部署排障 FAQ](./deployment-troubleshooting-faq)
- [长期运维 FAQ](./operations-maintenance-faq)
- [新项目复用清单](./project-reuse-checklist)
