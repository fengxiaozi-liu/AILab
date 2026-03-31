# Naming

## 作用范围

本文用于说明 Kratos 业务项目中围绕聚合根的命名规则，包括聚合根命名、从属实体命名、`repo` 命名、`usecase` 命名与领域文件命名。

当问题属于以下场景时，应优先查看本文：

- 判断聚合根、实体、关系视图应如何命名
- 判断 `repo` 应如何命名
- 判断 `usecase` 应如何命名
- 判断 `proto`、`service`、`message`、`rpc` 应如何命名
- 判断 `internal/biz`、`internal/data` 下文件应如何命名
- 判断当前命名是否仍围绕同一聚合根收敛

本文重点回答围绕聚合根的命名收敛问题，不展开聚合根识别本身、协议层命名细节与共享常量命名细节。

## 前置依赖

命名判断前，应先识别：

- 稳定业务对象
- 聚合根
- 从属实体
- 关系视图或临时投影

如果名字无法稳定映射回聚合根，优先判定为建模未收敛，而不是局部命名失误。

## 命名目标

命名的目标是让领域侧结构围绕同一聚合根稳定收敛，避免按页面、动作、临时返回体或第三方字段直接命名。

命名是领域组织的外显结果，不应脱离聚合根单独演化。

## 命名规则

### 聚合根先于命名

- 先确定聚合根，再命名 `repo`、`usecase`、领域文件与领域对象
- 不允许先有 `repo`、`usecase` 或文件名，再反推聚合根

### 聚合根命名

- 聚合根名称优先表达稳定业务主语
- 不按页面、接口动作、临时展示结构或第三方字段直接命名聚合根
- 聚合根名称确定后，相关领域对象都应围绕该名称展开

### 从属实体命名

- 从属实体围绕聚合根语义命名
- 从属实体名称应表达其在聚合内的角色，而不是页面动作
- 关系视图应明确表达“关系”或“视图”语义，不伪装成聚合根本体

### `repo` 命名

- 接口定义统一使用 `XXXRepo`
- 实现类型、文件名、构造函数名围绕同一聚合根收敛
- 即使承担第三方调用职责，只要在项目中表达为领域依赖抽象，也统一使用 `XXXRepo`

正例：

```go
type AccountRepo interface{}
type SumsubRepo interface{}
```

反例：

```go
type AccountClient interface{}
type SumsubService interface{}
```

### `usecase` 命名

- `usecase` 以聚合根或稳定领域动作命名
- 稳定领域动作命名仍应能回到聚合根语义
- 不按页面名、按钮动作名、临时查询名直接命名 `usecase`

正例：

```go
type AccountUseCase struct{}
type AccountReviewUseCase struct{}
```

反例：

```go
type OpenAccountPageUseCase struct{}
type SubmitAccountReplyUseCase struct{}
```

### `proto` 命名

- `proto` 文件、`service`、`message`、`rpc` 命名应围绕聚合根或稳定领域主题收敛
- `proto` 命名优先表达稳定业务语义，不按页面、按钮动作或临时展示结构命名
- `proto` 文件名优先使用聚合根或稳定领域主题的小写下划线形式
- `service` 名称应表达稳定服务边界，不平行发明脱离聚合根的新术语
- `rpc service` 命名应围绕聚合根或稳定领域主题命名，不按页面流、按钮流或临时展示流命名
- `message` 名称应表达围绕聚合根的请求或响应语义，而不是 UI 页面语义
- `rpc` 名称可表达稳定领域动作，但仍应能回到聚合根语义

正例：

```text
account.proto
account_review.proto
service AccountService
message GetAccountRequest
message GetAccountReply
rpc GetAccount
```

反例：

```text
open_account_page.proto
submit_account_reply_store.proto
service OpenAccountPageService
message SubmitAccountReplyStore
rpc ClickSubmitButton
```

### 领域文件命名

- `internal/biz` 下文件名优先回到聚合根或稳定领域主题
- `internal/data` 下 repo 实现文件名与构造函数名应与同一聚合根保持一致
- 不按临时页面、导出动作、展示结构命名领域文件
- 同一稳定主题下，`biz` 层的领域对象、`repo` 接口与 `usecase` 可以放在同一个文件中
- 文件拆分优先按领域主题细化，不按 `_usecase`、`_repo` 这类技术角色后缀拆分
- 只有在单文件过大或主题明显分叉时，才继续拆成更细的领域主题文件

正例：

```text
account.go
account_review.go
account_collect.go
account_kyc.go
```

反例：

```text
open_account_page.go
submit_account_reply_store.go
kyc_usecase.go
account_repo.go
```

## 成组收敛

当聚合根名称确定或变更时，应将同一聚合根相关命名成组收敛，而不是局部修改。

需要同步观察的对象包括：

- 聚合根结构体
- 从属实体
- `repo` 接口
- `usecase`
- `proto` 文件
- `service`
- `message`
- `rpc`
- `biz` 文件名
- `data` repo 实现名
- 构造函数名

同一聚合根不应在 `biz` 与 `data` 中长期并存两套名字。

## 判断提示

判断一个名字是否合理时，可优先观察：

- 当前对象能否稳定映射回某个聚合根
- 这个名字表达的是聚合根，还是页面 / 动作 / 临时结构
- `repo` 是否仍符合 `XXXRepo` 约定
- `usecase` 是否围绕聚合根或稳定领域动作命名
- `proto`、`service`、`message`、`rpc` 是否仍围绕聚合根或稳定领域主题命名
- 文件名是否仍然围绕领域概念，而不是 UI 或临时返回体

## 边界提示

如果问题继续细化，应转到更具体的知识文档：

- 聚合根识别 -> `aggregate-root.md`
- 领域组织、usecase、repo、data -> `domain.md`
- service / proto 结构设计与协议边界 -> `service.md`
- 共享枚举、错误语义与稳定字面量 -> `shared-conventions.md`
