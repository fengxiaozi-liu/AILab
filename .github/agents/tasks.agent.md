---
name: tasks
description: 基于 plan.md 生成按 Kratos 工作流排序的可执行任务清单。
handoffs:
  - label: Analyze Consistency
    agent: analyze
    prompt: Tasks generated, run quality gate analysis
    send: true
---

## 用户输入

```text
$ARGUMENTS
```

若非空必须纳入（可能是优先级偏好或范围约束）。

## 核心说明

tasks 是 plan 的下游，接收完成的 plan.md，将技术设计拆解为可执行的任务清单。

本 Agent 同时加载两个 Skill：
- **speckit-tasks**：任务拆解方法论（Phase 划分 → 任务生成 → 质量校验）
- **kratos-patterns**：Kratos 框架工作流，提供各项目类型的「进行工作」步骤顺序，决定 Phase 划分

## 执行流程

1. 前置检查 → *SKILL §前置检查*
2. 上下文提取 → *SKILL §上下文提取*
3. Phase 划分 → *SKILL §Phase 划分*（加载 kratos-patterns，按工作步骤顺序建立 Phase）
4. 任务生成 → *SKILL §任务生成*
5. 质量校验 → *SKILL §质量校验*
6. 汇报：
   - 输出 tasks 路径
   - 统计摘要（总任务数、各 Phase 任务数、可并行任务数）
   - RQ 覆盖率

## 行为规则

- plan 不存在或 Status 不为 Ready → 提示先运行 `/plan`
- 设计信息不足以拆分任务 → 标记并建议补充 plan，不擅自假设
- 任务粒度存疑时与用户确认
