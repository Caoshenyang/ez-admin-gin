# Execution Rules — 执行规则

## 每次开始任务

1. 阅读项目记忆文件（CLAUDE.md + docs/project/*）
2. 检查当前 Phase（PHASE_STATUS.md）
3. 检查当前任务是否属于当前 Phase
4. 输出执行计划：
   - 做什么
   - 预计修改哪些文件
   - 明确不做什么
5. 等计划清晰后再改代码

## 每次结束任务

1. 运行必要命令（make test-*）
2. 输出测试结果
3. 更新 PHASE_STATUS.md
4. 必要时更新 DECISION_LOG.md
5. 必要时更新 TESTING_STRATEGY.md
6. 总结剩余风险

## 防漂移规则

- 不得跳到后续 Phase
- 不得临时添加无关功能
- 不得扩大任务范围
- 不得重构无关模块
- 不得改变测试组织方式
- 不得把 TODO 伪装成完成
- 不得把 t.Skip 当成有效覆盖
- 不得在没有测试的情况下修改安全核心逻辑

## 任务分类

| 分类 | 处理方式 |
|------|----------|
| 当前 Phase 内任务 | 允许执行 |
| 当前 Phase 外任务 | 只记录到 backlog，不执行 |
| 高风险任务 | 必须先输出风险分析 |
| 架构变更 | 必须写 DECISION_LOG |
