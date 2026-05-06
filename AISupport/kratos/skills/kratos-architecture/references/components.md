# Components

## 作用范围

本文用于说明 Kratos 基础设施组件的接入知识，包括 `ent`、`consumer`、`listener`、`crontab`、`event`、`provider`、`wire` 与生命周期装配。

当问题属于以下场景时，应优先查看本文：

- 新增或修改基础设施组件
- 判断组件应该落在哪个目录
- 判断组件如何接线、注册与进入生命周期
- 判断 `ent`、`listener`、`consumer`、`crontab`、`event` 的职责边界
- 判断 provider、`wire`、注册入口、启动入口、关闭入口如何收口

## 总体规则

- 组件层只负责接入、注册、生命周期管理与最小适配，不承担业务编排、事务边界决策或领域主流程组织。
- `listener`、`consumer`、`crontab`、`server`、`event bus` 等运行时组件不得演化成隐藏的 `UseCase`。
- 修改 `ent`、`event`、`listener`、`consumer`、`provider`、`wire` 时，要同时评估生成链路、注入链路与运行时链路是否完整。
- 当前仓库中的 `eventbus` 按“进程内同步分发”理解，不按 MQ、异步任务或 eventual consistency 默认理解。
- 如果将 `Publish` 从事务内移动到事务外，必须明确说明是否改变失败传播、回滚语义、状态可见性和执行顺序。
- MySQL 下统一使用 `OnConflict()` 组织冲突更新，不使用 `OnConflictColumns(...)`。
- 组件层的目标是让接线、注册、生命周期与联动关系清晰、稳定、可复用，而不是为了抽象额外增加一层。
- 组件接入优先沿用仓库已有目录结构、`ProviderSet` 聚合方式与统一注册入口；新增组件前先确认是否已有同类实现模式，优先复用，不另起体系。
- provider、`wire`、注册入口、启动入口、关闭入口应尽量收敛；组件生命周期应由统一链路驱动，不应在零散位置分散管理。
- 新增组件时，不仅要实现组件本体，还要同时补齐 provider、`wire`、注册入口与生命周期链路。
- 运行时接入代码只做最小适配；`listener`、`consumer`、`crontab`、事件回调、任务入口都不应承载完整业务主流程。
- 稳定名称应集中定义在 enum 或稳定常量中，不在实现中散写字面量；同名事实在本地事件、MQ 事件、listener 名称之间必须明确区分，不得混用。
- payload、消息体、事件载荷优先使用稳定领域对象或显式结构；类型不匹配、解析失败、下游失败时默认直接返回错误，不吞错、不静默降级。
- 本地 `eventbus` 不作为异步 MQ、最终一致性机制或通用解耦手段使用。
- 如果事件发布位置、生命周期接线位置或组件启动方式发生变化，需要同时明确失败传播、事务语义、状态可见性与执行顺序是否改变。

## 组件定位

组件层负责：

- 基础设施接入
- Provider 聚合
- 注册收口
- 生命周期管理
- 最小协议适配

组件层不负责：

- 业务流程编排
- 事务边界决策
- 跨聚合状态机组织
- 抽象成新的通用业务层

组件的重点不是“多写一层”，而是让接线、注册、生命周期与联动关系清晰、稳定、可复用。

## 组件类型

Kratos 基础设施组件通常归为以下几类：

- `ent`
  - `ent/schema`、字段、edge、index、ORM 生成链路
- `consumer`
  - MQ consumer、topic / queue、消息回调、ack / retry 语义、生命周期接入
- `listener`
  - 本地 eventbus listener、handler 聚合注册、事件回调、生命周期联动
- `crontab`
  - scheduler、任务注册、启停、并发与幂等
- `event`
  - 本地事件、eventbus、topic、payload、发布/订阅关系

如果一个改动同时命中多个组件类型，优先看“注册归属与生命周期归属最明确的主改动点”。

## 组件接入规则

- 组件代码优先沿用仓库已有目录结构、`ProviderSet` 聚合方式与统一注册入口。
- 新增组件前，先确认仓库是否已有同类组件实现，优先复用模式，不另起体系。
- provider、`wire`、注册入口、启动入口、关闭入口应尽量收敛，不散落在多个模块。
- 组件自身的生命周期不要自行分散管理，应由统一生命周期链路驱动。
- 任何组件示例都应体现“当前仓库真实模式”，不要为了抽象而脱离实际结构。

