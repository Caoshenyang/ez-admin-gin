# 系统分析报告

> 本文档基于代码实际状态生成，用于指导文档体系重建。代码是唯一事实来源。

---

## 1. 当前项目真实能力

### 后端能力

| 能力 | 实现状态 | 关键路径 |
|------|---------|---------|
| JWT 登录认证 | ✅ 已完成 | `internal/platform/authn/` |
| RBAC 角色权限 | ✅ 已完成 | `internal/modules/iam/role/` |
| Casbin 接口授权 | ✅ 已完成 | `internal/platform/authz/` |
| 动态菜单（三级：目录/菜单/按钮） | ✅ 已完成 | `internal/modules/iam/menu/` |
| 部门树管理 | ✅ 已完成 | `internal/modules/iam/department/` |
| 岗位管理 | ✅ 已完成 | 代码存在，模型为 `sys_post` |
| 用户-角色-岗位关联 | ✅ 已完成 | `user_role`、`user_post` 关联表 |
| 数据权限（五级作用域） | ✅ 已完成 | `internal/platform/datascope/` |
| 用户管理 | ✅ 已完成 | `internal/modules/iam/user/` |
| 角色管理 | ✅ 已完成 | `internal/modules/iam/role/` |
| 菜单管理 | ✅ 已完成 | `internal/modules/iam/menu/` |
| 系统配置 | ✅ 已完成 | `internal/modules/system/config/` |
| 数据字典 | ✅ 已完成 | `internal/modules/system/dict/` |
| 文件上传 | ✅ 已完成 | `internal/modules/system/file/` |
| 附件管理 | ✅ 已完成 | `internal/modules/system/attachment/` |
| 操作日志 | ✅ 已完成 | `internal/modules/system/operationlog/` |
| 登录日志 | ✅ 已完成 | `internal/modules/system/loginlog/` |
| 公告管理 | ✅ 已完成 | `internal/modules/system/notice/` |
| 账户中心（个人信息/密码修改） | ✅ 已完成 | `internal/modules/auth/` |
| Dashboard 数据接口 | ✅ 已完成 | `internal/modules/auth/` |
| 系统初始化 (Setup) | ✅ 已完成 | `internal/modules/setup/` |
| 请求级操作日志 | ✅ 已完成 | `internal/platform/middleware/operation_log.go` |
| 登录限流 | ✅ 已完成 | `internal/platform/middleware/ratelimit.go` |
| CORS | ✅ 已完成 | `internal/platform/middleware/cors.go` |
| RequestID | ✅ 已完成 | `internal/platform/middleware/requestid.go` |
| Swagger | ✅ 已完成 | `internal/bootstrap/swagger.go` |

### 前端能力

| 能力 | 实现状态 | 关键路径 |
|------|---------|---------|
| 登录页 | ✅ 已完成 | `modules/auth/pages/LoginPage.vue` |
| Dashboard | ✅ 已完成 | `modules/auth/pages/DashboardHome.vue` |
| 账户中心 | ✅ 已完成 | `modules/auth/pages/AccountCenterPage.vue` |
| 后台壳子（侧边栏+头部+内容区+标签页） | ✅ 已完成 | `layouts/AdminLayout.vue` |
| 动态菜单渲染 | ✅ 已完成 | `router/dynamic-menu.ts` |
| 工作标签页 | ✅ 已完成 | `components/app-shell/` |
| 按钮权限控制 | ✅ 已完成 | `composables/usePermission.ts` |
| 用户管理页 | ✅ 已完成 | `modules/iam/pages/UserView.vue` |
| 角色管理页 | ✅ 已完成 | `modules/iam/pages/RoleView.vue` |
| 菜单管理页 | ✅ 已完成 | `modules/iam/pages/MenuView.vue` |
| 部门管理页 | ✅ 已完成 | `modules/iam/pages/DepartmentView.vue` |
| 岗位管理页 | ✅ 已完成 | `modules/iam/pages/PostView.vue` |
| 字典管理页 | ✅ 已完成 | `modules/system/pages/DictView.vue` |
| 配置管理页 | ✅ 已完成 | `modules/system/pages/ConfigView.vue` |
| 文件管理页 | ✅ 已完成 | `modules/system/pages/FileView.vue` |
| 操作日志页 | ✅ 已完成 | `modules/system/pages/OperationLogView.vue` |
| 登录日志页 | ✅ 已完成 | `modules/system/pages/LoginLogView.vue` |
| 通知公告页 | ✅ 已完成 | `modules/system/pages/NoticeView.vue` |

