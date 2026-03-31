# Layer

## 作用范围

本文用于说明当前仓库的整体目录结构、项目类型判断方式、各层职责与代码落位边界。

当问题属于以下场景时，应优先查看本文：

- 当前仓库属于哪一类 Kratos 项目
- 某段代码、职责或改动应放在哪一层
- 某个目录在当前仓库中承担什么作用
- 如何理解 `api`、`service`、`server`、`biz`、`data`、`listener`、`consumer`、`crontab`、`pkg`
- 仓库文档与实际代码结构不完全一致时，应如何判断

本文重点回答目录、层级与职责归属，不展开聚合根、命名、组件细节、公共能力子域细节与代码风格红线。

## 仓库总览

当前仓库是一个承载 Kratos 体系下协议定义、业务实现、基础设施接入与运行时组件的综合型仓库。

## 文档组织架构

从当前仓库的文档与代码组织上，可以先建立如下目录地图：

```text
api/
assets/
cmd/
configs/
internal/
├── api/
├── biz/
├── conf/
├── consumer/
├── crontab/
├── data/
├── enum/
├── error/
├── listener/
├── pkg/
├── server/
└── service/
manifests/
```

理解目录职责时，优先先看“它属于哪一层”，再看“它在该层中承担什么职责”。

这张目录地图的重点不是做第二套目录分类，而是帮助判断：

- 某段代码属于协议、应用、领域、数据、组件还是公共能力
- 某个目录是实现主体、装配入口、共享定义还是过程资产
- 当前问题应先落在哪个知识域，再进入更细 reference

## 项目类型判断

项目类型判断必须基于仓库事实，而不是只看单个名称。

Kratos 项目按主类型识别为以下三类之一：

- `BaseService`
  - 以共享协议、共享定义、公共契约与基础抽象为主
  - 通常更强调 `api`、`internal/enum`、`internal/error` 等共享资产
  - 不以完整业务闭环为核心
- `GatewayService`
  - 以外部接入、协议适配、聚合转发与下游代理为主
  - 通常更强调 `service`、`server`、gateway 代理与协议转换
  - 不以领域规则沉淀为主
- `BusinessService`
  - 以业务实现、领域规则、数据落地与运行时组件为主
  - 通常存在完整的 `biz`、`data`、`service`、`server`
  - 常伴随 `listener`、`consumer`、`crontab` 等运行时组件

识别时应综合以下信号判断：

- `.env.*`下的`SERVER_NAME`为主要判断依据, BaseService是BaseService, GatewayService定义GatewayService, 其他为业务服务
- 目录结构
- `wire` 装配范围
- 是否存在完整业务闭环
- 是否以共享定义、协议转发或业务实现为主

当命名信号与目录结构冲突时，以目录结构与装配事实为准。

## 目录职责

### `api`

定位：协议定义目录。

承载内容：

- proto 文件
- service / message / rpc 契约
- open / admin / inner 等 side 协议边界

边界提示：`api` 是协议事实源，不承载业务实现。

### `cmd`

定位：应用启动与装配入口目录。

承载内容：

- `main.go`
- `wire.go`
- `wire_gen.go`

边界提示：`cmd` 负责启动入口，不承载具体业务规则。

### `configs`

定位：配置模板与配置映射目录。

承载内容：

- 服务名
- 数据源配置
- 注册中心配置
- 运行时基础设施配置

边界提示：`configs` 提供配置来源，不承载业务语义。

### `internal`

定位：项目核心实现主体目录。

承载内容：

- `api`
- `service`
- `server`
- `biz`
- `conf`
- `data`
- `listener`
- `consumer`
- `crontab`
- `pkg`
- `enum`
- `error`

边界提示：`internal` 下各子目录按职责分层，不应把 `internal` 整体视为单一实现层。

## `internal` 下关键目录说明

### `internal/api`

定位：协议生成代码目录。

承载内容：

- `pb.go`
- `grpc.pb.go`
- `http.pb.go`
- `pb.validate.go`

边界提示：它是协议定义的生成结果，不是手工维护的协议事实源。

