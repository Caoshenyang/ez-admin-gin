---
title: 部署篇页面规划
description: "部署篇第一轮盘点结果：明确当前页面职责、重复位置、保留方案，以及未来部署说明站的 canonical 页面结构。"
search: false
---

# 部署篇页面规划

::: warning 维护页说明
这页是部署篇重构前的内部规划页，主要服务维护者统一页面归属。普通读者仍然以现有部署正文为准。
:::

这份规划只做一件事：

> 把当前仓库里和部署相关的教程页、参考页、历史页拉成一张清晰的映射表，确定哪些页面保留、哪些合并、哪些下放到参考手册。

## 为什么先做这张表

当前部署内容已经不少，但它们分散在三个区域：

- `tutorial/chapter-8/*`
- `reference/*`
- `tutorial/archive/*`

这会导致三个问题：

1. 主线入口很多，第一次部署的人不容易判断应该从哪一页开始
2. 同一件事在教程页和参考页之间边界不稳定
3. 有些页面其实已经是“参考资料”，却仍然混在教程主线里

所以部署篇第一步不是马上重写正文，而是先把页面职责定下来。

## 部署篇目标结构

本轮重构后，部署篇希望收成下面这组 canonical 页面：

| 目标页面 | 主要职责 |
| --- | --- |
| `deploy/index.md` | 部署篇首页，告诉读者怎么走完整部署链路 |
| `deploy/preflight-and-baseline.md` | 部署前准备、前置条件、版本与资源基线 |
| `deploy/compose-and-runtime.md` | Compose 结构、systemd、运行拓扑 |
| `deploy/env-and-init.md` | 环境变量、初始化数据、管理员初始化 |
| `deploy/nginx-and-https.md` | 入口层、静态资源、代理、HTTPS |
| `deploy/update-and-rollback.md` | 更新、回滚、回滚分级 |
| `deploy/troubleshooting.md` | 部署排障 FAQ |
| `deploy/operations-faq.md` | 长期运维 FAQ |
| `deploy/reuse-checklist.md` | 新项目复用清单 |

说明：

- 这里的 `deploy/` 是目标信息架构，不代表本轮立刻完成所有文件迁移。
- 当前阶段先用这张表明确“未来应该往哪收”，再逐页改正文。

## 当前页面映射表

### 教程主线页

| 当前页面 | 当前职责 | 当前问题 | 处理决定 | 未来归属 |
| --- | --- | --- | --- | --- |
| `tutorial/chapter-8/index.md` | 第 8 章部署总入口 | 仍然是教程章节入口，不够像成熟说明站首页 | 保留，但后续让它更多承担“教程里的部署章节入口” | `tutorial/chapter-8/index.md` |
| `tutorial/chapter-8/deployment-and-reuse.md` | 打包、上传、首次部署、更新、复用的总流程 | 主题过宽，既讲主线又讲复用，容易继续膨胀 | 拆分主职责，保留为“部署主线总览”的素材来源 | `deploy/index.md`、`deploy/preflight-and-baseline.md` |
| `tutorial/chapter-8/compose-and-service-layout.md` | 服务器运行拓扑、Compose 与 systemd 分工 | 职责清晰，适合保留为部署主线正文 | 基本保留，后续只做结构与表达优化 | `deploy/compose-and-runtime.md` |
| `tutorial/chapter-8/env-and-init-data.md` | 环境变量、迁移、初始化数据、管理员初始化 | 职责清晰，但和参考页存在天然交叉 | 保留为主线页，参数明细下放参考页 | `deploy/env-and-init.md` |
| `tutorial/chapter-8/nginx-and-https.md` | Nginx、代理、HTTPS | 职责清晰，适合作为部署主线正文 | 基本保留，后续补“怎么验证”与“常见误区” | `deploy/nginx-and-https.md` |
| `tutorial/chapter-8/deployment-variants.md` | 不同部署形态的比较说明 | 更像“部署决策补充”，不是首次部署主线 | 合并到部署首页或部署前准备页的“部署变体”小节 | `deploy/preflight-and-baseline.md` |
| `tutorial/chapter-8/update-and-rollback.md` | 更新与回滚主线 | 职责明确，但和回滚分级页需要进一步合并 | 保留为主线页 | `deploy/update-and-rollback.md` |
| `tutorial/chapter-8/rollback-strategy-levels.md` | 回滚分级说明 | 单独成页略碎，适合并入更新回滚主线 | 并回主线页的一个独立章节 | `deploy/update-and-rollback.md` |
| `tutorial/chapter-8/deployment-troubleshooting-faq.md` | 部署排障 FAQ | 职责清晰 | 保留 | `deploy/troubleshooting.md` |
| `tutorial/chapter-8/operations-maintenance-faq.md` | 运维 FAQ | 职责清晰 | 保留 | `deploy/operations-faq.md` |
| `tutorial/chapter-8/project-reuse-checklist.md` | 新项目复用清单 | 职责清晰 | 保留 | `deploy/reuse-checklist.md` |

