# Logging Rule

## Principles

- 日志必须结构化、可搜索、可脱敏。

## Specification

- 项目统一基于 Kratos `log.Logger` 与 `log.Helper` 组织日志，不散落构造底层 logger。
- 长生命周期对象优先注入 `log.Logger`，在对象内部转换为 `*log.Helper` 复用。
- 需要固定模块标识时，优先使用 `log.With(logger, "module", "...")` 追加结构化字段。
- 需要请求上下文时，优先使用 `helper.WithContext(ctx)` 输出日志。
- 错误日志至少包含动作语义、关键业务标识和错误对象。
- 生命周期日志、预加载日志、依赖调用日志应保持简短稳定，方便检索。

## Prohibit

- 禁止打印 token、secret、password、PII、证件号或完整请求体。
- 禁止在高频循环、批处理热点路径无控制地输出日志。
- 禁止同一模块混用多套日志风格，导致字段和上下文不一致。

## Checklist

- 是否统一使用 `log.Logger` / `log.Helper`？
- 是否通过 `WithContext(ctx)` 透传上下文？
- 是否避免敏感信息和高频日志风暴？
