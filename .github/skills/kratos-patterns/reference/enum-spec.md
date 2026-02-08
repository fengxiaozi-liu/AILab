# 枚举定义规范

## 概述

本项目枚举定义位于 `internal/enum/` 目录，按业务模块分类组织。

## 目录结构

```
internal/enum/
├── base/           # 基础枚举
│   ├── base.go     # 通用枚举
│   ├── amqp.go     # 消息队列相关
│   ├── eventbus.go # 事件总线相关
│   └── ...
├── admin/          # 管理模块枚举
│   ├── admin.go     # 管理模块主枚举
│   └── foo.go     # 管理模块其他枚举
├── system/         # 系统模块枚举
│   ├── system.go     # 系统模块主枚举
│   └── foo.go     # 系统模块其他枚举
├── user/           # 用户模块枚举
│   ├── user.go     # 用户模块主枚举
│   └── foo.go     # 用户模块其他枚举
├── gateway/        # 用户模块枚举
│   ├── gateway.go     # 用户模块主枚举
│   └── foo.go     # 系统模块其他枚举
└── ...
```

## 枚举定义规范

### 数值类型枚举

- 数值类型枚举使用具体的数值类型（如 `uint32`、`int` 等）作为底层类型。
- 数值类型枚举禁止使用`0`作为有效值，除非明确表示`未设置`或`未知`。

```go
package user

type UserStatus uint32

const (
	UserStatusNormal   UserStatus = 1 // 正常
	UserStatusDisabled UserStatus = 2 // 已禁用
)

func (s UserStatus) Value() uint32 {
	return uint32(s)
}

func UserStatusList() []UserStatus {
	return []UserStatus{UserStatusNormal, UserStatusDisabled}
}
```

### 消息队列 Topic 枚举

- 消息队列 Topic 使用`baseenum.AmqpTopic`作为基础类型。
- 消息队列 Topic 枚举名称格式为`AmqpTopic{Entity}{Action}`，如`AmqpTopicUserCreateAfter`。
- 消息队列 Topic 枚举值使用小驼峰命名法，如`userCreateAfter`。

```go
package base

type AmqpTopic string

const (
	AmqpTopicUserCreateAfter AmqpTopic = "userCreateAfter"
	AmqpTopicOrderNotify     AmqpTopic = "orderNotify"
)

func (t AmqpTopic) Value() string {
	return string(t)
}

```

### 本地事件枚举

- 本地事件使用`baseenum.LocalEvent`作为基础类型。
- 本地事件枚举名称格式为`LocalEvent{Entity}{Action}`，如`LocalEventLifecycleBeforeStart`。
- 本地事件枚举值使用小驼峰命名法，如`lifecycleBeforeStart`。

```go
package base

type LocalEvent string

const (
	LocalEventLifecycleBeforeStart LocalEvent = "lifecycleBeforeStart"
	LocalEventLifecycleAfterStart  LocalEvent = "lifecycleAfterStart"
	LocalEventLifecycleBeforeStop  LocalEvent = "lifecycleBeforeStop"
	LocalEventLifecycleAfterStop   LocalEvent = "lifecycleAfterStop"
)

func (e LocalEvent) Value() string {
	return string(e)
}
```

## 命名规范

### 类型命名

格式：`{Entity}{Field}` 或 `{Concept}`

```go
package main

// 好的命名
type UserStatus uint32
type OrderType string
type PaymentMethod string
type Language string

// 避免的命名
type Status uint32 // 太笼统
type Type string   // 太笼统

```

### 后缀语义（命名层面）

> 选择后缀的目的：让读代码的人通过类型名就能判断业务语义类别。

Good:
```go
type CustomerType uint8   // Type：类型/分类
type OrderStatus uint8    // Status：状态/阶段
type WebsocketEvent uint8 // Event：事件
type OrderApprovalReason uint32 // Reason：原因
type TenantPushChannel uint8 // Channel：渠道/通道
```

Bad:
```go
type CustomerStatus uint8 // BAD: 实际是“类型/分类”却叫 Status
type OrderType uint8      // BAD: 实际是“状态/阶段”却叫 Type
```

