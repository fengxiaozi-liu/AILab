# Metadata 参考

## 这个主题解决什么问题

统一 Kratos metadata 的读写、来源优先级和字段命名，保证语言、租户、viewer、client IP 等信息能稳定跨边界透传。

## 适用场景

- 新增 metadata 字段
- 需要在 client/server context 间读取或追加透传字段
- 需要和 middleware、context 联动维护请求级信息

## 推荐结构或实现方式

- 每个字段提供 `GetXxx/SetXxx` 成对 helper。
- 使用统一 `Config` 控制从 client/server 读取 metadata 的来源。

## 标准模板

```go
func GetLanguage(ctx context.Context, opts ...Option) baseenum.Language {
    language := baseenum.LanguageDefault()
    if md, ok := GetMetadata(ctx, opts...); ok {
        if mdLanguage := md.Get("x-md-global-language"); mdLanguage != "" {
            language = baseenum.Language(mdLanguage)
        }
    }
    return language
}

func SetLanguage(ctx context.Context, language baseenum.Language) context.Context {
    return metadata.AppendToClientContext(ctx, "x-md-global-language", string(language))
}
```

## Good Example

- `GetMetadata` 先合并 client metadata，再用 server metadata 覆盖，收口来源优先级。
- `SetTenantID`、`SetViewerID`、`SetPlatform` 都通过统一 helper 暴露给调用方。

## 常见坑

- 在 service 或 middleware 里手写 `x-md-*` key
- 同一语义字段存在多个别名，导致透传不稳定

## 相关 rule / 相关 reference

- `../rules/metadata-rule.md`
- `./context-spec.md`
- `./middleware-spec.md`
