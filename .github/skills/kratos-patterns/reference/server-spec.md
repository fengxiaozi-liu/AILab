# rpc服务注册规范

## 概述
本项目使用 gRPC 框架实现微服务间通信。为了确保服务的可发现性和可扩展性，所有 gRPC 服务必须遵循统一的注册规范。
## 服务注册步骤

### grpc

1. 确认服务实现文件位于 `internal/service/[service_type]/[version]/[service].go`，其中 `[service_type]` 可以是 `admin`、`inner` 或 `open`。
2. 在 `internal/server/grpc.go` 文件中，导入对应的服务实现包。

```go
package server

import (
	examplev1 "linksoft.cn/trader/internal/api/base/example/v1"
	"linksoft.cn/trader/internal/conf"
	"linksoft.cn/trader/internal/service/example"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport/grpc"

	adminapiv1 "linksoft.cn/trader/internal/api/order/admin/v1"
	innerapiv1 "linksoft.cn/trader/internal/api/order/inner/v1"
	traderapiv1 "linksoft.cn/trader/internal/api/order/trader/v1"
	adminservicev1 "linksoft.cn/trader/internal/service/admin/v1"
	innerservicev1 "linksoft.cn/trader/internal/service/inner/v1"
	traderservicev1 "linksoft.cn/trader/internal/service/trader/v1"
)

// NewGRPCServer new a gRPC server.
func NewGRPCServer(
	exampleService *example.ExampleService,
	adminOrderService *adminservicev1.OrderService,
	adminPositionService *adminservicev1.PositionService,
	adminTradeChargePackageService *adminservicev1.TradeChargePackageService,
	adminTradeChargeService *adminservicev1.TradeChargeService,
	adminOrderTradeChargeService *adminservicev1.OrderTradeChargeService,
	adminMarginPoolService *adminservicev1.MarginPoolService,
	adminMarginApplyService *adminservicev1.MarginApplyService,
	adminMarginQuotaService *adminservicev1.MarginQuotaService,
	adminMarginRepayRecordService *adminservicev1.MarginRepayRecordService,
	adminPositionHedgeService *adminservicev1.PositionHedgeService,
	innerOrderService *innerservicev1.OrderService,
	innerPositionService *innerservicev1.PositionService,
	innerTradeChargePackageService *innerservicev1.TradeChargePackageService,
	innerMarginQuotaService *innerservicev1.MarginQuotaService,
	traderOrderService *traderservicev1.OrderService,
	traderTradeChargeService *traderservicev1.OrderTradeChargeService,
	traderPositionService *traderservicev1.PositionService,
	traderMarginPoolService *traderservicev1.MarginPoolService,
	traderMarginApplyService *traderservicev1.MarginApplyService,
	traderMarginQuotaService *traderservicev1.MarginQuotaService,
	traderMarginRepayRecordService *traderservicev1.MarginRepayRecordService,
	c *conf.Server,
	_ log.Logger,
) *grpc.Server {
	var opts = []grpc.ServerOption{}
	if c.Grpc.Network != "" {
		opts = append(opts, grpc.Network(c.Grpc.Network))
	}
	if c.Grpc.Addr != "" {
		opts = append(opts, grpc.Address(c.Grpc.Addr))
	}
	if c.Grpc.Timeout != nil {
		opts = append(opts, grpc.Timeout(c.Grpc.Timeout.AsDuration()))
	}

	srv := grpc.NewServer(opts...)

	examplev1.RegisterExampleServiceServer(srv, exampleService)

	adminapiv1.RegisterOrderServiceServer(srv, adminOrderService)
	adminapiv1.RegisterPositionServiceServer(srv, adminPositionService)
	adminapiv1.RegisterTradeChargePackageServiceServer(srv, adminTradeChargePackageService)
	adminapiv1.RegisterTradeChargeServiceServer(srv, adminTradeChargeService)
	adminapiv1.RegisterOrderTradeChargeServiceServer(srv, adminOrderTradeChargeService)

	adminapiv1.RegisterMarginPoolServiceServer(srv, adminMarginPoolService)
	adminapiv1.RegisterMarginApplyServiceServer(srv, adminMarginApplyService)
	adminapiv1.RegisterMarginQuotaServiceServer(srv, adminMarginQuotaService)
	adminapiv1.RegisterMarginRepayRecordServiceServer(srv, adminMarginRepayRecordService)

	adminapiv1.RegisterPositionHedgeServiceServer(srv, adminPositionHedgeService)

	innerapiv1.RegisterOrderServiceServer(srv, innerOrderService)
	innerapiv1.RegisterPositionServiceServer(srv, innerPositionService)
	innerapiv1.RegisterTradeChargePackageServiceServer(srv, innerTradeChargePackageService)
	innerapiv1.RegisterMarginQuotaServiceServer(srv, innerMarginQuotaService)

	traderapiv1.RegisterOrderServiceServer(srv, traderOrderService)
	traderapiv1.RegisterOrderTradeChargeServiceServer(srv, traderTradeChargeService)
	traderapiv1.RegisterPositionServiceServer(srv, traderPositionService)

	traderapiv1.RegisterMarginPoolServiceServer(srv, traderMarginPoolService)
	traderapiv1.RegisterMarginApplyServiceServer(srv, traderMarginApplyService)
	traderapiv1.RegisterMarginQuotaServiceServer(srv, traderMarginQuotaService)
	traderapiv1.RegisterMarginRepayRecordServiceServer(srv, traderMarginRepayRecordService)
	return srv
}

```

