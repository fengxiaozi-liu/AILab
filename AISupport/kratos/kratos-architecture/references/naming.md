# Naming

## 作用范围

本文用于说明 Kratos 业务项目中围绕聚合根与实体的命名规则，包括聚合根命名、实体命名、`repo` 命名、`usecase` 命名与领域文件命名。

当问题属于以下场景时，应优先查看本文：

- 判断聚合根、实体、关系视图应如何命名
- 判断 `repo` 应如何命名
- 判断 `usecase` 应如何命名
- 判断 `proto`、`service`、`message`、`rpc` 应如何命名
- 判断 `internal/biz`、`internal/data`、`internal/service`、`api/*/*.proto` 下文件应如何命名
- 判断当前命名是否仍围绕同一聚合根与实体收敛

本文重点回答围绕聚合根的命名收敛问题，不展开聚合根识别本身、协议层命名细节与共享常量命名细节。

## 前置依赖

命名判断前，应先识别：

- 稳定业务对象
- 聚合根
- 其他实体
- 关系视图或临时投影

如果名字无法稳定映射回聚合根或实体，优先判定为建模未收敛，而不是局部命名失误。

## 命名目标

命名的目标是让领域侧结构围绕同一聚合根与实体稳定收敛，避免按页面、动作、临时返回体或第三方字段直接命名。

命名是领域组织的外显结果，不应脱离聚合根单独演化。

## 命名规则

### 实体与主题命名

- 先识别实体，再命名 `repo`、`usecase`、领域文件与领域对象
- 不允许先有 `repo`、`usecase` 或文件名，再反推实体或聚合根
- 实体名称优先表达稳定业务对象
- 不按页面、接口动作、临时展示结构或第三方字段直接命名实体
- 当某个实体被识别为聚合根后，相关领域对象都应围绕该实体名称展开
- 转换函数统一使用 `{Entity}Convert` 风格命名
- 不在转换函数名前附加 `inner`、`open`、`admin` 等 side 前缀
- 若同一实体存在多个转换目标，优先在 `{Entity}Convert` 基础上补充目标语义，而不是补充 side 前缀
- 转换函数命名仍应稳定映射回聚合根或实体，不按页面、接口动作或临时展示结构命名

### 实体命名

- 其他实体围绕聚合根语义命名
- 实体名称应表达其业务角色，而不是页面动作
- 关系视图应明确表达“关系”或“视图”语义，不伪装成聚合根本体
- 当某个实体已经形成稳定主题时，允许以该实体作为文件与命名主语
- 第三方若已经形成稳定围绕同一对象族操作的边界主题，也允许以该主题命名，例如 `sumsub.go`

### `repo` 命名

- 接口定义统一使用 `XXXRepo`
- 实现类型、文件名、构造函数名围绕同一聚合根、实体或稳定边界对象收敛
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

- `usecase` 以聚合根、稳定实体或稳定领域动作命名
- 稳定领域动作命名仍应能回到聚合根语义
- 不按页面名、按钮动作名、临时查询名直接命名 `usecase`
- 当某个实体需要对外暴露能力时，允许独立命名为 `EntityUseCase`

正例：

```go
type AccountUseCase struct{}
type AccountFlowPageUseCase struct{}
```

反例：

```go
type OpenAccountPageUseCase struct{}
type SubmitAccountReplyUseCase struct{}
```

### `proto` 命名

- `proto` 文件、`service`、`message`、`rpc` 命名应围绕聚合根、实体或稳定第三方对象边界收敛
- `proto` 命名优先表达稳定业务语义，不按页面、按钮动作或临时展示结构命名
- `proto` 文件名优先使用聚合根、实体或稳定第三方对象边界的小写下划线形式
- `service` 名称应表达稳定服务边界，不平行发明脱离聚合根的新术语
- `rpc service` 命名应围绕聚合根、实体或稳定边界对象命名，不按页面流、按钮流或临时展示流命名
- `message` 名称应表达围绕聚合根的请求或响应语义，而不是 UI 页面语义
- `rpc` 名称可表达稳定领域动作，但仍应能回到聚合根语义

正例：