## `ent`

定位：ORM schema 与数据生成组件。

承载内容：

- `internal/data/ent/schema`
- field、edge、index
- ORM 生成链路
- schema 注释与落位

边界提示：

- `ent` 的事实源是 schema，而不是生成物。
- 修改 schema、field、edge、index 时，要回到 `internal/data/ent/schema`。
- 不手改 `ent` 生成代码。
- schema、field、edge、index 的 `Comment` 统一使用中文简体，直接表达业务语义。
- 只有具备唯一约束或明确幂等键的字段，才能作为 upsert 冲突依据。

示例：

```text
internal/data/ent/schema/account.go
go generate ./internal/data/ent
```

### `ent schema` 变更闭环规范

- `ent/schema` 是 `ent` 模型的唯一事实来源，字段、边、索引、注释与类型定义都应回源到 `internal/data/ent/schema`
- 任何涉及 `field`、`edge`、`index`、字段类型或关系结构的变更，都应先保证 schema 本身完整，再执行生成链路，再清理旧引用或旧结构
- 双向 `edge`、反向 `edge` 与 `Ref(...)` 关系必须成对修改，不允许只改一侧后依赖生成器或运行期兜底
- 涉及字段替换、关系替换、类型替换时，优先采用“两阶段迁移”思路：先并存生成，再切换引用，最后清理旧结构
- 不应在生成前先拆坏旧关系、删除必需 schema 关联或制造半完成状态，否则容易导致生成失败与目录不一致
- 生成顺序默认应为：更新 schema、确认关系完整、执行 `ent generate`、切换业务引用、清理旧字段或旧关系、再次生成并验证
- `ent generate` 失败时，优先检查 schema 完整性、双向关系是否配对、生成目录是否一致，而不是先手改生成产物
- 修改 `ent/schema` 后，默认需要评估并执行对应生成与构建验证，确保 schema、生成物与业务代码保持同一版本状态

## `consumer`

定位：外部消息进入系统后的接入组件。

### 当前仓库事实

`consumer` 不是“单个消息处理函数”的统称，而是两层结构：

- 聚合容器：`internal/consumer/consumer.go`
- 具体业务 consumer：如 `openConsumer`、`CounterConsumer`、`AssetConsumer`

对应事实：

- `ProviderSet` 统一声明所有具体 consumer 构造器与聚合入口。
- `Consumer` 内部维护 `[]*amqp.Consumer`，由 `Run()` / `Stop()` 统一管理启停。
- 每个具体 consumer 实现 `register() (*amqp.Consumer, error)`，由聚合容器统一收口。
- `newConsumer(...)` 是底层 AMQP consumer 的统一构造入口，负责 topic、queue、连接参数、运行参数装配。
- `consumer` 本身不在 `main` 中手动启动，而是通过生命周期事件联动启动和停止。

关键参考：

- `internal/consumer/consumer.go`
- `internal/listener/lifecycle.go`
- `cmd/server/main.go`

### 职责边界

`consumer` 负责：

- 外部消息入口接入
- `event.Data` 到结构体的解析
- 最小入参组装
- 必要的上下文补齐
- 调用 usecase / repo 完成已有业务能力
- 将错误返回给底层消费框架，交由其处理 nack / retry

`consumer` 不负责：

- 主业务流程编排
- 新的事务边界决策
- 大量跨仓储协调
- 复杂状态机组织
- 在多个位置分散注册和启动 consumer

### 标准模式

聚合容器模式：

```go
type Consumer struct {
	consumers []*amqp.Consumer
}

type register interface {
	register() (*amqp.Consumer, error)
}

func NewConsumer(a *FooConsumer, b *BarConsumer) (*Consumer, error) {
	consumer := &Consumer{
		consumers: []*amqp.Consumer{},
	}

	if err := consumer.register(a, b); err != nil {
		return nil, err
	}

	return consumer, nil
}

func (c *Consumer) Run(_ context.Context) error {
	for _, item := range c.consumers {
		if err := item.Run(); err != nil {
			return err
		}
	}
	return nil
}

func (c *Consumer) Stop() error {
	for _, item := range c.consumers {
		_ = item.Close()
	}
	return nil
}
```

