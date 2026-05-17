## 未完成

- 当前阶段切换为质量与 UI 收口，详见 `docs/project/QUALITY_UI_STABILIZATION_PLAN.md`
- Phase 6 剩余功能暂停：
  - 业务模板
- 已移除：国际化、审批工作流、模块生成器 CLI（均无实际需求或已有替代方案）

## 当前下一步

1. 质量与 UI 收口计划 S0-S7 已完成，后续不再默认新增功能
2. 进入提交前整理：复查 diff、按计划建议拆分提交、确认文档与验证记录一致
3. 如继续优化，只做小范围回归和样式细节，不扩大业务功能面

## 阻塞点

- 无

## 最近一次执行记录

- **日期：** 2026-05-17
- **修改内容：**
  - **阶段切换：** 暂停 Phase 6 剩余功能补强，进入质量与 UI 收口阶段
  - **新增计划：** `QUALITY_UI_STABILIZATION_PLAN.md`
  - **S1 修复基础红线：**
    - 删除 `admin/src/modules/system/api/notification.ts` 中未使用的 `NotificationItem` 导入
    - 在根 `.gitignore` 中忽略 `test-results/`
    - 审阅当前后端认证 diff：`RequireUserID` / `Username` 仅集中处理当前用户上下文读取和 401 响应，未继续扩大行为面
  - **S2 设计 Token 收口（第一批）：**
    - `main.css` 新增 `--ez-radius-2xs/xs`、`--ez-brand-border`、`--ez-on-brand`、accent 和 sidebar scrollbar token
    - 全局 `.ez-card`、`.ez-toolbar`、`.ez-table-card`、`.ez-table-bar`、`.ez-modal-section` 改用 token
    - 工作标签、侧栏滚动条、登录页局部 Naive 变量、通知类型色、工作台指标色、文件图标色改用 `--ez-*`
    - 修复 `pnpm build` 暴露的严格类型问题：通知抽屉 slot / `NText depth`、WebSocket 回调入参、主题模式循环读取
  - **S2 设计 Token 收口（第二批）：**
    - 深色块文字统一使用 `--ez-on-dark*`
    - 登录页背景、字典面板英文标签、健康页深色状态卡、账户安全提示卡、操作日志详情抽屉改用语义变量
    - 操作日志详情抽屉移除固定 slate / rose 颜色，改用 `--ez-text-light`、`--ez-danger-*`、`--ez-panel-dark`
  - **S3 通用组件统一：**
    - `TableStatsBar` 增加具名操作区，配置、公告、文件、登录日志、操作日志、用户、岗位表格顶部统计栏统一复用
    - 表格底部分页统一改用 `.ez-table-footer`
    - `EmptyState` 接入占位页、附件空数据、字典项未选类型、工作台日志/登录/公告空区域
    - 新增 `PageHeader`，统一 IAM / System 高频页面标题、副标题和右侧操作区
  - **S4 页面级 UI 精修（第一批）：**
    - 页面头部 action 区支持窄屏换行，减少按钮挤压
    - `ez-table-bar` / `ez-table-footer` 增加响应式纵向排列
    - 附件表统一接入 `ez-table-card`、`TableStatsBar` 和 `.ez-table-footer`
    - 部门树表补齐 `TableStatsBar`，展示递归统计的部门节点数
    - 角色权限页在较窄屏幕切为单列，接口权限编辑行在移动宽度下纵向堆叠
    - 字典类型 / 字典项分页统一改用 `.ez-table-footer`
  - **S4 页面级 UI 精修（第二批）：**
    - 新增 `.ez-filter-actions`，统一筛选栏操作按钮组的间距和窄屏对齐
    - 操作日志筛选栏接入 `.ez-filter-actions`
    - 清理附件页、菜单表、部门表、操作日志表中的低信息注释
    - 表格顶部统计文案减少重复局部字号样式，交给 `TableStatsBar` 统一承接
  - **S4 页面级 UI 精修（第三批）：**
    - 用户表移除尚未接入行为的批量禁用、批量删除、列设置、密度按钮
    - 角色权限面板将误导性的“展开全部”文案改为“全选可用节点”
    - 占位页将“后续接入”从按钮改为状态标签，避免展示不可点击动作
    - `.ez-modal-footer` 统一按钮最小宽度、换行和小屏铺满行为，部门、附件编辑、附件上传弹窗接入该样式
  - **Chrome 视觉巡检记录：**
    - 已使用 Chrome 插件打开 `http://localhost:5174/login`
    - 登录页在 1817px 宽屏下无明显重叠、溢出或按钮挤压
    - 本轮仅访问本地 `localhost`，未读取 Cookie、密码、`localStorage` 等敏感浏览器存储
    - `127.0.0.1:8080` 后端未运行；Chrome 插件 evaluate 上下文不暴露 `localStorage`、`XMLHttpRequest` 或 `fetch`，无法临时 mock 受保护页接口
    - 用户、角色、菜单、字典、附件等受保护页需在 S7 前后端同时启动后复测
  - **S4 页面级 UI 精修（第四批）：**
    - 用户、岗位、部门、菜单、配置、公告、文件、附件、登录日志筛选栏的查询 / 重置按钮组统一接入 `.ez-filter-actions`
    - 用户、菜单、部门、附件、操作日志、登录日志表补充 `scroll-x`，减少窄屏下固定操作列或长文本列挤压主体列的风险
    - 附件表继续清理低信息注释，保留代码自解释结构
  - **S4 字典页视觉缺陷修复：**
    - 修复字典类型表格横向外溢到字典项面板下方的问题
    - 字典类型面板筛选区改为窄栏友好的垂直布局
    - 字典类型 / 字典项面板增加 `min-w-0`、`overflow-hidden` 和表格 `scroll-x`
    - 二次修正左侧窄栏可读性：压缩字典类型表列宽、取消固定操作列，分页改为简洁模式
    - 二次修正右侧筛选比例：字典项内容区改为独立三行网格，搜索输入限制最大宽度，避免被父级 grid 拉伸
    - 保留左侧表格形态，将双栏比例从 `420px / 1fr` 调整为 `520px / 1fr`，字典类型表 `scroll-x` 调整为 `500`
    - 字典类型筛选栏恢复一行展示，分页恢复 `show-size-picker`，与字典项分页风格保持一致
    - 修复双栏卡片高度过低导致表格内容被压缩、页面下方大面积空白的问题：双栏区域设置稳定最小高度，左右面板填满可用高度，表格区回到主体显示区域
  - **验证结果：**
    - `cd admin && pnpm exec oxlint .` 通过
    - `cd admin && pnpm exec vue-tsc --noEmit` 通过
    - `cd admin && pnpm exec vite build` 通过
    - `cd admin && pnpm build` 本轮曾在并发包装脚本处异常卡住；已清理该进程，拆分后的类型检查和 Vite 构建均通过
    - `cd server && go test ./...` 通过
    - `cd docs && pnpm docs:build` 通过（仅有既有 `conf` 语法高亮 fallback 提示）
  - **S5 前端代码质量收口：**
    - 巡检 `modules/*/components` 与 `modules/*/pages`，未发现展示组件或页面直接导入模块 API 的新增分层问题
    - 清理后台布局、IAM 高频 composable、System 高频 composable 和 System 类型文件中的低信息注释
    - 复查路由跳转、数组读取和 WebSocket 生命周期
    - `useWebSocket` 移除空主题 watcher 和无实际作用的 heartbeat interval，保留服务端 `ping` -> `pong` 响应和重连定时器
  - **S5 验证结果：**
    - `cd admin && pnpm exec oxlint .` 通过
    - `cd admin && pnpm exec vue-tsc --noEmit` 通过
    - `cd admin && pnpm build` 通过
  - **当前浏览器巡检状态：**
    - 本轮尝试使用干净 Playwright 浏览器巡检受保护页；沙箱放行后浏览器可启动，但本地 `5173/5174` 前端与 `8080` 后端均未监听，因此未继续截图
  - **S6 后端质量小收口：**
    - 审阅认证中间件和 auth handler 当前 diff：`RequireUserID` / `Username` 只收敛当前用户读取与 401 写响应
    - 将已有 `AccountLockChecker` 接入 auth `LoginService` 构造链路，让登录失败账号锁定配置实际生效
    - 删除 `server/tests/api/auth_api_test.go` 中登录限流 TODO，补充 `TestLoginRateLimiting` 集成测试
    - 更新 `TESTING_STRATEGY.md`：后端真实可跑测试数调整为 51，当前无已知测试 TODO
  - **S6 验证结果：**
    - `cd server && go vet ./...` 通过
    - `cd server && go test ./...` 通过
    - `cd server && go test -tags integration ./tests/api -run TestLoginRateLimiting -count=1` 通过
  - **S7 视觉验证和 E2E：**
    - 复用本机已启动的 `8080` 后端，`/health` 返回 database / redis / status 均为 `ok`
    - 修正视觉巡检脚本的页面路径：使用真实菜单路径 `/system/dicts`、`/system/users`、`/system/roles`、`/system/menus`、`/system/login-logs`、`/system/operation-logs`、`/system/attachments`
    - Playwright 截图巡检 light / dark 模式下的通知抽屉、工作台、系统状态、数据字典、登录日志、操作日志、用户管理、角色管理、菜单管理、附件中心
    - 截图产物位于 `/private/tmp/ez-admin-visual-<mode>-<page>.png`
    - 字典管理页 light / dark 下左右表格比例正常，左侧内容完整，搜索栏宽度稳定，分页风格一致
    - 登录日志 User-Agent ellipsis 文本在 DOM 指标中被识别为越界，但截图显示已被单元格裁剪和省略，不构成视觉问题
  - **S7 验证结果：**
    - `cd admin && pnpm exec playwright test` 通过，43 passed
    - 视觉巡检脚本无 console error / pageerror
  - **下一步：** 收口计划 S0-S7 已完成；进入提交前整理和小范围回归

- **日期：** 2026-05-15
- **修改内容：**
  - **WebSocket 实时通知系统：**
    - 后端：新建 `system/notification/` DDD 模块（domain/application/infra/api/ws）
    - WebSocket 库选用 `coder/websocket`（nhooyr.io/websocket）
    - 消息分发：本地 Hub + Redis Pub/Sub fan-out
    - 持久化：`sys_notification` 表（migration 000010）
    - REST API：列表、未读数、标记已读、全部已读
    - WS 端点：`/api/v1/system/notifications/ws?token=xxx`（query param 认证）
    - 前端：TS 类型、API、Pinia store、useWebSocket composable、NotificationDrawer
    - AppHeader 铃铛 + NBadge 未读数 + 打开通知抽屉
    - AdminLayout 挂载时自动连接 WS，卸载时断开
  - **设计决策（ADR-016）：** WebSocket 库选择、WS 认证方式、模块独立于 notice
