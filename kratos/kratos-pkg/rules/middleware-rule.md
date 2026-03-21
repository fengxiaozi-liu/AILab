# Middleware Rule

## Principles

- middleware 负责请求入口和出口的通用处理，不承载业务编排。
- middleware 与 `context`、`metadata` 必须协同工作，保持单向职责清晰。

## Specification

- middleware 可以负责恢复、错误格式化、localize 注入、viewer/tenant 注入、trace 和基础事务透传。
- middleware 应只做通用拦截、提取、补充、包装，不做聚合查询和业务状态判断。
- 中间件链顺序要稳定，先做上下文与基础信息注入，再做业务无关的通用包装。

## Prohibit

- 禁止在 middleware 中查询业务 repo、拼装业务 reply、执行业务状态流转。
- 禁止在不同 middleware 中重复解析相同 metadata 字段。
- 禁止让 middleware 成为 service/usecase 的旁路实现。

## Checklist

- 中间件逻辑是否只依赖基础能力。
- 是否复用了现有 metadata/context helper。
- 是否避免了业务查询和业务编排。
