# Proto Spec

## 创建顺序

先做聚合根建模，再做 proto，不要反过来。

| 步骤 | 做什么 |
|------|--------|
| 1 | 识别实体与聚合根 |
| 2 | 定义应用层输入输出对象（Biz 对象） |
| 3 | 基于同一聚合根定义 proto 文件、service、message |
| 4 | 补充嵌套 message 或跨聚合引用 |
| 5 | 更新 service 实现与 server 注册 |
| 6 | 更新 provider / wire，执行 codegen |

---

## Side 引用边界

决策表：允许 import 哪些 side。

| 当前 side | 可引用 | 禁止引用 |
|-----------|--------|----------|
| `admin` | `admin/*`、`base/*` | `open/*`、`inner/*` |
| `open` | `open/*`、`base/*` | `admin/*`、`inner/*` |
| `inner` | `inner/*`、`base/*` | `admin/*`、`open/*` |
| `base/*` | 无业务 side | — |

```proto
// ✅ open 引用 base
import "base/business/v1/business.proto";

// ❌ open 引用 admin —— 跨 side 禁止
import "admin/account/v1/account.proto";

// ❌ admin 引用 inner
import "inner/order/v1/order.proto";
```

跨 side 需要复用业务结构时：在各自 side 内单独定义，实现层做字段映射，不通过 import 建立协议依赖。

---

## Message 层级定义

决策表：结构放顶层还是嵌套。

| 条件 | 做法 |
|------|------|
| 聚合根主对象 | 顶层 message |
| 同文件多个 RPC 共用的从属实体 | 顶层 message，命名 `{Entity}Info` |
| 只服务于 1 个 RPC 的局部结构 | 嵌套 message |
| 跨 side / 跨聚合根复用 | 各自 side 内单独定义，不跨 import |

```proto
// ✅ 聚合根主对象 → 顶层
message Account {
  uint32 id = 1;
  string name = 2;
  repeated AccountCollectInfo account_collect_list = 3;
}

// ✅ 从属实体多 RPC 复用 → 顶层
message AccountCollectInfo {
  uint32 id = 1;
  string title = 2;
}

// ✅ 只服务于单个 RPC 的局部结构 → 嵌套
message ReviewAccountRequest {
  message RejectPageItem {
    string page_code = 1 [(validate.rules).string = {min_len: 1}];
    string reason_text = 2;
  }
  uint32 id = 1 [(validate.rules).uint32 = {gte: 1}];
  uint32 action = 2 [(validate.rules).uint32 = {gte: 1, lte: 2}];
  repeated RejectPageItem reject_page_list = 3;
}

// ❌ 只服务于单个 RPC 的局部结构 → 平铺顶层，错误
message RejectPageItem { ... }
```

边界案例：同一个局部结构被 2 个以上 RPC 引用时，提升为顶层 message。

```proto
// ⚠️ Filter 只被 ListAccount 用 → 嵌套
message ListAccountRequest {
  message Filter { uint32 status = 1; }
  Filter filter = 1;
}

// ⚠️ Filter 同时被 ListAccount 和 SearchAccount 用 → 提升为顶层
message AccountFilter { uint32 status = 1; }
message ListAccountRequest   { AccountFilter filter = 1; }
message SearchAccountRequest { AccountFilter filter = 1; }
```

---

## ID 字段命名

```proto
// ✅ 处于聚合根语境内 → 直接用 id
message GetAccountRequest {
  uint32 id = 1;
}

// ✅ 跨聚合引用 → 用 {aggregate}_id
message GetOrderRequest {
  uint32 account_id = 1;  // 引用了 Account 聚合根
}

// ❌ 当前语境已是 Account，仍追加聚合根前缀
message GetAccountRequest {
  uint32 account_id = 1;  // 冗余
}
```

---

## RPC 命名

```proto
// ✅
rpc GetAccount(GetAccountRequest) returns (GetAccountReply);
rpc ListAccount(ListAccountRequest) returns (ListAccountReply);
rpc CreateAccount(CreateAccountRequest) returns (CreateAccountReply);
rpc UpdateAccount(UpdateAccountRequest) returns (UpdateAccountReply);
rpc PageListAccount(PageListAccountRequest) returns (PageListAccountReply);

// ❌ 动词不一致、语义模糊
rpc QueryAccount(...)
rpc FetchAccountList(...)
rpc AccountCreate(...)
```

---

## 注释规范

```proto
// ✅ 不写任何注释，用命名和结构表达语义
message Account {
  uint32 id = 1;
  string name = 2;
}

// ❌ 写装饰性或说明性注释
// Account 账户信息
message Account {
  uint32 id = 1;   // 账户 ID
  string name = 2; // 名称
}
```

---

## 组合场景：管理侧分页接口完整示例

涵盖：side 引用、base 复用、RPC 命名、ID 命名、Reply 结构。

```proto
syntax = "proto3";
package admin.account.v1;

import "google/api/annotations.proto";
import "validate/validate.proto";
import "base/business/v1/business.proto";

service AccountService {
  rpc GetAccount(GetAccountRequest) returns (GetAccountReply) {
    option (google.api.http) = { get: "/admin.v1/account/{id}" };
  };
  rpc PageListAccount(PageListAccountRequest) returns (PageListAccountReply) {
    option (google.api.http) = { post: "/admin.v1/account/list/page" body: "*" };
  };
  rpc ReviewAccount(ReviewAccountRequest) returns (ReviewAccountReply) {
    option (google.api.http) = { post: "/admin.v1/account/review" body: "*" };
  };
}

// 聚合根主对象 → 顶层
message Account {
  uint32 id = 1;
  string name = 2;
  repeated AccountCollectInfo account_collect_list = 3;
}

// 从属实体多 RPC 复用 → 顶层
message AccountCollectInfo {
  uint32 id = 1;
  string title = 2;
}

message GetAccountRequest {
  uint32 id = 1 [(validate.rules).uint32 = {gte: 1}];
}
message GetAccountReply {
  Account account_info = 1;
}

message PageListAccountRequest {
  repeated uint32 open_status = 1;
  string user_code = 2;
  base.business.v1.TimeRange create_time = 10;
  base.business.v1.Paging paging = 13;
  base.business.v1.Sort sort = 14;
}
message PageListAccountReply {
  repeated Account list = 1;
  int32 count = 2;
}

// 只服务于单个 RPC 的局部结构 → 嵌套
message ReviewAccountRequest {
  message RejectPageItem {
    string page_code = 1 [(validate.rules).string = {min_len: 1}];
    string reason_text = 2;
  }
  uint32 id = 1 [(validate.rules).uint32 = {gte: 1}];
  uint32 action = 2 [(validate.rules).uint32 = {gte: 1, lte: 2}];
  repeated RejectPageItem reject_page_list = 3;
}
message ReviewAccountReply {}
```

---

## 常见错误模式

```proto
// ❌ Reply 把从属实体拍平，丢失聚合层级
message GetAccountReply {
  uint32 id = 1;
  string name = 2;
  uint32 collect_id = 3;       // 应该是 repeated AccountCollectInfo
  string collect_title = 4;
}

// ❌ message 机械复制数据库字段，而非聚合根语义
message Account {
  uint32 id = 1;
  string create_by = 2;        // 数据库审计字段，不属于协议语义
  string update_by = 3;
  int64 deleted_at = 4;
}

// ❌ 为复用少量字段跨 side import
import "open/account/v1/account.proto";  // admin side 引用 open side
```