### 部署能力

| 能力 | 实现状态 | 关键路径 |
|------|---------|---------|
| Docker Compose 本地开发 | ✅ 已完成 | `deploy/compose.local.yml` |
| Docker Compose Windows 适配 | ✅ 已完成 | `deploy/compose.local.win.yml` |
| 服务器部署（二进制 + Docker 基础设施） | ✅ 已完成 | `deploy/compose.server.yml` |
| 云端部署（全容器化） | ✅ 已完成 | `deploy/compose.deploy.yml` |
| 生产部署 | ✅ 已完成 | `deploy/compose.prod.yml` |
| Nginx HTTP/HTTPS 配置 | ✅ 已完成 | `deploy/nginx*.conf` |
| 环境变量模板 | ✅ 已完成 | `deploy/.env.example` |
| 一键部署脚本 | ✅ 已完成 | `scripts/deploy.sh` |
| 打包脚本（Linux/macOS/Windows） | ✅ 已完成 | `scripts/pack.sh`, `scripts/pack.ps1` |
| 服务器初始化脚本 | ✅ 已完成 | `scripts/setup-server.sh` |
| 更新脚本 | ✅ 已完成 | `scripts/update-server.sh` |
| Systemd 服务配置 | ✅ 已完成 | `deploy/ez-admin.service` |
| GitHub Actions 文档部署 | ✅ 已完成 | `.github/workflows/` |
| 迁移验证脚本 | ✅ 已完成 | `scripts/verify-realdb-migrations.sh` |

---

## 2. 当前系统模块

### 后端模块

```
server/internal/
├── modules/
│   ├── auth/          认证模块：登录、当前用户、菜单、Dashboard、账户中心
│   ├── iam/           身份与访问管理
│   │   ├── user/      用户管理
│   │   ├── role/      角色管理
│   │   ├── menu/      菜单管理
│   │   └── department/ 部门管理
│   ├── system/        系统模块
│   │   ├── config/    系统配置
│   │   ├── dict/      数据字典
│   │   ├── notice/    通知公告
│   │   ├── operationlog/ 操作日志
│   │   ├── loginlog/  登录日志
│   │   ├── file/      文件上传
│   │   └── attachment/ 附件管理
│   ├── setup/         系统初始化
│   └── modulekit/     模块工具包（公共注册接口）
├── platform/          平台层（跨模块基础设施）
│   ├── authn/         JWT 认证
│   ├── authz/         Casbin 授权
│   ├── config/        配置管理
│   ├── database/      数据库连接
│   ├── datascope/     数据权限
│   ├── logger/        日志
│   ├── middleware/     中间件
│   ├── migrate/       迁移管理
│   ├── model/         GORM 模型（18 个表模型）
│   └── redis/         Redis 连接
└── pkg/               公共工具包
    ├── errorsx/       错误处理
    ├── httpx/         HTTP 工具
    ├── actorx/        Actor 上下文
    └── paging/        分页
```

### 前端模块

```
admin/src/
├── modules/
│   ├── auth/          认证模块：登录页、Dashboard、账户中心
│   ├── iam/           IAM 模块：用户、角色、菜单、部门、岗位管理页
│   └── system/        系统模块：配置、字典、文件、日志、公告管理页
├── components/        全局组件
│   ├── app-shell/     布局壳子（侧边栏、头部、标签页）
│   └── brand/         品牌 Logo
├── layouts/           布局
├── router/            路由（含动态菜单注册）
├── stores/            Pinia 状态管理
├── composables/       组合式函数（含 usePermission）
├── api/               HTTP 客户端封装
├── styles/            全局样式
├── constants/         常量
├── types/             全局类型
└── utils/             工具函数（含 auth 工具）
```

---

## 3. 后端架构分层

### 分层结构

```
handler (api/)    → HTTP 请求处理、参数绑定、响应序列化
service           → 业务逻辑编排、权限校验、跨模块协调
repository (infra/) → 数据访问、GORM 查询、数据权限注入
domain            → 领域模型、常量、业务规则
```

### 模块内部结构（以 user 为例）

```
internal/modules/iam/user/
├── api/
│   ├── handlers.go    路由处理器
│   ├── dto.go         请求/响应 DTO
│   └── routes.go      路由注册
├── application/
│   └── user.service.go  业务逻辑
├── infra/
│   └── repository.go   数据访问
└── domain/
    └── types.go        领域类型和常量
```

### 请求链路

