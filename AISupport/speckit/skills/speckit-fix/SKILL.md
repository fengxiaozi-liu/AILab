---
name: speckit-fix
description: 基于当前 feature 的问题文档做受边界约束的问题修复。用于 `issue.md` 或 `review.md` 中的问题已记录，且需要按 Specify、Clarify、Plan、Tasks、Implement 阶段回写原问题文档并完成修复闭环时。
---

# Speckit Fix

## 何时使用

- 当前 feature 已存在 `specs/<feature>/issue.md` 或 `specs/<feature>/review.md`
- 问题文档已通过 `/issue`、`/code-review` 或人工方式准备完成
- 问题文档中存在一个或多个未完成的问题项
- 需要围绕这些问题做一次受控修复，而不是扩展成新需求

## 职责边界

- `speckit-fix` 只修复来源文档中记录的当前 feature 问题，来源文档可以是 `issue.md` 或 `review.md`
- `/fix` 必须在一次执行中按 `Specify -> Clarify -> Plan -> Tasks -> Implement` 的顺序完成闭环
- `/fix` 默认处理来源文档中所有未完成的问题项
- `/fix` 负责把澄清、计划、任务与执行结果持续回写到原问题文档
- `/fix` 使用 `ISSUE详情列表` 作为持久化结构；每个 ISSUE 自带状态、根因、影响范围、修复边界、决策分析和修复任务
- 禁止用新的全局 `## 修复计划` 或 `## 修复任务` 覆盖历史批次内容
- 仅当每个待修复 ISSUE 内已有对应执行计划与任务之后才允许进行 implement
- 修复工作不得偏离当前 feature 的既有目标和边界
- 若问题本质上是新需求或范围扩张，必须停止并要求单独建 feature

## 输入

- 必须：`specs/<feature>/issue.md` 或 `specs/<feature>/review.md`
- 必须：`specs/<feature>/spec.md`
- 可选：`specs/<feature>/plan.md`
- 可选：`specs/<feature>/tasks.md`
- 可选：另一份问题文档（`issue.md` 或 `review.md`）
- 完整文档模板：本技能目录下 `./assets/fix-template.md`

## 文档结构

- 完整文档骨架使用本技能目录下 `./assets/fix-template.md`
- `## 全局问题列表` 与 `## ISSUE详情列表` 的组织方式以模板为准
- `Specify / Clarify / Plan / Tasks` 阶段的信息全部回写到对应 `### ISSUE-xxx` 详情块
- 若来源文档已有旧结构，先迁移到模板定义的结构中，不得删除既有人工内容或历史分析

## 工作流

### 1. 前置检查

1. 识别本次修复的来源文档是 `issue.md` 还是 `review.md`
2. 读取来源文档，确认存在未完成的问题项
3. 读取原始 `spec.md`，把它作为本次修复的最高边界约束
4. 若存在 `plan.md` / `tasks.md` / 另一份问题文档，一并读取用于定位问题来源
5. 读取本技能目录下 `./assets/fix-template.md`
6. 明确说明本次修复只处理来源文档中列出的当前 feature 问题
7. 按来源文档中的问题顺序扫描所有未完成项，判断每项当前处于哪个阶段

### 2. Specify 阶段

对每一个未完成问题项：

- 先检查该 ISSUE 是否已具备后续修复所需的 issue 层基础字段
- 若以下字段存在缺失或表述不清，再结合来源问题与上下文进行补齐：
  - `问题说明`
  - `当前现象`
  - `期望行为`
  - `证据或复现`
- 若现有信息已足够支撑后续修复，则在不打断用户的前提下直接补清缺失字段
- 若无须澄清，则在该 ISSUE 内的 `澄清记录` 写入“无需澄清”的结果，并进入 `Clarify` 阶段
- 若仍存在关键歧义，必须：
  - 按本技能目录下 `./assets/fix-template.md` 在对应 ISSUE 内写入 `澄清记录`
  - 每个澄清问题必须提供互斥选项，便于用户直接二选一或多选一回复
  - 明确提醒用户进入对应问题文档补充澄清答案
  - 立即停止整个 `/fix` 运行，不得提前进入 `Plan`

