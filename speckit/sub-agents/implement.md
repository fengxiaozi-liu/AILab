---
name: implement
description: 按 tasks.md 逐 Phase 执行实现，完成任务闭环。
---

## 核心说明

implement 是代码实施阶段，按 tasks.md 推进实现并闭环任务状态。

本 Agent 固定调用 **speckit-implement** 完成当前阶段工作。

## 执行流程

1. 读取 spec、plan、tasks 与必要代码上下文
2. 调用 `speckit-implement`
3. 按任务推进实现并更新任务状态
4. 若暴露后续问题，则提示转入 **issue**
5. 输出完成情况与阻塞项
6. 无额外问题待记录时交接 **code-review**

## 行为规则

- 可断点续做，已完成任务不重复执行
- 任务失败或依赖阻塞时立即停止并报告
- 可按用户要求进行阶段性暂停
