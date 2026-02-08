# Proto API 规范

## 概述

本项目使用 Protocol Buffers (protobuf) 定义 gRPC API 接口。

## 文件组织

### 目录结构

```
api/
├── admin/           # 管理后台相关 API
│   ├── admin/       # 管理端接口
│   ├── inner/       # 服务间调用接口
│   └── open/        # C端接口
├── base/            # 基础公共 API
│   ├── business/    # 业务公共消息定义
│   └── example/     # 示例模块
├── gateway/         # 网关相关 API
│   ├── admin/       # 管理端接口
│   ├── inner/       # 服务间调用接口
│   └── open/        # C端接口
├── system/          # 系统级 API
│   ├── admin/       # 管理端接口
│   ├── inner/       # 服务间调用接口
│   └── open/        # C端接口
└── user/            # 用户相关 API
    ├── admin/       # 管理端接口
    ├── inner/       # 服务间调用接口
    └── open/        # C端接口
```

### 命名规则

- 文件名：小写下划线，如 `order_item.proto`
- 包名：小写点分隔，如 `user.inner.v1`
- go_package：使用项目模块路径

## Proto 文件模板

### admin/open
- http路由，需要根据 `/业务域.版本/服务名/接口` 展开
- 如果有公共消息，引用 `base/business/v1/business.proto`
- 如果有关联消息，需要定义内联message
- 不能直接引用其他业务模块的消息

```protobuf
syntax = "proto3";

package system.admin.v1;

import "google/api/annotations.proto";
import "validate/validate.proto";
import "base/business/v1/business.proto";

option go_package = "linksoft.cn/trader/internal/api/system/admin/v1;v1";

service CounterChannelAccountService {
  rpc GetCounterChannelAccount (GetCounterChannelAccountRequest) returns (CounterChannelAccount) {
    option (google.api.http) = {
      post: "/admin.v1/system/counter/channel/account/info"
      body: "*"
    };
  };
  rpc ListCounterChannelAccount (ListCounterChannelAccountRequest) returns (ListCounterChannelAccountReply) {
    option (google.api.http) = {
      post: "/admin.v1/system/counter/channel/account/list"
      body: "*"
    };
  };
  rpc PageListCounterChannelAccount (PageListCounterChannelAccountRequest) returns (PageListCounterChannelAccountReply) {
    option (google.api.http) = {
      post: "/admin.v1/system/counter/channel/account/list/page"
      body: "*"
    };
  };
  rpc CountCounterChannelAccount (CountCounterChannelAccountRequest) returns (CountCounterChannelAccountReply) {
    option (google.api.http) = {
      post: "/admin.v1/system/counter/channel/account/count"
      body: "*"
    };
  };
  rpc CreateCounterChannelAccount (CreateCounterChannelAccountRequest) returns (CounterChannelAccount) {
    option (google.api.http) = {
      post: "/admin.v1/system/counter/channel/account/create"
      body: "*"
    };
  };
  rpc UpdateCounterChannelAccount (UpdateCounterChannelAccountRequest) returns (base.business.v1.BlankReply) {
    option (google.api.http) = {
      post: "/admin.v1/system/counter/channel/account/update"
      body: "*"
    };
  };
}

message GetCounterChannelAccountRequest{
  uint32 id = 1 [(validate.rules).uint32.gte = 1];
}

message CreateCounterChannelAccountRequest {
  string account = 1;
  string name = 2;
  uint32 counter_channel_id = 3;
  string scene = 4;
}

message UpdateCounterChannelAccountRequest {
  string account = 1;
  string name = 2;
  uint32 counter_channel_id = 3;
  uint32 id = 4 [(validate.rules).uint32.gte = 1];
  string scene = 5;
}

message DeleteCounterChannelAccountRequest {
  uint32 id = 1[(validate.rules).uint32.gte = 1];
}

message ListCounterChannelAccountRequest {
  string account = 1;
  uint32 counter_channel_id = 2;
  uint32 status = 3;
  string name = 5;
  base.business.v1.TimeRange create_time = 6;
  base.business.v1.Paging paging = 7;
  base.business.v1.Sort sort = 8;
}

message PageListCounterChannelAccountRequest {
  string account = 1;
  uint32 counter_channel_id = 2;
  uint32 status = 3;
  string name = 5;
  base.business.v1.TimeRange create_time = 6;
  base.business.v1.Paging paging = 7;
  base.business.v1.Sort sort = 8;

}

message PageListCounterChannelAccountReply {
  repeated CounterChannelAccount list = 1;
  int32 count = 2;
}

message CountCounterChannelAccountRequest {
  string account = 1;
  uint32 counter_channel_id = 2;
  uint32 status = 3;
  string name = 5;
  base.business.v1.TimeRange create_time = 6;
  base.business.v1.Paging paging = 7;
  base.business.v1.Sort sort = 8;
}

message CountCounterChannelAccountReply {
  int32 count = 1;
}

message CounterChannelAccount {
  message CounterChannel{
    uint32 id = 1;
    string name = 2;
    uint32 status = 4;
    uint32 health_status = 5;
    uint32 create_time = 6;
  }
  uint32 id = 1;
  string account = 2;
  string name = 3;
  uint32 create_time = 5;
  uint32 update_time = 6;
  uint32 status = 7;
  uint32 counter_channel_id = 11;
  string scene = 12;
  CounterChannel counter_channel_info = 13;
}

message ListCounterChannelAccountReply {
  repeated CounterChannelAccount list = 1;
}

```

