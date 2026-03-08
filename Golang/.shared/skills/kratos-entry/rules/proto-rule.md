# Proto Rule

## Principles

- Contract First，跨边界协议先定义 `.proto`。

## Specification

- HTTP/gRPC/Event 等对外结构先改 proto，再生成，再实现。
- 协议变化后同步检查枚举、错误码和状态语义。
- proto side 只使用 `admin`、`open`、`inner` 三类目录。
- `admin` 只允许引用 `admin` 和 `base/*`。
- `open` 只允许引用 `open` 和 `base/*`。
- `inner` 只允许引用 `inner` 和 `base/*`。
- `admin` 与 `open` 禁止相互引用。
- `admin`、`open` 禁止引用 `inner`。
- `base/*` 属于基础公共协议层，可被任意 side 引用，但不参与业务 side 隔离冲突判定。
- `base/*` 只承载稳定公共结构，例如 `Paging`、`Sort`、`TimeRange`、`TransField`。
- 聚合根主对象使用顶层 message；仅服务于单个 RPC 的局部结构使用嵌套 message。
- 聚合根主 request 或 message 中的主对象 ID 优先使用 `id`，不要重复追加聚合根前缀。

## Prohibit

- 禁止用 Go struct 长期替代对外协议。
- 禁止手改 `*.pb.go`。
- 禁止跨 side import 其他 public proto。
- 禁止把 `base/*` 当业务 side 扩展目录使用。
- 禁止把只在单个 RPC 使用的局部结构平铺成全局顶层 message。
