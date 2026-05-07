---
title: 系统模块页面模式
description: "围绕当前真实的 pages/system/*，讲清后台管理页的几种固定模式、权限消费方式和验证方法。"
---

# 系统模块页面模式

第 7 章前面三页已经把前端运行时主线讲清楚了：

- 登录态怎么保存
- 动态菜单怎么生成
- 后台壳子怎么承接页面

接下来就轮到最贴近“日常开发”的部分了：

> 一个真实的系统页面，在当前后台里通常应该长成什么样。

这一页不再逐个讲某一个具体页面怎么复制代码，而是先把当前 `admin/src/pages/system/*` 已经稳定下来的页面模式收拢出来。

::: tip 🎯 本节目标
读完后，你应该能判断一个新页面更像哪一类，并知道它大概应该复用哪些前端结构，而不是每次都从空白 `.vue` 文件开始试。
:::

## 先看当前真实页面分布

现在系统页主要集中在：

```text
admin/src/pages/system/
├─ UserView.vue
├─ RoleView.vue
├─ MenuView.vue
├─ ConfigView.vue
├─ FileView.vue
├─ OperationLogView.vue
├─ LoginLogView.vue
├─ NoticeView.vue
└─ HealthView.vue
```

它们虽然业务不同，但并不是八套完全独立的写法。按当前项目的实际形态，大致可以归成下面四类：

| 页面类型 | 代表页面 | 典型特征 |
| --- | --- | --- |
| 列表型管理页 | `UserView`、`ConfigView`、`NoticeView` | 搜索 + 表格 + 分页 + 弹窗表单 |
| 权限编排页 | `RoleView`、`MenuView` | 除基础 CRUD 外，还要处理树、权限节点或层级结构 |
| 资源与审计页 | `FileView`、`OperationLogView`、`LoginLogView` | 更强调只读浏览、上传、详情抽屉等操作 |
| 状态看板页 | `HealthView` | 更像仪表页，而不是标准数据管理页 |

也就是说，第 7 章讲“系统模块页面”，重点不是记住每一个页面细节，而是先识别：

> 这个页面属于哪一种交互模式。

## 所有系统页几乎都共享同一套骨架

不管页面属于哪一类，当前系统页都在复用一条很稳定的实现骨架：

```text
types/<module>.ts
  ↓
api/<module>.ts
  ↓
pages/system/<Module>View.vue
  ↓
buttonPermissionCodes
  ↓
NDataTable / NModal / NPagination / NAlert
```

换句话说，一个系统页面通常不会把：

- 数据结构
- 请求逻辑
- 页面状态
- 按钮权限

混在一起写，而是让这几层先各自站稳。

## 当前页面里最常见的六个固定部件

只要你打开 `UserView.vue`、`ConfigView.vue` 或 `NoticeView.vue`，很快就会看到一批反复出现的部件：

1. `query`：保存筛选条件和分页参数
2. `loadXxx()`：拉取列表数据
3. `formVisible / formMode / formModel`：控制新增和编辑弹窗
4. `successText`：操作成功后的轻提示
5. `canUse(code)`：按钮权限判定
6. `columns`：表格列定义和操作按钮渲染

这几个部件反复出现，本质上说明当前后台已经形成了统一的管理页心智：

> 搜索条件是状态，表格是主画布，弹窗负责写操作，按钮显隐服从权限码。

## 第一类：列表型管理页

这一类页面是当前后台最常见的形态，代表文件有：

- `admin/src/pages/system/UserView.vue`
- `admin/src/pages/system/ConfigView.vue`
- `admin/src/pages/system/NoticeView.vue`

它们的共同点是：

- 顶部有标题和页面说明
- 中间有搜索区
- 主体是 `NDataTable`
- 底部有 `NPagination`
- 写操作通过 `NModal` + `NForm` 完成

### 为什么这是当前后台的默认页面形态

因为它刚好覆盖了大多数后台首版交付最常见的动作：

- 查
- 新增
- 编辑
- 启停

对中小型后台来说，这种模式比“列表页 + 单独详情页 + 单独编辑页”的多路由拆分更轻，也更稳定。

### 当前三个代表页各自多了什么

`UserView.vue` 在这套骨架上又往前走了一步：

- 支持角色筛选
- 支持岗位展示
- 支持多选行
- 支持“分配角色”这类附加动作

`ConfigView.vue` 则更像最标准的基础资源页：

- 查询字段简单
- 表格列固定
- 编辑动作清晰
- 配置键在编辑态下保持只读

`NoticeView.vue` 是更轻量的内容管理页：

- 表单字段更少
- 资源本身更接近“标题 + 内容 + 状态”
- 很适合作为新手理解标准管理页骨架的第一站

## 第二类：权限编排页

这一类页面的代表是：

- `admin/src/pages/system/RoleView.vue`
- `admin/src/pages/system/MenuView.vue`

它们和普通列表页最大的不同在于：

> 页面不只是维护一张资源表，还要管理一套结构化权限关系。

### `RoleView.vue` 为什么不是普通表格页

角色页当前更接近一个“小型权限工作台”：

- 左侧是角色列表
- 右侧是当前角色详情与权限面板
- 权限面板又拆成菜单、按钮、接口三个维度
- 菜单和按钮依赖 `NTree`
- 接口权限依赖独立的权限行集合

