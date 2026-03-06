# UseCase Reference

## 这个主题解决什么问题
说明 UseCase 如何围绕业务流程组织编排、事务、权限和 relation 需求，以及应用层输入输出对象如何作为聚合根的应用层投影存在。

## 适用场景

- 新增或修改 UseCase 方法
- 设计事务边界和状态流转
- 设计应用层输入输出对象
- 判断 relation 需求由谁决定

## 设计意图

UseCase 负责业务编排，不负责底层装配细节。

- UseCase 是应用层，不是协议层。
- 应用层输入输出对象与 `proto message` 是聚合根的同源平级产物。
- 应用层输入输出对象负责表达业务编排需要的输入输出。
- `proto message` 负责表达协议层契约。

## 实施提示

- 先按业务步骤描述流程，再决定哪些步骤放进事务。
- 先确定需要哪些 relation，再通过 `opts` 把需求传给 Repo。
- 如果应用层输入输出对象的字段设计只是机械复制 `proto`，通常说明应用层表达还没有围绕聚合根收敛。

## 推荐结构

UseCase 负责：

- 业务编排
- 权限校验
- 状态流转
- 事务边界
- relation 需求声明

UseCase 不直接实现 relation 查询细节，而是通过 `opts ...filter.Option` 影响 Repo 行为。

## 应用层输入输出对象的定位

常见顺序：

1. 先识别实体与聚合根
2. 再定义应用层输入输出对象
3. 再定义对应的 `proto message`
4. 最后在 service 层做应用层 DTO 与协议层 message 的转换

常见理解：

- 应用层输入输出对象是聚合根的应用层投影
- `proto message` 是聚合根的协议层投影
- 两者平级，不互相从属
- 如无歧义，应用层对象直接使用聚合根名称，不额外追加 `DTO` 后缀

## 标准模板

### 只读查询

```go
func (u *AccountUseCase) GetAccount(ctx context.Context, id uint32, opts ...filter.Option) (*Account, error) {
    return u.accountRepo.GetAccount(ctx, id, opts...)
}
```

### 业务强依赖 relation

```go
func (u *AccountUseCase) GetAccountWithRequiredRelation(ctx context.Context, id uint32) (*Account, error) {
    opts := []filter.Option{
        filter.WithRelation(openenum.AccountCollectRelation),
    }
    return u.accountRepo.GetAccount(ctx, id, opts...)
}
```

### 写操作与事务

```go
func (u *AccountUseCase) Review(ctx context.Context, id uint32, action openenum.ReviewAction) error {
    return u.tx.InTx(ctx, func(ctx context.Context) error {
        account, err := u.accountRepo.GetAccount(ctx, id)
        if err != nil {
            return err
        }

        account.Status = openenum.AccountStatusReviewed
        return u.accountRepo.UpdateAccount(ctx, account)
    })
}
```

## 显式入参模式

```go
func (s *AccountService) GetOpenStatus(ctx context.Context, req *v1.GetOpenStatusRequest) (*v1.GetOpenStatusReply, error) {
    userCode := strconv.FormatUint(uint64(metadata.GetViewerID(ctx)), 10)
    status, err := s.accountUseCase.GetOpenStatus(ctx, userCode)
    if err != nil {
        return nil, err
    }
    return &v1.GetOpenStatusReply{OpenStatus: uint32(status)}, nil
}
```

## 代码示例参考

```go
type Account struct {
    ID                  uint32                 `json:"id"`
    Status              openenum.AccountStatus `json:"status"`
    CreateTime          uint32                 `json:"create_time"`
    FirstCheckUserInfo  *adminbiz.AdminUser    `json:"first_check_user_info"`
    AccountCollectInfo  *AccountCollect        `json:"account_collect_info"`
    AccountFlowPageList []*AccountFlowPage     `json:"account_flow_page_list"`
}

func (u *AccountUseCase) GetAccountDetail(ctx context.Context, id uint32) (*Account, error) {
    opts := []filter.Option{
        filter.WithRelation(openenum.AccountCollectRelation),
        filter.WithRelation(openenum.AccountCheckUserRelation),
    }
    return u.accountRepo.GetAccount(ctx, id, opts...)
}
```

## 常见坑

- 在 UseCase 里直接补查 reviewer、customer、collect 等 relation
- 通过 `context` 偷读身份而不是显式入参
- 应用层对象和 `proto message` 同时各维护一套近义结构
- 让调用方通过判空去猜业务状态

## 相关 Rule

- `../rules/usecase-rule.md`
- `../rules/layer-rule.md`

## 相关 Reference

- `./aggregate-spec.md`
- `./repo-spec.md`
