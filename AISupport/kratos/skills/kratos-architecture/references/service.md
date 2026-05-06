# Service

## 作用范围

本文用于说明 Kratos 接入层与应用层暴露相关知识，包括 `proto`、`internal/service`、`internal/server`、gateway、wire 装配与 codegen 联动。

当问题属于以下场景时，应优先查看本文：

- 判断 `proto` 文件应如何组织
- 判断 `service`、`server`、gateway 各自承担什么职责
- 判断新增接口时接入层应如何收敛
- 判断 `wire`、provider 与接入层的装配关系
- 判断 codegen 与源定义之间的关系

本文重点回答接入层如何稳定组织，不展开领域建模、聚合根识别、公共能力下沉与代码风格红线。

## 规则

- 接入层主线保持 `proto -> service -> server`，其中 `service` 只做协议适配并调用 `usecase`
- `proto` 是契约事实源，`server / gateway` 只做暴露与转发，不承载业务编排
- 对外 side 的业务 proto 不允许跨业务域直接互引
- 涉及 `wire / provider / codegen` 的改动，必须同步评估生成链路与接入层装配闭环
- 不要把领域规则、状态推进或持久化细节塞进 `service / server / gateway`
- `service`下只允许调用`usecase`边界，不允许直接调用repo，eventbus或其他基础设置组件
- `proto` 优先围绕聚合根、实体或稳定第三方对象边界组织，不按页面动作或阶段动作拆文件
- 一旦某个主题已经独立定义 `UseCase` 且该能力需要对外暴露，就必须建立对应的 `service` 与 `proto`
- 不要出现“已有独立 `UseCase` 对外提供能力，但仍长期借挂在无关 `service/proto`”的情况；协议主题应与 `UseCase` 主题对齐
- 若某个实体或主题不直接对外暴露能力，则不要求单独新增 `service/proto`
- 新增 `inner` 能力时，必须在 `internal/biz/depend/{business}/{entity}.go` 中定义接口，并在 `internal/data/depend/{business}/{entity}.go` 中实现
- `biz` 只能依赖 `depend` 抽象，不直接依赖 generated inner client，也不允许把 inner client 直接注入到 `service` 或 `biz`
- `queryConfig.relation` 必须由 `service` 根据请求语义显式构造，并以清晰的嵌套 relation 结构向下传递
- `biz / repo / data` 应消费上游传入的 relation 结构，`data` 只负责执行，不负责写死 relation 关系或补造整棵 relation 树

## 接入层主线

Kratos 接入层通常围绕以下主线展开：

```text
proto -> service -> server
```

其中：

- `proto` 负责契约定义
- `service` 负责协议适配与调用 `usecase`
- `server` 负责对外暴露与注册

若存在 gateway，则 gateway 仍属于接入层，只负责代理、参数映射、响应适配与聚合转发。

## `proto`

定位：契约事实源。

承载内容：

- `service`
- `message`
- `rpc`
- side 边界
- 共享协议结构

边界提示：

- `proto` 设计优先围绕稳定服务边界，而不是围绕页面动作临时拆文件
- `proto` 文件优先围绕聚合根、实体或稳定第三方对象边界组织，而不是围绕页面动作临时拆文件
- side 之间不直接互引业务 proto，共享协议应回到公共位置收敛
- 修改接口结构时，优先回到 `proto` 源定义，而不是手改生成物

### `proto` 文件组织

- 引入新的核心聚合根、稳定实体或稳定第三方对象边界时，新建对应 `proto` 文件
- 同一聚合根在不同 side 暴露接口时，在对应 side 目录下维护对应 `proto`
- 只是为已有聚合补动作、查询或稳定接口时，优先追加到现有聚合 `proto`
  前提：这些动作仍属于该聚合既有 `UseCase/service` 主题，而不是已经独立成新的 `UseCase` owner
- 若已经形成独立 `UseCase` owner 并对外暴露能力，则应新建并收敛到对应 `proto/service`
- 若某个实体或主题不直接对外暴露能力，则不要求单独新增 `service/proto`

### OpenAPI v3 文档注解

