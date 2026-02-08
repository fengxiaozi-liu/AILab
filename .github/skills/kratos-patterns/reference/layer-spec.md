# 代码分层规范

## 概述

项目采用经典三层架构：**Service** -> **Biz** -> **Data**

```
┌─────────────────┐
│     Service     │  API 接口实现层
├─────────────────┤
│       Biz       │  业务逻辑层
├─────────────────┤
│      Data       │  数据访问层
└─────────────────┘
```

## 层级职责

### Service 层 (internal/service)

**职责**：API 接口实现，请求/响应转换

#### admin域

```go
package v1

import (
	"context"

	businessv1 "{module from go.mod}/internal/api/base/business/v1"
	v1 "{module from go.mod}/internal/api/system/admin/v1"
	systembiz "{module from go.mod}/internal/biz/system"
	baseenum "{module from go.mod}/internal/enum/base"
	systemenum "{module from go.mod}/internal/enum/system"
	"{module from go.mod}/internal/pkg/filter"
	"{module from go.mod}/internal/pkg/proto"
)

type CounterChannelAccountService struct {
	v1.UnimplementedCounterChannelAccountServiceServer
	counterChannelAccountUseCase *systembiz.CounterChannelAccountUseCase
}

func NewCounterChannelAccountService(counterChannelAccountUseCase *systembiz.CounterChannelAccountUseCase) *CounterChannelAccountService {
	return &CounterChannelAccountService{counterChannelAccountUseCase: counterChannelAccountUseCase}
}

func (c *CounterChannelAccountService) GetCounterChannelAccount(ctx context.Context, req *v1.GetCounterChannelAccountRequest) (*v1.CounterChannelAccount, error) {
	counterChannel, err := c.counterChannelAccountUseCase.Get(ctx, req.GetId(), filter.WithRelation(systemenum.RelationCounterChannelInfo))
	if err != nil {
		return nil, err
	}
	return counterChannelAccountConvert(counterChannel), nil
}

func (c *CounterChannelAccountService) ListCounterChannelAccount(ctx context.Context, req *v1.ListCounterChannelAccountRequest) (*v1.ListCounterChannelAccountReply, error) {
	list, err := c.counterChannelAccountUseCase.List(ctx, &systembiz.CounterChannelAccountFilter{
		Account:          req.GetAccount(),
		Name:             req.GetName(),
		CounterChannelID: req.GetCounterChannelId(),
		Status:           baseenum.Switch(req.GetStatus()),
		CreateTime:       proto.ParseTimeRange(req.GetCreateTime()),
		Paging:           proto.ParsePaging(req.GetPaging()),
		Sort:             proto.ParseSort(req.GetSort()),
	}, filter.WithRelation(systemenum.RelationCounterChannelInfo))
	if err != nil {
		return nil, err
	}
	res := &v1.ListCounterChannelAccountReply{List: make([]*v1.CounterChannelAccount, 0)}
	for _, channel := range list {
		res.List = append(res.List, counterChannelAccountConvert(channel))
	}
	return res, nil
}

func (c *CounterChannelAccountService) PageListCounterChannelAccount(ctx context.Context, req *v1.PageListCounterChannelAccountRequest) (*v1.PageListCounterChannelAccountReply, error) {
	res := &v1.PageListCounterChannelAccountReply{List: make([]*v1.CounterChannelAccount, 0)}

	list, count, err := c.counterChannelAccountUseCase.PageList(ctx, &systembiz.CounterChannelAccountFilter{
		Account:          req.GetAccount(),
		Name:             req.GetName(),
		CounterChannelID: req.GetCounterChannelId(),
		Status:           baseenum.Switch(req.GetStatus()),
		CreateTime:       proto.ParseTimeRange(req.GetCreateTime()),
		Paging:           proto.ParsePaging(req.GetPaging()),
		Sort:             proto.ParseSort(req.GetSort()),
	}, filter.WithRelation(systemenum.RelationCounterChannelInfo))
	if err != nil {
		return nil, err
	}
	for _, channel := range list {
		res.List = append(res.List, counterChannelAccountConvert(channel))
	}
	res.Count = int32(count)

	return res, nil

}

func (c *CounterChannelAccountService) CountCounterChannelAccount(ctx context.Context, req *v1.CountCounterChannelAccountRequest) (*v1.CountCounterChannelAccountReply, error) {
	count, err := c.counterChannelAccountUseCase.Count(ctx, &systembiz.CounterChannelAccountFilter{
		Account:          req.GetAccount(),
		Name:             req.GetName(),
		CounterChannelID: req.GetCounterChannelId(),
		Status:           baseenum.Switch(req.GetStatus()),
		CreateTime:       proto.ParseTimeRange(req.GetCreateTime()),
		Paging:           proto.ParsePaging(req.GetPaging()),
		Sort:             proto.ParseSort(req.GetSort()),
	})
	if err != nil {
		return nil, err
	}

	return &v1.CountCounterChannelAccountReply{Count: int32(count)}, nil
}

func (c *CounterChannelAccountService) CreateCounterChannelAccount(ctx context.Context, req *v1.CreateCounterChannelAccountRequest) (*v1.CounterChannelAccount, error) {
	create, err := c.counterChannelAccountUseCase.Create(ctx, &systembiz.CounterChannelAccount{
		Status:           baseenum.SwitchOn,
		Name:             req.GetName(),
		Account:          req.GetAccount(),
		CounterChannelID: req.GetCounterChannelId(),
		Scene:            systemenum.CounterChannelScene(req.GetScene()),
	})
	if err != nil {
		return nil, err
	}
	return counterChannelAccountConvert(create), nil
}

func (c *CounterChannelAccountService) UpdateCounterChannelAccount(ctx context.Context, req *v1.UpdateCounterChannelAccountRequest) (*businessv1.BlankReply, error) {
	err := c.counterChannelAccountUseCase.Update(ctx, &systembiz.CounterChannelAccount{
		ID:               req.GetId(),
		Name:             req.GetName(),
		Account:          req.GetAccount(),
		CounterChannelID: req.GetCounterChannelId(),
		Scene:            systemenum.CounterChannelScene(req.GetScene()),
	})
	if err != nil {
		return nil, err
	}
	return &businessv1.BlankReply{}, nil
}

func counterChannelAccountConvert(info *systembiz.CounterChannelAccount) *v1.CounterChannelAccount {
	if info == nil {
		return nil
	}
	res := &v1.CounterChannelAccount{
		Id:               info.ID,
		Account:          info.Account,
		Name:             info.Name,
		CreateTime:       info.CreateTime,
		UpdateTime:       info.UpdateTime,
		Status:           uint32(info.Status),
		CounterChannelId: info.CounterChannelID,
		Scene:            string(info.Scene),
	}
	if info.CounterChannelInfo != nil {
		res.CounterChannelInfo = &v1.CounterChannelAccount_CounterChannel{
			Id:           info.CounterChannelInfo.ID,
			Name:         info.CounterChannelInfo.Name,
			Status:       uint32(info.CounterChannelInfo.Status),
			HealthStatus: uint32(info.CounterChannelInfo.HealthStatus),
			CreateTime:   info.CounterChannelInfo.CreateTime,
		}
	}
	return res
}
```

