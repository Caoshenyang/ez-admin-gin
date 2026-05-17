# Quality & UI Stabilization Plan — 质量与界面收口计划

> 当前策略：停止新增功能，优先把已有能力打磨到稳定、统一、可验证。
> 本计划接替 Phase 6 的功能补强线，作为当前阶段的执行记忆。

## 为什么切换到收口阶段

项目已经完成测试、安全、可观测性、E2E、发布治理、暗色模式和 WebSocket 通知等关键能力。继续新增功能会扩大维护面，而当前更有价值的是：

- 让代码检查和测试长期保持全绿。
- 把前端样式从“能用”收敛到“统一、耐用、可复用”。
- 清理半成品 UI token 化、重复组件样式和文档状态漂移。
- 保持后台系统的工作台气质：紧凑、清晰、可扫描。

## 当前基线

检查日期：2026-05-16

| 项目 | 当前结果 | 说明 |
|------|----------|------|
| `pnpm exec vue-tsc --noEmit` | 通过 | 前端类型基线干净 |
| `go test ./...` | 通过 | 后端基础测试通过 |
| `pnpm exec oxlint .` | 通过 | S1 已清理未使用导入 |
| E2E | 未运行 | 需要前端、后端和测试环境 |
| 文档状态 | 有漂移 | `TESTING_STRATEGY.md` 声称 TODO 为 0，但实际仍有登录限流 TODO |

当前工作区改动主要集中在：

- `admin/src/styles/main.css`：新增 radius / font-size token。
- `admin/src/components/app-shell/*`：后台壳子样式变量化。
- `admin/src/modules/auth/*`：登录页、工作台视觉变量化。
- `admin/src/modules/system/*`、`admin/src/modules/iam/*`：部分表格、弹窗、面板样式变量化。
- `server/internal/modules/auth/api/handlers.go`、`server/internal/platform/middleware/auth.go`：少量认证相关改动，需要单独确认行为。

## 总体执行规则

- 不新增业务功能。
- 不推进 Phase 6 的“业务模板”。
- 不扩大模块边界，不顺手重构无关代码。
- UI 优化优先复用 Naive UI 和现有 `--ez-*` 变量。
- 每个阶段结束都更新 `PHASE_STATUS.md`。
- 涉及架构或长期方向变更时更新 `DECISION_LOG.md`。

## S0：冻结功能线

**目标：** 把项目从“继续补功能”切到“质量和 UI 收口”。

**任务：**

- 将 Phase 6 剩余功能标记为暂停。
- 明确当前阶段只处理质量、样式、验证和文档漂移。
- 清点当前未提交改动并分组：UI token、组件样式、后端认证、临时文件、测试产物。
- 清理或忽略 `test-results/` 等本地测试产物。

**验收标准：**

- `QUALITY_UI_STABILIZATION_PLAN.md` 存在并作为当前阶段记忆。
- `PHASE_STATUS.md` 的当前下一步指向本计划。
- 后续任务不得再默认接入新功能。

**状态：** 已完成（2026-05-16）

## S1：修复基础红线

**目标：** 先让项目回到检查全绿的最低稳定线。

**任务：**

- 删除 `admin/src/modules/system/api/notification.ts` 中未使用的 `NotificationItem` 导入。
- 确认 `admin/src/api/contract-check.ts` 是否需要保留，若保留则补充用途说明或接入检查链路。
- 单独审阅后端认证相关改动，避免 UI 收口任务中混入未验证的安全语义变化。
- 检查 `test-results/` 是否需要加入忽略规则。

**验收标准：**

- `cd admin && pnpm exec oxlint .` 通过。
- `cd admin && pnpm exec vue-tsc --noEmit` 通过。
- `cd server && go test ./...` 通过。

**状态：** 已完成（2026-05-16）

## S2：设计 Token 收口

**目标：** 完成当前已经做了一半的字号、圆角、颜色变量化。

**任务：**

- 统一 `--ez-radius-*` 的使用，减少页面散落的 `rounded-lg`、`rounded-2xl`、`rounded-[16px]`。
- 统一 `--ez-text-*` 的使用，避免重复硬编码 `13px`、`14px`、`28px`。
- 收敛业务组件中的品牌色、状态色、浅背景色，不在组件里直接散落色值。
- 梳理 `admin/src/styles/main.css` 与 `admin/src/styles/theme.ts` 的职责：
  - `main.css` 负责自定义组件、Tailwind 补充和 CSS 变量。
  - `theme.ts` 负责 Naive UI 主题桥接。
