# Logging Spec

## logger 与 helper

| 条件 | 做法 |
|------|------|
| 模块初始化 | 统一用 `log.Logger` / `log.Helper` |
| 需要模块标识 | 用 `log.With(logger, "module", "...")` |

```go
// ✅
type AccountUseCase struct {
    log *log.Helper
}

func NewAccountUseCase(logger log.Logger) *AccountUseCase {
    return &AccountUseCase{
        log: log.NewHelper(log.With(logger, "module", "account/usecase")),
    }
}
```

```go
// ❌ 直接散落底层 logger
zap.L().Error("load account failed")
```

---

## 上下文与错误日志

| 条件 | 做法 |
|------|------|
| 需要请求上下文 | 用 `WithContext(ctx)` |
| 记录错误 | 带动作和 `err`，不要只打一行模糊信息 |

```go
// ✅
if err != nil {
    u.log.WithContext(ctx).Errorf("load account failed: %v", err)
    return nil, err
}
```

```go
// ❌ 无上下文无动作
log.Error(err)
```

---

## 日志量与脱敏

| 条件 | 做法 |
|------|------|
| 高频循环、批量处理 | 控制日志量级，必要时只记录摘要 |
| 包含敏感信息 | 脱敏或不记录 |

```go
// ✅ 只记录批量摘要
l.log.Infof("sync accounts finished, count=%d", len(list))
```

```go
// ❌ 循环内逐条打印敏感内容
for _, item := range list {
    l.log.Infof("account=%+v", item)
}
```

---

## 组合场景

```go
type Consumer struct {
    log *log.Helper
}

func NewConsumer(logger log.Logger) *Consumer {
    return &Consumer{log: log.NewHelper(log.With(logger, "module", "account/consumer"))}
}

func (c *Consumer) Handle(ctx context.Context, msg *Message) error {
    c.log.WithContext(ctx).Infof("consume message, id=%s", msg.ID)
    if err := c.process(ctx, msg); err != nil {
        c.log.WithContext(ctx).Errorf("consume message failed, id=%s err=%v", msg.ID, err)
        return err
    }
    return nil
}
```

这个组合场景同时满足：

- 统一 helper
- module 字段稳定
- 带上下文日志
- 错误可搜索

---

## 常见错误模式

```go
// ❌ 直接用底层 logger
```

```go
// ❌ 高频循环日志风暴
```

```go
// ❌ 记录敏感信息
```