#### inner域

```go
package v1

import (
	"context"

	v1 "{module from go.mod}/internal/api/system/inner/v1"
	systembiz "{module from go.mod}/internal/biz/system"
	baseenum "{module from go.mod}/internal/enum/base"
	"{module from go.mod}/internal/pkg/proto"
)

type CounterChannelAccountService struct {
	v1.UnimplementedCounterChannelAccountServiceServer
	counterChannelAccountUseCase *systembiz.CounterChannelAccountUseCase
}

func NewCounterChannelAccountService(counterChannelAccountUseCase *systembiz.CounterChannelAccountUseCase) *CounterChannelAccountService {
	return &CounterChannelAccountService{counterChannelAccountUseCase: counterChannelAccountUseCase}
}

func (c *CounterChannelAccountService) GetCounterChannelAccount(ctx context.Context, req *v1.GetCounterChannelAccountRequest) (*v1.CounterChannelAccount, error) {
	counterChannel, err := c.counterChannelAccountUseCase.Get(ctx, req.GetId(), proto.ParseFilterConfig(req.GetFilterConfig())...)
	if err != nil {
		return nil, err
	}
	return counterChannelAccountConvert(counterChannel), nil
}

func (c *CounterChannelAccountService) GetCounterChannelAccountByAccount(ctx context.Context, req *v1.GetCounterChannelAccountByAccountRequest) (*v1.CounterChannelAccount, error) {
	counterChannel, err := c.counterChannelAccountUseCase.GetByAccount(ctx, req.GetAccount(), proto.ParseFilterConfig(req.GetFilterConfig())...)
	if err != nil {
		return nil, err
	}
	return counterChannelAccountConvert(counterChannel), nil
}

func (c *CounterChannelAccountService) ListCounterChannelAccount(ctx context.Context, req *v1.ListCounterChannelAccountRequest) (*v1.ListCounterChannelAccountReply, error) {
	list, err := c.counterChannelAccountUseCase.List(ctx, &systembiz.CounterChannelAccountFilter{
		IDList:           req.GetIdList(),
		Account:          req.GetAccount(),
		Name:             req.GetName(),
		CounterChannelID: req.GetCounterChannelId(),
		Status:           baseenum.Switch(req.GetStatus()),
		CreateTime:       proto.ParseTimeRange(req.GetCreateTime()),
		Paging:           proto.ParsePaging(req.GetPaging()),
		Sort:             proto.ParseSort(req.GetSort()),
		AccountList:      req.GetAccountList(),
	}, proto.ParseFilterConfig(req.GetFilterConfig())...)
	if err != nil {
		return nil, err
	}
	res := &v1.ListCounterChannelAccountReply{List: make([]*v1.CounterChannelAccount, 0)}
	for _, channel := range list {
		res.List = append(res.List, counterChannelAccountConvert(channel))
	}
	return res, nil
}

func (c *CounterChannelAccountService) MapCounterChannelAccount(ctx context.Context, req *v1.MapCounterChannelAccountRequest) (*v1.MapCounterChannelAccountReply, error) {
	dataMap, err := c.counterChannelAccountUseCase.Map(ctx, &systembiz.CounterChannelAccountFilter{
		IDList:           req.GetIdList(),
		Account:          req.GetAccount(),
		Name:             req.GetName(),
		CounterChannelID: req.GetCounterChannelId(),
		Status:           baseenum.Switch(req.GetStatus()),
		CreateTime:       proto.ParseTimeRange(req.GetCreateTime()),
		Paging:           proto.ParsePaging(req.GetPaging()),
		Sort:             proto.ParseSort(req.GetSort()),
	}, proto.ParseFilterConfig(req.GetFilterConfig())...)
	if err != nil {
		return nil, err
	}
	res := &v1.MapCounterChannelAccountReply{Map: map[uint32]*v1.CounterChannelAccount{}}
	for _, channel := range dataMap {
		res.Map[channel.ID] = counterChannelAccountConvert(channel)
	}
	return res, nil
}

func counterChannelAccountConvert(channel *systembiz.CounterChannelAccount) *v1.CounterChannelAccount {
	if channel == nil {
		return nil
	}

	return &v1.CounterChannelAccount{
		Id:                 channel.ID,
		Account:            channel.Account,
		Name:               channel.Name,
		Scene:              string(channel.Scene),
		CreateTime:         channel.CreateTime,
		UpdateTime:         channel.UpdateTime,
		Status:             uint32(channel.Status),
		CounterChannelId:   channel.CounterChannelID,
		CounterChannelInfo: counterChannelConvert(channel.CounterChannelInfo),
	}
}
```
**规范要点**：
- 只做请求参数校验和响应转换
- 不包含业务逻辑
- 依赖 Biz 层的 UseCase
- 使用 `proto` 包工具函数转换通用参数

