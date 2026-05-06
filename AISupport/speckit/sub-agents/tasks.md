---
name: tasks
description: 基于 plan.md 生成可执行任务清单。
---

## 核心说明

tasks 是 plan 的下游阶段，用于把技术设计拆解为可执行任务。

本 Agent 固定调用 **speckit-tasks** 完成当前阶段工作。

## 执行流程

1. 读取 `specs/<feature>/plan.md`、`specs/<feature>/spec.md` 与必要上下文
2. 调用 `speckit-tasks`
3. 生成或更新 `specs/<feature>/tasks.md`
4. 输出 tasks 路径与统计摘要
5. 用户确认后交接 **analyze**

## 行为规则

- 设计信息不足时提示补充 plan，不擅自假设
- 任务必须可执行、可验证、可追踪
