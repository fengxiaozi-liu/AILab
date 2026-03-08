# Proto Reference

## 这个主题解决什么问题
说明 proto 文件如何按 side、聚合根和 message 结构组织，以及如何和聚合根派生出来的应用层输入输出对象对齐。

## 适用场景

- 新增或修改对外协议
- 设计 RPC、Request、Reply、嵌套 message
- 检查 proto 与聚合根、应用层输入输出对象是否同构

## 设计意图

Proto 负责表达协议层契约，不负责反向定义领域模型。

- proto 文件和 message 应围绕聚合根创建。
- `proto message` 与应用层输入输出对象都是聚合根的投影。
- `proto message` 是协议层投影，应用层输入输出对象是应用层投影。
- 两者平级，不是哪个指导哪个。

## 实施提示

- 先明确对外能力，再定义 message 和 service。
- 优先让请求响应表达业务语义，而不是贴着数据库结构命名。
- 如果 proto 文件名、service 名或 message 名无法稳定对应到聚合根，应先回看 domain 建模。

## 推荐结构

- 按 side 组织 proto：`admin`、`open`、`inner`
- 基础公共协议放在 `base/*`，供各 side 复用
- 聚合根为主组织 RPC 和 message
- 局部子结构优先作为当前 RPC 的嵌套 message

## 基于聚合根创建 Proto

创建顺序通常是：

1. 先识别实体与聚合根
2. 再基于聚合根定义应用层输入输出对象
3. 再基于同一个聚合根定义 `proto` 文件、`service`、`message`
4. 最后补充局部嵌套 message 或跨聚合引用

常见映射关系：

- 一个聚合根通常对应一个主 `proto` 文件
- 一个主 `proto` 文件下可以有多个围绕该聚合根的 RPC
- `Reply` 内的核心对象优先围绕聚合根对象命名
- 从属实体优先作为 `{Entity}Info` 或 `{Entity}List` 出现在聚合根 reply 中

## Side 引用边界

- `admin` 侧 proto 只允许引用 `admin` 侧 proto 和 `base/*`。
- `open` 侧 proto 只允许引用 `open` 侧 proto 和 `base/*`。
- `admin` 与 `open` 之间禁止相互引用。
- `admin`、`open` 禁止引用 `inner`。
- `inner` 只允许引用 `inner` 和 `base/*`。

`base/*` 属于基础公共协议层，可被 `admin`、`open`、`inner` 复用，例如 `Paging`、`Sort`、`TimeRange`、`TransField`。`base/*` 不承载具体业务聚合 message，也不参与业务 side 隔离冲突判定。

如果需要跨 side 复用业务结构，应在各自 side 内定义稳定 message，并在实现层完成字段映射，不通过跨 side import 建立协议依赖。

## 公共基础协议

- `base/*` 只承载跨业务复用的稳定基础结构。
- `base/*` 不承载具体业务聚合、业务主 message 或业务 side RPC。
- 业务聚合仍定义在各自 `api/{business}/{side}/v1/*.proto` 中。

## 注释说明

Proto 文件默认不添加说明性注释。

- 文件主题、模块边界、服务分组通过目录、文件名、service 名和 message 结构表达
- message 语义优先通过稳定命名表达，不通过注释补充
- 不写字段注释、不写章节注释、不写装饰性注释

## RPC 命名模式

```text
GetAccount
ListAccount
CreateAccount
UpdateAccount
```

## Message 设计模式

```proto
message GetAccountReply {
  Account account_info = 1;
}

message Account {
  uint32 id = 1;
  string name = 2;
  repeated AccountCollectInfo account_collect_list = 3;
}
```

## Message 定义层级

- 聚合根主对象使用顶层 message，例如 `Account`。
- 在同一聚合根内可复用的从属实体，使用同文件顶层 message，例如 `AccountCollectInfo`。
- 只服务于单个 RPC 的局部结构，优先定义为当前 `Request` 或 `Reply` 的嵌套 message。
- 不要为了少量字段复用，跨 side 或跨聚合根 import 其他 proto 的 message。

## 代码示例参考

```proto
service AccountService {
  rpc GetAccount(GetAccountRequest) returns (GetAccountReply);
}

message GetAccountRequest {
  uint32 id = 1;
}

message GetAccountReply {
  Account account_info = 1;
}
```

## 管理侧分页协议示例

```proto
import "google/api/annotations.proto";
import "validate/validate.proto";
import "base/business/v1/business.proto";

service AccountService {
  rpc PageListAccount(PageListAccountRequest) returns (PageListAccountReply) {
    option (google.api.http) = { post: "/admin.v1/open/account/list/page" body: "*" };
  };
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
```

## 管理侧嵌套 message 示例

```proto
message ReviewAccountRequest {
  message RejectPageItem {
    string page_code = 1 [(validate.rules).string = {min_len: 1}];
    string reason_text = 2;
  }

  uint32 id = 1 [(validate.rules).uint32 = {gte: 1}];
  uint32 action = 2 [(validate.rules).uint32 = {gte: 1, lte: 2}];
  repeated RejectPageItem reject_page_list = 3;
}
```

当 request 或 message 已经明确处在聚合根语境里时，主对象 ID 直接使用 `id`。
只有跨聚合根引用、关联字段或过滤条件需要区分对象归属时，才使用 `account_id` 这类命名。

## 开放侧基础协议复用示例

```proto
import "base/business/v1/business.proto";

message GetLicenseeReply {
  string logo = 1;
  string code = 2;
  base.business.v1.TransField name = 3;
}
```

## Proto 与应用层输入输出对象对齐

- 聚合根结构与应用层输入输出对象使用同一业务语义
- relation 字段、嵌套结构、列表层级尽量一一对应
- 字段顺序更适合按核心字段、时间字段、关系字段阅读

## 文件创建联动

当聚合根发生新增或拆分时，proto 的调整通常与应用层输入输出对象同步发生，并先于 `service`、`server` 和 `wire`。

推荐联动顺序：

1. 先调整实体、聚合根和 Biz 对象
2. 再平级调整应用层输入输出对象与 `proto` 文件、message、service
3. 再更新 service 实现与 server 注册
4. 最后更新 provider / wire 和生成物

## 常见坑

- 先写 proto，再倒推聚合根
- 为跨聚合复用 message 而引入复杂 import
- `Reply` 结构把从属实体拍平成多个并列字段，导致聚合层级不清
- `proto message` 机械复制数据库结构，而不是聚合根语义

## 相关 Rule

- `../rules/proto-rule.md`
- `../../kratos-domain/rules/naming-rule.md`
- `../../kratos-conventions/rules/comment-rule.md`
