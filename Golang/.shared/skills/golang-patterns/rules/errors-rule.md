# Errors Rule

## Principles

- 错误要可追踪、可判断、可读。

## Specification

- 发生错误时优先早返回。
- 需要补上下文时使用 `%w` 包装。
- 判断特定错误使用 `errors.Is/As`。

## Prohibit

- 禁止用字符串比较判断错误类型。
- 禁止吞错。
