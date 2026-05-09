# Process

本文件现在只记录一件事：`docs/` 这一轮“文档说明站点对齐”的执行计划与状态。

- 更新时间：`2026-05-09`
- 适用范围：`docs/`、`docs/.vitepress/`、文档导航与信息架构
- 当前目标：先把“部署篇”整理成一组成熟、可上线的说明页面，再回头统一整站大纲、导航和内容分层
- 当前约束：代码结构本轮已经定稿，后续以文档对齐和少量内容修正为主，不再发散到大规模代码重构

## 状态标记

- `[x]` 已完成
- `[-]` 进行中
- `[ ]` 未开始
- `[!]` 阻塞 / 需要先确认

## 当前判断

### 已确认事实

- [x] 当前站点已经具备 `guide / tutorial / reference` 三条主线，内容量足够支撑“成熟说明站点”升级
- [x] 部署相关内容目前分散在 `tutorial/chapter-7`、`tutorial/chapter-8` 与 `reference/` 多个页面中，信息存在重复和入口分散问题
- [x] 当前导航更偏“教程编排视角”，还没有完全切到“产品说明站点视角”
- [x] 当前文档读者至少有三类：
  - 新读者：想快速理解项目能做什么、怎么跑起来
  - 落地者：想把项目部署上线、维护、回滚
  - 复用者：想把这套后台底座带到自己的项目里
- [x] 这轮工作不以“新增大量知识点”为目标，而以“重组信息、拉直入口、减少重复、提升交付感”为目标

### 当前核心问题

- [-] 部署内容已经有料，但没有形成一条清晰的“上线主线”
- [-] 站点当前是“教程很强，说明站属性还不够强”
- [-] 同一主题在教程页与参考页之间的边界还不够稳定
- [ ] 还没有一份明确的“整站信息架构升级计划”作为唯一执行基线

## 本轮总目标

把 `EZ Admin Gin` 的文档站点从“内容很多的教程仓库”提升为“可直接交付给用户阅读的成熟说明站点”。

本轮只按下面的优先级推进：

1. 先整理部署篇
2. 再统一整站信息架构
3. 再清理重复页面与导航口径
4. 最后补文档基线验证与维护规则

## 目标站点结构

本轮规划采用下面这套站点分层，后续页面重组都以它为准：

### 1. 开始使用（Getting Started）

目标：

- 让第一次进入项目的人知道这是什么、适合谁、最快怎么跑起来

建议承载内容：

- 项目简介与定位
- 快速启动
- 项目结构
- Java → Go 的理解桥接

现有主要来源：

- `docs/index.md`
- `docs/guide/index.md`
- `docs/guide/project-structure.md`
- `docs/guide/java-to-go-structure.md`
- `docs/guide/enterprise-architecture.md`

### 2. 部署与运维（Deployment）

目标：

- 让准备落地的人能从“环境准备”一路走到“上线、验证、更新、回滚、排障”

建议承载内容：

- 部署首页
- 部署前准备
- Compose 单机部署
- Nginx / HTTPS / 反向代理
- 环境变量与初始化数据
- 更新与回滚
- 常见部署问题
- 运维 FAQ
- 项目复用清单

现有主要来源：

- `docs/tutorial/chapter-8/index.md`
- `docs/tutorial/chapter-8/env-and-init-data.md`
- `docs/tutorial/chapter-8/deployment-and-reuse.md`
- `docs/tutorial/chapter-8/compose-and-service-layout.md`
- `docs/tutorial/chapter-8/nginx-and-https.md`
- `docs/tutorial/chapter-8/deployment-variants.md`
- `docs/tutorial/chapter-8/update-and-rollback.md`
- `docs/tutorial/chapter-8/rollback-strategy-levels.md`
- `docs/tutorial/chapter-8/deployment-troubleshooting-faq.md`
- `docs/tutorial/chapter-8/operations-maintenance-faq.md`
- `docs/tutorial/chapter-8/project-reuse-checklist.md`
- `docs/reference/deploy-artifacts-reference.md`
- `docs/reference/environment-variables-reference.md`
- `docs/reference/init-data-reference.md`
- `docs/reference/nginx-config-reference.md`

