# Context Rule

## Principles

- `context` 用于请求级取消、超时和元信息传递。

## Specification

- `ctx context.Context` 作为跨层函数第一参数。
- `WithTimeout/WithCancel` 后及时 `cancel()`。

## Prohibit

- 禁止在请求路径随意改用 `context.Background()`。
- 禁止长期持有 `context.Context` 到结构体字段。
