# Seata Spec

## 事务 helper 边界

| 条件 | 做法 |
|------|------|
| 通用事务封装 | 放 seata/helper |
| 业务补偿、状态判断、流程分支 | 不放事务 helper |

```go
// ✅ 通用事务封装
type Transaction interface {
    InTx(ctx context.Context, fn func(ctx context.Context) error) error
}
```

```go
// ❌ helper 里塞业务补偿
func InTx(ctx context.Context, fn func(context.Context) error) error {
    if order.Status == ... { ... }
}
```

---

## 使用边界

| 条件 | 做法 |
|------|------|
| 事务边界定义 | 仍由 UseCase 决定 |
| Repo | 只负责原子读写，不定义业务事务步骤 |

---

## 常见错误模式

```text
// ❌ seata/helper 中加入领域状态判断
```

```text
// ❌ 用事务 helper 替代业务编排
```
