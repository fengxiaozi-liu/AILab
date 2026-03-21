# Enum Rule

## Principles

- 枚举和值域变化必须保持状态流转和边界一致。

## Specification

- 新增枚举时同步检查 proto、DB、分支逻辑。
- switch 分支保持完整。

## Prohibit

- 禁止用裸字符串替代业务枚举。
