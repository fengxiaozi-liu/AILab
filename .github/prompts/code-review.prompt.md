---
description: 对已实现代码做完成度、架构合规、性能与安全审查。
agent: code-review
---

# /code-review 命令

调用 **code-review** agent 对已实现代码做质量审查。

## 何时使用

- `/implement` 执行完毕后，提交前做代码质量审查
- 修复问题后再次审查，直到通过

## 前置条件

已完成代码实现：
- `specs/<feature>/tasks.md` 存在，且有 `[x]` 标记的已完成任务
- `specs/<feature>/plan.md` 存在

## 产出物

审查报告 `specs/<feature>/review.md`，包含：
- 4 维度评分（完成度、架构合规、性能、安全，各 10 分）
- 发现清单（按严重级别排序）
- 修复建议优先级

## 工作流程

/code-review [关注维度（可选）]
    ↓
前置检查 → 确认 tasks.md + plan.md 存在，有已完成任务
    ↓
项目类型判定 → 确定适用的架构规范
    ↓
加载文档 + 读取实际代码文件
    ↓
4 维度审查（完成度 → 架构合规 → 性能 → 安全）
    ↓
生成评分报告 + 下一步建议
