# I18n Spec

## Key 设计

| 条件 | 做法 |
|------|------|
| 错误文案 | `ERROR_{MODULE}_{DESCRIPTION}` |
| 枚举展示文案 | `ENUM_{TYPE}_{VALUE}` |
| 业务提示文案 | 按语义命名，不按页面路径命名 |

```go
// ✅ Key 按语义稳定命名
&i18n.Message{ID: "ERROR_OPEN_ACCOUNT_NOT_FOUND", Other: "开户申请不存在"}
&i18n.Message{ID: "ENUM_OPEN_STATUS_INIT", Other: "初始化"}
&i18n.Message{ID: "ACCOUNT_OPEN_STATUS_PENDING", Other: "待开户"}

// ❌ Key 按页面路径或临时描述命名，不稳定
&i18n.Message{ID: "page_open_form_error_1", Other: "不存在"}
&i18n.Message{ID: "status_text_3", Other: "待审核"}
```

---

## 翻译生成流程

```text
// ✅ 正确流程
goi18n extract -outdir ./assets/i18n -sourceLanguage zh-CN
-> python tools/translate/translate.py ./assets/i18n/active.zh-CN.toml zh-CN en
-> 产出 assets/i18n/active.zh-CN-en.toml
-> 合并/重命名为 active.en.toml

// ❌ 直接手改 active.*.toml 产物
// 下次 extract 后改动会被覆盖
```

---

## 标准用法

```go
// ✅ 在错误函数中使用
func ErrorPageNotFound(ctx context.Context) *errors.Error {
    localizer := context2.GetLocalize(ctx)
    return errors.New(400, "OPEN_PAGE_NOT_FOUND", localize.LocalizeOrUseOther(localizer, &i18n.LocalizeConfig{
        DefaultMessage: &i18n.Message{
            Description: "页面不存在",
            ID:          "ERROR_OPEN_PAGE_NOT_FOUND",
            Other:       "页面不存在",
        },
    }))
}

// ✅ 在枚举 Localize 方法中使用
func (s AccountOpenStatus) Localize(localizer *i18n.Localizer) string {
    return localize.LocalizeOrUseOther(localizer, &i18n.LocalizeConfig{
        DefaultMessage: &i18n.Message{ID: "ENUM_OPEN_STATUS_FILLING", Other: "填写中"},
    })
}
```

---

## Proto 多语言字段

```proto
// ✅ 多语言字段使用 TransField
import "base/business/v1/business.proto";
message GetLicenseeReply {
  base.business.v1.TransField name = 1;
}

// ❌ 每种语言单独一个字段
message GetLicenseeReply {
  string name_zh = 1;
  string name_en = 2;
}
```

---

## 组合场景

```go
// 完整：key 定义 → 错误函数 → Localize 中间件 → 使用
// 1. 定义错误函数，key=ERROR_OPEN_ACCOUNT_NOT_FOUND
func ErrorAccountNotFound(ctx context.Context) *errors.Error { ... }

// 2. middleware 注入 localizer 到 ctx
func Localize() middleware.Middleware {
    return func(handler middleware.Handler) middleware.Handler {
        return func(ctx context.Context, req interface{}) (interface{}, error) {
            ctx = context2.SetLocalize(ctx, i18n.NewLocalizer(bundle, metadata.GetLanguage(ctx).ToI18nLanguage()))
            return handler(ctx, req)
        }
    }
}

// 3. 在 data 层触发转换
if ent.IsNotFound(err) {
    return nil, openerror.ErrorAccountNotFound(ctx)
}
```

---

## 常见错误模式

```go
// ❌ 原始文案散落业务代码
return errors.New(404, "NOT_FOUND", "账号不存在")

// ❌ 手改 active.en.toml
// extract 后变动丢失

// ❌ Key 不稳定，带版本号
&i18n.Message{ID: "error_account_v2_not_found"}
```
