# Error And Enum

## 作用范围

本文用于说明 Kratos 项目中的错误语义、`not found`、枚举、typed const、稳定值域与稳定协议字面量的收敛方式。

当问题属于以下场景时，应优先查看本文：

- 判断错误语义、`reason`、`not found` 应如何统一
- 判断错误 helper 应如何命名与落位
- 判断 enum、typed const、稳定值域应如何收敛
- 判断稳定协议字面量、header、path、query key、事件名是否应提取
- 判断第三方返回值如何映射成项目内稳定语义

本文重点回答错误与枚举这组共享语义如何统一，不展开 i18n、logging、comment 与代码红线。

## 规则

### 错误语义

- 稳定错误语义必须集中定义，`not found` 必须表达真实未找到，不得伪装成成功
- 稳定错误默认落在 `internal/error/<module>`
- 通用公共错误可落在 `internal/error/base`
- `DefaultMessage.Description` 与 `DefaultMessage.Other` 统一使用中文简体
- `Error{Module}{Description}`、`{MODULE}_{DESCRIPTION}`、`ERROR_{MODULE}_{DESCRIPTION}` 保持同一语义收敛

### 枚举与稳定值域

- 稳定业务值域应优先收敛为项目内 enum、typed const 或具名常量，不使用裸 `string` / `int` 参与业务判断
- 实体、事件、审核结果、filter 固定选项等业务消费对象中的稳定值域字段，应直接使用 enum 或 typed const
- 第三方原始响应可暂保留原始值类型；一旦进入 `biz` 层实体、领域事件、审核结果对象或其他业务消费对象，必须转换为项目内语义
- 不要只在局部 `switch` 或分支中补枚举，而让字段长期保持裸值类型
- 若字段保留原始类型，必须明确仅用于透传、记录或原始 payload 存储，不参与业务判断
- 稳定业务值域与稳定协议字面量统一落在 `internal/enum/<domain>`
- 不允许在 `internal/biz/<domain>`、`internal/data/<domain>`、`internal/service/<domain>` 或单个 repo/usecase 文件中直接定义稳定枚举或稳定共享常量
- 不得以“只在一个专项”或“只在当前文件使用”为理由保留裸值；专项值域或专项协议常量也应在 `internal/enum/<domain>` 下独立文件收口
- 不要把专项值域继续堆入通用聚合文件；例如已有 `open.go` 时，`account_kyc`、`sumsub` 应优先拆到独立文件
- 稳定协议字面量进入 `internal/enum/<domain>` 后，应按语义选择“业务 enum”或“具名常量”表达
- 如果枚举是数字类型，必须通过 `iota` 定义，且合法业务值不能从 `0` 开始；`0` 保留给“全部”或未指定筛选态

### 枚举本地化

- enum 的 `i18n` `ID` 必须使用固定语义常量，不通过字符串拼接动态构造
- 若文案需要动态变量，占位部分通过 `TemplateData` 传入，不用拼接 `ID`
- 不要把第三方原始值直接拼进 `ID`，而是先收敛为项目内 enum，再由 enum 输出固定 `ID`
- 对业务枚举，`Localize()` 属于三段式的一部分，应与枚举值定义、`List()` 一起表达

## 错误语义

定位：稳定业务错误语义与共享错误表达。

承载内容：

- 模块级共享错误语义
- `reason`
- `not found`
- 统一错误 helper

边界提示：

- 错误语义应与领域语义一致，不改写成模糊通用错误
- `not found` 必须表达“确实未找到”，不伪装成成功
- 共享错误定义默认落在 `internal/error/<module>`

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

### `not found`

`not found` 是错误语义的一部分，不是成功分支。

规则：

- data 层或 repo 层识别到 not found 时，应转换成统一业务 not found 错误
- 上层若需要区分 not found，应判断统一错误语义
- 禁止 `nil, nil`
- 禁止空对象伪装成功

示例：

```go
if ent.IsNotFound(err) {
	return nil, openerror.ErrorAccountNotFound(ctx)
}
```

## 枚举与稳定值域

### 职责描述

枚举与稳定值域用于承载项目内可复用、可判断、可展示的共享语义，覆盖状态、原因、结果、类型以及稳定协议常量。

它的职责是：

- 为业务判断提供稳定值域表达
- 为跨层传递提供统一类型语义
- 为本地化展示提供固定语义出口
- 为第三方协议接入提供统一收敛点

定义上，以下情况默认属于枚举或稳定值域的适用范围：

- 字段值来自有限集合，例如状态、原因、结果、类型、阶段、动作
- 字段会参与 `switch`、分支判断、状态流转或结果映射
- 一般字符串、数字只要承载稳定共享语义并参与状态流转、分支判断、结果映射或协议映射，也应收敛为 enum、typed const 或具名常量
- 对可列举的业务枚举，默认使用三段式表达：枚举值定义、`List()`、`Localize()`；如需透出底层值，再补 `Value()`

