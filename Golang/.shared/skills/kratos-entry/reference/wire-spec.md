# Wire Reference

## 这个主题解决什么问题
说明 provider 集合和 Wire 依赖注入如何组织，以及它与聚合根、`usecase DTO`、proto、service 之间的联动顺序。

## 适用场景

- 新增 provider
- 调整模块依赖
- 排查 wire 生成失败

## 设计意图

Wire 是注入层收口，不是领域建模或协议设计入口。

- Wire 不负责设计 proto，也不负责定义 `usecase DTO`。
- Wire 依赖聚合根、应用层 DTO、proto/service 契约和 provider 的稳定结果。
- 更适合作为最后收口的文件，而不是最先创建的文件。

## 实施提示

- 先找清楚新增依赖属于 data、biz 还是 service/provider 侧。
- 沿现有 provider set 增量扩展，而不是另起一套注入入口。
- 如果 `usecase DTO`、proto 或 service 还没稳定，就先不要急着改 wire。

## 推荐结构

- 每层维护自己的 `ProviderSet`
- 构造函数显式声明依赖
- 配置对象优先按父配置注入，再在构造函数中读取子字段

## 与 Proto 的联动顺序

更稳定的顺序通常是：

1. 先完成实体识别与聚合根建模
2. 再平级完成 `usecase DTO` 与 proto/service 契约
3. 再补 repo / usecase / service 的 provider
4. 最后更新 `ProviderSet` 和 wire 生成

## 文件创建理解

常见顺序是：

- 先有 domain 文件
- 再有 `usecase DTO` / proto / service 文件
- 再有 repo / usecase / service provider
- 最后才有 wire 入口更新

## 标准模板

```go
var ProviderSet = wire.NewSet(
    NewAccountRepo,
    NewAccountUseCase,
    NewAccountService,
)
```

## Good Example

```go
func NewAccountUseCase(repo biz.AccountRepo, tx biz.Tx, logger log.Logger) *biz.AccountUseCase {
    return &biz.AccountUseCase{repo: repo, tx: tx, log: log.NewHelper(logger)}
}
```

## 代码示例参考

```go
var ProviderSet = wire.NewSet(
    data.NewData,
    data.NewAccountRepo,
    biz.NewAccountUseCase,
    service.NewAccountService,
)
```

## 项目通用入口示例

```go
func wireApp(
    serverConf *conf.Server,
    dataConf *conf.Data,
    registryConf *conf.Registry,
    remoteConf *conf.RemoteConfig,
    thirdParty *conf.ThirdParty,
    logger log.Logger,
) (*kratos.App, func(), error) {
    panic(wire.Build(
        server.ProviderSet,
        data.ProviderSet,
        biz.ProviderSet,
        service.ProviderSet,
        listener.ProviderSet,
        consumer.ProviderSet,
        crontab.ProviderSet,
        newApp,
    ))
}
```

## 分层 ProviderSet 示例

```go
var ProviderSet = wire.NewSet(
    NewAccountRepo,
    NewAccountFlowPageRepo,
)

var ProviderSet = wire.NewSet(
    NewAccountUseCase,
    NewAccountFlowPageUseCase,
)

var ProviderSet = wire.NewSet(
    NewAccountService,
    NewAccountFlowPageService,
)
```

## 常见坑

- proto / `usecase DTO` 还没定型就先改 wire
- provider 分散在多个无关文件里，难以检索
- 构造函数依赖过多，模块边界不清
- 注入层先改了，但契约和实现还没收敛

## 相关 Rule

- `../rules/wire-rule.md`
- `../rules/codegen-rule.md`
