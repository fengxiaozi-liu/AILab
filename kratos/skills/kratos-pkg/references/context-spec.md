# Context Spec

## typed helper

| 条件 | 做法 |
|------|------|
| 请求级基础信息透传 | 用 `WithXxx` + `GetXxx` typed helper |
| 大对象、聚合关系、临时结果 | 不放进 context |

```go
// ✅
type BusinessOption func(context.Context) context.Context

func WithTenant(id uint32) BusinessOption {
    return func(ctx context.Context) context.Context {
        return metadata.SetTenantID(ctx, id)
    }
}
```

```go
// ❌ 裸 WithValue key
context.WithValue(ctx, "tenant_id", id)
```

---

## 与 metadata 联动

| 条件 | 做法 |
|------|------|
| 字段需要跨边界透传 | 同步评估 metadata/middleware |
| 只在进程内短链路使用 | 不强行进 metadata |

---

## 常见错误模式

```text
// ❌ 把业务聚合对象塞进 context
```

```text
// ❌ 裸 key 到处散落
```
