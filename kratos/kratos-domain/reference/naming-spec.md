# Naming Reference

## 这个主题解决什么问题
说明 Kratos 项目里文件、Proto、UseCase、Repo、Service、结构体和字段命名如何保持同一套业务语义体系，让实现者在跨层实现和重构时能稳定复用已有术语，而不是每层重新发明一套名字。

## 适用场景

- 新增聚合根、实体、Repo、UseCase、Proto message
- 新增 RPC、事件、DTO、结构体字段
- 重命名历史漂移的领域符号
- 统一多层之间的术语

## 设计意图

命名参考的重点是让不同层围绕同一套业务词汇沟通，而不是每层发明一套自己的名字。

- 稳定命名能在 `aggregate`、`repo`、`usecase`、`proto` 间快速建立映射关系。
- 业务名词统一后，搜索、重构和审查成本都会下降。
- 当名称被接口动作带偏时，领域模型通常会跟着碎片化。

## 实施提示

- 先确定业务主体名词，再派生 `UseCase`、`Repo`、`Service`、`Message` 名称。
- 看到 `Get`、`Handle`、`Do` 这类动作词时，优先把它们放到方法名而不是对象名。
- 如果一个词只能解释单个接口，而不能解释业务对象本身，通常不适合做领域核心命名。

## 命名总览

- 聚合根、实体、Repo、UseCase 优先复用同一领域词根。
- 目录已经表达的语义，不必在类型名里重复一遍。
- 方法名表达动作，对象名表达领域概念。
- Proto service、RPC、message、field 应该和领域模型保持同一语义来源。

## 文件命名

### 文件命名参考

| 类型 | 推荐形式 | 示例 |
|------|------|------|
| Go 源文件 | 小写下划线 | `account.go`, `account_flow_page.go` |
| Go 测试文件 | `_test.go` 后缀 | `account_test.go` |
| Go 生成文件 | `_gen.go` 后缀 | `wire_gen.go` |
| Proto 文件 | 小写下划线 | `account.proto`, `account_collect.proto` |
| YAML 配置 | 小写下划线 | `config.yaml` |

### 文件与类型对应示例

正确：

- `account.go` -> `Account`, `AccountRepo`, `AccountUseCase`
- `account_flow_page.go` -> `AccountFlowPage`, `AccountFlowPageRepo`, `AccountFlowPageUseCase`
- `account_collect.go` -> `AccountCollect`, `AccountCollectRepo`, `AccountCollectUseCase`

容易漂移的写法：

- `review.go` -> `ReviewUseCase`
- `open_flow.go` -> `OpenFlowUseCase`

### 包命名参考

统一使用小写字母，简短且有意义：

```go
package account
package admin
package open
```

## Proto 命名

### Service 命名

```proto
service AccountService {
  rpc GetAccount(GetAccountRequest) returns (GetAccountReply);
}
```

### RPC 命名思路

- 方法名通常采用 `{Operate}{Domain}` 形式
- `Get` 更适合表达详情查询
- `List` 更适合表达非分页列表
- `PageList` 更适合表达分页查询
- 筛选条件更适合放进 `Request` 字段，而不是直接进方法名

更稳定的写法：

- `GetAccount`
- `ListAccount`
- `PageListAccount`
- `DeleteAccountFlow`

容易让语义变散的写法：

- `GetAccountDetail`
- `DeleteAccountFlowByAccountID`

### Message 命名

```proto
message GetAccountRequest {
  uint32 id = 1;
}

message ListAccountReply {
  repeated Account list = 1;
}
```

推荐围绕聚合根和稳定实体命名：

- `GetAccountFlowPageRequest`
- `GetAccountCollectReply`
- `AccountFlowPage`

容易漂移的写法：

- `FlowPageDataResult`
- `OpenInfoRequest`
- `ReviewDetailReply`

### Field 命名

字段通常使用 `snake_case`：

```proto
message Account {
  uint32 id = 1;
  string account_no = 2;
  uint32 create_time = 3;
}
```

当聚合根语境已经明确时，主对象 ID 直接使用 `id`，不要重复追加聚合根前缀。

正例：

```proto
message DeleteAccountRequest {
  uint32 id = 1;
}
```

反例：

```proto
message DeleteAccountRequest {
  uint32 account_id = 1;
}
```

### 字段后缀示例

- 单对象字段常见为 `_info`
- 列表字段常见为 `_list`

更稳定的写法：

- `account_info`
- `first_check_user_info`
- `page_list`

容易歧义的写法：

- `account`
- `first_check_user`
- `pages`
- `account_detail`

### 结构体注释与 JSON Tag

- 聚合根 DTO、实体 DTO、关键入参 DTO、关键事件 DTO 更适合补全结构体注释。
- 参与协议对齐、JSON 编解码或事件投递的结构体字段，建议统一补全 `json` tag，并使用 `snake_case`。

示例：

```go
type Account struct {
    ID                 uint32          `json:"id"`
    Status             AccountStatus   `json:"status"`
    CreateTime         uint32          `json:"create_time"`
    FirstCheckUserInfo *AdminUser      `json:"first_check_user_info"`
    AccountFlowPageList []*AccountFlowPage `json:"account_flow_page_list"`
}
```

