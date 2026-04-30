---
title: Token 签发与解析
description: "把用户名密码登录升级成真正的登录态，让认证模块返回 access token，并为后续中间件和权限链路提供统一身份载荷。"
---

# Token 签发与解析

上一节只是确认“用户名密码对不对”。这一节开始让登录成功真正变成一个可持续使用的登录态：服务端签发 `access_token`，后续接口通过它识别当前用户。

::: tip 🎯 本节目标
完成后，`/api/v1/auth/login` 会返回 `access_token`、`token_type` 和 `expires_at`；认证模块会把 Token 签发职责收敛到统一的 `authn` 管理器中，而不是把 JWT 逻辑散落在登录 Handler 里。
:::

## 本节会改什么

当前主线里，这一节主要对应下面这些位置：

```text
server/
├─ configs/
│  └─ config.yaml
├─ internal/
│  ├─ config/
│  │  └─ config.go
│  ├─ platform/
│  │  └─ authn/
│  │     └─ authn.go
│  ├─ token/
│  │  └─ jwt.go
│  └─ module/
│     └─ auth/
│        ├─ dto.go
│        ├─ login_service.go
│        └─ login_handler.go
```

| 位置 | 用途 |
| --- | --- |
| `config.yaml` | 声明 JWT 密钥、签发方和有效期 |
| `config.go` | 读取 `auth` 配置段 |
| `token/jwt.go` | 封装 Token 生成和解析 |
| `platform/authn/authn.go` | 给最终版结构提供统一 `authn` 命名空间 |
| `module/auth/login_service.go` | 登录成功后签发 Token |
| `module/auth/login_handler.go` | 只负责协议层绑定与输出 |

## 为什么 Token 不应该直接写死在登录 Handler 里

从“能登录”走到“能长期维护”，最重要的一步就是把登录态能力单独收口。

如果 JWT 逻辑直接散在 Handler 里，后面你一旦再补：

- Refresh Token
- 多端登录策略
- 登录失效原因分类
- 自定义 claims

就会发现认证代码越来越难整理。  
当前主线里，这部分职责已经被拆成：

- `token/jwt.go` 负责底层签发和解析
- `platform/authn` 提供统一命名空间
- `module/auth/login_service.go` 决定什么时候签发 Token

## Token 里现在装了什么

当前项目的 access token 至少承载这些信息：

- `user_id`
- `username`
- `issuer`
- `issued_at`
- `expires_at`

这样做的好处是，后续认证中间件不用每次先查数据库才能知道“当前是谁”，而是可以先从 Token 中拿到最基础的身份载荷。

::: warning ⚠️ Token 不是用户快照
这一节的 Token 只承载最基础的身份信息，不把完整用户资料、角色树或权限列表直接塞进去。真正会变化的组织和权限信息，后面仍然通过中间件和数据库链路补齐。
:::

## 配置层为什么要单独给 `auth` 一段

当前主线里，JWT 相关配置不是散在代码常量里，而是统一从 `auth` 段读取，例如：

- `jwt_secret`
- `access_token_ttl`
- `issuer`

这一步看起来只是“多了一段配置”，但它会直接影响后续能力是否好扩展：

- 本地开发可以直接跑
- 生产环境可以通过环境变量覆盖
- 认证策略不会和业务代码耦死在一起

## 登录接口现在对外长什么样

当前 `/api/v1/auth/login` 登录成功后，返回结构至少包含：

```json
{
  "user_id": 1,
  "username": "admin",
  "nickname": "系统管理员",
  "access_token": "xxx.yyy.zzz",
  "token_type": "Bearer",
  "expires_at": "2026-04-30T10:00:00Z"
}
```

其中最关键的是：

- `access_token`：后续所有受保护接口的身份凭证
- `token_type`：明确前端用 `Bearer` 方式拼接请求头
- `expires_at`：让前端和用户都知道登录态何时失效

## `module/auth` 在这一节里承担了什么职责

认证模块现在的结构不是“为了分层而分层”，而是为了给后面留扩展空间。  
这一节最值得记住的是登录链路已经明确拆成：

- `LoginHandler`
  负责请求绑定和返回响应
- `LoginService`
  负责用户名密码校验、签发 Token、写登录日志
- `authn.Manager`
  负责真正的 Token 生成和解析

也就是说，现在登录成功不再只是“查表后返回用户信息”，而是已经形成了完整的登录态生成流程。

## 怎么验证这一节已经做成

### 1. 后端构建通过

在 `server/` 目录执行：

```bash
go test ./...
```

应该能看到认证模块和 `platform/authn` 已经进入编译链路。

### 2. 登录成功能返回完整登录态

调用：

```text
POST /api/v1/auth/login
```

成功后，响应里应该已经包含：

- `access_token`
- `token_type`
- `expires_at`

### 3. Token 结构看起来像标准 JWT

`access_token` 一般会是三段式结构：

```text
header.payload.signature
```

也就是中间会出现两个 `.`。

::: details 为什么这一节先不讲 Refresh Token
这套教程当前先把企业后台里最核心、最稳定的 access token 主线收稳。等登录、认证中间件、菜单和数据权限都串起来之后，再继续补更重的会话策略，会更顺。
:::

## 本节最关键的收获

这一节真正建立的判断标准是：

> 登录成功不等于认证体系完成。只有当系统能稳定签发、解析和传递 access token，后续的认证中间件和权限链路才真正成立。

Token 这一节的价值，不只是“多返回一个字符串”，而是把后台底座从“会校验密码”推进到了“拥有正式登录态”。

下一节继续把登录态接进请求链路：[登录校验中间件](./auth-middleware)。
