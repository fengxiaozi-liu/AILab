# Proto Spec

## 文件创建判定（何时新建 .proto 文件）

| 判定条件 | 做法 |
|------|------|
| **引入了全新的核心聚合根/业务实体**（如：新增了 `Invoice` 或 `Coupon` 模块） | **新建 proto 文件**。通常按被操作实体的下划线命名法处理，如 `invoice.proto`。 |
| **同一个实体，但是需要在另外的端（Side）暴露一组全新接口**（如新增了 Admin 后台独有的用户列表功能） | **在对应端目录下新建配套文件**（若该实体在该目录下尚不存在）。如 `api/admin/v1/user.proto`。 |
| **只是为已有的聚合增加一个特例方法、页面动作或临时查询节点**（如 `GetUserInfo`, `ExportUserList`） | **绝对禁止新建**！直接向此端现存的该实体 proto 文件（如 `user.proto`）中追加 Service 动作与 Message 模型。 |
| 定义跨越多个聚合或多个端的纯公共数据结构（如 `Paging`, `ErrorReason`） | **必须新建**。归入 `api/base/v1` 等公用业务下的统一结构体文件中。 |

### ❌ 反例：滥建文件（MUST NOT）
```proto
// 错误：按孤立的动作方法去建文件，碎片化严重
api/open/v1/user_login.proto
api/open/v1/user_logout.proto

// 错误：按前端给定的纯 UI 页面视角去堆砌文件
api/admin/v1/user_export_page.proto 
api/admin/v1/user_detail_page.proto
```

### ✅ 正例：按单一核心实体与端（Side）隔离（MUST）
```proto
// 正确：在一个文件内，囊括所有 C 端开放接口涉及的 User 操作（Login, Logout, Detail 等全汇聚于此）
api/open/v1/user.proto

// 正确：涵盖所有内部 B 端专有的 User 管理动作（Export, Auditing）
api/admin/v1/user.proto
```

---

## Side 边界与引用

| 条件 | 做法 |
|------|------|
| `admin`、`open`、`inner` 各自协议 | 只引用本 side 和 `base/*` |
| 多个 side 共用协议 | 提取到 `base/*` |

```text
// ✅
api/open/v1/account.proto
api/admin/v1/account.proto
api/base/v1/paging.proto
```

```proto
// ❌ side 互引业务 proto
import "api/open/v1/account.proto";   // in admin proto
import "api/inner/v1/user.proto";     // in open proto
```

```proto
// ⚠️ 真正的公共结构才提升到 base/*
import "api/base/v1/paging.proto";
```

---

## 聚合根组织

| 条件 | 做法 |
|------|------|
| proto 文件、service、message | 围绕聚合根或稳定业务主题组织 |
| 只有单个接口动作视角 | 不作为文件主语义 |

```proto
// ✅ 围绕 Account 聚合根组织
service AccountService {
  rpc GetAccount(GetAccountRequest) returns (GetAccountReply);
}

message AccountInfo {
  uint32 id = 1;
  string status = 2;
}
```

```proto
// ❌ 围绕页面或动作组织
service AccountPageHandler {}
message GetAccountPageRow {}
```

---

## Message 结构与嵌套

| 条件 | 做法 |
|------|------|
| 子结构只服务单个 RPC | 嵌套在当前 request/reply 中 |
| 子结构跨多个 RPC 复用 | 提升为顶层 message 或 `base/*` |

```proto
// ✅ 只服务单个 RPC 的过滤条件嵌套
message ListAccountRequest {
  message Filter {
    uint32 status = 1;
  }
  Filter filter = 1;
}
```

```proto
// ❌ 本地子结构平铺成无必要顶层
message AccountListFilter {
  uint32 status = 1;
}

message ListAccountRequest {
  AccountListFilter filter = 1;
}
```

```proto
// ⚠️ 多 RPC 共用时提升为顶层
message AccountFilter {
  uint32 status = 1;
}
```

---

## 字段命名与语义

| 条件 | 做法 |
|------|------|
| 当前语境已是聚合根主对象 ID | 直接用 `id` |
| 跨聚合引用 | 用 `{aggregate}_id` |

```proto
// ✅
message GetAccountRequest {
  uint32 id = 1;
}

message AccountInfo {
  uint32 user_id = 1;
}
```

```proto
// ❌ 当前语境已是 Account，仍使用 account_id
message GetAccountRequest {
  uint32 account_id = 1;
}
```

---

## 组合场景

```proto
syntax = "proto3";

package open.v1;

option go_package = "xxx/internal/service/open/v1;v1";

service AccountService {
  rpc ListAccount(ListAccountRequest) returns (ListAccountReply);
}

message ListAccountRequest {
  message Filter {
    uint32 status = 1;
  }
  uint32 page = 1;
  uint32 page_size = 2;
  Filter filter = 3;
}

message ListAccountReply {
  repeated AccountInfo list = 1;
}

message AccountInfo {
  uint32 id = 1;
  uint32 user_id = 2;
}
```

这个组合场景同时满足：

- side 清晰
- service 和 message 围绕 `Account`
- 局部过滤条件使用嵌套 message
- ID 字段命名符合聚合语境

---

## 常见错误模式

```proto
// ❌ side 互引
import "api/open/v1/account.proto";
```

```proto
// ❌ 动作名污染主语义
service GetAccountHandler {}
```

```proto
// ❌ 字段命名冗余
uint32 account_id = 1;
```
