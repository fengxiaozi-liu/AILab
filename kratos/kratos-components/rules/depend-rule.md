# Depend Rule

## Principles

- 跨服务依赖统一封装，优先批量调用。

## Specification

- InnerRPC 调用必须设置超时。
- 重试仅用于幂等操作，并采用退避策略。
- 批量收集 ID，避免逐条远程调用。

## Prohibit

- 禁止散落直连上游。
- 禁止循环里逐条远程请求。
