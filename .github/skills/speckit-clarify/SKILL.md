---
name: speckit-clarify
description: 对已有 spec 进行深度澄清，逐题对话并回写文档。
---

# Spec Kit Clarify Skill（中文）

## 何时使用

- `specs/<feature>/spec.md` 已存在，存在待澄清项（CQ）需要深度对话。
- specify 完成后，作为下游澄清工具接管。

## 输入

- `specs/<feature>/spec.md`（由 specify 生成）
- 用户给出的澄清方向或约束

## 工作流

### §spec 加载

1. 读取 `specs/<feature>/spec.md`。
2. 解析已有 CQ 条目，识别状态为 `[待澄清]` 的项。
3. 若 spec 不存在，终止并提示先运行 `/specify`。

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
   - specify 阶段遗留的 `[待澄清]` CQ
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
   - 最多 **8 个** CQ（含 specify 遗留 + 新增总计）
   - 优先生成影响架构和建模的维度问题

### §澄清回写

用户每确认一个答案，立即增量回写：

1. 将答案填入对应 CQ 的"澄清结果"字段
2. 同步更新**需求描述区**的相关章节内容
3. 若与旧表述冲突，替换旧表述而非重复堆叠
4. 检查答案是否触发新的澄清需求，如有则追加 CQ（不超过 8 个总量）
5. 全部 CQ 澄清后，将文档 Status 改为 `Ready`

### §写后校验

每次回写后执行：

1. CQ 条目与已确认答案一一对应，无重复
2. CQ 总数不超过 8
3. 目标歧义已消除、无残留冲突
4. 需求描述区与澄清结果一致（无矛盾内容）
5. Markdown 结构有效

## 输出

- 更新后的 `specs/<feature>/spec.md`