### 3. 从零搭建（Tutorial）

目标：

- 保留完整的教学路线，服务“跟着项目一步一步做出来”的读者

建议承载内容：

- 章节导读
- 各章正文
- 教程大纲

原则：

- 教程保持章节化
- 但不再承担“说明站入口”职责

### 4. 参考手册（Reference）

目标：

- 服务查阅，不讲故事

建议承载内容：

- 环境变量参考
- 初始化数据参考
- Deploy 文件参考
- Nginx 配置参考
- 错误码、目录规范、模块规范

原则：

- 参考页只回答“查什么”
- 不承担“第一次怎么做”的主线职责

### 5. 项目说明（Project）

目标：

- 承载路线图、更新日志、文档基线、维护说明

建议承载内容：

- 更新日志
- 路线图
- 文档基线清单
- 当前执行计划页

## 本轮执行策略

### Phase A：部署篇优先重组

目标：

先产出一组可以单独成立的部署说明页面，让站点先具备“可上线交付”的说明能力。

状态：

- [-] 进行中

本阶段要做的事：

1. 先做部署信息盘点
   - 列出教程页、参考页里所有和部署相关的页面
   - 标注每页职责：主线、补充、FAQ、纯参考
   - 标注重复与冲突位置

2. 先确定部署篇目标结构
   - `部署首页`
   - `环境准备`
   - `Compose 与服务运行结构`
   - `环境变量与初始化数据`
   - `Nginx 与 HTTPS`
   - `更新与回滚`
   - `部署排障 FAQ`
   - `运维 FAQ`
   - `项目复用清单`

3. 再确定每页写作口径
   - 主线页：讲顺序、讲验证、讲失败排查
   - 参考页：讲参数、文件、模板，不讲长流程
   - FAQ 页：讲问题和判断，不重复主线正文

4. 最后再落第一版文案更新
   - 先更新部署首页与部署总导航
   - 再逐页改正文
   - 完成后再统一站内链接与侧边栏

完成标准：

- 读者不需要翻很多章，也能找到完整部署路径
- 部署主线和部署参考页边界清楚
- 部署篇可以单独作为“上线说明”阅读

当前进展：

- [x] 已完成部署篇页面盘点，并产出 `docs/guide/deployment-document-plan.md`
- [x] 已完成部署首页与部署导航第一轮重写
- [-] 下一步进入“环境准备 / Compose / 环境变量”三张主链路页面重写

### Phase B：整站导航与大纲重组

目标：

把站点入口从“章节入口优先”调整成“读者任务优先”。

状态：

- [ ] 未开始

本阶段要做的事：

- 重审 `docs/.vitepress/config.mts` 的 `nav` 和 `sidebar`
- 判断是否新增独立的部署导航入口
- 调整首页与 Guide 首页文案，让它们更像说明站入口而不是教程跳板
- 明确每个一级区块的职责：
  - 开始使用
  - 部署与运维
  - 从零搭建
  - 参考手册
  - 项目说明

完成标准：

- 新读者第一次进站，不需要先理解教程章节结构
- “我要上线”“我要学习”“我要查资料”三类需求都能快速定位入口

### Phase C：内容去重与页面归位

目标：

把部署、初始化、复用、FAQ 这些重复话题重新归位，减少站内重复解释。

状态：

- [ ] 未开始

本阶段要做的事：

- 清理教程部署页与参考页的重复段落
- 决定哪些内容保留在教程，哪些上移到部署篇
- 对历史页面做显式说明或归档处理
- 清理“教程里兼做说明页”的页面

完成标准：

- 同一个部署问题不会在 3 个页面里各讲一遍
- 页面职责清楚，读者知道该看主线、FAQ 还是参考

### Phase D：文档质量基线固化

目标：

让后续文档维护有统一标准，不再每次重排都重新定义规则。