### Biz 层 (internal/biz)

**职责**：核心业务逻辑，领域模型定义
```go
package system

import (
	"context"

	"{module from go.mod}/internal/biz/kit"
	"{module from go.mod}/internal/enum/base"
	systemenum "{module from go.mod}/internal/enum/system"
	systemerror "{module from go.mod}/internal/error/system"
	"{module from go.mod}/internal/pkg/filter"
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
	// 拓展参数
	Extra map[string]interface{} `json:"extra"`
	// 通道信息
	CounterChannelInfo *CounterChannel `json:"counter_channel_info"`
}

type CounterChannelAccountFilter struct {
	IDList             []uint32
	IDNotInList        []uint32
	Account            string
	Name               string
	CounterChannelID   uint32
	Status             base.Switch
	Scene              systemenum.CounterChannelScene
	CreateTime         filter.TimeRange
	Paging             filter.Paging
	Sort               filter.Sort
	AccountList        []string
	TradeAccountIdList []uint32
}

type CounterChannelAccountRepo interface {
	ListCounterChannelAccount(ctx context.Context, filter *CounterChannelAccountFilter, opts ...filter.Option) ([]*CounterChannelAccount, error)
	GetCounterChannelAccount(ctx context.Context, ID uint32, opts ...filter.Option) (*CounterChannelAccount, error)
	GetCounterChannelAccountByAccount(ctx context.Context, account string, opts ...filter.Option) (*CounterChannelAccount, error)
	GetUsableCounterChannelAccount(ctx context.Context, ID uint32) (*CounterChannelAccount, error)
	PageListCounterChannelAccount(ctx context.Context, filter *CounterChannelAccountFilter, opts ...filter.Option) ([]*CounterChannelAccount, int, error)
	CountCounterChannelAccount(ctx context.Context, filter *CounterChannelAccountFilter, opts ...filter.Option) (int, error)
	CreateCounterChannelAccount(ctx context.Context, info *CounterChannelAccount) (*CounterChannelAccount, error)
	UpdateCounterChannelAccount(ctx context.Context, info *CounterChannelAccount) error
	DeleteCounterChannelAccount(ctx context.Context, info *CounterChannelAccount) error
	CheckCounterChannelAccountExists(ctx context.Context, filter *CounterChannelAccountFilter) error
	MapCounterChannelAccount(ctx context.Context, filter *CounterChannelAccountFilter, opts ...filter.Option) (map[uint32]*CounterChannelAccount, error)
	BatchDeleteCounterChannelAccountCache(ctx context.Context, IDList []uint32) error
}

type CounterChannelAccountUseCase struct {
	counterChannelAccountRepo CounterChannelAccountRepo
	transaction               kit.Transaction
	counterChannelUseCase     *CounterChannelUseCase
}

func NewCounterChannelAccountUseCase(
	counterChannelAccountRepo CounterChannelAccountRepo,
	transaction kit.Transaction,
	counterChannelUseCase *CounterChannelUseCase,
) *CounterChannelAccountUseCase {
	return &CounterChannelAccountUseCase{
		counterChannelAccountRepo: counterChannelAccountRepo,
		transaction:               transaction,
		counterChannelUseCase:     counterChannelUseCase,
	}
}

func (c *CounterChannelAccountUseCase) Get(ctx context.Context, ID uint32, opts ...filter.Option) (*CounterChannelAccount, error) {
	return c.counterChannelAccountRepo.GetCounterChannelAccount(ctx, ID, opts...)
}

func (c *CounterChannelAccountUseCase) GetByAccount(ctx context.Context, account string, opts ...filter.Option) (*CounterChannelAccount, error) {
	return c.counterChannelAccountRepo.GetCounterChannelAccountByAccount(ctx, account, opts...)
}

func (c *CounterChannelAccountUseCase) List(ctx context.Context, filter *CounterChannelAccountFilter, opts ...filter.Option) ([]*CounterChannelAccount, error) {
	return c.counterChannelAccountRepo.ListCounterChannelAccount(ctx, filter, opts...)
}

func (c *CounterChannelAccountUseCase) PageList(ctx context.Context, filter *CounterChannelAccountFilter, opts ...filter.Option) ([]*CounterChannelAccount, int, error) {
	return c.counterChannelAccountRepo.PageListCounterChannelAccount(ctx, filter, opts...)
}

func (c *CounterChannelAccountUseCase) Count(ctx context.Context, filter *CounterChannelAccountFilter, opts ...filter.Option) (int, error) {
	return c.counterChannelAccountRepo.CountCounterChannelAccount(ctx, filter, opts...)
}

func (c *CounterChannelAccountUseCase) Create(ctx context.Context, info *CounterChannelAccount) (res *CounterChannelAccount, err error) {
	channelInfo, err := c.counterChannelUseCase.Get(ctx, info.CounterChannelID)
	if err != nil {
		return nil, err
	}

	if channelInfo.Status == base.SwitchOff {
		return nil, systemerror.ErrorCounterChannelDisabled(ctx)
	}
	//err = c.counterChannelAccountRepo.CheckCounterChannelAccountExists(ctx, &CounterChannelAccountFilter{Account: info.Account, CounterChannelID: info.CounterChannelID, Scene: info.Scene})
	//if err != nil {
	//	return nil, err
	//}

	err = c.transaction.InTx(ctx, func(ctx context.Context) error {
		res, err = c.counterChannelAccountRepo.CreateCounterChannelAccount(ctx, info)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	err = c.counterChannelAccountRepo.BatchDeleteCounterChannelAccountCache(ctx, []uint32{res.ID})
	if err != nil {
		return nil, err
	}

	return
}

func (c *CounterChannelAccountUseCase) Update(ctx context.Context, info *CounterChannelAccount) error {
	channelAccountInfo, err := c.Get(ctx, info.ID, filter.WithRelation(systemenum.RelationCounterChannelInfo))
	if err != nil {
		return err
	}

	if channelAccountInfo.Name != info.Name || channelAccountInfo.Account != info.Account {
		// 修改状态  不校验
		if channelAccountInfo.CounterChannelInfo.Status == base.SwitchOff {
			return systemerror.ErrorCounterChannelDisabled(ctx)
		}
	}

	//err = c.counterChannelAccountRepo.CheckCounterChannelAccountExists(ctx, &CounterChannelAccountFilter{Account: info.Account, CounterChannelID: info.CounterChannelID, IDNotInList: []uint32{info.ID}, Scene: info.Scene})
	//if err != nil {
	//	return err
	//}

	channelAccountInfo.Name = info.Name
	channelAccountInfo.Account = info.Account
	channelAccountInfo.CounterChannelID = info.CounterChannelID
	channelAccountInfo.Scene = info.Scene
	err = c.counterChannelAccountRepo.UpdateCounterChannelAccount(ctx, channelAccountInfo)
	if err != nil {
		return err
	}

	err = c.counterChannelAccountRepo.BatchDeleteCounterChannelAccountCache(ctx, []uint32{info.ID})
	if err != nil {
		return err
	}

	return nil
}

func (c *CounterChannelAccountUseCase) Delete(ctx context.Context, info *CounterChannelAccount) error {
	err := c.counterChannelAccountRepo.DeleteCounterChannelAccount(ctx, info)
	if err != nil {
		return err
	}

	err = c.counterChannelAccountRepo.BatchDeleteCounterChannelAccountCache(ctx, []uint32{info.ID})
	if err != nil {
		return err
	}

	return nil
}

func (c *CounterChannelAccountUseCase) Map(ctx context.Context, filter *CounterChannelAccountFilter, opts ...filter.Option) (map[uint32]*CounterChannelAccount, error) {
	return c.counterChannelAccountRepo.MapCounterChannelAccount(ctx, filter, opts...)
}

func (c *CounterChannelAccountUseCase) BatchCheckChannelAccountStatus(ctx context.Context, IDList []uint32) error {
	if IDList == nil {
		return nil
	}
	channelAccounts, err := c.counterChannelAccountRepo.ListCounterChannelAccount(ctx, &CounterChannelAccountFilter{IDList: IDList}, filter.WithRelation(systemenum.RelationCounterChannelInfo))
	if err != nil {
		return err
	}

	for _, channelAccount := range channelAccounts {
		if channelAccount.Status == base.SwitchOff {
			return systemerror.ErrorCounterChannelAccountDisabled(ctx, map[string]interface{}{
				"account": channelAccount.Account,
			})
		}
		if channelAccount.CounterChannelInfo == nil {
			return systemerror.ErrorCounterChannelNotFound(ctx)
		}
		if channelAccount.CounterChannelInfo.Status == base.SwitchOff {
			return systemerror.ErrorCounterChannelDisabled(ctx)
		}
	}

	return nil
}
```
**规范要点**：
- 定义领域模型（Domain Model）
- 定义 Repository 接口
- 实现业务用例（UseCase）
- 不依赖具体数据实现
- 业务逻辑集中在此层

