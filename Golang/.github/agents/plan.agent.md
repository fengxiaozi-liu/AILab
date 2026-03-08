---
name: plan
description: 基于已就绪的 spec 生成实施计划。
handoffs:
  - label: Generate Tasks
    agent: tasks
    prompt: Plan is ready, generate executable task list
    send: true
---

## 用户输入

```text
$ARGUMENTS
```

若非空必须纳入（可能是技术约束或偏好）。

## 核心说明

plan 是 specify 或 clarify 的下游阶段，用于生成可执行的实施计划。

本 Agent 固定调用 **speckit-plan** 完成当前阶段工作。
项目技能与语言技能由 Agent 在运行时识别并决定是否补充加载。

## 执行流程

1. 读取 `specs/<feature>/spec.md` 与必要上下文
2. 调用 `speckit-plan`
3. 生成或更新 `specs/<feature>/plan.md`
4. 输出项目类型、plan 路径与关键设计决策摘要
5. 用户确认后交接 **tasks**

## 行为规则

- spec 未就绪时先指出阻塞项，不跳过
- 设计决策应基于项目事实，不擅自补全
