---
title: Vue 管理台项目初始化
description: "初始化 admin 子项目，为 Vue 3 后台管理台准备基础工程。"
---

# Vue 管理台项目初始化

上一节已经让后端服务可以通过 `/health` 验证。现在初始化 `admin/` 子项目，为后续管理台页面、路由和状态管理准备基础工程。

::: tip 🎯 本节目标
用最新稳定版 Vue 工程脚手架初始化 `admin/`，并确认管理台开发服务可以正常启动。
:::

本节会启用这些基础能力：

| 能力 | 用途 |
| --- | --- |
| TypeScript | 给前端代码提供类型约束 |
| Vue Router | 后续管理台页面需要路由 |
| Pinia | 后续登录状态、用户信息、菜单状态会放进全局状态 |
| ESLint | 保持代码质量 |
| Prettier / Oxfmt | 保持格式统一 |

::: details 当前版本参考
截止 2026-04-21，npm registry 查询到的当前版本是：

- `vue@3.5.32`
- `vite@8.0.9`
- `create-vue@3.22.3`

本节使用 `pnpm create vue@latest`，实际安装版本以你执行命令时 registry 解析到的最新稳定版为准。
:::

## 🛠️ 确认 Node 与 pnpm

Vite 8 要求 Node.js `20.19+` 或 `22.12+`。先确认本机版本：

```bash
node -v
pnpm -v
```

如果没有安装 Node.js，先到官方页面安装：

- Node.js 官网：[https://nodejs.org](https://nodejs.org)

如果已经安装 Node.js，但没有 `pnpm`，可以先启用 Corepack：

```bash
corepack enable
corepack prepare pnpm@latest --activate
```

::: warning ⚠️ Node 版本过低
如果 `node -v` 显示低于 `v20.19.0`，后面启动 Vite 可能会失败。先升级 Node，再继续初始化管理台项目。
:::

## 🛠️ 初始化 admin 项目

确认当前位于项目根目录，然后先移除 `admin/` 里的占位文件：

::: warning ⚠️ 先确认当前位置
下面的初始化命令必须在项目根目录执行，也就是能看到 `server`、`admin`、`docs`、`deploy`、`scripts` 的那一层。不要进入 `admin/` 目录后再执行。
:::

::: code-group

```powershell [Windows PowerShell]
Remove-Item .\admin\.gitkeep -ErrorAction SilentlyContinue
```

```bash [macOS / Linux]
rm -f admin/.gitkeep
```

:::

按照 Vue 官方 Quick Start 的方式，使用最新脚手架初始化 `admin/` 目录：

```bash
pnpm create vue@latest admin
```

脚手架会进入交互式流程。按下面的方向选择：

| 提示项 | 选择 |
| --- | --- |
| 是否使用 TypeScript 语法？ | Yes |
| 请选择要包含的功能 | Router、Pinia、Linter、Prettier |
| 选择要包含的试验特性 | 使用 Oxfmt 替代 Prettier |
| 跳过所有示例代码，创建一个空白的 Vue 项目？ | No |

::: details 为什么选择 Oxfmt
Oxfmt 是脚手架提供的实验性格式化选项，目标是用更快的格式化工具替代 Prettier。这里跟随脚手架提供的最新选项先启用它；如果后续团队协作或生态兼容需要回到 Prettier，再统一调整格式化配置。
:::

::: info 跟随官方脚手架
如果你看到新的交互选项，先保持默认选择；后续章节真正用到时，再单独接入。这样既能跟随官方最新项目模板，也不会在初始化阶段引入过多工具。
:::

执行完成后，终端会提示接下来进入项目目录、安装依赖、格式化并启动开发服务：

```bash
cd admin
pnpm install
pnpm format
pnpm dev
```

::: warning ⚠️ 不要再次初始化 Git
脚手架最后会提示“可选：使用以下命令在项目目录中初始化 Git”。本项目已经在根目录初始化过 Git，这里不要在 `admin/` 目录里再次执行 `git init`。
:::

按照提示进入 `admin/` 并安装依赖：

::: code-group

```powershell [Windows PowerShell]
Set-Location .\admin
pnpm install
pnpm format
```

```bash [macOS / Linux]
cd admin
pnpm install
pnpm format
```

:::

完成后，`admin/` 目录里会出现 `package.json`、`src/`、`vite.config.ts`、`tsconfig.json` 等管理台工程文件。

## 脚本说明

初始化完成后，`package.json` 里会生成一些常用脚本：

| 命令 | 作用 |
| --- | --- |
| `pnpm dev` | 启动本地开发服务，用来在浏览器里预览管理台 |
| `pnpm build` | 先做类型检查，再打包生产环境文件 |
| `pnpm build-only` | 只执行 Vite 打包，不做类型检查 |
| `pnpm preview` | 本地预览打包后的产物 |
| `pnpm lint` | 执行代码检查，发现潜在问题 |
| `pnpm format` | 格式化代码，保持代码风格统一 |

::: details 为什么同时有 `build` 和 `build-only`
脚手架生成的 `build` 通常会先执行类型检查，再调用 `build-only` 完成 Vite 打包。

也就是说，日常验收优先使用：

```bash
pnpm build
```

只有在已经单独确认类型没问题、只想快速验证打包流程时，才需要直接执行：

```bash
pnpm build-only
```
:::

::: details 为什么没有写 `pnpm run dev`
`pnpm dev` 是 `pnpm run dev` 的简写，常见脚本可以直接省略 `run`。

两种写法都可以执行 `package.json` 里的脚本：

```bash
pnpm dev
pnpm run dev
```

本教程默认使用更短的写法。
:::

::: tip 当前只做初始化
这一节只确认管理台工程能安装、启动和构建。接口请求、登录状态、菜单权限等内容，等后端基础能力完成后再统一对接。
:::

## ✅ 启动并验证

在 `admin/` 目录下启动开发服务：

```bash
pnpm dev
```

终端会输出本地访问地址，通常是：

```text
Local: http://localhost:5173/
```

打开这个地址，能看到 Vue 默认页面，就说明前端开发服务启动成功。

::: warning ⚠️ 端口被占用
如果 `5173` 被占用，Vite 会自动尝试下一个可用端口。以终端实际输出的地址为准。
:::

验证完成后，回到运行 `pnpm dev` 的终端，按 `Ctrl + C` 停止服务。

## ✅ 构建检查

继续在 `admin/` 目录下执行：

```bash
pnpm build
```

如果看到构建成功信息，说明当前管理台工程可以正常打包。

再执行代码检查：

```bash
pnpm lint
```

这一步用于确认脚手架生成的 ESLint 配置可以正常工作。

## ✅ 确认 Git 状态

回到项目根目录：

::: code-group

```powershell [Windows PowerShell]
Set-Location ..
git status
```

```bash [macOS / Linux]
cd ..
git status
```

:::

应该能看到 `admin/` 下新增了一批管理台工程文件，并且原来的 `admin/.gitkeep` 已经被删除。

下一节开始初始化文档项目：[VitePress 文档项目初始化](./docs-init)。
