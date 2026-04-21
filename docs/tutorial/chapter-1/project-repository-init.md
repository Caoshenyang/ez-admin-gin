---
title: 项目仓库初始化
description: "创建 ez-admin-gin 项目目录，初始化 Git 仓库，并建立第一层基础目录。"
---

# 项目仓库初始化

这一节先把项目的起点固定下来：创建项目目录，初始化 Git 仓库，再建立后续章节都会用到的一级目录。

::: tip 🎯 本节目标
完成后，项目根目录会包含 `server`、`admin`、`docs`、`deploy`、`scripts` 五个目录，并且可以通过 `git status` 看到它们已经进入 Git 工作区。
:::

::: info Windows 终端约定
本教程在 Windows 下默认使用 PowerShell。相比传统 `cmd`，PowerShell 更适合现代开发环境，也方便后续执行脚本、查看端口和排查问题。
:::

## 创建项目目录

在你平时放项目的目录下执行：

::: code-group

```powershell [Windows PowerShell]
Set-Location D:\A
New-Item -ItemType Directory -Name ez-admin-gin
Set-Location .\ez-admin-gin
```

```bash [macOS / Linux]
cd ~/Projects
mkdir ez-admin-gin
cd ez-admin-gin
```

:::

::: warning ⚠️ 如果目录已经存在
直接进入已有目录即可，不需要重复创建同名目录。
:::

## 初始化 Git 仓库

先确认本机已经安装 Git：

```bash
git --version
```

如果提示找不到 `git` 命令，先安装 Git：

- Git 官网：[https://git-scm.com](https://git-scm.com)

安装完成后，在 `ez-admin-gin` 目录下执行：

```bash
git init
```

::: details 为什么一开始就初始化 Git
Git 用来记录项目的每一次变更。后面新增目录、初始化后端、初始化前端、调整配置时，都可以通过 Git 看清楚改了什么，也方便在出错时回到上一个可用状态。
:::

## 创建基础目录

接下来创建第一层目录：

::: code-group

```powershell [Windows PowerShell]
New-Item -ItemType Directory -Path server, admin, docs, deploy, scripts -Force

foreach ($dir in 'server', 'admin', 'docs', 'deploy', 'scripts') {
  New-Item -ItemType File -Path "$dir/.gitkeep" -Force
}
```

```bash [macOS / Linux]
mkdir -p server admin docs deploy scripts
touch server/.gitkeep admin/.gitkeep docs/.gitkeep deploy/.gitkeep scripts/.gitkeep
```

:::

完成后，项目根目录会是这样：

```text
ez-admin-gin/
├─ server/
├─ admin/
├─ docs/
├─ deploy/
└─ scripts/
```

::: warning ⚠️ 为什么要加 `.gitkeep`
Git 默认不会记录空目录。`.gitkeep` 只是占位文件，用来让这些目录先被提交。后面目录里有了真实文件后，可以删掉对应的 `.gitkeep`。
:::

## 目录职责

| 目录 | 放什么 | 不放什么 |
| --- | --- | --- |
| `server/` | 后端服务代码、配置、数据库访问、接口实现 | 前端页面、文档站配置 |
| `admin/` | 管理台前端代码、页面、组件、接口调用 | 后端业务逻辑、门户页面、部署脚本 |
| `docs/` | 项目文档、教程、参考手册、VitePress 配置 | 后端或前端运行时代码 |
| `deploy/` | 部署相关文件，例如 Docker、Nginx、环境示例 | 日常开发脚本、业务代码 |
| `scripts/` | 可重复执行的辅助脚本，例如检查、构建、启动 | 只执行一次的临时命令 |

::: tip 目录边界先简单一点
这一节只固定第一层目录。`server/` 和 `admin/` 里面的细分结构，等对应项目初始化时再展开。
:::

## ✅ 确认结果

查看当前目录：

::: code-group

```powershell [Windows PowerShell]
Get-ChildItem
```

```bash [macOS / Linux]
ls
```

:::

应该能看到 `server`、`admin`、`docs`、`deploy`、`scripts`。

再查看 Git 状态：

```bash
git status
```

应该能看到 5 个 `.gitkeep` 文件处于未跟踪状态，说明这些目录已经可以被 Git 记录。

下一节开始初始化后端项目：[Go 后端项目初始化](./backend-init)。
