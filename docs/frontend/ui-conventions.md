---
title: 前端 UI 规范
description: "统一后台页面中搜索区、列表区、表单区和页面编排的组件入口，避免同类页面出现多套写法。"
---

# 前端 UI 规范

这页用于约束后台页面的视觉与结构写法。新增模块时，优先照这里的组件入口和 class 命名组织页面；已有页面改造时，也按同一套规范逐步收口。

::: tip 核心判断
同一种区域只保留一种主写法：搜索区用 `EzSearchPanel`，列表区用 `EzDataTable`，表单区用 `ez-modal-form / ez-modal-section / ez-modal-footer`。
:::

## 品牌组件风格

后台业务组件统一走“紧凑控制台”视觉：中性面、清晰层级、稳定网格和低噪声边界。EZ 的品牌感不依赖重复色条、渐变或网格背景，而是通过更收紧的表单节奏、更一致的工具栏密度、更可扫描的表格结构和克制的主色使用建立。

| 区域 | 视觉入口 | 使用方式 |
| --- | --- | --- |
| 页头 | `PageHeader` / `EzPageHeader` | 标题、说明和操作区保持同一行节奏，不手写额外标题块 |
| 搜索区 | `EzSearchPanel` | 只放筛选字段和查询动作，避免重新包一层卡片 |
| 列表区 | `EzDataTable` | 普通表格、顶部工具条、刷新按钮和分页底栏由公共组件负责 |
| 弹窗表单 | `FormModalHeader` + `ez-modal-*` | 使用分组和网格收紧节奏，不在业务组件里重写弹窗头 |
| 抽屉 | `NDrawer` + `ez-drawer` | 只用于通知、日志、审计详情等只读侧栏 |
| 空状态 / 指标卡 | `EzEmptyState` / `EzMetricCard` | 使用共享组件，不再手写普通白卡 |

::: warning 不要复制装饰样式
业务页面不要再手写渐变、网格背景、彩色顶线、多层边框或大阴影。品牌表达集中在公共组件的间距、字号、边界、表格密度和主操作状态里，局部页面只补充必要布局。
:::

## 页面骨架

业务页面默认采用一屏后台布局，页面内部滚动由内容区承接，不让浏览器出现默认滚动条。

```vue
<template>
  <main class="admin-page">
    <section class="admin-page-section">
      <PageHeader title="用户管理" description="维护后台账号、启停状态和角色绑定。">
        <template #actions>
          <NButton type="primary">新增</NButton>
        </template>
      </PageHeader>

      <UserFilterBar />
      <UserTable />
    </section>
  </main>
</template>
```

| 区域 | 固定入口 | 说明 |
| --- | --- | --- |
| 页面容器 | `admin-page` | 负责承接一屏高度和纵向布局 |
| 页面内容 | `admin-page-section` | 负责统一页头、搜索区、列表区之间的间距 |
| 页头 | `PageHeader` / `EzPageHeader` | 标题、描述和主操作按钮 |
| 搜索区 | `EzSearchPanel` | 筛选项和查询/重置按钮 |
| 列表区 | `EzDataTable` | 表格、顶部统计条、底部分页 |

## 弹窗与抽屉

新增、编辑、分配、上传这类会改变业务数据的主操作，统一使用 `NModal preset="card"`。抽屉只用于旁路信息，例如通知列表、日志详情、审计明细这类不打断当前页面上下文的查看场景。

| 场景 | 固定入口 | 当前示例 |
| --- | --- | --- |
| 新增 / 编辑基础资料 | `NModal` + `FormModalHeader` + `ez-modal-form` | 用户、角色、部门、岗位、菜单、配置、公告、字典 |
| 分配关系 | `NModal` + `FormModalHeader` | 用户分配角色 |
| 上传 / 编辑附件 | `NModal` + `ez-modal-form` | 附件上传、附件编辑 |
| 无附加字段的直接动作 | 页头按钮或 `NUpload` 包裹按钮 | 文件管理的直接上传 |
| 详情查看 | `NDrawer` | 操作日志详情 |
| 全局通知 | `NDrawer` | 顶部通知抽屉 |

::: warning 不再新增抽屉表单
不要为新增、编辑、上传再创建 `DrawerForm` 一类组件。表单字段再多，也优先通过分组、两列网格和弹窗宽度处理；只有只读详情或全局侧栏信息才使用抽屉。
:::

## 搜索区

搜索区只做“条件输入 + 查询动作”，不要混入新增、导出、批量删除这类列表操作。列表操作放到表格顶部工具条。

```vue
<template>
  <EzSearchPanel>
    <NInput v-model:value="query.keyword" clearable placeholder="名称 / 编码" class="w-56" />
    <NSelect v-model:value="query.status" :options="STATUS_FILTER_OPTIONS" class="w-36" />

    <template #actions>
      <EzActionButton kind="search" label="查询" type="primary" @click="emit('search')" />
      <EzActionButton kind="reset" label="重置" @click="emit('reset')" />
    </template>
  </EzSearchPanel>
</template>
```

搜索区约定：

- 关键词输入宽度优先用 `w-56` / `w-64`。
- 状态、类型、扩展名这类短选择器优先用 `w-36` / `w-40`。
- 查询按钮始终放在 actions 插槽第一位，重置按钮第二位，使用 `EzActionButton` 的图标 + 短文字样式。
- 支持回车查询的字段绑定 `@keyup.enter="emit('search')"`。
- 不再新增手写 `<div class="ez-toolbar">` 筛选条。

## 列表区

