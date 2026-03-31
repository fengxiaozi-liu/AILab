---
name: issue
description: 创建、补齐或追加当前 feature 的 issue.md。
---

## 核心说明

issue 是问题记录阶段，用于维护当前 feature 的 `issue.md`，为后续 **fix** 提供标准化修复文档。

本 Agent 固定调用 **speckit-issue** 完成当前问题文档工作。

## 执行流程

1. 读取当前 feature 上下文与已有 `issue.md`
2. 调用 `speckit-issue`
3. 若无问题说明，则初始化或补齐 spec 风格的 issue 模板
4. 若有问题说明，先写入 `用户问题`，再立即结构化为新的或已有的 `ISSUE-xxx` 单元
5. 输出更新结果，并提醒用户后续执行 **fix**

## 行为规则

- 保留已有 issue 单元与人工内容，不做整文重写
- 无问题说明时不自动创建空白 issue 条目
- 有问题说明时不得只追加原文，必须按顺序追加或补全结构化 `ISSUE-xxx` 内容
- slash 命令参数中的问题描述视为有效输入，必须在本次执行中直接完成规格化
- issue 文档只负责记录问题与修复工作面板，不直接执行修复
