# 错误处理规范

## 概述

本项目使用 Kratos 框架的错误处理机制，结合国际化实现多语言错误消息。

## 目录结构

```
internal/error/
├── base/           # 基础错误
│   └── base.go
├── admin/          # 管理模块错误
├── system/         # 系统模块错误
└── user/           # 用户模块错误
```

## 错误定义规范

### 基础错误示例

```go
package base

import (
    "context"
    "github.com/go-kratos/kratos/v2/errors"
    "github.com/nicksnyder/go-i18n/v2/i18n"
    context2 "{module from go.mod}/internal/pkg/context"
    "{module from go.mod}/internal/pkg/localize"
)

// ErrorFailed 系统异常
func ErrorFailed(ctx context.Context) *errors.Error {
    localizer := context2.GetLocalize(ctx)
    return errors.New(500, "BUSINESS_FAILED", localize.LocalizeOrUseOther(localizer, &i18n.LocalizeConfig{
        DefaultMessage: &i18n.Message{
            Description: "系统异常",
            ID:          "ERROR_BUSINESS_FAILED",
            Other:       "系统异常",
        },
    }))
}

// ErrorBadRequest 客户端错误
func ErrorBadRequest(ctx context.Context) *errors.Error {
    localizer := context2.GetLocalize(ctx)
    return errors.New(400, "BUSINESS_BAD_REQUEST", localize.LocalizeOrUseOther(localizer, &i18n.LocalizeConfig{
        DefaultMessage: &i18n.Message{
            Description: "请求信息错误",
            ID:          "ERROR_BUSINESS_BAD_REQUEST",
            Other:       "请求信息错误",
        },
    }))
}

// ErrorNeedLogin 需要登录
func ErrorNeedLogin(ctx context.Context) *errors.Error {
    localizer := context2.GetLocalize(ctx)
    return errors.New(401, "BUSINESS_NEED_LOGIN", localize.LocalizeOrUseOther(localizer, &i18n.LocalizeConfig{
        DefaultMessage: &i18n.Message{
            Description: "需要登录",
            ID:          "ERROR_BUSINESS_NEED_LOGIN",
            Other:       "需要登录",
        },
    }))
}

// ErrorPermissionDenied 权限不足
func ErrorPermissionDenied(ctx context.Context) *errors.Error {
    localizer := context2.GetLocalize(ctx)
    return errors.New(403, "BUSINESS_PERMISSION_DENIED", localize.LocalizeOrUseOther(localizer, &i18n.LocalizeConfig{
        DefaultMessage: &i18n.Message{
            Description: "权限不足",
            ID:          "ERROR_BUSINESS_PERMISSION_DENIED",
            Other:       "权限不足",
        },
    }))
}
```

### 业务错误示例

```go
package user

import (
    "context"
    "github.com/go-kratos/kratos/v2/errors"
    "github.com/nicksnyder/go-i18n/v2/i18n"
    context2 "{module from go.mod}/internal/pkg/context"
    "{module from go.mod}/internal/pkg/localize"
)

// ErrorUserNotFound 用户不存在
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

// ErrorUserAlreadyExists 用户已存在
func ErrorUserAlreadyExists(ctx context.Context) *errors.Error {
    localizer := context2.GetLocalize(ctx)
    return errors.New(400, "USER_ALREADY_EXISTS", localize.LocalizeOrUseOther(localizer, &i18n.LocalizeConfig{
        DefaultMessage: &i18n.Message{
            Description: "用户已存在",
            ID:          "ERROR_USER_ALREADY_EXISTS",
            Other:       "用户已存在",
        },
    }))
}
```

## HTTP 状态码规范

| HTTP 状态码 | 说明 | 使用场景 |
|------------|------|---------|
| 400 | Bad Request | 请求参数错误、业务校验失败 |
| 401 | Unauthorized | 未登录、Token 过期 |
| 403 | Forbidden | 权限不足 |
| 404 | Not Found | 资源不存在 |
| 500 | Internal Server Error | 服务器内部错误 |

## 错误码命名规范

格式：`{MODULE}_{ERROR_TYPE}`