具体业务 consumer 模式：

```go
type OrderConsumer struct {
	cfg    *conf.Data_Amqp
	logger log.Logger
	uc     *biz.OrderUseCase
}

func NewOrderConsumer(data *conf.Data, logger log.Logger, uc *biz.OrderUseCase) *OrderConsumer {
	return &OrderConsumer{
		cfg:    data.GetAmqp(),
		logger: logger,
		uc:     uc,
	}
}

func (c *OrderConsumer) register() (*amqp.Consumer, error) {
	consumer, err := newConsumer(domainenum.AmqpTopicOrder, pushenum.AmqpQueuePush, 1, c.cfg, c.logger)
	if err != nil {
		return nil, err
	}

	consumer.RegisterConsumeFunc(domainenum.AmqpEventOrderCreated.Value(), c.orderCreated)
	consumer.RegisterConsumeFunc(domainenum.AmqpEventOrderCanceled.Value(), c.orderCanceled)

	return consumer, nil
}
```

消息回调模式：

```go
func (c *OrderConsumer) orderCreated(ctx context.Context, event *amqp.Event) error {
	ctx = contextpkg.NewBusinessContext(ctx)

	var payload biz.OrderCreatedPayload
	if err := helper.MapToStruct(event.Data, &payload); err != nil {
		return err
	}

	return c.uc.HandleOrderCreated(ctx, &biz.HandleOrderCreatedCommand{
		OrderID: payload.OrderID,
		UserID:  payload.UserID,
	})
}
```

### 新增 `consumer` 的统一规范

- 新类型必须实现 `register() (*amqp.Consumer, error)`。
- 必须接入 `ProviderSet` 与聚合入口 `NewConsumer(...)`。
- topic、queue、event 名称必须来自 enum 或稳定常量，不在实现里散写字面量。
- 消息 handler 只做入口适配，不在 `consumer` 内堆叠完整业务流程。
- 解析失败、类型错误、下游失败默认直接返回错误，不吞错、不静默降级。
- 若有明确业务语义允许忽略异常，必须在代码和文档里显式说明原因，不能作为通用模板。
- 需要幂等、重试、死信、顺序消费时，先定义消息语义，再决定放在哪一层实现。
- 不在 `main`、`init`、`wire_gen.go` 之外的零散位置手动调用 `Run()` 或注册 AMQP handler。

### topic / queue / event 定义规范

AMQP 相关稳定名称统一放在 enum 中定义，不在 consumer、listener、usecase 中散写字符串。

建议落位：
- `internal/enum/<domain>/<domain>.go` 或同域 enum 文件定义具体业务域的 `topic / queue / event` 常量，不把业务语义常量放进 `base`

业务域定义示例：

```go
package order

import baseenum "linksoft.cn/alpha/internal/enum/base"

const (
	// AmqpTopicOrder 表示消息由 order 域发出。
	AmqpTopicOrder baseenum.AmqpTopic = "order"
)

const (
	// AmqpQueueOrder 表示消息由 order 这一侧消费。
	AmqpQueueOrder baseenum.AmqpQueue = "order"
)

const (
	// AmqpEventOrderCreated 表示订单已创建。
	AmqpEventOrderCreated baseenum.AmqpEvent = "orderCreated"
	// AmqpEventOrderCanceled 表示订单已取消。
	AmqpEventOrderCanceled baseenum.AmqpEvent = "orderCanceled"
)
```

统一定义原则：

- 只要常量名带有明确业务语义，如 `Order`、`Open`、`Push`、`UserMessageCreate`，就应定义在对应业务域 enum 包
- `topic` 定义在发送方业务域 enum 中，表达“谁发”
- `queue` 定义在消费方业务域 enum 中，表达“谁收”
- `event` 定义在发送方业务域 enum 中，表达“发生了什么”
- 新增 consumer 时，优先复用“上游域 `AmqpTopicXxx` + 当前服务域 `AmqpQueueXxx` + 上游域 `AmqpEventXxx`”这一组合模式
- 不要因为底层类型来自 `baseenum`，就把业务域常量也错误地下沉到 `internal/enum/base`
- 不要复用仓库中历史上把 `AmqpEvent`、`LocalEvent`、`Listener` 误声明成错误基础类型的写法；新增代码必须保证名字语义和底层类型一致

