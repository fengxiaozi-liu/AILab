---
name: analyze
description: 对 spec/plan/tasks 三件套做只读一致性与质量分析。
handoffs:
  - label: Start Implementation
    agent: implement
    prompt: Analysis passed, start implementation by phases
    send: true
---

## 用户输入

```text
$ARGUMENTS
```

若非空必须纳入（可能是特定关注点）。

## 核心说明

analyze 是 tasks 的下游质量门，对 spec/plan/tasks 三件套做只读交叉检查，不修改任何文件。

本 Agent 加载 **speckit-analyze** skill。

## 执行流程

> 领域方法论详见 **SKILL: speckit-analyze**。

1. 前置检查 → *SKILL §前置检查*
2. 文档加载 → *SKILL §文档加载*
3. 语义模型构建 → *SKILL §语义模型构建*
4. 检测流程 → *SKILL §检测流程*（7 类检测，最多 50 条）
5. 输出报告 → *SKILL §输出报告*
6. 下一步建议 → *SKILL §下一步建议*

## 行为规则

- **严格只读**：禁止修改任何文件
- 可给出修复建议，但不自动执行
- 宪法冲突一律 CRITICAL
