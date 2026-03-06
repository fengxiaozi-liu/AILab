# Defer Resource Rule

## Principles

- 资源与锁要成对释放。

## Specification

- 获取资源或加锁成功后第一时间 `defer` 释放。

## Prohibit

- 禁止多分支手写释放导致遗漏。
- 禁止在高频循环里无必要滥用 `defer`。
