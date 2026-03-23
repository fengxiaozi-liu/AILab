# Crontab Spec

## 注册与执行分离

```go
// ✅ 注册与执行分离，Run 函数只做调度编排
type SyncAccountJob struct {
    repo biz.AccountRepo
    log  *log.Helper
}

func (c *Crontab) register() error {
    return c.cron.Register("sync_account", "0 */5 * * * *", c.syncJob.Run)
}

func (j *SyncAccountJob) Run(ctx context.Context) error {
    list, err := j.repo.ListNeedSyncAccount(ctx, &biz.AccountFilter{
        Status: openenum.AccountStatusPending,
    })
    if err != nil {
        return err
    }
    for _, item := range list {
        if err := j.repo.MarkSyncing(ctx, item.ID); err != nil {
            j.log.WithContext(ctx).Errorf("mark syncing failed: account=%d err=%v", item.ID, err)
            continue
        }
    }
    return nil
}

// ❌ 注册和执行混写，任务函数直接包含大量状态操作
func NewCrontab(...) {
    cron.Register("sync", "0 */5 * * * *", func(ctx context.Context) error {
        rows, _ := db.Query("SELECT * FROM account WHERE status=1")
        for rows.Next() { /* 直接写库 */ }
        return nil
    })
}
```

---

## 幂等设计

| 任务类型 | 幂等要求 |
|----------|---------|
| 补偿/重跑任务 | MUST 幂等：重复执行不破坏数据 |
| 巡检/统计任务 | SHOULD 幂等或只做读操作 |
| 周期性写入 | MUST 在上线前确认幂等保证（如使用唯一索引或状态机控制） |

```go
// ✅ 幂等写：Upsert 保证重复执行安全
err = r.data.Db.AccountSnapshot(ctx).Create().
    SetAccountID(accountID).
    SetSnapshotData(data).
    OnConflict().
    UpdateSnapshotData().
    Exec(ctx)

// ❌ 无幂等保证的直接 INSERT
_, err = r.data.Db.AccountSnapshot(ctx).Create().
    SetAccountID(accountID).
    SetSnapshotData(data).
    Save(ctx)
```

---

## Crontab 初始化

```go
// ✅ 标准 Crontab 初始化
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

---

## 组合场景

```go
// 完整任务：注册 + 幂等补偿 + 明确并发策略（分布式锁由 crontab 框架保证）
func (c *Crontab) register() error {
    return c.cron.Register("compensate_account", "0 */10 * * * *", c.compensateJob.Run)
}

func (j *CompensateAccountJob) Run(ctx context.Context) error {
    // 明确范围：只处理 pending 超过 10 分钟的记录
    deadline := uint32(time.Now().Add(-10 * time.Minute).Unix())
    list, err := j.repo.ListStaleAccount(ctx, &biz.AccountFilter{
        Status:         openenum.AccountStatusPending,
        CreateTimeLte:  deadline,
    })
    if err != nil {
        return err
    }
    for _, item := range list {
        // ✅ 幂等：先检查状态，再执行
        if err := j.uc.TryCompensate(ctx, item.ID); err != nil {
            j.log.WithContext(ctx).Errorf("compensate failed: account=%d err=%v", item.ID, err)
        }
    }
    return nil
}
```

---

## 常见错误模式

```go
// ❌ 无幂等保证的周期写
func (j *Job) Run(ctx context.Context) error {
    _, err = db.Exec("INSERT INTO snapshots (...) VALUES (...)")
    return err
}

// ❌ 任务内可变状态（第二次运行使用了上次的缓存结果）
type Job struct { lastResult []*Account }

// ❌ 缺少超时和失败处理，长时间阻塞
func (j *Job) Run(ctx context.Context) error {
    for _, id := range ids {
        doSlowRemoteCall(id) // 无超时
    }
}
```
