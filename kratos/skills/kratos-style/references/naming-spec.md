# Naming Spec

## 领域词根统一

| 条件 | 做法 |
|------|------|
| 同一聚合根跨层出现 | 使用同一领域词根 |
| 只是接口动作、页面形态、临时视图 | 放到方法名或 DTO 用途后缀，不作为核心对象名 |

```go
// ✅ 同一词根跨层统一
type Account struct{}
type AccountUseCase struct{}
type AccountRepo interface{}
type AccountService struct{}
```

```go
// ❌ 每层各发明一套近义词
type AccountEntity struct{}
type UserAccountUseCase struct{}
type OpenPageRepo interface{}
type AccountPageService struct{}
```

```go
// ⚠️ 视图对象可以带用途后缀，但核心词根仍应稳定
type AccountDetail struct{}
type AccountListItem struct{}
```

---

## 对象名与方法名分工

| 条件 | 做法 |
|------|------|
| 命名对象 | 用领域名词 |
| 命名行为 | 用动作词 |

```go
// ✅ 对象名表达领域，方法名表达动作
type AccountUseCase struct{}

func (u *AccountUseCase) GetAccount(ctx context.Context, id uint32) (*Account, error) {
    ...
}
```

```go
// ❌ 对象名带动作，方法名失去语义
type GetAccountData struct{}
type HandleAccountRepo struct{}

func (r *HandleAccountRepo) Do(ctx context.Context) error {
    ...
}
```

---

## 文件名与目录下命名

| 条件 | 做法 |
|------|------|
| `biz`、`data`、`service` 下核心文件 | 使用简短 `snake_case`，围绕聚合根 |
| repo 实现文件 | 使用 `{entity}.go` |

```text
// ✅
internal/biz/{业务文件夹}/account.go
internal/data/{业务文件夹}/account.go
internal/service/account.go
api/{业务文件夹}/account.proto
```

```text
// ❌
internal/biz/{业务文件夹}/get_account_page_data.go
internal/data/{业务文件夹}/handle_account_repo.go
api/get_account.proto
```

```text
// ⚠️ 只在子领域确实独立时追加限定词
internal/biz/{业务文件夹}/account_risk.go
```

---

## Proto 与服务命名

| 条件 | 做法 |
|------|------|
| proto 文件 | 用聚合根或稳定业务主题命名 |
| message | `{Domain}{Purpose}` |
| service | `{Domain}Service` |
| rpc | 用动作表达行为 |

```proto
// ✅ 围绕聚合根组织
syntax = "proto3";

service AccountService {
  rpc GetAccount(GetAccountRequest) returns (GetAccountReply);
}

message AccountInfo {}
message GetAccountRequest { uint32 id = 1; }
message GetAccountReply { AccountInfo account = 1; }
```

```proto
// ❌ 围绕页面和处理器命名
service AccountPageHandler {
  rpc Handle(GetAccountPageRequest) returns (GetAccountPageReply);
}

message GetAccountPageRow {}
```

```proto
// ⚠️ 同一 proto 文件中允许多个动作 RPC，但 service 仍应围绕聚合根
service AccountService {
  rpc GetAccount(GetAccountRequest) returns (GetAccountReply);
  rpc FreezeAccount(FreezeAccountRequest) returns (FreezeAccountReply);
}
```

---

## 临时结构与边界命名

| 条件 | 做法 |
|------|------|
| 临时 DTO / VO 仅服务某个场景 | 用 `{Domain}{Purpose}`，不要上升为核心词根 |
| 名称无法映射回聚合根 | 先收敛建模，再命名 |

```go
// ✅ 临时结构保留领域词根
type AccountExportRow struct{}
type AccountRiskSnapshot struct{}
```

```go
// ❌ 页面/接口反推核心领域
type OpenAccountPageRow struct{}
type HandleAccountResult struct{}
```

---

## 组合场景

```go
// internal/biz/account.go
type Account struct{}

type AccountRepo interface {
    GetAccount(ctx context.Context, id uint32) (*Account, error)
}

type AccountUseCase struct {
    repo AccountRepo
}

func (u *AccountUseCase) GetAccount(ctx context.Context, id uint32) (*Account, error) {
    return u.repo.GetAccount(ctx, id)
}
```

```proto
// api/account.proto
service AccountService {
  rpc GetAccount(GetAccountRequest) returns (GetAccountReply);
}

message AccountInfo {}
message GetAccountRequest { uint32 id = 1; }
message GetAccountReply { AccountInfo account = 1; }
```

这个组合场景同时满足：

- 文件、对象、repo、usecase、service、proto 使用统一词根 `Account`
- 对象名是领域名词，动作放在 `GetAccount`
- proto/service 没有被页面名或 handler 语义污染

---

## 常见错误模式

```go
// ❌ 对象名带动作
type GetAccountData struct{}
```

```go
// ❌ 每层用不同词根
type UserAccountUseCase struct{}
type AccountRepo interface{}
type OpenAccountService struct{}
```

```proto
// ❌ service 按 handler 或页面命名
service AccountPageHandler {}
```

```text
// ❌ 文件名围绕接口动作组织
get_account_page.go
```
