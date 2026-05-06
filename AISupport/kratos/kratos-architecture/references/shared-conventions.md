# Shared Conventions

## 作用范围

本文用于说明 Kratos 项目中的横切共享约定，包括 `i18n` key、`logging` 字段、`comment` 规范。
当问题属于以下场景时，应优先查看本文：

- 判断 `i18n` key 与多语言文案应如何组织
- 判断 `logging` 字段与日志表达应如何统一
- 判断 `comment` 应如何表达职责、边界与语义
- 判断共享表达应落在 `i18n`、`logging` 还是 `comment`

本文重点回答共享表达如何统一，不展开业务编排、组件接线、协议层结构设计与代码红线。

## 规则

- `i18n / logging / comment` 只解决跨模块共享表达统一，不承担业务编排或领域建模职责
- 同一类共享表达应使用统一口径、统一字段和统一语言风格，不为局部实现破例
- 注释、日志与多语言文案应真实表达职责与语义，不得掩盖失败路径或边界漂移
- 共享约定一旦已在仓库中形成稳定写法，新实现应优先复用，不重发明局部风格

## 共享约定定位

共享约定解决的是“多个模块应如何说同一种话”。
它关注的是：

- `i18n` 表达统一
- `logging` 口径统一
- 注释表达统一

它不关注的是：

- 业务流程本身
- 组件接线
- 聚合建模
- 代码风格红线

## `i18n`

定位：稳定语义的多语言 key 与文案模板。
承载内容：

- `i18n` key
- 文案模板
- 多语言资源
- `localize` 调用口径

边界提示：

- `i18n` key 与多语言资源默认落在项目既有 `i18n` 资源目录
- `i18n` key 应表达稳定语义，不直接使用页面型临时 key
- enum / typed const 的 `i18n` `ID` 约束与正确示例，统一回到 `error-enum.md`
- 默认文案放在 `DefaultMessage`，不散落在业务流程
- 修改 `i18n` 资源时，走生成或资源同步链路，不手改派生产物

示例：

```toml
[OPEN_ACCOUNT_SUBMIT_SUCCESS]
description = "开户提交成功"
other = "开户提交成功"
```

## `logging`

定位：统一日志字段与日志表达。
承载内容：

- `log.Logger`
- `log.Helper`
- 可选的 `module` 字段
- 上下文日志
- 错误日志表达

边界提示：

- logging 配置与模块日志默认沿项目既有日志初始化与 helper 构造方式收敛
- 默认直接使用 `log.NewHelper(logger)`，不要为了统一模板强行追加 `module`
- `log.Logger` 作为框架注入依赖可信传递，组件内只在构造阶段派生一次 `log.Helper`
- 不在运行时链路重复包装 logger，也不在同一对象内同时维护多套等价 logger/helper
- 只有组件确实需要稳定的模块字段，且仓库内同类组件已形成一致模式时，才使用 `log.With(logger, "module", "...")`
- 使用 `module` 时，命名必须复用稳定组件路径，不围绕临时文件名、实现细节或同义别名发明新值
- 记录错误时带动作和 `err`
- 高频循环控制日志量，敏感信息脱敏或不记录

示例：

```go
func NewAccountUseCase(logger log.Logger) *AccountUseCase {
	return &AccountUseCase{
		log: log.NewHelper(logger),
	}
}
```

## `comment`

定位：职责、边界与语义注释规范。
承载内容：

- 对象级职责注释
- 字段语义注释
- 流程关键转折点注释
- 生成物注释来源约束
- 字段类型旁的注释说明
- 方法内关键注释说明

边界提示：

- `comment` 默认跟随对象定义、字段定义与关键流程转折点落位
- `comment` 补充边界、约束与语义，不重复代码表面行为
- `Service` / `UseCase` / `Repo` / `Provider` 保留简短职责注释
- 字段类型的注释说明放在同一行后面，直接解释字段承载的职责或语义，不单独拆成下一行长注释
- 方法内只在关键转折、关键分支或关键约束处加注释，说明采用简洁短句，直接描述“这里在做什么”
- ent schema / field 保留业务语义注释，统一使用中文
- 生成物注释通过源文件生成，不直接手改生成物
- 不要用 `comment` 补救糟糕命名、糟糕分层或糟糕抽象

示例：

```go
// AccountUseCase 负责编排账户相关业务流程。
type AccountUseCase struct {
	log *log.Helper // 账户用例日志，记录账户业务编排中的关键节点与异常
}
```

```go
field.String("status").Comment("账户当前审核状态")
```

```go
func (uc *AccountUseCase) Submit(ctx context.Context, req *SubmitRequest) error {
	// 先落库，再进入后续提交流程，避免后续步骤缺少账户主记录。
	if err := uc.repo.Create(ctx, req); err != nil {
		return err
	}

	// 审核流只在开户记录创建成功后触发。
	return uc.workflow.Start(ctx, req)
}
```

## 判断提示

判断一个表达是否应纳入共享约定时，可优先观察：

- 它是否会跨函数、跨文件或跨层复用
- 它是否承载稳定共享语义
- 它是否已经在多个位置出现并开始漂移
- 它应落在 `i18n`、`logging` 还是 `comment`

## 边界提示

如果问题继续细化，应转到更具体的知识文档：

- 聚合根识别 -> `aggregate-root.md`
- 围绕聚合根的命名收敛 -> `naming.md`
- 错误语义、枚举、稳定值域 -> `error-enum.md`
- service / proto 结构与协议边界 -> `service.md`
- usecase、repo、data、事务边界 -> `domain.md`
- 代码红线与反防御式编程 -> `code-style.md`
