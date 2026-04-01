# Code Style

## 作用范围

本文用于说明 Kratos 项目中的实现约束与反防御式编程规则。
当问题属于以下场景时，应优先查看本文：
- 判断某种写法是否违背当前项目默认实现约束
- 判断代码中是否存在冗余防御、伪成功返回、吞错、重复 helper 或坏味道
- 判断内部输入、注入依赖、错误路径、DTO 写法是否符合项目风格
- 判断为了测试或“看起来更稳”而引入的额外抽象是否合理

本文重点回答实现红线与坏味道治理，不展开聚合边界、分层职责、组件接线与共享语义收敛。

## 规则

### 实现约束

- 上游已经保证的条件，下游不再重复兜底
- 只有当前层真的在接收外部不可信输入时，才保留校验和兜底
- 对构造函数注入的依赖，不写额外 `nil` 判空
- `wire_gen.go` 与 `NewXxx(...)` 构造签名共同构成装配事实源；只要某依赖通过构造参数传入并进入 Wire 链路，就默认属于必选依赖
- 只有显式可选依赖，才允许判空
- 必选注入依赖的缺失属于构造期 / 装配期问题，必须在 `NewXxx(...)` 或 provider 阶段暴露，不得下沉到 `biz / data / service` 方法内做运行期 `nil` 兜底
- 禁止在业务方法、repo 方法、client 方法内针对 Wire 注入字段写这类分支：`if r.dep == nil { ... }`、`if u.repo == nil { ... }`
- 禁止以“可能忘记生成 Wire / wire_gen 未更新”为理由保留运行期判空；这属于代码未完成或装配错误，不是合法业务路径
- 如果依赖是可选能力，必须在构造结构、命名或对象表达上显式说明它是“可选”，而不是将其伪装成必选注入后再在运行期判空
- 失败路径必须真实表达，不得通过 `nil, nil`、吞错、静默成功或只记日志不返回来伪装成功
- 禁止 `nil, nil`
- 禁止 `_ = err`
- 禁止只打日志不返回错误
- 降级、忽略、兜底必须有明确业务契约支持

### 反防御式编程

- 为了“看起来更稳”而增加的防御式代码、重复 helper 或局部例外，必须先证明其确有必要
- 同类问题优先遵循仓库既有实现约束，不要为局部便利重发明一套写法
- 不要为了测试方便，把 `time.Now`、随机数、`UUID` 这类基础函数包装成结构体字段
- 不要为了“看起来更稳”补重复默认值
- 不要对内部输入重复 `TrimSpace`、重复大小写转换或重复规范化
- 错误如果会影响调用结果，就不应被吞掉
- 不要因为单次调用、局部逻辑或“看起来统一”就额外抽一层 helper 或 wrapper
- 无语义增量的判空、兜底、格式化不应继续包装成复用 helper
- `struct`、`interface`、typed const、具名响应对象等稳定声明，默认放在文件顶部
- 仅在声明明显只服务于一个极短局部逻辑，且提到顶部会降低可读性时，才允许保留在临近位置
- 不要把会被多个函数复用的响应结构、请求结构、转换结构放到文件尾部

## 第一部分：实现约束

### 上游契约可信

上游 `usecase`、`service`、`proto` 校验、状态机和主流程分支已经保证的前置条件，下游不要再次做默认值回填。

正例：
```go
levelName := info.LevelName
```

反例：
```go
if levelName == "" {
	levelName = r.config.GetLevelName()
}
```

### 注入依赖可信

通过构造函数或 `Wire` 注入的依赖，默认视为必然存在。

正例：
```go
func NewSumsubRepo(data *kit.Data, tp *conf.ThirdParty) openbiz.SumsubRepo {
	return &sumsubRepo{
		data:   data,
		config: tp.GetSumsub(),
	}
}
```

正例：
```go
func NewKycRepo(data *kit.Data, tp *conf.ThirdParty) (*kycRepo, error) {
	sumsub := tp.GetSumsub()
	return &kycRepo{data: data, sumsub: sumsub}, nil
}
```

反例：
```go
func (r *sumsubRepo) CreateApplicant(ctx context.Context, req *CreateApplicantRequest) error {
	if r == nil || r.data == nil || r.config == nil {
		return nil
	}
	return nil
}
```

