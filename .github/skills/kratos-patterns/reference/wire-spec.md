# Wire / DI 规范（Kratos 通用）

> 本文档描述 Kratos 项目中使用 Google Wire 的通用规范。
> 
> “目录与聚合方式”属于 `project-patterns`；这里仅描述 Wire 的通用落地与推荐约束。

## 核心原则

1. **依赖注入只负责装配，不承载业务逻辑**
2. **构造函数显式依赖（参数列表明确）**
3. **Data 层构造函数返回 Biz 层接口（面向接口编程）**
4. **生成文件禁止手改**：`wire_gen.go`

## ProviderSet 的组织建议

- `internal/biz/`：聚合 UseCase ProviderSet
- `internal/data/`：聚合 Repo ProviderSet
- `internal/service/`：聚合 Service ProviderSet
- `cmd/server/`：顶层 `wire.go` 负责把 server/service/biz/data 装配成 app

## Go Demo（结构示例）

### 1) Data 构造函数返回 Biz 接口

```go
package data

import (
	"github.com/go-kratos/kratos/v2/log"
)

// RepoInterface 模拟 biz 层接口（真实项目里定义在 biz 包）。
type RepoInterface interface {
	Ping() error
}

type repoImpl struct {
	log *log.Helper
}

// NewRepo 推荐：返回接口类型，而不是实现类型。
func NewRepo(logger log.Logger) RepoInterface {
	return &repoImpl{log: log.NewHelper(logger)}
}

func (r *repoImpl) Ping() error { return nil }
```

### 2) ProviderSet 聚合

```go
package data

import "github.com/google/wire"

var ProviderSet = wire.NewSet(
	NewRepo,
)
```

### 3) 顶层 wire.go（仅示意）

```go
//go:build wireinject

package main

import (
	"github.com/google/wire"
	"linksoft.cn/node/internal/biz"
	"linksoft.cn/node/internal/data"
	"linksoft.cn/node/internal/server"
	"linksoft.cn/node/internal/service"
)

func initApp() (*App, error) {
	panic(wire.Build(
		server.ProviderSet,
		service.ProviderSet,
		biz.ProviderSet,
		data.ProviderSet,
	))
}
```

## 常见陷阱

- 构造函数返回具体实现类型，导致跨层耦合、难以 mock
- 手改 `wire_gen.go`，后续生成被覆盖
- ProviderSet 循环依赖（通常是分层边界被破坏）
