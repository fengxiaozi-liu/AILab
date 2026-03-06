# Comment Reference

## 这个主题解决什么问题

说明 Kratos 项目中的代码注释、Proto 注释和 Ent 注释通常如何组织，避免“无注释”与“无效注释”两种极端。

## 适用场景

- 新增领域对象、UseCase、Repo、Service
- 新增 Ent schema、field、table comment
- 重构代码后同步修正注释

## 设计意图

注释不是为了补救糟糕命名，而是为了表达以下代码本身不容易直接看出的信息：

- 业务角色
- 能力边界
- 生命周期职责
- 状态或字段的业务含义
- 协议或表结构的外部可见语义

稳定做法是：

- service / usecase / repo / provider 保留简洁的职责注释
- Ent schema 保留业务语义注释
- 阶段性长流程使用少量章节注释分段

这类注释通常都很短，一般一行即可说明功能、规则或关键语义。
在语言选择上，注释默认优先中文，只有外部接口名、标准术语或英文表达更稳定时才保留英文。

## 实施提示

- 先把命名稳定下来，再补注释
- 注释优先解释“为什么存在”“负责什么”，不是解释“这一行在赋值”
- 默认写短注释，一般一行足够；只有边界复杂时才补多行
- 注释优先落在对象定义、入口函数、关键分支、流程阶段切换点
- 默认优先中文注释，同一文件内风格尽量统一
- 生成物注释只通过源文件生成，不在生成物手改

## 推荐实现方式

### 1. Service / Provider 注释

```go
// ProviderSet 服务提供集合
var ProviderSet = wire.NewSet(NewAccountService)

// AccountService 实现 admin 侧账户服务接口。
type AccountService struct {
    ...
}

// NewAccountService 创建账户服务。
func NewAccountService(...) *AccountService {
    ...
}
```

### 2. Domain / UseCase 注释

```go
// Account 领域模型
type Account struct {
    ...
}

// AccountUseCase 账户聚合根能力
type AccountUseCase struct {
    ...
}

// Prepare 创建或刷新开户主体。
func (u *AccountUseCase) Prepare(ctx context.Context, userCode string) (*Account, error) {
    ...
}
```

### 3. 关键位置注释

```go
// ----- Commit flow helpers -----
func (u *AccountFlowPageUseCase) buildPageStore(...) *Store {
    ...
}
```

### 4. Ent 注释

```go
func (Account) Fields() []ent.Field {
    return []ent.Field{
        field.Uint32("id").Comment("primary id"),
        field.String("user_code").MaxLen(64).Default("").Comment("user code"),
        field.Uint8("open_status").Default(0).Comment("open status"),
    }
}

func (Account) Annotations() []schema.Annotation {
    return []schema.Annotation{
        schema.Comment("account table"),
    }
}
```

## 代码示例参考

### 工具能力注释

```go
// Parallel 实现了一个并发任务池，使用 goroutine 来执行任务，支持控制并发度、异常处理和任务序列号
type Parallel struct {
    ...
}
```

## Good Example

- `AccountService implements the admin AccountServiceServer interface.`
- `AccountUseCase 账户聚合根能力`
- `// ----- Commit flow helpers -----`
- `schema.Comment("开户页面流程状态表")`
- `field.String("page_code").Comment("页面编码")`

## 常见坑

- 只有导出名，没有职责注释
- Ent schema 没有业务语义注释，只剩技术字段名
- 在 proto 文件中添加任何说明性注释，破坏协议文件整洁性
- 注释写成大段说明书，掩盖真正的代码结构
- 每一小段代码都加注释，导致关键点反而不突出
- 大量写“获取数据”“执行逻辑”这种无效注释
- 重构后代码变了，注释没同步更新

## 相关 Rule

- `../rules/comment-rule.md`

## 相关 Reference

- `../../kratos-components/reference/ent-spec.md`
- `../../kratos-entry/reference/proto-spec.md`
