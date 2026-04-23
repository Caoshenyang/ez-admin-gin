---
title: 接口风格决策
description: "说明 EZ Admin Gin 为什么采用企业后台常见的 GET + POST 接口风格，以及路径命名如何表达业务动作。"
---

# 接口风格决策

EZ Admin Gin 的接口风格优先贴近企业后台系统：**读操作使用 `GET`，写操作默认使用 `POST`，路径表达业务动作和资源归属**。

::: tip 项目决策
本项目不追求纯 RESTful 写法，而是采用企业后台里更常见、更容易落地的风格：查询走 `GET`，新增、编辑、删除、状态变更、授权、导入导出等写操作统一走 `POST`。
:::

::: info 阅读建议
如果只关心项目规范，直接看“本项目的划分规则”和“推荐路径格式”；如果想理解为什么这么定，再回来看前面的取舍说明。
:::

## 为什么不采用纯 RESTful

纯 RESTful 通常会这样设计：

| 操作 | RESTful 写法 |
| --- | --- |
| 查询列表 | `GET /api/v1/system/users` |
| 创建用户 | `POST /api/v1/system/users` |
| 编辑用户 | `PUT /api/v1/system/users/:id` |
| 修改状态 | `PATCH /api/v1/system/users/:id/status` |
| 删除用户 | `DELETE /api/v1/system/users/:id` |

这种写法语义清晰，适合开放 API、平台 API、资源模型稳定的服务。

但在企业后台系统里，接口往往还要服务更具体的工程约束。

### 历史兼容和前端表单惯性

HTML 原生 `<form>` 长期只直接支持 `GET` 和 `POST`。如果不使用 Ajax、Fetch 或额外的 method override 机制，浏览器不能直接提交 `PUT`、`PATCH`、`DELETE` 请求。

虽然现代前端已经可以轻松发起任意 HTTP 方法，但后台系统、低代码平台、老项目和内部工具里仍然保留了大量 `GET / POST` 的使用习惯。

### 网络设施和内网环境差异

企业系统经常经过代理、网关、防火墙、负载均衡、统一鉴权平台等多层基础设施。大多数环境都天然支持 `GET` 和 `POST`，但部分旧网关或安全策略可能会限制 `PUT`、`PATCH`、`DELETE`。

对后台底座来说，接口风格越稳定，部署到不同公司、不同内网环境时遇到的兼容性问题就越少。

### 团队协作和权限审计成本

严格 RESTful 要求团队对资源建模、幂等性、`PUT` 与 `PATCH` 的区别保持一致。如果团队成员理解不一致，最后反而会出现一部分接口用 `PUT`，一部分接口用 `POST /update`，风格更乱。

后台系统的权限、审计、操作日志也更关注“谁执行了什么业务动作”，例如“禁用账号”“分配角色”“删除菜单”。这类动作通常比纯资源动词更适合放进路径里表达。

所以本项目选择的是一种更务实的后台接口风格：**路径保持资源化，写操作统一 `POST`，具体动作由路径表达**。

## GET / POST 和 RESTful 的取舍

<table>
  <colgroup>
    <col width="120">
    <col>
    <col>
  </colgroup>
  <thead>
    <tr>
      <th>维度</th>
      <th>GET / POST 后台风格</th>
      <th>严格 RESTful 风格</th>
    </tr>
  </thead>
  <tbody>
    <tr>
      <td><strong>语义表达</strong></td>
      <td>由路径表达业务动作，例如 <code>/:id/update</code>、<code>/:id/delete</code></td>
      <td>由 HTTP 方法表达资源操作，例如 <code>PUT</code>、<code>DELETE</code></td>
    </tr>
    <tr>
      <td><strong>兼容性</strong></td>
      <td>更高，几乎所有网关、代理、表单工具都支持</td>
      <td>依赖环境完整支持多种 HTTP 方法</td>
    </tr>
    <tr>
      <td><strong>团队成本</strong></td>
      <td>较低，只需统一“查询 GET、写入 POST”</td>
      <td>较高，需要准确区分 <code>PUT</code>、<code>PATCH</code>、幂等性等概念</td>
    </tr>
    <tr>
      <td><strong>权限审计</strong></td>
      <td>更适合按“路径 + 方法”定位业务动作</td>
      <td>更适合资源模型稳定的开放接口</td>
    </tr>
    <tr>
      <td><strong>常见风险</strong></td>
      <td>如果路径命名混乱，<code>POST</code> 语义会变模糊</td>
      <td>如果团队理解不一致，接口方法会混用</td>
    </tr>
  </tbody>
</table>

这个选择不是否定 RESTful，而是为后台底座选择更低维护成本、更强环境适应性的方案。

## 本项目的划分规则

| 场景 | 方法 | 命名建议 | 示例 |
| --- | --- | --- | --- |
| 查询列表 | `GET` | 资源复数名 | `GET /api/v1/system/users` |
| 查询详情 | `GET` | 资源 ID | `GET /api/v1/system/users/:id` |
| 创建资源 | `POST` | 资源复数名 | `POST /api/v1/system/users` |
| 编辑资源 | `POST` | `/:id/update` | `POST /api/v1/system/users/:id/update` |
| 修改状态 | `POST` | `/:id/status` | `POST /api/v1/system/users/:id/status` |
| 删除资源 | `POST` | `/:id/delete` | `POST /api/v1/system/users/:id/delete` |
| 批量删除 | `POST` | `/batch-delete` | `POST /api/v1/system/users/batch-delete` |
| 分配关系 | `POST` | `/:id/roles`、`/:id/menus` | `POST /api/v1/system/users/:id/roles` |
| 登录登出 | `POST` | 动作名 | `POST /api/v1/auth/login` |
| 上传文件 | `POST` | 资源名 | `POST /api/v1/system/files` |
| 导入导出 | `POST` | `/import`、`/export` | `POST /api/v1/system/users/import` |

::: warning 命名要稳定
既然写操作统一走 `POST`，路径就必须把动作表达清楚。不要出现同一个路径靠请求体里的 `action` 字段分发不同业务的设计。
:::

## 和 RESTful 的关系

这个决策不是说 RESTful 不规范，而是本项目选择另一种更适合后台系统的规范。

可以这样理解：

| 风格 | 核心关注点 | 更适合 |
| --- | --- | --- |
| RESTful | HTTP 方法表达资源操作语义 | 开放 API、平台 API、资源服务 |
| 企业后台风格 | 路径表达业务动作，写操作统一 `POST` | 管理后台、内网系统、权限审计较重的系统 |

两种风格都可以规范，关键是团队内部要统一。

## 推荐路径格式

后台管理模块默认采用下面的路径结构：

```text
/api/v1/system/<resource>
/api/v1/system/<resource>/:id/update
/api/v1/system/<resource>/:id/status
/api/v1/system/<resource>/:id/delete
/api/v1/system/<resource>/batch-delete
```

关系类操作放在主资源下面：

```text
POST /api/v1/system/users/:id/roles
POST /api/v1/system/roles/:id/menus
POST /api/v1/system/roles/:id/permissions
```

动作类接口直接使用动作名：

```text
POST /api/v1/auth/login
POST /api/v1/auth/logout
POST /api/v1/system/permissions/reload
```

## 需要避免的写法

不建议把所有写操作都塞进一个模糊入口：

```text
POST /api/v1/system/users/action
```

也不建议用请求体里的 `action` 决定真实业务：

```json
{
  "action": "delete",
  "id": 1
}
```

这种写法会让权限控制、日志审计、接口文档和问题排查都变得不清楚。
