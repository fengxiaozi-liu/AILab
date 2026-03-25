# Schema Spec

## schema 工具边界

| 条件 | 做法 |
|------|------|
| 通用结构提取、反射辅助、字段检查 | 可以放 schema |
| 业务校验、业务配置判断 | 不放 schema |

```go
// ✅ 通用结构提取/检查 helper
func ExtractFields(v any) []string { ... }
```

```go
// ❌ 业务专用校验
func ValidateAccountOpenSchema(v any) error { ... }
```

---

## 复用门槛

| 条件 | 做法 |
|------|------|
| 多处结构处理都需要 | 才抽 schema helper |
| 一次性业务场景 | 留在原位置 |

---

## 常见错误模式

```text
// ❌ schema helper 混入业务配置判断
```

```text
// ❌ 为一次性场景创建不可复用 schema 工具
```
