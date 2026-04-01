# Domain

## 作用范围

本文用于说明 Kratos 业务项目中的核心领域组织方式，重点回答以下问题：

- `usecase`、`repo`、`kit`、`data` 各自是什么
- 它们分别负责什么，不负责什么
- 它们之间的依赖方向是什么
- 事务、relation、第三方依赖、输入输出对象应如何围绕领域主线组织

当问题属于以下场景时，应优先查看本文：

- 新增或重构业务能力
- 判断逻辑应落在 `usecase`、`repo`、`kit` 还是 `data`
- 判断事务边界应放在哪一层
- 判断 relation 应在哪一层收口
- 判断第三方依赖应直接归属 `repo`，还是应下沉到 `kit` 提供支撑能力
- 判断应用层输入输出对象与领域对象的关系

本文默认已完成实体与聚合根识别；实体与聚合根的定义以 `aggregate-root.md` 为准。

## 规则

### 文件

- 实体必须有自己的单独文件，必须定义对应 `Repo`
- `UseCase` 只在实体需要对外暴露业务能力时定义
- 一旦某个实体定义了 `UseCase`，该实体、对应 `Repo`、对应 `UseCase` 必须收敛在同一文件
- 工程文件中不要出现无业务语义的函数，如有应下沉到 `internal/pkg`

### `service`

- 领域主线保持 `service -> usecase -> repo -> kit/data/depend`，禁止越层调用与职责漂移

### `usecase`

- 同一业务能力默认只保留一个 `UseCase owner`；发起、重试、回调消费、结果同步、后置处理仍属同一 `UseCase owner`
- `usecase` 不按流程阶段、回调方向、结果类型或读写动作拆成多个平级 `UseCase`
- 如果新增类型只是承接同一业务能力中的“发起”“重试”“回调消费”“结果同步”“后置处理”“补偿查询”“状态推进”等动作，通常属于过度拆分，应并回原 `UseCase`
- 事务边界优先留在 `usecase`

### `repo`

- relation 统一收口在 `repo`
- 第三方依赖先判断是领域依赖还是基础设施能力；稳定值域一旦进入 `biz`，应收敛为项目内部语义

### `kit`

- 已存在于 `kit` 的能力，`repo` 应优先通过注入复用
- 如果现有 `kit` 能力不足以承载真实语义，应先补齐 `kit` 的边界，再回到 `repo` 使用
- `repo` 可以持有业务归属的第三方配置，但不应为了复制 `kit` 的全局基础设施装配而再次注入 `conf.Data`
- 问题如果本质是 HTTP client / cache / transaction / eventbus 等全局基础设施能力缺口，应补 `kit`，不归某个 `repo` 私有实现

### `data`

- `data` 只承接 `repo` 的落地实现与依赖访问，不反向承载业务编排

## 领域主线

Kratos 业务项目中的领域主线通常围绕以下关系展开：

```text
Service -> UseCase -> Repo -> Kit/Data/Depend
```

其中：

- `service` 负责协议适配
- `usecase` 负责业务编排
- `repo` 负责领域依赖抽象、查询组织与 relation 收口
- `kit` 负责为 `repo/data` 提供基础设施能力支撑
- `data` 负责 repo 落地实现与依赖访问

分层的目标不是目录整齐，而是让协议变化、业务变化、基础设施变化和数据变化可以分开演进。

## 文件组织

识别完实体、聚合根或稳定第三方边界主题后，应继续确定对应的文件组织方式。

### 实体文件

- 一个实体必须有一个对应文件
- 实体结构与对应 `Repo` 接口应定义在同一文件中
- 实体可以是聚合根，也可以不是聚合根
- 当该实体需要对外暴露业务能力时，可在同一文件中继续定义对应 `UseCase`
- `service` 通过调用 `UseCase` 对外暴露能力，不直接调用 `Repo`

### 第三方边界主题文件

- 如果某个第三方依赖已经形成稳定边界主题，也应独立建立对应文件
- 该文件用于收敛对应 `Repo` 接口、请求对象、响应对象和相关稳定结构
- 这类文件不等同于本地领域实体文件，但同样应保持主题收敛

## 四个核心角色

下面四个角色是本文的主结构。阅读顺序应固定为：

1. `usecase`
2. `repo`
3. `kit`
4. `data`

后续事务、relation、第三方依赖等规则，都是围绕这四个角色展开，而不是替代这四个角色本身。

### `usecase`

#### 定位

`usecase` 是业务编排中心。

#### 负责什么

- 业务编排
- 事务边界
- 权限校验
- 状态流转
- 协调多个 `repo`
- 接收业务 filter 和 relation 需求

