# Error And Enum

## 作用范围

本文用于说明 Kratos 项目中的错误语义、`not found`、枚举、typed const、稳定值域与稳定协议字面量的收敛方式。

当问题属于以下场景时，应优先查看本文：

- 判断错误语义、reason、`not found` 应如何统一
- 判断错误 helper 应如何命名与落位
- 判断 enum、typed const、稳定值域应如何收敛
- 判断稳定协议字面量、header、path、query key、事件名是否应提取
- 判断第三方返回值如何映射成项目内稳定语义

本文重点回答错误与枚举这组共享语义如何统一，不展开 i18n、logging、comment 与代码红线。

## 错误语义

定位：稳定业务错误语义与共享错误表达。

承载内容：

- 模块级共享错误语义
- `reason`
- `not found`
- 统一错误 helper

边界提示：

- 错误语义应与领域语义一致，不改写成模糊通用错误
- `not found` 必须表达“确实未找到”，不伪装成功
- 共享错误定义默认落在 `internal/error/<module>`
- 通用公共错误可落在 `internal/error/base`
- `DefaultMessage.Description` 与 `DefaultMessage.Other` 统一使用中文简体，不使用英文说明

命名模式：

- 函数名：`Error{Module}{Description}`
- reason：`{MODULE}_{DESCRIPTION}`
- i18n ID：`ERROR_{MODULE}_{DESCRIPTION}`

示例：

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

```go
func ErrorSumsubRequestFail(ctx context.Context) *errors.Error {
	localizer := context2.GetLocalize(ctx)
	return errors.New(500, "OPEN_SUMSUB_REQUEST_FAIL", localize.LocalizeOrUseOther(localizer, &i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:          "ERROR_OPEN_SUMSUB_REQUEST_FAIL",
			Description: "Sumsub 请求失败",
			Other:       "Sumsub 请求失败",
		},
	}))
}
```

## `not found`

`not found` 是错误语义的一部分，不是成功分支。

规则：

- data 层或 repo 层识别到 not found 时，转成统一业务 not found 错误
- 上层若需要区分 not found，判断统一错误语义
- 禁止 `nil, nil`
- 禁止空对象伪装成功

示例：

```go
if ent.IsNotFound(err) {
	return nil, openerror.ErrorAccountNotFound(ctx)
}
```

## 枚举与稳定值域

定位：共享值域、typed const 与稳定协议常量。

承载内容：

- 状态
- 原因
- 结果
- 类型
- 协议常量

边界提示：

- 稳定业务值域优先类型化
- 第三方返回值若进入业务判断，应先收敛成项目内 enum 或常量
- 实体、事件、审核结果、filter 固定选项中的稳定值域字段，应直接使用 enum 或 typed const，而不是裸 `string` / `int`
- 协议常量与业务 enum 分层表达，不混放
- 共享值域默认落在 `internal/enum/<module>`

### 实体字段枚举化约束

当字段满足以下任一条件时，应优先定义为枚举类型，而不是裸类型：

- 字段值来自有限集合，例如状态、结果、原因、类型、阶段、动作
- 字段会参与 `switch`、分支判断、状态流转或结果映射
- 字段会在多个函数、文件或层之间被重复消费

规则如下：

- 原始第三方响应可以暂时保留原始值类型
- 一旦进入 `biz` 层实体、领域事件、审核结果对象或其他可被业务消费的对象，应转换为项目内 enum 或 typed const
- 不要只在 `switch` 处补枚举，而让实体字段继续长期保持裸 `string` / `int`
- 如果字段保留原始类型，必须明确它只用于透传、记录或原始 payload 存储，不参与业务判断

### 枚举落位约束

枚举的落位是强约束，不是建议项。

规则如下：

- 稳定业务值域只能定义在 `internal/enum/<domain>`
- 不允许在 `internal/biz/<domain>`、`internal/data/<domain>`、`internal/service/<domain>` 或单个 repo/usecase 文件中直接定义稳定枚举
- 一般字符串、数字，只要承载稳定共享语义并参与状态流转、分支判断、结果映射或协议映射，也应收敛为 enum、typed const 或具名常量
- 不能因为“只在一个专项里使用”就直接在业务文件中定义裸值
- 如果属于专项值域或专项协议常量，就在 `internal/enum/<domain>` 下单独建立一个文件收口，而不是内联在业务文件里
- 不要为了省事把新枚举继续堆进通用聚合文件，例如已有 `open.go` 时，`account_kyc`、`sumsub` 这类专项值域应优先新建 `account_kyc.go`、`sumsub.go`
- 稳定协议字面量同样进入 `internal/enum/<domain>`，但根据语义选择“业务 enum”或“具名常量”表达，而不是散落在业务文件中
- 不要把“`const` 只是当前文件使用”当作不收敛的理由；只要它是稳定共享语义常量，就应回到 `internal/enum/<domain>`

示例：

```text
// ✅
internal/enum/open/account_review.go
internal/enum/open/account_callback.go
internal/enum/open/sumsub.go
internal/enum/order/settlement.go

// ❌
internal/biz/open/account.go
internal/data/open/account_repo.go
internal/service/open/account.go
```

示例：

```go
type ReviewStatus string

const (
	ReviewStatusPassed ReviewStatus = "reviewed-pass"
	ReviewStatusFailed ReviewStatus = "reviewed-fail"
)
```

```go
const HeaderWebhookSignature = "X-Payload-Digest"
```

```go
type KycReviewResult struct {
	ReviewAnswer AccountKycReviewAnswer `json:"review_answer"`
}
```

## 稳定协议字面量

当以下字面量具备稳定共享语义时，应优先收敛：

- header 名
- path 模板
- query key
- 事件名
- 固定状态值

