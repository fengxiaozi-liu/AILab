---
name: speckit-clarify
description: |
  用于 Speckit 研发流程的需求澄清与探讨，包括分析需求中不清晰或矛盾的点并收敛生成 CQ (Clarification Questions) 问题清单，确认后回写 spec。适用于面对逻辑缺失或极度模糊的需求，需要通过提问与用户互动以消除歧义并产出最小可行草案的场景。触发关键词包括 clarify、澄清需求、CQ、追问、消除歧义。
---

# Spec Kit Clarify Skill（中文）

## 何时使用

- 有 `specs/<feature>/spec.md`：存在待澄清项（CQ）需要深度对话与回写。
- 没有 spec：用户只有自然语言需求/想法，需要先通过澄清问题把边界、约束、验收口径问清楚。
- specify 完成后可接管；specify 之前也可先用本技能做澄清，再转 `/specify` 落 spec。

## 输入

- 可选：`specs/<feature>/spec.md`（若存在则读取并回写）
- 可选：`specs/<feature>/checklists/*.md`（若存在则合并 checklist 未完成项）
- 必需：用户给出的需求描述/澄清方向或约束（无 spec 时作为主输入）
- 用户给出的澄清方向或约束

## 工作流

### §spec 加载

本技能支持两种模式：

1. **Spec 澄清模式**（存在 `specs/<feature>/spec.md`）：
   - 读取 `specs/<feature>/spec.md`
   - 解析已有 CQ 条目，识别状态为 `[待澄清]` 的项
2. **输入澄清模式**（spec 不存在）：
   - 以用户输入为主生成 CQ 清单
   - 不强制生成/回写 spec，但可以在用户同意后将澄清结果交给 `/specify` 生成最小 spec 草案

### §覆盖评估

按 9 大评估维度扫描 spec，标记各维度覆盖状态：

| # | 评估维度 | 状态标记 |
|---|---------|---------|
| 1 | 功能范围与行为 | Clear / Partial / Missing |
| 2 | 领域与数据模型 | Clear / Partial / Missing |
| 3 | 交互与 UX 流程 | Clear / Partial / Missing / 不适用 |
| 4 | 非功能属性 | Clear / Partial / Missing |
| 5 | 外部集成依赖 | Clear / Partial / Missing |
| 6 | 边界与失败处理 | Clear / Partial / Missing |
| 7 | 约束与权衡 | Clear / Partial / Missing |
| 8 | 术语一致性 | Clear / Partial / Missing |
| 9 | 完成信号 | Clear / Partial / Missing |

- 仅对当前特性**实际涉及的维度**评估，不涉及的标记"不适用"
- 后端项目中"交互与 UX 流程"通常标记"不适用"

### §问题生成

1. 合并两个来源的待澄清项：
   - specify 阶段遗留的 `[待澄清]` CQ / checklist阶段未标记完成的`[x]`CHK
   - 覆盖评估新发现的 Partial/Missing 维度
2. 生成/补充 CQ 条目：
   - 编号：接续已有 CQ 编号（如已有 CQ-001~003，从 CQ-004 开始）
   - 类别：对应的评估维度名称
   - 详细描述：为什么需要澄清
   - 可选项：2-4 个合理选项 + 推荐项及理由
   - 澄清结果：初始为 `[待澄清]`
3. 约束：
   - 只问会影响**架构、数据建模、任务拆分、测试设计、运维、合规**的问题
   - 每题必须可用"选项"或"≤5 词短答"回答
   - 优先生成影响架构和建模的维度问题

### §澄清回写

用户每确认一个答案，立即增量回写或沉淀：

1. 将答案填入对应 CQ 的"澄清结果"字段
2. Spec 澄清模式：同步更新**需求描述区**的相关章节内容
3. 输入澄清模式：把答案汇总到“澄清结论摘要”（对话输出即可），用于后续 `/specify`/`/plan`/`/tasks`
3. 若与旧表述冲突，替换旧表述而非重复堆叠
4. 检查答案是否触发新的澄清需求，如有则追加 CQ
5. Spec 澄清模式：全部 CQ 澄清后，将文档 Status 改为 `Ready`
6. 若待澄清来源于 checklist，澄清完成后标记 check 项为 `[x]`

### §写后校验

每次回写后执行：

1. CQ 条目与已确认答案一一对应，无重复
2. 目标歧义已消除、无残留冲突
3. 需求描述区与澄清结果一致（无矛盾内容）
4. Markdown 结构有效

## 输出

- Spec 澄清模式：更新后的 `specs/<feature>/spec.md`
- 输入澄清模式：CQ 澄清清单 + 澄清结论摘要（可选：转交 `/specify` 生成最小 spec）

