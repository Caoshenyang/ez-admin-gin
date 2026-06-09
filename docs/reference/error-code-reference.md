---
title: 错误码参考
description: "集中记录 EZ Admin Gin 当前统一响应体里的业务错误码、HTTP 状态码映射和最常见的返回场景。"
---

# 错误码参考

这页只做快速查表，不展开讲完整教程。你在联调接口、排查前端弹错或补新模块时，优先以这里的错误码约定为准。

::: tip 🎯 先记住一条主线
当前后端把“HTTP 状态码”和“业务错误码”分成两层：

- HTTP 状态码负责表达协议层语义
- `code` 字段负责表达前后端都能稳定识别的业务结果
:::

## 当前统一响应体

所有接口都走同一份响应结构：

```json
{
  "code": 0,
  "message": "ok",
  "data": {}
}
```

对应实现位于：

- `server/internal/pkg/httpx/response.go`

字段职责可以直接这样理解：

| 字段 | 用途 |
| --- | --- |
| `code` | 前后端稳定判断结果的业务码 |
| `message` | 可以直接展示或记录日志的提示语 |
| `data` | 成功时返回的业务数据 |

## 当前固定错误码

当前业务错误码集中定义在：

- `server/internal/pkg/errorsx/errors.go`

固定值如下：

| 业务码 | HTTP 状态码 | 含义 |
| --- | --- | --- |
| `0` | `200` | 请求成功 |
| `40000` | `400` | 请求参数错误 |
| `40100` | `401` | 未登录或登录已过期 |
| `40300` | `403` | 没有权限访问资源 |
| `40400` | `404` | 资源不存在 |
| `42900` | `429` | 请求过于频繁 |
| `50300` | `503` | 数据库、Redis 等依赖不可用 |
| `50000` | `500` | 服务器内部错误 |

::: info 为什么 `code` 不直接等于 HTTP 状态码
因为前端真正稳定消费的是业务语义，不是协议数字本身。

例如两个接口都返回 `400`，但前端和排障人员更关心的是“这是参数错误”还是“这是内部兜底错误”。当前项目先用一组固定业务码把这层语义稳定下来。
:::

## 当前错误构造方式

当前后端不会在 Handler 里手写一堆 `c.JSON(...)`，而是先构造应用错误，再交给统一响应层输出。

常用构造函数：

| 函数 | 用途 |
| --- | --- |
| `errorsx.BadRequest(message)` | 参数校验失败 |
| `errorsx.Unauthorized(message)` | 未登录、登录失效、用户名密码错误 |
| `errorsx.TooManyRequests(message)` | 请求过于频繁 |
| `errorsx.Forbidden(message)` | 权限不足 |
| `errorsx.NotFound(message)` | 资源不存在 |
| `errorsx.ServiceUnavailable(message, err)` | 依赖服务异常 |
| `errorsx.Internal(message, err)` | 未归类服务内部错误 |

对应的统一输出入口是：

- `httpx.Success(c, data)`
- `httpx.Error(c, err, log)`

## 最常见返回场景

### `40000` 请求参数错误

这是当前最常见的一类错误，通常出现在：

- DTO 归一化失败
- 列表筛选参数非法
- 路径参数不是合法 ID
- 业务前置条件不满足

典型例子：

- 用户名为空
- 菜单编码重复
- 公告状态不正确
- 不能禁用当前登录用户

这类错误多数在下面几层产生：

- `server/internal/modules/*/*/api/dto.go`
- `server/internal/modules/*/*/application/service.go`
- 少量路径参数解析处的 `api/handler.go`

### `40100` 未登录或登录失效

通常来自：

- `platform/middleware.Auth`
- `/auth/login` 用户名密码校验失败
- `/auth/me`、`/auth/menus`、`/auth/dashboard` 访问时没有有效登录态

典型提示语包括：

- `请先登录`
- `登录已过期，请重新登录`
- `用户名或密码错误`
- `登录状态无效，请重新登录`

### `40300` 权限不足

通常来自：

- `platform/middleware.Permission`

它表达的是：

- 已经登录
- 但当前角色没有访问该接口的 Casbin 权限

当前典型提示语是：

- `没有权限访问`

### `40400` 资源不存在

通常来自 Repository 在把 `gorm.ErrRecordNotFound` 转成业务错误时。

典型资源包括：

- 用户不存在
- 菜单不存在
- 配置不存在
- 公告不存在

### `42900` 请求过于频繁

通常来自登录接口限流或账号锁定逻辑。

它表达的是：

- 当前请求格式可能没错
- 但触发了短时间内重复请求或连续失败保护

### `50300` 依赖不可用

这一类当前出现得比 `50000` 少，但语义很明确：

- 应用逻辑本身未必写错
- 只是数据库、Redis 或其他依赖暂时不可用

更适合在下面场景里使用：

- 初始化连接依赖失败
- 外部基础设施短暂不可用

### `50000` 服务器内部错误

这类错误表示：

- 当前请求失败
- 但不应该把底层内部细节直接回给前端

如果错误不是 `*errorsx.Error`，`httpx.Error(...)` 会：

1. 把原始错误打到日志里
2. 对前端统一返回 `50000 + 服务器内部错误`

## 当前推荐的抛错分工

如果你在补新模块，当前更稳的分工是：

| 层 | 更适合做什么 |
| --- | --- |
| `dto.go` | 字段格式、长度、枚举值校验 |
| `service.go` | 业务前置条件、事务内规则判断 |
| `repository.go` | `not found`、唯一性、持久化细节转换 |
| `handler.go` | 参数绑定失败、路径参数解析失败、把错误统一回写 |

::: warning 不要把数据库原始错误直接透给前端
当前主线要求前端看到的是可理解、可预期的业务语义，而不是数据库错误字符串或堆栈细节。

所以未归类错误应该走：

```go
httpx.Error(c, errorsx.Internal("查询用户失败", err), log)
```
:::

## 一个新模块最少要对齐什么

如果你要新增一个模块，至少把下面这几类错误路径补齐：

1. 参数错误统一返回 `40000`
2. 未登录和登录失效统一返回 `40100`
3. 权限不足统一返回 `40300`
4. 资源不存在统一返回 `40400`
5. 请求过于频繁统一返回 `42900`
6. 未分类异常统一兜底到 `50000`

这样前端、日志和接口排障才能保持一致。

## 相关代码与文档

- `server/internal/pkg/errorsx/errors.go`
- `server/internal/pkg/httpx/response.go`
- `server/internal/platform/middleware/auth.go`
- `server/internal/platform/middleware/permission.go`
- [接口风格决策](./api-style-decision)
- [权限码约定](./permission-code-conventions)
