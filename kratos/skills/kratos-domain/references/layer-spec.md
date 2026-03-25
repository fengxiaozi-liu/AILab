# Layer Reference

## 约束先看
必须遵守：

- Service 只做协议适配与参数转换
- UseCase 只做业务编排、事务边界、权限和状态流转
- Repo 只做查询组织、relation 装配和数据访问边界
- relation 统一收口在 Repo，Service 和 UseCase 不补查
- 禁止越层调用、循环依赖、职责下沉错位

## 使用说明

说明 Kratos 项目中 `biz`、`data`、`service`、`server` 四层通常分别承担什么工作，以及常见调用路径如何组织。

## 常见场景

- 新增模块或移动代码文件
- 判断逻辑应放在 UseCase、Repo 还是 Service
- 分析跨层依赖和代码落位

## 使用边界

分层不是为了目录整齐，而是让协议变化、业务变化和数据变化分开演进。

- `service/server` 更接近接入和协议
- `biz/data` 更接近领域与持久化
- 层边界清楚后，改接口时不会顺手把业务逻辑带进接入层，改 Repo 时也不会反向污染 UseCase

## 分层职责表

| 层 | 应承担职责 | 不应承担职责 |
|------|------|------|
| `server` | 注册 gRPC/HTTP 服务与中间件 | 承载业务编排 |
| `service` | 协议适配、参数转换、显式入参准备 | 补查 relation、状态流转、复杂业务决策 |
| `biz/usecase` | 业务编排、事务边界、权限校验、状态流转 | 直接写 DB 查询细节、拼协议 DTO |
| `data/repo` | 查询组织、relation 装配、DB/远程依赖访问 | 承担整聚合更新编排、把 relation 补查回抛给上层 |

## 实施提示

- 先判断当前改动属于协议适配、业务编排还是数据装配
- 再决定代码应落在 `service/server`、`biz` 还是 `data`
- 如果一个函数同时依赖协议细节和数据装配细节，通常值得重新拆层

## 推荐结构

- `service/`：协议适配、参数转换、显式入参准备
- `biz/`：业务编排、事务、权限、状态流转
- `data/`：DB、缓存、远程依赖、relation 装配
- `server/`：注册 gRPC/HTTP 服务与中间件

## 典型调用路径

```text
Request -> Service -> UseCase -> Repo -> DB/Depend
                              -> EventBus/Tx
```

## 示例

- Service 只做请求解析与响应转换
- UseCase 负责决定是否需要事务和 relation
- Repo 负责查询、装配和依赖调用

## 常见坑

- 在 Service 中直接补查 relation
- 在 Repo 中处理状态流转
- 在 UseCase 中拼接协议层 DTO