状态：

- [ ] 未开始

本阶段要做的事：

- 固定部署页写作模板：
  - 前置条件
  - 操作步骤
  - 期望结果
  - 失败排查
  - 下一步
- 固定 FAQ 页写法
- 固定参考页写法
- 每次实质性修改后执行 `pnpm docs:build`

完成标准：

- 后续补一篇文档时，不需要再临时决定页面写法
- 站点逐渐形成稳定的“说明站语气”和结构感

## 部署篇具体编写计划

这一轮先按下面顺序写，不同时开太多页面。

### Step 1：部署篇盘点与页面职责表

目标：

先做一张映射表，明确现有页面保留、合并、下放还是只做参考。

输出物：

- 一个部署篇页面映射草案

涉及页面：

- `tutorial/chapter-8/*`
- `reference/deploy-artifacts-reference.md`
- `reference/environment-variables-reference.md`
- `reference/init-data-reference.md`
- `reference/nginx-config-reference.md`

状态：

- [x] 已完成

输出结果：

- `docs/guide/deployment-document-plan.md`
- 已明确教程主线页、参考页、历史页三类部署页面的处理决定

### Step 2：部署首页与章节入口重写

目标：

先把“部署篇怎么读”讲清楚。

输出物：

- 部署首页新版
- 明确的页面阅读顺序

优先级：

- 最高

状态：

- [x] 已完成

输出结果：

- 新增 `docs/deploy/index.md` 作为部署与运维独立入口
- `docs/.vitepress/config.mts` 已补部署导航与部署侧边栏
- `docs/index.md` 与 `docs/guide/index.md` 已补部署入口
- `docs/tutorial/chapter-8/index.md` 已补入口分流说明

### Step 3：环境准备 / Compose / 环境变量 三页打通

目标：

把最核心的上线主链路收成一套顺序明确的页面。

输出物：

- 环境准备页
- Compose 与服务运行结构页
- 环境变量与初始化数据页

状态：

- [-] 进行中

当前进展：

- [x] 已完成 `tutorial/chapter-8/deployment-and-reuse.md` 第一轮重写，定位为“部署主线总览页”
- [ ] 还需继续重写 `compose-and-service-layout.md`
- [ ] 还需继续重写 `env-and-init-data.md`

### Step 4：Nginx / HTTPS / 更新回滚

目标：

补齐部署进入正式生产环境前后的关键动作。

输出物：

- Nginx 与 HTTPS 页
- 更新与回滚页
- 回滚分级说明页

状态：

- [ ] 未开始

### Step 5：排障 FAQ / 运维 FAQ / 复用清单

目标：

把“上线后遇到的问题”集中收好，避免散落在教程正文里。

输出物：

- 部署排障 FAQ
- 运维 FAQ
- 项目复用清单

状态：

- [ ] 未开始

## 当前优先级

只按下面顺序执行：

1. [x] 把 `process.md` 切到文档说明站计划
2. [x] 先做部署篇盘点与页面职责表
3. [x] 先改部署首页与部署导航
4. [-] 再推进部署主链路三页
5. [ ] 再推进 Nginx / HTTPS / 更新回滚
6. [ ] 最后补 FAQ 与复用清单
7. [ ] 部署篇收稳后再重构整站信息架构

## 完成定义

当下面条件同时满足时，可以认为这一轮“说明站对齐”的第一阶段完成：

- [ ] 部署篇已经形成独立可读的说明入口
- [ ] 部署主线和部署参考页边界清楚
- [ ] 站点导航能明确区分“开始使用 / 部署 / 教程 / 参考”
- [ ] 重复内容开始明显减少
- [ ] `pnpm docs:build` 持续通过

## 默认约束

- 默认工作路径：`D:\A\ez-admin-gin`
- 默认以 `docs/` 与 `process.md` 为唯一执行基线
- 默认这轮只做文档规划、重写、归位和导航对齐
- 默认不再把 `process.md` 用作代码重构状态页
