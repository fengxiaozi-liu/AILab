# Consumer Spec

## Consumer 边界

| 条件 | 做法 |
|------|------|
| 处理消息消费 | 解析消息、组装入参、调用 UseCase、返回错误 |
| 业务状态机、repo 聚合装配 | 回到 UseCase/Repo |

```go
// ✅
func (c *SyncConsumer) handle(ctx context.Context, body []byte) error {
    var payload Payload
    if err := json.Unmarshal(body, &payload); err != nil {
        return err
    }
    return c.useCase.Process(ctx, &payload)
}
```

```go
// ❌ 解析失败吞错
if err := json.Unmarshal(body, &payload); err != nil {
    return nil
}
```

---

## topic 与 queue 语义

| 条件 | 做法 |
|------|------|
| `topic` | 表示监听哪个业务域的消息主题 |
| `queue` | 表示当前服务自己的消费队列/消费组标识 |

```go
// ✅
newConsumer(assetenum.AmqpTopicAsset, openenum.AmqpQueueOpen, 1, cfg, logger)
```

```go
// ❌ 把 topic 和 queue 当成同一语义随意混放
newConsumer(openenum.AmqpQueueOpen, openenum.AmqpQueueOpen, 1, cfg, logger)
```

---

## 注册与生命周期

| 条件 | 做法 |
|------|------|
| 新增 consumer | 通过 register 接口接入统一容器 |
| 启停管理 | 跟随生命周期统一 Run/Stop |

```go
// ✅
type syncRegister struct {
    cfg    *conf.Data_Amqp
    logger log.Logger
    c      *SyncConsumer
}

func (r *syncRegister) register() (*amqp.Consumer, error) {
    return newConsumer(assetenum.AmqpTopicAsset, openenum.AmqpQueueOpen, 1, r.cfg, r.logger)
}
```

```go
// ❌ 单个 consumer 自己管理生命周期
c.Run(ctx)
defer c.Stop()
```

---

## 组合场景

```text
Consumer -> parse payload -> UseCase
register() -> consumer container -> lifecycle Run/Stop
```

这个组合场景同时满足：

- 解析失败会返回错误
- topic/queue 语义清晰
- 生命周期统一收口

---

## 常见错误模式

```text
// ❌ handler 里写业务流程
```

```text
// ❌ 吞错假成功
```

```text
// ❌ 单 consumer 自己启停
```
