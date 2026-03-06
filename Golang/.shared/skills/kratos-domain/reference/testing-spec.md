# Testing Reference

## 这个主题解决什么问题

说明 Kratos 业务代码中如何组织 UseCase、Repo、Service 的测试，以及如何覆盖边界场景。

## 适用场景

- 业务变更后补测试
- 新增回归测试
- 设计用例覆盖范围

## 设计意图

测试参考的作用是把“业务变化”映射成“可验证行为”，帮助在改动后补足真正能保护回归的测试。

- 测试组织清楚后，后续重构时更容易判断哪些行为不能被意外改掉。
- 业务边界、错误语义和状态流转如果没有测试样例，重构时更容易漏掉隐含约束。
- 把测试视为行为文档，也能反向帮助理解聚合、UseCase 和 Repo 的职责边界。

## 实施提示

- 先列出正常路径、边界路径和失败路径。
- 再判断这些场景更适合放在 UseCase、Repo 还是 Service 层验证。
- 如果某个约束只存在于注释里，优先把它转成一个可执行测试场景。

## 推荐结构

- UseCase 测试关注业务分支、事务、状态流转
- Repo 测试关注查询条件、relation 装配、分页和远程补全
- Service 测试关注协议转换和参数透传

## 典型实现方式

```text
account.go
account_test.go
```

## 项目通用测试骨架

```go
type mockAccountRepo struct {
    getAccountFn func(context.Context, uint32) (*Account, error)
}

func (m *mockAccountRepo) GetAccount(ctx context.Context, id uint32, opts ...filter.Option) (*Account, error) {
    if m.getAccountFn == nil {
        return nil, openerror.ErrorAccountNotFound(ctx)
    }
    return m.getAccountFn(ctx, id)
}

func setupTestDeps(t *testing.T) *testDeps {
    return &testDeps{
        accountRepo: &mockAccountRepo{},
        flowPageRp:  &mockAccountFlowPageRepo{},
    }
}
```

### UseCase 用例示例

```text
- 正常路径
- 权限不足
- 状态流转失败
- 必填 relation 缺失
```

## UseCase 测试示例

```go
func TestReview_RejectPagesEmpty(t *testing.T) {
    deps := setupTestDeps(t)
    ctx := ctxWithViewer(999)
    uc := deps.newUseCase()

    deps.accountRepo.getAccountFn = func(_ context.Context, _ uint32) (*Account, error) {
        return &Account{ID: 1, OpenStatus: openenum.AccountOpenStatusFirstChecking}, nil
    }

    err := uc.Review(ctx, 1, openenum.ReviewActionReject, nil)
    if err == nil {
        t.Fatal("expected error when reject_pages is empty")
    }
}
```

### Repo 用例示例

```text
- parseFilter 排序和分页
- queryConfig 本地 relation
- serviceRelation 批量补全
- not found 语义
```

## 行为断言示例

```go
func TestPageListAccount_PassesFilterAndCollectRelation(t *testing.T) {
    deps := setupTestDeps(t)
    uc := deps.newUseCase()

    deps.accountRepo.pageListAccountFn = func(_ context.Context, f *AccountFilter, opts ...filter.Option) ([]*Account, int, error) {
        cfg := filter.NewConfig(opts...)
        if _, ok := cfg.Relations[openenum.AccountCollectRelation]; !ok {
            t.Fatal("expected account_collect relation preload option")
        }
        return []*Account{{ID: 1}}, 1, nil
    }

    list, count, err := uc.account.PageListAccount(
        context.Background(),
        &AccountFilter{},
        filter.WithRelation(openenum.AccountCollectRelation),
    )
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    if count != 1 || len(list) != 1 {
        t.Fatalf("unexpected result: count=%d len=%d", count, len(list))
    }
}
```

## 常见坑

- 只测 happy path
- 只测 Service，不测 UseCase 或 Repo 核心逻辑
- relation 变更后没有回归测试

## 相关 Rule

- `../rules/testing-rule.md`
