# Config Reference

## 这个主题解决什么问题

说明配置项如何落位、如何设计默认值、以及如何在代码中读取和使用配置。

## 适用场景

- 新增配置项
- 调整环境差异
- 设计配置结构体

## 设计意图

Config 参考解释的是配置如何表达运行时差异，而不是把常量简单搬进配置文件。

- 配置通常承接环境差异、开关和可调参数，而不是业务规则本身。
- 先理解配置的业务背景，就更容易判断一个值该放配置、常量还是数据库。
- 配置结构清楚后，默认值、示例、校验和回滚路径会更容易保持一致。

## 实施提示

- 先判断这个值为什么需要配置化，以及由谁调整。
- 再决定它的层级、命名、默认值和示例位置。
- 如果配置影响行为分支，顺手补上启动校验或最小验证样例。

## 推荐结构

- 配置按模块分组
- 结构体字段名与配置语义保持一致
- 读取配置后在构造函数中使用

## 标准模板

```go
type Depend struct {
    Timeout string `json:"timeout" yaml:"timeout"`
}
```

```go
type AccountConfig struct {
    SyncCron string `json:"sync_cron" yaml:"sync_cron"`
    PageSize int32  `json:"page_size" yaml:"page_size"`
}

func (c *AccountConfig) Normalize() {
    if c.PageSize <= 0 {
        c.PageSize = 100
    }
}
```

## Good Example

- 配置结构按服务或模块组织
- 默认值和环境差异在加载层或初始化层体现

## 常见坑

- 配置字段名只服务当前实现细节
- 一个配置项同时被多个无关模块随意复用
- 配置结构层级过深，初始化时难以理解

## 相关 Rule

- `../rules/config-rule.md`
