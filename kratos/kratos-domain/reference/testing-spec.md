# Testing Spec

## 测试范围决策

| 测试层 | 覆盖内容 | 不覆盖 |
|--------|---------|--------|
| UseCase Test | 状态流转、分支、事务顺序 | 具体 SQL 语句 |
| Repo Test | 四段式执行、not found 行为 | 业务状态校验 |
| Service Test | 入参转换、出参 proto 构建 | 业务逻辑分支 |

---

## Mock 函数字段模式

```go
// ✅ 用函数字段 mock，按方法独立注入
type mockAccountRepo struct {
    getAccountFn      func(ctx context.Context, id uint32, opts ...filter.Option) (*biz.Account, error)
    updateStatusFn    func(ctx context.Context, id uint32, status openenum.AccountOpenStatus) error
}

func (m *mockAccountRepo) GetAccount(ctx context.Context, id uint32, opts ...filter.Option) (*biz.Account, error) {
    return m.getAccountFn(ctx, id, opts...)
}
func (m *mockAccountRepo) UpdateStatus(ctx context.Context, id uint32, status openenum.AccountOpenStatus) error {
    return m.updateStatusFn(ctx, id, status)
}

// ❌ 全局 mock 变量（并发不安全，且不能按测试定制行为）
var mockAccount *biz.Account
func (m *mockAccountRepo) GetAccount(...) (*biz.Account, error) { return mockAccount, nil }
```

---

## 测试覆盖分支

```go
// ✅ 每个分支独立子测试
func TestAccountUseCase_Review(t *testing.T) {
    tests := []struct {
        name    string
        account *biz.Account
        pass    bool
        wantErr bool
    }{
        {"pass review", &biz.Account{ID: 1, OpenStatus: openenum.AccountOpenStatusPending}, true, false},
        {"reject review", &biz.Account{ID: 1, OpenStatus: openenum.AccountOpenStatusPending}, false, false},
        {"wrong status", &biz.Account{ID: 1, OpenStatus: openenum.AccountOpenStatusOpened}, true, true},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            repo := &mockAccountRepo{
                getAccountFn:   func(...) (*biz.Account, error) { return tt.account, nil },
                updateStatusFn: func(...) error { return nil },
            }
            uc := NewAccountUseCase(repo, ...)
            err := uc.Review(context.Background(), 1, tt.pass)
            if (err != nil) != tt.wantErr {
                t.Errorf("got err=%v, wantErr=%v", err, tt.wantErr)
            }
        })
    }
}

// ❌ 只测 happy path
func TestReview(t *testing.T) {
    uc.Review(ctx, 1, true)
    // ❌ 不测 wrong status、reject 分支
}
```

---

## 测试文件组织

```text
// ✅ 每个实现文件对应一个测试文件
biz/account_usecase.go → biz/account_usecase_test.go
data/account.go        → data/account_test.go

// ❌ 所有测试堆在一个文件
biz/all_usecase_test.go  // ❌
```

---

## Repo 测试（使用 enttest in-memory）

```go
// ✅ 使用 enttest 提供 in-memory DB
func TestAccountRepo_GetAccount(t *testing.T) {
    client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
    defer client.Close()

    repo := newAccountRepo(&Data{Db: client})
    // 先写入种子数据
    created, _ := client.Account(context.Background()).Create().
        SetUserCode("U001").
        SetStatus(openenum.AccountOpenStatusPending).
        Save(context.Background())

    got, err := repo.GetAccount(context.Background(), uint32(created.ID))
    assert.NoError(t, err)
    assert.Equal(t, "U001", got.UserCode)
}

// ✅ not found 返回具体业务错误
func TestAccountRepo_GetAccount_NotFound(t *testing.T) {
    ...
    _, err := repo.GetAccount(context.Background(), 9999)
    assert.ErrorContains(t, err, "ACCOUNT_NOT_FOUND")
}
```

---

## 组合场景

```go
// 完整 UseCase 测试：状态守卫 + 事务 + 函数字段 mock
func TestAccountUseCase_PassReview(t *testing.T) {
    updateCalled := false
    publishCalled := false

    repo := &mockAccountRepo{
        getAccountFn: func(ctx context.Context, id uint32, opts ...filter.Option) (*biz.Account, error) {
            return &biz.Account{ID: 1, OpenStatus: openenum.AccountOpenStatusPending}, nil
        },
        updateStatusFn: func(ctx context.Context, id uint32, status openenum.AccountOpenStatus) error {
            updateCalled = true
            assert.Equal(t, openenum.AccountOpenStatusOpened, status)
            return nil
        },
    }
    bus := &mockEventBus{
        publishFn: func(ctx context.Context, event *eventbus.Event) error {
            publishCalled = true
            return nil
        },
    }
    uc := NewAccountUseCase(repo, bus, ...)

    err := uc.PassReview(context.Background(), 1)
    assert.NoError(t, err)
    assert.True(t, updateCalled)
    assert.True(t, publishCalled)
}
```

---

## 常见错误模式

```go
// ❌ 测试仅 assert nil error，不验证副作用
err := uc.PassReview(ctx, 1)
assert.NoError(t, err)  // ❌ 未验证 updateStatus 被调用

// ❌ 测试之间共享状态
var sharedRepo = &mockAccountRepo{ ... }  // ❌ 并发测试时相互影响

// ❌ mock 无法覆盖错误分支（固定返回 nil）
func (m *mockAccountRepo) GetAccount(...) (*biz.Account, error) {
    return &biz.Account{}, nil  // ❌ 无法模拟 not found
}
```
