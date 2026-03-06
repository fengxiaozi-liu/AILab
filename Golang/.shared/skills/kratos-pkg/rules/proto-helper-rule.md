# Proto Helper Rule

## Principles

- `internal/pkg/proto` 只负责稳定的协议辅助转换，不负责业务协议设计。
- proto helper 的输入输出应围绕公共 `base/business` 结构和内部 `filter` 等基础类型。

## Specification

- `TransField`、`Paging`、`Sort`、`TimeRange`、`GroupBy`、`Compare`、`FilterConfig` 这类通用协议转换可以放在 proto helper。
- 命名使用 `ParseXxx` 和 `BuildXxx` 成对表达方向。
- 协议字段与内部 filter 字段的命名转换要集中处理，不在各个 service/repo 重复实现。

## Prohibit

- 禁止在 proto helper 中拼装业务 reply 或业务聚合。
- 禁止把某个单一业务 proto 的专用组装逻辑下沉到 `internal/pkg/proto`。
- 禁止在 service、repo 中复制 `Paging/Sort/TimeRange/TransField` 转换代码。

## Checklist

- 这段转换是否对多个模块稳定复用。
- 是否使用了 `ParseXxx/BuildXxx` 成对命名。
- 是否只处理公共协议与基础过滤能力。
