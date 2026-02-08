# 命名规范

## 概述

本文档定义了项目中各类命名的规范，包括文件、包、变量、函数、常量等。

## 文件命名

### Go 文件

| 类型 | 规则 | 示例 |
|------|-----|------|
| 普通文件 | 小写下划线 | `user.go` |
| 测试文件 | `*_test.go` | `user_test.go` |
| 生成文件 | `*_gen.go` | `wire_gen.go` |

### Proto 文件

| 类型 | 规则 | 示例 |
|------|-----|------|
| API 定义 | 小写下划线 | `user.proto` |
| 公共消息 | 小写下划线 | `business.proto` |

### 配置文件

| 类型 | 规则 | 示例 |
|------|-----|------|
| YAML | 小写下划线 | `config.yaml` |
| TOML | 小写点分隔 | `active.zh-CN.toml` |

## 包命名

### 规则

- 小写字母
- 单词之间不使用分隔符
- 简短且有意义
- 避免使用通用名称

```go
// 好的命名
package user
package order
package base

// 避免的命名
package util        // 太通用
package common      // 太通用
package userService // 不要使用驼峰
package user_error  // 不要使用下划线
```

### 常见包命名

| 用途 | 命名 | 示例 |
|------|-----|------|
| 业务层 | 业务名 | `user`, `order` |
| 枚举 | 业务名 | `user`, `order` |
| 错误 | 业务名 | `user`, `order` |
| API 版本 | `v` + 数字 | `v1`, `v2` |

## 变量命名

### 驼峰命名

```go
// 局部变量 - 小驼峰
var userID uint32
var orderList []*Order
var isActive bool

// 导出变量 - 大驼峰
var ProviderSet = wire.NewSet(...)
var DefaultTimeout = 30 * time.Second
```

### 常用缩写

| 缩写 | 全称 | 使用方式 |
|------|-----|---------|
| ID | Identifier | `userID`, `GetUserID()` |
| URL | Uniform Resource Locator | `imageURL` |
| API | Application Programming Interface | `apiVersion` |
| HTTP | HyperText Transfer Protocol | `httpClient` |
| JSON | JavaScript Object Notation | `jsonData` |

### 接收者命名

```go
// 使用类型名首字母小写
func (u *UserUseCase) Get(ctx context.Context, id uint32) (*User, error)
func (r *userRepo) GetUser(ctx context.Context, id uint32) (*User, error)
func (srv *UserService) GetUser(ctx context.Context, req *v1.GetUserRequest) (*v1.User, error)
```

补充约定（便于见名知意，且保持简短）：

- `r`：Repo/Data 实现
- `u`：UseCase
- `srv`：Service（对外接口实现）
- `s`：Server/状态类对象（视上下文）

避免：

- 在同一文件/同一包内把 `r` 同时用作 Repo 与 Service 的接收者（语义混乱）

## 函数命名

### 构造函数

```go
// New + 类型名
func NewUserUseCase(repo UserRepo) *UserUseCase
func NewUserRepo(data *kit.Data) UserRepo
func NewUserService(uc *UserUseCase) *UserService
```

### CRUD 方法

| 操作 | Service | UseCase | Repository |
|------|---------|---------|------------|
| 单条查询 | `GetUser` | `Get` | `GetUser` |
| 列表查询 | `ListUser` | `List` | `ListUser` |
| 分页查询 | `PageListUser` | `PageList` | `PageListUser` |
| 计数 | `CountUser` | `Count` | `CountUser` |
| 创建 | `CreateUser` | `Create` | `CreateUser` |
| 更新 | `UpdateUser` | `Update` | `UpdateUser` |
| 删除 | `DeleteUser` | `Delete` | `DeleteUser` |
| 批量删除 | `BatchDeleteUser` | `BatchDelete` | `BatchDeleteUser` |

### 布尔方法

```go
// Is/Has/Can/Should 开头
func (u *User) IsActive() bool
func (u *User) HasPermission(perm string) bool
func (o *Order) CanCancel() bool
func (p *Payment) ShouldRetry() bool
```

## 常量命名

### 普通常量

```go
// 全大写下划线分隔
const (
    MAX_RETRY_COUNT = 3
    DEFAULT_TIMEOUT = 30 * time.Second
)
```

### 枚举常量

```go
// 类型名 + 值名
const (
    UserStatusActive   UserStatus = 1
    UserStatusInactive UserStatus = 2
    
    OrderTypeNormal OrderType = "normal"
    OrderTypeRefund OrderType = "refund"
)
```

## 结构体命名

### 领域模型

```go
// 实体名称，大驼峰
type User struct {
    ID       uint32
    Name     string
    Email    string
    Status   UserStatus
}
```

### 过滤器

```go
// 实体名 + Filter
type UserFilter struct {
    IDList   []uint32
    Name     string
    Status   UserStatus
    Paging   filter.Paging
    Sort     filter.Sort
}
```

### 接口

```go
// 动词/名词 或 描述性名称
type UserRepo interface {
    GetUser(ctx context.Context, id uint32) (*User, error)
}

type Transaction interface {
    InTx(context.Context, func(ctx context.Context) error) error
}
```

补充约定（命名层面）：

- 接口类型名使用 `PascalCase`，依赖/仓储类接口常用 `XxxRepo` 后缀
- 禁止 `I` 前缀（如 `IUserRepo`）与 `Interface` 后缀（如 `UserRepoInterface`）

## Proto 命名

### Service

```protobuf
// 实体名 + Service
service UserService {
    rpc GetUser (GetUserRequest) returns (User) {};
}
```

### Message

```protobuf
// 操作 + 实体 + Request/Reply
message GetUserRequest {
    uint32 id = 1;
}

message ListUserReply {
    repeated User list = 1;
}

message PageListUserReply {
    repeated User list = 1;
    int32 count = 2;
}
```

### Field

```protobuf
// 小写下划线 (snake_case)
message User {
    uint32 id = 1;
    string user_name = 2;
    string email_address = 3;
    uint32 create_time = 4;
}
```

## 数据库命名

### 表名

```go
// 前缀 + 小写下划线
entsql.Annotation{Table: "link_user"}
entsql.Annotation{Table: "link_order_item"}
```

### 字段名

```go
// 小写下划线
field.String("user_name")
field.Uint32("create_time")
field.Uint32("order_id")
```

### 索引名

```go
// idx_ 或 uniq_ 前缀
index.Fields("status").StorageKey("idx_status")
index.Fields("user_id", "create_time").StorageKey("idx_userid_createtime")
index.Fields("code").Unique().StorageKey("uniq_code")
```

## 特殊命名约定

### 私有实现

```go
// 小写开头，实现接口
type userRepo struct {  // 私有，实现 UserRepo 接口
    data *kit.Data
}
```

### Wire Provider

```go
// ProviderSet 固定名称
var ProviderSet = wire.NewSet(
    NewUserUseCase,
    NewUserRepo,
)
```

### Context Key

```go
// 私有结构体类型
type localizerKey struct{}
type tenantKey struct{}

ctx = context.WithValue(ctx, localizerKey{}, localizer)
```

## 最佳实践

1. **一致性**：项目内保持一致的命名风格
2. **可读性**：名称要能清楚表达意图
3. **简洁性**：在不损失清晰度的前提下尽量简短
4. **避免缩写**：除了常用缩写外，避免自创缩写
5. **上下文**：名称要在其使用的上下文中有意义