- 把 `.ez-card`、`.ez-toolbar`、`.ez-table-card`、`.ez-table-bar` 里的固定圆角、字号、阴影迁移到 token。

**验收标准：**

- 新增 UI 不再直接写散落圆角、字号和品牌色。
- 明暗主题下卡片、表格、弹窗、抽屉没有明显浅色残留。
- `rg "rounded-\\[|text-\\[|#[0-9A-Fa-f]{3,8}|rgba\\(" admin/src` 的剩余结果都能解释为 token 定义、Naive UI theme 或必要特例。

**状态：** 已完成（2026-05-16，浏览器视觉巡检留到 S7）

## S3：通用组件统一

**目标：** 减少页面重复拼 class，让后台页面更像一套产品。

**任务：**

- 完善 `EmptyState`，统一表格、占位页和搜索无结果状态。
- 完善 `TableStatsBar`，替换手写 `.ez-table-bar` 的重复结构。
- 抽取或规范页面头部：标题、副标题、右侧操作区。
- 继续使用 `FormModalHeader`，但将渐变、关闭按钮、间距纳入 token。
- 统一表格容器写法，减少 `NCard + content-class` 的分叉。

**验收标准：**

- IAM 与 System 模块的表格顶部、底部分页、空状态视觉一致。
- 新模块可直接复用通用组件，不需要复制旧页面结构。
- 组件不直接调用 API，继续保持 `api / types / composables / components / pages` 分层。

**状态：** 已完成（2026-05-16，首轮收口；后续页面细节进入 S4）

**已完成：**

- `TableStatsBar` 增加具名操作区，替换配置、公告、文件、登录日志、操作日志、用户、岗位表格顶部统计栏。
- 表格底部分页统一改用 `.ez-table-footer`。
- `EmptyState` 接入占位页、附件空数据、字典项未选类型、工作台日志/登录/公告空区域。
- 新增 `PageHeader`，统一 IAM / System 高频页面标题、副标题和右侧操作区。

**保留到 S4 的细节：**

- 继续统一表格容器写法，减少 `NCard + content-class` 的分叉。
- 复查小尺寸区域是否仍适合保留 Naive UI 原生 `NEmpty`。

## S4：页面级 UI 精修

**目标：** 按后台系统气质统一高频页面体验。

**优先级：**

1. 后台壳子：`AdminLayout`、`AppSidebar`、`AppHeader`、`WorkTabs`。
2. 高频业务页：用户、角色、菜单、部门、字典、附件。
3. 系统页：健康检查、登录日志、操作日志、占位页。
4. 登录页和工作台：保留质感，同时控制视觉重量。

**检查项：**

- 一屏布局稳定，业务区内部滚动正常。
- 表格列宽、操作按钮、筛选栏紧凑可扫读。
- hover、focus、active 状态一致。
- 窄屏下不出现明显文字溢出或按钮挤压。
- 暗色模式下文本、边框、面板层级可读。

**验收标准：**

- 核心页面在 light / dark 下都可用。
- 管理台没有营销页式大 hero、大卡片堆叠。
- 页面内部没有明显重复的视觉模式实现。

**状态：** 进行中（2026-05-17 已完成第四批高频页精修，并补充 Chrome 登录页视觉巡检记录）

**已完成：**

- 页面头部 action 区支持窄屏换行，避免按钮挤压。
- `ez-table-bar` / `ez-table-footer` 增加响应式纵向排列，表格统计和分页在窄屏下更稳定。
- 附件表统一接入 `ez-table-card`、`TableStatsBar` 和 `.ez-table-footer`。
- 部门树表补齐 `TableStatsBar`，展示递归统计的部门节点数。
- 角色权限页在较窄屏幕切为单列，接口权限编辑行在移动宽度下纵向堆叠。
- 字典类型 / 字典项分页统一改用 `.ez-table-footer`。
- 新增 `.ez-filter-actions`，统一筛选栏操作按钮组的间距和窄屏对齐。
- 操作日志筛选栏接入 `.ez-filter-actions`。
- 清理附件页、菜单表、部门表、操作日志表中的低信息注释。
- 表格顶部统计文案减少重复局部字号样式，交给 `TableStatsBar` 统一承接。
- 用户表移除尚未接入行为的批量禁用、批量删除、列设置、密度按钮。
- 角色权限面板将误导性的“展开全部”文案改为“全选可用节点”。
- 占位页将“后续接入”从按钮改为状态标签，避免展示不可点击动作。
- `.ez-modal-footer` 统一按钮最小宽度、换行和小屏铺满行为，部门、附件编辑、附件上传弹窗接入该样式。
- 使用 Chrome 插件打开 `http://localhost:5174/login` 做真实浏览器宽屏视觉巡检：
  - 登录页整体布局、表单、说明卡片和页脚在 1817px 宽屏下无明显重叠、溢出或按钮挤压。
  - 本轮仅访问本地 `localhost`，未读取 Cookie、密码、`localStorage` 等敏感浏览器存储。
  - `127.0.0.1:8080` 后端未运行；Chrome 插件当前的 evaluate 上下文不暴露 `localStorage`、`XMLHttpRequest` 或 `fetch`，无法在真实 Chrome 中做临时 API mock，因此用户、角色、菜单、字典、附件等受保护页留到 S7 在“前后端同时启动”条件下复测。
