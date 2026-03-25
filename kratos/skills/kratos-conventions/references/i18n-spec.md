# I18n Spec

## key 命名

| 条件 | 做法 |
|------|------|
| 稳定业务语义 | 用稳定 key |
| 临时页面文案 | 不直接上升为全局 key |

```toml
# ✅
[ERROR_ORDER_NOT_FOUND]
description = "订单不存在"
other = "订单不存在"
```

```toml
# ✅ 带模板参数
[ERROR_ORDER_AMOUNT_EXCEED]
description = "订单金额超限"
other = "订单金额超过最大限制 {{.maxAmount}}"
```

```toml
# ❌ 页面型临时 key
[ORDER_PAGE_POPUP_TEXT]
other = "订单不存在"
```

---

## 代码中引用

| 条件 | 做法 |
|------|------|
| 需要本地化文案 | 通过 localize 配置引用 key |
| 需要默认文案 | 放在 `DefaultMessage`，不散落在业务流程里 |

```go
// ✅
msg := localize.LocalizeOrUseOther(localizer, &i18n.LocalizeConfig{
    DefaultMessage: &i18n.Message{
        ID:          "ERROR_ORDER_NOT_FOUND",
        Description: "订单不存在",
        Other:       "订单不存在",
    },
})
```

```go
// ❌ 中文文案散落
return nil, errors.New("订单不存在")
```

---

## 生成流程

| 条件 | 做法 |
|------|------|
| 修改源语言文件 | 走生成流程同步产物 |
| 只改派生产物 | 不允许 |

```text
// ✅
goi18n extract
-> 更新源语言文件
-> 翻译脚本补齐目标语言
```

```text
// ❌
直接手改 en.toml / zh-CN.toml 派生产物
```

---

## 组合场景

```go
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

这个组合场景同时满足：

- i18n key 表达稳定语义
- 默认文案不散落在业务流程
- 错误语义和多语言 key 对齐

---

## 常见错误模式

```go
// ❌ 中文文案写死在业务逻辑
```

```toml
# ❌ 每个页面造一套 key
```

```text
// ❌ 手改翻译产物
```