```
HTTP Request
  → gin.Engine
  → CORS Middleware
  → RequestID Middleware
  → Logger Middleware
  → Recovery Middleware
  → Route Group (/api/v1)
  → Auth Middleware (JWT 验证)
  → LoadActor Middleware (加载用户上下文+权限)
  → Permission Middleware (Casbin 策略检查)
  → OperationLog Middleware (操作日志记录)
  → Handler (参数绑定)
  → Service (业务逻辑)
  → Repository (数据访问，含数据权限过滤)
  → Response (统一格式)
```

### 统一响应格式

```json
{
  "code": 0,
  "message": "ok",
  "data": { ... }
}
```

错误码体系：`0` 成功、`40000` 请求错误、`40100` 未认证、`40300` 无权限、`40400` 未找到、`50000` 内部错误。

### 配置系统

- 配置文件：`server/configs/config.yaml`
- 环境变量覆盖：前缀 `EZ_`，如 `EZ_AUTH_JWT_SECRET`
- 配置结构体：`internal/platform/config/`
- 支持的配置段：app、server、database、redis、log、auth、swagger、cors、rate_limit、upload

---

## 4. 前端架构分层

### 技术栈

| 技术 | 版本 | 用途 |
|------|------|------|
| Vue | 3.5.32 | UI 框架 |
| Naive UI | 2.44.1 | 组件库 |
| Pinia | 3.0.4 | 状态管理 |
| Vite | 8.0.8 | 构建工具 |
| Tailwind CSS | 4.2.4 | CSS 工具类 |
| Axios | 1.15.2 | HTTP 客户端 |
| TypeScript | 全量使用 | 类型安全 |

### 模块分层约定

每个业务模块遵循固定三层职责：

```
modules/{module}/
├── api/            接口调用，只做 HTTP 请求和类型转换
├── types/          TypeScript 类型定义
├── composables/    状态管理 + 副作用逻辑（useXxxPage）
├── components/     展示组件，通过 props/events 通信，不导入 api
└── pages/          编排层，只做拼装和路由绑定
```

### API 层

- HTTP 客户端封装：`admin/src/api/http.ts`
- 请求拦截器：自动注入 `Authorization: Bearer <token>`
- 响应拦截器：401/403 自动跳转登录页
- API 代理：Vite 开发模式代理 `/api` → `http://localhost:8080`

### 状态管理

Pinia store `admin-shell` 管理：
- 活跃菜单项
- 展开的菜单组
- 工作标签列表
- 侧边栏折叠状态
- 路由刷新 nonce

---

## 5. 权限链路

### 完整权限链路

```
1. 用户登录 → JWT Token 签发
   authn.GenerateToken(userID, username)

2. 请求携带 Token → Auth 中间件验证
   middleware/auth.go: 解析 JWT，提取 userID + username

3. LoadActor 中间件 → 加载用户完整上下文
   middleware/auth.go#LoadActor:
   - 根据 userID 查询用户信息
   - 查询用户角色列表
   - 查询角色关联的菜单 ID 集合
   - 查询角色关联的按钮权限码集合
   - 将上下文注入 gin.Context

4. Permission 中间件 → Casbin 策略匹配
   middleware/permission.go:
   - 从 Actor 上下文获取角色列表
   - 构造 Casbin 请求 (角色, URL路径, HTTP方法)
   - 逐一匹配 Casbin 策略
   - 任一角色匹配即放行

5. 数据权限 → datascope 在 Repository 层注入
   - 根据 Actor 的角色数据范围类型
   - 注入不同的 GORM scope 过滤数据
```

### Casbin RBAC 模型

```
[request_definition]
r = sub, obj, act

[policy_definition]
p = sub, obj, act

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = r.sub == p.sub && keyMatch2(r.obj, p.obj) && (r.act == p.act || p.act == "*")
```

- sub = 角色编码（如 `super_admin`）
- obj = URL 路径模式（如 `/api/v1/iam/users/*`）
- act = HTTP 方法（GET/POST/PUT/DELETE，`*` 表示所有）
- 策略存储在 `casbin_rule` 数据库表

### 按钮权限

- 菜单表中 `type=3` 的记录为按钮权限
- 每个按钮权限有唯一 `code`（如 `system:user:create`）
- 前端通过 `usePermission().canUse(code)` 判断按钮是否显示
- 后端在 Casbin 中同时注册按钮级 API 权限

---

## 6. 动态菜单机制

### 后端菜单结构

