# Error Rule

## Principles

- 统一错误体系优先于原生错误直接外抛。

## Specification

- not found 按业务语义转换。
- 对外或跨层返回项目统一错误。
- 保持错误可观测、可定位。

## Prohibit

- 禁止把“必须存在”的查询返回 `nil, nil`。
- 禁止跨边界直接返回 `errors.New` 或 `fmt.Errorf` 作为最终错误。