- 高频筛选栏的查询 / 重置按钮组统一接入 `.ez-filter-actions`：
  - 用户、岗位、部门、菜单、配置、公告、文件、附件、登录日志筛选栏已统一。
  - 保留操作日志筛选栏既有网格布局和 `.ez-filter-actions`。
- 用户、菜单、部门、附件、操作日志、登录日志表补充 `scroll-x`，避免窄屏下固定操作列或长文本列挤压主体列。
- 附件表继续清理低信息注释，保留代码自解释结构。
- 修复字典管理页双栏布局溢出：
  - 字典类型面板改为窄栏友好的垂直筛选布局。
  - 字典类型 / 字典项面板增加 `min-w-0` 和 `overflow-hidden`，限制内部表格外溢。
  - 字典类型 / 字典项表格补充 `scroll-x`，固定操作列只在自身卡片内横向滚动。
  - 二次修正左侧窄栏可读性：字典类型表压缩列宽、取消固定操作列，分页改为简洁模式。
  - 二次修正右侧筛选比例：字典项内容区改为独立三行网格，搜索输入限制最大宽度，避免被父级 grid 拉伸。
  - 保留左侧表格形态，将双栏比例从 `420px / 1fr` 调整为 `520px / 1fr`，并将字典类型表 `scroll-x` 调整为 `500`，让左侧表头、数据和操作列完整呈现。
  - 字典类型筛选栏在 520px 左栏下恢复一行展示：关键字输入、状态选择和查询 / 重置按钮同排；分页恢复 `show-size-picker`，与右侧分页风格一致。
  - 双栏内容区增加稳定最小高度，字典类型 / 字典项卡片填满可用高度，避免表格内容被挤在卡片上半段、页面下方留下大面积空白。

**待继续：**

- 在后端可用后，浏览器巡检用户、角色、菜单、字典、附件等页面的 light / dark 实际视觉。
- 继续收口表格行高、操作列按钮密度和暗色层级。
- 复查工作台、健康检查、日志页的暗色层级和滚动区域。

## S5：前端代码质量收口

**目标：** UI 调整后保持代码边界清晰、类型稳定。

**任务：**

- 检查 page 是否只做编排，业务副作用留在 composable。
- 检查 component 是否只通过 props / emits 通信，不直接导入 API。
- 清理无意义注释，例如“xxx 函数。”这类低信息注释。
- 统一 `void router.push(...)` / `void router.replace(...)`。
- 检查数组读取和异步状态是否有明确 undefined 处理。
- 复查通知 WebSocket store 与 `AdminLayout` 生命周期，避免连接泄漏或状态错乱。

**验收标准：**

- `cd admin && pnpm exec oxlint .` 通过。
- `cd admin && pnpm exec vue-tsc --noEmit` 通过。
- `cd admin && pnpm build` 通过。

**状态：** 已完成（2026-05-17，前端代码质量首轮收口）

**已完成：**

- 巡检 `modules/*/components` 与 `modules/*/pages`，未发现展示组件或页面直接导入模块 API 的新增分层问题。
- 清理后台布局、IAM 高频 composable、System 高频 composable 和 System 类型文件中的低信息注释，保留解释状态同步、重置基线和权限边界的注释。
- 复查路由跳转调用：普通导航继续使用 `void router.push(...)` / `void router.replace(...)`，登录提交流程保留 `await router.push('/dashboard')` 以保持提交时序。
- 复查数组读取：当前命中的索引读取均已有 `?.` 或兜底值。
- 复查通知 WebSocket 生命周期：`AdminLayout` 挂载时连接、卸载时断开；`notification` store 防重复连接并在 reset 时断开；`useWebSocket` 移除空主题 watcher 和无实际作用的 heartbeat interval，仅保留服务端 `ping` -> `pong` 响应和重连定时器。

**验证结果：**

- `cd admin && pnpm exec oxlint .` 通过。
- `cd admin && pnpm exec vue-tsc --noEmit` 通过。
- `cd admin && pnpm build` 通过。