所以角色页虽然也有新增和编辑弹窗，但真正复杂的地方其实在：

- 选中哪个角色
- 这个角色当前有哪些菜单和按钮
- 这个角色当前有哪些接口权限

它不是“资源 CRUD 页”，而是“权限编排页”。

### `MenuView.vue` 为什么要用树形表格

菜单页要同时表达：

- 目录
- 页面菜单
- 按钮节点

这天然就是层级结构，所以当前实现选择了：

- 树形表格
- 展开 / 收起
- 父节点选择
- 子级新增

这也是为什么菜单页的核心难点不在输入框数量，而在：

> 你是否能把一个后台页面入口、一个按钮节点和它们的层级关系同时表达清楚。

## 第三类：资源与审计页

这一类页面更强调“浏览”和“辅助动作”，代表文件有：

- `admin/src/pages/system/FileView.vue`
- `admin/src/pages/system/OperationLogView.vue`
- `admin/src/pages/system/LoginLogView.vue`

### `FileView.vue`：上传是第一动作

文件页虽然也有表格和分页，但它的第一动作不是“新增表单”，而是：

- 上传文件
- 查看文件类型
- 复制 URL

因此它和普通列表页最大的区别是：

- 入口按钮是 `NUpload`
- 结果反馈强调上传成功
- 列表更关注文件类型、大小、链接等展示信息

### `OperationLogView.vue`：详情抽屉比编辑表单更重要

操作日志页是典型的只读审计页面：

- 没有新增
- 没有编辑
- 没有删除

它最重要的交互不是表单，而是：

- 风险级别着色
- 行级背景提示
- “详情”抽屉展开完整上下文

这说明当前后台并没有把所有页面都硬塞进 CRUD 模型，而是会按资源特性调整交互重点。

### `LoginLogView.vue`：最纯粹的查询页

登录日志页更简单：

- 只有筛选
- 只有表格
- 只有分页

它几乎可以看作“最小查询页模板”。当你未来要接一个只读列表资源时，它会比角色页、菜单页更值得参考。

## 第四类：状态看板页

当前 `HealthView.vue` 不属于标准管理页，它更接近系统状态看板。

这种页面的特点通常是：

- 关注运行状态和摘要信息
- 不强调分页和表单
- 更像“概览页”而不是“资源列表页”

这也提醒我们：

> 第 7 章讲的页面模式虽然有主流骨架，但不是所有页面都必须长成同一张表格。

## 按钮权限在这些页面里是怎么落地的

当前大多数系统页都会写一个很轻的辅助函数：

```ts
function canUse(code: string) {
  return buttonPermissionCodes.value.includes(code)
}
```

再把它用在：

- 新增按钮
- 编辑按钮
- 状态切换按钮
- 上传按钮
- 分配角色按钮

这一层很关键，因为它让页面组件只关心：

> 我需要哪个权限码。

而不需要在页面里自己推导“当前是不是管理员”“当前是不是某个角色”。

## 当前页面模式和后端模块结构为什么能一一对上

前面第 6 章已经把后端模块结构收口到了：

```text
dto / entity / repository / service / handler / routes
```

前端这边对应也逐渐稳定成：

```text
types / api / pages/system / dynamic-menu
```

这两套结构能对上的好处是：

- `types/*` 对应后端 DTO / Response
- `api/*` 对应模块 handler 暴露的资源接口
- `pages/system/*` 对应模块真正的后台入口
- `dynamic-menu.ts` 对应后端菜单 `component`

所以一个模块前后端都成熟后，读代码时会更容易形成稳定映射关系。

## 看完总览后，应该继续进入哪些详细页

这一页负责先帮你识别页面模式。确认类型后，最自然的继续阅读顺序是：

- 想看标准管理页：进入 [用户管理页实现要点](./user-management-page-detail)
- 想看权限编排页：进入 [角色与菜单页实现要点](./role-and-menu-page-detail)
- 想看资源页和上传页：进入 [配置与文件页实现要点](./config-and-file-page-detail)
- 想看审计查询页：进入 [日志查询页实现要点](./audit-log-pages)

## 一份页面选型速查表

如果你准备新增一个系统页，可以先用这张表判断该参考哪类现有页面：

| 你的页面需求 | 优先参考 |
| --- | --- |
| 有查询、分页、新增、编辑、启停 | `UserView` / `ConfigView` / `NoticeView` |
| 有树结构、层级关系或权限节点 | `RoleView` / `MenuView` |
| 需要上传文件或复制资源链接 | `FileView` |
| 只读审计、支持详情查看 | `OperationLogView` |
| 只读查询、结构简单 | `LoginLogView` |
| 展示系统状态摘要 | `HealthView` |

## 本页读完后，建议怎么继续

推荐按下面顺序继续：

1. 先回头结合 [后台壳子、动态菜单与按钮权限](./admin-shell-and-dynamic-menu)，把页面入口和按钮权限再串一遍
2. 再按需要进入对应的详细案例页，看某一类页面的完整实现
3. 最后进入 [第 6 章：核心系统模块](../chapter-6/)，把这些页面模式放进完整模块接入流程里
