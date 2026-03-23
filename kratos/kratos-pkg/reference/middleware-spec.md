# Middleware Spec

## 职责范围决策

| 允许在 Middleware 中做 | 禁止在 Middleware 中做 |
|----------------------|----------------------|
| 读取 metadata / context 字段 | 查询 Repo / 数据库 |
| 注入 localize / tenant / viewer | 做业务判断（权限、状态流转） |
| 格式化错误响应（FormatError） | 重复解析已有 context helper 的字段 |
| Recovery panic 处理 | 在 Middleware 里组装业务 reply |

---

## Middleware 链声明

```go
// ✅ 聚合函数统一暴露 middleware 链
func Server() []middleware.Middleware {
    return []middleware.Middleware{
        Recovery(),
        FormatError(),
        Localize(),
        Tenant(),
        Viewer(),
    }
}

func Client() []middleware.Middleware {
    return []middleware.Middleware{
        Recovery(),
        FormatError(),
        SeataTxClient(),
        metadata.Client(),
        tracing.Client(),
    }
}

// ❌ 每个 server 自己拼装 middleware 顺序
func NewHTTPServer(...) *http.Server {
    return http.NewServer(listener,
        http.Middleware(
            Recovery(), Localize(), Tenant(), Viewer(),  // ❌ 每处都写，顺序可能不一致
        ),
    )
}
```

---

## Middleware 单一职责

```go
// ✅ Localize middleware 只注入 localizer，不做翻译
func Localize() middleware.Middleware {
    return func(handler middleware.Handler) middleware.Handler {
        return func(ctx context.Context, req interface{}) (interface{}, error) {
            ctx = context2.SetLocalize(ctx, i18n.NewLocalizer(
                localize.GetBundle(),
                metadata.GetLanguage(ctx).ToI18nLanguage(),
            ))
            return handler(ctx, req)
        }
    }
}

// ❌ Middleware 中查询 DB 补充用户权限
func Auth() middleware.Middleware {
    return func(handler middleware.Handler) middleware.Handler {
        return func(ctx context.Context, req interface{}) (interface{}, error) {
            viewerID := metadata.GetViewerID(ctx)
            user, err := userRepo.GetUser(ctx, viewerID)  // ❌ middleware 不依赖 Repo
            if err != nil { return nil, err }
            ctx = context2.WithViewer(user)(ctx)
            return handler(ctx, req)
        }
    }
}

// ❌ Middleware 做业务鉴权判断
func Auth() middleware.Middleware {
    return func(handler middleware.Handler) middleware.Handler {
        return func(ctx context.Context, req interface{}) (interface{}, error) {
            if !hasPermission(ctx, req) {  // ❌ 业务鉴权应在 UseCase 中
                return nil, errors.New(403, "FORBIDDEN", "无权限")
            }
            return handler(ctx, req)
        }
    }
}
```

---

## 组合场景

```go
// 完整：Middleware 链注入 context + Depend 透传 metadata

// 1. server middleware 链（按顺序注册）
func Server() []middleware.Middleware {
    return []middleware.Middleware{
        Recovery(),          // panic 兜底，最外层
        FormatError(),       // 错误响应格式化
        Localize(),          // 注入 localizer（依赖 metadata.GetLanguage，需在 Tenant 之前）
        Tenant(),            // 注入 tenant（从 metadata 读 tenantID）
        Viewer(),            // 注入 viewerID
    }
}

// 2. Tenant middleware 实现
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
```

---

## 常见错误模式

```go
// ❌ 每个 server 各自拼 middleware 链（顺序不一致）
// grpc_server.go 和 http_server.go 的 Middleware 列表不同 → 行为不一致

// ❌ middleware 中访问 Repo
func Viewer() middleware.Middleware {
    return func(handler middleware.Handler) middleware.Handler {
        return func(ctx context.Context, req interface{}) (interface{}, error) {
            id := metadata.GetViewerID(ctx)
            user, _ := viewerRepo.GetUser(ctx, id)  // ❌ 不应在 middleware 访问 Repo
            ...
        }
    }
}

// ❌ 在 middleware 中手写 metadata key
md.Get("x-md-global-language")  // ❌ 应用 metadata.GetLanguage(ctx)
```
