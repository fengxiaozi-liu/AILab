# Layer Reference

## 这个主题解决什么问题

说明 Kratos 项目中 biz、data、service、server 四层通常分别承担什么工作，以及常见调用路径如何组织。

## 适用场景

- 新增模块或移动代码文件
- 判断逻辑应放在 UseCase、Repo 还是 Service
- 分析跨层依赖和代码落位

## 设计意图

分层参考不是为了目录整齐，而是让协议变化、业务变化和数据变化分开演进。

- `service/server` 更接近接入和协议，`biz/data` 更接近领域与持久化。
- 先理解每层“为什么存在”，后续落代码时更容易把逻辑放到稳定位置。
- 层边界清楚后，改接口时不会顺手把业务逻辑带进接入层，改 Repo 时也不会反向污染 UseCase。

## 实施提示

- 先判断当前改动属于协议适配、业务编排还是数据装配。
- 再决定代码应该落在 `service/server`、`biz` 还是 `data`。
- 如果一个函数同时依赖协议细节和数据装配细节，通常值得重新拆层。

## 推荐结构

- `service/`：协议适配、参数转换、显式入参准备
- `biz/`：业务编排、事务、权限、状态流转
- `data/`：DB、缓存、远程依赖、relation 装配
- `server/`：注册 gRPC/HTTP 服务与中间件

## 典型实现方式

```text
Request -> Service -> UseCase -> Repo -> DB/Depend
                              -> EventBus/Tx
```

## 标准模板

```go
func (s *AccountService) Get(ctx context.Context, req *v1.GetAccountRequest) (*v1.GetAccountReply, error) {
    info, err := s.uc.GetAccount(ctx, req.Id)
    if err != nil {
        return nil, err
    }
    return convertToReply(info), nil
}
```

## Good Example

- Service 只做请求解析与响应转换
- UseCase 负责决定是否需要事务和 relation
- Repo 负责查询、装配和依赖调用

## 常见坑

- 在 Service 中直接补查 relation
- 在 Repo 中处理状态流转
- 在 UseCase 中拼接协议层 DTO

## 相关 Rule

- `../rules/layer-rule.md`
- `../rules/usecase-rule.md`
- `../rules/repo-rule.md`
