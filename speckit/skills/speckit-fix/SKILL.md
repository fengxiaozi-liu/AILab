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

- `speckit-fix` 修复问题来源文档中记录的当前 feature 问题，来源文档可以是 `issue.md` 或 `review.md`
- `/fix` 必须在一次执行中按 `Specify -> Clarify -> Plan -> Tasks -> Implement` 的顺序完成闭环
- 仅当来源文档中有相关执行计划与任务之后才允许进行 implement，实现工程更改
- `/fix` 默认处理来源文档中所有未完成的问题项
- `/fix` 负责把澄清、计划、任务与执行结果持续回写到原问题文档
- `/fix` 使用参考模板生成 `澄清记录`、`修复计划`、`修复任务` 章节
- 修复工作不得偏离当前 feature 的既有目标和边界
- 若问题本质上是新需求或范围扩张，必须停止并要求单独建 feature

## 输入

- 必须：`specs/<feature>/issue.md` 或 `specs/<feature>/review.md`
- 必须：`specs/<feature>/spec.md`
- 可选：`specs/<feature>/plan.md`
- 可选：`specs/<feature>/tasks.md`
- 可选：另一份问题文档（`issue.md` 或 `review.md`）
- 参考工作流：`./references/fix-workflow.md`
- 输出模板：`./assets/clarify-template.md`
- 输出模板：`./assets/plan-template.md`
- 输出模板：`./assets/tasks-template.md`

## 工作流

### 1. 前置检查

1. 识别本次修复的来源文档是 `issue.md` 还是 `review.md`
2. 读取来源文档，确认存在未完成的问题项
3. 读取原始 `spec.md`，把它作为本次修复的最高边界约束
4. 若存在 `plan.md` / `tasks.md` / 另一份问题文档，一并读取用于定位问题来源
5. 读取 `./references/fix-workflow.md`、`./assets/clarify-template.md`、`./assets/plan-template.md`、`./assets/tasks-template.md`
6. 明确说明本次修复只处理来源文档中列出的当前 feature 问题
7. 按来源文档中的问题顺序扫描所有未完成项，判断每项当前处于哪个阶段

### 2. Specify 阶段

对每一个未完成问题项：

- 先根据用户原生问题，对该问题做结构化
- 将原生问题整理并回填到：
  - `问题说明`
  - `当前现象`
  - `期望行为`
  - `证据或复现`
- 若已能直接理解，则在不打断用户的前提下直接写清这些字段
- 若无须澄清，则为该问题写入“无需澄清”的结果，并进入 `Clarify` 阶段
- 若仍存在关键歧义，必须：
  - 按 `./assets/clarify-template.md` 在来源文档中写入 `澄清记录`
  - 每个澄清问题必须提供互斥选项，便于用户直接二选一或多选一回复
  - 明确提醒用户进入对应问题文档补充澄清答案
  - 立即停止整个 `/fix` 运行，不得提前进入 `Plan`

### 3. Clarify 阶段

仅在以下条件满足时进入：

- 当前问题无待处理澄清项，或
- 用户已在来源文档的 `澄清记录` 中补充了答案

执行要求：

- 将用户填写的澄清答案回填为已澄清结论，并同步更新对应问题项的结构化字段
- 若澄清信息仍不足以安全规划，则继续按 `澄清记录` 模板补充互斥问题并停止
- 仅在当前问题已达到可规划状态时，才允许进入 `Plan` 阶段

### 4. Plan 阶段

按 `./assets/plan-template.md` 在来源文档中回写真正的修复计划：

- 计划必须覆盖本次修复涉及的问题、根因、影响范围、修复边界、决策依据与验证方式
- `根因分析` 必须说明历史 AI 决策为何形成，并明确其主要依据来源：`SKILL`、仓库事实或二者共同作用
- 对存在偏差的历史决策，必须继续定位偏差来源，包括 `SKILL` 约束、模板表达、仓库事实或二者映射关系
- `决策分析` 必须说明本次修复方案为何成立，并明确其主要依据来源：`SKILL`、仓库事实或二者共同作用
- 不引入问题范围外的新功能或大规模无关重构
- 明确写出：
  - 修复目标
  - 受影响的问题、模块或文件
  - 修复边界
  - 验证方法
  - 明确不做
- 仅当 plan 内容回写到来源文档时，认为 plan 阶段完成
- `Plan` 阶段完成后，进入 `Tasks` 阶段

### 5. Tasks 阶段

按 `./assets/tasks-template.md` 在来源文档中回写真正的 todo 项：

- 任务必须映射回对应问题项
- 任务必须是实现导向、可验证、可打勾的 todo 项
- 不允许把 unrelated cleanup 混入修复任务
- todo 项直接使用 `- [ ] T001 ...` 这种格式
- 仅当 tasks 内容回写到来源文档时，认为 task 阶段完成
- `Tasks` 阶段完成后，进入 `Implement` 阶段

### 6. Implement 阶段

按 `修复任务` 区块中的 todo 项执行代码修改与验证：

- 仅修改解决该问题所需的最小范围
- 按真实完成情况将对应任务从 `- [ ]` 更新为 `- [✅️]`
- 全部任务完成后，将对应问题在全局问题列表中的状态更新为 `[✅️]`
- 若来源文档存在 `待解决问题` / `未解决问题` 分组，则将已解决问题移入 `已解决问题`
- 完成信号以任务打钩为准，不再单独维护额外的任务完成字段

## 停止与恢复

- 若在 `Specify` 或 `Clarify` 阶段发现仍需用户澄清，必须停止当前 `/fix`
- 停止前必须先把以下内容写回来源文档：
  - 问题摘要
  - `澄清记录`
  - 当前阶段
- 用户补充后，下一次 `/fix` 不得从头重建问题项
- 应从当前问题的 `Clarify` 阶段继续，再进入 `Plan -> Tasks -> Implement`

## 输出

- 更新后的来源问题文档：`specs/<feature>/issue.md` 或 `specs/<feature>/review.md`
- 对已完成任务的 todo 项打 `[✅️]`
- 若某个问题已完成修复，则同步更新全局问题列表状态
- 若全部闭环完成，提示进入 `/code-review`
