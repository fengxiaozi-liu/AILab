# Comment Spec

## 注释位置决策

| 位置 | 是否需要注释 |
|------|------------|
| 导出类型（struct、interface） | MUST，表达业务角色和能力 |
| 导出函数/方法 | MUST，用"对象/能力 + 作用"式注释 |
| Ent schema 表/字段 | MUST，表达业务含义 |
| proto message/field | 视情况，不加说明性注释和装饰线 |
| 私有简单赋值 | 不需要 |
| 关键分支/阶段切换 | 允许短注释 |

---

## 导出类型注释

```go
// ✅ 表达业务角色，一行简洁说明
// AccountUseCase 账户聚合根能力
type AccountUseCase struct { ... }

// ✅ 方法注释表达动作语义
// Prepare 创建或刷新开户主体。
func (u *AccountUseCase) Prepare(ctx context.Context, userCode string) (*Account, error) { ... }

// ❌ 重复函数名，无业务含义
// GetAccount gets the account.
func (u *AccountUseCase) GetAccount(...) { ... }

// ❌ 空泛注释
// 处理数据
func (r *accountRepo) ParseFilter(...) { ... }
```

---

## Ent schema 注释

```go
// ✅ 表级注释 + 字段业务含义注释
func (Account) Fields() []ent.Field {
    return []ent.Field{
        field.Uint32("id").Comment("primary id"),
        field.String("user_code").MaxLen(64).Default("").Comment("user code"),
        field.Uint8("open_status").Default(0).Comment("open status 1=init 2=filling"),
        field.Uint32("create_time").Comment("create timestamp"),
    }
}

func (Account) Annotations() []schema.Annotation {
    return []schema.Annotation{schema.Comment("account table")}
}

// ❌ 字段无注释
field.Uint8("open_status").Default(0)
```

---

## proto 注释禁止项

```proto
// ❌ 加说明性注释
// 账户服务接口，包含创建、查询、审核功能
service AccountService {
  rpc GetAccount(GetAccountRequest) returns (GetAccountReply);
}

// ❌ 加装饰线
// ========================
// Account RPC
// ========================

// ✅ 不加注释，或只加极简字段说明（可选）
service AccountService {
  rpc GetAccount(GetAccountRequest) returns (GetAccountReply);
  rpc ListAccount(ListAccountRequest) returns (ListAccountReply);
}
```

---

## 阶段性注释

```go
// ✅ 只在明显阶段边界加短注释
func (u *AccountFlowPageUseCase) CommitPage(ctx context.Context, req *CommitPageInput) error {
    // ----- 1. 权限校验 -----
    if err := u.checkPermission(ctx, req); err != nil { return err }

    // ----- 2. 构建 store -----
    store, err := u.buildPageStore(ctx, req)
    if err != nil { return err }

    // ----- 3. 提交 -----
    return u.doCommit(ctx, store)
}

// ❌ 每行都加注释，正文被注释淹没
// 获取 account
account, err := u.repo.GetAccount(ctx, id)
// 检查 err
if err != nil {
    // 返回 err
    return err
}
```

---

## 组合场景

```go
// ProviderSet + Service + UseCase + Repo 完整注释示例
// ProviderSet 服务提供集合
var ProviderSet = wire.NewSet(NewAccountService)

// AccountService 实现 admin 侧账户服务接口。
type AccountService struct {
    v1.UnimplementedAccountServer
    uc *biz.AccountUseCase
}

// AccountUseCase 账户聚合根能力
type AccountUseCase struct { ... }

// Review 执行账户审核流程。
func (u *AccountUseCase) Review(ctx context.Context, id uint32, action openenum.ReviewAction) error { ... }
```

---

## 常见错误模式

```go
// ❌ 过期实现细节长期保留
// TODO: 临时方案，等接口对齐后删除（2022-01-01）

// ❌ 在生成物上手改注释（ent/account.go、wire_gen.go）

// ❌ 注释与现有命名不一致
// GetUser 返回账户信息  ← 函数名叫 GetAccount，注释说 User

// ❌ proto 加分隔线装饰
// ==================== Account ====================
```
