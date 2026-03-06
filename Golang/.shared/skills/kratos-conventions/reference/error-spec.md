# Error Reference

## 这个主题解决什么问题

说明 Kratos 项目中统一错误函数如何定义、命名、组织，以及 Data/Biz 层如何使用这些错误函数。

## 适用场景

- 新增业务错误
- 调整错误函数模板
- 统一 HTTP 状态码、Reason 和 i18n 消息

## 设计意图

错误参考不只是说明如何 new 一个错误，而是解释错误体系如何承接业务语义、对外表达和多语言文案。

- 错误函数名、Reason 和 i18n ID 一致时，定位问题和检索日志会更直接。
- Data/Biz 层共享同一套错误表达后，协议层不需要再临时猜测返回语义。
- 理解错误是业务语义对象后，更不容易把内部异常直接暴露给外层。

## 实施提示

- 先确定业务语义，再映射成错误函数、Reason 和文案 ID。
- 先看当前模块已有命名，再补同风格的新错误。
- 如果错误需要展示给用户，顺手检查是否需要与 i18n key 联动。

## 目录结构

```text
internal/error/{module}/{module}.go
internal/error/base/base.go
```

## 错误函数模板

```go
func ErrorOrderNotFound(ctx context.Context) *errors.Error {
    localizer := context2.GetLocalize(ctx)
    return errors.New(404, "ORDER_NOT_FOUND", localize.LocalizeOrUseOther(localizer, &i18n.LocalizeConfig{
        DefaultMessage: &i18n.Message{
            Description: "订单不存在",
            ID:          "ERROR_ORDER_NOT_FOUND",
            Other:       "订单不存在",
        },
    }))
}
```

## 项目通用错误示例

```go
func ErrorAccountNotFound(ctx context.Context) *errors.Error {
    localizer := context2.GetLocalize(ctx)
    return errors.New(404, "OPEN_ACCOUNT_NOT_FOUND", localize.LocalizeOrUseOther(localizer, &i18n.LocalizeConfig{
        DefaultMessage: &i18n.Message{
            Description: "开户申请不存在",
            ID:          "ERROR_OPEN_ACCOUNT_NOT_FOUND",
            Other:       "开户申请不存在",
        },
    }))
}

func ErrorDuplicateRequest(ctx context.Context) *errors.Error {
    localizer := context2.GetLocalize(ctx)
    return errors.New(429, "OPEN_DUPLICATE_REQUEST", localize.LocalizeOrUseOther(localizer, &i18n.LocalizeConfig{
        DefaultMessage: &i18n.Message{
            Description: "重复请求",
            ID:          "ERROR_OPEN_DUPLICATE_REQUEST",
            Other:       "请求过于频繁，请稍后重试",
        },
    }))
}
```

## 带模板参数

```go
func ErrorOrderAmountExceed(ctx context.Context, templateData map[string]interface{}) *errors.Error {
    localizer := context2.GetLocalize(ctx)
    return errors.New(400, "ORDER_AMOUNT_EXCEED", localize.LocalizeOrUseOther(localizer, &i18n.LocalizeConfig{
        DefaultMessage: &i18n.Message{
            Description: "订单金额超限",
            ID:          "ERROR_ORDER_AMOUNT_EXCEED",
            Other:       "订单金额超过最大限额 {{.maxAmount}}",
        },
        TemplateData: templateData,
    }))
}
```

## 命名模式

| 元素 | 格式 | 示例 |
|------|------|------|
| 函数名 | `Error{Module}{Description}` | `ErrorOrderNotFound` |
| Reason | `{MODULE}_{DESCRIPTION}` | `ORDER_NOT_FOUND` |
| i18n ID | `ERROR_{MODULE}_{DESCRIPTION}` | `ERROR_ORDER_NOT_FOUND` |

## Code 与 Reason 的分工

- `code` 用于表达错误大类，通常使用通用 HTTP/Kratos 错误码分层。
- `reason` 用于表达稳定的业务错误语义，作为主要区分标识。
- 业务侧优先通过 `reason` 判断具体错误，而不是继续细分大量数字错误码。

常见分层：

- `400`：参数错误、状态不允许、业务前置条件不满足
- `404`：资源不存在
- `500`：内部处理失败、依赖调用失败、系统异常

示例：

```go
func ErrorDuplicateRequest(ctx context.Context) *errors.Error {
    localizer := context2.GetLocalize(ctx)
    return errors.New(400, "OPEN_DUPLICATE_REQUEST", localize.LocalizeOrUseOther(localizer, &i18n.LocalizeConfig{
        DefaultMessage: &i18n.Message{
            Description: "重复请求",
            ID:          "ERROR_OPEN_DUPLICATE_REQUEST",
            Other:       "请求过于频繁，请稍后重试",
        },
    }))
}
```

这里的重点是：

- `400` 只表示“频控/重复请求”这一大类
- `OPEN_DUPLICATE_REQUEST` 才是稳定的业务区分标识
- `ERROR_OPEN_DUPLICATE_REQUEST` 对应 i18n 文案键

## 典型使用方式

### Data 层

```go
if ent.IsNotFound(err) {
    return nil, ordererror.ErrorOrderNotFound(ctx)
}
```

### Biz 层

```go
if exists {
    return nil, ordererror.ErrorOrderAlreadyExists(ctx)
}
```

### UseCase 中的 not found 语义收敛

```go
account, err := u.accountRepo.GetAccountByUserCode(ctx, userCode)
if err != nil {
    if stderrors.Is(err, openerror.ErrorAccountNotFound(ctx)) {
        return openenum.AccountOpenStatusInit, nil
    }
    return 0, err
}
```

## 常见坑

- 一个业务模块里同时混用多套 Reason 命名
- 错误函数模板和 i18n key 对不上
- 错误描述直接写死在业务逻辑里

## 相关 Rule

- `../rules/error-rule.md`
- `../rules/logging-rule.md`
