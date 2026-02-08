# InnerRpc 调用封装规范

## 概述

本文件主要针对innerrpc调用进行封装规范说明，旨在统一调用方式，提高代码可读性和维护性。

## 适用范围
本规范适用于所有使用InnerRpc进行服务间调用的项目代码。

## 文件模板

### innerrpc

- 以.proto文件形式定义服务接口和消息格式。

```protobuf
syntax = "proto3";

package system.inner.v1;

import "base/business/v1/business.proto";
import "system/inner/v1/counter_channel.proto";

option go_package = "linksoft.cn/trader/internal/api/system/inner/v1;v1";

service CounterChannelAccountService {
  rpc GetCounterChannelAccount (GetCounterChannelAccountRequest) returns (CounterChannelAccount);
  rpc GetCounterChannelAccountByAccount (GetCounterChannelAccountByAccountRequest) returns (CounterChannelAccount);
  rpc ListCounterChannelAccount (ListCounterChannelAccountRequest) returns (ListCounterChannelAccountReply);
  rpc MapCounterChannelAccount (MapCounterChannelAccountRequest) returns (MapCounterChannelAccountReply);
}

message CounterChannelAccount {
  uint32 id = 1;
  string account = 2;
  string name = 3;
  uint32 create_time = 5;
  uint32 update_time = 6;
  uint32 status = 7;
  uint32 counter_channel_id = 11;
  string scene = 12;
  CounterChannel counter_channel_info = 13;
}

message GetCounterChannelAccountRequest{
  uint32 id = 1;
  base.business.v1.FilterConfig filter_config = 2;
}

message GetCounterChannelAccountByAccountRequest{
  string account = 1;
  base.business.v1.FilterConfig filter_config = 2;
}

message MapCounterChannelAccountRequest {
  repeated uint32 id_list = 1;
  string account = 2;
  repeated string account_list = 3;
  uint32 counter_channel_id = 4;
  uint32 status = 5;
  string name = 6;
  base.business.v1.TimeRange create_time = 7;
  base.business.v1.Paging paging = 8;
  base.business.v1.Sort sort = 9;
  base.business.v1.FilterConfig filter_config = 10;
}
message MapCounterChannelAccountReply {
  map<uint32, CounterChannelAccount> map = 1;
}

message ListCounterChannelAccountRequest {
  repeated uint32 id_list = 1;
  string account = 2;
  repeated string account_list = 3;
  uint32 counter_channel_id = 4;
  uint32 status = 5;
  string name = 6;
  base.business.v1.TimeRange create_time = 7;
  base.business.v1.Paging paging = 8;
  base.business.v1.Sort sort = 9;
  base.business.v1.FilterConfig filter_config = 10;

}
message ListCounterChannelAccountReply {
  repeated CounterChannelAccount list = 1;
}

```

### internal/biz/depend

- 定义领域模型和Repo接口。
- 指导其他服务如何与该领域交互，可以使用哪些接口。

```go
package system

import (
	"context"

	systemenum "linksoft.cn/trader/internal/enum/system"

	"linksoft.cn/trader/internal/enum/base"
	"linksoft.cn/trader/internal/pkg/filter"
)

type CounterChannelAccount struct {
	ID uint32 `json:"id"`
	// 创建时间
	CreateTime uint32 `json:"create_time"`
	// 更新时间
	UpdateTime uint32 `json:"update_time"`
	// 通道名称
	CounterChannelID uint32 `json:"counter_channel_id"`
	// 状态 1 启用 2 禁用
	Status base.Switch `json:"status"`
	// 账户号
	Account string `json:"account"`
	// 账户名称
	Name string `json:"name"`
	// 通道场景
	Scene systemenum.CounterChannelScene `json:"scene"`
	// 通道信息
	CounterChannelInfo *CounterChannel `json:"counter_channel_info"`
}

type CounterChannelAccountFilter struct {
	IDList           []uint32
	IDNotInList      []uint32
	Account          string
	Name             string
	CounterChannelID uint32
	Status           base.Switch
	CreateTime       filter.TimeRange
	Paging           filter.Paging
	Sort             filter.Sort
	AccountList      []string
}

type CounterChannelAccountRepo interface {
	ListCounterChannelAccount(ctx context.Context, filter *CounterChannelAccountFilter, opts ...filter.Option) ([]*CounterChannelAccount, error)
	GetCounterChannelAccount(ctx context.Context, ID uint32, opts ...filter.Option) (*CounterChannelAccount, error)
	GetCounterChannelAccountByAccount(ctx context.Context, account string, opts ...filter.Option) (*CounterChannelAccount, error)
	MapCounterChannelAccount(ctx context.Context, filter *CounterChannelAccountFilter, opts ...filter.Option) (map[uint32]*CounterChannelAccount, error)
}

```

### internal/data/depend

- 实现 Repo 接口，封装 innerrpc 调用细节。
- 转换领域模型与 protobuf 消息格式。