### 3. Clarify 阶段

仅在以下条件满足时进入：

- 当前问题无待处理澄清项，或
- 用户已在对应 ISSUE 的 `澄清记录` 中补充了答案

执行要求：

- 将用户填写的澄清答案回填为已澄清结论，并同步更新对应 ISSUE 的结构化字段
- 若澄清信息仍不足以安全规划，则继续在该 ISSUE 的 `澄清记录` 中补充互斥问题并停止
- 仅在当前问题已达到可规划状态时，才允许进入 `Plan` 阶段

### 4. Plan 阶段

按本技能目录下 `./assets/fix-template.md` 更新每个待修复 ISSUE 的详情字段：

- 计划信息必须写入对应 `### ISSUE-xxx` 内，计划信息不准游离在ISSUE之外
- `根因分析` 必须说明历史 AI 决策为何形成，并明确其主要依据来源：`SKILL`、仓库事实或二者共同作用
- 对存在偏差的历史决策，必须继续定位偏差来源，包括 `SKILL` 约束、模板表达、仓库事实或二者映射关系
- `决策分析` 必须说明本次修复方案为何成立，并明确其主要依据来源：`SKILL`、仓库事实或二者共同作用
- `影响范围`、`修复边界` 必须严格围绕该 ISSUE，不引入问题范围外的新功能或大规模无关重构
- `修复边界` 必须包含至少一条明确不做项
- 仅当 plan 内容回写到对应 ISSUE 详情时，认为该 ISSUE 的 plan 阶段完成
- `Plan` 阶段完成后，进入 `Tasks` 阶段

### 5. Tasks 阶段

按本技能目录下 `./assets/fix-template.md` 更新每个待修复 ISSUE 的 `修复任务` 字段：

- 任务必须写入对应 ISSUE 内， 任务列表不允许脱离ISSUE
- 一个 ISSUE 可以对应 N 个修复任务
- 每个任务必须显式标注所属 `ISSUE-xxx`
- 任务必须是实现导向或验证导向、可执行、可验证、可打勾
- 不允许把 unrelated cleanup 混入修复任务
- todo 项直接使用 `- [ ] T001 [ISSUE-xxx] ...` 格式
- 仅当 tasks 内容回写到对应 ISSUE 详情时，认为 task 阶段完成
- `Tasks` 阶段完成后，进入 `Implement` 阶段

### 6. Implement 阶段

按每个 ISSUE 内 `修复任务` 的待办项执行代码修改与验证：

- 仅修改解决该 ISSUE 所需的最小范围
- 按真实完成情况将对应任务从 `- [ ]` 更新为 `- [✅️]`
- 某个 ISSUE 的全部任务完成后，将该 ISSUE 的 `问题状态` 更新为 `[✅️]`
- 某个 ISSUE 的全部任务完成后，将其在 `全局问题列表` 中的状态同步更新为 `[✅️]`
- 完成信号以 ISSUE 内任务打钩为准，不再维护额外的全局任务完成字段

## 停止与恢复

- 若在 `Specify` 或 `Clarify` 阶段发现仍需用户澄清，必须停止当前 `/fix`
- 停止前必须先把以下内容写回来源文档：
  - 问题摘要
  - 对应 ISSUE 内的 `澄清记录`
  - 当前阶段
- 用户补充后，下一次 `/fix` 不得从头重建问题项
- 应从当前问题的 `Clarify` 阶段继续，再进入 `Plan -> Tasks -> Implement`
- 恢复执行时不得重建已有 ISSUE，也不得覆盖已完成 ISSUE 的 `根因分析`、`决策分析` 或 `修复任务`

## 输出

- 更新后的来源问题文档：`specs/<feature>/issue.md` 或 `specs/<feature>/review.md`
- 对已完成任务的 todo 项打 `[✅️]`
- 若某个问题已完成修复，则同步更新 ISSUE 内 `问题状态` 与全局问题列表状态
- 若全部闭环完成，提示进入 `/code-review`
