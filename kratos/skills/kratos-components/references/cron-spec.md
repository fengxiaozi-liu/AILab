# Cron Spec

## Job 边界

| 条件 | 做法 |
|------|------|
| 定时触发现有业务能力 | Job 只做触发与错误返回 |
| 复杂业务流程和状态判断 | 回到 UseCase |

```go
// ✅
type SyncJob struct {
    useCase *biz.SyncUseCase
}

func (j *SyncJob) Run(ctx context.Context) error {
    return j.useCase.RunSync(ctx)
}
```

```go
// ❌ cron 回调里写业务流程
func (j *SyncJob) Run(ctx context.Context) error {
    // 查询、判断、补偿、状态流转全写在这里
    ...
}
```

---

## 注册与生命周期

| 条件 | 做法 |
|------|------|
| 新增 job | 通过注册器接入统一 crontab 容器 |
| 启停控制 | 跟随生命周期统一 Start/Stop |

```go
// ✅ 注册器
type syncRegister struct {
    job *SyncJob
}

func (r *syncRegister) register(cron *crontab.Crontab) error {
    return cron.AddFunc("0 */5 * * * *", r.job.Run)
}
```

```go
// ❌ 在业务层到处 Start/Stop
cron.Start()
defer cron.Stop()
```

---

## 幂等与并发

| 条件 | 做法 |
|------|------|
| 重复触发有风险 | 明确幂等键、并发策略、超时 |
| 没有幂等设计 | 不直接上线定时写操作 |

```text
// ✅ 先明确:
job name
cron expr
timeout
concurrency policy
idempotency key
```

---

## 组合场景

```text
Cron -> Job -> UseCase
ProviderSet -> Register -> Lifecycle Start/Stop
```

这个组合场景同时满足：

- Job 足够薄
- 注册统一
- 生命周期统一
- 幂等和并发策略先于实现

---

## 常见错误模式

```text
// ❌ 业务流程塞进 cron 回调
```

```text
// ❌ 启停散落
```

```text
// ❌ 无幂等策略
```