使用方式示例：

```go
consumer, err := newConsumer(orderenum.AmqpTopicOrder, pushenum.AmqpQueuePush, 1, c.cfg, c.logger)
if err != nil {
	return nil, err
}

consumer.RegisterConsumeFunc(orderenum.AmqpEventOrderCreated.Value(), c.orderCreated)
consumer.RegisterConsumeFunc(orderenum.AmqpEventOrderCanceled.Value(), c.orderCanceled)
```

定义约束：

- `topic` 表示消息的发送侧归属或发布主题，重点回答“这条消息是谁发出来的、属于哪个发送方域”。
- `queue` 表示消息的接收侧归属，重点回答“这条消息最终由谁消费、落到哪个消费方队列”。
- `event` 表示具体业务事实，命名应表达“已发生的事情”，而不是动作命令。
- 同一领域的 topic / event / local event / listener 名称应统一放在同一个 enum 文件或同一组 enum 文件中维护。
- 同名事实如果同时存在本地事件与 MQ 事件，应明确区分 `LocalEventXxx` 与 `AmqpEventXxx`，不要混用。

### topic / queue / event 语义说明

可以按“谁发送、谁接收、发生了什么”理解：

- `topic`
  - 表示发送方主题
  - 重点描述消息来源域、来源系统或来源能力
  - 更偏向“生产者是谁”
- `queue`
  - 表示接收方队列
  - 重点描述消息由哪个消费方收取
  - 更偏向“消费者是谁”
- `event`
  - 表示在该消息上的具体事实
  - 重点描述“发生了什么”

组合理解：

```text
topic = 谁发
queue = 谁收
event = 发了什么
```

例如当前仓库中的消费写法：

```go
consumer, err := newConsumer(openenum.AmqpTopicOpen, pushenum.AmqpQueuePush, 1, c.cfg, c.logger)
```

其语义是：

- `openenum.AmqpTopicOpen`
  - 说明这类消息来自 `open` 领域或开户侧发送方
- `pushenum.AmqpQueuePush`
  - 说明这类消息由 `push` 服务这一侧接收和消费
- `openenum.AmqpEventOpenFirstReviewPass`
  - 说明这条消息表达的是“开户一审通过”这一已发生事实

### 当前仓库中的发送方 / 接收方事实

当前仓库里，AMQP 的发送方和接收方要分开看：

- 发送方 topic
  - 当前服务作为生产者发送消息时，producer topic 来自配置 `conf.Amqp.ProducerTopic`
  - 该值在 `internal/data/kit/data.go` 中创建 producer 时传入
  - 因此“当前服务对外发送时属于哪个发送主题”首先受部署配置控制
- 接收方订阅
  - 当前服务作为消费者时，通过 `newConsumer(topic, queue, ...)` 订阅指定 topic，并绑定到本服务自己的 queue
  - 这里的 `topic` 通常表示上游发送方域，例如 `open`、`asset`、`counter`
  - 这里的 `queue` 表示当前服务以什么消费方身份接收，例如 `push`

因此文档里要避免把 `topic` 简化成“只是消费者订阅名”：

- 对发送方而言，`topic` 是发布主题，是发送侧身份的一部分
- 对接收方而言，`topic` 是订阅来源，是消费时匹配的上游主题
- `queue` 始终更接近接收方 / 消费方身份

### 设计建议

- 定义 `topic` 时，优先按发送方领域或发送方系统命名，例如 `open`、`asset`、`counter`。
- 定义 `queue` 时，优先按消费方系统或消费方职责命名，例如 `push`、`order`、`market`。
- 定义 `event` 时，优先按业务事实命名，例如 `orderCreated`、`openFirstReviewPass`，不要写成命令式名称。
- 如果当前服务既是发送方又是接收方，应分别说明它在哪些链路中作为生产者，在哪些链路中作为消费者。
- 如果 producer topic 受配置控制，文档中要注明“发送主题的最终取值来自配置”，而不是假设它一定等于某个 enum 常量。

