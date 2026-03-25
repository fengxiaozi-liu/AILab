# 问题记录：[FEATURE NAME]

**Feature**: `specs/[feature]`
**Created**: [DATE]
**Status**: Open | Clarifying | Ready for Fix | Fixed
**Source**: 本次需求实施闭环后的问题归档

---

## 记录范围

- 仅记录本次 feature 直接相关的问题
- 不记录与当前 feature 无关的新需求
- 不把重构愿望、未来优化或旁支需求混入本文件

---

## 上下文引用

- Spec: `specs/[feature]/spec.md`
- Plan: `specs/[feature]/plan.md`
- Tasks: `specs/[feature]/tasks.md`
- Review: `specs/[feature]/review.md`

---

## 问题清单

### ISSUE-001: [问题标题]

> **人工填写说明**:
> 只填写以下 4 块：问题标题、现象、期望行为、复现或证据。
> 其余字段由 AI 在后续修复流程中补全。

**现象**:
[描述当前出现了什么问题]

**期望行为**:
[描述原本应该如何表现]

**复现或证据**:
- [日志、报错、审查结论、人工验证结果]

---

> **以下内容由 AI 补全，不要求人工填写**

**类型**:
[Bug | Gap | Regression | Ambiguity | Task Miss]

**严重级别**:
[Blocker | Major | Minor]

**相关任务**:
[T00x, T00y | 无]

**相关需求**:
[RQ-00x, RQ-00y | 无]

**相关文件**:
[`[path/to/file]` | 无]

**影响范围**:
[谁会受影响、影响到什么流程]

**修复边界**:
- [允许修改的范围]
- [明确不做的范围，防止偏离本次需求]

**待澄清问题（如有）**:
- [需要先确认的问题；没有则写“无”]

**修复状态**:
[待澄清] | [待计划] | [待拆解] | [待实施] | [已修复]

---

## 修复约束

- 修复必须服务于当前 feature 的既有目标
- 修复不得扩展为新的独立需求
- 若问题本质上属于新需求，必须退出本流程并单独建 feature
