# UseCase 参考

## 先看约束

- UseCase 负责业务编排、事务边界、权限校验和状态流转。
- relation 需求由 UseCase 声明，由 Repo 实现。
- UseCase 不得直接编写 DB 查询细节，也不得手工组装 relation。
- 对于查询、列表、统计类场景，业务 filter 应由 Service/Handler 构建后传入 UseCase。
- UseCase 不得基于 proto/request 原始字段临时拼装查询 filter。
- 在引入专用 filter 之前，优先复用已有的业务 filter 类型。
- 不要用隐式上下文替代显式业务入参。

## 适用场景

- 新增或修改 UseCase 方法。
- 设计事务边界或状态流转。
- 设计应用层输入输出对象。
- 在查询链路中澄清 Service、UseCase、Repo 的职责边界。

## UseCase 职责

- 业务编排。
- 权限校验。
- 状态流转。
- 事务边界。
- 接收业务 filter 并协调 repo 调用。

## 实施指引

- 先描述业务流程，再判断哪些步骤需要放进事务。
- 先确定 relation 需求，再通过 `opts` 传给 Repo。
- 在查询场景中，由 Service 将协议输入映射为业务 filter，UseCase 直接接收该 filter。
- 当已有业务 filter 可以清晰表达查询语义时，优先复用。
- 只有在复用会让语义变得模糊或别扭时，才新增独立 filter。

## 输入与输出对象

- 应用层对象表达编排所需的输入与输出。
- `proto message` 表达协议层契约。
- 二者可以共享同一领域语义，但默认不应视为同一个对象。

## 示例

### 查询 filter 透传

```go
func (u *AccountUseCase) PageList(ctx context.Context, f *AccountFilter) ([]*Account, int, error) {
    return u.accountRepo.PageListAccount(ctx, f)
}
```


## 常见坑

- 在 UseCase 内直接编写 DB 查询细节。
- 在 UseCase 内根据 request/proto 原始字段拼装查询 filter。
- 将协议映射逻辑吸入 UseCase。
- 对已由下层负责的通用校验再做一层重复防御。