| 模块 | 前缀 | 示例 |
|------|-----|------|
| 基础 | `BUSINESS_` | `BUSINESS_FAILED` |
| 用户 | `USER_` | `USER_NOT_FOUND` |
| 订单 | `ORDER_` | `ORDER_NOT_FOUND` |
| 系统 | `SYSTEM_` | `SYSTEM_CONFIG_ERROR` |

## 命名规范（函数 / Reason / i18n ID）

> 目标：见名知意 + 一致性。错误定义的命名需要稳定，方便检索与国际化维护。

### 1) 错误构造函数命名：`ErrorXxx`

Good:
```go
func ErrorBadRequest(ctx context.Context) *errors.Error { return nil }
func ErrorFailed(ctx context.Context) *errors.Error { return nil }
func ErrorUserNotFound(ctx context.Context) *errors.Error { return nil }
```

Bad:
```go
func BadRequest(ctx context.Context) *errors.Error { return nil } // BAD: 缺少 Error 前缀
func ErrBadRequest(ctx context.Context) *errors.Error { return nil } // BAD: 项目错误构造函数不使用 ErrXxx
```

### 2) Reason 命名：全大写下划线 + 模块前缀

Good:
```text
BUSINESS_FAILED
BUSINESS_BAD_REQUEST
USER_NOT_FOUND
USER_ALREADY_EXISTS
```

Bad:
```text
badRequest
BusinessBadRequest
userNotFound
```

### 3) 国际化消息 ID 命名：`ERROR_{MODULE}_{ERROR_TYPE}`

Good:
```text
ERROR_BUSINESS_FAILED
ERROR_BUSINESS_BAD_REQUEST
ERROR_USER_NOT_FOUND
```

Bad:
```text
BusinessFailed
error_user_not_found
```

## 国际化消息 ID 规范

格式：`ERROR_{MODULE}_{ERROR_TYPE}`

```go
package main

import "github.com/nicksnyder/go-i18n/v2/i18n"

var msg = &i18n.Message{
    Description: "描述信息",           // 开发者可读描述
    ID:          "ERROR_USER_NOT_FOUND", // 国际化 ID
    Other:       "用户不存在",          // 默认中文消息
}
```

## 错误使用方式

### 在 data 层使用

```go
package user

import (
	"context"
	"{module from go.mod}/internal/data/ent"
	"{module from go.mod}/internal/data/kit"
	"{module from go.mod}/internal/pkg/filter"
	userbiz "{module from go.mod}/internal/biz/user"
	usererror "{module from go.mod}/internal/error/user"
)

type userRepo struct {
	data *kit.Data
}

func (r *userRepo) GetUser(ctx context.Context, ID uint32, opts ...filter.Option) (*userbiz.User, error) {
	query := r.data.Db.User(ctx).Query().Where(entorder.IDEQ(ID))

	query = r.queryConfig(query, opts...)

	info, err := query.First(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, usererror.ErrorUserNotFound(ctx)
		}
		return nil, err
	}

	res := r.queryRelation(userConvert(info), info.Edges)

	err = r.serviceRelation(ctx, res, opts...)
	if err != nil {
		return nil, err
	}

	return res, nil
}

```

### 在 biz 层使用

```go
package user

import (
    "context"
	usererror "{module from go.mod}/internal/error/user"
)

type User struct {
    Email string
    // 其他字段...
}

type UserRepo interface {
    ExistsByEmail(ctx context.Context, email string) (bool, error)
    Create(ctx context.Context, user *User) (*User, error)
}

type UserUseCase struct{
    userRepo UserRepo
}

func (u *UserUseCase) Create(ctx context.Context, user *User) (*User, error) {
    // 业务校验
    exists, _ := u.userRepo.ExistsByEmail(ctx, user.Email)
    if exists {
        return nil, usererror.ErrorUserAlreadyExists(ctx)
    }
    
    return u.userRepo.Create(ctx, user)
}
```

## 带参数的错误消息

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

## 最佳实践

1. **按模块组织**：错误定义按业务模块分目录
2. **国际化支持**：所有错误消息都要支持国际化
3. **有意义的错误码**：错误码要能描述错误类型
4. **正确的 HTTP 状态码**：根据错误类型返回合适的状态码
5. **不暴露敏感信息**：错误消息不要包含系统内部信息
6. **日志记录**：严重错误需要记录详细日志
