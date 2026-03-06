# Middleware 参考

## 这个主题解决什么问题

统一请求链路中的基础拦截、上下文补充、错误格式化和事务透传，避免在业务层重复处理这些通用逻辑。

## 适用场景

- 新增 client/server middleware
- 调整 localize、viewer、tenant、format error、recovery、seata 事务相关中间件

## 推荐结构或实现方式

- middleware 目录按主题拆分，例如 `localize.go`、`tenant.go`、`viewer.go`、`format_error.go`。
- 再通过聚合函数统一暴露中间件链。

## 标准模板

```go
func Client() []middleware.Middleware {
    return []middleware.Middleware{
        Recovery(),
        FormatError(),
        SeataTxClient(),
        metadata.Client(),
        tracing.Client(),
    }
}
```

```go
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
```

## Good Example

- localize middleware 只负责读取语言并构造 localizer，不做业务字段翻译。
- client middleware 统一挂载 tracing、metadata、recovery，调用方不逐处拼链。

## 常见坑

- 在 middleware 中直接查询 repo 或做业务判断
- 不复用 metadata/context helper，重复解析相同字段

## 相关 rule / 相关 reference

- `../rules/middleware-rule.md`
- `./context-spec.md`
- `./metadata-spec.md`
