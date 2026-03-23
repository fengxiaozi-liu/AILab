# Error Spec

## 错误函数模板

```go
// ✅ 函数名、Reason、i18n ID 统一，使用 localizer
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

// ❌ 直接用 errors.New，无 Reason，无 i18n，跨边界后无法定位
return errors.New("account not found")

// ❌ 直接用 fmt.Errorf 外抛到协议层
return nil, fmt.Errorf("account %d not found", id)
```

---

## 命名规则

| 元素 | 格式 | 示例 |
|------|------|------|
| 函数名 | `Error{Module}{Description}` | `ErrorAccountNotFound` |
| Reason | `{MODULE}_{DESCRIPTION}` | `OPEN_ACCOUNT_NOT_FOUND` |
| i18n ID | `ERROR_{MODULE}_{DESCRIPTION}` | `ERROR_OPEN_ACCOUNT_NOT_FOUND` |

---

## HTTP code 选择

| 场景 | code |
|------|------|
| 资源不存在 | 404 |
| 参数错误、状态不允许、前置条件不满足 | 400 |
| 频控、重复请求 | 400 / 429 |
| 内部处理失败、依赖调用失败 | 500 |

---

## not found 处理

```go
// ✅ Data 层判断 ent.IsNotFound → 转换为业务错误
func (r *accountRepo) GetAccount(ctx context.Context, id uint32) (*biz.Account, error) {
    info, err := query.First(ctx)
    if err != nil {
        if ent.IsNotFound(err) {
            return nil, openerror.ErrorAccountNotFound(ctx)
        }
        return nil, err
    }
    return accountConvert(info), nil
}

// ❌ not found 返回 nil, nil，调用方需要猜测状态
func (r *accountRepo) GetAccount(ctx context.Context, id uint32) (*biz.Account, error) {
    info, err := query.First(ctx)
    if ent.IsNotFound(err) {
        return nil, nil  // 调用方需要 if obj != nil 判断
    }
    return accountConvert(info), err
}
```

---

## 带模板参数的错误

```go
// ✅
func ErrorOrderAmountExceed(ctx context.Context, templateData map[string]interface{}) *errors.Error {
    localizer := context2.GetLocalize(ctx)
    return errors.New(400, "ORDER_AMOUNT_EXCEED", localize.LocalizeOrUseOther(localizer, &i18n.LocalizeConfig{
        DefaultMessage: &i18n.Message{
            ID:    "ERROR_ORDER_AMOUNT_EXCEED",
            Other: "订单金额超过最大限额 {{.maxAmount}}",
        },
        TemplateData: templateData,
    }))
}
```

---

## 组合场景

```go
// Data 层转换 not found + Biz 层检查业务状态
func (r *accountRepo) GetAccount(ctx context.Context, id uint32) (*biz.Account, error) {
    info, err := r.data.Db.Account(ctx).Query().Where(entaccount.IDEQ(id)).First(ctx)
    if err != nil {
        if ent.IsNotFound(err) {
            return nil, openerror.ErrorAccountNotFound(ctx)
        }
        return nil, err
    }
    return accountConvert(info), nil
}

func (u *AccountUseCase) Review(ctx context.Context, id uint32) error {
    account, err := u.accountRepo.GetAccount(ctx, id)
    if err != nil {
        return err  // ✅ 直接返回，不吞错
    }
    if !account.CanReview() {
        return openerror.ErrorReviewStatusInvalid(ctx)  // ✅ 业务错误语义清晰
    }
    return nil
}
```

---

## 常见错误模式

```go
// ❌ 必须存在的查询返回 nil, nil
if ent.IsNotFound(err) { return nil, nil }

// ❌ 吞错
_, _ = u.repo.UpdateAccount(ctx, account)

// ❌ 跨边界外抛原生错误
return nil, fmt.Errorf("db error: %w", err)

// ❌ 方法名/Reason/i18n ID 不统一
return errors.New(404, "NOT_FOUND_ACCOUNT", "账号不存在")
// 函数名叫 ErrorAccountMissing，Reason 叫 NOT_FOUND_ACCOUNT
```