### 评审检查点

- 是否沿用了聚合容器模式，而不是自己分散启动 consumer
- 是否使用统一 `newConsumer(...)`
- 是否使用 enum 声明 topic / queue / event
- handler 是否只做协议适配与最小转换
- 错误处理是否清晰，是否存在无说明的吞错
- 生命周期是否接到了 `LifeCycleListener`

## `listener`

定位：本地事件响应组件。

### 当前仓库事实

当前仓库中的 `listener` 是“进程内 eventbus listener 的统一注册体系”，不是“任何事件回调代码都算 listener”。

对应事实：

- `internal/listener/listener.go` 是统一聚合注册入口。
- 每个具体 listener 只负责自己的事件注册与 handler 实现。
- 所有 listener 通过 `NewListener(...)` 内部统一调用 `register(...)` 收口。
- 不应在其他零散位置直接调用 `eventBus.RegisterListener(...)`。
- 生命周期相关联动也通过 listener 实现，而不是由 `main` 直接启动组件。

关键参考：

- `internal/listener/listener.go`
- `internal/listener/lifecycle.go`
- `internal/listener/user_message.go`

### 职责边界

`listener` 负责：

- 本地同步事件的监听
- listener 统一注册
- payload 类型适配
- 必要的后置动作
- 必要时将本地事件转发为外部 MQ 消息

`listener` 不负责：

- 直接承担主业务流程
- 充当新的 usecase 编排层
- 大量 repo 协调
- 复杂事务决策
- 分散注册 eventbus handler

### 两类 listener

1. 生命周期 listener

适用场景：

- 初始化
- 预热
- 启停联动
- `crontab` / `consumer` 启停收口

当前仓库模式：

- Kratos 生命周期钩子发布本地生命周期事件
- `LifeCycleListener` 监听这些事件
- 在 `beforeStart` 中执行初始化、启动 `crontab`、启动 `consumer`
- 在 `afterStop` 中停止 `consumer`、停止 `crontab`

2. 业务 listener

适用场景：

- 对领域事实做本地同步响应
- 最小后置处理
- 将本地事件桥接到 MQ

当前仓库典型链路：

```text
usecase -> eventbus.Publish -> listener -> producer.Publish
```

例如：

- `UserMessageUseCase.Create()` 发布 `LocalEventUserMessageCreate`
- `UserMessageListener.create()` 监听该事件
- listener 内完成结构转换并转发 `AmqpEventUserMessageCreate`

### 标准模式

聚合注册模式：

```go
var ProviderSet = wire.NewSet(
	NewLifeCycleListener,
	NewOrderListener,
	NewListener,
)

type Listener struct {
	eventBus *eventbus.EventBus
}

type register interface {
	register(eventBus *eventbus.EventBus)
}

func NewListener(
	lifeCycleListener *LifeCycleListener,
	orderListener *OrderListener,
	eventBus *eventbus.EventBus,
) *Listener {
	listener := &Listener{
		eventBus: eventBus,
	}

	listener.register(
		lifeCycleListener,
		orderListener,
	)

	return listener
}
```

生命周期 listener 模式：

```go
func (l *LifeCycleListener) register(eventBus *eventbus.EventBus) {
	listeners := []*eventbus.Listener{
		{
			Name:   baseenum.ListenerLifecycleBeforeStart.Value(),
			Topics: []string{baseenum.LocalEventLifecycleBeforeStart.Value()},
			Handle: l.beforeStart,
		},
		{
			Name:   baseenum.ListenerLifecycleAfterStop.Value(),
			Topics: []string{baseenum.LocalEventLifecycleAfterStop.Value()},
			Handle: l.afterStop,
		},
	}

	for _, item := range listeners {
		eventBus.RegisterListener(item)
	}
}

func (l *LifeCycleListener) beforeStart(ctx context.Context, _ interface{}) error {
	if err := l.initializer.Init(ctx); err != nil {
		return err
	}
	if err := l.crontab.Start(); err != nil {
		return err
	}
	return l.consumer.Run(ctx)
}

func (l *LifeCycleListener) afterStop(_ context.Context, _ interface{}) error {
	_ = l.consumer.Stop()
	_ = l.crontab.Stop()
	return nil
}
```

