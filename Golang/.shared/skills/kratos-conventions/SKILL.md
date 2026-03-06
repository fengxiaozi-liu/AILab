---
name: kratos-conventions
description: |
  用于 Kratos 横切规范，包括错误码、错误语义、枚举、状态流转、i18n、日志与注释一致性。适用于新增或修改 error、enum、i18n、logging、comment、not found 处理、状态机分支或多语言文案的场景。触发关键词包括 error、错误码、not found、enum、状态、i18n、locale、logging、日志、comment、注释、一致性。
---

# Kratos Conventions

## 必读规则

- `./rules/error-rule.md`
- `./rules/enum-rule.md`
- `./rules/i18n-rule.md`
- `./rules/logging-rule.md`
- `./rules/comment-rule.md`

## 按需参考

- Error：`./reference/error-spec.md`
- Enum：`./reference/enum-spec.md`
- i18n：`./reference/i18n-spec.md`
- Logging：`./reference/logging-spec.md`
- Comment：`./reference/comment-spec.md`

## 读取顺序

先读 `./rules/*.md` 明确语义、日志、注释与一致性约束，再按当前任务读取必要的 `./reference/*.md`。

## 何时使用

- 新增或修改错误码、错误返回语义、not found 处理
- 新增或调整枚举值、状态机分支
- 新增或修改 i18n key、多语言文案
- 新增或修改日志打印方式、日志字段、日志脱敏策略
- 新增或修改 proto、Ent、Service、Repo、UseCase 等注释说明

## 核心约束

1. 错误语义优先，禁止把“必须存在”的 not found 传染成 `nil, nil`。
2. 枚举和值域变化必须同步检查 proto、DB、分支逻辑和状态流转完整性。
3. i18n key 命名保持稳定，避免把文案直接散落进业务流程。
4. 日志必须结构化、可搜索、可脱敏，避免日志风暴和敏感信息泄漏。
5. 注释用于表达业务语义、边界和职责，不用于重复代码字面意思。

## 实施流程

1. 先定义业务语义，再决定错误码、枚举、文案 key、日志字段和注释内容。
2. 核对各层返回值、分支、序列化边界和日志上下文是否一致。
3. 如果涉及 proto、Ent、生成物或对外契约变更，联动 `kratos-entry` 或 `kratos-components` 做验证。

## 强制输出

开始前输出：

- `ConventionScope:` error / enum / i18n / logging / comment
- `SemanticChange:` 本次语义变更摘要

提交前输出：

- not found 是否按语义转换（Yes/No）
- 枚举分支是否完整覆盖（Yes/No）
- i18n key 和文案是否与业务语义一致（Yes/No）
- 日志是否避免敏感信息和日志风暴（Yes/No）
- 注释是否与当前语义一致且未污染生成物（Yes/No）