如果问题涉及以下内容，请转到 `openapi-v3-spec.md`：
- `openapi.v3.document / operation / schema / property` 注解覆盖
- 中文标题、标签、描述与 `operation_id` 约定
- 字段文档说明与 `validate.rules` 对齐
- OpenAPI v3 常见错误模式与完整示例

`service.md` 负责 proto / service / server / gateway 的分层与职责边界，
`openapi-v3-spec.md` 负责 OpenAPI v3 文档注解与展示规范。

示例：

```text
api/open/open/v1/account.proto
api/open/inner/v1/account.proto
api/base/business/v1/business.proto
```

### side 边界与引用规则

`proto` 的 side 边界必须稳定，避免业务协议在不同 side 之间直接耦合。

通常 side 包括：

- `open`
- `admin`
- `inner`
- `base`

引用规则如下：

- `open`、`admin` 以及其他对外 side 只允许在各自 side 内部引用业务 proto
- `open`、`admin` 以及其他对外 side 之间不允许直接互引业务 proto
- `inner` side 之间允许按需相互引入，用于服务间依赖契约组织
- `inner` proto 代表服务间依赖契约，不应被 `open`、`admin` 或其他对外 side 直接当作对外协议复用

允许的情况：

- 同一对外 side 内部围绕稳定业务域的引用
- `inner` side 之间围绕服务依赖的引用
- 引用 `api/base/*` 下的公共结构

不允许的情况：

- `open` 直接引用 `admin` 业务 proto
- `admin` 直接引用 `open` 业务 proto
- 其他对外 side 直接引用 `open` 或 `admin` 业务 proto
- `open`、`admin` 或其他对外 side 直接把 `inner` proto 当对外协议使用
- 为了复用少量字段，在 side 之间直接互引业务 proto

正例：

```proto
import "api/base/business/v1/business.proto";
```

```proto
import "api/user/inner/v1/user.proto";
```

反例：

```proto
import "api/open/open/v1/account.proto";
import "api/user/inner/v1/user.proto";
```

### `inner` proto 联动

当新增 `api/*/inner/v1/*.proto` 时，不能只停留在协议层，必须同步补齐领域依赖抽象与数据层实现。

需要同步补齐的内容包括：

- 在 `internal/biz/depend/<domain>/` 下定义对应依赖抽象
- 在 `internal/data/depend/<domain>/` 下实现对应依赖
- 在对应 `depend` 聚合入口中纳入 provider
- 按需补齐 `internal/service/inner/` 的 service 暴露与注册联动

典型映射关系如下：

```text
api/open/inner/v1/account.proto
-> internal/biz/depend/open/account.go
-> internal/data/depend/open/account.go

api/system/inner/v1/captcha.proto
-> internal/biz/depend/system/captcha.go
-> internal/data/depend/system/captcha.go
```

边界提示：

- `biz/depend` 只定义依赖抽象与必要领域对象
- `data/depend` 负责基于 generated inner client 落地实现
- 不把 inner client 直接泄漏到业务层
- 若只新增 `inner proto` 而未补 depend 抽象与实现，视为联动不完整

## `service`

定位：协议适配与应用层入口。

承载内容：

- 请求参数转换
- 调用 `usecase`
- 响应结果映射
- 显式构造业务输入对象

边界提示：

- `service` 只做协议适配，不承担业务编排
- `service` 不直接访问 `repo`
- `service` 不维护状态机、事务边界或领域规则
- `service` 暴露的是 `UseCase` 能力
- 当某个能力已经独立为单独 `UseCase` 且需要对外暴露时，应有与之对应的 `service`，不要长期复用其他主题的 `service` 承载该能力
- `queryConfig.relation` 属于查询输入的一部分，应由 `service` 显式构造成 relation 树后继续向下传递
- relation 传递时应保持嵌套结构，清楚表达父子 relation 关系，而不是只传递扁平开关
- 若 relation 查询策略来源于接口语义，应在 `service` 层体现该构造过程，而不是在 `data` 层隐式写死
- 在 relation 分工中，`service` 只负责根据协议语义构造 `filter / opts / relation tree`，不在 `service` 层做 relation 补全