#### 不负责什么

- 不直接编写 DB 查询细节
- 不手工组装 relation
- 不吸收协议映射逻辑
- 不重复实现基础设施能力

#### 示例

```go
func (u *AccountUseCase) Submit(ctx context.Context, req *SubmitAccountRequest) error {
	return u.tx.ExecTx(ctx, func(ctx context.Context) error {
		account, err := u.accountRepo.GetAccount(ctx, req.ID)
		if err != nil {
			return err
		}
		if err := account.Submit(); err != nil {
			return err
		}
		return u.accountRepo.UpdateAccount(ctx, account)
	})
}
```

#### 文件收敛说明

`usecase` 围绕实体对外暴露的业务能力组织。
如果某个实体定义了 `UseCase`，则该实体的结构、对应 `Repo` 接口、对应 `UseCase` 结构与核心方法应继续收敛在同一实体文件内，避免把同一主题机械拆成 `xxx.go` / `xxx_usecase.go`。

简例：

```text
account_flow_page.go
├─ type Accountflow_page
├─ type AccountFlowPageRepo interface
├─ type AccountFlowPageUseCase struct
└─ func (u *AccountFlowPageUseCase) ...
```

### `repo`

#### 定位

`repo` 是领域依赖抽象与查询组织中心。

#### 负责什么

- 查询组织
- relation 装配
- 数据回填
- 领域依赖访问
- 远程 relation 补全

#### 不负责什么

- 不承载整聚合更新编排
- 不主导业务状态流转
- 不主导事务边界
- 不在内部重复 new 应由 `kit` 提供的基础设施

#### 示例

```go
func (u *AccountUseCase) GetAccount(ctx context.Context, id uint32, opts ...filter.Option) (*Account, error) {
	return u.accountRepo.GetAccount(ctx, id, opts...)
}
```

```go
func (r *accountRepo) GetAccount(ctx context.Context, id uint32, opts ...filter.Option) (*openbiz.Account, error) {
	query := r.data.Db.Account(ctx).Query()
	query = r.queryConfig(query, &openbiz.AccountFilter{IDList: []uint32{id}}, opts...)

	info, err := query.First(ctx)
	if err != nil {
		return nil, err
	}

	res := r.queryRelation(accountConvert(info), info.Edges)
	if err := r.serviceRelation(ctx, res, opts...); err != nil {
		return nil, err
	}
	return res, nil
}
```

#### Repo 组织方式

Repo 查询通常围绕以下几个稳定阶段组织：

- 查询条件收敛
- 本地 relation 配置
- 本地 relation 装配
- 远程 relation 补全

这些阶段可以分别表现为：

- `parseFilter`
- `queryConfig`
- `queryRelation`
- `serviceRelation`

关键约束如下：

- 业务 filter 只表达业务查询语义
- relation 是否加载由显式 relation 配置控制
- 本地 relation 与远程 relation 分开处理
- 远程 relation 必须优先采用批量收集、批量查询、批量回填

#### `queryRelation`

`queryRelation` 用于把本地查询结果中的结构化 relation 映射回业务对象。

适用场景：

- relation 已经建模在 ent schema 中
- 查询阶段已经通过 `WithXxx()` 预加载了对应 edges
- 当前工作是把 `info.Edges` 转成业务对象字段，而不是再发起额外查询

不适用场景：

- 关系并没有建模成 ent edge
- 需要在主查询之后再访问远程依赖或其他 repo 才能补齐数据
- 当前逻辑本质是“查询后补充”，而不是“读取已加载 edges”

使用要点：

- `queryRelation` 只做本地 edges 到业务对象的映射
- 不在 `queryRelation` 中发起远程调用
- 不在 `queryRelation` 中承载业务编排
- 它和 `queryConfig` 成对出现，前者负责声明加载哪些本地 relation，后者负责把已加载结果映射回来

示例
```go
func (r *accountFlowPageRepo) queryConfig(query *ent.AccountFlowPageQuery, opts ...filter.Option) *ent.AccountFlowPageQuery {
	cfg := filter.NewConfig(opts...)
	if _, ok := cfg.Relations[openenum.AccountFlowPageFieldRelation]; ok {
		query = query.WithFields()
	}
	return query
}

func (r *accountFlowPageRepo) queryRelation(info *openbiz.AccountFlowPage, edges ent.AccountFlowPageEdges) *openbiz.AccountFlowPage {
	if info == nil {
		return nil
	}
	if len(edges.Fields) > 0 {
		info.FieldList = make([]*openbiz.AccountFlowField, 0, len(edges.Fields))
		for _, fieldInfo := range edges.Fields {
			info.FieldList = append(info.FieldList, flowFieldConvert(fieldInfo))
		}
	} else if info.FieldList == nil {
		info.FieldList = []*openbiz.AccountFlowField{}
	}
	return info
}
```

