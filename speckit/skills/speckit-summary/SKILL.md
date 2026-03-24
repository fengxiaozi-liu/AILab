---
name: speckit-summary
description: |
  用于 Speckit 研发流程的闭环复盘与交付确认，包括汇总梳理研发全过程 (spec/plan/tasks/review) 并导出 summary.md 交付报告。适用于单阶段需求任务做收尾，需向上级汇报或跨团队交接核心变更点清单、接口文档索引总结和关键验证结果的场景。触发关键词包括 summary、总结、复盘、交付报告、变更汇总。
---

# Spec Kit Summary Skill

## 何时使用

- `/code-review` 审查通过后，spec 流程的最后一步。
- 产出物 `specs/<feature>/summary.md` 与 spec/plan/tasks/review 并列存放。

## 输入

| 文档 | 必需 | 提取内容 |
|------|------|---------|
| `specs/<feature>/spec.md` | ✅ | 需求背景、目标用户、RQ 列表、约束与边界 |
| `specs/<feature>/plan.md` | ✅ | 接口清单、状态机、关键决策（DR-xxx）、依赖服务 |
| `specs/<feature>/tasks.md` | ✅ | 产出文件路径、任务完成率、RQ 覆盖率 |
| `specs/<feature>/review.md` | ✅ | 审查评分、遗留问题（CR-xxx） |

## 约束

- **只读提取**：不修改任何源文档。
- **忌重复叙述**：summary 是压缩归档，不是四份文档的翻译；丢弃推导过程，保留决策结论。
- **结构固定**：严格按下方模板生成，不增减章节。
- **数字精确**：任务数、RQ 数、文件数、评分均从源文档直接读取，不估算。

## 工作流

### §前置检查

1. 确认四份源文档均存在，任一缺失则终止并说明原因。
2. 若 tasks.md 有未完成任务（`[ ]`），警告并注明，仍可继续生成。
3. 若 review.md 无审查结论，警告并注明。

### §上下文加载

按顺序读取四份文档，提取关键信息：

| 来源 | 提取项 |
|------|-------|
| spec.md | 背景动机段落、RQ 编号+标题（不含详细描述）、约束与边界列表 |
| plan.md | 项目类型、接口设计表（方法/URL/说明）、依赖服务表、DR-xxx 决策结论、状态流转规则 |
| tasks.md | Phase 列表、已完成/总任务数、每个任务的产出文件路径、RQ 覆盖数 |
| review.md | 各维度评分、总分、是否通过、CR-xxx 遗留问题（级别+描述） |

### §文档生成

1. 读取模板 `./templates/summary-template.md`。
2. 按下方填充规则逐节替换占位符，写入 `specs/<feature>/summary.md`。

**填充规则**：

| 章节 | 来源 | 处理方式 |
|------|------|---------|
| 需求背景 | spec.md 背景动机 | 压缩为 2-3 句话 |
| 产出文件清单 | tasks.md 产出路径 | 按分类建表，末尾附统计行 |
| 接口清单 | plan.md 接口设计表 | 保留方法名、URL、一句话说明；按 C端/Admin/Inner 分组 |
| 状态机 | plan.md 状态流转 | ASCII 图或文字描述，仅含状态值和流转路径 |
| 关键设计决策 | plan.md DR-xxx | 每条只保留结论一行，省略理由和排除方案 |
| 依赖服务 | plan.md 依赖服务表 | 直接复制，精简到服务/用途/状态三列 |
| 本期边界 | spec.md 约束与边界 + 澄清结果 | bullet 列表，只含明确"不做"的项 |
| 质量结论 | review.md 审查摘要 | 评分表 + CR-xxx 遗留列表（级别/描述/建议） |

## 输出

- 模板路径：`./templates/summary-template.md`
- 完工总结：`specs/<feature>/summary.md`
- 完成后提示用户 review `summary.md`，确认无误后本 feature 开发流程结束。
