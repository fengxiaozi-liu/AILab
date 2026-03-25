---
name: speckit-fix
description: 基于当前 feature 的 issue.md 做受边界约束的问题修复。用于 issue 已记录，且需要一次性完成 clarify、plan、tasks、implement 修复闭环时。
---

# Speckit Fix

## 何时使用

- 当前 feature 已存在 `specs/<feature>/issue.md`
- `issue.md` 中记录了本次需求实施过程中暴露的问题
- 需要围绕这些问题做一次受控修复，而不是扩展成新需求

## 职责边界

- `speckit-fix` 只修复 `issue.md` 中记录的当前 feature 问题
- `/fix` 必须在一次执行中按 Clarify -> Plan -> Tasks -> Implement 的顺序完成闭环
- 修复工作不得偏离当前 feature 的既有目标和边界
- 若 issue 本质上是新需求或范围扩张，必须停止并要求单独建 feature

## 输入

- 必需：`specs/<feature>/issue.md`
- 必需：`specs/<feature>/spec.md`
- 可选：`specs/<feature>/plan.md`
- 可选：`specs/<feature>/tasks.md`
- 可选：`specs/<feature>/review.md`
- 参考：`./references/fix-workflow.md`

## 工作流

### 1. 前置检查

1. 读取 `issue.md`，确认存在待修复问题
2. 读取原始 `spec.md`，把它作为本次修复的最高边界约束
3. 若存在 `plan.md` / `tasks.md` / `review.md`，一并读取用于定位问题来源
4. 明确说明本次修复只处理 `issue.md` 中列出的当前 feature 问题
5. 优先从人工填写的最小信息中提取问题事实，再补全 issue 的其余结构字段

### 2. Clarify 阶段

对 `issue.md` 中每个待修复问题做边界澄清：

- 问题是否描述清楚
- 期望行为是否明确
- 修复边界是否明确
- 是否会引出新的非本次需求内容

若存在未决项：

- 先在修复工作上下文中补齐澄清结论
- 不允许跳过澄清直接规划和实施
- Clarify 完成后直接进入 Plan，不要求用户额外再次调用命令

### 3. Plan 阶段

为本次修复生成受控修复计划：

- 聚焦问题根因、受影响任务、受影响文件和验证方式
- 计划只覆盖 issue 中的问题，不引入额外功能
- 必须显式列出“修复边界”与“明确不做”
- 计划可作为内部推演过程使用，不强制单独落盘

### 4. Tasks 阶段

把修复计划拆成最小可执行任务：

- 任务应能映射回 issue 项
- 每个任务都要注明目标文件或产物
- 不允许把 unrelated cleanup 混入修复任务
- 任务可作为内部推演过程使用，不强制单独落盘

### 5. Implement 阶段

按修复任务执行代码修改与验证：

- 仅修改修复问题所需的最小范围
- 按任务状态真实回写
- 完成后同步更新 `issue.md` 中对应问题的修复状态

## 输出

- 更新后的 `specs/<feature>/issue.md`
- 完成后提示进入 `/code-review`
