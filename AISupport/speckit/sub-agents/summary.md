---
name: summary
description: 读取 spec/plan/tasks/review 四份文档，生成特性完工总结 summary.md。
---

## 用户输入

若非空必须纳入（可能是 feature 名或交付范围）。

## 核心说明

summary 是交付总结阶段，用于从已完成的流程文档中提取关键信息并生成总结。

本 Agent 固定调用 **speckit-summary** 完成当前阶段工作。

## 执行流程

1. 读取 spec、plan、tasks、review 与必要上下文
2. 调用 `speckit-summary`
3. 生成或更新 `specs/<feature>/summary.md`
4. 输出 summary 路径与主要章节

## 行为规则

- 只读源文档，只写 summary.md
- 不对源文档做主观评价，只做信息提取与压缩
