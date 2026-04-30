---
title: Go vs Java 工程结构
description: "帮助 Java 工程师理解 EZ Admin Gin 在 v2 阶段采用的 Go 单体后台结构，以及它和常见 Java 分层的映射关系。"
---

# Go vs Java 工程结构

::: tip 🎯 这页解决什么
如果你熟悉的是 `controller / service / repository` 这一套 Java 分层，这页会帮你把 `EZ Admin Gin` 的 Go 结构快速对应起来。
:::

## 先说结论

Go 项目不是不能分层，而是不强调“层越多越专业”。  
更常见的做法是：

- 先按**领域模块**收拢代码
- 再在模块内部按职责拆分
- 技术基础设施统一放到 `platform`

这和 Java 里“按技术层全局分目录”的思路不完全一样。

## 一张对照表先看懂

| Java 常见结构 | 在 `EZ Admin Gin` v2 中对应什么 | 主要职责 |
| --- | --- | --- |
| `Controller` | `handler.go` | 参数绑定、调用 service、返回响应 |
| `Service` | `service.go` | 业务规则、事务边界、跨仓储协作 |
| `Repository / Mapper` | `repository.go` | GORM 查询、持久化、查询拼装 |
| `Entity / DO` | `entity.go` | 数据库实体和状态枚举 |
| `DTO / VO` | `dto.go` | 请求参数和响应结构 |
| `Security` | `platform/authn`、`platform/authz` | 登录态、接口权限 |
| `Data Permission` | `platform/datascope` + 模块内 `datascope.go` | 数据范围解析和资源级过滤 |
| `Spring Boot 启动类` | `cmd/server` + `bootstrap` | 应用启动、依赖装配、模块注册 |

## 为什么不继续用 v1 的结构

`v1` 阶段的 `handler` 直接操作 `gorm.DB`，对个人项目和第一版后台来说很高效。  
但只要开始做这些能力：

- 部门树
- 岗位
- 角色数据范围
- 自定义部门授权
- 查询级数据过滤

`handler` 就会很快同时承担太多职责。

在 Java 项目里，你通常会自然想到把它们拆进 `service` 和 `repository`。  
在 Go 里也是一样，只是拆分的方式更强调：

- 不为了形式而多加层
- 只在复杂度真的上来时再补职责边界

`v2` 正是在这个节点上开始补分层。

## Go 里更推荐“按模块收拢”

很多 Java 项目喜欢全局这样分：

```text
controller/
service/
repository/
entity/
```

Go 在企业单体里更常见的做法是：

```text
module/
  iam/
    user/
      entity.go
      dto.go
      repository.go
      service.go
      handler.go
      routes.go
```

这样做的好处是：

- 一个资源的代码尽量放在一起
- 新增模块时不需要在全局四五个目录来回跳
- 模块边界更清晰
- 后续接入数据权限时，更容易把资源自己的过滤规则放在同一目录

## `platform` 在 Go 里相当于什么

可以把 `platform` 理解成“项目级技术底座”，它大致对应 Java 里这些内容的组合：

- `config`
- `common`
- `infrastructure`
- `security`
- `starter`

但和 Java 不同的是，`platform` 不应该承载业务逻辑。  
它只放这些通用能力：

- 配置
- 日志
- 数据库
- Redis
- 登录态
- 接口权限
- 数据权限基础设施

## 数据权限为什么不能只放中间件

这是 Java 转 Go 时最容易误判的一点。

接口权限适合放中间件，因为它只需要判断：

- 当前请求是否允许进入

但数据权限通常依赖：

- 当前资源怎么归属
- 当前查询按部门过滤还是按创建人过滤
- 当前角色的数据范围是部门树还是自定义部门

这些规则和具体资源强相关，所以更适合放在：

- `service`
- `repository`
- 模块内 `datascope.go`

中间件只负责把“当前登录人是谁”这类上下文准备好，不直接拼业务查询。

## 推荐心智模型

如果你来自 Java，可以先用这个心智模型理解 `EZ Admin Gin` v2：

- `bootstrap` 像应用装配层
- `platform` 像基础设施层
- `module/*/*/handler.go` 像 controller
- `module/*/*/service.go` 像 service
- `module/*/*/repository.go` 像 repository
- `module/*/*/datascope.go` 像每个资源自己的数据权限扩展点

换句话说，差别不是“不分层”，而是：

- Java 更习惯按全局技术层组织
- Go 更适合先按领域收拢，再在模块里补职责边界

## 怎么继续读

- 想看为什么 `v2` 要做这次结构升级：看 [企业级架构升级](/guide/enterprise-architecture)
- 想继续从整体上了解仓库布局：看 [项目结构](/guide/project-structure)
- 想看数据权限模型本身：看 [数据权限模型](/reference/data-scope-model)
