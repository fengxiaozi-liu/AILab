# Util 参考

## 这个主题解决什么问题

沉淀跨模块复用的无状态基础工具，例如并发执行、编码转换、数值处理、链式处理等。

## 适用场景

- 某个纯工具函数已经在多个模块重复出现
- 需要稳定的并发、链式、数值、载荷辅助函数

## 推荐结构或实现方式

- util 函数保持短小、纯净、可测试。
- 命名直接体现能力，如 `Parallel`、`Chain`、`DecodeBase64`，不要使用模糊入口。

## 标准模板

```go
results, err := util.Parallel(ctx, tasks...)
if err != nil {
    return err
}
```

## Good Example

- `parallel.go`、`base64.go`、`decimal.go` 这类稳定工具按主题拆分，并配套测试。

## 常见坑

- 为单一业务流程创建工具函数后直接放入 util
- 在 util 中访问 repo、config 或业务 context

## 相关 rule / 相关 reference

- `../rules/util-rule.md`