列表区统一使用 `EzDataTable`。普通列表直接传 `columns / data / loading / page / pageSize / total`；需要勾选、树表、额外事件或自定义空态时，用 `body` 插槽接管内部 `NDataTable`，但仍复用统一工具条、刷新、密度、列设置和分页外壳。

```vue
<template>
  <EzDataTable
    remote
    :columns="columns"
    :data="items"
    :loading="loading"
    :page="query.page"
    :page-size="query.page_size"
    :row-key="(row) => row.id"
    :total="total"
    @page-change="emit('pageChange', $event)"
    @page-size-change="emit('pageSizeChange', $event)"
    @refresh="emit('refresh')"
  />
</template>
```

树表、勾选表或自定义空态接管表体时，必须继续使用 `body` 插槽提供的 `tableColumns / tableScrollX / tableSize`，否则列显隐、横向滚动和密度设置不会生效：

```vue
<template>
  <EzDataTable :columns="columns" :data="items" :loading="loading">
    <template #toolbarActions>
      <EzActionButton kind="expand" label="展开全部" secondary size="small" @click="emit('expandAll')" />
    </template>

    <template #body="{ tableColumns, tableScrollX, tableSize }">
      <NDataTable
        class="ez-table-fill-table"
        :columns="tableColumns"
        :data="items"
        :loading="loading"
        :scroll-x="tableScrollX"
        :size="tableSize"
        :pagination="false"
        :bordered="false"
        flex-height
      />
    </template>
  </EzDataTable>
</template>
```

列表区约定：

- 表格不直接使用 Naive UI 分页，普通列表由 `EzDataTable` 统一渲染分页。
- 表格容器不再新增 `NCard class="ez-table-card"` 写法。
- 顶部统计条只放统计、刷新、展开/收起、批量动作，不放搜索条件。
- 刷新、密度、列设置由 `EzDataTable` 统一提供，并作为右侧图标型视图控制组展示；业务侧扩展按钮放入 `toolbarActions`，优先使用 `EzActionButton` 的图标 + 短文字样式，不要覆盖整段 `actions`。
- `EzTableCard` 只作为底层外壳使用，业务列表不要再手写 `EzTableCard + TableStatsBar + NDataTable + NPagination`。
- 树表、勾选表和带自定义空态的表格可以通过 `EzDataTable` 的 `body` 插槽接管 `NDataTable`。
- 表头统一保持 13px / 700 的信息层级，避免比正文弱太多。
- 长表格必须配置稳定列宽，避免操作列挤压正文列。
- `EzDataTable` 会根据可见列宽自动计算 `scroll-x`；业务自定义表体必须透传 `tableScrollX`，确保列展示不全时有横向滚动条。
- 操作列统一固定在右侧，按按钮数量配置稳定宽度，并使用 `ez-row-actions` 承接统一按钮样式。
- 操作列按钮优先使用 `EzActionButton icon-only size="tiny" secondary`，通过 tooltip / `aria-label` 暴露完整动作名称；按钮数量超过四个时再考虑下拉菜单收纳。

## 表单区

表单区优先使用 Naive UI 表单组件，但布局和间距必须走项目公共类。

```vue
<template>
  <NForm class="ez-modal-form" label-placement="left" label-width="76">
    <section class="ez-modal-section ez-modal-section--soft">
      <div class="ez-modal-section__head">
        <h3>基础信息</h3>
        <p>编码类字段创建后应保持稳定。</p>
      </div>

      <div class="ez-form-grid ez-form-grid--2">
        <NFormItem label="名称" path="name">
          <NInput placeholder="请输入名称" />
        </NFormItem>
        <NFormItem label="状态" path="status">
          <NSelect />
        </NFormItem>
      </div>
    </section>

    <div class="ez-modal-footer">
      <NButton>取消</NButton>
      <NButton type="primary">保存</NButton>
    </div>
  </NForm>
</template>
```

表单区约定：

- 弹窗表单使用 `NModal preset="card"` 和 `FormModalHeader` 作为外壳。
- 新增、编辑、分配、上传统一使用 `NModal`，不使用抽屉表单。
- 如果动作没有任何附加字段，例如直接选择文件上传，可以放在页头主按钮中，不额外套表单弹窗。
- 表单主体使用 `ez-modal-form`，分组使用 `ez-modal-section`。
- 两列字段使用 `ez-form-grid ez-form-grid--2`，三列以内才允许同屏并排。
- 长文本、备注、说明字段单独占满一行。
- 底部按钮统一使用 `ez-modal-footer`，主按钮在右侧。
- 表单 label 默认用 `label-placement="left"`，字段少、移动端优先或登录页可用 `top`。

## 新增页面检查表

新增或重构页面前，先按这张表自查：

| 检查项 | 通过标准 |
| --- | --- |
| 页面骨架 | `admin-page > admin-page-section > PageHeader + Search + Table` |
| 搜索区 | 使用 `EzSearchPanel`，没有手写 `.ez-toolbar` |
| 列表区 | 使用 `EzDataTable`，特殊表格通过 `body` 插槽接管表体 |
| 表单区 | 使用 `ez-modal-form`、`ez-modal-section`、`ez-modal-footer` |
| 宽度 | 输入、选择器、表格 `scroll-x` 有稳定尺寸 |
| 行为 | 查询、重置、刷新、分页事件由组件 emit 给页面/composable |

::: warning 不要把规范写成“局部特例”
如果一个页面需要特殊布局，先判断它是不是新的公共模式。只有确实无法复用时，才在局部 scoped style 里做补充，并避免覆盖全局组件样式。
:::
