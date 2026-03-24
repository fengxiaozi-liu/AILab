---
name: speckit-analyze
description: 用于 Speckit 研发流程的一致性与质量校验 (只读)，包括跨文档扫描 spec.md、plan.md 和 tasks.md 之间的内容对齐度与逻辑矛盾。适用于任务拆解完成但未开始写代码前，需要做整体方案审计与防呆防漏的质量门检查，挖掘潜在实施风险的场景。触发关键词包括 analyze、校验、一致性、对齐检查、缺漏分析、逻辑矛盾。
---

# Spec Kit Analyze Skill

## 何时使用

- `tasks.md` 生成后、实施前做只读质量门检查

## 输入

- 必需：`specs/<feature>/spec.md`、`plan.md`、`tasks.md`
- 模板：`./templates/analyze-template.md`
- 可选：当前 agent 指令文件与关键 skill rules

## 目标

在实施前识别三件套中的不一致、覆盖缺口、歧义和欠规范项。

## 约束

- 严格只读：禁止修改任何文件，仅输出报告
- 上位入口优先：与当前 agent 指令或显式 skill rules 冲突一律视为 CRITICAL

## 工作流

### 1. 前置检查

1. 确认 `spec.md`、`plan.md`、`tasks.md` 均存在
2. 若任一缺失，终止并提示先运行对应命令
3. 若当前 agent 指令文件存在，一并加载

### 2. 文档加载

按需读取关键段落：

| 文档 | 提取内容 |
|------|---------|
| spec.md | RQ、关键实体、约束与边界 |
| plan.md | Phase 0 分析、Phase 1 设计表、风险项 |
| tasks.md | Task ID、Phase、RQ 引用、文件路径 |
| 当前 agent 指令文件 | 上位执行要求、角色边界、资源入口与治理规则 |

### 3. 语义模型构建

1. RQ 追踪链：`RQ -> 设计项 -> Task`
2. 实体一致性：spec / plan / tasks 术语与实体是否一致
3. 沌理规则集：上位执行要求 + MUST/SHOULD 语句

### 4. 检测类别

- 覆盖缺口
- 不一致
- 重复
- 歧义
- 欠规范
- Agent 与规则对齐
- 风险交叉

### 5. 严重级别

| 级别 | 定义 |
|------|------|
| CRITICAL | 违反当前 agent 指令 / skill rules MUST；核心 RQ 零覆盖；三件套缺失关键文档 |
| HIGH | RQ 冲突；关键设计项无任务；关键 Phase 顺序错误 |
| MEDIUM | 术语漂移；路径不规范；非关键覆盖不足 |
| LOW | 表述优化；小型格式问题 |

### 6. 输出报告

- 基于 `./templates/analyze-template.md`
- 写入 `specs/<feature>/checklists/analyze.md`
- 输出统计、发现表、RQ 覆盖汇总、Agent 与规则对齐、下一步建议

## 输出

- `specs/<feature>/checklists/analyze.md`
- 完成后提示用户确认分析结论，无遗漏后继续执行 `/specify` 生成正式 spec 或 `/plan` 制定计划。
