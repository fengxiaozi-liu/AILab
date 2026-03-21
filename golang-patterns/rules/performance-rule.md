# Performance Rule

## Principles

- 先正确清晰，再做低风险优化。

## Specification

- 可预估容量时为 slice/map 预分配。
- 热路径避免明显重复分配。

## Prohibit

- 禁止在热循环里无脑扩容或重复构造临时对象。
