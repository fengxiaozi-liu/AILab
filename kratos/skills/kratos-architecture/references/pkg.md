# Pkg

## 作用范围

本文用于说明 `internal/pkg` 的公共能力边界、子域划分、复用原则与下沉判断。
当问题属于以下场景时，应优先查看本文：

- 判断某段 helper 是否应进入 `internal/pkg`
- 判断能力应落在哪个 `pkg` 子域
- 判断应复用现有公共能力，还是新增子域
- 判断某段逻辑是公共技术能力，还是业务语义实现
- 判断公共函数是否已经命中下沉条件

本文重点回答公共能力如何稳定下沉，不展开业务编排、组件接入、领域建模与代码红线。

## `internal/pkg` 定位

`internal/pkg` 承载的是脱离具体业务语义后仍成立，并且会被多个模块稳定复用的公共技术能力。
它处理的是：

- 通用技术能力
- 跨业务复用能力
- 稳定 helper
- 技术侧上下文、协议辅助与链路工具

它不处理的是：

- 业务流程本身
- 单一业务场景抽象
- 组件接入与生命周期注册
- 带强业务语义的协议拼装

## 什么叫公共能力

一个能力适合进入 `internal/pkg`，通常同时满足以下条件：

- 脱离具体业务语义后仍然成立
- 预期会被多个模块稳定复用
- 解决的是一个稳定能力，而不是某个业务流程局部步骤
- 能用目录名和函数名直接表达边界

如果一段代码离开当前业务场景就失去意义，它通常就不是 `internal/pkg` 级别的公共能力。

## 公共能力判断规则

判断是否应下沉到 `internal/pkg` 时，应优先遵守以下规则：

- 先检索 `internal/pkg` 是否已有可复用实现
- 能复用就复用，不能复用再判断是否值得新增子域
- 每个公共函数只解决一个稳定能力，不承载业务流程编排
- 若一段代码同时包含“通用算法 / 技术计算逻辑”和“协议拼装 / 业务边界逻辑”，只把前者视为公共能力
- 目录名应直接表达边界，避免 `common`、`helper`、`misc` 这类模糊名字
- 遇到重复出现的无状态技术片段时，要显式判断“是否已经命中公共函数下沉条件”，不能只因为当前文件还能放下就继续内联

例如：

- `hmac`、`sha256`、`hex`
- 通用 URL 处理
- 稳定字符串 / 字节算法

更适合进入 `internal/pkg`。

而以下内容不应伪装成公共能力下沉：

- 第三方路径拼装
- 协议字段选择
- 业务状态判断
- 单一业务流程中的临时 helper

## 已有子域

当前仓库中的 `internal/pkg` 已有以下稳定子域：

- `util`
- `context`
- `metadata`
- `middleware`
- `proto`
- `schema`
- `seata`
- `filter`
- `localize`

这些子域已经表达出当前仓库的公共能力边界。新增能力前，优先判断是否应复用现有子域。

## 子域说明

### `util`

定位：无状态通用函数子域。
承载内容：

- 字符串处理
- URL 处理
- 编码与解码
- 并发与小型算法能力

边界提示：只放无业务语义、无状态、可稳定复用的小型通用能力。

示例：

```go
func JoinBaseURL(baseURL, uri string) string {
	return strings.TrimRight(baseURL, "/") + uri
}
```

### `context`

定位：请求级上下文 helper 子域。
承载内容：

- typed `WithXxx`
- typed `GetXxx`
- 进程内上下文辅助

边界提示：不把聚合对象、大型关系对象或临时结果塞进 context。

### `metadata`

定位：metadata 读写与透传子域。
承载内容：

- metadata key
- `GetXxx/SetXxx` helper
- 来源透传能力

边界提示：同一语义只保留一套稳定 key，不把 key 散落在 `service` / `repo` / `middleware` 中硬编码。

### `middleware`

定位：通用链路处理中间件子域。
承载内容：

- recovery
- error format
- metadata 注入
- tracing 链路聚合

边界提示：`middleware` 只处理通用链路问题，不查 `repo`，不做业务判断。

### `proto`

定位：协议辅助转换子域。
承载内容：

- paging 转换
- sort / filter / time range 辅助转换
- 稳定协议辅助工具

边界提示：不把单一业务 reply 拼装下沉到 `internal/pkg/proto`。

### `schema`

定位：结构提取与反射辅助子域。
承载内容：

- 通用结构提取
- 反射辅助
- 字段检查工具

边界提示：不放业务校验或业务配置判断。

### `seata`

定位：事务辅助子域。
承载内容：

- 通用事务封装
- seata 辅助能力

边界提示：事务边界仍由 `usecase` 决定，`seata` helper 不承载业务补偿或状态判断。

### `filter`

定位：稳定查询过滤与 relation 配置子域。
承载内容：

- paging
- sort
- compare
- group by
- relation option

边界提示：`filter` 表达稳定查询与 relation 配置能力，不直接拼业务协议对象。

### `localize`

定位：本地化与语言包辅助子域。
承载内容：

- i18n bundle
- localize helper
- 多语言资源访问辅助

边界提示：它负责本地化基础能力，不承载业务文案流程编排。

## 新增子域

新增 `internal/pkg` 子域时，应优先遵守以下规则：

- 现有子域无法准确表达能力边界时，再新增目录
- 新子域应承载一类稳定能力面，而不是多个主题混放
- 优先 small surface area，只暴露少量稳定 API
- 至少用一个真实调用点证明复用价值

不要为了“看起来通用”或“结构更漂亮”就额外新增目录。

## 判断提示

判断一个能力是否适合进入 `internal/pkg` 时，可优先观察：

- 它是否脱离具体业务语义后仍成立
- 它是否会被多个模块稳定复用
- 它更像技术能力，还是某个业务步骤
- 现有 `pkg` 子域是否已经有相同或相邻能力
- 它是否开始混入 `repo`、业务 `context`、领域对象或组件职责

## 边界提示

如果问题继续细化，应转到更具体的知识文档：

- usecase、repo、data、事务边界 -> `domain.md`
- ent、listener、consumer、crontab、event、kit -> `components.md`
- 共享枚举、错误语义与稳定字面量 -> `shared-conventions.md`
- 代码红线与反防御式编程 -> `code-style.md`
