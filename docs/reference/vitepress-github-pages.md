---
title: VitePress 部署到 GitHub Pages
description: "从零开始将 VitePress 文档站点部署到 GitHub Pages：配置 base 路径、编写 GitHub Actions 工作流、处理包管理器差异和常见报错。"
---

# VitePress 部署到 GitHub Pages

这一页帮你把一个本地能跑的 VitePress 项目，变成一个通过 GitHub Actions 自动构建、部署到 GitHub Pages 的在线文档站点。全程不需要手动上传文件，也不需要额外的服务器。

::: tip 🎯 读完这页你能做到什么
- 修改 VitePress 配置，让站点在项目子路径下正常工作
- 创建 GitHub Actions 工作流，推送代码后自动部署
- 处理 npm / pnpm / yarn 的差异，避免缓存和安装报错
- 排查部署过程中最常见的几个错误
:::

## 前置条件

- 有一个 GitHub 仓库，里面已经有可运行的 VitePress 项目（`docs/` 目录下）
- 本地 `pnpm docs:build`（或对应的构建命令）能正常生成 `docs/.vitepress/dist`
- 仓库有 `main` 分支

## 第一步：配置 base 路径

VitePress 默认假设站点部署在域名根路径 `/`。GitHub Pages 项目站点的实际地址是 `https://<username>.github.io/<repo>/`，所以必须把 `base` 改成仓库名：

```ts
// docs/.vitepress/config.mts
export default defineConfig({
  base: '/<repo>/',  // [!code focus]
  // ...其余配置
})
```

把 `<repo>` 替换成你的 GitHub 仓库名。例如仓库名是 `ez-admin-gin`，就写：

```ts
base: '/ez-admin-gin/',
```

::: warning ⚠️ 前后带斜杠
`base` 的值必须以 `/` 开头并以 `/` 结尾，例如 `/ez-admin-gin/`。写成 `ez-admin-gin` 或 `/ez-admin-gin` 都会导致静态资源加载失败。
:::

::: details 什么时候 base 写 `/`
只有一种情况：你使用的是 **用户站点**（`<username>.github.io`），仓库本身就是 `<username>.github.io`，站点直接部署在域名根路径。这时 `base` 保持默认的 `/` 即可。

本文只覆盖**项目站点**（`<username>.github.io/<repo>/`）的情况。
:::

## 第二步：创建 GitHub Actions 工作流

在仓库根目录创建 `.github/workflows/deploy-docs.yml`。下面给出完整文件，再逐段解释。

### 使用 pnpm（推荐）

```yaml
# .github/workflows/deploy-docs.yml
name: Deploy Docs

on:
  push:
    branches: [main]
    paths:
      - 'docs/**'
      - '.github/workflows/deploy-docs.yml'
  workflow_dispatch:

permissions:
  contents: read
  pages: write
  id-token: write

concurrency:
  group: pages
  cancel-in-progress: false

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - uses: pnpm/action-setup@v4
        with:
          version: 9

      - uses: actions/setup-node@v4
        with:
          node-version: 22
          cache: pnpm
          cache-dependency-path: docs/pnpm-lock.yaml

      - run: pnpm install --frozen-lockfile
        working-directory: docs

      - run: pnpm run docs:build
        working-directory: docs

      - uses: actions/upload-pages-artifact@v3
        with:
          path: docs/.vitepress/dist

  deploy:
    environment:
      name: github-pages
      url: ${{ steps.deployment.outputs.page_url }}
    runs-on: ubuntu-latest
    needs: build
    steps:
      - name: Deploy to GitHub Pages
        id: deployment
        uses: actions/deploy-pages@v4
```

### 使用 npm

```yaml
# .github/workflows/deploy-docs.yml
name: Deploy Docs

on:
  push:
    branches: [main]
    paths:
      - 'docs/**'
      - '.github/workflows/deploy-docs.yml'
  workflow_dispatch:

permissions:
  contents: read
  pages: write
  id-token: write

concurrency:
  group: pages
  cancel-in-progress: false

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - uses: actions/setup-node@v4
        with:
          node-version: 22
          cache: npm
          cache-dependency-path: docs/package-lock.json

      - run: npm ci
        working-directory: docs

      - run: npm run docs:build
        working-directory: docs

      - uses: actions/upload-pages-artifact@v3
        with:
          path: docs/.vitepress/dist

  deploy:
    environment:
      name: github-pages
      url: ${{ steps.deployment.outputs.page_url }}
    runs-on: ubuntu-latest
    needs: build
    steps:
      - name: Deploy to GitHub Pages
        id: deployment
        uses: actions/deploy-pages@v4
```

::: details 使用 yarn 的写法
```yaml
- uses: actions/setup-node@v4
  with:
    node-version: 22
    cache: yarn
    cache-dependency-path: docs/yarn.lock

- run: yarn install --frozen-lockfile
  working-directory: docs

- run: yarn docs:build
  working-directory: docs
```
:::

## 第三步：启用 GitHub Pages

1. 打开仓库 → **Settings** → **Pages**
2. **Build and deployment** → **Source** 选择 **GitHub Actions**

::: warning ⚠️ 什么时候能看到这个选项
如果工作流文件还没推送到 `main` 分支，这里可能只显示 "Deploy from a branch"。先把工作流文件推送上去，GitHub 检测到后会自动出现 GitHub Actions 选项。
:::

## 第四步：推送并验证

1. 将 `base` 配置和工作流文件一起提交到 `main` 分支
2. 打开仓库 → **Actions** 标签页，应该能看到 `Deploy Docs` 工作流正在运行
3. 等待构建和部署完成（通常 1-2 分钟）
4. 访问 `https://<username>.github.io/<repo>/` 查看站点

