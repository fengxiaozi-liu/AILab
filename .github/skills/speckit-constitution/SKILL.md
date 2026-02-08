---
name: speckit-constitution
description: 创建或更新项目宪法（.github/.specify/memory/constitution.md），并同步校验关键运行时文档。
---

# 项目宪法 Skill

## 适用场景

- 项目初始化，首次建立治理原则
- 原则变更（安全、质量门禁、发布策略、流程约束）

## 输入

- 用户给出的治理要求或修订指令
- 现有宪法：`.github/.specify/memory/constitution.md`
- 项目规则：`.github/rules/project/*.md`
- 仓库事实：`README.md`、当前项目具体内容

## 输出

- 更新后的宪法：`.github/.specify/memory/constitution.md`
- 更新后的规则：`.github/rules/project/*.md`
- 更新后的运行资产：`.github/agents/*.md`、`.github/prompts/*.md`
- 宪法顶部 Sync Impact Report（HTML 注释）
- 变更总结：
  - 版本变化（old -> new）与升级理由
  - 修改的原则列表（含 Principle-ID）
  - 受影响文件与待办项

## 目标

- 维护“项目宪法（Project Constitution）”实例文档，不输出空泛模板。
- 宪法仅保留原则与治理，不重复 rules 细则。
- 确保宪法可被其他 Agent 精确引用（`[KRATOS_PRINCIPLE_*]`）。
- 确保 Rule Application 明确指向 `.github/rules/project/*.md` 与其 Good/Bad demo。

## 执行流程

1. 读取与盘点
- 读取现有宪法（若不存在则创建实例骨架：核心原则 + 治理）。
- 识别现有 Principle-ID、版本号、批准日期、最后修订日期。
- 读取用户本次输入的新增/修订规则（若有）。
- 读取 `.github/rules/project/*.md` 提取可执行约束与 demo 语义。

2. 归纳原则（项目实例化）
- 用户输入优先；上下文用于补全和约束校验。
- 将 project rules 归并到 5 条核心原则（或用户指定数量）。
- 每条原则必须包含：
  - `Principle-ID: [KRATOS_PRINCIPLE_*]`
  - `说明`（为什么）
  - Rule Source（且仅一个文件，映射到 `.github/rules/project/*.md`）
- 若原则与现有宪法冲突，优先保留兼容性并记录修订原因。

3. 治理落地
- 治理章节必须包含：
  - 宪法优先级与作用域
  - Rule Application（`.github/rules/project/*.md` + Good/Bad demo 为落地依据）
  - 更新联动要求（宪法更新时联动更新 rules/project 与运行资产）
  - 跨 Agent 引用规范（必须用 Principle-ID）
  - 修订流程
  - 版本策略（MAJOR/MINOR/PATCH）

4. 版本与日期决策
- MAJOR：原则删除/重定义/不兼容治理变更。
- MINOR：新增原则/章节或实质扩展约束。
- PATCH：措辞澄清、排版修复、非语义修订。
- 日期统一 `YYYY-MM-DD`：
  - 批准日期（首次建立后保持不变）
  - 最后修订日期（有变更则更新为当天）

5. 写入 Sync Impact Report
- 写入内容至少包含：
  - Version change（old -> new）
  - Change type（MAJOR/MINOR/PATCH）
  - Principles changed（按 Principle-ID）
  - Governance changes
  - Deferred TODOs（如有）

6. 一致性传播检查
- 核对 `.github/copilot-instructions.md` 是否声明宪法路径与定位。
- 核对 `.github/rules/project/*.md` 是否与最新宪法章节映射一致。
- 核对关键 agent（如 `planner`）是否要求执行宪法检查。
- 核对相关 prompts/commands 是否存在过期引用。

7. 最终校验
- 宪法中不得存在未解释的占位符。
- Principle-ID 唯一且稳定。
- Rule Application 必须存在。
- 每条原则的 Rule Source 必须唯一（不得多文件并列）。
- 宪法中不得重复 rules 细则内容。
- 版本号与 Sync Impact Report 一致。
- 语言必须具体、可执行，避免空话。

## 强约束

- 不得输出与现有项目规则冲突的新“平行规则体系”。
- 宪法是治理抽象，不替代具体编码规范细则。
- 治理章节必须包含 Rule Application：明确 `.github/rules/project/*.md` 的规则与 demo 为落地依据。
- 不得将 `.github/rules/ai/*` 作为项目宪法原则来源。
- 缺失关键信息时，使用 `TODO(<FIELD>): ...` 并在 Sync Impact Report 中列出。
