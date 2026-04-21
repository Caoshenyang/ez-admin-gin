---
title: Markdown 扩展示例
description: "集中说明 VitePress 文档中常用的 Markdown 扩展能力，包括提示块、代码组、行高亮、代码导入、数学公式和图片懒加载。"
---

# Markdown 扩展示例

这一页是 `EZ Admin Gin` 文档写作的语法速查。新增教程页、参考页或 FAQ 时，优先从这里选择合适的语法。

::: tip 使用建议
不要为了“看起来高级”而堆扩展。每一个扩展都应该降低读者理解成本。
:::

## 全局能力

VitePress 常用能力包括：

- 代码块行号
- 行高亮
- 代码组
- GitHub Alerts
- VitePress 自定义容器
- `[[toc]]` 目录
- 代码片段导入
- Markdown include
- 数学公式
- Markdown 图片懒加载

具体是否启用，最终以 `docs/.vitepress/config.mts` 为准。

## GitHub 风格警报

```md
> [!TIP]
> 适合放推荐做法或第一次阅读时应该抓住的重点。

> [!NOTE]
> 适合放补充说明、实现背景或章节之间的关系。

> [!WARNING]
> 适合放容易踩坑的地方，例如路径、环境变量、大小写和生产配置差异。
```

渲染效果：

> [!TIP]
> 适合放推荐做法或第一次阅读时应该抓住的重点。

> [!NOTE]
> 适合放补充说明、实现背景或章节之间的关系。

> [!WARNING]
> 适合放容易踩坑的地方，例如路径、环境变量、大小写和生产配置差异。

## VitePress 容器

```md
::: tip 推荐做法
这里写推荐路径。
:::

::: warning 注意
这里写高频风险。
:::

::: details 可选深入
这里写不影响主流程的补充内容。
:::
```

::: tip 推荐做法
主流程中的关键建议可以放在 `tip` 里。
:::

::: warning 注意
高频踩坑、版本差异和生产风险可以放在 `warning` 里。
:::

::: details 可选深入
不影响第一次跑通的背景知识可以放在 `details` 里。
:::

## 代码块增强

### 行高亮

````md
```ts{2,4-6}
export function createAuthHeader(token: string) {
  const normalized = token.trim()
  return {
    Authorization: `Bearer ${normalized}`,
    'Content-Type': 'application/json',
    Accept: 'application/json'
  }
}
```
````

效果：

```ts{2,4-6}
export function createAuthHeader(token: string) {
  const normalized = token.trim()
  return {
    Authorization: `Bearer ${normalized}`,
    'Content-Type': 'application/json',
    Accept: 'application/json'
  }
}
```

### 聚焦、差异、警告和错误

````md
```ts
const columns = [
  { title: '用户名', key: 'username' }, // [!code ++]
  { title: '邮箱', key: 'email' }, // [!code focus]
  { title: '手机号', key: 'phone' } // [!code --]
]

submitForm(payload) // [!code warning]
unsafeEval(payload) // [!code error]
```
````

效果：

```ts
const columns = [
  { title: '用户名', key: 'username' }, // [!code ++]
  { title: '邮箱', key: 'email' }, // [!code focus]
  { title: '手机号', key: 'phone' } // [!code --]
]

submitForm(payload) // [!code warning]
unsafeEval(payload) // [!code error]
```

说明：

- `focus`：突出当前最重要的代码。
- `++ / --`：表达新增或移除。
- `warning / error`：强调潜在风险和错误用法。

## 代码组

多平台、多包管理器、多调用方式时使用 `code-group`。

````md
::: code-group

```bash [pnpm]
pnpm install
pnpm docs:dev
```

```bash [npm]
npm install
npm run docs:dev
```

```bash [yarn]
yarn
yarn docs:dev
```

:::
````

效果：

::: code-group

```bash [pnpm]
pnpm install
pnpm docs:dev
```

```bash [npm]
npm install
npm run docs:dev
```

```bash [yarn]
yarn
yarn docs:dev
```

:::

## 目录

正文内目录使用：

```md
[[toc]]
```

默认建议：

- 右侧大纲正常显示时，不额外加正文目录。
- FAQ 或特殊汇总页可以加正文目录。
- 如果配置关闭了右侧大纲，长文可以加正文目录。

## 导入代码片段

导入完整文件：

```md
<<< @/snippets/example.http
```

导入指定区域：

```md
<<< @/snippets/router.ts#auth-routes
```

使用要求：

- 被导入文件必须真实存在。
- 片段路径应尽量稳定。
- 不要在文档里留下只有示例意义、实际会构建失败的导入。

## 包含 Markdown 文件

重复出现的说明可以抽成独立 Markdown 片段，再用 include 复用：

```md
<!--@include: ./_includes/doc-review-checklist.md-->
```

适合 include 的内容：

- 发布前检查表
- 通用风险提醒
- 环境说明
- 多页共用的术语解释

## 数学公式

只有当公式能帮助读者更快理解时才使用。

行内公式：

```md
缓存命中率可以写成 $hit\ rate = \frac{hit}{hit + miss}$。
```

块级公式：

```md
$$
T_{response} = T_{network} + T_{service} + T_{database}
$$
```

效果：

缓存命中率可以写成 $hit\ rate = \frac{hit}{hit + miss}$。

$$
T_{response} = T_{network} + T_{service} + T_{database}
$$

## 图片

Markdown 图片语法：

```md
![系统架构图](/images/architecture.svg)
```

使用建议：

- 图片必须服务理解，不要只是装饰。
- 图下方最好说明读者应该看哪里。
- 同一页不要连续堆太多大图。

## 推荐约定

::: info 写作优先级
教程型文档优先使用 `tip / warning / details / code-group` 这些高频、稳定、低学习成本的语法；只有在确实能提升理解效率时，再引入更进阶的写法。
:::
