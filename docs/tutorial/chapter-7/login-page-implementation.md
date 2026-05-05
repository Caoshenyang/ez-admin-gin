---
title: 登录页实现细节
description: "对齐当前 LoginPage.vue、auth API 和本地登录态工具，讲清登录页 UI、提交流程和记住登录。"
---

# 登录页实现细节

前面一页 [登录态与会话流转](./login-and-session-flow) 已经把登录主链路讲清楚了，这一页继续往下落到页面本身：

> 当前登录页到底是怎么把视觉结构、登录提交、记住登录和本地状态保存接起来的。

::: tip 🎯 本节目标
读完后，你应该能看懂当前 `admin/src/pages/auth/LoginPage.vue` 为什么长这样，也知道以后要改登录页时，哪些地方属于 UI，哪些地方属于登录协议和本地状态。
:::

## 先看当前真实文件

当前登录页相关代码主要集中在四处：

```text
admin/src/pages/auth/LoginPage.vue
admin/src/api/auth.ts
admin/src/api/http.ts
admin/src/utils/auth.ts
```

它们之间的关系是：

```text
LoginPage.vue
  ↓
api/auth.ts -> POST /auth/login
  ↓
utils/auth.ts -> setAuthSession(...)
  ↓
router.push('/dashboard')
```

## 当前登录页在页面层做了什么

`LoginPage.vue` 现在承担的是典型的“协议层页面”职责：

- 渲染品牌区和登录卡片
- 管理用户名、密码、验证码占位和记住登录状态
- 调用 `login(...)`
- 成功后写入本地登录态并跳转到 `/dashboard`

它并不会自己处理：

- Axios 实例创建
- `Authorization` 请求头
- 401 失效清理

这些都已经交给更底层的通用能力了。

## 当前视觉结构为什么是“双栏 + 卡片”

登录页现在不是一张裸表单，而是：

- 左侧深色品牌区
- 右侧登录卡片

这样做的目的不是单纯“更好看”，而是让这个后台从一进门开始就有更明确的产品感。

左侧品牌区当前主要承载：

- `BrandLogo`
- 产品定位说明
- 四条能力摘要

右侧登录卡片则承载：

- 用户名
- 密码
- 验证码占位
- 记住登录
- 忘记密码入口
- 登录按钮

这也让页面逻辑变得更清楚：

> 左侧负责产品表达，右侧负责登录动作。

## 为什么验证码现在只保留 UI 占位

当前 `LoginPage.vue` 里虽然有：

- `captchaText`
- `refreshCaptcha()`
- 验证码输入框

但它并没有真正进入登录请求体。

原因很简单：当前后端登录接口仍然只接受：

- `username`
- `password`

::: warning ⚠️ 不要让前端表单先跑在后端协议前面
如果页面现在就把 `captcha` 当成真实登录字段提交，前后端契约会立刻失真。

所以当前实现保留了验证码位置和交互感，但没有把它伪装成已经接通的安全能力。
:::

## `login(...)` 这一层只做了一件事

`admin/src/api/auth.ts` 非常薄，它只负责：

- 调用 `POST /api/v1/auth/login`
- 返回后端统一响应里的 `data`

这层故意保持很轻，是为了让页面层拿到的是明确的资源函数，而不是自己手写请求细节。

## “记住登录”是怎么落地的

这部分真正的核心不在页面，而在：

- `admin/src/utils/auth.ts`

当前实现有一个很明确的约定：

- 勾选“记住登录”时，登录态写入 `localStorage`
- 不勾选时，登录态写入 `sessionStorage`

也就是说，“记住登录”不是一个装饰复选框，而是真的决定了登录态的存储位置。

## 当前本地登录态到底保存了什么

`setAuthSession(...)` 现在会保存三类信息：

| 键 | 内容 |
| --- | --- |
| `ez-admin-access-token` | 访问令牌 |
| `ez-admin-token-type` | Token 类型，默认 `Bearer` |
| `ez-admin-user-info` | 当前用户的 `userId / username / nickname / expiresAt` |

这样后续前端就能在多个地方复用：

- `api/http.ts` 读取 Token 拼接请求头
- `AdminLayout.vue` 读取当前登录用户昵称
- 路由守卫判断当前是否已登录

## 当前提交动作怎么跑完闭环

`handleSubmit()` 的真实顺序现在很稳定：

1. 先跑 Naive UI 表单校验
2. 调用 `login({ username, password })`
3. 用 `setAuthSession(result, rememberLogin)` 保存本地登录态
4. 提示“登录成功”
5. 跳转到 `/dashboard`

如果失败，则根据 Axios 错误结构优先取后端返回的 `message`。

这意味着登录页虽然是 UI 页面，但它已经和当前后端错误格式、登录响应格式完整对齐了。

## 为什么页面里一进来就会判断 `hasAccessToken()`

当前实现里，`LoginPage.vue` 会在进入时检查：

- 本地是否已经有 Token

如果已经有，就直接跳转到 `/dashboard`。

它的意义是避免用户在“已经登录”的情况下又看到登录页。

不过这只是页面层的第一道短路；真正更完整的登录态守卫，还是在 `router/index.ts` 里。

## 忘记密码入口为什么现在只提示未接入

当前页面保留了“忘记密码？”这个入口，但点击只会弹出一条提示。

这是一种很实在的处理方式，因为它明确区分了两件事：

- 交互入口已经预留
- 后续找回密码流程还没有接通

这比做一个无效链接或假页面更诚实，也更方便后面继续扩展。

## 这一页和登录链路总览页的关系

可以这样理解两页的分工：

- [登录态与会话流转](./login-and-session-flow)：讲“登录这条链路怎么跑”
- 这一页：讲“登录页这个具体页面怎么实现”

前者偏链路，后者偏页面实现。

## 最小验收方式

这一页改动或理解完成后，至少可以按下面顺序手动验证：

1. 打开 `/login`，确认看到双栏登录页
2. 输入正确账号密码，点击登录
3. 登录成功后跳到 `/dashboard`
4. 勾选和取消“记住登录”分别登录一次，确认刷新后的保留行为符合预期
5. 返回 `/login`，确认已有 Token 时会直接回到后台

## 下一步

- 想继续看登录后页面容器怎么接住这些状态，下一页读 [后台布局与工作标签](./admin-layout-and-worktabs)
- 想继续看登录后菜单和页面路由怎么注册，继续读 [动态菜单注册与按钮权限](./dynamic-route-registration)