示例：

```go
func (s *AccountService) GetAccount(ctx context.Context, req *v1.GetAccountRequest) (*v1.GetAccountReply, error) {
	account, err := s.uc.GetAccount(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &v1.GetAccountReply{Info: convertAccount(account)}, nil
}
```

## `server`

定位：服务暴露与注册入口。

承载内容：

- HTTP / gRPC server 注册
- 路由暴露
- middleware 挂载
- 启动相关的暴露装配

边界提示：

- `server` 负责注册与暴露，不承载业务逻辑
- 服务注册应集中在统一入口，不散落在业务文件或 `init()`

示例：

```go
func NewGRPCServer(account *service.AccountService) *grpc.Server {
	srv := grpc.NewServer()
	v1.RegisterAccountServiceServer(srv, account)
	return srv
}
```

## gateway

定位：代理、协议转换与聚合转发入口。

承载内容：

- 参数映射
- 下游调用
- 响应适配
- 多下游聚合转发

边界提示：

- gateway 只做代理和聚合适配
- gateway 不直接访问 `repo`
- gateway 不维护领域状态或业务状态流转

示例：

```go
func (s *GatewayService) GetOverview(ctx context.Context, req *v1.GetOverviewRequest) (*v1.GetOverviewReply, error) {
	account, err := s.accountClient.GetAccount(ctx, &openv1.GetAccountRequest{Id: req.Id})
	if err != nil {
		return nil, err
	}
	user, err := s.userClient.GetUser(ctx, &userv1.GetUserRequest{Id: account.Info.UserId})
	if err != nil {
		return nil, err
	}
	return convertOverview(account, user), nil
}
```

## wire 与 provider

定位：接入层依赖装配知识。

承载内容：

- `ProviderSet`
- 构造函数依赖声明
- 统一 `wire` 入口
- 接入层模块装配

边界提示：

- provider、构造函数依赖和 `ProviderSet` 应收敛在既有装配入口
- `wire` 负责依赖装配，不反推业务设计
- 新增能力时先稳定 domain / proto / service，再补 provider 与 wire

示例：

```go
var ProviderSet = wire.NewSet(NewAccountService)
```

```go
func initApp(*conf.Bootstrap, log.Logger) (*kratos.App, func(), error) {
	panic(wire.Build(
		server.ProviderSet,
		service.ProviderSet,
		biz.ProviderSet,
		data.ProviderSet,
		newApp,
	))
}
```

## codegen

定位：接入层与装配链路变化后的生成联动知识。

承载内容：

- proto 生成
- wire 生成
- schema / ent 相关生成联动
- 生成后的最小构建验证

边界提示：

- `codegen` 是联动结果，不是手改入口
- 修改生成物应回到源定义
- 修改 proto、wire、schema 或注册链路后，要同步考虑生成与构建验证

示例：

```text
修改 api/open/v1/account.proto
-> make generate
-> 检查生成产物 diff
-> 如有 provider 联动，再执行 wire
-> go build ./...
```

## 接入层规则

围绕接入层组织时，应优先遵守以下规则：

- `proto` 是契约事实源
- `service` 只做协议适配与 `usecase` 调用
- `server` 只做服务暴露与注册
- gateway 只做代理、映射与聚合适配
- provider 与 `wire` 收敛在统一装配入口
- 生成物不手改，变更回到源定义

## 判断提示

判断接入层问题时，可优先观察：

- 当前改动属于 `proto`、`service`、`server`、gateway、`wire` 还是 `codegen`
- 当前职责是在定义契约、做协议适配、做服务注册还是做依赖装配
- 当前逻辑是否开始承载业务编排或领域规则
- 当前改动是否已经触发生成链路联动

## 边界提示

如果问题继续细化，应转到更具体的知识文档：

- 聚合根识别 -> `aggregate-root.md`
- 围绕聚合根的命名收敛 -> `naming.md`
- usecase、repo、data、事务边界 -> `domain.md`
- ent、listener、consumer、crontab、组件接线 -> `components.md`
- 共享枚举、错误语义与稳定字面量 -> `shared-conventions.md`
