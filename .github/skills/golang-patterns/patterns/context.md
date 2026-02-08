# Go Context（语言层）

## Scope

只覆盖 Go `context` 的语言层用法：参数传递、取消/超时、不要滥用全局/结构体存储 ctx。

### 1) `context.Context` 作为第一参数传递

```go
type Repo interface {
    Get(ctx context.Context, id uint32) (*T, error)
}
```

### 2) `WithTimeout/WithCancel` + `defer cancel()`

```go
ctx, cancel := context.WithTimeout(parent, timeout)
defer cancel()
```

## Rule hints

- `ctx` 仅用于传递取消信号/超时/请求级元信息；不要作为业务数据容器滥用
- 除非明确生命周期，否则不要把 `context.Context` 存进 struct 字段长期持有
