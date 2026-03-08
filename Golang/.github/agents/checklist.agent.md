---
name: checklist
description: 为当前特性生成需求质量或验收检查清单。
handoffs:
  - label: Generate Plan
    agent: plan
    prompt: Checklist complete, continue with implementation planning
    send: true
---

## 用户输入

```text
$ARGUMENTS
```

若非空必须纳入（可能是检查域、验收重点或约束）。

## 核心说明

checklist 用于从需求视角生成可检查、可复审的质量清单。

本 Agent 固定调用 **speckit-checklist** 完成当前阶段工作。
项目技能与语言技能由 Agent 在运行时识别并决定是否补充加载。

## 执行流程

1. 读取 spec 与必要上下文
2. 调用 `speckit-checklist`
3. 生成或更新 `specs/<feature>/checklists/<domain>.md`
4. 输出清单路径与检查范围
5. 完成后交接 **plan**

## 行为规则

- 检查项必须可检查、可复审
- 只审需求质量，不验证实现行为
