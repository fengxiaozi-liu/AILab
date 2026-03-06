---
description: 按 tasks.md 逐 Phase 执行实现，完成任务闭环。
agent: implement
---

# /implement 命令

调用 **implement** agent 按任务清单逐步实施代码。

## 何时使用

- tasks.md 已生成，analyze 质量门通过
- 准备开始正式编码

## 前置条件

已生成完整三件套：
- `specs/<feature>/spec.md`
- `specs/<feature>/plan.md`
- `specs/<feature>/tasks.md`

## 产出物

- 代码实现改动（按 tasks.md 中的文件路径）
- 更新后的 `tasks.md`（完成的任务标记 `[x]`）

## 工作流程

```
/implement [范围约束（可选）]
    ↓
前置检查 → 确认三件套齐全
    ↓
加载上下文（tasks + plan + kratos-patterns reference）
    ↓
检测断点（已有 [x] 则续做）
    ↓
┌─ Phase N 执行 ──────────────┐
│  读取对应 kratos-patterns 规范  │
│  逐任务编写代码                │
│  每任务完成后 tasks.md 打勾     │
│  Phase 完成 → 暂停汇报        │
│  等用户确认 → 下一 Phase       │
└─────────────────────────────┘
    ↓
完成校验 → 全部任务打勾、编译通过
```

## 示例

```
User: /implement

Agent:
已加载 tasks：`specs/manage-portfolio/tasks.md`
项目类型：业务项目
总任务：19，已完成：0

开始 Phase 1：Ent Schema（2 个任务）

✅ T001 [RQ-001] 编写 Portfolio Ent Schema `internal/data/ent/schema/portfolio.go`
✅ T002 [RQ-002] 编写 PortfolioItem Ent Schema `internal/data/ent/schema/portfolio_item.go`

Phase 1 完成（2/2）。下一步：Phase 2 Ent 代码生成。
继续？

WAIT_CONFIRM
```

```
User: 继续

Agent:
Phase 2：Ent 代码生成（1 个任务）

✅ T003 执行 Ent 代码生成 — 生成成功，无报错

Phase 2 完成（1/1）。下一步：Phase 3 Biz 层实现。
继续？

WAIT_CONFIRM
```