规则：

- 业务状态值优先定义为 enum 或 typed const
- 协议固定字面量定义为具名常量
- 第三方协议返回值先收敛到项目内语义，再进入业务判断
- 稳定协议字面量默认也落在 `internal/enum/<domain>`，不散落在 repo / service / biz 文件中
- 第三方 client、repo、service 文件中新增的 header、path、query key、事件名与固定状态值，先判断是否属于稳定共享语义常量；若是，先收敛再继续实现
- 第三方 endpoint、header、path、query key 如果明显围绕某个专项能力，例如 `sumsub`、`account_kyc`，应优先落到对应专项 enum 文件，而不是继续追加到通用 `open.go`

示例：

```go
const (
	sumsubHeaderAppToken              = "X-App-Token"
	sumsubHeaderAccessSig             = "X-App-Access-Sig"
	sumsubPathApplicantByExternalUser = "/resources/applicants/-;externalUserId=%s/one"
)
```

### 枚举本地化约束

当 enum / typed const 提供 `Localize` 能力时，`i18n` `ID` 仍然属于稳定共享语义的一部分，应继续遵守 enum 收敛规则。

规则如下：

- enum 的 `i18n` `ID` 必须使用固定语义常量，不通过字符串拼接动态构造
- 若文案需要动态变量，占位部分通过 `TemplateData` 传入，不用拼接 `ID`
- 不要把第三方原始值直接拼进 `ID`，而是先收敛为项目内 enum，再由 enum 输出固定 `ID`

正确示例：

```go
func (s AccountKycReviewStatus) Localize(localizer *i18n.Localizer) string {
	switch s {
	case AccountKycReviewStatusPending:
		return localize.LocalizeOrUseOther(localizer, &i18n.LocalizeConfig{
			DefaultMessage: &i18n.Message{
				ID:          "ENUM_OPEN_KYC_REVIEW_STATUS_PENDING",
				Description: "KYC review status-pending",
				Other:       "pending",
			},
		})
	default:
		return localize.LocalizeOrUseOther(localizer, &i18n.LocalizeConfig{
			DefaultMessage: &i18n.Message{
				ID:          "ENUM_OPEN_KYC_REVIEW_STATUS_INIT",
				Description: "KYC review status-init",
				Other:       "init",
			},
		})
	}
}
```

```go
func (e SystemSmsEvent) SMSSubject(localizer *i18n.Localizer, templateData map[string]interface{}) string {
	return localize.LocalizeOrUseOther(localizer, &i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:          "ENUM_SYSTEM_SMS_EVENT_REGISTER_SUBJECT",
			Description: "注册短信标题",
			Other:       "{{.code}} 为您的验证码",
		},
		TemplateData: templateData,
	})
}
```

错误示例：

```go
func (s AccountKycReviewStatus) Localize(localizer *i18n.Localizer) string {
	return localize.LocalizeOrUseOther(localizer, &i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:          "ENUM_OPEN_KYC_REVIEW_STATUS_" + string(s),
			Description: "KYC review status",
			Other:       string(s),
		},
	})
}
```

## 共享规则

围绕错误与枚举收敛时，应优先遵守以下规则：

- 错误语义与领域语义保持一致
- `not found` 不伪装成功
- 稳定值域优先类型化
- 稳定协议字面量集中定义并统一落在 `internal/enum/<domain>`
- 错误、reason、enum、常量都应回到共享目录收敛

## 判断提示

判断一个表达是否应纳入错误或枚举收敛时，可优先观察：

- 它是否会跨函数、跨文件或跨层复用
- 它是否承载稳定业务语义或稳定协议语义
- 它是否已经在多个位置出现并开始漂移
- 它更像错误语义、业务值域、协议常量还是第三方边界值
- 若它承载稳定共享语义，优先收敛到 `internal/error/<domain>` 或 `internal/enum/<domain>`

## 边界提示

如果问题继续细化，应转到更具体的知识文档：

- 聚合根识别 -> `aggregate-root.md`
- 围绕聚合根的命名收敛 -> `naming.md`
- i18n、logging、comment -> `shared-conventions.md`
- 代码红线与反防御式编程 -> `code-style.md`
## 第三方 DTO 与稳定协议字面量的边界

处理第三方集成代码时，禁止按“当前用了几处”来决定落位，必须先按语义类型判断。

### 判定顺序

1. 先判断它是 DTO，还是稳定协议字面量。
2. 如果是 DTO，再判断是否应留在 `internal/data/...`。
3. 如果是稳定协议字面量，直接收敛到 `internal/enum/<domain>`。

### DTO

以下内容可以留在 `internal/data/...`：
- 请求体结构体
- 响应体结构体
- payload 结构体
- 原始 JSON 字段名与 tag
- 仅用于解析、透传、记录原始 payload 的原始值

### 稳定协议字面量

以下内容只要承载稳定共享语义，就必须收敛到 `internal/enum/<domain>`：
- header 名
- path 模板
- query key
- content-type
- 第三方固定状态值
- 第三方固定结果值
- 第三方固定动作名
- 第三方固定事件名

### 禁止的误判

- 不要因为它当前只出现在一个 repo 文件里，就把它当成局部实现细节。
- 不要把 DTO 结构体和协议常量混成一类做同一次落位判断。
- 不要把“已经集中在一个文件里”当成“已经收敛到规定共享目录”。

### Checklist 复核问题

复核第三方集成代码时，必须显式回答：
- 这是 DTO，还是稳定协议字面量？
- 如果它是稳定协议字面量，为什么还没有进入 `internal/enum/<domain>`？
- 我是否正在用“当前只在一个文件里使用”作为例外理由？
