# 网关代理规范

## 介绍
网关代理规范定义了如何配置和使用API网关来管理和路由请求到后端服务。它涵盖了代理设置、路由规则、安全性和性能优化等方面。

## 代理配置

### 代理文件编写
- 对应api文件，例如`api/order/admin/v1/order_trade_charge.proto`，代理文件放在`internal/proxy/admin/v1/order`目录下，命名为`order_trade_charge.go`
```go
package order

import (
	"context"

	v1 "linksoft.cn/trader/internal/api/order/admin/v1"
)

type OrderTradeChargeProxy struct {
	v1.UnimplementedOrderTradeChargeServiceServer
	orderTradeChargeServiceClient v1.OrderTradeChargeServiceClient
}

func NewOrderTradeChargeProxy(orderTradeChargeServiceClient v1.OrderTradeChargeServiceClient) *OrderTradeChargeProxy {
	return &OrderTradeChargeProxy{orderTradeChargeServiceClient: orderTradeChargeServiceClient}
}

func (srv *OrderTradeChargeProxy) GetOrderTradeCharge(ctx context.Context, req *v1.GetOrderTradeChargeRequest) (*v1.OrderTradeCharge, error) {
	return srv.orderTradeChargeServiceClient.GetOrderTradeCharge(ctx, req)
}

func (srv *OrderTradeChargeProxy) ListOrderTradeCharge(ctx context.Context, req *v1.ListOrderTradeChargeRequest) (*v1.ListOrderTradeChargeReply, error) {
	return srv.orderTradeChargeServiceClient.ListOrderTradeCharge(ctx, req)
}

func (srv *OrderTradeChargeProxy) PageListOrderTradeCharge(ctx context.Context, req *v1.PageListOrderTradeChargeRequest) (*v1.PageListOrderTradeChargeReply, error) {
	return srv.orderTradeChargeServiceClient.PageListOrderTradeCharge(ctx, req)
}

func (srv *OrderTradeChargeProxy) CountOrderTradeCharge(ctx context.Context, req *v1.CountOrderTradeChargeRequest) (*v1.CountOrderTradeChargeReply, error) {
	return srv.orderTradeChargeServiceClient.CountOrderTradeCharge(ctx, req)
}

```

### 注册代理服务
- 分别注册open，admin等业务域的代理服务。例如在`internal/proxy/admin/v1/admin.go`中注册代理服务
```go
package admin

import (
	"github.com/google/wire"
	v1order "linksoft.cn/trader/internal/proxy/admin/v1/order"
)

var ProviderSet = wire.NewSet(
	v1order.NewOrderTradeChargeProxy,
)

```

### 注册代理路由
- 例如在`internal/server/http.go`中注册admin域的代理路由。

```go
package server

import (
	"strings"

	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/go-kratos/kratos/v2/transport/http/pprof"
	orderadminv1api "linksoft.cn/trader/internal/api/order/admin/v1"
	"linksoft.cn/trader/internal/conf"
	orderadminv1proxy "linksoft.cn/trader/internal/proxy/admin/v1/order"
)

// NewHTTPServer new an HTTP server.
func NewHTTPServer(
	orderAdminV1OrderTradeChargeProxy *orderadminv1proxy.OrderTradeChargeProxy,
	c *conf.Server,
) *http.Server {
	var opts = []http.ServerOption{}
	if c.Http.Network != "" {
		opts = append(opts, http.Network(c.Http.Network))
	}
	if c.Http.Addr != "" {
		opts = append(opts, http.Address(c.Http.Addr))
	}
	if c.Http.Timeout != nil {
		opts = append(opts, http.Timeout(c.Http.Timeout.AsDuration()))
	}

	srv := http.NewServer(opts...)

	if strings.EqualFold(c.Pprof.Enable, "true") {
		srv.Handle("/debug/pprof", pprof.NewHandler())
	}

	orderadminv1api.RegisterOrderTradeChargeServiceHTTPServer(srv, orderAdminV1OrderTradeChargeProxy)

	return srv
}

```
