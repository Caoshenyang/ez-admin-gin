## 未完成

- Phase 6 其余项目：
  - 模块生成器（进行中，约 50%）
  - WebSocket 通知
  - 业务模板
- 已移除：国际化、审批工作流（无实际需求）

## 当前下一步

1. 完成模块生成器剩余工作：
   - 后端模板：app_service.go.tpl, infra_repo.go.tpl, routes.go.tpl, services.go.tpl（4/10 待完成）
   - 前端模板：全部 8 个（types_ts, api_ts, composable_ts, utils_ts, form_modal_vue, filter_bar_vue, table_vue, view_vue）
   - 示例配置 gen/examples/product.yaml
   - Makefile `make gen` 目标
   - 端到端测试（go build + pnpm type-check）
2. Phase 6 剩余项目按优先级推进

## 阻塞点

- 无

## 最近一次执行记录

- **日期：** 2026-05-14
- **修改内容：**
  - **模块生成器 WIP（约 50%）：**
    - 创建 CLI 入口 server/cmd/gen/main.go（YAML 配置解析、项目根目录查找）
    - 创建配置类型 server/cmd/gen/config.go（Config/Field 结构体、验证、辅助方法）
    - 创建模板引擎 server/cmd/gen/generator.go（embed FS、模板渲染、文件输出、接入指引打印）
    - 创建后端模板 6/10：
      - model.go.tpl（GORM 模型，含状态类型生成）
      - domain_types.go.tpl（DTO、权限常量、归一化、BuildResponse）
      - api_dto.go.tpl（API 层类型别名）
      - api_routes.go.tpl（路由注册）
      - api_handler.go.tpl（HTTP 处理器，含 Swagger 注释）
      - app_ports.go.tpl（Repository 接口）
    - 后端模板待完成：app_service.go.tpl, infra_repo.go.tpl, routes.go.tpl, services.go.tpl
    - 前端模板全部待完成：8 个 .tpl 文件
  - **设计决策：**
    - 使用 Go embed.FS 嵌入模板文件，单二进制分发
    - 已存在文件不覆盖（安全策略）
    - 生成后打印手动接入指引（路由注册、菜单图标、种子数据）
    - 不自动修改现有文件（避免破坏性操作）
  - **QUALITY_ROADMAP 更新：** 移除国际化、审批工作流
