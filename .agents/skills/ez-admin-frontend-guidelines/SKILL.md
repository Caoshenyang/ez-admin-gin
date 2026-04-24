---
name: ez-admin-frontend-guidelines
description: 当需要在本仓库中编写、修改、审阅或同步 Vue 前端代码、Naive UI 管理台页面、Tailwind CSS 4 样式、第五章前端教程文档时使用。尤其适用于登录页、后台布局、工作标签、菜单、空状态页、表单紧凑化、一屏高度和前端规范一致性维护。
---

# EZ Admin 前端编写规范

本 skill 用于维护本仓库 `admin` 前端和第五章教程文档的一致性。目标是让页面贴近原型、代码可维护、组件优先复用 Naive UI，并且保持一屏后台应用的交互体验。

## 技术边界

- Vue 文件默认使用 Vue 3 Composition API、`<script setup lang="ts">` 和 TypeScript。
- 后台组件优先使用 Naive UI：`NLayout`、`NLayoutSider`、`NLayoutHeader`、`NLayoutContent`、`NMenu`、`NForm`、`NInput`、`NButton`、`NCard`、`NEmpty`、`NDropdown`。
- Tailwind CSS 4 负责页面尺寸约束、栅格补充、间距、背景、颜色微调和响应式，不替代成熟 Naive UI 组件。
- 不为了“全 Tailwind”重复实现表单、按钮、菜单、布局、标签、下拉、空状态等已有组件。
- 自定义原生元素只用于品牌块、验证码占位、少量装饰和 Naive UI 没有直接表达的细节。

## 页面布局

- 后台主框架默认限制在一屏：根布局使用 `h-screen`，浏览器级滚动关闭。
- `html`、`body`、`#app` 使用上一行注释写法：

```css
/* 关闭浏览器默认滚动条 */
overflow: hidden;
```

- 页面超出一屏时，优先让业务内容区内部滚动，不让浏览器出现默认滚动条。
- 使用 `NLayout` 体系承接后台壳子，不手写大段 `aside/header/main` 来替代组件库能力。
- 左侧菜单使用 `NMenu`，菜单数据先抽成配置数组，后续便于替换成后端动态菜单。
- 工作标签以原型视觉为准。`NTabs` 视觉不合适时，可以保留轻量自实现的小标签按钮，但要继续使用 Naive UI 承接布局、菜单和操作按钮。

## 视觉与交互

- 后台系统要偏工作台气质：紧凑、清晰、可扫描，避免营销页式的大 hero、大卡片堆叠。
- 登录页、后台布局、空状态页可以用 Tailwind 做页面氛围，但组件交互仍优先走 Naive UI。
- 表单要控制纵向节奏。Naive UI 表单间距过大时，优先用局部 class 调整组件 CSS 变量，不全局污染。
- 可点击元素必须有小手反馈。全局可以覆盖：

```css
button:not(:disabled),
[role='button']:not([aria-disabled='true']),
.n-button:not(.n-button--disabled) {
  cursor: pointer;
}
```

- 验证码、忘记密码等未接后端能力时，可以保留占位 UI，但文档和页面必须明确“后续接入”。
- 当前阶段默认管理员账号为 `admin / EzAdmin@123456`，不要写成原型里的 `admin / 123456`。

## Tailwind 使用

- 能用 Tailwind 标准 spacing 表达时，不使用任意值，例如优先 `gap-4.5` 而不是 `gap-[18px]`。
- 颜色、复杂 grid、阴影、精确字号这类没有清晰等价类时，可以保留任意值。
- 不拼接动态 Tailwind 类名；需要动态样式时使用对象 class、明确分支或行内 style。
- 不把大量设计细节写进全局 CSS。页面局部修正放到 SFC scoped style。

## TypeScript 与 Vue

- 不留下未使用导入。模板使用如果被工具误判，优先换成普通 HTML 或更直接的 Naive UI 组件写法。
- 数组读取必须处理 `undefined`，即使运行时已有兜底，也要让 TypeScript 明确知道。
- 路由跳转使用 `void router.push(...)` / `void router.replace(...)`，避免未处理 Promise 提示。
- 存储、认证、接口封装放到 `utils` / `api` / `types`，页面只承接交互和展示。

## 文档同步

- 修改前端实现时，同步更新第五章对应教程代码块，避免教程和真实代码分叉。
- 文档里要写清楚当前能力边界：哪些已接真实接口，哪些只是占位。
- 代码示例要能直接复制到项目中使用，避免残留未使用导入、过期类名或和真实文件不一致的逻辑。
- 如果用户明确说“只写文档”，只改文档；如果用户要求页面视觉或验证问题，代码和文档一起维护。

## 验证清单

前端改动完成后，优先运行：

```bash
cd admin
pnpm exec oxlint .
pnpm exec vue-tsc --noEmit
```

如果 `pnpm exec` 解析不到二进制，可以直接调用 `admin/node_modules/.bin` 下的对应命令。视觉验证由用户负责时，在最终回复里明确说明未跑浏览器验证。
