# Process

只跟踪代码闭环，不以 `docs/` 为准。

- 更新时间：`2026-05-08`
- 状态标记：
  - `[x]` 已完成
  - `[-]` 进行中 / 已确认存在问题但未收口
  - `[ ]` 未开始
  - `[!]` 阻塞 / 当前失败，需要优先处理
- 评分摘要：
  - 后端落地度：`8.5 / 10`
  - 前端落地度：`6.7 / 10`
  - 整体代码质量：`7.7 / 10`

## 当前基线状态

- [x] 后端主链路已具备运行闭环：启动、路由装配、认证、RBAC、数据权限、系统模块均已落地
- [x] 后端构建当前通过，启动链路已验证到数据库、迁移、Redis 和路由注册阶段
- [x] 前端 `vue-tsc` 当前通过
- [x] 前端构建验证当前通过
- [x] 代码层面已形成 `server / admin / deploy` 的实际工程分工
- [x] 文档与代码当前脱钩，本阶段不以文档状态作为完成依据
- [-] 前端正在补齐列表空值显示与一屏布局约束；根布局已收口为内容区内部滚动，表格卡片的高度来源已统一补齐，仍需继续人工走查工作台和少量页面余量

## 当前收口判断

- [x] 后端主链路形成闭环
- [x] 前端类型检查通过
- [x] 前端构建通过
- [-] 前端模块结构已形成，但边界仍需收口
- [-] 后端结构成熟，但实现风格仍需统一
- [-] 整体代码闭环正在总验收中，静态校验已重新确认通过

## 执行阶段

### Phase 1：恢复前端工程健康

目标：先恢复 `admin` 的工程健康状态，让前端重新具备稳定的类型和构建闭环。

状态：

- [x] 修复 `admin` 当前 TypeScript 报错
- [x] 保证 `pnpm -C admin type-check` 通过
- [x] 保证页面层、子组件 props、composable 返回值契约一致（已修复本轮暴露问题）
- [x] 清理明显的 `string | undefined` / `status | undefined` 传递问题
- [x] 修复 `MenuView` 与 `useMenuPage` 的返回值漂移
- [x] 前端构建验证通过

验收标准：

- `pnpm -C admin type-check` 通过
- 不依赖模板隐式兜底
- 当前已存在页面至少可静态通过类型检查

### Phase 2：收紧前端模块边界

目标：在类型恢复后，收紧前端模块边界，降低页面继续膨胀和接入新页面时漏闭环的风险。

状态：

- [x] 菜单、部门、登录日志、操作日志、用户、角色、配置、公告、字典、文件页面的状态管理已进一步收紧；页面常量、默认查询、表单回填、payload、提示关闭、详情开关和展示工具已持续外提，成功提示已复用统一的轻量 composable，`system` 侧的 `config / notice / file / loginlog / operationlog / dict` 也已补齐这套页面边界
- [x] 拆分过大的 composable；`menu / department / user / role / config / notice / file / loginlog / operationlog` 的默认查询、回填和 payload 逻辑已从页面 composable 中继续拆出
- [x] 主要管理页已明确“页面负责装配，composable 负责状态，组件负责展示”的边界，`LoginLogView` 等页面已收成纯装配壳子
- [x] 分页、筛选、表单 model 的类型约束已在主要管理页收口；`menu / department / user / role / config / notice / dict / file / loginlog / operationlog` 已形成可复用的一致页面模式
- [x] 让动态菜单映射与页面注册关系更不容易漏改（已改为从页面文件自动发现组件，并回灌菜单表单选项）

验收标准：

- 核心页面 composable 不再同时承担过多职责
- 新增页面时，类型与装配路径清晰可复用
- 主要管理页具备一致的页面模式

### Phase 4：压实后端可维护性

目标：在后端已形成闭环的基础上，继续压实结构一致性和代码可维护性。

状态：

- [x] 后端主链路已形成闭环
- [-] 模块实现风格还未完全统一
- [x] 统一模块装配方式，已进一步收口 `post / notice / dict / auth-account / setup`，并把 `auth`、`loginlog`、`operationlog` 的 application service 切到 ports；`iam / system` 顶层路由已共用受保护 system group 组装器，`file / attachment / notice / dict / post / department / menu / role / user / loginlog / operationlog / setup` 已下沉统一的 service 构造入口
- [x] 新增模块时的最小后端接入模板已基本明确：`setup` 已补齐 `ports.go -> service.go -> repository.go -> routes.go` 结构，`modulekit` 已提供通用 protected group 壳子，`file / notice / dict / post / department / menu / role / user` 已提供模块级依赖装配入口供路由或关联模块复用
- [-] 保持 `datascope / authn / authz / bootstrap` 作为平台层稳定边界（`auth/login` 已改为依赖 `TokenIssuer` 接口，`attachment` 已去掉对 `file/application` 的直接类型别名依赖，并改为走 `file` 模块入口获取资产服务）

验收标准：

- 新增后端模块时不需要重新决定结构
- 服务层和仓储层模式趋于一致

### Phase 5：形成代码闭环完成条件

目标：统一做总验收，把“后端能跑、前端能交付、状态已同步”真正收成一套可信代码基线。当前阶段不维护单测/集成测试资产，后续验证策略转为 E2E + 人工链路回归。

状态：

- [x] 前端类型检查通过
- [x] 前端构建通过
- [-] 前后端关键链路待人工走通；本轮已重新确认 `pnpm -C admin type-check`、`pnpm -C admin build-only`、`go build ./...` 全部通过，后端启动链路仍以人工联调为最终口径
- [x] 根目录 `process.md` 状态全部同步到最新

验收标准：

- 代码本体可作为当前唯一可信基线
- 后续再回头整理文档时，以代码和 `process.md` 为准

## 实施顺序

1. [x] 先创建根目录 `process.md`
2. [x] 先填写“当前基线状态”和 5 个阶段的初始状态
3. [x] 优先推进 `Phase 1`
4. [-] `Phase 1` 完成后推进 `Phase 2`
5. [-] 并行推进 `Phase 4`
6. [-] 最后统一更新 `Phase 5` 总验收状态

## 假设与默认值

- 默认文件路径为仓库根目录：`D:\A\ez-admin-gin\process.md`
- 默认 `process.md` 只追踪代码，不追踪 `docs/`
- 默认当前“完成状态”以本轮验证结果为准：
  - 前端类型检查通过
  - 前端构建通过
  - 当前阶段不维护单测 / 集成测试资产
- 默认后续每完成一个阶段，都要同步改写 `process.md` 对应状态，不额外维护第二份 checklist