#### `serviceRelation`

`serviceRelation` 用于在主查询完成后，补齐无法通过本地 ent edges 直接表达的 relation。

适用场景：

- schema 中不存在合适的 ent `edge`
- 数据来源是远程依赖、第三方服务，或本服务内另一类更适合作为独立 repo 查询能力的数据
- 需要在主查询完成后，根据已有字段批量收集、批量查询、批量回填

不适用场景：

- 当前关系已经适合通过 ent schema 的 `edge` 表达
- 当前只是把 `info.Edges` 映射成业务对象
- 当前逻辑已经开始承载业务编排、状态流转或事务主导

使用要点：

- `serviceRelation` 必须由显式 relation 配置触发，不能隐式补查
- `serviceRelation` 只做 relation 补全，不做业务编排
- 优先使用批量收集、批量查询、批量回填，避免 N+1
- 数据来源既可以是 `depend repo`，也可以是本服务内部但不适合建 edge 的 repo 能力
- `serviceRelation` 中允许调用其他 `Repo` 做查询、批量回填与关系补全，但这种 repo-to-repo 调用不得演化为业务编排、状态流转或事务主导

示例
```go
func (r *accountRepo) serviceRelation(ctx context.Context, data interface{}, opts ...filter.Option) error {
	cfg := filter.NewConfig(opts...)
	if _, ok := cfg.Relations[openenum.AccountCheckUserRelation]; !ok || r.adminUserRepo == nil {
		return nil
	}

	list := helper.SliceNormalize[*openbiz.Account](data)
	idList := collectCheckUserIDs(list)
	userMap, err := r.adminUserRepo.MapAdminUser(ctx, &adminbiz.AdminUserFilter{IDList: idList})
	if err != nil {
		return err
	}

	for _, item := range list {
		fillCheckUsers(item, userMap)
	}

	return nil
}
```

区分标准：

- 有稳定结构关系、适合进入 ent schema 的，用 `queryConfig + queryRelation + edges`
- 没有稳定 edge、但需要查询后补全 relation 的，用 `serviceRelation`

### `kit`

#### 定位

`kit` 是为 `repo/data` 提供基础设施能力支撑的稳定入口，不是业务能力。

#### 负责什么

- 基础设施聚合
- 可注入的复用能力
- 生命周期与关闭收口
- 非业务语义的稳定抽象

#### 仓库事实

基于当前仓库实现，`internal/data/kit` 已经统一聚合以下能力：

- `Db`
- `Redis`
- `RedisLock`
- `Snowflake`
- `Producer`
- `EventBus`
- `HttpClient`
- `ProviderSet`

同时，`internal/biz/kit` 承载事务等稳定抽象接口。

#### 不负责什么

- 不承载业务语义
- 不承载业务编排
- 不承载领域状态流转
- 不为某个单一 `repo` 定制私有逻辑

#### 示例

正例：

```go
type flow_pageRepo struct {
	data   *kit.Data
	sumsub *conf.Sumsub
}
```

```go
func Newflow_pageRepo(data *kit.Data, sumsub *conf.Sumsub) openbiz.flow_pageRepo {
	return &flow_pageRepo{
		data:   data,
		sumsub: sumsub,
	}
}
```

反例：

```go
func Newflow_pageRepo(data *kit.Data, dataConf *conf.Data, sumsub *conf.Sumsub) openbiz.flow_pageRepo {
	return &flow_pageRepo{
		data: data,
		// This duplicates infrastructure assembly that should be owned by kit.
		client: &http.Client{
			Timeout: dataConf.HttpClient.Timeout.AsDuration(),
		},
		sumsub: sumsub,
	}
}
```

### `data`

#### 定位

`data` 是 `repo` 的落地实现层。

#### 负责什么

- DB 访问
- ent 集成
- 远程依赖访问
- 查询落地
- relation 装配实现

#### 不负责什么

- 不反向拥有业务主流程
- 不主导事务边界
- 不重复初始化应由 `kit` 统一管理的基础设施

#### 示例