### http/tcp/websocket

1. 确认服务实现文件位于 `internal/service/[service_type]/[version]/[service].go`，其中 `[service_type]` 可以是 `admin`、`inner` 或 `open`。
2. 只有gateway类型项目需要注册HTTP服务。
3. 在 `internal/server/{http|tcp|websocket}.go` 文件中，导入对应的服务实现包。

```go
package server

import (
	"strings"

	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/go-kratos/kratos/v2/transport/http/pprof"
	adminadminv1api "linksoft.cn/trader/internal/api/admin/admin/v1"
	assetadminv1api "linksoft.cn/trader/internal/api/asset/admin/v1"
	orderadminv1api "linksoft.cn/trader/internal/api/order/admin/v1"
	riskcontroladminv1api "linksoft.cn/trader/internal/api/riskcontrol/admin/v1"
	settlementadminv1api "linksoft.cn/trader/internal/api/settlement/admin/v1"
	systemadminv1api "linksoft.cn/trader/internal/api/system/admin/v1"
	systemtraderv1api "linksoft.cn/trader/internal/api/system/trader/v1"
	traderadminv1api "linksoft.cn/trader/internal/api/trader/admin/v1"
	"linksoft.cn/trader/internal/conf"
	httppkg "linksoft.cn/trader/internal/pkg/transport/http"
	adminadminv1proxy "linksoft.cn/trader/internal/proxy/admin/v1/admin"
	assetadminv1proxy "linksoft.cn/trader/internal/proxy/admin/v1/asset"
	orderadminv1proxy "linksoft.cn/trader/internal/proxy/admin/v1/order"
	riskcontroladminv1proxy "linksoft.cn/trader/internal/proxy/admin/v1/riskcontrol"
	settlementadminv1proxy "linksoft.cn/trader/internal/proxy/admin/v1/settlement"
	systemadminv1proxy "linksoft.cn/trader/internal/proxy/admin/v1/system"
	traderadminv1proxy "linksoft.cn/trader/internal/proxy/admin/v1/trader"
	systembusinessv1proxy "linksoft.cn/trader/internal/proxy/business/v1/system"
	systemtraderv1proxy "linksoft.cn/trader/internal/proxy/trader/v1/system"
)

// NewHTTPServer new an HTTP server.
func NewHTTPServer(
	adminAdminV1PermissionProxy *adminadminv1proxy.PermissionProxy,
	adminAdminV1AdminRoleProxy *adminadminv1proxy.AdminRoleProxy,
	adminAdminV1AdminUserProxy *adminadminv1proxy.AdminUserProxy,
	systemAdminV1ChargePersonProxy *systemadminv1proxy.ChargePersonProxy,
	orderAdminV1OrderProxy *orderadminv1proxy.OrderProxy,
	orderAdminV1PositionProxy *orderadminv1proxy.PositionProxy,
	orderAdminV1TradeChargeProxy *orderadminv1proxy.TradeChargeProxy,
	orderAdminV1TradeChargePackageProxy *orderadminv1proxy.TradeChargePackageProxy,
	orderAdminV1OrderTradeChargeProxy *orderadminv1proxy.OrderTradeChargeProxy,
	orderAdminV1MarginPoolProxy *orderadminv1proxy.MarginPoolProxy,
	orderAdminV1MarginApplyProxy *orderadminv1proxy.MarginApplyProxy,
	orderAdminV1MarginQuotaProxy *orderadminv1proxy.MarginQuotaProxy,
	orderAdminV1MarginRepayRecordProxy *orderadminv1proxy.MarginRepayRecordProxy,
	orderAdminV1PositionHedgeProxy *orderadminv1proxy.PositionHedgeProxy,
	traderAdminV1TraderProxy *traderadminv1proxy.TraderProxy,
	traderAdminV1TraderGroupProxy *traderadminv1proxy.TraderGroupProxy,
	traderAdminV1TraderRiskControlProxy *traderadminv1proxy.TraderRiskControlProxy,
	systemAdminV1CounterChannelProxy *systemadminv1proxy.CounterChannelProxy,
	systemAdminV1CounterChannelAccountProxy *systemadminv1proxy.CounterChannelAccountProxy,
	systemAdminV1FileProxy *systemadminv1proxy.FileProxy,
	systemAdminV1CountryCodeProxy *systemadminv1proxy.CountryCodeProxy,
	systemAdminV1SystemConfigProxy *systemadminv1proxy.SystemConfigProxy,
	systemAdminV1ExchangeRateProxy *systemadminv1proxy.ExchangeRateProxy,
	systemAdminV1ExchangeRateRecordProxy *systemadminv1proxy.ExchangeRateRecordProxy,
	systemAdminV1EnumProxy *systemadminv1proxy.EnumProxy,
	systemAdminV1SecSecurityBasicProxy *systemadminv1proxy.SecurityBasicProxy,
	systemAdminV1ClientVersionProxy *systemadminv1proxy.ClientVersionProxy,
	settlementAdminV1StatisticProxy *settlementadminv1proxy.StatisticProxy,
	settlementAdminV1StatisticSecurityProxy *settlementadminv1proxy.StatisticSecurityProxy,
	settlementAdminV1SummaryProxy *settlementadminv1proxy.SummaryProxy,
	assetAdminV1TradeAccountProxy *assetadminv1proxy.TradeAccountProxy,
	assetAdminV1TradeAccountRecordProxy *assetadminv1proxy.TradeAccountRecordProxy,
	assetAdminV1TradeAccountRiskControlProxy *assetadminv1proxy.TradeAccountRiskControlProxy,
	riskControlAdminV1TradeAccountRiskControlRecordProxy *riskcontroladminv1proxy.TradeAccountRiskControlRecordProxy,
	riskControlAdminV1TraderRiskControlRecordProxy *riskcontroladminv1proxy.TraderRiskControlRecordProxy,
	systemTraderV1EnumProxy *systemtraderv1proxy.EnumProxy,
	systemTraderV1SystemConfigProxy *systemtraderv1proxy.SystemConfigProxy,
	systemTraderV1FileProxy *systemtraderv1proxy.FileProxy,
	systemTraderV1ClientVersionProxy *systemtraderv1proxy.ClientVersionProxy,
	systemBusinessV1CounterChannelOrderProxy *systembusinessv1proxy.CounterChannelOrderProxy,
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

	route := srv.Route("/")

	route.POST("/business.v1/system/counter/channel/order/notify", httppkg.WrapHandlerFunc(systembusinessv1proxy.OperationCounterChannelOrderServiceNotifyCounterChannelOrder, systemBusinessV1CounterChannelOrderProxy.NotifyCounterChannelOrder))

	route.POST("/admin.v1/system/enum/map", httppkg.WrapHandlerFunc(systemadminv1proxy.OperationEnumServiceMapEnum, systemAdminV1EnumProxy.MapEnum))
	route.POST("/admin.v1/system/config/map", httppkg.WrapHandlerFunc(systemadminv1proxy.OperationSystemConfigServiceMapSystemConfig, systemAdminV1SystemConfigProxy.MapSystemConfig))
	route.POST("/admin.v1/system/remote/config/map", httppkg.WrapHandlerFunc(systemadminv1proxy.OperationSystemConfigServiceMapSystemConfig, systemAdminV1SystemConfigProxy.MapRemoteSystemConfig))

	route.GET("/trader.v1/system/client/download/latest.yml", httppkg.WrapHandlerFunc(systemtraderv1proxy.OperationClientVersionServiceDownloadLatest, systemTraderV1ClientVersionProxy.DownloadLatest))
	route.GET("/trader.v1/system/client/download/latest-mac.yml", httppkg.WrapHandlerFunc(systemtraderv1proxy.OperationClientVersionServiceDownloadLatestMac, systemTraderV1ClientVersionProxy.DownloadLatestMac))
	route.POST("/trader.v1/system/enum/map", httppkg.WrapHandlerFunc(systemtraderv1proxy.OperationEnumServiceMapEnum, systemTraderV1EnumProxy.MapEnum))
	route.POST("/trader.v1/system/config/map", httppkg.WrapHandlerFunc(systemtraderv1proxy.OperationSystemConfigServiceMapSystemConfig, systemTraderV1SystemConfigProxy.MapSystemConfig))

	adminadminv1api.RegisterPermissionServiceHTTPServer(srv, adminAdminV1PermissionProxy)
	adminadminv1api.RegisterAdminRoleServiceHTTPServer(srv, adminAdminV1AdminRoleProxy)
	adminadminv1api.RegisterAdminUserServiceHTTPServer(srv, adminAdminV1AdminUserProxy)

	systemadminv1api.RegisterChargePersonServiceHTTPServer(srv, systemAdminV1ChargePersonProxy)

	orderadminv1api.RegisterOrderServiceHTTPServer(srv, orderAdminV1OrderProxy)
	orderadminv1api.RegisterPositionServiceHTTPServer(srv, orderAdminV1PositionProxy)
	orderadminv1api.RegisterTradeChargeServiceHTTPServer(srv, orderAdminV1TradeChargeProxy)
	orderadminv1api.RegisterTradeChargePackageServiceHTTPServer(srv, orderAdminV1TradeChargePackageProxy)
	orderadminv1api.RegisterOrderTradeChargeServiceHTTPServer(srv, orderAdminV1OrderTradeChargeProxy)

	orderadminv1api.RegisterMarginPoolServiceHTTPServer(srv, orderAdminV1MarginPoolProxy)
	orderadminv1api.RegisterMarginApplyServiceHTTPServer(srv, orderAdminV1MarginApplyProxy)
	orderadminv1api.RegisterMarginQuotaServiceHTTPServer(srv, orderAdminV1MarginQuotaProxy)
	orderadminv1api.RegisterMarginRepayRecordServiceHTTPServer(srv, orderAdminV1MarginRepayRecordProxy)

	orderadminv1api.RegisterPositionHedgeServiceHTTPServer(srv, orderAdminV1PositionHedgeProxy)

	traderadminv1api.RegisterTraderServiceHTTPServer(srv, traderAdminV1TraderProxy)
	traderadminv1api.RegisterTraderGroupServiceHTTPServer(srv, traderAdminV1TraderGroupProxy)
	traderadminv1api.RegisterTraderRiskControlServiceHTTPServer(srv, traderAdminV1TraderRiskControlProxy)

	systemadminv1api.RegisterCounterChannelServiceHTTPServer(srv, systemAdminV1CounterChannelProxy)
	systemadminv1api.RegisterCounterChannelAccountServiceHTTPServer(srv, systemAdminV1CounterChannelAccountProxy)
	systemadminv1api.RegisterChargePersonServiceHTTPServer(srv, systemAdminV1ChargePersonProxy)
	systemadminv1api.RegisterFileServiceHTTPServer(srv, systemAdminV1FileProxy)
	systemadminv1api.RegisterCountryCodeServiceHTTPServer(srv, systemAdminV1CountryCodeProxy)
	systemadminv1api.RegisterSystemConfigServiceHTTPServer(srv, systemAdminV1SystemConfigProxy)
	systemadminv1api.RegisterExchangeRateServiceHTTPServer(srv, systemAdminV1ExchangeRateProxy)
	systemadminv1api.RegisterExchangeRateRecordServiceHTTPServer(srv, systemAdminV1ExchangeRateRecordProxy)
	systemadminv1api.RegisterSecurityBasicServiceHTTPServer(srv, systemAdminV1SecSecurityBasicProxy)
	systemadminv1api.RegisterClientVersionServiceHTTPServer(srv, systemAdminV1ClientVersionProxy)

	settlementadminv1api.RegisterStatisticServiceHTTPServer(srv, settlementAdminV1StatisticProxy)
	settlementadminv1api.RegisterStatisticSecurityServiceHTTPServer(srv, settlementAdminV1StatisticSecurityProxy)
	settlementadminv1api.RegisterSummaryServiceHTTPServer(srv, settlementAdminV1SummaryProxy)

	assetadminv1api.RegisterTradeAccountServiceHTTPServer(srv, assetAdminV1TradeAccountProxy)
	assetadminv1api.RegisterTradeAccountRecordServiceHTTPServer(srv, assetAdminV1TradeAccountRecordProxy)
	assetadminv1api.RegisterTradeAccountRiskControlServiceHTTPServer(srv, assetAdminV1TradeAccountRiskControlProxy)

	riskcontroladminv1api.RegisterTradeAccountRiskControlRecordServiceHTTPServer(srv, riskControlAdminV1TradeAccountRiskControlRecordProxy)
	riskcontroladminv1api.RegisterTraderRiskControlRecordServiceHTTPServer(srv, riskControlAdminV1TraderRiskControlRecordProxy)

	systemtraderv1api.RegisterFileServiceHTTPServer(srv, systemTraderV1FileProxy)
	systemtraderv1api.RegisterClientVersionServiceHTTPServer(srv, systemTraderV1ClientVersionProxy)

	return srv
}

```