### Data 层 (internal/data)

**职责**：数据访问实现，Repository 实现
```go
package system

import (
	"context"
	"errors"
	"time"

	"gitlab.linksoft.cn/TeamA/GoPackage/helper"

	jsoniter "github.com/json-iterator/go"
	systembiz "{module from go.mod}/internal/biz/system"
	"{module from go.mod}/internal/data/ent"
	"{module from go.mod}/internal/data/ent/counterchannelaccount"
	"{module from go.mod}/internal/data/kit"
	baseenum "{module from go.mod}/internal/enum/base"
	systemenum "{module from go.mod}/internal/enum/system"
	systemerror "{module from go.mod}/internal/error/system"
	"{module from go.mod}/internal/pkg/filter"
)

type counterChannelAccountRepo struct {
	counterChannelRepo systembiz.CounterChannelRepo
	data               *kit.Data
}

func NewCounterChannelAccountRepo(
	counterChannelRepo systembiz.CounterChannelRepo,
	data *kit.Data,
) systembiz.CounterChannelAccountRepo {
	return &counterChannelAccountRepo{
		counterChannelRepo: counterChannelRepo,
		data:               data,
	}
}

func (r *counterChannelAccountRepo) GetCounterChannelAccount(ctx context.Context, ID uint32, opts ...filter.Option) (res *systembiz.CounterChannelAccount, err error) {
	cfg := filter.NewConfig(opts...)
	if !cfg.FromCache {
		res, err = r.getCounterChannelAccount(ctx, ID, opts...)
	} else {
		res, err = r.getCacheCounterChannelAccount(ctx, ID, opts...)
	}
	if err != nil {
		return nil, err
	}

	err = r.serviceRelation(ctx, res, opts...)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (r *counterChannelAccountRepo) GetCounterChannelAccountByFilter(ctx context.Context, filter *systembiz.CounterChannelAccountFilter, opts ...filter.Option) (res *systembiz.CounterChannelAccount, err error) {
	query := r.data.Db.CounterChannelAccount(ctx).Query()

	query = r.parseFilter(query, filter)
	query = r.queryConfig(query, opts...)

	info, err := query.First(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, systemerror.ErrorCounterChannelAccountNotFound(ctx)
		}
		return nil, err
	}

	res = r.queryRelation(counterChannelAccountConvert(info), info.Edges)

	err = r.serviceRelation(ctx, res, opts...)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (r *counterChannelAccountRepo) GetCounterChannelAccountByAccount(ctx context.Context, account string, opts ...filter.Option) (*systembiz.CounterChannelAccount, error) {
	query := r.data.Db.CounterChannelAccount(ctx).Query().Where(counterchannelaccount.AccountEQ(account))
	query = r.queryConfig(query, opts...)
	first, err := query.First(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, systemerror.ErrorCounterChannelAccountNotFound(ctx)
		}
		return nil, err
	}

	res := r.queryRelation(counterChannelAccountConvert(first), first.Edges)

	err = r.serviceRelation(ctx, res, opts...)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (r *counterChannelAccountRepo) ListCounterChannelAccount(ctx context.Context, f *systembiz.CounterChannelAccountFilter, opts ...filter.Option) (res []*systembiz.CounterChannelAccount, err error) {
	cfg := filter.NewConfig(opts...)
	if !cfg.FromCache {
		res, err = r.listCounterChannelAccount(ctx, f, opts...)
	} else {
		res, err = r.listCacheCounterChannelAccount(ctx, f, opts...)
	}
	if err != nil {
		return nil, err
	}

	err = r.serviceRelation(ctx, res, opts...)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (r *counterChannelAccountRepo) PageListCounterChannelAccount(ctx context.Context, filter *systembiz.CounterChannelAccountFilter, opts ...filter.Option) ([]*systembiz.CounterChannelAccount, int, error) {
	count, err := r.CountCounterChannelAccount(ctx, filter, opts...)
	if err != nil {
		return nil, 0, err
	}

	if count == 0 {
		return []*systembiz.CounterChannelAccount{}, 0, nil
	}
	list, err := r.ListCounterChannelAccount(ctx, filter, opts...)

	if err != nil {
		return nil, 0, err
	}

	return list, count, nil
}

func (r *counterChannelAccountRepo) CountCounterChannelAccount(ctx context.Context, filter *systembiz.CounterChannelAccountFilter, opts ...filter.Option) (int, error) {
	filter.Paging.Disable()
	defer filter.Paging.Enable()
	query := r.data.Db.CounterChannelAccount(ctx).Query()
	query = r.parseFilter(query, filter)
	count, err := query.Count(ctx)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r *counterChannelAccountRepo) MapCounterChannelAccount(ctx context.Context, filter *systembiz.CounterChannelAccountFilter, opts ...filter.Option) (map[uint32]*systembiz.CounterChannelAccount, error) {
	channelList, err := r.ListCounterChannelAccount(ctx, filter, opts...)
	if err != nil {
		return nil, err
	}

	res := map[uint32]*systembiz.CounterChannelAccount{}
	for _, info := range channelList {
		res[info.ID] = info
	}

	return res, nil
}

func (r *counterChannelAccountRepo) CreateCounterChannelAccount(ctx context.Context, info *systembiz.CounterChannelAccount) (*systembiz.CounterChannelAccount, error) {
	save, err := r.data.Db.CounterChannelAccount(ctx).
		Create().
		SetCounterChannelID(info.CounterChannelID).
		SetStatus(info.Status).
		SetAccount(info.Account).
		SetName(info.Name).
		SetScene(info.Scene).
		SetExtra(info.Extra).
		Save(ctx)
	if err != nil {
		return nil, err
	}
	return counterChannelAccountConvert(save), nil
}

func (r *counterChannelAccountRepo) UpdateCounterChannelAccount(ctx context.Context, info *systembiz.CounterChannelAccount) error {
	_, err := r.data.Db.CounterChannelAccount(ctx).
		Update().
		Where(counterchannelaccount.IDEQ(info.ID)).
		SetCounterChannelID(info.CounterChannelID).
		SetAccount(info.Account).
		SetName(info.Name).
		SetScene(info.Scene).
		SetExtra(info.Extra).
		Save(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (r *counterChannelAccountRepo) DeleteCounterChannelAccount(ctx context.Context, info *systembiz.CounterChannelAccount) error {
	err := r.data.Db.CounterChannelAccount(ctx).DeleteOneID(info.ID).Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (r *counterChannelAccountRepo) GetUsableCounterChannelAccount(ctx context.Context, accountID uint32) (*systembiz.CounterChannelAccount, error) {
	channelAccountInfo, err := r.GetCounterChannelAccount(
		ctx,
		accountID,
		filter.WithFromCache(true),
		filter.WithRelation(systemenum.RelationCounterChannelInfo, filter.WithRelationConfigOpts(
			filter.WithFromCache(true),
		)),
	)
	if err != nil {
		return nil, err
	}
	if channelAccountInfo.Status != baseenum.SwitchOn {
		return nil, systemerror.ErrorCounterChannelAccountDisabled(ctx, map[string]interface{}{"account": channelAccountInfo.Account})
	}

	if channelAccountInfo.CounterChannelInfo == nil {
		return nil, systemerror.ErrorCounterChannelNotFound(ctx)
	}
	if channelAccountInfo.CounterChannelInfo.HealthStatus != systemenum.HealthStatusNormal {
		return nil, systemerror.ErrorCounterChannelUnhealthy(ctx)
	}

	if channelAccountInfo.CounterChannelInfo.Status != baseenum.SwitchOn {
		return nil, systemerror.ErrorCounterChannelDisabled(ctx)
	}

	return channelAccountInfo, nil
}

func (r *counterChannelAccountRepo) CheckCounterChannelAccountExists(ctx context.Context, filter *systembiz.CounterChannelAccountFilter) error {
	query := r.data.Db.CounterChannelAccount(ctx).Query()
	query = query.Where(counterchannelaccount.AccountEQ(filter.Account)).Where(counterchannelaccount.CounterChannelIDEQ(filter.CounterChannelID)).Where(counterchannelaccount.SceneEQ(filter.Scene))
	if filter.IDNotInList != nil {
		query = query.Where(counterchannelaccount.IDNotIn(filter.IDNotInList...))
	}

	exist, err := query.Exist(ctx)
	if err != nil {
		return err
	}
	if exist {
		return systemerror.ErrorCounterChannelAccountExists(ctx)
	}
	return nil
}

func (r *counterChannelAccountRepo) BatchDeleteCounterChannelAccountCache(ctx context.Context, IDList []uint32) error {
	keyList := helper.SliceConvert(IDList, func(item uint32) string {
		return systemenum.RedisKeyCounterChannelAccountInfo.Build(item)
	})

	_, err := r.data.Redis.Del(ctx, keyList...)
	if err != nil {
		return err
	}
	return nil
}

func (r *counterChannelAccountRepo) getCounterChannelAccount(ctx context.Context, ID uint32, opts ...filter.Option) (*systembiz.CounterChannelAccount, error) {
	query := r.data.Db.CounterChannelAccount(ctx).Query().Where(counterchannelaccount.ID(ID))
	query = r.queryConfig(query, opts...)
	first, err := query.First(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, systemerror.ErrorCounterChannelAccountNotFound(ctx)
		}
		return nil, err
	}

	res := r.queryRelation(counterChannelAccountConvert(first), first.Edges)

	return res, nil
}

func (r *counterChannelAccountRepo) getCacheCounterChannelAccount(ctx context.Context, ID uint32, opts ...filter.Option) (*systembiz.CounterChannelAccount, error) {
	key := systemenum.RedisKeyCounterChannelAccountInfo.Build(ID)

	cache, err := r.data.Redis.GetOrInit(ctx, key, func() (string, error) {
		res := ""

		info, err := r.getCounterChannelAccount(ctx, ID, append(opts, filter.WithTransField(true))...)
		if err != nil && !errors.Is(err, systemerror.ErrorCounterChannelAccountNotFound(ctx)) {
			return "", err
		}

		if info != nil {
			res, err = jsoniter.MarshalToString(info)
			if err != nil {
				return "", err
			}
		}
		return res, nil
	}, 24*time.Hour, time.Minute)
	if err != nil {
		return nil, err
	}
	if cache == "" {
		return nil, systemerror.ErrorCounterChannelAccountNotFound(ctx)
	}

	res := &systembiz.CounterChannelAccount{}

	err = jsoniter.UnmarshalFromString(cache, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (r *counterChannelAccountRepo) listCounterChannelAccount(ctx context.Context, filter *systembiz.CounterChannelAccountFilter, opts ...filter.Option) ([]*systembiz.CounterChannelAccount, error) {
	query := r.data.Db.CounterChannelAccount(ctx).Query()

	query = r.parseFilter(query, filter)
	query = r.queryConfig(query, opts...)

	list, err := query.All(ctx)
	if err != nil {
		return nil, err
	}

	var res []*systembiz.CounterChannelAccount

	for _, channel := range list {
		res = append(res, r.queryRelation(counterChannelAccountConvert(channel), channel.Edges))
	}

	return res, nil
}

func (r *counterChannelAccountRepo) listCacheCounterChannelAccount(ctx context.Context, f *systembiz.CounterChannelAccountFilter, opts ...filter.Option) ([]*systembiz.CounterChannelAccount, error) {
	if f.IDList == nil || len(f.IDList) == 0 {
		return []*systembiz.CounterChannelAccount{}, nil
	}

	keyList := helper.SliceConvert(f.IDList, func(item uint32) string {
		return systemenum.RedisKeyCounterChannelAccountInfo.Build(item)
	})

	cacheList, err := r.data.Redis.ListOrInit(ctx, keyList, func(keyList []string) (map[string]string, error) {
		IDList := make([]uint32, len(keyList))
		for index, key := range keyList {
			err := systemenum.RedisKeyCounterChannelAccountInfo.Parse(key, &IDList[index])
			if err != nil {
				return nil, err
			}
		}

		list, err := r.listCounterChannelAccount(ctx, &systembiz.CounterChannelAccountFilter{IDList: IDList}, append(opts, filter.WithTransField(true))...)
		if err != nil {
			return nil, err
		}

		res := map[string]string{}

		for _, info := range list {
			key := systemenum.RedisKeyCounterChannelAccountInfo.Build(info.ID)

			res[key], err = jsoniter.MarshalToString(info)
			if err != nil {
				return nil, err
			}
		}
		return res, nil
	}, 24*time.Hour, time.Minute)
	if err != nil {
		return nil, err
	}

	res := make([]*systembiz.CounterChannelAccount, 0, len(cacheList))
	for _, cache := range cacheList {
		info := &systembiz.CounterChannelAccount{}
		err = jsoniter.UnmarshalFromString(cache, info)
		if err != nil {
			return nil, err
		}
		res = append(res, info)
	}

	return res, nil
}

func (r *counterChannelAccountRepo) parseFilter(query *ent.CounterChannelAccountQuery, filter *systembiz.CounterChannelAccountFilter) *ent.CounterChannelAccountQuery {
	if filter.IDList != nil {
		query = query.Where(counterchannelaccount.IDIn(filter.IDList...))
	}
	if filter.IDNotInList != nil {
		query = query.Where(counterchannelaccount.IDNotIn(filter.IDNotInList...))
	}

	if filter.Status != baseenum.SwitchAll {
		query = query.Where(counterchannelaccount.StatusEQ(filter.Status))
	}

	if filter.Name != "" {
		query = query.Where(counterchannelaccount.NameContains(filter.Name))
	}

	if filter.Account != "" {
		query = query.Where(counterchannelaccount.AccountContains(filter.Account))
	}

	if filter.CounterChannelID != 0 {
		query = query.Where(counterchannelaccount.CounterChannelIDEQ(filter.CounterChannelID))
	}

	if filter.AccountList != nil {
		query = query.Where(counterchannelaccount.AccountIn(filter.AccountList...))
	}

	query.Modify(
		filter.CreateTime.ModifyFn(counterchannelaccount.FieldCreateTime),
		filter.Sort.ModifyFn(counterchannelaccount.ValidColumn),
		filter.Paging.ModifyFn(),
	)
	return query
}

func (r *counterChannelAccountRepo) serviceRelation(ctx context.Context, data interface{}, opts ...filter.Option) error {
	if data == nil {
		return nil
	}
	cfg := filter.NewConfig(opts...)

	if info, ok := data.(*systembiz.CounterChannelAccount); ok {
		data = []*systembiz.CounterChannelAccount{info}
	}
	list := data.([]*systembiz.CounterChannelAccount)

	if relation, ok := cfg.Relations[systemenum.RelationCounterChannelInfo]; ok {
		channelIDList := helper.SliceNonDuplicateColumn(list, func(item *systembiz.CounterChannelAccount) uint32 {
			return item.CounterChannelID
		})

		if len(channelIDList) > 0 {
			counterChannelMap, err := r.counterChannelRepo.MapCounterChannel(ctx, &systembiz.CounterChannelFilter{IDList: channelIDList}, relation.ConfigOpts...)
			if err != nil {
				return err
			}
			for _, info := range list {
				info.CounterChannelInfo = counterChannelMap[info.CounterChannelID]
			}
		}
	}

	return nil
}

func (r *counterChannelAccountRepo) queryConfig(query *ent.CounterChannelAccountQuery, opts ...filter.Option) *ent.CounterChannelAccountQuery {
	return query
}

func (r *counterChannelAccountRepo) queryRelation(info *systembiz.CounterChannelAccount, edges ent.CounterChannelAccountEdges) *systembiz.CounterChannelAccount {
	return info
}

func counterChannelAccountConvert(info *ent.CounterChannelAccount) *systembiz.CounterChannelAccount {
	if info == nil {
		return nil
	}
	return &systembiz.CounterChannelAccount{
		ID:               info.ID,
		CreateTime:       info.CreateTime,
		UpdateTime:       info.UpdateTime,
		CounterChannelID: info.CounterChannelID,
		Status:           info.Status,
		Account:          info.Account,
		Name:             info.Name,
		Scene:            info.Scene,
		Extra:            info.Extra,
	}
}
```