```go
type accountRepo struct {
	data *kit.Data
}

func (r *accountRepo) PageListAccount(ctx context.Context, f *openbiz.AccountFilter, opts ...filter.Option) ([]*openbiz.Account, int, error) {
	query := r.data.Db.Account(ctx).Query()
	query = r.parseFilter(query, f)
	query = r.queryConfig(query, f, opts...)

	count, err := query.Clone().Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	list, err := query.All(ctx)
	if err != nil {
		return nil, 0, err
	}

	result := make([]*openbiz.Account, 0, len(list))
	for _, info := range list {
		result = append(result, r.queryRelation(accountConvert(info), info.Edges))
	}

	if err := r.serviceRelation(ctx, result, opts...); err != nil {
		return nil, 0, err
	}
	return result, count, nil
}
```

## 横切关注

下面这些关注点不是新的层，而是围绕 `usecase / repo / kit / data` 四个角色展开的横切约束。

### 事务边界

事务边界优先放在 `usecase`。

原因：

- 事务属于业务动作边界
- 事务是否需要开启，取决于业务编排，而不是单一数据访问
- `service` 离业务规则太远
- `repo` 只负责数据边界，不应主导业务事务

因此：

- 不将事务边界上浮到 `service`
- 不将事务边界下沉到 `repo`
- 若本地同步 `event bus` 的 `Publish` 位于事务内，则其 listener 执行结果默认跟随当前事务语义
- 若把事务内的 `Publish` 移到事务外，应直接视为事务边界变更，必须重新评估失败时是否仍应触发整体回滚

### relation 收口

relation 应统一收口在 `repo`，由 `usecase` 声明需求，由 `repo` 负责实现。

其中：

- `usecase` 负责表达“需要哪些 relation”
- `repo` 负责查询、装配与补全 relation
- `service` 不补查 relation
- `usecase` 不手工补 relation

### 第三方依赖与 `repo/kit/data` 边界

当 `repo` 需要访问第三方 HTTP / SDK 时，默认先判断它是否天然属于某个 `repo` 的领域依赖。

通常规则如下：

- 若第三方调用只服务单一 `repo`，且没有独立鉴权、重试、限流、错误映射或复用价值，优先直接在 `repo` 中使用已注入的 `data/http client/config`
- 若外部能力本身已经形成稳定边界对象，且被多个 `repo/usecase` 复用，或承载独立策略，则可独立建模
- 若第三方已经形成稳定围绕同一对象族操作的边界对象，即使本地不落库，也允许在 `biz` 中独立定义对应主题文件，收敛其 `Repo` 接口与请求/响应对象
- 若问题本质是“现有基础设施能力缺口”，优先补齐 `kit`，不要在 `repo` 内重复造一套相同能力

不要为了“结构看起来更分层”而额外拆一个只做透传的薄包装 `client`。

### 输入与输出对象

应用层输入输出对象表达业务编排所需的输入与输出。

它们与 `proto message` 的关系是：

- 共享同一领域语义
- 不是同一个对象
- 默认不直接互相替代

如果输入输出对象、领域实体、审核结果对象或事件对象中的字段承载稳定值域，应在进入 `biz` 可消费对象时完成枚举化，而不是继续保留第三方原始字符串参与业务判断。

查询场景中：

- 业务 filter 由 `service` 或 handler 根据协议输入构建
- `opts` 由 `service` 或 handler 显式构建并传递
- `usecase` 接收业务 filter
- `usecase` 接收 `opts`
- `repo` 使用业务 filter 组织查询

不要在 `usecase` 内根据 `proto/request` 原始字段临时拼装查询 filter，也不要在 `usecase` 内临时组装 relation `opts`。

### 领域规则

- `service` 只做协议适配，不吸收业务编排
- `usecase` 只做业务编排、事务边界、权限和状态流转
- `repo` 只做查询组织、relation 装配和依赖访问边界
- `kit` 只做基础设施能力支撑，不吸收业务语义
- `data` 只做 repo 落地实现与依赖访问
- `relation` 统一收口在 `repo`
- 依赖方向保持 `service -> usecase -> repo -> kit/data/depend`

### 判断提示

判断某段逻辑应放在哪一层时，可优先观察：

- 它是在做协议适配、业务编排、依赖抽象，还是基础设施能力支撑
- 它是否需要事务边界或状态流转
- 它是在声明 relation 需求，还是在装配 relation
- 它是在组织查询条件，还是在访问 DB / 远程依赖
- 它是否已经属于通用基础设施能力，而不是单一 `repo` 私有逻辑

## 边界提示

如果问题继续细化，应转到更具体的知识文档：

- 聚合根识别 -> `aggregate-root.md`
- 围绕聚合根的命名收敛 -> `naming.md`
- service / proto 结构与协议边界 -> `service.md`
- ent、wire、listener、consumer、crontab -> `components.md`
- 共享枚举、错误语义与稳定字面量 -> `shared-conventions.md`
