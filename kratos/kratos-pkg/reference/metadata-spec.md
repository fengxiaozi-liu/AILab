# Metadata Spec

## 来源优先级

| 来源 | 优先级 | 说明 |
|------|--------|------|
| Server metadata | 高 | 由网关注入，可信 |
| Client metadata | 低 | 发出端附带，可被 server 覆盖 |

```go
// ✅ GetMetadata：先合并 client，再用 server 覆盖
func GetMetadata(ctx context.Context, opts ...Option) (metadata.Metadata, bool) {
    cfg := &Config{FromClient: true, FromServer: true}
    for _, opt := range opts { opt(cfg) }

    if cfg.FromClient && cfg.FromServer {
        clientMD, _ := metadata.FromClientContext(ctx)
        serverMD, _ := metadata.FromServerContext(ctx)
        merged := metadata.Join(clientMD, serverMD)  // server 优先，后写的 key 覆盖前
        return merged, true
    }
    ...
}
```

---

## GetXxx / SetXxx 成对 Helper

```go
// ✅ 提供成对 helper，调用方不直接操作 key
func GetLanguage(ctx context.Context, opts ...Option) baseenum.Language {
    language := baseenum.LanguageDefault()
    if md, ok := GetMetadata(ctx, opts...); ok {
        if v := md.Get("x-md-global-language"); v != "" {
            language = baseenum.Language(v)
        }
    }
    return language
}

func SetLanguage(ctx context.Context, language baseenum.Language) context.Context {
    return metadata.AppendToClientContext(ctx, "x-md-global-language", string(language))
}

// ❌ 调用方直接写死 metadata key
ctx = metadata.AppendToClientContext(ctx, "x-md-global-language", "zh-CN")  // ❌

// ❌ 多处用不同 key 表达同一语义
ctx = metadata.AppendToClientContext(ctx, "language", "zh-CN")    // ❌ key 不一致
ctx = metadata.AppendToClientContext(ctx, "x-language", "zh-CN")  // ❌
```

---

## 字段清单

```go
// ✅ 已有字段的标准 helper（不要手写 key），例如：
tenantID  := metadata.GetTenantID(ctx)
viewerID  := metadata.GetViewerID(ctx)
language  := metadata.GetLanguage(ctx)
platform  := metadata.GetPlatform(ctx)
clientIP  := metadata.GetClientIP(ctx)

// ✅ 新字段需要同时提供 GetXxx 和 SetXxx，并对齐 x-md-global-* 前缀
func GetCustomField(ctx context.Context) string {
    md, ok := GetMetadata(ctx)
    if !ok { return "" }
    return md.Get("x-md-global-custom-field")
}
func SetCustomField(ctx context.Context, val string) context.Context {
    return metadata.AppendToClientContext(ctx, "x-md-global-custom-field", val)
}
```

---

## 组合场景

```go
// 完整：middleware 设置 → Depend 透传 → 跨服务接收
// 1. Server middleware（接收端）
func Tenant() middleware.Middleware {
    return func(handler middleware.Handler) middleware.Handler {
        return func(ctx context.Context, req interface{}) (interface{}, error) {
            tenantID := metadata.GetTenantID(ctx)         // ✅ 用 helper 读
            ctx = context2.NewBusinessContext(ctx,
                context2.WithTenant(tenantID),
            )
            return handler(ctx, req)
        }
    }
}

// 2. Depend（发出端）透传 tenant/language
func (d *AdminUserDepend) GetUsers(ctx context.Context, ids []uint32) ([]*AdminUser, error) {
    ctx = metadata.SetLanguage(ctx, context2.GetLanguage(ctx))  // ✅ 透传 language
    reply, err := d.client.GetUsers(ctx, &req)
    ...
}
```

---

## 常见错误模式

```go
// ❌ 手写 metadata key
md.Get("x-language")  // ❌ 应用 metadata.GetLanguage(ctx)

// ❌ 同一语义多个 key 并存
// service A 用 "x-md-global-tenant-id"
// service B 用 "tenant_id"
// ❌ 导致透传断链

// ❌ 只写 typed context，不写 metadata（跨边界丢失）
ctx = context2.SetTenant(ctx, tenant)  // ❌ 没有 metadata.SetTenantID 那么跨服务读不到
```