### 常量命名

格式：`{Type}{Value}`

```go
package main

type UserStatus uint32
type OrderType string

// 好的命名
const (
	UserStatusActive   UserStatus = 1
	UserStatusInactive UserStatus = 2
)
const (
	OrderTypeNormal OrderType = "normal"
	OrderTypeRefund OrderType = "refund"
)

// 避免的命名
const (
	Active UserStatus = 1        // 缺少类型前缀
	Normal OrderType  = "normal" // 缺少类型前缀
)

```

### 业务枚举必须用语义化类型表达（禁止用原始类型传递语义）

Good:
```go
type TradeAccountAssetChangeType uint8

const (
	TradeAccountAssetChangeTypePayment  TradeAccountAssetChangeType = iota + 1
	TradeAccountAssetChangeTypeWithdraw
)

func ApplyChange(t TradeAccountAssetChangeType) {}
```

Bad:
```go
// BAD: 用 uint8/int32 等原始类型传递“类型/状态”等业务枚举语义，见名不知意
func ApplyChange(t uint8) {}
```

## 必要方法

每个枚举类型应该实现以下方法：

```go
package userenum

import "github.com/nicksnyder/go-i18n/v2/i18n"

type UserStatus uint32

// Value 返回底层值
func (s UserStatus) Value() uint32 {
	return uint32(s)
}

// UserStatusList 返回所有枚举值
func UserStatusList() []UserStatus {
	return []UserStatus{UserStatusInactive, UserStatusActive, UserStatusDisabled}
}

// UserStatusDefault 返回默认值（可选）
func UserStatusDefault() UserStatus {
	return UserStatusInactive
}

// Localize 国际化显示（需要时）
func (s UserStatus) Localize(localizer *i18n.Localizer) string {
	// ...
}

// IsValid 验证是否有效（可选）
func (s UserStatus) IsValid() bool {
	for _, v := range UserStatusList() {
		if v == s {
			return true
		}
	}
	return false
}

```

## 国际化消息 ID 规范

枚举值的国际化消息 ID 格式：`ENUM_{MODULE}_{TYPE}_{VALUE}`

```go
package user

var msg = &i18n.Message{
	ID:    "ENUM_MARKET_SECURITY_STATUS_LISTED", // 格式: ENUM_{服务}_{类型}_{值}
	Other: "已激活",
}

```

## 使用示例

### 在 Biz 层使用

```go
package user

import "context"
import userenum "{module from go.mod}/internal/enum/user"

type User struct {
	Status userenum.UserStatus
}

type UserRepo interface {
    Create(ctx context.Context, info *User) (*User, error)
}

type UserUseCase struct{
	userRepo UserRepo
}

func (u *UserUseCase) Create(ctx context.Context, info *User) (*User, error) {
	// 设置默认状态
	info.Status = userenum.UserStatusInactive
	return u.userRepo.Create(ctx, info)
}

```

### 在 Service 层使用（带国际化）

```go
package user

import (
	"context"

	v1 "{module from go.mod}/internal/api/user/open/v1"
	userbiz "{module from go.mod}/internal/biz/user"
)

type UserService struct {
	userUseCase *userbiz.UserUseCase
}

func (srv *UserService) CreateUser(ctx context.Context, req *v1.CreateUserRequest) (*v1.User, error) {
	res, err := srv.userUseCase.Create(ctx, &userbiz.User{
		// ...
	})
	if err != nil {
		return nil, err
	}
	return userConvert(res), nil
}

func userConvert(info *User) *v1.User {
	return &v1.User{
		Status: uint32(info.Status.Value()),
	}
}

```

## 最佳实践

1. **类型安全**：使用自定义类型而非原始类型
2. **完整性**：提供 List、Default 等辅助方法
3. **国际化**：需要展示的枚举实现 Localize 方法
4. **注释**：每个枚举值都要有注释说明
5. **按模块组织**：枚举按业务模块分目录存放
6. **一致性**：同一类型的枚举值命名风格要一致
