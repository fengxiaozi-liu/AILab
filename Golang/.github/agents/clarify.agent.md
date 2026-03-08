---
name: clarify
description: 对已有 spec 进行深度澄清，逐题与用户对话并将答案回写到 spec。
handoffs:
  - label: Review Requirement Quality
    agent: checklist
    prompt: Spec is ready, run requirement quality checklist before planning
    send: true
  - label: Generate Plan
    agent: plan
    prompt: Spec is ready, generate implementation plan
    send: true
---

## 用户输入

```text
$ARGUMENTS
```

若非空必须纳入（可能是澄清方向或约束）。

## 核心说明

clarify 是 specify 的下游阶段，用于收敛 spec 中的关键歧义。

本 Agent 固定调用 **speckit-clarify** 完成当前阶段工作。
项目技能与语言技能由 Agent 在运行时识别并决定是否补充加载。

## 执行流程

1. 读取 `specs/<feature>/spec.md` 与必要上下文
2. 调用 `speckit-clarify`
3. 增量回写 spec，保留已解决项与未解决项
4. 输出覆盖摘要与 spec 路径
5. spec 就绪后交接 **checklist** 或 **plan**

## 行为规则

- 一次聚焦一个关键问题
- 无关键歧义则明确反馈可继续
- 尊重用户提前结束
