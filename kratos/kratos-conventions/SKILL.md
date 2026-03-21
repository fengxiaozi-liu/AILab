---
name: kratos-conventions
description: |
  用于 Kratos 项目的横切规范，包括错误语义、防御式编程、枚举、i18n、日志与注释一致性。适用于修改错误处理、not found 行为、nil nil 返回、吞错、重复参数校验、枚举/状态流转、i18n、日志或注释的场景。触发关键词包括 error、not found、nil nil、吞错、重复校验、校验职责、error 直接返回、enum、状态、i18n、locale、logging、comment、一致性。
---

# Kratos Conventions

## 必读规则

- `./rules/error-rule.md`
- `./rules/defensive-programming-rule.md`
- `./rules/enum-rule.md`
- `./rules/i18n-rule.md`
- `./rules/logging-rule.md`
- `./rules/comment-rule.md`

## 参考资料

- Error: `./reference/error-spec.md`
- Enum: `./reference/enum-spec.md`
- i18n: `./reference/i18n-spec.md`
- Logging: `./reference/logging-spec.md`
- Comment: `./reference/comment-spec.md`

## 读取顺序

先读 `./rules/*.md`，锁定语义、边界与一致性约束，再按任务需要读取相关 `reference` 文件。

如果当前任务涉及防御式编程，例如 `nil nil`、吞错、重复校验、校验职责、`return err directly`，优先同时读取：

- `./rules/defensive-programming-rule.md`
- `./rules/error-rule.md`

职责分工：

- `defensive-programming-rule.md`：负责返回值语义、校验职责和禁止吞错
- `error-rule.md`：负责业务错误语义和错误转换

## 适用场景

- 新增或修改错误码、错误返回语义、not found 处理
- 新增或修改 nil/err 返回约定、吞错处理、重复参数校验
- 新增或修改枚举值、状态流转分支
- 新增或修改 i18n key、多语言文案
- 新增或修改日志结构、字段、脱敏策略
- 新增或修改 proto、Ent、Service、Repo、UseCase 中的注释

## 核心约束

1. 错误语义优先，不要把必须存在的 not found 传染成 `nil, nil`。
2. 防御式编程语义要收敛：默认避免重复校验、默认直接返回 `error`、禁止吞错。
3. 枚举和值域变化要联动检查 proto、DB、分支逻辑和状态流转完整性。
4. i18n key 要保持稳定，不要把原始文案散落到业务流程里。
5. 日志要结构化、可搜索、可脱敏，避免日志风暴和敏感信息泄漏。
6. 注释用于表达业务语义、边界和职责，不要重复显而易见的代码。

## 实施流程

1. 先定义业务语义，再决定错误码、返回语义、枚举使用、i18n key、日志字段和注释内容。
2. 核对跨层返回、分支、序列化边界和日志上下文是否一致。
3. 如果改动涉及 proto、Ent、生成物或对外契约，联动 `kratos-entry` 或 `kratos-components` 验证。

## 强制输出

开始前输出：

- `ConventionScope:` error / defensive / enum / i18n / logging / comment
- `SemanticChange:` 本次语义变更摘要

提交前输出：

- not found 是否按业务语义完成转换？(Yes/No)
- nil/err 返回是否消除了 `nil, nil` 歧义？(Yes/No)
- 是否存在吞错或忽略错误返回？(Yes/No)
- 是否存在重复参数校验？(Yes/No)
- 枚举分支是否完整覆盖？(Yes/No)
- i18n key 和文案是否与业务语义一致？(Yes/No)
- 日志是否避免敏感信息和日志风暴？(Yes/No)
- 注释是否与当前语义一致且未污染生成物？(Yes/No)