### inner
- 不需要 http 路由
- 如果有公共消息，引用 `base/business/v1/business.proto`
- 关联message可以直接引用其他inner模块的message

```protobuf
syntax = "proto3";

package system.inner.v1;

import "base/business/v1/business.proto";
import "system/inner/v1/counter_channel.proto";

option go_package = "linksoft.cn/trader/internal/api/system/inner/v1;v1";

service CounterChannelAccountService {
  rpc GetCounterChannelAccount (GetCounterChannelAccountRequest) returns (CounterChannelAccount);
  rpc GetCounterChannelAccountByAccount (GetCounterChannelAccountByAccountRequest) returns (CounterChannelAccount);
  rpc ListCounterChannelAccount (ListCounterChannelAccountRequest) returns (ListCounterChannelAccountReply);
  rpc MapCounterChannelAccount (MapCounterChannelAccountRequest) returns (MapCounterChannelAccountReply);
}

message CounterChannelAccount {
  uint32 id = 1;
  string account = 2;
  string name = 3;
  uint32 create_time = 5;
  uint32 update_time = 6;
  uint32 status = 7;
  uint32 counter_channel_id = 11;
  string scene = 12;
  CounterChannel counter_channel_info = 13;
}

message GetCounterChannelAccountRequest{
  uint32 id = 1;
  base.business.v1.FilterConfig filter_config = 2;
}

message GetCounterChannelAccountByAccountRequest{
  string account = 1;
  base.business.v1.FilterConfig filter_config = 2;
}

message MapCounterChannelAccountRequest {
  repeated uint32 id_list = 1;
  string account = 2;
  repeated string account_list = 3;
  uint32 counter_channel_id = 4;
  uint32 status = 5;
  string name = 6;
  base.business.v1.TimeRange create_time = 7;
  base.business.v1.Paging paging = 8;
  base.business.v1.Sort sort = 9;
  base.business.v1.FilterConfig filter_config = 10;
}
message MapCounterChannelAccountReply {
  map<uint32, CounterChannelAccount> map = 1;
}

message ListCounterChannelAccountRequest {
  repeated uint32 id_list = 1;
  string account = 2;
  repeated string account_list = 3;
  uint32 counter_channel_id = 4;
  uint32 status = 5;
  string name = 6;
  base.business.v1.TimeRange create_time = 7;
  base.business.v1.Paging paging = 8;
  base.business.v1.Sort sort = 9;
  base.business.v1.FilterConfig filter_config = 10;

}
message ListCounterChannelAccountReply {
  repeated CounterChannelAccount list = 1;
}

```

## 公共消息类型

`base/business/v1/business.proto` 定义了公共消息：

```protobuf
// 多语言字段
message TransField {
  string zh_CN = 1 [json_name = "zh-CN"];
  string tc = 2;
  string en = 3;
}

// 时间范围
message TimeRange {
  string start = 1;
  string end = 2;
}

// 分页参数
message Paging {
  uint32 page = 1;
  uint32 size = 2;
}

// 排序参数
message Sort {
  string order = 1;   // asc/desc
  string field = 2;   // 排序字段
}

// 分组参数
message GroupBy {
  repeated string fields = 1;
}

// 比较条件
message Compare {
  uint32 operator = 1;  // 操作符
  double value = 2;     // 比较值
}

// 空请求
message BlankRequest {}

// 空响应
message BlankReply {}
```

