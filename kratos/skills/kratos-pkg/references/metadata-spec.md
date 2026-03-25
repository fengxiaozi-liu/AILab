# Metadata Spec

## key 收口

| 条件 | 做法 |
|------|------|
| 一个 metadata 字段 | 提供成对 `GetXxx/SetXxx` helper |
| 同一语义多个近义 key | 不允许 |

```go
// ✅
func GetLanguage(ctx context.Context, opts ...Option) baseenum.Language { ... }

func SetLanguage(ctx context.Context, language baseenum.Language) context.Context {
    return metadata.AppendToClientContext(ctx, "x-md-global-language", string(language))
}
```

```go
// ❌ 在 service/repo/middleware 里硬编码 key
metadata.AppendToClientContext(ctx, "x-md-global-language", lang)
```

---

## 来源优先级

| 条件 | 做法 |
|------|------|
| client/server 都可能携带字段 | 明确优先级 |
| 来源不清晰 | 不要直接新增 helper |

---

## 常见错误模式

```text
// ❌ 同一语义多个 key
```

```text
// ❌ 读写入口散落
```