## 工作流逐段解析

### 触发条件

```yaml
on:
  push:
    branches: [main]
    paths:
      - 'docs/**'
      - '.github/workflows/deploy-docs.yml'
  workflow_dispatch:
```

| 字段 | 作用 |
| --- | --- |
| `branches: [main]` | 只有 `main` 分支的推送才触发 |
| `paths` | 只在 `docs/` 目录或工作流文件有变更时触发，避免无关改动浪费 Actions 时间 |
| `workflow_dispatch` | 允许在 Actions 页面手动触发，用于初次部署或紧急重发布 |

### 权限

```yaml
permissions:
  contents: read
  pages: write
  id-token: write
```

| 权限 | 作用 |
| --- | --- |
| `contents: read` | 读取仓库代码 |
| `pages: write` | 写入 GitHub Pages 部署 |
| `id-token: write` | 用于 OIDC 身份验证，`deploy-pages` Action 需要这个权限 |

### 并发控制

```yaml
concurrency:
  group: pages
  cancel-in-progress: false
```

确保同一时间只有一个部署任务在跑。`cancel-in-progress: false` 表示新的推送不会取消正在进行的部署，而是排队等待。

### 依赖安装与缓存

pnpm 示例中的关键部分：

```yaml
- uses: pnpm/action-setup@v4
  with:
    version: 9

- uses: actions/setup-node@v4
  with:
    node-version: 22
    cache: pnpm
    cache-dependency-path: docs/pnpm-lock.yaml
```

| 步骤 | 作用 |
| --- | --- |
| `pnpm/action-setup@v4` | 在 Runner 上安装 pnpm（npm 不需要这步） |
| `cache: pnpm` | 根据 lock 文件缓存依赖，加速后续构建 |
| `cache-dependency-path` | 指向 lock 文件的路径，确保缓存能正确命中 |
| `--frozen-lockfile` | 严格按 lock 文件安装，CI 环境不允许自动更新依赖 |

::: warning ⚠️ lock 文件必须存在且已提交
`cache-dependency-path` 指向的文件必须已经提交到仓库。如果 lock 文件不存在，`setup-node` 会报错 `Some specified paths were not resolved, unable to cache dependencies`。
:::

### 构建与上传

```yaml
- run: pnpm run docs:build
  working-directory: docs

- uses: actions/upload-pages-artifact@v3
  with:
    path: docs/.vitepress/dist
```

`docs:build` 命令会将站点生成到 `docs/.vitepress/dist` 目录。`upload-pages-artifact` 把这个目录打包，交给下一步部署。

### 部署

```yaml
deploy:
  environment:
    name: github-pages
    url: ${{ steps.deployment.outputs.page_url }}
  runs-on: ubuntu-latest
  needs: build
  steps:
    - name: Deploy to GitHub Pages
      id: deployment
      uses: actions/deploy-pages@v4
```

`deploy` Job 等待 `build` 完成后执行。`deploy-pages` Action 从 artifact 中取出构建产物并发布到 GitHub Pages。部署完成后，`page_url` 会输出站点地址，你可以在 Actions 日志中看到。

## 常见问题

### 构建报错：`unable to cache dependencies`

```
Error: Some specified paths were not resolved, unable to cache dependencies.
```

**原因**：`cache-dependency-path` 指向的 lock 文件不存在。

**解决**：确保对应的 lock 文件已提交到仓库：

| 包管理器 | 需要提交的文件 |
| --- | --- |
| pnpm | `docs/pnpm-lock.yaml` |
| npm | `docs/package-lock.json` |
| yarn | `docs/yarn.lock` |

如果本地没有 lock 文件，先在 `docs/` 目录下执行一次安装命令生成它：

::: code-group

```bash [pnpm]
cd docs && pnpm install
```

```bash [npm]
cd docs && npm install
```

```bash [yarn]
cd docs && yarn install
```

:::

### 页面样式错乱或 404

**原因**：`base` 配置与实际访问路径不匹配。

**排查**：

1. 确认 `base` 值为 `'/<repo>/'`，前后都有斜杠
2. 确认仓库名与 `base` 中的名称完全一致（区分大小写）
3. 打开浏览器开发者工具，查看静态资源（CSS/JS）请求路径是否以 `/<repo>/` 开头

### Actions 没有自动触发

**排查清单**：

- 确认工作流文件在 `main` 分支上（而不是只在功能分支）
- 确认推送的提交包含了 `docs/` 目录的变更（受 `paths` 过滤）
- 打开 **Settings → Pages → Build and deployment**，确认 Source 是 **GitHub Actions**
- 检查 **Settings → Actions → General**，确认没有禁用 Actions

### 自定义域名

如果你有自己的域名，可以绑定到 GitHub Pages：

1. 在 `docs/public/` 下创建 `CNAME` 文件，内容为你的域名（如 `docs.example.com`）
2. 在域名服务商处添加 CNAME 记录，指向 `<username>.github.io`
3. 等待 DNS 生效

使用自定义域名后，`base` 可以改回 `/`，因为站点不再部署在子路径下了。

## Node.js 版本选择

GitHub Actions Runner 上 Node.js 20 已进入弃用周期，建议使用 Node.js 22：

```yaml
- uses: actions/setup-node@v4
  with:
    node-version: 22  # [!code focus]
```

如果你仍需使用 Node.js 20，可以在工作流中设置环境变量临时绕过：

```yaml
env:
  ACTIONS_ALLOW_USE_UNSECURE_NODE_VERSION: true
```

但这只应该作为过渡方案，长期建议升级到 Node.js 22。
