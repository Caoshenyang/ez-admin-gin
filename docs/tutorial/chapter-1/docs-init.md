---
title: VitePress 文档项目初始化
description: "说明 docs 文档站是可选模块，愿意维护文档时可按 VitePress 官方流程初始化。"
---

# VitePress 文档项目初始化

`docs/` 不是后台底座运行所必需的模块。它适合用来放使用说明、接口说明、部署记录或项目教程；如果你暂时不想维护文档，可以跳过这一节，后面再补也不影响后端和管理台继续开发。

::: tip 🎯 本节怎么读
想写文档，就按 VitePress 官方流程初始化；不想写文档，保留 `docs/` 目录或先删掉都可以。
:::

## 是否需要 docs

可以按下面的方式判断：

| 情况 | 建议 |
| --- | --- |
| 项目会交给别人使用 | 建议保留 `docs/`，至少写安装、启动和部署说明 |
| 项目需要长期复用 | 建议保留 `docs/`，记录目录约定、模块接入方式和常见问题 |
| 只是自己快速验证功能 | 可以先跳过，等功能稳定后再补 |

::: info 官方文档
VitePress 官网：[https://vitepress.dev](https://vitepress.dev)

初始化说明：[https://vitepress.dev/guide/getting-started](https://vitepress.dev/guide/getting-started)
:::

## 如果你要初始化

如果决定维护文档，可以进入 `docs/` 目录，按官网流程创建 VitePress 项目：

::: code-group

```powershell [Windows PowerShell]
# 进入文档目录并移除占位文件
Set-Location .\docs
Remove-Item .gitkeep -ErrorAction SilentlyContinue
# 安装 VitePress，并按向导初始化文档站
pnpm add -D vitepress@next
pnpm vitepress init
```

```bash [macOS / Linux]
# 进入文档目录并移除占位文件
cd docs
rm -f .gitkeep
# 安装 VitePress，并按向导初始化文档站
pnpm add -D vitepress@next
pnpm vitepress init
```

:::

初始化完成后，按向导提示运行本地预览和构建命令即可。

下一节开始准备本地基础环境：[Docker Compose 基础环境](./docker-compose-env)。
