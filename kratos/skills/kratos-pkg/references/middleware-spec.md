# Middleware Spec

## 中间件边界

| 条件 | 做法 |
|------|------|
| 通用链路处理 | 放 middleware，如恢复、错误格式化、metadata 注入 |
| 业务判断、repo 查询 | 不放 middleware |

```go
// ✅
func Client() []middleware.Middleware {
    return []middleware.Middleware{
        Recovery(),
        FormatError(),
        metadata.Client(),
        tracing.Client(),
    }
}
```

```go
// ❌ middleware 里查 repo
func Auth() middleware.Middleware { ... }
```

---

## 聚合方式

| 条件 | 做法 |
|------|------|
| 多个中间件按主题拆分 | 用聚合函数统一暴露链路 |
| 业务层到处拼链路 | 不允许 |

---

## 常见错误模式

```text
// ❌ 重复解析已由 context/metadata helper 收口的字段
```

```text
// ❌ 在 server 外随意拼 middleware 链
```
