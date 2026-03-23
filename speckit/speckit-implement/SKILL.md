---
name: speckit-implement
description: USE FOR：按任务清单推进编码落地、逐 Phase 执行实现与验证、断点续做已有 tasks.md、构建验证与闭环。
  触发关键词：implement、实现任务、落地、写代码、编码、按任务执行、开始开发。
  DO NOT USE FOR：需求分析与拆解（用 speckit-plan / speckit-tasks）、架构设计（用 speckit-design）、代码评审（用 speckit-code-review）。
---

# Spec Kit Implement Skill

## 输入

- 必需：用户输入，包括目标、范围、约束、可跳过项、验收口径
- 可选：`specs/<feature>/tasks.md`
- 可选：`specs/<feature>/plan.md`
- 可选：`specs/<feature>/spec.md`
- 用户约束，例如范围、跳过的 Phase 等

缺少必需输入时，MUST 先向用户提问收集目标与验收口径，不得跳过直接编码。

## 工作流

### 1. 读取 tasks.md 状态

1. 若存在 `specs/<feature>/tasks.md`：
   - 解析所有 Phase、TaskID、描述、文件路径、完成标记 `[x]` / `[ ]`
   - 定位第一个 `[ ]` 任务，确定续做起点
   - 若存在已完成任务，汇报续做起点（跳过了哪些任务）
2. 若 `tasks.md` 不存在：
   - 基于用户输入和可选 `plan.md` / `spec.md` 生成最小任务清单草案
   - 写入 `specs/<feature>/tasks.md`，再进入执行循环
3. 若存在 `plan.md` 或 `spec.md`：读取设计边界与验收口径，作为每条 task 实施时的约束

### 2. 确定本轮执行范围

根据用户输入确定执行模式和范围：

**执行模式：**

| 用户输入 | 执行模式 | Phase 边界行为 |
|---------|---------|---------------|
| 未指定 / 「逐阶段」/ 「按 Phase」 | 逐阶段模式 | 每个 Phase 完成后暂停，等待确认 |
| 「一次性」/ 「全部执行」/ 「不用确认」 | 连续模式 | 所有 Phase 连续执行，完成后统一汇报 |

**执行范围：**

| 用户输入 | 执行范围 |
|---------|----------|
| 未指定 | 从第一个 `[ ]` 任务开始，执行全部剩余任务 |
| 指定 Phase | 仅执行该 Phase 内的 `[ ]` 任务 |
| 指定 TaskID | 仅执行该单个任务 |
| 指定跳过项 | 跳过对应 task，继续执行其余 `[ ]` 任务 |

### 3. 逐 task 执行循环

对每个待执行的 `[ ]` 任务，按以下步骤处理：

1. **读取 task**：从 `tasks.md` 中取出 TaskID、描述、文件路径
2. **加载技能**：根据 task 描述和文件路径，识别该 task 所需的项目适配 skill 或语言 skill，并加载对应 SKILL.md；若无法识别，则按仓库事实和 `plan.md` 约束继续
3. **实施**：依据 task 描述、文件路径、已加载技能的 MUST 规则及 `plan.md` 约束完成编码
4. **立即标记**：在 `tasks.md` 中将该 task 改为 `[x]`，然后再处理下一个 task
5. **失败处理**：若 task 执行失败，立即停止，报告错误详情，等待用户决策

### 4. Phase 边界处理

每个 Phase 内所有 `[ ]` 任务执行完毕后，根据执行模式处理：

**逐阶段模式：**
1. 汇报本 Phase 完成的 TaskID 列表和产出文件
2. 暂停，等待用户确认后再进入下一 Phase
3. 若用户要求调整（修改某 task 或跳过），回到步骤 3 继续循环

**连续模式：**
1. 记录本 Phase 完成的 TaskID 列表和产出文件
2. 直接进入下一 Phase，不暂停
3. 所有 Phase 全部完成后，统一汇报各 Phase 产出和完成校验结果

### 5. 完成校验

所有 Phase 的 task 均为 `[x]` 后执行：

| # | 校验项 | 标准 |
|---|--------|------|
| 1 | 任务完成 | `tasks.md` 中所有任务均为 `[x]` |
| 2 | 需求覆盖 | 若 spec 存在，每条 RQ 均有对应 task 落地；若无 spec，每条临时需求点均有对应落地 |
| 3 | 设计一致 | 若 plan 存在，代码与 `plan.md` 一致；若无 plan，代码与任务清单一致 |
| 4 | 构建通过 | 代码可编译、构建或通过必要验证 |

## 输出

- 代码实现改动
- 更新后的 `specs/<feature>/tasks.md`

## 约束

### MUST
- 每完成一个 task，MUST 立即将 `tasks.md` 中对应条目改为 `[x]`，不得攒完再批量标记
- 每个 task 实施前，MUST 先识别并加载该 task 所需的技能，再按技能 MUST 规则约束实现
- task 执行失败时 MUST 立即停止并报告错误，不得跳过继续执行后续 task
- 逐阶段模式下，每个 Phase 内所有 task 完成后，MUST 暂停等待用户确认再进入下一 Phase
- 连续模式下，MUST 在所有 Phase 完成后统一汇报，不得中途暂停

### MUST NOT
- MUST NOT 在 `tasks.md` 不存在且用户目标不清晰时直接开始编码，必须先生成最小任务清单草案
- MUST NOT 实施任何 task 前不先读取 `tasks.md` 中的 `[ ]`/`[x]` 状态
- MUST NOT 跨 Phase 预执行（当前 Phase 未完成前不开始下一 Phase）

### SHOULD
- SHOULD 在 Phase 汇报中列出完成的 TaskID 和产出文件
- SHOULD 在续做时明确告知跳过了哪些已标记 `[x]` 的任务及续做起点

## 强制输出

进入执行循环前输出：

```json
{
  "tasksFile": "specs/<feature>/tasks.md 或 draft-generated",
  "mode": "逐阶段 | 连续",
  "resumeFrom": "TaskID 或 Phase 名称",
  "skippedCompleted": ["已跳过的 TaskID 列表"],
  "scope": "全部 / 指定 Phase / 指定 TaskID"
}
```

每个 Phase 完成后输出：

```json
{
  "phase": "Phase 名称",
  "completedTasks": ["T001", "T002"],
  "outputFiles": ["文件路径"],
  "nextPhase": "下一 Phase 名称或 done",
  "waitingForConfirmation": true
}
```

## 参考文件

| 文件 | 适用场景 | 加载时机 |
|------|----------|----------|
| `specs/<feature>/tasks.md` | 任务清单、Phase 列表、完成标记 | 前置检查时必读（若存在） |
| `specs/<feature>/plan.md` | 设计边界、验收口径、执行约束 | 上下文加载时读取（若存在） |
| `specs/<feature>/spec.md` | 需求与验收口径 | 上下文加载时读取（若存在） |
