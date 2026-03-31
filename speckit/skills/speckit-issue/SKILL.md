---
name: speckit-issue
description: 创建、初始化、补齐或追加当前 feature 的 issue.md。用于维护 `specs/[feature]/issue.md` 的用户问题、全局问题列表与待解决问题，并为后续 `/fix` 提供 ISSUE-xxx 问题单元。
---

# Speckit Issue

## 何时使用

- 需要为当前 feature 创建 `specs/<feature>/issue.md`
- 需要按 spec 风格初始化问题修复文档
- 需要向已有 `issue.md` 追加新的 `ISSUE-xxx`
- 用户希望先记录问题，再交给 `/fix` 做完整修复闭环

## 职责边界

- `speckit-issue` 只负责创建和维护当前 feature 的 `issue.md`
- 负责把问题文档组织成可供 `/fix` 消费的结构化问题清单
- `issue.md` 只保留问题登记相关内容：`用户问题`、`全局问题列表`、`待解决问题`
- `issue.md` 必须保持 `./templates/issue-template.md` 模板风格，若不符合更新为模板风格
- 不负责实际修复，不负责执行 plan、tasks 或 implement
- 新问题以追加方式写入新的 `ISSUE-xxx`
- 已解决问题只保留在全局问题列表中，不保留在 `待解决问题` 章节中

## 输入

- 必需：当前 feature 上下文，对应 `specs/<feature>/`
- 可选：`specs/<feature>/spec.md`
- 可选：`specs/<feature>/review.md`
- 可选：用户提供的问题说明
- 模板：`./templates/issue-template.md`

## 工作流

### 1. 前置检查

1. 确认当前 feature 上下文已经明确，目标路径为 `specs/<feature>/issue.md`
2. 若存在已有 `issue.md`，先读取现有内容，保留已有人类编辑与既有 `ISSUE-xxx`
3. 按需读取 `spec.md` / `review.md`，仅用于补齐上下文引用
4. 明确说明本次输出基于 `./templates/issue-template.md`

### 2. 问题规格化与同步

- 根据用户输出与 `issue.md` 中现有未解决问题进行对照，判断是补充已有问题还是新增问题
- 若 `issue.md` 不存在，则按模板 `./templates/issue-template.md` 初始化文档结构
- 若 `issue.md` 已存在，则只补齐缺失的模板段落与基础元数据，不覆盖已有人工内容
- 若用户通过 `/issue 问题描述` 直接传入问题参数，必须先将该参数写入 `用户问题`，再立即按 ISSUE 单元结构完成规格化，不得只做原文追加
- 按照 `./templates/issue-template.md` 内容对文档进行规格化

### 3. ISSUE 单元规则

- 若用户输出明显对应已有未解决问题，则补充该 issue 的结构化字段，不重复新建编号
- 若用户输出构成新问题，则按现有最大 `ISSUE-xxx` 顺延创建新的问题单元
- 每个新增 issue 单元必须规格化为以下结构：
  - 问题说明
  - 当前现象
  - 期望行为
  - 证据或复现

### 4. 合并规则

- 既有 `ISSUE-xxx` 保持原样
- 新问题始终以追加方式写入新的 `ISSUE-xxx`
- 编号按现有最大 `ISSUE-xxx` 顺延
- 全局问题列表必须保留所有问题，并按 `ISSUE`、`问题简述`、`状态` 的列顺序展示
- 全局问题列表中的状态必须使用 `- [ ]` 表示未解决，使用 `- [✅️]` 表示已解决
- `待解决问题` 章节只保留未解决的 `ISSUE-xxx`
- 人工填写内容视为事实来源，只需补结构与最小必要元数据

## 5. 输出报告
对于 `specs/<feature>/issue.md`
- 模板中只保留：`用户问题`、`全局问题列表`、`待解决问题`
- 全局问题列表用于汇总所有问题及其解决状态，列顺序固定为 `ISSUE`、`问题简述`、`状态`
- `待解决问题` 中不得出现已解决的问题
- 每个新增 issue 单元必须包含：
  - 问题说明
  - 当前现象
  - 期望行为
  - 证据或复现

## 输出

- 新建或更新后的 `specs/<feature>/issue.md`
- 若未提供问题说明，提醒用户继续补充用户问题
- 若已提供问题说明，说明已完成问题规格化，并同步更新全局问题列表与待解决问题