### 参考手册页

| 当前页面 | 当前职责 | 当前问题 | 处理决定 | 未来归属 |
| --- | --- | --- | --- | --- |
| `reference/environment-variables-reference.md` | 环境变量详细字段说明 | 很适合查阅，但不适合作为第一次部署主线 | 保留在参考手册 | `reference/environment-variables-reference.md` |
| `reference/init-data-reference.md` | 初始化数据与种子说明 | 适合查阅，不适合承担部署主线 | 保留在参考手册 | `reference/init-data-reference.md` |
| `reference/deploy-artifacts-reference.md` | 打包产物与部署文件清单 | 纯参考，适合被部署主线引用 | 保留在参考手册 | `reference/deploy-artifacts-reference.md` |
| `reference/nginx-config-reference.md` | Nginx 配置详情与模板参考 | 适合查配置，不适合作为第一次部署正文 | 保留在参考手册 | `reference/nginx-config-reference.md` |
| `reference/ssh-tunnel-database.md` | SSH 隧道连服务器数据库 | 属于运维补充，不属于首次部署主线 | 暂留参考手册，后续可在运维 FAQ 中挂链接 | `reference/ssh-tunnel-database.md` |

### 历史与归档页

| 当前页面 | 当前职责 | 当前问题 | 处理决定 | 未来归属 |
| --- | --- | --- | --- | --- |
| `tutorial/archive/chapter-7/deployment-and-reuse.md` | 历史部署正文 | 已不是主线，继续暴露会分散入口 | 明确保留为历史页，不再参与主线导航 | `tutorial/archive/*` |
| `tutorial/archive/chapter-7/env-and-init-data.md` | 历史环境与初始化正文 | 已被新主线取代 | 保留归档，不再继续修正文案 | `tutorial/archive/*` |

## 页面职责原则

后续部署篇重写时，统一按下面规则裁内容：

### 主线页负责什么

- 讲顺序
- 讲前置条件
- 讲执行动作
- 讲预期结果
- 讲失败排查入口

### 参考页负责什么

- 讲字段
- 讲模板
- 讲参数
- 讲文件清单

### FAQ 页负责什么

- 回答“遇到问题时先看哪里”
- 不重复完整主线步骤
- 优先用问答式短段落组织

## 当前阶段结论

部署篇已经有足够多的原始素材，不缺内容，真正缺的是：

- 一个明确的总入口
- 稳定的主线与参考边界
- 更像“上线说明站”的阅读顺序

所以这一轮的核心不是继续堆新页面，而是先把这些页面重新排成一条可交付主线。

## 下一步

部署篇下一步直接进入：

1. 重写部署首页
2. 给部署篇补清晰的阅读顺序
3. 再改“环境准备 / Compose / 环境变量”三张主链路页面