边界上，业务 enum 与稳定协议字面量都进入 `internal/enum/<domain>`，但根据语义分层表达，不混成一类实现。

### 枚举本地化约束

当 enum / typed const 提供 `Localize` 能力时，`i18n` `ID` 仍然属于稳定共享语义的一部分，应继续遵守 enum 收敛规则。

规则如下：

- enum 的 `i18n` `ID` 必须使用固定语义常量，不通过字符串拼接动态构造
- 若文案需要动态变量，占位部分通过 `TemplateData` 传入，不用拼接 `ID`
- 不要把第三方原始值直接拼进 `ID`，而是先收敛为项目内 enum，再由 enum 输出固定 `ID`

### 示例

#### 文件位置与命名示例

```text
// ✅
internal/enum/open/account.go
internal/enum/open/account_kyc.go
internal/enum/order/settlement.go

// ❌
internal/biz/open/account.go
internal/data/open/account_repo.go
internal/service/open/account.go
```

#### 业务枚举

```go
type AccountKycRejectType uint8

const (
	AccountKycRejectTypeNone AccountKycRejectType = iota + 1
	AccountKycRejectTypeRetry
	AccountKycRejectTypeFinal
)

func AccountKycRejectTypeList() []AccountKycRejectType {
	return []AccountKycRejectType{
		AccountKycRejectTypeNone,
		AccountKycRejectTypeRetry,
		AccountKycRejectTypeFinal,
	}
}

func (t AccountKycRejectType) Localize(localizer *i18n.Localizer) string {
	switch t {
	case AccountKycRejectTypeRetry:
		return localize.LocalizeOrUseOther(localizer, &i18n.LocalizeConfig{
			DefaultMessage: &i18n.Message{
				ID:          "ENUM_OPEN_ACCOUNT_KYC_REJECT_TYPE_RETRY",
				Description: "KYC 拒绝类型-可重试",
				Other:       "可重试",
			},
		})
	case AccountKycRejectTypeFinal:
		return localize.LocalizeOrUseOther(localizer, &i18n.LocalizeConfig{
			DefaultMessage: &i18n.Message{
				ID:          "ENUM_OPEN_ACCOUNT_KYC_REJECT_TYPE_FINAL",
				Description: "KYC 拒绝类型-最终失败",
				Other:       "最终失败",
			},
		})
	default:
		return localize.LocalizeOrUseOther(localizer, &i18n.LocalizeConfig{
			DefaultMessage: &i18n.Message{
				ID:          "ENUM_OPEN_ACCOUNT_KYC_REJECT_TYPE_NONE",
				Description: "KYC 拒绝类型-未返回",
				Other:       "未返回",
			},
		})
	}
}

func (t AccountKycRejectType) Value() uint8 { return uint8(t) }
```

#### 稳定协议字面量

```go
const (
	sumsubHeaderAppToken              = "X-App-Token"
	sumsubHeaderAccessSig             = "X-App-Access-Sig"
	sumsubPathApplicantByExternalUser = "/resources/applicants/-;externalUserId=%s/one"
)
```

## 稳定语义常量落位规则

当一个常量用于表达查询装配、关联加载、聚合展开、共享过滤、统一状态映射等可跨多个实现点复用的稳定业务语义时，应将其视为业务域级共享常量。

业务域级共享常量应收口到该业务域统一的 `enum/constants` 入口文件中集中管理，不应散落在局部实现文件中。

局部实现文件只承载当前实现范围内私有、不可复用、语义不稳定的常量，不承载业务域级共享常量。

常量是否收口到统一入口，判断依据是“复用范围”和“语义稳定性”，而不是“它最先在哪个文件中出现”。

### 判定标准

- 如果一个常量会在多个 `repo / usecase / service / assembler / proto convert` 中重复使用，按业务域级共享常量处理
- 如果一个常量表达的是稳定业务语义，而不是单次实现细节，按业务域级共享常量处理
- 如果一个常量只服务于单个文件内部实现，且没有跨层复用价值，可以保留在局部文件
- 常量落位优先服从语义层级，再考虑创建位置

### 反例模式

- 不要因为某个常量最早在某个局部实现中出现，就把它长期定义在该局部文件中
- 不要把已经跨多个实现点复用的共享常量继续放在局部实现文件里

## 边界提示

如果问题继续细化，应转到更具体的知识文档：

- 聚合根识别 -> `aggregate-root.md`
- 围绕聚合根的命名收敛 -> `naming.md`
- i18n、logging、comment -> `shared-conventions.md`
- 代码红线与反防御式编程 -> `code-style.md`
