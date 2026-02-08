# 测试规范

## 概述

本文档定义 Kratos 微服务项目中单元测试和集成测试的编写规范。

## 文件命名与放置

### 单元测试

| 层 | 测试文件路径 | 命名规则 |
|----|------------|---------|
| Biz | `internal/biz/<domain>/<domain>_test.go` | 与被测文件同目录 |
| Data | `internal/data/<domain>/<domain>_test.go` | 与被测文件同目录 |
| Service | `internal/service/<domain>/<domain>_test.go` | 与被测文件同目录 |

### 集成测试

- 放置于 `internal/biz/<domain>/integration_test.go` 或独立的 `tests/` 目录
- 使用 `//go:build integration` 构建标签隔离

## Biz 层单测

Biz 层是业务逻辑核心，测试优先级最高。

### Mock 策略

- Mock biz 层定义的 Repo 接口（data 层实现）
- Mock depend 包中的 InnerRPC 依赖
- 推荐使用 `go.uber.org/mock/mockgen` 生成 mock

### Mock 生成命令

```bash
# 为 biz 层 Repo 接口生成 mock
mockgen -source=internal/biz/<domain>/repo.go -destination=internal/biz/<domain>/mock_repo_test.go -package=<domain>

# 为 depend 接口生成 mock
mockgen -source=internal/biz/depend/<depend>.go -destination=internal/biz/depend/mock_<depend>_test.go -package=depend
```

### 测试模式

```go
package <domain>

import (
    "context"
    "testing"

    "go.uber.org/mock/gomock"
)

func TestXxxUsecase_MethodName(t *testing.T) {
    // 表驱动测试
    tests := []struct {
        name    string
        input   InputType
        mock    func(ctrl *gomock.Controller) *MockRepo
        want    OutputType
        wantErr bool
    }{
        {
            name:  "正常场景",
            input: InputType{...},
            mock: func(ctrl *gomock.Controller) *MockRepo {
                m := NewMockRepo(ctrl)
                m.EXPECT().Method(gomock.Any(), gomock.Any()).Return(result, nil)
                return m
            },
            want:    expectedOutput,
            wantErr: false,
        },
        {
            name:  "异常场景-xxx",
            input: InputType{...},
            mock: func(ctrl *gomock.Controller) *MockRepo {
                m := NewMockRepo(ctrl)
                m.EXPECT().Method(gomock.Any(), gomock.Any()).Return(nil, errXxx)
                return m
            },
            wantErr: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            ctrl := gomock.NewController(t)
            defer ctrl.Finish()

            repo := tt.mock(ctrl)
            uc := NewXxxUsecase(repo, logger)

            got, err := uc.MethodName(context.Background(), tt.input)
            if (err != nil) != tt.wantErr {
                t.Errorf("MethodName() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if !tt.wantErr && got != tt.want {
                t.Errorf("MethodName() = %v, want %v", got, tt.want)
            }
        })
    }
}
```

## Data 层单测

### 策略

- 使用 Ent 官方的 `enttest` 包创建内存 SQLite 数据库
- 测试真实的数据库操作，不 mock ORM

### 测试模式

```go
package <domain>

import (
    "context"
    "testing"

    _ "github.com/mattn/go-sqlite3"
    "xxx/internal/data/ent/enttest"
)

func TestXxxRepo_MethodName(t *testing.T) {
    client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&_fk=1")
    defer client.Close()

    repo := NewXxxRepo(client, logger)

    tests := []struct {
        name    string
        setup   func(ctx context.Context) // 准备测试数据
        input   InputType
        want    OutputType
        wantErr bool
    }{
        {
            name: "正常查询",
            setup: func(ctx context.Context) {
                client.Xxx.Create().SetField("value").SaveX(ctx)
            },
            input:   InputType{...},
            want:    expectedOutput,
            wantErr: false,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            ctx := context.Background()
            if tt.setup != nil {
                tt.setup(ctx)
            }
            got, err := repo.MethodName(ctx, tt.input)
            if (err != nil) != tt.wantErr {
                t.Errorf("MethodName() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if !tt.wantErr && got != tt.want {
                t.Errorf("MethodName() = %v, want %v", got, tt.want)
            }
        })
    }
}
```

## Service 层单测

- 通常 service 层较薄（DTO 转换 + 调用 biz），可选择性测试
- 重点测试参数校验逻辑和 DTO 映射正确性
- Mock biz 层的 Usecase 接口

## 测试优先级

| 优先级 | 测试目标 | 说明 |
|--------|---------|------|
| P0 | Biz 层核心业务逻辑 | 分支多、规则复杂的方法 |
| P1 | Data 层复杂查询 | 多条件查询、聚合、事务 |
| P2 | Biz 层简单 CRUD | 直通型方法 |
| P3 | Service 层 DTO 映射 | 可选 |

## 命令

```bash
# 运行所有单元测试
go test ./internal/...

# 运行特定包测试
go test ./internal/biz/<domain>/...

# 运行测试并生成覆盖率
go test -coverprofile=coverage.out ./internal/...
go tool cover -html=coverage.out

# 运行集成测试
go test -tags=integration ./...
```

## 注意事项

- 测试文件与被测文件同包，可访问未导出成员
- 使用 `t.Parallel()` 标记可并行的测试
- 避免测试间共享状态
- Data 层测试每个用例使用独立事务或清理数据
- 测试方法命名：`Test<Struct>_<Method>`（如 `TestOrderUsecase_Create`）
