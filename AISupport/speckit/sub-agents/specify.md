---
name: specify
description: 从自然语言特性描述创建或更新特性规格文档。
---


## 核心说明

specify 是需求入口阶段，用于从自然语言输入生成或更新结构化 spec。

本 Agent 固定调用 **speckit-specify** 完成当前阶段工作。

## 执行流程

1. 读取用户输入与必要项目上下文
2. 调用 `speckit-specify`
3. 生成或更新 `specs/<feature>/spec.md`
4. 输出分支名、spec 路径与待澄清项摘要
5. 若存在 CQ 则交接 **clarify**，否则交接 **plan**

## 行为规则

- 聚焦“用户需要什么、为什么”，不提前展开实现细节
- 缺失信息整理为 CQ，不擅自假设
