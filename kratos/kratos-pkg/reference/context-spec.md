# Context 参考

## 这个主题解决什么问题

统一业务上下文的构造与读取，避免租户、语言、viewer 等请求级基础信息在各层散落传递。

## 适用场景

- 需要把 tenant、language、viewer 等请求级信息放入上下文
- 需要在 middleware、service、biz、depend 之间传递稳定上下文

## 推荐结构或实现方式

- 在 `internal/pkg/context` 中提供 `NewBusinessContext` 和 `WithXxx` 入口。
- `context` 负责 typed helper，底层 metadata 透传通过 `internal/pkg/metadata` 收口。

## 标准模板

```go
func NewBusinessContext(ctx context.Context, opts ...BusinessOption) context.Context {
    if GetLocalize(ctx) == nil {
        ctx = WithLocalize(baseenum.LanguageDefault())(ctx)
    }
    for _, opt := range opts {
        ctx = opt(ctx)
    }
    return ctx
}

func WithTenant(id uint32) BusinessOption {
    return func(ctx context.Context) context.Context {
        ctx = metadata.SetTenantID(ctx, id)
        ctx = SetTenant(ctx, &Tenant{ID: id})
        return ctx
    }
}
```

## Good Example

- `WithTenant` 同时维护 typed context 与 metadata，调用方只感知统一入口。
- `WithLocalize` 通过语言枚举创建 localizer，而不是在业务层手动拼装。

## 常见坑

- 直接在业务代码里写裸 `context.WithValue`
- 同一语义只写入 typed context，不同步 metadata，导致跨边界丢失

## 相关 rule / 相关 reference

- `../rules/context-rule.md`
- `./metadata-spec.md`
- `./middleware-spec.md`