**规范要点**：
- 实现 Biz 层定义的 Repository 接口
- 封装对数据源的访问
- 使用 Ent 进行数据库操作
- 负责 Ent 实体与领域模型的转换
- 封装查询过滤逻辑（parseFilter 方法）
- 支持缓存读取（如有需要）
- 支持远程数据同步（如有需要）
- 支持多语言字段处理（如有需要）
- 封装对内部数据关联的处理（queryConfig，queryRelation）
- 封装对外部数据关联的处理（serviceRelation）

## 依赖注入 (Wire)

每层都需要定义 `ProviderSet`：

### biz

- internal/biz/{service name}/{service name}.go

```go
package system

var ProviderSet = wire.NewSet(
    NewCounterChannelAccountUseCase,
)

```

### data
- internal/data/{service name}/{service name}.go

```go
package system

import "github.com/google/wire"

var ProviderSet = wire.NewSet(
	NewCounterChannelAccountRepo,
)

```

### service

#### admin
- internal/service/admin/admin.go

```go
package admin

import (
	"github.com/google/wire"
	"{module from go.mod}/internal/service/admin/v1"
)

var ProviderSet = wire.NewSet(
	v1.NewCounterChannelAccountService,
)

```

#### inner
- internal/service/inner/inner.go

```go
package admin

import (
	"github.com/google/wire"
	"{module from go.mod}/internal/service/inner/v1"
)

var ProviderSet = wire.NewSet(
	v1.NewCounterChannelAccountService,
)

```

