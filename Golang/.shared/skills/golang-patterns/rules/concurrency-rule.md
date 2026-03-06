# Concurrency Rule

## Principles

- goroutine 必须可退出，共享状态必须同步。

## Specification

- 启动 goroutine 前定义退出机制。
- 使用 mutex 或 atomic 保护共享状态。

## Prohibit

- 禁止无退出条件 goroutine。
- 禁止未明确所有权就关闭 channel。
