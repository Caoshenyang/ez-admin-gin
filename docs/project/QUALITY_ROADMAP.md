# Quality Roadmap — 长期路线图

> 从约 8.0 分提升到 9.0+ 的分阶段计划。
> 当前阶段：质量与 UI 收口进行中（Phase 6 剩余功能暂停）
> 执行计划：`docs/project/QUALITY_UI_STABILIZATION_PLAN.md`

---

## Phase 1：集中式测试做实 (当前)

**目标：** 把 server/tests 从测试骨架变成真实可重复运行的质量门禁。

**范围：**
- API 黑盒测试（auth、user、role、menu）
- RBAC 流程测试（权限授予/拒绝、super_admin bypass）
- datascope 流程测试
- OpenAPI 契约测试增强
- 测试数据隔离（seed/reset）

**验收标准：**
- `make test-contract` 通过
- `make test-api` 通过
- `make test-rbac` 通过（非 t.Skip）
- `make test-integration` 通过
- CI integration job 通过
- RBAC / datascope 测试不再大面积 t.Skip
- 测试数据可重复初始化

**状态：** 已完成

---

## Phase 2：安全基线升级

**目标：** 提升认证和授权安全性。

**范围：**
- Access Token + Refresh Token 双 token
- HttpOnly Secure Cookie
- Refresh Token rotation
- 服务端会话撤销
- 登录限流增强
- 生产配置强校验
- CORS 生产校验
- 上传安全增强
- 安全响应头

**状态：** 已完成

---

## Phase 3：可观测性和运维能力

**目标：** 提升生产环境可观测性和运维能力。

**范围：**
- X-Request-ID
- 结构化日志
- /healthz + /readyz
- Prometheus metrics
- 错误码规范
- 运维排障文档
- 备份恢复文档

**状态：** 已完成

---

## Phase 4：前端质量和 E2E

**目标：** 建立前端 E2E 测试体系。

**范围：**
- Playwright E2E
- 登录流程
- 菜单权限
- 按钮权限
- 用户管理
- 角色授权
- 无权限页面
- token 过期处理
- OpenAPI 生成前端类型
- CI 检查前后端契约一致

**状态：** 已完成

---

## Phase 5：发布治理和文档成熟

**目标：** 建立发布流程和文档规范。

**范围：**
- CHANGELOG
- Release note
- Migration guide
- Issue templates / PR template
- SECURITY.md / CONTRIBUTING.md
- Production checklist
- Testing strategy 文档
- Architecture decision records
- Demo 环境

**状态：** 已完成

---

## Phase 6：产品能力补强

**目标：** 补强产品级能力。

**范围：**
- 暗色模式 / 主题系统 ✅
- 国际化（已移除：内部工具无多语言需求）
- 模块生成器（已替换为 AI Skill：`.agents/skills/module-generator/SKILL.md`）
- WebSocket 通知 ✅
- 业务模板
- 审批工作流（已移除：当前无实际需求）

**状态：** 暂停继续新增功能（暗色模式 ✅、模块生成器已替换为 Skill、WebSocket 通知 ✅）

---

## Stabilization：质量与 UI 收口（当前）

**目标：** 不再新增功能，优先完善代码质量、统一 UI 样式、补齐验证闭环。

**范围：**
- 修复当前基础红线
- 收敛前端设计 token
- 统一通用后台组件
- 精修核心页面体验
- 收口前端代码质量
- 小范围确认后端认证质量风险
- 执行视觉验证和 E2E

**验收标准：**
- `cd admin && pnpm exec oxlint .` 通过
- `cd admin && pnpm exec vue-tsc --noEmit` 通过
- `cd admin && pnpm build` 通过
- `cd server && go test ./...` 通过
- 核心后台页面 light / dark 模式无明显视觉问题
- `TESTING_STRATEGY.md` 与真实测试状态一致

**状态：** 进行中，详见 `QUALITY_UI_STABILIZATION_PLAN.md`