## UseCase 命名

### 类型命名

| 类型 | 推荐命名 | 示例 |
|------|------|------|
| UseCase | `{Domain}UseCase` | `AccountUseCase` |
| 方法 | `{Action}` | `Review`, `Prepare`, `GetMeta` |

### 构造函数

```go
func NewAccountUseCase(repo AccountRepo) *AccountUseCase
```

### 行为命名示例

更贴近领域的写法：

- `AccountUseCase.Review`
- `AccountUseCase.Prepare`
- `AccountFlowPageUseCase.GetMeta`
- `AccountCollectUseCase.GetAccountCollect`

容易漂移的写法：

- `ReviewUseCase.Pass`
- `OpenFlowUseCase.Commit`

### 命名理解

- `UseCase` 名称优先表达“这个业务对象负责什么流程”
- 方法名再表达具体动作
- 如果 `UseCase` 自身已经是动作词，通常说明领域对象还没有被识别清楚
- 在 `AccountUseCase.DeleteAccount(ctx, id uint32)` 这类聚合根语境中，主对象标识直接使用 `id`
- 只有 `DeleteByAccountID` 这类跨对象引用或关联资源场景才使用 `AccountID`

## Repo 命名

### 类型命名

| 类型 | 推荐命名 | 示例 |
|------|------|------|
| Repo 接口 | `{Domain}Repo` | `AccountRepo` |
| Repo 实现 | `{domain}Repo` | `accountRepo` |
| Repo 文件 | `{domain}.go` | `account.go` |

### 构造函数

```go
func NewAccountRepo(data *kit.Data) AccountRepo
```

### CRUD 方法参考

| 操作 | Repo 方法 | 示例 |
|------|------|------|
| 单条查询 | `Get{Domain}` | `GetAccount` |
| 列表查询 | `List{Domain}` | `ListAccount` |
| 分页查询 | `PageList{Domain}` | `PageListAccount` |
| 创建 | `Create{Domain}` | `CreateAccount` |
| 更新 | `Update{Domain}` | `UpdateAccount` |
| 删除 | `Delete{Domain}` | `DeleteAccount` |

### 命名理解

- Repo 命名优先表达“维护哪个领域对象”
- 不要把筛选条件、页面语义、协议语义塞进 Repo 类型名
- 如果 Repo 名称只对应某个接口动作，后续通常很难复用
- 在 `AccountRepo.GetAccount(ctx, id uint32)` 这类聚合根语境中，主对象标识直接使用 `id`
- 只有跨对象引用、从属资源删除或过滤条件才使用 `AccountID`

## Service 命名

### 类型命名

目录已经区分 `inner/admin/open` 后，Service 名一般不再重复目录语义。

更稳定的写法：

- `AccountService`
- `AccountFlowPageService`
- `AccountCollectService`

容易冗余的写法：

- `AccountInnerService`
- `OpenAccountAdminService`

### 构造函数

```go
func NewAccountService(uc *AccountUseCase) *AccountService
```

### 命名理解

- Service 更像协议入口，因此名称通常围绕领域对象
- 协议侧的 `admin/open/inner` 等差异，通常由目录和 proto side 表达，不必再重复进类型名

## 通用命名矩阵

| 层 | 推荐命名 | 示例 |
|------|------|------|
| 聚合/实体 | `{Domain}` | `Account`, `AccountFlowPage` |
| Repo | `{Domain}Repo` | `AccountRepo` |
| UseCase | `{Domain}UseCase` | `AccountUseCase` |
| Service | `{Domain}Service` | `AccountService` |
| Request | `{Operate}{Domain}Request` | `GetAccountRequest` |
| Reply | `{Operate}{Domain}Reply` | `GetAccountReply` |
| Event | `{Domain}{Action}` | `AccountCreated` |

## 代码示例参考

```go
type Account struct{}

type AccountRepo interface {
    GetAccount(ctx context.Context, id uint32) (*Account, error)
    ListAccount(ctx context.Context, filter *AccountFilter) ([]*Account, error)
}

type AccountUseCase struct {
    repo AccountRepo
}

func (u *AccountUseCase) Review(ctx context.Context, id uint32) error {
    _, err := u.repo.GetAccount(ctx, id)
    return err
}
```

```proto
message GetAccountRequest {
  uint32 id = 1;
}

message GetAccountReply {
  AccountInfo account_info = 1;
}
```

## Good Example

- `AccountRepo`
- `AccountFlowPageUseCase`
- `AccountCollectRelation`
- `GetAccountRequest`
- `ACCOUNT_OPEN_STATUS_PENDING`

## 常见坑

- 同一概念在不同层使用不同词
- 复用动作词作为实体或 Repo 名称
- 在 RPC 或方法名里携带筛选条件
- 先按页面或接口动作命名，后续再倒推领域模型
- 已有聚合能表达语义时，又新增一个近义 DTO 或结构体

## 相关 Rule

- `../rules/naming-rule.md`
- `../rules/aggregate-rule.md`

## 相关 Reference

- `./aggregate-spec.md`
- `./repo-spec.md`
- `./usecase-spec.md`
