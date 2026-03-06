# Design Patterns Reference

## 这个主题解决什么问题

说明在 Go 代码里常见的轻量设计模式使用场景，例如选项模式、适配器、构造函数注入。

## 适用场景

- 需要可扩展初始化参数
- 需要协议适配
- 需要减少直接耦合

## 设计意图

设计模式参考的作用是帮助识别什么时候抽象真的能降低复杂度，而不是机械地往代码里套名词。

- 模式不是目标，稳定的职责分离和扩展点才是目标。
- 先理解变化点和复用点，再选择策略、工厂或适配器，代码会更贴近实际问题。
- 适度抽象能减少重复，过度抽象会让简单流程变得难读，这正是解释层需要补回来的部分。

## 实施提示

- 先判断当前痛点是扩展困难、依赖耦合还是实现重复。
- 再选最小的抽象方式，而不是一开始引入完整模式族。
- 如果模式没有带来明显复用或隔离价值，通常说明抽象时机还不成熟。

## 推荐模式

- Functional Options
- Adapter
- Constructor Injection

## 标准模板

```go
type Option func(*Config)

func WithTimeout(d time.Duration) Option {
    return func(c *Config) { c.Timeout = d }
}
```

## 常见坑

- 为简单场景引入过重抽象
- 同时叠加多种模式导致阅读成本过高

## 相关 Reference

- `./style-organization.md`
