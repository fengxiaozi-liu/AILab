# Kit Spec

## data kit 与 biz kit

| 条件 | 做法 |
|------|------|
| data 层基础设施聚合 | 放 `internal/data/kit` |
| biz 层稳定抽象接口 | 放 `internal/biz/kit` |

```go
// ✅ data kit
type Data struct {
    Db       *ent.Client
    EventBus *eventbus.EventBus
}
```

```go
// ✅ biz kit
type Transaction interface {
    Transaction(ctx context.Context, fn func(ctx context.Context) error) error
}
```

```go
// ❌ 把 data 聚合和 biz 抽象混成一个目录概念
```

---

## 复用优先

| 条件 | 做法 |
|------|------|
| 已有 `internal/data/kit` / `internal/biz/kit` 能力 | 先复用 |
| 只有单一业务调用点 | 不下沉成 kit |

```text
// ✅ 先检索
rg -n "type .*Data|Transaction|EventBus|Client|ProviderSet" internal/data/kit internal/biz/kit
```

```go
// ❌ 为一个业务流程硬抽 kit
type AccountSyncKit struct {}
```

---

## 调用边界

| 条件 | 做法 |
|------|------|
| data 层依赖基础设施 | 依赖 `*kit.Data` |
| biz 层依赖事务抽象 | 依赖 `kit.Transaction` |

```go
// ✅
func NewAccountRepo(data *kit.Data) openbiz.AccountRepo {
    return &accountRepo{data: data}
}
```

```go
// ✅
type AccountUseCase struct {
    tx kit.Transaction
}
```

```go
// ❌ kit 里写业务状态判断
if account.Status == ...
```

---

## 组合场景

```text
internal/data/kit -> infrastructure aggregation
internal/biz/kit  -> stable interface abstraction
callers reuse existing capability before adding new one
```

这个组合场景同时满足：

- kit 边界清晰
- 没有业务流程下沉
- 复用优先于抽象

---

## 常见错误模式

```text
// ❌ data/biz kit 混用
```

```text
// ❌ kit 中塞业务流程
```

```text
// ❌ 明明已有能力还重造
```
