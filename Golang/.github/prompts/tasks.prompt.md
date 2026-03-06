---
description: 基于 plan.md 生成按 Kratos 工作流排序的可执行任务清单。
agent: tasks
---

# /tasks 命令

调用 **tasks** agent 将技术设计拆解为可执行任务。

## 何时使用

- plan.md 已完成（Status 为 Ready），需要拆分为具体可执行任务
- 需要明确实施顺序和并行机会

## 前置条件

已使用 `/plan` 生成：
- `specs/<feature>/plan.md`（Status 为 Ready）

## 产出物

```
specs/<feature>/
├── spec.md              ← 特性规格（specify + clarify 产出）
├── plan.md              ← 实施计划（plan 产出）
└── tasks.md             ← 任务清单（本命令产出）
```

## 工作流程

```
/tasks [优先级偏好（可选）]
    ↓
前置检查 → 确认 plan Ready
    ↓
提取上下文（plan 设计表 + spec RQ 列表）
    ↓
Phase 划分（按 kratos-patterns 工作步骤顺序）
    ↓
任务生成（RQ 追踪 + 并行标记 + 文件路径）
    ↓
质量校验 → 汇报结果
```
