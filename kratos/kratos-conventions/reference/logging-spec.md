# Logging Reference

## 这个主题解决什么问题

说明 Kratos 项目中的日志对象如何创建、如何在不同层传递，以及错误日志、上下文日志、资源生命周期日志通常如何书写。

## 适用场景

- 初始化全局 logger
- 在 UseCase、Repo、Listener、Crontab 中记录日志
- 记录依赖调用、资源关闭、预加载、失败重试等过程日志
- 需要把请求上下文带入日志

## 设计意图

日志承担两类职责：

- 为运行问题提供可搜索、可定位的上下文
- 为关键流程节点提供最小必要的观测信息

稳定做法不是散落地直接调用底层 zap，而是统一使用 Kratos `log.Logger` / `log.Helper`，需要上下文时再通过 `WithContext(ctx)` 附着上下文信息。

## 实施提示

- 组件初始化阶段优先创建 `log.NewHelper(logger)` 或 `log.NewHelper(log.With(...))`
- 需要模块标识时，用 `log.With(logger, "module", "...")` 先补充固定字段
- 需要请求上下文时，用 `helper.WithContext(ctx)` 输出日志
- 流程节点日志写清动作和关键业务标识，错误日志补 `err`
- 高频循环、批量处理和调试 SQL 日志要控制量级，避免日志风暴

## 推荐实现方式

### 1. 全局 Logger 初始化

```go
zLogger := kratoszap.NewLogger(zap.New(core))
logger := log.With(
    zLogger,
    "ts", log.DefaultTimestamp,
    "caller", log.DefaultCaller,
    "service.name", Name,
    "service.version", Version,
)

log.SetLogger(logger)
```

### 2. 模块级 Helper

```go
type AccountListener struct {
    log *log.Helper
}

func NewAccountListener(logger log.Logger) *AccountListener {
    return &AccountListener{
        log: log.NewHelper(log.With(logger, "module", "listener/account")),
    }
}
```

### 3. 上下文日志

```go
r.log.WithContext(ctx).Info("schema preload: start")
r.log.WithContext(ctx).Warn("schema preload: no base_url configured, skipping")
r.log.WithContext(ctx).Errorf("schema preload: download manifest %s failed: %v", manifestURL, err)
```

### 4. 生命周期日志

```go
helper := log.NewHelper(logger)
helper.Info("closing the data resources")
helper.Info("data resources closed")
```

## 代码示例参考

### Repo / 组件日志

```go
type SchemaRepo struct {
    log *log.Helper
}

func NewSchemaRepo(logger log.Logger) *SchemaRepo {
    return &SchemaRepo{
        log: log.NewHelper(logger),
    }
}

func (r *SchemaRepo) PreloadAll(ctx context.Context) error {
    r.log.WithContext(ctx).Info("schema preload: start")

    if err := r.downloadManifest(ctx); err != nil {
        r.log.WithContext(ctx).Errorf("schema preload: download manifest failed: %v", err)
        return err
    }

    r.log.WithContext(ctx).Info("schema preload: done")
    return nil
}
```

### Listener 日志

```go
func (l *AccountListener) Handle(ctx context.Context, payload []byte) {
    if err := json.Unmarshal(payload, &body); err != nil {
        l.log.Errorf("unmarshal account event failed: %v", err)
        return
    }
}
```

### SQL 调试日志

```go
sqlDrv := dialect.DebugWithContext(drv, func(ctx context.Context, i ...interface{}) {
    helper.WithContext(ctx).Info(i...)
})
```

## Good Example

- 全局 logger 统一初始化，局部只注入 `log.Logger`
- 模块对象持有 `*log.Helper`
- 错误日志写清动作、目标对象和错误原因
- 上下文相关日志通过 `WithContext(ctx)` 统一透传

## 常见坑

- 直接在业务代码里散落构造 zap logger
- 同一模块同时混用 `log.Logger`、`*log.Helper`、全局 `log`，风格不统一
- 错误日志只输出 `err`，没有动作和业务标识
- 把完整请求体、密钥、证件号直接写入日志

## 相关 Rule

- `../rules/logging-rule.md`

## 相关 Reference

- `./error-spec.md`
