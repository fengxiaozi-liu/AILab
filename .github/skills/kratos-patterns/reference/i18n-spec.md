# 国际化规范

## 概述

本项目使用 [go-i18n](https://github.com/nicksnyder/go-i18n) 实现国际化，支持中文简体、中文繁体、英文三种语言。

## 文件结构

```
assets/i18n/
├── active.zh-CN.toml    # 中文简体
├── active.zh-TW.toml    # 中文繁体
└── active.en.toml       # 英文
```

## 支持的语言

```go
package base

type Language string

const (
    LanguageZH Language = "zh-CN"  // 中文简体
    LanguageTC Language = "tc"     // 中文繁体 (映射到 zh-TW)
    LanguageEN Language = "en"     // 英文
)

// ToI18nLanguage 转换为 i18n 语言标识
func (l Language) ToI18nLanguage() string {
    switch l {
    case LanguageZH:
        return "zh-CN"
    case LanguageTC:
        return "zh-TW"
    default:
        return "en"
    }
}
```

## 消息定义规范

### 在代码中定义

```go
package main

import (
    "github.com/nicksnyder/go-i18n/v2/i18n"
    "{module from go.mod}/internal/pkg/localize"
)

// 基本消息
var msg = localize.LocalizeOrUseOther(localizer, &i18n.LocalizeConfig{
    DefaultMessage: &i18n.Message{
        Description: "用户不存在",              // 描述信息（开发用）
        ID:          "ERROR_USER_NOT_FOUND",   // 消息ID（唯一）
        Other:       "用户不存在",              // 默认消息（中文）
    },
})

// 带参数的消息
var msg = localize.LocalizeOrUseOther(localizer, &i18n.LocalizeConfig{
    DefaultMessage: &i18n.Message{
        Description: "订单金额超限",
        ID:          "ERROR_ORDER_AMOUNT_EXCEED",
        Other:       "订单金额超过最大限额 {{.maxAmount}}",
    },
    TemplateData: map[string]interface{}{
        "maxAmount": "10000",
    },
})
```

### TOML 文件格式

```toml
# active.zh-CN.toml
[ERROR_USER_NOT_FOUND]
other = "用户不存在"

[ERROR_ORDER_AMOUNT_EXCEED]
other = "订单金额超过最大限额 {{.MaxAmount}}"

[ENUM_USER_STATUS_ACTIVE]
other = "已激活"

[ENUM_CURRENCY_CNY]
other = "人民币"
```

```toml
# active.en.toml
[ERROR_USER_NOT_FOUND]
other = "User not found"

[ERROR_ORDER_AMOUNT_EXCEED]
other = "Order amount exceeds maximum limit {{.MaxAmount}}"

[ENUM_USER_STATUS_ACTIVE]
other = "Active"

[ENUM_CURRENCY_CNY]
other = "CNY"
```

## 消息 ID 命名规范

| 类型 | 格式                             | 示例 |
|------|--------------------------------|------|
| 错误消息 | `ERROR_{MODULE}_{ERROR_TYPE}`  | `ERROR_USER_NOT_FOUND` |
| 枚举值 | `ENUM_{MODULE}_{TYPE}_{VALUE}` | `ENUM_USER_STATUS_ACTIVE` |

## 使用方式

### 获取 Localizer

```go
package main

import (
    context2 "linksoft.cn/fin/internal/pkg/context"
)

func (srv *UserService) GetUser(ctx context.Context, req *v1.GetUserRequest) (*v1.User, error) {
    // 从 context 获取 localizer
    localizer := context2.GetLocalize(ctx)
    
    // 使用 localizer 进行国际化
    message := localize.LocalizeOrUseOther(localizer, &i18n.LocalizeConfig{
        DefaultMessage: &i18n.Message{
            ID:    "MSG_USER_FOUND",
            Other: "找到用户",
        },
    })
}
```

### 在错误中使用

```go
package user

import (
    "context"
    "github.com/go-kratos/kratos/v2/errors"
    "github.com/nicksnyder/go-i18n/v2/i18n"
    context2 "linksoft.cn/fin/internal/pkg/context"
    "linksoft.cn/fin/internal/pkg/localize"
)

func ErrorUserNotFound(ctx context.Context) *errors.Error {
    localizer := context2.GetLocalize(ctx)
    return errors.New(404, "USER_NOT_FOUND", localize.LocalizeOrUseOther(localizer, &i18n.LocalizeConfig{
        DefaultMessage: &i18n.Message{
            Description: "用户不存在",
            ID:          "ERROR_USER_NOT_FOUND",
            Other:       "用户不存在",
        },
    }))
}
```

### 在枚举中使用

```go
package user

func (s UserStatus) Localize(localizer *i18n.Localizer) string {
    switch s {
    case UserStatusActive:
        return localize.LocalizeOrUseOther(localizer, &i18n.LocalizeConfig{
            DefaultMessage: &i18n.Message{
                Description: "用户状态-已激活",
                ID:          "ENUM_USER_STATUS_ACTIVE",
                Other:       "已激活",
            },
        })
    // ... 其他状态
    default:
        return string(s)
    }
}
```

## 提取国际化文件

从代码中提取所有国际化消息到 TOML 文件：

```shell
goi18n extract -outdir .\assets\i18n -sourceLanguage zh-CN
```

## 合并翻译文件

当有新的消息需要翻译时：

```shell
# 切换到翻译工具目录
cd tools/translate
# 使用翻译脚本生成 translate.*.toml 文件
python .\translate.py ..\..\assets\i18n\active.zh-CN.toml zh-CN {目标语言zh-TW|en} --template ..\..\assets\i18n\active.{目标语言zh-TW|en}.toml
# 将生成的文件重命名为正确的文件名
mv -f ..\..\assets\i18n\active.zh-CN_{目标语言zh-TW|en}.toml ..\..\assets\i18n\active.{目标语言zh-TW|en}.toml
```

## Context 传递语言

语言信息通过 gRPC metadata 传递，在中间件中设置 Localizer：

```go
// internal/pkg/middleware/localize.go
func Localize() middleware.Middleware {
    return func(handler middleware.Handler) middleware.Handler {
        return func(ctx context.Context, req interface{}) (interface{}, error) {
            // 从 metadata 获取语言设置
            lang := getLanguageFromMetadata(ctx)
            
            // 创建 localizer 并放入 context
            localizer := i18n.NewLocalizer(bundle, lang)
            ctx = context.WithValue(ctx, localizerKey{}, localizer)
            
            return handler(ctx, req)
        }
    }
}
```

## 最佳实践

1. **统一消息 ID**：使用规范的消息 ID 命名
2. **默认消息**：始终提供 `Other` 作为默认消息（中文）
3. **参数化消息**：使用模板语法 `{{.Param}}` 处理动态内容
4. **定期提取**：定期运行 `goi18n extract` 更新翻译文件
5. **完整翻译**：确保所有消息在所有语言中都有翻译
6. **上下文传递**：通过 context 传递 localizer，避免全局状态