业务 listener 模式：

```go
type UserMessageListener struct {
	producer *amqp.Producer
}

func (l *UserMessageListener) register(eventBus *eventbus.EventBus) {
	listeners := []*eventbus.Listener{
		{
			Name:   pushenum.ListenerUserMessageCreate.Value(),
			Topics: []string{pushenum.LocalEventUserMessageCreate.Value()},
			Handle: l.create,
		},
	}

	for _, item := range listeners {
		eventBus.RegisterListener(item)
	}
}

func (l *UserMessageListener) create(ctx context.Context, payload interface{}) error {
	info, ok := payload.(*biz.UserMessage)
	if !ok || info == nil {
		return errors.New("invalid user message payload")
	}

	data := map[string]interface{}{}
	if err := helper.StructConvert(info, &data); err != nil {
		return err
	}

	return l.producer.Publish(&amqp.Event{
		Event: pushenum.AmqpEventUserMessageCreate.Value(),
		Data:  data,
	})
}
```

### 新增 `listener` 的统一规范

- 新类型必须实现 `register(eventBus *eventbus.EventBus)`。
- 必须通过 `ProviderSet` 和 `NewListener(...)` 聚合注册。
- 不在 `main`、`usecase`、`service` 或其他零散位置直接 `RegisterListener(...)`。
- listener 名称与本地事件 topic 必须来自 enum 或稳定常量。
- payload 优先使用稳定领域对象或显式 payload 结构，不默认退化为 `map[string]interface{}`。
- payload 类型断言失败时应立即返回错误，不静默忽略。
- 一个 listener 可以监听多个 topic，但前提是它们属于同一职责闭环。
- 如果一个 listener 开始同时处理多类无关事件，应拆分为多个 listener 类型。

### 评审检查点

- 是否通过聚合入口统一注册
- 是否把生命周期联动写进了 `LifeCycleListener`
- 是否把业务后置动作错塞进了 usecase 主流程
- 是否错误地把 listener 当成异步 MQ 消费器
- payload 类型是否稳定清晰
- 是否存在无说明的强制断言、吞错或分散注册

## `event`

定位：本地事件与 eventbus 使用规范。

承载内容：

- 本地事件 topic / event 的稳定命名
- listener 名称与订阅关系
- payload 结构选择
- 发布位置与失败传播语义
- eventbus 基础设施落位与注入方式
- 本地事件与 MQ 事件的边界划分
- 本地后置流程编排与生命周期联动

### 当前仓库事实

当前仓库中的 `eventbus` 默认按“进程内同步分发”理解：

- `Publish()` 返回前，相关 listener 已同步执行
- listener 返回错误会沿调用链同步传播
- 它更接近“本地钩子 / 本地回调分发器”，而不是 MQ

同时，`eventbus` 在当前仓库中也是明确的基础设施组件，而不只是概念说明：

- 在 `internal/data/kit/data.go` 中统一创建 `eventBus`
- 通过 `NewEventBus(data *Data)` 提供给业务层和 listener 装配使用
- 作为依赖注入对象参与 usecase、listener、应用生命周期链路
- 生命周期事件由 `cmd/server/main.go` 中的 Kratos 钩子发布到 eventbus
- 通过“发布本地事实 -> listener 响应”的方式参与部分后置流程编排

本地事件相关命名在当前仓库中有明确落位：

- 通用生命周期事件定义在 `internal/enum/base/eventbus.go`
- 业务域本地事件定义在各自 enum 中，例如 `internal/enum/open/event.go`
- `Listener` 与 `LocalEvent` 都应作为稳定字面量集中维护，而不是在发布方和监听方重复散写

当前仓库典型链路：

```text
usecase -> eventbus.Publish -> listener.Handle
```

示例：

