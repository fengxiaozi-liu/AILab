---
name: summary
description: 读取 spec/plan/tasks/review 四份文档，生成特性完工总结 summary.md。
---

## 用户输入

```text
$ARGUMENTS
```

若非空，则为 feature 分支名（如 `account-open`），用于定位 `specs/<feature>/` 目录。
若为空，则从当前编辑文件路径或对话上下文推断 feature 名。

## 核心说明

summary 是 spec 流程的终态步骤，在 code-review 通过后执行。

本 Agent 加载一个 Skill：
- **speckit-summary**：提取规则、填充逻辑与输出模板

## 执行流程

1. 解析 feature 名 → 确定 `specs/<feature>/` 路径
2. 前置检查 → *SKILL §前置检查*
3. 上下文加载 → *SKILL §上下文加载*（读取四份源文档）
4. 压缩提取 → 按 SKILL 填充规则从各文档提取关键信息
5. 生成文档 → *SKILL §文档生成*（写入 `specs/<feature>/summary.md`）
6. 输出确认 → 告知用户文件路径与主要章节

## 行为规则

- tasks.md 有未完成任务 → 警告但不终止，summary 中注明未完成项
- review.md 无审查结论 → 警告但不终止，质量结论章节注明"待审查"
- **禁止修改源文档**：只读 spec/plan/tasks/review，只写 summary.md
- summary.md 已存在时覆盖写入（每次生成完整文档）
- 不对源文档内容做主观评价，只做信息提取与压缩