### `internal/conf`

定位：配置协议与配置结构目录。

承载内容：

- 配置 proto
- 配置生成代码
- 运行时配置结构定义

边界提示：它用于承接配置结构，不承担业务规则。

### `internal/service`

定位：application 层与协议适配目录。

承载内容：

- 从协议请求转换到 usecase 调用
- 组织应用层返回结果
- 保持 service 代码薄而清晰

边界提示：`service` 负责应用层适配，不拥有业务主流程、数据访问实现与运行时组件注册职责。

### `internal/server`

定位：服务暴露与 server 注册目录。

承载内容：

- HTTP / gRPC server 注册
- 对外路由暴露
- middleware 挂载
- 与应用启动相关的 server 侧装配

边界提示：`server` 负责暴露与注册，不拥有业务规则、usecase 编排与数据持久化职责。

### `internal/biz`

定位：领域组织与业务规则目录。

承载内容：

- usecase
- 领域对象
- repo 抽象
- 业务不变量
- 事务边界与业务编排

边界提示：`biz` 是业务主线，不承载协议定义与持久化实现细节。

### `internal/data`

定位：持久化实现与外部依赖实现目录。

承载内容：

- repo 实现
- ent 集成
- 下游 service client 实现
- 数据模型转换
- 基础设施支持下的数据访问

边界提示：`data` 负责实现，不反向拥有业务主流程。

### `internal/listener`

定位：本地事件监听目录。

承载内容：

- listener 注册
- 本地事件回调
- 事件驱动触发入口

边界提示：`listener` 负责运行时事件入口，不承载业务主流程组织。

### `internal/consumer`

定位：消息消费入口目录。

承载内容：

- MQ consumer 注册
- topic / queue 接入
- payload 接收
- 向应用 / 业务层转交

边界提示：`consumer` 负责消费入口，不拥有协议定义与领域主线。

### `internal/crontab`

定位：定时任务入口目录。

承载内容：

- cron 注册
- 调度入口
- 最小任务触发闭环

边界提示：`crontab` 负责调度入口，不承担领域建模职责。

### `internal/pkg`

定位：公共技术能力目录。

承载内容：

- `context`
- `filter`
- `localize`
- `metadata`
- `middleware`
- `proto`
- `schema`
- `seata`
- `util`

边界提示：只有脱离具体业务语义后仍成立的能力，才进入 `internal/pkg`。

### `internal/enum`

定位：共享枚举与稳定值域目录。

承载内容：

- typed const
- 稳定值域
- 共享枚举定义

### `internal/error`

定位：错误语义与共享错误表达目录。

承载内容：

- 错误语义
- 共享错误表达
- 跨模块错误定义

## 分层边界

进行层级判断时，应优先遵守以下边界：

- 协议定义放 `api`
- application 适配放 `internal/service`
- 暴露与注册放 `internal/server`
- 业务规则放 `internal/biz`
- 持久化与下游依赖实现放 `internal/data`
- 运行时事件与消费入口放 `listener` / `consumer` / `crontab`
- 脱离业务语义后的公共技术能力才进入 `internal/pkg`

如果某段代码看起来可以放在多个层，应优先放到“职责真正拥有者”所在层，而不是“调用最方便”的那一层。

## 目录判断原则

在当前仓库中做目录归属判断时，应优先遵守以下静态规则：

- 已形成稳定模式的目录优先沿用，不平行发明新落位
- 目录归属优先看职责拥有权，而不是调用便利性
- 带业务语义的实现优先留在业务相关层，不伪装成公共能力
- 适配职责归适配层，职责拥有权归拥有者层

## 边界提示

如果问题继续细化，应转到更具体的知识文档：

- 聚合根识别 -> `aggregate-root.md`
- 领域组织 -> `domain.md`
- 命名收敛 -> `naming.md`
- service / proto / server -> `service.md`
- ent / wire / listener / consumer / cron -> `components.md`
- `internal/pkg` 判断 -> `pkg.md`
- error / enum / logging / comment -> `shared-conventions.md`
- 代码红线 -> `code-style.md`