```
sys_menu
├── id
├── parent_id       父菜单 ID（树结构）
├── type             1=目录, 2=菜单, 3=按钮
├── code             唯一权限编码（如 system:user:list）
├── title            显示名称
├── path             前端路由路径
├── component        前端组件路径
├── icon             菜单图标
├── sort             排序值
├── status           状态（启用/禁用）
├── permission       后端权限标识
├── created_at
├── updated_at
└── deleted_at
```

### 菜单加载流程

```
1. 前端登录成功 → 调用 GET /api/v1/auth/menus
2. 后端根据用户角色 → 查询关联的菜单 ID 集合
3. 过滤 type=1,2 的菜单 → 构建菜单树
4. 提取 type=3 的按钮 → 返回按钮权限码列表
5. 前端收到菜单树 → 动态注册 Vue Router 路由
6. 前端收到按钮权限码 → 存入 usePermission composable
7. NMenu 组件渲染菜单树 → 侧边栏展示
8. 页面中通过 canUse(code) 控制按钮显隐
```

### 动态路由注册

```typescript
// admin/src/router/dynamic-menu.ts
- 遍历后端返回的菜单树
- 根据 component 字段匹配前端组件（白名单映射）
- 使用 router.addRoute() 动态添加路由
- type=1(目录) → 嵌套路由父级
- type=2(菜单) → 具体页面路由
- type=3(按钮) → 提取权限码
```

---

## 7. 数据权限机制

### 五级数据作用域

| 作用域 | 值 | 含义 |
|--------|---|------|
| all | 1 | 查看所有数据 |
| dept | 2 | 仅本部门数据 |
| dept_and_children | 3 | 本部门及下级部门数据 |
| self | 4 | 仅本人数据 |
| custom_dept | 5 | 自定义部门范围 |

### 数据权限注入点

```
Repository 层 → datascope.ApplyScopes(db, actor)
```

- 根据 Actor 的角色数据范围类型，注入不同的 GORM scope
- `all` → 不追加过滤条件
- `dept` → `WHERE department_id = actor.department_id`
- `dept_and_children` → `WHERE department_id IN (本部门及所有子部门ID)`
- `self` → `WHERE creator_id = actor.userID`
- `custom_dept` → `WHERE department_id IN (sys_role_data_scope 配置的部门ID列表)`

### 数据权限相关表

- `sys_role.data_scope`：角色数据范围类型
- `sys_role_data_scope`：自定义部门范围关联表
- `sys_department.ancestors`：部门祖先路径（逗号分隔），用于快速查询子树

---

## 8. Migration 与 Seed 初始化机制

### 迁移文件（9 个版本）

```
server/migrations/
├── mysql/
│   ├── 000001_init_schema.up.sql           基础表结构
│   ├── 000002_seed_data.up.sql             种子数据
│   ├── 000003_enterprise_foundation.up.sql  企业基础（部门/岗位/数据权限）
│   ├── 000004_phase4_dict_schema.up.sql     字典表结构
│   ├── 000005_phase4_dict_seed_data.up.sql  字典种子数据
│   ├── 000006_phase4_attachment_schema.up.sql 附件表结构
│   ├── 000007_phase4_attachment_seed_data.up.sql 附件种子数据
│   ├── 000008_phase4_org_menu_seed_data.up.sql  组织菜单种子数据
│   └── 000009_phase4_menu_icon_alignment.up.sql 菜单图标对齐
└── postgres/
    └── (同结构，PostgreSQL 方言)
```

每个迁移都有对应的 `.down.sql` 回滚文件。

### 种子数据内容

- `super_admin` 角色（全权限）
- 完整的系统管理菜单树（目录 + 菜单 + 按钮）
- 所有 API 权限分配给 `super_admin`
- 角色-菜单绑定
- 字典类型和字典项
- 附件配置

### 数据库表清单（18 张表）

| 表名 | 用途 |
|------|------|
| `sys_user` | 用户 |
| `sys_role` | 角色 |
| `sys_menu` | 菜单/权限 |
| `sys_department` | 部门 |
| `sys_post` | 岗位 |
| `sys_user_role` | 用户-角色关联 |
| `sys_role_menu` | 角色-菜单关联 |
| `sys_user_post` | 用户-岗位关联 |
| `sys_role_data_scope` | 角色自定义数据范围 |
| `sys_config` | 系统配置 |
| `sys_dict_type` | 字典类型 |
| `sys_dict_item` | 字典项 |
| `sys_login_log` | 登录日志 |
| `sys_operation_log` | 操作日志 |
| `sys_notice` | 通知公告 |
| `sys_file` | 文件 |
| `sys_attachment` | 附件 |
| `casbin_rule` | Casbin 策略 |

