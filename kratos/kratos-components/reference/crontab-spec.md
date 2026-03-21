# Crontab Reference

## 这个主题解决什么问题

说明 Kratos 定时任务如何组织注册、执行函数和运行时状态，以及如何设计可维护的调度任务。

## 适用场景

- 新增定时任务
- 重构现有任务
- 设计补偿任务或巡检任务

## 设计意图

Crontab 参考主要解释调度入口和任务执行逻辑为什么要分开，以及定时任务为什么更适合设计成可重复运行的业务入口。

- 调度注册是运行时入口，任务执行是业务行为，两者分开后更容易复用和观测。
- 定时任务经常伴随补偿、批处理和重跑，需要比普通接口更容易看出执行范围与状态来源。
- 先理解任务目标和执行粒度，就不容易把任务写成无边界的大扫描脚本。

## 实施提示

- 先说明任务是巡检、补偿、同步还是定时聚合。
- 再把注册入口、执行函数、状态读取和结果写回拆开组织。
- 如果一次运行会处理大量对象，优先先设计批次、游标或范围切分方式。

## 推荐结构

- 任务对象负责注册和调度入口
- 运行时状态、快照、缓存等下沉到 Repo 或持久层

## 标准模板

```go
type SyncAccountJob struct {
    repo biz.AccountRepo
}

func (j *SyncAccountJob) register(cron *crontab.Crontab) error {
    return cron.Register("sync_account", "0 */5 * * * *", j.Run)
}
```

```go
func (j *SyncAccountJob) Run(ctx context.Context) error {
    list, err := j.repo.ListNeedSyncAccount(ctx, &biz.AccountFilter{Status: openenum.AccountStatusPending})
    if err != nil {
        return err
    }
    for _, item := range list {
        if err := j.repo.MarkSyncing(ctx, item.ID); err != nil {
            return err
        }
    }
    return nil
}
```

## Good Example

- Run 函数只负责调度与编排
- 业务状态变更和查询委托给 Repo

## 项目通用注册入口示例

```go
var ProviderSet = wire.NewSet(
    NewCrontab,
)

type Crontab struct {
    cron *crontab.Crontab
    cfg  *conf.Data_Crontab
}

func NewCrontab(redis *redis.Redis, conf *conf.Data, logger log.Logger) (*Crontab, error) {
    cron, err := crontab.NewCrontab(redis, constant.NewConfig(
        constant.WithSecondLevel(true),
        constant.WithMutexKeyPrefix("crontab"),
        constant.WithLogger(log.NewHelper(logger)),
    ))
    if err != nil {
        return nil, err
    }

    c := &Crontab{cron: cron, cfg: conf.Crontab}
    if err := c.register(); err != nil {
        return nil, err
    }
    return c, nil
}
```

## 生命周期示例

```go
func (c *Crontab) Start() error {
    if strings.EqualFold(c.cfg.Enable, "true") {
        c.cron.Start()
    }
    return nil
}

func (c *Crontab) Stop() error {
    if strings.EqualFold(c.cfg.Enable, "true") {
        return c.cron.Stop()
    }
    return nil
}
```

## 常见坑

- 任务内部持有大量可变状态
- 定时任务直接做复杂数据扫描和写入而无清晰边界
- 多个任务共用一套模糊的日志和观测字段

## 相关 Rule

- `../rules/crontab-rule.md`
