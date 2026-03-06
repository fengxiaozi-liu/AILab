# I18n Reference

## 这个主题解决什么问题

说明错误文案和业务文案的 i18n key 如何设计，以及翻译文件通常如何生成。

## 适用场景

- 新增 i18n key
- 调整文案结构
- 维护翻译生成流程

## 设计意图

i18n 参考主要解释文案为什么要围绕稳定 key 组织，而不是围绕原文直接散落在代码里。

- 稳定 key 能把业务语义、错误表达和多语言文案统一起来。
- 理解 key 是语义标识后，更容易避免把界面文案直接嵌入业务流程。
- 文案与 key 分离后，翻译、审校和替换描述都会更轻量。

## 实施提示

- 先确定 key 想表达的业务语义，再写默认文案。
- 先复用已有 key 前缀和模块划分，再新增条目。
- 如果文本会同时出现在错误、提示和 UI 上，优先确认是否共享同一语义 key。

## 推荐结构

- key 按语义命名，不按页面临时文案命名
- 由源语言文件生成其他语言版本

## 典型流程

```text
goi18n extract
-> translate.py (--template)
-> 覆盖目标翻译文件
```

## 脚本调用示例

从源语言文件抽取：

```powershell
goi18n extract -outdir .\assets\i18n -sourceLanguage zh-CN
```

从中文模板翻译英文：

```powershell
python .\tools\translate\translate.py .\assets\i18n\active.zh-CN.toml zh-CN en
```

带已有模板增量翻译：

```powershell
python .\tools\translate\translate.py .\assets\i18n\active.zh-CN.toml zh-CN en --template .\assets\i18n\active.en.toml
```

输出结果：

```text
assets/i18n/active.{fromLanguage}-{toLanguage}.toml
```

## Key 设计示例

```text
ERROR_ORDER_NOT_FOUND
ACCOUNT_OPEN_STATUS_PENDING
```

## 标准模板

```go
func AccountStatusText(ctx context.Context) string {
    localizer := context2.GetLocalize(ctx)
    return localize.LocalizeOrUseOther(localizer, &i18n.LocalizeConfig{
        DefaultMessage: &i18n.Message{
            ID:    "ACCOUNT_OPEN_STATUS_PENDING",
            Other: "待开户",
        },
    })
}
```

## 项目通用错误文案示例

```go
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
```

## Proto 中的多语言字段

```proto
import "base/business/v1/business.proto";

message GetLicenseeReply {
  string logo = 1;
  string code = 2;
  base.business.v1.TransField name = 3;
}
```

## 本地化上下文示例

```go
func Localize() middleware.Middleware {
    return func(handler middleware.Handler) middleware.Handler {
        return func(ctx context.Context, req interface{}) (interface{}, error) {
            ctx = context2.SetLocalize(ctx, i18n.NewLocalizer(bundle, "zh-CN"))
            return handler(ctx, req)
        }
    }
}
```

## 常见坑

- key 直接绑定展示文案，后期难复用
- 多语言文件来源不一致
- 错误信息和业务文案混用同一套命名策略但语义不清

## 相关 Rule

- `../rules/i18n-rule.md`