反例：
```go
func (r *kycRepo) sumsubRequest(ctx context.Context) error {
	if r.sumsub == nil {
		return fmt.Errorf("sumsub config is required")
	}
	return nil
}
```

### 错误路径必须真实

返回签名里有 `error` 时，失败默认返回真实错误。

正例：
```go
func (r *accountRepo) GetAccount(ctx context.Context, id uint32) (*biz.Account, error) {
	info, err := r.query(ctx, id)
	if err != nil {
		return nil, err
	}
	return info, nil
}
```

反例：
```go
func (r *accountRepo) GetAccount(ctx context.Context, id uint32) (*biz.Account, error) {
	info, err := r.query(ctx, id)
	if err != nil {
		r.log.Errorf("query account failed: %v", err)
		return nil, nil
	}
	return info, nil
}
```

## 第二部分：反防御式编程

### 不要为测试污染生产结构

不要为了测试方便，把 `time.Now`、随机数、UUID 这类基础函数包装成结构体字段。

正例：
```go
nowTime := uint32(time.Now().Unix())
```

反例：
```go
type sumsubRepo struct {
	now func() time.Time
}
```

### 不要重复默认值兜底

不要为了“看起来更稳”补重复默认值。

正例：
```go
status := account.Status
```

反例：
```go
if account.Status == "" {
	account.Status = "init"
}
```

### 不要做无意义格式化

不要对内部输入重复 `TrimSpace`、重复大小写转换或重复规范化。

正例：
```go
levelName := info.LevelName
```

反例：
```go
levelName := strings.TrimSpace(strings.ToLower(info.LevelName))
```

### 不要吞错

错误如果会影响调用结果，就不应被吞掉。

正例：
```go
if err := helper.StructConvert(store.Account, &data); err != nil {
	return err
}
```

反例：
```go
if err := helper.StructConvert(store.Account, &data); err != nil {
	log.Error(err)
}
```

### 不要制造伪抽象

不要因为单次调用、局部逻辑或“看起来统一”就额外抽一层 helper 或 wrapper。

正例：
```go
if account.Status == openenum.AccountStatusInit {
	return account.Submit()
}
```

反例：
```go
func processAccountStatus(account *Account) error {
	if account.Status == openenum.AccountStatusInit {
		return account.Submit()
	}
	return nil
}
```

### 不要把无信息增量包装成 helper

无语义增量的判空、兜底、格式化不应继续包装成复用 helper。

正例：
```go
email := req.Email
```

反例：
```go
func normalizeEmail(email string) string {
	return strings.TrimSpace(strings.ToLower(email))
}
```

### 文件内声明顺序要稳定

文件内的稳定声明应优先集中放在顶部，避免类型定义散落在函数之后，增加阅读和定位成本。

正例：
```go
type sumsubApplicantResponse struct {
	ID string `json:"id"`
}

func (r *accountKycRepo) CreateApplicant(ctx context.Context) error {
	return nil
}
```

反例：
```go
func (r *accountKycRepo) CreateApplicant(ctx context.Context) error {
	return nil
}

type sumsubApplicantResponse struct {
	ID string `json:"id"`
}
```

## 判断提示

判断某段实现是否违反 code style 时，可优先观察：
- 它是在表达真实约束，还是在增加心理安全感
- 它是否重复了上游已保证的前置条件
- 它是否把失败路径伪装成成功
- 它是否为了测试、统一或“看起来优雅”引入了无价值抽象
- 它是否把稳定结构写成匿名或临时形态
- 它是否把稳定声明散落在函数之后，导致文件结构不稳定

## 边界提示

如果问题继续细化，应转到更具体的知识文档：
- 聚合根识别 -> `aggregate-root.md`
- 围绕聚合根的命名收敛 -> `naming.md`
- usecase、repo、data、事务边界 -> `domain.md`
- service / proto 结构与协议边界 -> `service.md`
- ent、listener、consumer、crontab、event、kit -> `components.md`
- `internal/pkg` 公共能力边界 -> `pkg.md`
- 错误语义、枚举、稳定值域 -> `error-enum.md`
