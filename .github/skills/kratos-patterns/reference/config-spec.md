# 配置体系规范（Kratos 通用）

> 本文档描述 Kratos 项目中“配置契约 + 配置文件 + 注入”的通用模式。

## 核心原则

1. **Contract First**：配置结构优先用 proto 定义（可生成强类型结构体）
2. **配置只读**：业务运行时不要修改配置对象
3. **默认值可控**：缺省值写在 config.yaml 或构造函数中，避免分散

## Go Demo（结构示例）

### 1) conf.proto 定义（示意）

```proto
syntax = "proto3";
package demo.conf;

message Bootstrap {
  Server server = 1;
}

message Server {
  string http_addr = 1;
}
```

### 2) 读取配置并注入（示意）

```go
package server

import (
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
	"github.com/go-kratos/kratos/v2/log"
)

// Bootstrap 为示意：真实项目是由 proto 生成。
type Bootstrap struct {
	Server struct {
		HTTPAddr string `json:"http_addr"`
	} `json:"server"`
}

func LoadConfig(path string, logger log.Logger) (*Bootstrap, error) {
	c := config.New(
		config.WithSource(file.NewSource(path)),
	)
	if err := c.Load(); err != nil {
		return nil, err
	}
	defer c.Close()

	var bc Bootstrap
	if err := c.Scan(&bc); err != nil {
		return nil, err
	}
	_ = logger
	return &bc, nil
}
```

## 常见陷阱

- 直接在业务里解析 yaml/json，丢失契约与校验
- 把配置当“动态开关”随业务写入（会导致不可预期行为）
- 多处重复定义默认值（上线后难以追溯）
