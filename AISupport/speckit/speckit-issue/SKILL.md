---
name: speckit-issue
description: 创建、初始化、补齐或追加当前 feature 的 issue.md。用于维护 `specs/[feature]/issue.md` 的用户问题、全局问题列表与 ISSUE详情列表，并为后续 `/fix` 提供 ISSUE-xxx 问题单元。
---

# Speckit Issue

## 何时使用

- 需要为当前 feature 创建 `specs/<feature>/issue.md`
- 需要按本技能目录下 `issue-template.md` 风格初始化问题修复文档
- 需要向已有 `issue.md` 追加新的 `ISSUE-xxx`
- 用户希望先记录问题，再交给 `/fix` 做完整修复闭环

## 职责边界

- `speckit-issue` 只负责创建和维护当前 feature 的 `issue.md`
- 负责把问题文档组织成可供 `/fix` 消费的结构化问题清单
- `issue.md` 使用与 `/fix` 相同的 `ISSUE-xxx` 容器结构，但 `/issue` 只负责问题登记层字段
- `issue.md` 必须保持本技能目录下 `./templates/issue-template.md` 的最小模板风格；若不符合则补齐缺失结构，但不得按模板整体覆盖已有 ISSUE
- 不负责实际修复，不负责执行 plan、tasks 或 implement
- 新问题以追加方式写入新的 `ISSUE-xxx`
- `问题状态` 属于 issue 层字段；新问题默认写为 `[ ]`
- 已解决问题仍可保留在 `ISSUE详情列表` 中，但其状态必须通过 ISSUE 内 `问题状态` 与全局问题列表同步表达
- 若已有 ISSUE 中存在 `/fix` 回填的 `问题状态`、`澄清记录`、`根因分析`、`影响范围`、`修复边界`、`决策分析`、`修复任务`，必须原样保留，不得删除、清空、重写或按模板回退

## 输入

- 必需：当前 feature 上下文，对应 `specs/<feature>/`
- 可选：`specs/<feature>/spec.md`
- 可选：`specs/<feature>/review.md`
- 可选：`/issue xxx` 中直接传入的问题说明 `xxx`
- 可选：`issue.md` 模板中的 `用户问题`
- 模板：本技能目录下 `./templates/issue-template.md`

## 工作流

### 1. 前置检查

1. 确认当前 feature 上下文已经明确，目标路径为 `specs/<feature>/issue.md`
2. 若存在已有 `issue.md`，先读取现有内容，保留已有人类编辑与既有 `ISSUE-xxx`
3. 按需读取 `spec.md` / `review.md`，仅用于补齐上下文引用
4. 明确说明本次输出基于本技能目录下 `./templates/issue-template.md`
5. 明确原始问题输入源有两类：`/issue xxx` 中的 `xxx`，或 `issue.md` 中 `用户问题` 区的文本
6. 若仅调用 `/issue` 且既没有 `xxx`、也没有待处理的 `用户问题`，则仅初始化模板并提醒用户去 `用户问题` 下补充描述

### 2. 对用户输入或 `用户问题` 进行结构化

- 若存在 `/speckit-issue {user_argument}` 的 `{user_argument}`，则将其视为本次原始用户问题输入
- 若不存在 `{user_argument}`，则从 `issue.md` 的 `用户问题` 区读取待结构化的问题输入
- 根据原始问题输入与现有 ISSUE 进行对照，判断是补充已有问题还是新增问题
- 原始问题输入一旦被成功沉淀为 `ISSUE-xxx`，应从 `用户问题` 区移除或避免重复写入，不得长期与结构化结果并存

### 3. 若 `issue.md` 不符合结构则先结构化原始文档

- 若 `issue.md` 不存在，则按本技能目录下模板 `./templates/issue-template.md` 初始化文档结构
- 若 `issue.md` 已存在但不符合当前模板结构，则只补齐缺失的模板段落与基础元数据
- 若原始 `issue.md` 不符合当前模板结构，则按本技能目录下 `./templates/issue-template.md` 初始化并整理为模板形式；结构化时应将原有问题内容迁移到模板对应位置，已符合结构的 ISSUE 保持不动，且不得回退 `/fix` 已补充的扩展字段

### 4. 不改动符合结构要求的既有 ISSUE，新增问题按结构化方式追加

- 既有 `ISSUE-xxx` 若已符合当前 ISSUE 结构则保持原样；若不符合则先结构化到当前模板要求后再继续补充
- 若用户输出明显对应已有 ISSUE，则补充该 issue 的结构化字段，不重复新建编号
- 若用户输出构成新问题，则按现有最大 `ISSUE-xxx` 顺延创建新的问题单元
- 若命中已有被 `/fix` 扩展过的 ISSUE，只允许补充或修正 issue 层字段，不得改动 fix 阶段字段
- 每个新增 issue 单元必须规格化为以下结构：
  - 问题状态
  - 问题说明
  - 当前现象
  - 期望行为
  - 证据或复现
- `/issue` 可调整的 issue 层字段仅限：`问题状态`、`问题说明`、`当前现象`、`期望行为`、`证据或复现`
- 新问题始终以追加方式写入新的 `ISSUE-xxx`
- 编号按现有最大 `ISSUE-xxx` 顺延
- 全局问题列表必须保留所有问题，并按 `ISSUE`、`问题简述`、`状态` 的列顺序展示
- 全局问题列表中的状态必须在表格单元格内使用 `[ ]` 表示未解决，使用 `[✅️]` 表示已解决
- 人工填写内容视为事实来源，只需补结构与最小必要元数据
- `/issue` 只触达 `用户问题`、`全局问题列表` 以及 ISSUE 内的 issue 层字段；不得删除或重排 `/fix` 已补充内容

## 5. 输出报告
对于 `specs/<feature>/issue.md`
- 模板定义的最小内容为：`用户问题`、`全局问题列表`、`ISSUE详情列表`
- 全局问题列表用于汇总所有问题及其解决状态，列顺序固定为 `ISSUE`、`问题简述`、`状态`
- 每个新增 issue 单元必须包含：
  - 问题状态
  - 问题说明
  - 当前现象
  - 期望行为
  - 证据或复现
- 若 ISSUE 中已存在由 `/fix` 补充的澄清、分析、边界或任务信息，必须完整保留

## 6. 输出

- 新建或更新后的 `specs/<feature>/issue.md`
- 若仅初始化模板，提醒用户去 `用户问题` 下补充描述
- 若已提供问题说明，说明已完成问题规格化，并同步更新全局问题列表与 `ISSUE详情列表`
