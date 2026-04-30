# EZ Admin Gin Enterprise Foundation V2 Handoff

## 1. 当前基线

- 仓库路径：`D:\A\ez-admin-gin`
- 当前分支：`enterprise-foundation-v2`
- 最近提交：`f7e97a5`
- 最近提交信息：`feat: reshape enterprise foundation and tutorial`

这轮工作的目标，不再是继续维护一个“可用但偏 demo 化”的后台模板，而是把项目直接收敛为：

- 面向 Java 转 Go 工程师的企业级通用后台管理系统底座
- 与最终版代码结构一致的企业级完整版 0-1 教程

核心原则已经明确：

- 不再写“简化版 -> 企业版”的迁移叙事
- 教程直接讲最终结构
- 代码和教程同步推进
- 保持单体架构，不拆微服务

## 2. 当前已经完成的内容

### 2.1 后端总体结构

后端已经从旧的扁平结构开始收敛到最终结构：

```text
server/
├── cmd/server/
├── internal/bootstrap/
├── internal/platform/
├── internal/module/
└── migrations/
```

已经落地的关键目录：

- `server/cmd/server/main.go`
- `server/internal/bootstrap/router.go`
- `server/internal/bootstrap/run.go`
- `server/internal/platform/authn`
- `server/internal/platform/authz`
- `server/internal/platform/config`
- `server/internal/platform/database`
- `server/internal/platform/datascope`
- `server/internal/platform/logger`
- `server/internal/platform/migrate`
- `server/internal/platform/redis`

### 2.2 组织体系与数据权限基础

已经完成：

- `Actor` 上下文加载
- 数据权限基础模型与聚合逻辑
- 部门、岗位、用户岗位、角色自定义部门范围模型
- 用户模型 `department_id`
- 角色模型 `data_scope`
- PostgreSQL / MySQL 企业级升级迁移脚本

关键文件：

- `server/internal/middleware/actor.go`
- `server/internal/platform/datascope/datascope.go`
- `server/internal/model/department.go`
- `server/internal/model/post.go`
- `server/internal/model/user_post.go`
- `server/internal/model/role_data_scope.go`
- `server/internal/model/user.go`
- `server/internal/model/role.go`
- `server/migrations/postgres/000003_enterprise_foundation.*`
- `server/migrations/mysql/000003_enterprise_foundation.*`

### 2.3 已进入最终结构的模块

#### 认证模块

已经收进 `server/internal/module/auth/`：

- `login`
- `me`
- `menus`
- `dashboard`

认证模块现在已经不再只是空壳路由聚合，已经拥有：

- `dto.go`
- `repository.go`
- `login_service.go`
- `me_service.go`
- `menus_service.go`
- `dashboard_service.go`
- 对应 handler
- `routes.go`

#### IAM 模块

已经进入最终结构：

- `server/internal/module/iam/user`
- `server/internal/module/iam/role`
- `server/internal/module/iam/department`
- `server/internal/module/iam/post`
- `server/internal/module/iam/menu`

当前用户模块还额外完成了：

- `post_ids` 正式接入用户接口链路
- 用户与岗位关系进入前后端真实能力

#### System 模块

已经进入最终结构：

- `server/internal/module/system/config`
- `server/internal/module/system/file`
- `server/internal/module/system/operationlog`
- `server/internal/module/system/loginlog`
- `server/internal/module/system/notice`
- `server/internal/module/system/routes.go`

说明：

- 操作日志已经明确为“中间件写入 + 模块查询”
- 登录日志已经明确为“登录入口写入 + 模块查询”
- 公告模块已经作为轻量内容管理模块进入统一边界

### 2.4 路由聚合现状

当前真正应该继续沿着走的是：

- `server/internal/module/auth/routes.go`
- `server/internal/module/system/routes.go`

旧目录中的 `server/internal/router/router.go` 和部分 legacy handler 仍然存在，但当前主线已经明显转向 `module/*`。

后续继续推进时，应优先扩最终结构，不要再扩大旧的扁平 handler 路径。

### 2.5 前端已完成事项

已经做过的前端改动：

- 用户管理页接入岗位多选和岗位展示
- 新增岗位 API 与类型

关键文件：

- `admin/src/pages/system/UserView.vue`
- `admin/src/api/post.ts`
- `admin/src/types/post.ts`
- `admin/src/types/user.ts`

### 2.6 文档已完成事项

#### Guide / Reference

已新增或重写：

- `docs/guide/enterprise-architecture.md`
- `docs/guide/java-to-go-structure.md`
- `docs/guide/execution-plan.md`
- `docs/reference/data-scope-model.md`

#### 教程主线

已完成或明显收敛的页面：