### 初始化流程

```
1. 首次启动 → Setup API (POST /api/v1/setup/init)
2. 创建管理员用户
3. 迁移通过 golang-migrate 自动执行
4. 种子数据在迁移 SQL 中写入
5. 后续启动自动跳过已执行的迁移
```

---

## 9. 部署结构

### 部署变体

| 变体 | 配置文件 | 适用场景 |
|------|---------|---------|
| 本地开发 | `compose.local.yml` | macOS/Linux 开发环境 |
| Windows 开发 | `compose.local.win.yml` | Windows 开发环境 |
| 服务器部署 | `compose.server.yml` | 二进制运行在宿主机，基础设施容器化 |
| 云端部署 | `compose.deploy.yml` | 全容器化，Docker Hub 镜像 |
| 生产部署 | `compose.prod.yml` | 生产环境 |

### Nginx 配置

- `nginx.conf`：HTTP 反向代理
- `nginx-ssl.conf`：HTTPS + SSL 终止
- `nginx-native.conf`：宿主机原生 Nginx
- `nginx-native-ssl.conf`：宿机原生 Nginx HTTPS

所有配置包含：SPA fallback、静态资源 1 年缓存、Gzip 压缩、安全头。

### 部署脚本

| 脚本 | 用途 |
|------|------|
| `scripts/deploy.sh` | 一键部署到远程服务器 |
| `scripts/pack.sh` | 构建并打包（Linux/macOS） |
| `scripts/pack.ps1` | 构建并打包（Windows） |
| `scripts/setup-server.sh` | 服务器初始化 |
| `scripts/update-server.sh` | 更新已部署服务 |

### 环境变量

- 模板文件：`deploy/.env.example`
- 覆盖机制：环境变量 > config.yaml
- 前缀：`EZ_`
- 关键变量：数据库连接、Redis、JWT 密钥、CORS、限流

---

## 10. 当前文档与代码不一致的地方

### 严重不一致

1. **文档描述为"从零搭建教程"，代码已是一个完整的成品系统**
   - 整个 `docs/tutorial/` 目录（8 章 + 归档）都是教学式内容
   - 代码已经是成熟的企业后台底座，不需要教程

2. **README 声称"9 章"教程，VitePress 配置只显示 8 章**
   - README 写"当前教程主线已经稳定为 **9 章**"
   - 实际 config.mts 侧边栏只有 chapter-1 到 chapter-8

3. **文档定位与代码能力不匹配**
   - 文档面向"Java 转 Go 工程师"的教学场景
   - 代码已经是一个可直接使用的后台底座
   - 应该转向"项目介绍 + 架构设计 + 二次开发"定位

4. **部署文档散落在 tutorial/chapter-8 和 deploy/ 两处**
   - 部署侧边栏直接链接到 `tutorial/chapter-8/` 的页面
   - `deploy/index.md` 存在但内容可能只是跳转页

5. **guide/ 下的"部署篇页面规划"和"文档基线盘点"是内部规划文档**
   - `guide/deployment-document-plan.md` — 内部规划文档，不应面向用户
   - `guide/document-baseline-inventory.md` — 文档盘点清单，内部用途
   - `guide/execution-plan.md` — 执行计划，内部用途

6. **reference/ 中存在过时的框架入门文档**
   - `gorm-quick-start.md`、`casbin-quick-start.md` — 应该是外部文档链接，不需要自建

### 中度不一致

7. **VitePress 描述仍然是"面向 Java 转 Go 工程师"**
   - `description: '面向 Java 转 Go 工程师的企业级通用后台管理系统底座。'`
   - 应该改为更通用的项目定位

8. **README 缺少重要信息**
   - 没有 License 标识（只写了 MIT，没有 LICENSE 文件）
   - 没有截图
   - 没有 Contributing 指引
   - 没有 Roadmap 链接指向具体内容

9. **tutorial/chapter-5 中混入了前端页面内容**
   - `login-page.md`、`user-pages.md`、`role-menu-pages.md`、`config-file-pages.md`、`log-pages.md`、`dynamic-menu.md`、`admin-layout.md`、`vue-project-init.md` — 这些是前端页面教程
   - 第 5 章标题是"组织体系与数据权限"，但内容包含前端实现

10. **archive/ 下的文档仍然占用 VitePress 构建资源**
    - 6 个归档文件，可能不需要继续构建

---

