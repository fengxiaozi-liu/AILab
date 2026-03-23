---
name: kratos-components
description: 用于 Kratos 组件与基础设施能力的使用规范，包括 Ent、EventBus、Crontab、Depend、Config。
  适用于新增或修改 schema、索引、edge、事件、监听器、定时任务、跨服务依赖封装、配置项、默认值或校验逻辑的场景。
  触发关键词：ent schema、eventbus、listener、cron、crontab、depend、InnerRPC、config、幂等、重试、批量装配。
  DO NOT USE FOR：业务逻辑编排（→ kratos-domain）、接入层协议适配（→ kratos-entry）、横切规范（→ kratos-conventions）。
---

# Kratos Components

## 输入

- 必需：变更目标描述（组件类别 / 变更内容）
- 可选：相关 Ent schema 路径、EventBus 事件定义、config 字段
- 可选：`specs/<feature>/tasks.md`

缺少必需输入时，MUST 先向用户提问，不得猜测继续。

## 工作流

1. 识别本次变更类型：`ent` / `eventbus` / `crontab` / `depend` / `config`
2. 输出开始前结构化状态（见强制输出）
3. 按需加载对应参考文件（见参考文件清单）
4. 执行变更
5. IF 涉及 Ent schema → 评估 Repo relation 装配和 codegen 影响，联动 `kratos-entry` 执行 `ent generate`
   IF 涉及 Depend / EventBus → 检查批量化策略和超时设置
   IF 涉及 Crontab → 检查幂等、重试、并发控制设计
6. 输出完成后结构化状态（见强制输出）

## 约束

### MUST
- 组件能力 MUST 按职责收口，不把业务编排塞进基础设施层
- Depend / EventBus MUST 优先批量收集和统一封装，避免散落直连与逐条远程调用
- Crontab MUST 显式考虑幂等、重试、并发控制和可观测性
- InnerRPC 调用 MUST 设置超时；重试仅用于幂等操作并采用退避策略
- Ent schema 变更 MUST 同步评估 Repo relation 装配和 codegen 影响
- Config 新增配置项 MUST 提供默认值或必填校验，并区分环境差异
- Config MUST 可校验、可回滚、可观测
- EventBus Listener MUST 保持幂等和可重试
- EventBus 传递 MUST 优先复用已有稳定对象，不额外新建近义传递结构
- MySQL 冲突写 MUST 使用 `OnConflict()`，不使用 `OnConflictColumns(...)`

### MUST NOT
- MUST NOT 手改 Ent 生成产物
- MUST NOT 在 Listener 中隐藏核心业务编排
- MUST NOT 仅为 EventBus 传递场景新建 `EventPayload`、`ContextDTO`、`XxxData` 等近义壳结构
- MUST NOT 在循环中逐条远程请求（N+1 depend 调用）
- MUST NOT 散落直连上游（不经 Depend 封装）
- MUST NOT 无幂等保证的 Crontab 周期性写操作直接上线
- MUST NOT 新增配置项但没有校验或默认值说明

### SHOULD
- Ent schema 变化后 SHOULD 联动评估 Repo relation 和生成物影响
- 批量操作 SHOULD 保持批量收集 ID、批量查询、批量回填模式

## 强制输出

开始前输出：

```json
{
  "componentScope": "ent | eventbus | crontab | depend | config",
  "infraChange": "变更摘要",
  "riskControl": "幂等 / 重试 / 批量化 / 校验等控制点"
}
```

完成后输出：

```json
{
  "noBatchN1": true,
  "failureHandled": true,
  "configValidated": true,
  "codegenTriggered": true
}
```

## 参考文件

| 文件 | 适用场景 | 加载时机 |
|------|----------|----------|
| `reference/ent-spec.md` | 新增或修改 Ent schema、edge、index、annotation | 涉及 Ent 变更时 |
| `reference/eventbus-spec.md` | 新增或修改 EventBus 事件、发布与 Listener | 涉及 EventBus 变更时 |
| `reference/crontab-spec.md` | 新增或修改定时任务、补偿任务、调度策略 | 涉及 Crontab 变更时 |
| `reference/depend-spec.md` | 调整 InnerRPC 或跨服务依赖封装 | 涉及 Depend 变更时 |
| `reference/config-spec.md` | 新增或变更配置项、默认值、校验、环境差异 | 涉及 Config 变更时 |