```go
func (uc *UserMessageUseCase) Create(ctx context.Context, info *UserMessage) (*UserMessage, error) {
	res, err := uc.userMessageRepo.CreateUserMessage(ctx, info)
	if err != nil {
		return nil, err
	}

	if err := uc.eventBus.Publish(ctx, &eventbus.Event{
		Topic:   pushenum.LocalEventUserMessageCreate.Value(),
		Payload: res,
	}); err != nil {
		return nil, err
	}

	return res, nil
}
```

### 规则

- event 的职责不只是贯通 `listener`、`consumer`，还包括沉淀“稳定事件命名、payload 约定、发布语义、订阅边界”。
- eventbus 不只是事件分发器，也承担一部分“本地同步后置流程编排”的职责，例如生命周期联动、业务完成后的后置动作串联。
- 但 eventbus 参与的是“后置流程编排”和“运行时联动”，不是主业务流程编排，不应演化成新的业务调度中心。
- topic / event 名称统一定义在 enum 或稳定常量中。
- 发布的是“已经发生的事实”，不是给订阅方定制的流程命令。
- payload 优先使用稳定领域对象或显式 payload 结构。
- 本地 event 不替代事务边界设计。
- 不要因为“看起来像事件驱动”就把本地同步 eventbus 当作异步解耦手段。
- 如果 `Publish` 位于事务内，要明确 listener 失败是否参与事务失败传播。
- 如果将 `Publish` 从事务内移到事务外，要明确这是语义变更，而不是普通重构。

### `event -> listener -> consumer` 贯通规范

- `event`
  - 负责定义本地同步事实、稳定 topic、listener 名称、payload 约定，以及必要的本地后置流程编排入口
- `listener`
  - 负责本地同步监听、最小后置逻辑、必要时桥接 MQ
- `consumer`
  - 负责外部消息进入系统后的解析、校验和 usecase 调用

### 评审检查点

- 本地事件名称是否集中定义在 enum，而不是散写字面量
- payload 是否稳定，是否能被订阅方长期理解和复用
- `Publish` 的位置是否清楚表达了失败传播和事务语义
- 是否把本地 eventbus 误写成异步消息系统
- 是否把 eventbus 的后置流程编排职责滥用成主流程编排中心
- eventbus 是否仍通过统一基础设施入口创建和注入
- 是否错误地把 MQ event、LocalEvent、Listener 名称混为一谈

统一判断原则：

- 进程内同步后置逻辑，优先使用 `event + listener`
- 外部系统投递消息、跨进程消费、需要 queue / ack / retry 语义，使用 `consumer`
- 启停、初始化、预热等运行时联动，收口到生命周期 `listener`
- 如果实现开始在 `event`、`listener`、`consumer` 中堆积领域决策、状态机和 repo 协调，应回退到 `usecase`

## `crontab`

定位：定时任务接入组件。

承载内容：

- cron 注册
- 调度入口
- 生命周期启停
- 最小任务触发闭环

边界提示：

- Job 要足够薄，只做触发与错误返回。
- 复杂业务流程应回到 `usecase`。
- 幂等策略与并发策略要先于具体实现定义。

## provider、`wire` 与生命周期

定位：组件接入的统一装配知识。

承载内容：

- provider 聚合
- `ProviderSet`
- 注册器
- 生命周期 Run / Stop
- 启动入口与关闭入口

边界提示：

- 新增组件不仅要写组件本体，还要补齐 provider、`wire`、注册入口与生命周期链路。
- 组件自己的生命周期不要单独散落管理。
- 修改 schema、provider、`wire`、注册入口时，要同步评估生成、构建与启动链路。

## 判断提示

判断某个改动是否属于组件问题时，可优先检查：

- 当前改动是在做基础设施接入，还是在做业务编排
- 是否涉及 provider、`wire`、注册入口、生命周期
- 当前仓库是否已有稳定目录与聚合入口
- 当前逻辑是最小适配，还是已经开始承担业务主线
- 是否已经补齐生成链路、构建链路与运行链路

## 边界延伸

如果问题继续细化，应转到更具体的知识文档：

- usecase、repo、data、`kit`、事务边界 -> `domain.md`
- service、proto 结构与协议边界 -> `service.md`
- `internal/pkg` 公共能力边界 -> `pkg.md`
- 共享枚举、错误语义、稳定字面量 -> `shared-conventions.md`
