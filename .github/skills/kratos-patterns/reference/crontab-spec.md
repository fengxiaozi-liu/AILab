# 定时任务（crontab）规范

> 本文档描述 Kratos 项目中 `internal/crontab/` 的组织方式与最低运行约束。

## 适用范围

- 目录：`internal/crontab/`
- 典型任务：`customer_cache_task.go`

## 基本原则

1. **任务只负责调度与编排，不承载业务状态**
   - 任务结构体允许注入 UseCase/Repo/Server，但不应把可变状态放在任务 struct 里。
   - 运行时快照/缓存应下沉到 Repo 或专用存储层。

2. **任务必须可控：可启停、可超时、可降级**
   - 每个任务的执行函数必须使用 `context.WithTimeout` 设定上限（避免拖垮调度器）。
   - 外部依赖失败（DB/Redis/上游订阅）应明确是“失败中止”还是“best-effort 记录日志继续”。

3. **禁止在热路径执行 O(n) 扫描**
   - 允许在 cron 任务中做 O(n) 扫描（例如刷新后同步 WS 状态），但必须在注释中明确成本与边界。

4. **幂等与互斥策略必须明确**
   - 若任务可能在多实例部署下重复执行，需要明确：
     - 是否允许多实例同时跑（OnOneServer=false）
     - 是否需要单实例互斥（Single=true）
   - 建议优先使用调度框架自身提供的互斥能力，避免自行造锁。

## 目录结构约定

- `internal/crontab/crontab.go`
  - 负责初始化调度器、集中注册所有任务。
- `internal/crontab/*_task.go`
  - 每个文件对应一个任务实体（一个 struct + register + handler）。

## 任务注册约定

以 `CustomerCacheTask` 为例（`internal/crontab/customer_cache_task.go`）：

- 任务必须实现 `register(cron *crontab.Crontab) error` 方法，并通过 `Crontab.register(...)` 统一注册。
- 注册信息必须包含：
  - `Name`：稳定、可读、可用于定位日志
  - `Rule`：cron 表达式
  - `Single`：是否同一时刻只跑一个
  - `OnOneServer`：是否限制在某一台机器执行（取决于部署策略）

## 任务执行约定

1. **统一补业务上下文**
   - 执行入口尽早 `ctx = contextpkg.NewBusinessContext(ctx)`（与项目中间件/上下文保持一致）。

2. **必须设置超时**
   - 示例：
     - `const refreshTimeout = 30 * time.Second`
     - `timeoutCtx, cancel := context.WithTimeout(ctx, refreshTimeout)`

3. **错误处理**
   - 业务刷新失败：返回 err（让调度框架感知失败）
   - best-effort 逻辑失败：记录 Error 日志但不阻塞任务整体完成

4. **任务成本必须可观测**
   - 至少记录：扫描量、影响量、耗时（ms）
   - 示例字段（项目已使用）：`connected/kicked/scanned_session/released_sub_key/upstream_unsub/cost_ms`

## 与 Redis 的关系

- 当前调度器初始化依赖 `*redis.Redis`（见 `internal/crontab/crontab.go`），用于调度框架的互斥/协调能力。
- 业务缓存/快照是否进入 Redis **不在本规范讨论范围**；若涉及，请结合对应 feature 的设计文档与 data 层实现。

## 代码示例（现有实现引用）

- 调度器初始化：`internal/crontab/crontab.go`
- 示例任务：`internal/crontab/customer_cache_task.go`

## 常见陷阱

- 未设置超时：任务卡死导致后续堆积
- 热路径执行全量扫描：请求延迟抖动
- 无幂等/互斥：多实例重复执行造成资源风暴
- best-effort 与强一致混用：该返回 err 时只打日志，导致问题被吞掉

## Go Demo（可复制模板）

> 说明：以下 demo 展示 QuoteNode 推荐的任务组织方式。
> 
> 约束：任务只负责调度与编排；运行时状态下沉到 Repo；必须设置超时；必须记录耗时与影响面。

### 1) 标准任务骨架（struct + register + handle）

```go
package demo

import (
	"context"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"gitlab.linksoft.cn/TeamA/GoPackage/crontab"
	contextpkg "linksoft.cn/node/internal/pkg/context"
)

// ExampleTask 示例任务：仅负责调度与编排，不持有运行时可变状态。
// 注意：状态/缓存/快照必须下沉到 Repo。
type ExampleTask struct {
	uc  UseCase
	log *log.Helper
}

// UseCase 仅示例：真实项目直接注入具体 UseCase。
type UseCase interface {
	Refresh(ctx context.Context) error
}

func NewExampleTask(uc UseCase, logger log.Logger) *ExampleTask {
	return &ExampleTask{
		uc:  uc,
		log: log.NewHelper(logger),
	}
}

// register 由 Crontab 统一调用注册。
func (t *ExampleTask) register(cron *crontab.Crontab) error {
	// 启动时 best-effort 预热一次（失败不阻塞启动）。
	_ = t.handle(contextpkg.NewBusinessContext(context.Background()))

	return cron.Register(crontab.Entity{
		Name:        "example_task",
		Rule:        "0 */5 * * * *", // 每 5 分钟一次
		OnOneServer: false,            // 是否限制在单机执行，取决于部署策略
		Single:      true,             // 同一时刻只跑一个，避免重入
		Handle:      t.handle,
	})
}

func (t *ExampleTask) handle(ctx context.Context) error {
	ctx = contextpkg.NewBusinessContext(ctx)

	const timeout = 30 * time.Second
	timeoutCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	start := time.Now()

	if err := t.uc.Refresh(timeoutCtx); err != nil {
		// 失败：返回 err，让调度框架感知失败
		t.log.WithContext(ctx).Errorf("op=%s err=%+v", "example_task_refresh_failed", err)
		return err
	}

	// 成功：记录耗时（必要时也记录影响面，比如 refreshed_count）
	t.log.WithContext(ctx).Infof(
		"msg=%s cost_ms=%d",
		"[ExampleTask.handle] refresh_ok",
		time.Since(start).Milliseconds(),
	)
	return nil
}
```

### 2) 反例（禁止项）

```go
// ❌ Bad: 不设置超时，可能卡死整个调度器；把运行时状态放在任务 struct 里；热路径做 O(n) 扫描
// 仅作为反例不要复制。
type BadTask struct {
	cache map[string]string // ❌ 禁止：可变状态
}

func (t *BadTask) handle(ctx context.Context) error {
	// ❌ 禁止：没有超时
	_ = ctx

	// ❌ 禁止：热路径做全量扫描
	for i := 0; i < 1_000_000; i++ {
		// ...
	}
	return nil
}
```
