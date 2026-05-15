## 未完成

- Phase 6 其余项目：
  - 业务模板
- 已移除：国际化、审批工作流、模块生成器 CLI（均无实际需求或已有替代方案）

## 当前下一步

1. Phase 6 剩余项目按优先级推进（业务模板）

## 阻塞点

- 无

## 最近一次执行记录

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