## proto规范

### 设计规范

- 例如 `Order` 模块不能包含操作 `OrderCharge` 的接口。`OrderService` 只管 `Order` 元数据，`OrderChargeService` 管 `OrderCharge`。
- `Order` 是主动方（聚合根），`OrderCharge` 是被动方。`OrderCharge Message` 可以包含 `order_id`，但 `OrderCharge` 服务不能反向操作 `Order`。
- 逻辑拆分了，物理文件也要拆分。`order.proto` 和 `order_charge.proto` 必须分离成两个文件。
- 应提供通用的资源查询接口，通过参数控制。例如 `ListOrder` 支持按状态、时间范围等过滤，而不是提供 `ListPendingOrders`、`ListCompletedOrders` 等多个接口。
- 避免过度设计。不要为了“可能的将来需求”而设计复杂的过滤器或查询条件。
- 消息设计应以实际业务需求为导向，避免冗余字段和复杂嵌套。
- 保持消息简洁明了，易于理解和使用。
- 遵循一致的命名规范，确保代码可读性和维护性。
- 优先使用基础类型，避免不必要的嵌套消息，优先复用 `base` 模块。
- 确保消息的向前和向后兼容性，避免破坏现有客户端。

### 惯用法
- 单条 Create/Update: 输入数据**明确展开字段**而不复用实体 Message（解耦输入与存储结构），Create 接口直接返回主结构体 `Order`，Update接口直接返回`BlankReply`。
- 批量 Create/Update: 输入数据复用实体 Message（`repeated Order`）以减少冗余定义，批量Create返回值使用 `ListOrderReply`，批量Update返回值使用 `BlankReply`。
- time类型字段统一使用 `uint32` 表示 Unix 时间戳（秒）。
- 金额字段使用 `double` 表示，避免精度问题。
- 枚举类型优先使用 `uint32`，避免字符串枚举带来的问题，除非枚举本身就是字符串（如状态码）。
- JSON 对象使用嵌套 message 定义，避免使用字符串存储 JSON。
- 列表使用 `repeated` 定义，确保类型一致性。
- 时间类的查询条件默认要使用 `TimeRange`，避免单点时间查询带来的歧义。
- 分页查询优先使用 `Paging`，避免一次性返回大量数据。
- 排序使用 `Sort` 定义，确保排序字段和顺序明确。
- 分组使用 `GroupBy` 定义，支持多字段分组。
- 比较条件使用 `Compare` 定义，支持多种比较操作。
- 关联查询配置使用 `FilterConfig` 定义，支持灵活的关联查询配置。
- 不能确定值域大小的字段，如`amount`，`extra`之类，除非用户明确要求，否则不能作为查询条件。

### 字段类型规范

| 业务类型 | Proto 类型            | 说明          |
|---------|---------------------|-------------|
| 主键 ID | `uint32`            | 无符号整数       |
| 字符串 | `string`            | UTF-8 编码    |
| 布尔值 | `bool`              | -           |
| 时间戳 | `uint32`            | Unix 时间戳（秒） |
| 金额 | `double`            | 双精度浮点数      |
| 枚举 | `uint32` 或 `string` | 推荐使用 uint32 |
| JSON 对象 | 嵌套 message          | -           |
| 列表 | `repeated`          | -           |

### RPC 方法命名规范

| 操作类型 | 命名格式 | 示例 |
|---------|---------|------|
| 单条查询 | `Get{Entity}` | `GetExample` |
| 列表查询 | `List{Entity}` | `ListExample` |
| 分页查询 | `PageList{Entity}` | `PageListExample` |
| 计数查询 | `Count{Entity}` | `CountExample` |
| 创建 | `Create{Entity}` | `CreateExample` |
| 更新 | `Update{Entity}` | `UpdateExample` |
| 批量删除 | `BatchDelete{Entity}` | `BatchDeleteExample` |

## 代码生成

执行以下命令生成 Go 代码：

```shell
make api
```

生成的代码位于 `internal/api/` 目录下。

## 版本管理

- 使用 `v1`, `v2` 等版本号目录
- 不兼容变更需要升级版本
- 旧版本需要维护到客户端迁移完成
