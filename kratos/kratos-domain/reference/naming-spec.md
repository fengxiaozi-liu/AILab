# Naming Spec

## 文件命名规则

| 文件类型 | 命名格式 | 示例 |
|---------|---------|------|
| 领域对象 + UseCase 接口 | `{aggregate}.go` | `account.go` |
| UseCase 实现 | `{aggregate}_usecase.go` | `account_usecase.go` |
| Repo 实现 | `{aggregate}.go` (data 包) | `data/account.go` |
| Service 实现 | `{aggregate}_service.go` | `account_service.go` |
| Proto 文件 | `{aggregate}.proto` | `account.proto` |

```go
// ✅ 聚合根名统一用领域概念
// biz/account.go  data/account.go  service/account_service.go

// ❌ 按接口动作命名文件
// biz/get_account.go  biz/review_account.go
// data/open_account_repo.go
```

---

## RPC 方法命名规则

| 前缀 | 语义 | 返回 |
|------|------|------|
| `Get` | 单对象查询 | 单个聚合根 |
| `List` | 全量/条件列表 | 列表 |
| `PageList` | 分页列表 | 分页+列表 |
| `Create` | 新建 | 主键或对象 |
| `Update` | 修改 | 空 or 对象 |
| `Delete` | 删除 | 空 |

```go
// ✅ 方法语义清晰
rpc GetAccount(GetAccountRequest) returns (GetAccountReply)
rpc ListAccount(ListAccountRequest) returns (ListAccountReply)
rpc PageListAccount(PageListAccountRequest) returns (PageListAccountReply)

// ❌ 后缀行为词或冗余 Detail
rpc GetAccountDetail(...)   // ❌ Detail 属于展示层概念
rpc QueryAccounts(...)      // ❌ Query 与 List 语义重叠
rpc FetchAccount(...)       // ❌ Fetch 不是标准前缀
```

---

## 字段命名规则

```go
// ✅ 聚合根上下文内省略主语
type Account struct {
    ID         uint32 `json:"id"`       // ✅ 不写 AccountID
    UserCode   string `json:"user_code"` // ✅ 不写 AccountUserCode
}

// ❌ 冗余主语
type Account struct {
    AccountID       uint32 `json:"account_id"`     // ❌
    AccountStatus   string `json:"account_status"` // ❌
}

// ✅ 关联 ID 使用被关联对象标识
type AccountFlowPage struct {
    AccountID        uint32 `json:"account_id"`        // ✅ 外键显式带主语
    CheckAdminUserID uint32 `json:"check_admin_user_id"` // ✅ 关联 ID
}
```

---

## Go 标识符命名

```go
// ✅ UseCase/Repo/Service 接口名包含角色
type AccountUseCase interface { ... }
type AccountRepo interface { ... }

// ✅ 实现结构体用小写 + Repo/UseCase
type accountUseCase struct { ... }
type accountRepo struct { ... }

// ❌ 无角色后缀
type Account interface { ... }   // 不知是 UseCase 还是 Repo
type AccountImpl struct { ... }  // 应用 accountUseCase
```

---

## 组合场景

```go
// 完整：proto 方法 → Service → UseCase → Repo（命名一致）

// proto
rpc GetAccount(GetAccountRequest) returns (GetAccountReply) {}
rpc PageListAccount(PageListAccountRequest) returns (PageListAccountReply) {}

// Service
type AccountService struct { uc *biz.AccountUseCase }

// UseCase 接口（biz/account.go）
type AccountRepo interface {
    GetAccount(ctx context.Context, id uint32, opts ...filter.Option) (*Account, error)
    PageListAccount(ctx context.Context, pg *page.Page, opts ...filter.Option) ([]*Account, error)
}

// Data 实现 (data/account.go)
type accountRepo struct { data *Data }
func (r *accountRepo) GetAccount(ctx context.Context, id uint32, opts ...filter.Option) (*biz.Account, error) { ... }
func (r *accountRepo) PageListAccount(ctx context.Context, pg *page.Page, opts ...filter.Option) ([]*biz.Account, error) { ... }
```

---

## 常见错误模式

```go
// ❌ 动作词命名聚合根文件
// biz/open_flow.go → 应为 biz/account.go
// biz/review_handler.go → 应为 biz/account_usecase.go

// ❌ Get/List 混用
rpc ListAccountById(...)  // ❌ 按 ID 查是 Get
rpc GetAllAccounts(...)   // ❌ 全量是 List

// ❌ 聚合根字段带自身主语
type Account struct { AccountID uint32 }  // ❌ 应为 ID
```