```text
account.proto
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

#### `rpc` 命名公式

`rpc` 统一采用：

`{Operate}{AggregateRoot|Entity}`

其中：

- `Operate` 表示稳定业务动作
- `AggregateRoot|Entity` 表示聚合根或稳定实体

#### 查询类动作

- 单条获取：`Get{AggregateRoot|Entity}`
- 列表获取：`List{AggregateRoot|Entity}`
- 分页获取：`PageList{AggregateRoot|Entity}`
- 统计获取：`Get{AggregateRoot|Entity}Statistics`

正例：

```text
GetAccount
ListAccountReview
PageListAccount
GetAccountStatistics
```

反例：

```text
QueryAccount
QueryAccountListByPage
FetchAccountInfo
SearchAccount
```

#### 通用写入类动作

- 创建：`Create{AggregateRoot|Entity}`
- 更新：`Update{AggregateRoot|Entity}`
- 删除：`Delete{AggregateRoot|Entity}`

正例：

```text
CreateAccount
UpdateAccount
```

#### 业务动作类

业务动作类统一采用：

`{Operate}{AggregateRoot|Entity}`

推荐 `Operate` 词表包括：

- `Prepare`
- `Commit`
- `Submit`
- `Review`
- `Retry`
- `Start`
- `Finish`
- `Cancel`
- `Enable`
- `Disable`
- `Approve`
- `Reject`

正例：

```text
PrepareAccount
CommitAccount
ReviewAccount
CancelAccountReview
```

反例：

```text
HandleAccountCallback
ClickSubmitAccount
ProcessReviewResult
```

### 领域文件命名

- `internal/biz` 下文件名优先回到聚合根、实体或稳定边界对象
- `internal/data` 下 repo 实现文件名与构造函数名应与同一聚合根、实体或稳定边界对象保持一致
- 不按临时页面、导出动作、展示结构命名领域文件
- 稳定实体必须有对应文件；一旦某个实体定义了 `UseCase`，则该实体、对应 `Repo`、对应 `UseCase` 必须在同一个文件中
- 文件拆分优先按实体或稳定边界对象细化，不按 `_usecase`、`_repo` 这类技术角色后缀拆分
- 只有在单文件过大或边界明显分叉时，才继续拆成更细的文件

正例：

```text
account.go
account_collect.go
```

反例：

```text
open_account_page.go
submit_account_reply_store.go
account_repo.go
```

### 目录与命名映射

命名规范在目录层面应这样收敛：

- `internal/biz/{business}`：聚合根、实体、`UseCase`、`Repo`
- `internal/data/{business}`：repo 实现、构造函数、数据访问文件
- `internal/service/{side}/{version}`：`Service` 命名
- `api/{side}/{domain}/{version}`：`proto`、`service`、`message`、`rpc`
- `internal/listener/*`、`internal/consumer/*`：事件入口文件与类型命名

### 聚合根与实体的标准组织

当某个实体已经形成稳定主题时，推荐组织为：

```text
internal/biz/open/account.go
internal/data/open/account.go
api/open/open/v1/account.proto
```

其中：
- `account.go` 定义聚合根 `Account`
- `Account` 通过 `AccountFlowPageInfo *AccountFlowPage` 与其他实体建立关系表达
- `data` 层继续使用 `account_flow_page.go` 对齐实体，而不是改成 `account_flow_page_repo.go`

## 成组收敛

当聚合根名称确定或变更时，应将同一聚合根相关命名成组收敛，而不是局部修改。

需要同步观察的对象包括：

- 聚合根结构体
- 其他实体
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
- `usecase` 是否围绕聚合根、稳定实体或稳定领域动作命名
- `proto`、`service`、`message`、`rpc` 是否仍围绕聚合根、实体或稳定边界对象命名
- 文件名是否仍然围绕领域概念，而不是 UI 或临时返回体
- `rpc` 是否遵循 `{Operate}{AggregateRoot|Entity}`

## 边界提示

如果问题继续细化，应转到更具体的知识文档：

- 聚合根识别 -> `aggregate-root.md`
- 领域组织、usecase、repo、data -> `domain.md`
- service / proto 结构设计与协议边界 -> `service.md`
- 共享枚举、错误语义与稳定字面量 -> `shared-conventions.md`
