# Schema 参考

## 这个主题解决什么问题

提供通用结构字段提取和结构分析能力，便于基础设施层做通用处理。

## 适用场景

- 需要从结构体提取字段定义
- 需要做通用字段映射、列提取、结构检查

## 推荐结构或实现方式

- schema helper 的输入输出保持通用，不绑定具体业务聚合。
- 优先围绕字段提取、结构分析、通用反射辅助组织代码。

## 标准模板

```go
fields, err := Extract(SomeStruct{})
if err != nil {
    return err
}
```

## Good Example

- 提供稳定的 `extract.go` 与对应测试，调用方按需消费提取结果。

## 常见坑

- 为一次性业务场景创建专用 schema helper
- 在 schema helper 中混入业务校验和业务配置判断

## 相关 rule / 相关 reference

- `../rules/schema-rule.md`
