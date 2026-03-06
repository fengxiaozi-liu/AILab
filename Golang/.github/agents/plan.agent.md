---
name: plan
description: 基于已就绪的 spec 生成 Kratos 微服务实施计划。
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

plan 是 specify/clarify 的下游，接收 Status 为 Ready 的 spec，产出可执行的技术实施计划。

本 Agent 同时加载两个 Skill：
- **speckit-plan**：计划生成方法论（Phase 0 → Phase 1 → 框架适配）
- **kratos-patterns**：Kratos 框架规范（命名、Proto、Ent、分层等），用于扩充 Phase 1 技术设计的框架考量

## 执行流程

> 领域方法论详见 **SKILL: speckit-plan**，框架规范详见 **SKILL: kratos-patterns**。

1. 前置检查 → *SKILL §前置检查*
2. 项目类型判定 → *SKILL §项目类型判定*
3. Phase 0 领域调研 → *SKILL §Phase 0*
4. Phase 1 技术设计 → *SKILL §Phase 1*
5. 框架适配 → *SKILL §框架适配*（加载 kratos-patterns，扩充技术设计）
6. 质量校验 → *SKILL §质量校验*
7. 汇报：
   - 输出项目类型、spec 路径、plan 路径
   - 列出关键设计决策摘要
8. 用户确认
   - 用户审核汇报内容，确认设计决策和计划概要
9. 交接
   - 用户确认后交接 **tasks** agent 拆分可执行任务

## 行为规则

- 调研中发现 spec 有遗漏 → 标记并建议补充，不擅自假设
- 设计决策需与用户确认后再写入