- `docs/tutorial/index.md`
- `docs/tutorial/curriculum.md`
- `docs/tutorial/chapter-3/index.md`
- `docs/tutorial/chapter-3/user-model-and-login.md`
- `docs/tutorial/chapter-3/jwt-auth.md`
- `docs/tutorial/chapter-3/auth-middleware.md`
- `docs/tutorial/chapter-3/menu-permission.md`
- `docs/tutorial/chapter-4/index.md`
- `docs/tutorial/chapter-4/user-management.md`
- `docs/tutorial/chapter-4/role-management.md`
- `docs/tutorial/chapter-4/menu-management.md`
- `docs/tutorial/chapter-4/system-config.md`
- `docs/tutorial/chapter-4/file-upload.md`
- `docs/tutorial/chapter-4/operation-logs.md`
- `docs/tutorial/chapter-4/login-logs.md`
- `docs/tutorial/chapter-4/notice-management.md`
- `docs/tutorial/chapter-5/organization-model-design.md`
- `docs/tutorial/chapter-5/role-data-scope-and-query-scopes.md`
- `docs/tutorial/chapter-5/department-tree-and-management.md`
- `docs/tutorial/chapter-5/post-management-and-user-affiliation.md`
- `docs/tutorial/chapter-8/index.md`
- `docs/tutorial/chapter-9/index.md`

说明：

- 第 3 章最关键的认证主线已经改成最终版叙事
- 第 4 章高频系统模块主线已经基本完整
- 第 5 章已经有数据权限核心正文，但还没有完全收透

## 3. 当前还没完成的内容

### 3.1 代码层未完成

仍然未完成或未彻底收口的事项：

- 更稳定的 `ScopeResolver` 抽象尚未单独沉淀
- 部门、岗位、真实业务模块的数据权限接入规范还需要进一步统一
- `data dict / account / attachment center / business sample module` 还没做
- 关键链路自动化测试还没补
- 旧的部分 legacy handler / router 仍留在仓库中，虽然主线已切，但尚未彻底清理

### 3.2 文档层未完成

仍然未完成或需要继续重写的重点：

- `docs/tutorial/chapter-3/rbac-model.md`
- `docs/tutorial/chapter-3/casbin-permission.md`
- 第 5 章剩余页面继续深化
- 第 6 章正文重写
- 第 7 章正文重写
- 第 8、9 章补实内容，而不只是章节入口
- 参考手册继续补齐：
  - 错误码
  - 权限码规范
  - 环境变量
  - 模块接入清单

## 4. 当前建议的继续顺序

为了降低换电脑后继续开发时的认知负担，建议按下面顺序继续：

### Step 1：先继续收教程第 3 章

优先处理：

- `docs/tutorial/chapter-3/rbac-model.md`
- `docs/tutorial/chapter-3/casbin-permission.md`

原因：

- 认证主线已经重写到一半
- 如果第 3 章继续保留旧写法，会和已经完成的 `module/auth` 形成明显错位

### Step 2：继续补第 5 章数据权限正文

优先方向：

- 把 `Actor`
- `gorm.Scopes(...)`
- 多角色并集
- 部门范围过滤

写成一条真正能跟做的主线，而不是只停留在概念页

### Step 3：继续做企业常用模块

建议顺序：

1. 数据字典
2. 账户中心
3. 附件中心
4. 非 `system` 分组的真实业务示例模块

### Step 4：补测试

优先测试场景：

- 登录成功 / 失败
- 接口权限拦截
- 数据权限过滤
- 部门树
- 初始化迁移

## 5. 下一步最推荐直接做什么

如果换电脑后要无缝接上，最推荐的“下一步第一件事”是：

### 直接重写这两页

1. `docs/tutorial/chapter-3/rbac-model.md`
2. `docs/tutorial/chapter-3/casbin-permission.md`

原因：

- 这是当前最容易继续保持主线一致的地方
- 风险小，不需要先动数据库
- 可以继续把第 3 章整体收成最终版叙事

如果第 3 章收完，再继续进入：

3. `docs/tutorial/chapter-5/` 深化数据权限正文

## 6. 开发时的注意事项

### 6.1 代码侧

- 继续坚持 `module/*` 最终结构，不要再往旧的全局 handler 里扩
- 手动修改文件时继续使用补丁方式，避免无意覆盖大块代码
- 注释保持“解释职责和边界”，不要写无意义注释

### 6.2 文档侧

- 教程直接讲最终结构
- 不再写“从旧结构迁移到新结构”的过渡文案
- 优先保留这套叙事节奏：
  - 本节目标
  - 当前边界
  - 模块职责
  - 为什么这样设计
  - 怎么验证

### 6.3 本地文件注意

当前工作区里有一个未提交本地文件：

- `.claude/settings.local.json`

它没有被纳入提交，后续继续开发时也不应该把它误提交进版本库。

## 7. 当前可用验证命令

后续继续开发时，至少保持这两条验证命令常跑：

### 后端

```bash
cd server
go test ./...
```

### 文档

```bash
cd docs
pnpm docs:build
```

如果涉及前端页面改动，再补：

```bash
cd admin
pnpm build
```

## 8. 当前状态一句话总结

当前项目已经完成了“企业级底座 v2”最难的第一段：

- 代码结构已经开始真正脱离 demo 形态
- 组织体系和数据权限基础已经落地
- 高频系统模块已经大面积进入最终结构
- 教程主线已经开始按最终代码结构重写

接下来最重要的，不是再做零散新功能，而是：

> 把第 3 章和第 5 章继续收透，再在这个稳定基线之上补齐企业常用模块和自动化测试。