## 11. 哪些旧文档应该删除

### 应该整体删除

| 路径 | 原因 |
|------|------|
| `docs/tutorial/` (全部) | 教学式教程，与项目定位冲突。部分有价值的架构/设计内容需要提取到新文档中后再删除 |
| `docs/guide/deployment-document-plan.md` | 内部规划文档 |
| `docs/guide/document-baseline-inventory.md` | 内部盘点文档 |
| `docs/guide/execution-plan.md` | 内部执行计划 |

### 可以提取内容后删除

以下教程页面包含有价值的架构/设计说明，应提取到新文档后删除：

| 当前路径 | 可提取的内容 | 目标位置 |
|----------|-------------|---------|
| `tutorial/chapter-5/organization-model-design.md` | 组织模型设计 | `architecture/organization.md` |
| `tutorial/chapter-5/role-data-scope-and-query-scopes.md` | 数据权限模型 | `architecture/data-permission.md` |
| `tutorial/chapter-5/actor-and-grant-merge.md` | Actor 上下文 | `architecture/data-permission.md` |
| `tutorial/chapter-5/request-flow-walkthrough.md` | 请求权限走读 | `architecture/request-flow.md` |
| `tutorial/chapter-4/rbac-model.md` | RBAC 模型 | `architecture/rbac.md` |
| `tutorial/chapter-4/casbin-permission.md` | Casbin 接入 | `backend/casbin.md` |
| `tutorial/chapter-6/module-structure.md` | 模块结构 | `backend/module-structure.md` |
| `tutorial/chapter-6/backend-module-flow.md` | 模块接入流程 | `backend/module-development.md` |
| `tutorial/chapter-7/dynamic-route-registration.md` | 动态路由 | `frontend/dynamic-menu.md` |
| `tutorial/chapter-7/login-and-session-flow.md` | 登录态流转 | `frontend/auth-flow.md` |
| `tutorial/chapter-8/compose-and-service-layout.md` | Compose 结构 | `deployment/compose.md` |
| `tutorial/chapter-8/nginx-and-https.md` | Nginx 配置 | `deployment/nginx.md` |

### 可以保留并调整的

| 路径 | 处理方式 |
|------|---------|
| `docs/guide/index.md` | 重写为 Getting Started |
| `docs/guide/project-structure.md` | 更新为当前结构 |
| `docs/guide/enterprise-architecture.md` | 调整为架构设计文档 |
| `docs/guide/java-to-go-structure.md` | 保留，移至 reference/ |
| `docs/guide/changelog.md` | 保留，移至 reference/ |
| `docs/guide/roadmap.md` | 保留，更新 |
| `docs/reference/` 大部分文件 | 保留，按需更新 |

### 应该保留不变的 reference 文件

| 文件 | 原因 |
|------|------|
| `reference/api-style-decision.md` | API 约定，仍有参考价值 |
| `reference/permission-code-conventions.md` | 权限码约定 |
| `reference/error-code-reference.md` | 错误码 |
| `reference/database-ddl.md` | 数据库 DDL |
| `reference/environment-variables-reference.md` | 环境变量 |
| `reference/directory-conventions.md` | 目录约定 |
| `reference/module-conventions.md` | 模块规范 |
| `reference/query-and-pagination-conventions.md` | 查询分页 |
| `reference/data-scope-model.md` | 数据权限模型 |
| `reference/deploy-artifacts-reference.md` | 部署文件参考 |
| `reference/nginx-config-reference.md` | Nginx 参考 |
| `reference/init-data-reference.md` | 初始化数据参考 |
| `reference/dynamic-menu-component-reference.md` | 菜单组件白名单 |
| `reference/button-permission-consumption-example.md` | 按钮权限示例 |
| `reference/upload-public-path-reference.md` | 上传路径参考 |
| `reference/module-init-template.md` | 模块模板 |
| `reference/migration-tool-selection.md` | 迁移工具选型 |
| `reference/logical-delete-and-unique-index.md` | 逻辑删除 |
| `reference/ssh-tunnel-database.md` | SSH 隧道 |
| `reference/vitepress-github-pages.md` | 文档部署 |

### 可考虑删除的 reference 文件

| 文件 | 原因 |
|------|------|
| `reference/gorm-quick-start.md` | 外部框架教程，链接到官方文档即可 |
| `reference/casbin-quick-start.md` | 外部框架教程，链接到官方文档即可 |

---

*本分析基于 `enterprise-foundation-v2` 分支代码状态，生成时间 2026-05-10。*
