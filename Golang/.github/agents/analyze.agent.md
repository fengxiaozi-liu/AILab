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

analyze 是 tasks 的下游质量门，用于对 spec、plan、tasks 做只读交叉检查。

本 Agent 固定调用 **speckit-analyze** 完成当前阶段工作。
项目技能与语言技能由 Agent 在运行时识别并决定是否补充加载。

## 执行流程

1. 读取 spec、plan、tasks
2. 调用 `speckit-analyze`
3. 输出分析报告
4. 若通过则交接 **implement**

## 行为规则

- 严格只读，不修改任何文件
- 可给出修复建议，但不自动执行
