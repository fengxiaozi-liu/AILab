# 项目结构规范

## 概述

本项目采用 [Kratos](https://go-kratos.dev/) 微服务框架，遵循 DDD (领域驱动设计) 思想，采用分层架构设计。

## 目录结构

```
├── api/                 # API proto 定义 (接口契约)
│   ├── admin/           # 管理后台相关 API
│   │   ├── admin/       # 管理端接口
│   │   ├── inner/       # 服务间调用接口
│   │   └── open/        # C端接口
│   ├── base/            # 基础公共 API
│   │   ├── business/    # 业务公共消息定义
│   │   └── example/     # 示例模块
│   ├── gateway/         # 网关相关 API
│   │   ├── admin/       # 管理端接口
│   │   ├── inner/       # 服务间调用接口
│   │   └── open/        # C端接口
│   ├── system/          # 系统级 API
│   │   ├── admin/       # 管理端接口
│   │   ├── inner/       # 服务间调用接口
│   │   └── open/        # C端接口
│   └── user/            # 用户相关 API
│       ├── admin/       # 管理端接口
│       ├── inner/       # 服务间调用接口
│       └── open/        # C端接口
├── assets/              # 资源文件
│   ├── i18n/            # 国际化文件
│   └── seata/           # Seata 配置
├── cmd/                 # 项目启动入口
│   └── server/
│       ├── main.go      # 主入口
│       ├── wire.go      # Wire 依赖注入定义
│       └── wire_gen.go  # Wire 自动生成代码
├── configs/             # 配置文件
│   └── config.yaml      # 本地配置
├── docs/                # 文档目录
│   ├── project/         # 项目文档
│   ├── scaffold/        # 脚手架文档
│   └── spec/            # 开发规范
├── internal/            # 内部模块（不对外暴露）
│   ├── api/             # API proto 自动生成的 Go 代码
│   ├── biz/             # 业务逻辑层（UseCase）
│   ├── conf/            # 配置结构定义
│   ├── consumer/        # AMQP 消费者
│   ├── crontab/         # 定时任务
│   ├── data/            # 数据访问层（Repository）
│   ├── enum/            # 枚举定义
│   ├── error/           # 错误定义
│   ├── listener/        # 内部事件总线监听器
│   ├── pkg/             # 内部公共代码
│   ├── server/          # HTTP/gRPC 服务器配置
│   └── service/         # API 实现层（Controller）
│       ├── admin/       # 管理端接口
│       ├── inner/       # 服务间调用接口
│       └── open/        # C端接口
├── manifests/           # K8s 部署文件
│   ├── dev/             # 开发环境
│   └── prod/            # 生产环境
└── third_party/         # 第三方 proto 定义
```

## 新增业务模块步骤

假设新增 `order` 订单模块：

### 1. 定义 API

- 创建 `api/order/admin/v1/order.proto`
- 创建 `api/order/inner/v1/order.proto`
- 创建 `api/order/open/v1/order.proto`

### 2. 生成代码

```shell
make api
```

### 3. 创建 Biz 层

- 创建 `internal/biz/order/order.go`文件，编写Order领域模型 + OrderFilter查询条件构造器 + OrderRepo + OrderUseCase

### 4. 创建 Data 层

- 创建 `internal/data/order/order.go`文件，编写orderRepo struct，实现 biz 下定义的 OrderRepo 接口

### 5. 创建 Service 层

- 创建 `internal/service/admin/order.go` 文件，实现管理端 API
- 创建 `internal/service/inner/order.go` 文件，实现服务间调用 API
- 创建 `internal/service/open/order.go` 文件，实现 C 端 API

### 6. 创建 Ent Schema

- 创建 `internal/data/ent/schema/order.go`
- 执行 Ent 代码生成

```shell
cd internal/data/ent
ent generate ./schema --feature privacy,entql,sql/lock,sql/modifier,intercept,schema/snapshot,sql/upsert --template ./template
```

### 7. 注册依赖

在各层的 ProviderSet 中注册：

- `internal/biz/biz.go`
- `internal/data/data.go`
- `internal/service/admin/admin.go`
- `internal/service/inner/inner.go`
- `internal/service/open/open.go`

### 8. 注册 gRPC 服务

在 `internal/server/grpc.go` 中注册服务

## 目录命名规范

| 类型   | 命名规则     | 示例              |
|------|----------|-----------------|
| 包名   | 小写单词     | `order`, `user` |
| 版本目录 | `v` + 数字 | `v1`, `v2`      |
| 配置目录 | 小写单词     | `dev`, `prod`   |

## 文件命名规范

| 类型       | 命名规则  | 示例                 |
|----------|-------|--------------------|
| Go 文件    | 小写下划线 | `order_item.go`    |
| Proto 文件 | 小写下划线 | `order_item.proto` |
| 配置文件     | 小写下划线 | `config.yaml`      |
