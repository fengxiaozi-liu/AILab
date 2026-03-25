# Error Spec

## 错误语义收敛

| 条件 | 做法 |
|------|------|
| 已有稳定业务语义 | 用统一错误函数、reason 和 i18n ID 表达 |
| 只是底层技术异常 | 不强行包装成新的业务错误 |

```go
// ✅ 统一错误函数
func ErrorOrderNotFound(ctx context.Context) *errors.Error {
    localizer := context2.GetLocalize(ctx)
    return errors.New(500, "ORDER_NOT_FOUND", localize.LocalizeOrUseOther(localizer, &i18n.LocalizeConfig{
        DefaultMessage: &i18n.Message{
            ID:          "ERROR_ORDER_NOT_FOUND",
            Description: "订单不存在",
            Other:       "订单不存在",
        },
    }))
}
```

```go
// ❌ 业务错误散落
return nil, fmt.Errorf("order not found")
```

---

## not found 处理

| 条件 | 做法 |
|------|------|
| ent/db not found | 转成统一业务 not found 错误 |
| 上层要判断 not found | 判断统一错误语义，不伪装成功 |

```go
// ✅
if ent.IsNotFound(err) {
    return nil, openerror.ErrorAccountNotFound(ctx)
}
```

```go
// ❌ 伪装成功
if ent.IsNotFound(err) {
    return nil, nil
}
```

```go
// ⚠️ 上层需要区分 not found，也应保留错误语义
if err != nil && stderrors.Is(err, openerror.ErrorAccountNotFound(ctx)) {
    return nil, err
}
```

---

## 命名模式

| 元素 | 做法 |
|------|------|
| 函数名 | `Error{Module}{Description}` |
| reason | `{MODULE}_{DESCRIPTION}` |
| i18n ID | `ERROR_{MODULE}_{DESCRIPTION}` |

```go
// ✅
ErrorOrderNotFound
ORDER_NOT_FOUND
ERROR_ORDER_NOT_FOUND
```

```go
// ❌ 模式漂移
NewOrder404
ORDER_NOT_EXIST
ORDER_NOT_FOUND_MSG
```

---

## 组合场景

```go
func (r *orderRepo) GetOrder(ctx context.Context, id uint32) (*biz.Order, error) {
    info, err := r.data.Db.Order(ctx).Query().Where(order.IDEQ(id)).Only(ctx)
    if err != nil {
        if ent.IsNotFound(err) {
            return nil, ordererror.ErrorOrderNotFound(ctx)
        }
        return nil, err
    }
    return orderConvert(info), nil
}
```

这个组合场景同时满足：

- Data/Biz 共用稳定错误语义
- not found 没被伪装成功
- 函数名、reason、i18n ID 同步

---

## 常见错误模式

```go
// ❌ nil, nil
return nil, nil
```

```go
// ❌ 裸错误字符串
errors.New("account not found")
```

```go
// ❌ 一个模块多套 reason 命名
ORDER_NOT_FOUND / ORDER_NOT_EXIST
```
