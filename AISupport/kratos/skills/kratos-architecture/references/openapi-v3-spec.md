# OpenAPI V3 Spec

## 规则

- 对外 HTTP proto 交付时，OpenAPI v3 注解应与 `google.api.http`、message 字段和校验语义保持一致
- 文档注解缺失、标题不统一、字段描述与真实校验不一致，都视为不完整交付
- 组合场景下优先保持注解与实际响应结构一致，不为文档展示方便篡改协议语义
- 常见错误模式应直接修正文档源定义，不手工维护与 proto 脱节的外部说明

## 注解覆盖层级

| 条件 | 做法 |
|------|------|
| 对外 HTTP proto | 同步补齐 document、operation、schema、property 注解 |
| 只有 `google.api.http` | 视为路由已定义，但文档未完整交付 |

```proto
// ✅ 文件、RPC、message、字段都有 OpenAPI v3 注解
option (openapi.v3.document) = {
  info: { title: "账户管理接口" version: "v1" description: "账户管理相关接口" }
  tags: { name: "账户管理" description: "账户查询与操作" }
};

service AccountService {
  rpc GetAccount(GetAccountRequest) returns (GetAccountReply) {
    option (google.api.http) = { get: "/admin/v1/accounts/{id}" };
    option (openapi.v3.operation) = {
      tags: "账户管理"
      summary: "查询账户详情"
      description: "查询指定账户详情"
      operation_id: "AccountService_GetAccount"
    };
  }
}

message GetAccountRequest {
  option (openapi.v3.schema) = {
    title: "账户管理-查询详情-请求"
    description: "查询详情请求"
  };
  uint32 id = 1 [(openapi.v3.property) = {description: "账户ID"}];
}
```

```proto
// ❌ 只有 http 路由，没有文档注解
service AccountService {
  rpc GetAccount(GetAccountRequest) returns (GetAccountReply) {
    option (google.api.http) = { get: "/admin/v1/accounts/{id}" };
  }
}
```

---

## 中文展示与标题约定

| 条件 | 做法 |
|------|------|
| 文档展示字段 | 使用中文标题、标签和描述 |
| 稳定系统标识 | `operation_id` 保持英文稳定值 |

```proto
// ✅ 中文展示 + 英文 operation_id
option (openapi.v3.document) = {
  info: { title: "库存券商安全管理接口" version: "v1" description: "库存券商安全相关接口" }
  tags: { name: "库存券商安全管理" description: "库存券商安全相关接口" }
};

option (openapi.v3.operation) = {
  tags: "库存券商安全管理"
  summary: "分页查询库存券商安全列表"
  description: "分页查询库存券商安全数据"
  operation_id: "InventorySecurityService_PageListInventorySecurity"
};
```

```proto
// ❌ 展示名继续沿用内部英文代号
option (openapi.v3.document) = {
  info: { title: "InventorySecurityService API" }
  tags: { name: "inventory-security" }
};
```

```proto
// ⚠️ summary 可以简短，但 title/tag 要保持同一中文模块前缀
option (openapi.v3.schema) = {
  title: "账户管理-查询详情-响应"
  description: "查询详情响应"
};
```

---

## 字段描述与校验对齐

| 条件 | 做法 |
|------|------|
| 字段有枚举、范围、示例、格式约束 | 写进 `property.description` 或对应文档字段 |
| 字段有 `validate.rules` | 文档说明必须和校验一致 |

```proto
// ✅ 描述包含语义和值域
uint32 security_type = 1 [
  (validate.rules).uint32 = {gte: 1, lte: 3},
  (openapi.v3.property) = {description: "券商安全类型 1-普通 2-重点 3-禁用"}
];

string broker_code = 2 [
  (validate.rules).string = {min_len: 2, max_len: 32},
  (openapi.v3.property) = {description: "券商编码 示例 CICC"}
];
```

```proto
// ❌ 有校验没文档语义
uint32 security_type = 1 [(validate.rules).uint32 = {gte: 1, lte: 3}];
```

---

## 组合场景

```proto
import "google/api/annotations.proto";
import "openapi/v3/annotations.proto";
import "validate/validate.proto";

option (openapi.v3.document) = {
  openapi: "3.0.3"
  info: { title: "账户管理接口" version: "v1" description: "账户管理相关接口" }
  tags: { name: "账户管理" description: "账户查询与操作" }
};

service AccountService {
  rpc GetAccount(GetAccountRequest) returns (GetAccountReply) {
    option (google.api.http) = { get: "/admin/v1/accounts/{id}" };
    option (openapi.v3.operation) = {
      tags: "账户管理"
      summary: "查询账户详情"
      description: "查询指定账户详情"
      operation_id: "AccountService_GetAccount"
    };
  }
}

message GetAccountRequest {
  option (openapi.v3.schema) = {
    title: "账户管理-查询详情-请求"
    description: "查询详情请求"
  };
  uint32 id = 1 [
    (validate.rules).uint32 = {gt: 0},
    (openapi.v3.property) = {description: "账户ID"}
  ];
}
```

这个组合场景同时满足：

- 文档层级完整
- 中文展示统一
- `operation_id` 稳定
- 字段说明与校验一致

---

## 常见错误模式

```proto
// ❌ 只写 google.api.http
option (google.api.http) = { ... };
```

```proto
// ❌ 标签和标题混用英文内部名
tags: { name: "inventory-security" }
```

```proto
// ❌ schema.title 沿用 message 英文名
title: "GetAccountReply"
```
