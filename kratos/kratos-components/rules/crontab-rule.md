# Crontab Rule

## Principles

- 定时任务必须先考虑幂等、重试、并发与补偿。

## Specification

- 任务执行要有超时和失败处理。
- 并发执行策略必须显式。
- 补偿任务要能重复执行而不破坏数据。

## Prohibit

- 禁止无幂等保证的周期性写操作直接上线。
