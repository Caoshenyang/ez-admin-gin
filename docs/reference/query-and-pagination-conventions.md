---
title: 查询与分页约定
description: "集中说明 EZ Admin Gin 当前列表接口统一使用的 page / page_size、keyword、status 和模块扩展筛选约定，以及后端归一化边界。"
---

# 查询与分页约定

这页只管快速查阅，不展开某一个具体模块。它回答的是：

> 当前后台列表接口的查询参数到底有哪些共性，分页边界怎么定，哪些筛选字段属于模块自己扩展。

## 当前列表接口的统一骨架

当前系统里大多数列表接口都沿用一套非常稳定的结构：

| 参数 | 作用 |
| --- | --- |
| `page` | 当前页码 |
| `page_size` | 每页条数 |
| `keyword` | 关键字模糊搜索 |
| `status` | 启停状态筛选 |

典型模块包括：

- `user`
- `role`
- `config`
- `file`
- `notice`
- `operationlog`
- `loginlog`

这些模块通常都会在自己的 `dto.go` 里定义：

```go
type ListQuery struct {
	Page     int    `form:"page"`
	PageSize int    `form:"page_size"`
	Keyword  string `form:"keyword"`
	Status   int    `form:"status"`
}
```

## 当前分页参数的固定边界

现在多个模块都把分页归一化成同一套规则：

| 情况 | 最终结果 |
| --- | --- |
| `page < 1` | 修正为 `1` |
| `page_size < 1` | 修正为 `10` |
| `page_size > 100` | 修正为 `100` |

也就是当前默认分页约定是：

- 默认第一页
- 默认每页 10 条
- 单次最多 100 条

::: tip 为什么这条规则值得统一
如果每个模块各自决定默认页大小和上限，前端分页组件、联调习惯和压测边界都会变得混乱。

当前主线把它固定下来，是为了让所有列表接口在使用体验和资源消耗上保持一致。
:::

## 当前分页返回结构

大多数列表返回都会长成这样：

```json
{
  "items": [],
  "total": 0,
  "page": 1,
  "page_size": 10
}
```

这意味着前端分页表格基本可以复用同一套读法：

- `items` 渲染表格
- `total` 驱动分页总数
- `page` 和 `page_size` 回填当前状态

## `keyword` 当前怎么用

当前 `keyword` 不是全文检索，它主要是：

- 对几个常用文本字段做 `LIKE` 模糊匹配

例如：

| 模块 | 常见匹配字段 |
| --- | --- |
| 用户 | `username` / `nickname` |
| 角色 | `code` / `name` |
| 部门 | `name` / `code` |
| 岗位 | `code` / `name` |
| 配置 | `group_code` / `config_key` / `name` |
| 文件 | 文件名等基础字段 |
| 公告 | `title` |

当前实现习惯上会先：

1. `TrimSpace(keyword)`
2. 非空时再拼接 `%keyword%`

所以它更适合：

- 后台管理台常规模糊搜索

而不是：

- 高级检索
- 多字段布尔组合查询

## `status` 当前怎么用

`status` 当前主要服务于“启用 / 禁用”这类枚举筛选。

常见模式是：

- `0` 表示不筛选
- 其他值再转换成模块自己的状态类型

例如：

- `UserStatus`
- `RoleStatus`
- `SystemConfigStatus`
- `NoticeStatus`

如果值不在合法枚举范围内，当前统一返回：

- `40000`
- 对应模块的“状态不正确”提示语

## 模块扩展筛选字段怎么处理

在统一骨架之外，允许模块补自己真正需要的筛选字段。

当前比较典型的例子是：

| 模块 | 扩展字段 | 作用 |
| --- | --- | --- |
| 配置 | `group_code` | 按配置分组筛选 |
| 登录日志 | `status` 语义偏向登录结果 |
| 操作日志 | `success`、时间范围等 | 按成功状态和时间筛选 |
| 部门 | 树形结构下的关键字与状态 | 更偏组织管理场景 |

这说明当前规范并不是“所有列表都只能四个参数”，而是：

- 先有统一骨架
- 再按模块补扩展筛选

## 当前更推荐的后端分工

如果你要新增一个列表接口，当前更稳的分工是：

| 层 | 推荐职责 |
| --- | --- |
| `dto.go` | 定义 `ListQuery`，归一化分页和筛选参数 |
| `service.go` | 决定这次列表查询的业务边界 |
| `repository.go` | 把筛选条件真正翻译成 SQL / GORM 查询 |

也就是说：

- `NormalizePage(...)` 放 DTO
- 状态合法性转换也优先放 DTO
- `Where(...)` 细节放 Repository

## 一个新列表接口最稳的起步模板

如果你要新建一个模块列表，当前最稳的最小模板可以直接照这个思路：

1. `ListQuery` 先放 `page / page_size / keyword / status`
2. 写一个模块自己的 `NormalizePage(...)`
3. 状态筛选单独写 `NormalizeStatusFilter(...)`
4. Repository 里先支持 `keyword + status`
5. 返回统一的 `items / total / page / page_size`

## 当前为什么每个模块都各自写 `NormalizePage`

你会发现现在很多模块各自都有一份 `NormalizePage(...)`，而不是抽成一个全局工具。

当前这样做的现实原因是：

- 逻辑虽然一致
- 但还没正式抽成平台级分页工具

这也意味着当前对外约定已经稳定，但内部仍有进一步收口空间。

## 最常见的查询问题

### 1. 列表不传分页参数会怎样

当前会自动落到：

- `page = 1`
- `page_size = 10`

### 2. 前端一口气要 1000 条怎么办

当前后端会把 `page_size` 压到：

- `100`

### 3. 状态传错为什么直接 400

因为当前 `status` 不是自由文本，而是受模块枚举约束的筛选值。

### 4. 为什么有的模块没有 `page/page_size`

像部门树、菜单树这类更偏树结构的接口，当前通常不会走标准分页骨架，而是直接返回整棵树或整组节点。

## 相关教程与参考页

- [模块规范](./module-conventions)
- [错误码参考](./error-code-reference)
- [第 6 章：核心系统模块](../tutorial/chapter-6/)
- [第 8 章：模块化接入规范](../tutorial/chapter-8/)
