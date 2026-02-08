# Generics（语言层）

## Scope

只覆盖 Go 泛型（1.18+）在本项目的“语言层使用边界”：用于通用工具/结构，避免在业务域滥用。

### 1) 约束（constraints）表达可用操作

```go
type Compare[T constraints.Integer | constraints.Float] struct{ Value T }
```

### 2) 泛型用于通用 util（责任链/工具结构）

```go
type Handle[T, V any] func(ctx context.Context, req T, resp V) (bool, error)
```

## Rule hints

- 泛型优先用于 `internal/pkg` 这类通用层；业务实体/服务接口谨慎使用
- 约束尽量收敛（不要 `any` 到处传播，除非确实是通用容器）
