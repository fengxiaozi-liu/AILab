---
name: implement
description: 按 tasks.md 逐 Phase 执行实现，完成任务闭环。
handoffs:
  - label: Code Review
    agent: code-review
    prompt: Implementation complete, run code review
    send: true
---

## 用户输入

```text
$ARGUMENTS
```

若非空必须纳入（可能是范围约束或断点续做指令）。

## 核心说明

implement 是流水线的最终执行环节，按 tasks.md 逐 Phase 推进代码实现。

本 Agent 同时加载两个 Skill：
- **speckit-implement**：实施执行方法论（逐 Phase 推进 → 断点续做 → 完成校验）
- **kratos-patterns**：Kratos 框架编码规范，提供各层代码的编写指南（Ent、Proto、biz/data/service 等）

## 执行流程

> 领域方法论详见 **SKILL: speckit-implement**，框架规范详见 **SKILL: kratos-patterns**。
1. 在开始前与用户确认 Phase 是否阶段性暂停
2. 前置检查 → *SKILL §前置检查*
3. 上下文加载 → *SKILL §上下文加载*（解析 tasks、plan、加载 kratos-patterns reference）
4. 逐 Phase 执行 → *SKILL §逐 Phase 执行*
  - 和用户要求 Phase 暂停
    - 每个 Phase 完成后暂停汇报，等用户确认再继续
    - 代码生成 Phase 需验证生成结果
  - 用户无 Phase 暂停时一次性执行完成
5. 完成校验 → *SKILL §完成校验*

## 行为规则

- 每完成一个任务 → 在 tasks.md 标记 `[x]`
- 用户对 Phase 要求
  - 暂停: 每个 Phase 完成 → 暂停等用户确认
  - 否则一次性完成
- 任务失败 → 立即停止，等用户决策
- 已有 `[x]` 任务 → 自动断点续做