## S6：后端质量小收口

**目标：** 不加业务功能，只处理当前可见的后端质量风险。

**任务：**

- 单独确认认证中间件和 auth handler 的当前未提交改动。
- 处理登录限流 TODO：
  - 如果功能已存在，补测试。
  - 如果功能暂不做，移动到明确 backlog，不把 TODO 伪装成完成。
- 更新 `TESTING_STRATEGY.md` 中关于 TODO / t.Skip 的真实状态。
- 必要时补充契约或认证相关测试。

**验收标准：**

- `cd server && go vet ./...` 通过。
- `cd server && go test ./...` 通过。
- 如本地 DB / Redis 可用，`make test-integration` 通过或记录明确跳过原因。
- 测试策略文档与真实代码状态一致。

**状态：** 已完成（2026-05-17，登录限流 TODO 收口）

**已完成：**

- 审阅认证中间件和 auth handler 当前 diff：`RequireUserID` / `Username` 只收敛当前用户读取与 401 写响应，没有扩大认证语义。
- 将已有 `AccountLockChecker` 接入 auth `LoginService` 构造链路，让 `rate_limit.login_lockout_threshold` / `login_lockout_sec` 配置实际生效。
- 删除 `server/tests/api/auth_api_test.go` 中登录限流 TODO，补充 `TestLoginRateLimiting` 集成测试，验证达到失败阈值后登录接口返回 429。
- 更新 `TESTING_STRATEGY.md`：后端真实可跑测试数从 50 调整为 51，TODO 状态改为当前无已知测试 TODO，仅保留 DB / Redis / Casbin 环境不可用时的条件跳过说明。

**验证结果：**

- `cd server && go vet ./...` 通过。
- `cd server && go test ./...` 通过。
- `cd server && go test -tags integration ./tests/api -run TestLoginRateLimiting -count=1` 通过。

## S7：视觉验证和 E2E

**目标：** 用浏览器和自动化测试验证 UI 收口没有破坏核心流程。

**任务：**

- 启动后端和前端。
- 跑登录、菜单、权限、用户、角色相关 Playwright 用例。
- 手工或截图巡检核心页面：
  - 登录页
  - 工作台
  - 用户管理
  - 角色管理
  - 菜单管理
  - 字典管理
  - 附件管理
  - 健康检查
  - 通知抽屉
- 分别检查 light / dark 模式。

**验收标准：**

- `cd admin && pnpm exec playwright test` 通过，或记录明确失败原因。
- 核心页面无明显文本溢出、按钮挤压、滚动异常、暗色模式残留。
- 截图检查结果记录到 `PHASE_STATUS.md`。

**状态：** 已完成（2026-05-17，E2E 与 light / dark 视觉巡检通过）

**已完成：**

- 复用本机已启动的后端服务，健康检查返回 database / redis / status 均为 `ok`。
- 运行完整 Playwright E2E：登录、菜单、权限、用户、角色相关 43 个用例全部通过。
- 使用 Playwright 截图巡检核心页面的 light / dark 模式：
  - 通知抽屉、工作台、系统状态、数据字典、登录日志、操作日志、用户管理、角色管理、菜单管理、附件中心。
  - 截图输出到 `/private/tmp/ez-admin-visual-<mode>-<page>.png`。
- 修正视觉巡检脚本的路径假设：字典、用户、角色、菜单、日志、附件等页面使用后端菜单种子中的真实复数路径，例如 `/system/dicts`、`/system/users`。

**验证结果：**

- `cd admin && pnpm exec playwright test` 通过，43 passed。
- 视觉巡检脚本无 console error / pageerror。
- 字典管理页 light / dark 下左右表格比例正常，左侧字典类型表内容完整，右侧筛选栏宽度稳定，左右分页风格一致。
- 表格巡检未发现实际横向滚动溢出；登录日志中 User-Agent ellipsis 文本会被 DOM 指标识别为越界，但截图显示已被单元格裁剪和省略，不构成视觉问题。

## 建议提交拆分

1. `quality baseline`：修 lint、临时文件、文档状态漂移。
2. `ui tokens`：收敛 `main.css`、`theme.ts` 和散落样式变量。
3. `shared admin components`：统一空状态、表格工具条、页面头部、弹窗头部。
4. `page polish`：按模块做页面级 UI 精修。
5. `verification`：补测试、跑 E2E、更新阶段状态。

## 当前下一步

当前收口计划 S0-S7 已完成。后续不再默认新增功能；如继续推进，优先做小范围回归清单、提交拆分和发布前文档校准。
