---
title: 仓库创建与 Git 初始化
description: "创建 ez-admin-gin 项目目录，并初始化本地 Git 仓库。"
---

# 仓库创建与 Git 初始化

这一节只完成两件事：

1. 创建 `ez-admin-gin` 项目目录
2. 初始化本地 Git 仓库

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

如果目录已经存在，直接进入即可。

## 初始化 Git 仓库

Git 用来记录项目的每一次变更。后面新增目录、初始化后端、初始化前端、调整配置时，都可以通过 Git 看清楚改了什么，也方便在出错时回退到上一个可用状态。

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

## 确认结果

```bash
git status
```

能看到当前分支和工作区状态，就说明 Git 仓库初始化完成。

下一节开始创建基础目录：[单仓库目录结构](./directory-structure)。