#### open
- internal/service/open/open.go

```go
package admin

import (
	"github.com/google/wire"
	"{module from go.mod}/internal/service/open/v1"
)

var ProviderSet = wire.NewSet(
	v1.NewCounterChannelAccountService,
)

```

## 依赖方向

```
Service -> Biz <- Data
             ↑
         Interface
```

- Service 依赖 Biz 的 UseCase
- Data 实现 Biz 定义的 Repository 接口
- Biz 层不依赖 Data 层的具体实现

## 事务处理

使用 `kit.Transaction` 接口进行事务管理：

```go
// internal/biz/kit/transaction.go
package kit

import "context"

type Transaction interface {
	InTx(context.Context, func(ctx context.Context) error) error
}

// 在 UseCase 中使用
func (u *OrderUseCase) CreateOrder(ctx context.Context, order *Order) error {
	return u.tx.InTx(ctx, func(ctx context.Context) error {
		// 事务内的操作
		if _, err := u.orderRepo.Create(ctx, order); err != nil {
			return err
		}
		// 更多操作...
		return nil
	})
}

```

## 标准 CRUD 方法命名

| 操作 | Biz 层方法 | Repo 接口方法 |
|------|-----------|--------------|
| 单条查询 | `Get` | `GetXxx` |
| 列表查询 | `List` | `ListXxx` |
| 分页查询 | `PageList` | `PageListXxx` |
| 计数查询 | `Count` | `CountXxx` |
| 创建 | `Create` | `CreateXxx` |
| 更新 | `Update` | `UpdateXxx` |
| 批量删除 | `BatchDelete` | `BatchDeleteXxx` |
