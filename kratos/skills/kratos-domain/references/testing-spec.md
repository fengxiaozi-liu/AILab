# Testing Reference

## 约束先看
必须遵守：

- 业务行为变化必须补测试
- 测试优先覆盖边界、失败路径和状态流转，不只覆盖 happy path
- UseCase、Repo、Service 的测试职责要分层
- 注释里的约束如果重要，应转成可执行测试

## 使用说明

说明 Kratos 业务代码中如何组织 UseCase、Repo、Service 的测试，以及如何覆盖关键边界场景。

## 常见场景

- 业务变更后补测试
- 新增回归测试
- 设计用例覆盖范围

## 分层关注点

| 层 | 测试重点 |
|------|------|
| UseCase | 业务分支、事务、状态流转、权限 |
| Repo | 查询条件、relation 装配、分页、远程补全 |
| Service | 协议转换、参数透传、响应映射 |

## 实施提示

- 先列正常路径、边界路径、失败路径
- 再判断这些场景更适合放在 UseCase、Repo 还是 Service 验证
- 如果某个约束只存在于注释里，优先转成测试用例

## 示例

### 文件组织

```text
account.go
account_test.go
```

### UseCase 测试模板

```go
func TestAccountUseCase_Submit(t *testing.T) {
    repo := &mockAccountRepo{
        getAccountFn: func(ctx context.Context, id uint32, opts ...filter.Option) (*biz.Account, error) {
            return &biz.Account{ID: id, Status: openenum.AccountStatusInit}, nil
        },
    }

    uc := NewAccountUseCase(repo, nil, nil, nil)
    err := uc.Submit(context.Background(), &SubmitAccountRequest{ID: 1})
    require.NoError(t, err)
}
```

### Repo 测试关注点

- 查询条件是否生效
- relation 是否装配完整
- 分页是否稳定排序
- 远程 relation 是否批量补齐

## 常见坑

- 只补 happy path
- 把 Service 测试写成业务测试
- 业务状态变化没有回归测试
