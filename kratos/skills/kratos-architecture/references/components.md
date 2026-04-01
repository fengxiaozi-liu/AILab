# Components

## 作用范围

本文用于说明 Kratos 基础设施组件的接入知识，包括 `ent`、`consumer`、`listener`、`crontab`、`event`、`provider`、`wire` 与生命周期装配。

当问题属于以下场景时，应优先查看本文：

- 新增或修改基础设施组件
- 判断组件应落在哪个目录
- 判断组件如何接线、注册与进入生命周期
- 判断 `ent`、`listener`、`consumer`、`crontab`、`event` 的职责边界
- 判断 provider / `wire` / 注册入口应如何收敛

## 规则

- 组件层只负责接入、注册、生命周期管理与最小适配，不承担业务编排、事务裁决或领域主线组织
- `listener / consumer / crontab / server / event bus` 等运行时组件不得演化成隐藏的 `UseCase`
- `ent / wire / provider / event / listener` 的变更应同时评估生成链路、注册链路与运行时影响
- event bus 默认按进程内同步钩子理解；迁移其事务内外位置时，应直接按事务语义变化处理
- MySQL 下统一使用 `OnConflict()` 组织冲突更新，不使用 `OnConflictColumns(...)`

## 组件定位

组件层负责：

- 基础设施接入
- 注册
- 生命周期管理
- 最小适配

组件层不负责：

- 业务流程编排
- 事务边界决策
- 领域主线组织
- 通用公共能力下沉

组件的重点不是“多写一层”，而是让接线、注册、生命周期与联动关系清晰稳定。

## 组件类型

Kratos 基础设施组件通常可归为以下几类：

- `ent`
  - `ent/schema`、字段、edge、index、ORM 生成链路
- `consumer`
  - MQ consumer、topic / queue、消费回调、ack / retry、注册入口
- `listener`
  - 本地 listener、handler 聚合注册、事件回调
- `crontab`
  - scheduler、定时任务注册、启停、幂等与并发
- `event`
  - local event / eventbus、topic、payload、发布订阅关系

如果一个改动同时命中多个组件类型，优先看“注册归属与生命周期归属最明确的主改动点”。

## 组件规则

围绕组件接入与装配时，应优先遵守以下规则：

- 组件层只做接入、注册、生命周期管理和最小适配
- 业务编排、事务边界、状态机判断留在 `usecase`
- 组件代码优先沿用项目已有目录结构、provider 聚合和生命周期收口方式
- 新增组件前，先看项目是否已有同类组件实现
- provider、`wire`、注册器、启动入口、关闭入口应尽量收敛，不散落在多个模块

## `ent`

定位：ORM schema 与数据生成组件。

承载内容：

- `internal/data/ent/schema`
- 字段、edge、index
- ORM 生成链路
- schema 注释与落位

边界提示：

- `ent` 的事实源是 schema，不是生成物
- 修改 schema、field、edge、index 时，回到 `internal/data/ent/schema`
- 不手改 `ent` 生成代码
- schema、field、edge、index 的 `Comment` 统一使用中文简体，直接表达业务语义
- MySQL 下统一使用 `OnConflict()` 组织冲突更新，不使用 `OnConflictColumns(...)`
- 只有具备唯一约束或明确幂等键的字段，才能作为 upsert 冲突依据

示例：

```text
internal/data/ent/schema/account.go
go generate ./internal/data/ent
```

## `consumer`

定位：消息消费入口组件。

承载内容：

- payload 解析
- 最小入参组装
- 调用 `usecase`
- topic / queue 接入
- 注册与生命周期接入

边界提示：

- `consumer` 负责消费入口，不编排完整业务流程
- 解析失败应返回错误，不吞错
- topic 与 queue 语义保持清晰

## `listener`

定位：本地事件监听组件。

承载内容：

- listener 注册
- 本地事件回调
- 最小转发逻辑

边界提示：

- `listener` 负责事件响应入口，不承担 repo 编排
- 注册应收敛在统一 listener 聚合入口
- 错误直接返回，不吞错

## `crontab`

定位：定时任务入口组件。

承载内容：

- cron 注册
- 调度入口
- 生命周期启停
- 最小任务触发闭环

边界提示：

- Job 足够薄，只做触发与错误返回
- 复杂业务流程回到 `usecase`
- 幂等策略与并发策略先于实现

## `event`

定位：本地事件与事件总线组件知识。

承载内容：

- topic / event 值域
- payload 选择
- 发布与订阅边界
- eventbus 基础设施落位

边界提示：

- 稳定 topic / key 集中定义在 enum 或常量
- payload 优先复用稳定对象，必要时单独定义 payload
- 发布方只发布事实，不为订阅方额外拼装
- 事件不替代事务边界
- 当前仓库中的 event bus 默认视为进程内同步分发机制，不要按异步 MQ / eventual consistency 理解
- event bus 更接近钩子函数或本地回调分发，用于抽离部分后置逻辑，不天然形成新的异步边界
- 若 `Publish` 发生在事务闭包内，则 listener 的执行与失败传播属于当前事务语义的一部分
- 若把原本事务内的 `Publish` 迁到事务外，必须显式评估是否改变回滚语义、失败传播路径、状态可见性与执行顺序
- 不要仅因“看起来更像事件驱动”就把本地同步 event bus 当成异步解耦手段；先判断是否会改变原有业务一致性

## provider、`wire` 与生命周期

定位：组件接入的统一装配知识。

承载内容：

- provider 聚合
- `ProviderSet`
- 注册器
- 生命周期 Run / Stop
- 启动与关闭入口

边界提示：

- 新增组件不仅要写组件本体，还要补齐注册、provider、生命周期与最小联动闭环
- 组件自己的生命周期不单独散落管理
- 修改 schema、provider、`wire`、注册入口时，要同步考虑生成、构建或启动链路

## 判断提示

判断某个改动是否属于组件问题时，可优先观察：

- 当前问题是在做基础设施接入，还是在做业务编排
- 当前问题是否涉及注册、provider、`wire`、生命周期
- 当前组件是否已有稳定目录与聚合入口
- 当前逻辑是最小适配，还是已经开始承载业务主线
- 当前联动是否已经补齐生成链路、构建链路或启动链路

## 边界提示

如果问题继续细化，应转到更具体的知识文档：

- usecase、repo、data、`kit`、事务边界 -> `domain.md`
- service / proto 结构与协议边界 -> `service.md`
- `internal/pkg` 公共能力边界 -> `pkg.md`
- 共享枚举、错误语义与稳定字面量 -> `shared-conventions.md`
