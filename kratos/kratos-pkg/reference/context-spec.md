# Context Spec

## Typed Helper 决策

| 所需能力 | 做法 |
|---------|------|
| 读写 tenant / viewer / language 等 | 使用 `internal/pkg/context` 的 `WithXxx`/`GetXxx` |
| 跨服务传递（RPC/EventBus） | 同时更新 metadata，不仅写 typed context |
| middleware 注入请求级基础信息 | 在 `NewBusinessContext` 中通过 opts 注入 |
| 业务层读取上下文字段 | 只调用 typed helper，不直接 `context.Value` 取 |

---

## WithXxx 双写模式

```go
// ✅ WithTenant 同时维护 typed context 和 metadata（保证跨边界传递）
func WithTenant(id uint32) BusinessOption {
    return func(ctx context.Context) context.Context {
        ctx = metadata.SetTenantID(ctx, id)   // ← 写入 metadata，跨 RPC 透传
        ctx = SetTenant(ctx, &Tenant{ID: id}) // ← 写入 typed context，本进程读取
        return ctx
    }
}

// ❌ 只写入 typed context，不写 metadata
func WithTenant(id uint32) BusinessOption {
    return func(ctx context.Context) context.Context {
        return SetTenant(ctx, &Tenant{ID: id})  // ❌ 跨 RPC 后对方读不到
    }
}

// ❌ 裸 context.WithValue
ctx = context.WithValue(ctx, "tenant_id", id)  // ❌ 类型不安全，无法枚举
```

---

## NewBusinessContext 初始化

```go
// ✅ middleware 通过 NewBusinessContext 统一注入
func Tenant() middleware.Middleware {
    return func(handler middleware.Handler) middleware.Handler {
        return func(ctx context.Context, req interface{}) (interface{}, error) {
            tenantID := metadata.GetTenantID(ctx)
            ctx = context2.NewBusinessContext(ctx,
                context2.WithTenant(tenantID),
            )
            return handler(ctx, req)
        }
    }
}

// ✅ 语言默认值在 NewBusinessContext 中兜底
func NewBusinessContext(ctx context.Context, opts ...BusinessOption) context.Context {
    if GetLocalize(ctx) == nil {
        ctx = WithLocalize(baseenum.LanguageDefault())(ctx)  // ✅ 默认语言兜底
    }
    for _, opt := range opts {
        ctx = opt(ctx)
    }
    return ctx
}
```

---

## 组合场景

```go
// 完整：middleware 初始化 context → 业务层只用 typed helper
// 1. middleware 层
func Tenant() middleware.Middleware {
    return func(handler middleware.Handler) middleware.Handler {
        return func(ctx context.Context, req interface{}) (interface{}, error) {
            tenantID := metadata.GetTenantID(ctx)
            ctx = context2.NewBusinessContext(ctx,
                context2.WithTenant(tenantID),
                context2.WithViewer(metadata.GetViewerID(ctx)),
            )
            return handler(ctx, req)
        }
    }
}

// 2. 业务层（UseCase/Service）
func (u *AccountUseCase) GetAccount(ctx context.Context, id uint32) (*Account, error) {
    tenant := context2.GetTenant(ctx)      // ✅ 用 typed helper
    viewer := context2.GetViewer(ctx)      // ✅ 用 typed helper
    _ = tenant
    _ = viewer
    return u.accountRepo.GetAccount(ctx, id)
}
```

---

## 常见错误模式

```go
// ❌ 裸 context.WithValue
ctx = context.WithValue(ctx, "language", "zh-CN")

// ❌ 只写 typed context 不写 metadata（跨边界丢失）
ctx = SetTenant(ctx, &Tenant{ID: id})  // ❌ 没有 metadata.SetTenantID

// ❌ 在 Repo 层直接解析 metadata，绕过 context helper
func (r *accountRepo) GetAccount(ctx context.Context, id uint32) (*Account, error) {
    md, _ := metadata.FromServerContext(ctx)
    tenantID := md.Get("x-md-global-tenant-id")  // ❌ 应用 context2.GetTenant(ctx)
    ...
}
```