```go
package system

import (
	"context"
	"linksoft.cn/trader/internal/biz/depend/system"
	systemenum "linksoft.cn/trader/internal/enum/system"

	"gitlab.linksoft.cn/TeamA/GoPackage/helper"

	v1 "linksoft.cn/trader/internal/api/system/inner/v1"
	baseenum "linksoft.cn/trader/internal/enum/base"
	"linksoft.cn/trader/internal/pkg/filter"
	"linksoft.cn/trader/internal/pkg/proto"
)

type counterChannelAccountRepo struct {
	counterChannelAccountClient v1.CounterChannelAccountServiceClient
}

func NewCounterChannelAccountRepo(counterChannelAccountClient v1.CounterChannelAccountServiceClient) system.CounterChannelAccountRepo {
	return &counterChannelAccountRepo{counterChannelAccountClient: counterChannelAccountClient}
}

func (c *counterChannelAccountRepo) ListCounterChannelAccount(ctx context.Context, filter *system.CounterChannelAccountFilter, opts ...filter.Option) ([]*system.CounterChannelAccount, error) {
	reply, err := c.counterChannelAccountClient.ListCounterChannelAccount(ctx, &v1.ListCounterChannelAccountRequest{
		Account:          filter.Account,
		CounterChannelId: filter.CounterChannelID,
		Status:           uint32(filter.Status),
		Name:             filter.Name,
		CreateTime:       proto.BuildTimeRange(filter.CreateTime),
		Paging:           proto.BuildPaging(filter.Paging),
		Sort:             proto.BuildSort(filter.Sort),
		IdList:           filter.IDList,
		FilterConfig:     proto.BuildFilterConfig(opts...),
		AccountList:      filter.AccountList,
	})
	if err != nil {
		return nil, err
	}

	return helper.SliceConvert(reply.List, CounterChannelAccountConvert), err
}

func (c *counterChannelAccountRepo) GetCounterChannelAccount(ctx context.Context, ID uint32, opts ...filter.Option) (*system.CounterChannelAccount, error) {
	info, err := c.counterChannelAccountClient.GetCounterChannelAccount(ctx, &v1.GetCounterChannelAccountRequest{Id: ID, FilterConfig: proto.BuildFilterConfig(opts...)})
	if err != nil {
		return nil, err
	}
	return CounterChannelAccountConvert(info), nil
}

func (c *counterChannelAccountRepo) GetCounterChannelAccountByAccount(ctx context.Context, account string, opts ...filter.Option) (*system.CounterChannelAccount, error) {
	info, err := c.counterChannelAccountClient.GetCounterChannelAccountByAccount(ctx, &v1.GetCounterChannelAccountByAccountRequest{Account: account, FilterConfig: proto.BuildFilterConfig(opts...)})
	if err != nil {
		return nil, err
	}
	return CounterChannelAccountConvert(info), nil
}

func (c *counterChannelAccountRepo) MapCounterChannelAccount(ctx context.Context, filter *system.CounterChannelAccountFilter, opts ...filter.Option) (map[uint32]*system.CounterChannelAccount, error) {
	list, err := c.counterChannelAccountClient.MapCounterChannelAccount(ctx, &v1.MapCounterChannelAccountRequest{
		Account:          filter.Account,
		CounterChannelId: filter.CounterChannelID,
		Status:           uint32(filter.Status),
		Name:             filter.Name,
		CreateTime:       proto.BuildTimeRange(filter.CreateTime),
		Paging:           proto.BuildPaging(filter.Paging),
		Sort:             proto.BuildSort(filter.Sort),
		IdList:           filter.IDList,
		FilterConfig:     proto.BuildFilterConfig(opts...),
	})
	if err != nil {
		return nil, err
	}
	res := make(map[uint32]*system.CounterChannelAccount, len(list.Map))
	for k, i := range list.Map {
		res[k] = CounterChannelAccountConvert(i)
	}

	return res, nil
}

func CounterChannelAccountConvert(channel *v1.CounterChannelAccount) *system.CounterChannelAccount {
	if channel == nil {
		return nil
	}
	return &system.CounterChannelAccount{
		ID:                 channel.Id,
		CreateTime:         channel.CreateTime,
		UpdateTime:         channel.UpdateTime,
		CounterChannelID:   channel.CounterChannelId,
		Scene:              systemenum.CounterChannelScene(channel.Scene),
		Status:             baseenum.Switch(channel.Status),
		Account:            channel.Account,
		Name:               channel.Name,
		CounterChannelInfo: CounterChannelConvert(channel.GetCounterChannelInfo()),
	}
}

func CounterChannelAccountRestore(info *system.CounterChannelAccount) *v1.CounterChannelAccount {
	if info == nil {
		return nil
	}
	return &v1.CounterChannelAccount{
		Id:                 info.ID,
		CreateTime:         info.CreateTime,
		UpdateTime:         info.UpdateTime,
		CounterChannelId:   info.CounterChannelID,
		Scene:              string(info.Scene),
		Status:             uint32(info.Status),
		Account:            info.Account,
		Name:               info.Name,
		CounterChannelInfo: CounterChannelRestore(info.CounterChannelInfo),
	}
}

```

### 服务注册

- 在/internal/data/depend/{所属服务}/{所属服务}.go中注册依赖实现，以及底层实际调用的`innerrpc ServiceClient`。

```go
package system

import (
	"github.com/google/wire"
	v1 "linksoft.cn/trader/internal/api/system/inner/v1"
	"linksoft.cn/trader/internal/data/kit"
)

var ProviderSet = wire.NewSet(
	NewCounterChannelAccountRepo,
	NewV1CounterChannelAccountServiceClient,
)

func NewV1CounterChannelAccountServiceClient(client *kit.SystemServiceClient) v1.CounterChannelAccountServiceClient {
	return v1.NewCounterChannelAccountServiceClient(client)
}
```
