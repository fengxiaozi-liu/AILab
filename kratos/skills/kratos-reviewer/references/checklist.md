# Reviewer Checklist

## 作用

这份 checklist 只负责防漏。

它只回答三件事：

- 本次改动触达了哪些审查维度
- 每个已触发维度至少不能漏看什么
- 每个涉及文件都要顺手扫哪些默认问题

## 使用原则

- 先做变化画像，再展开 checklist
- checklist 用于防漏，不代替 facts 与 `规则` 的回查
- 已触发的维度必须完整覆盖，未触发的维度可以不展开
- 每条 finding 都应能回收到某个审查维度

## 默认文件级扫描

每个涉及文件都顺手检查：

- 是否存在 `nil, nil`、吞错、静默成功、空结果无错误
- 是否存在已脱离当前业务语义、具备稳定复用价值、应考虑下沉到 `internal/pkg` 的函数或逻辑
- 构造注入与 `Wire` 已保证的前提，构造注入依赖是否被运行期判空
- 是否存在用历史写法为新偏离背书的情况

## 审查维度 1：aggregate-root

至少不要漏看：

- 是否把同一主题拆成多个平级 `UseCase` 或平级 owner
- 是否按页面、阶段、回调方向、结果类型重新发明主题入口
- 是否仍然能稳定映射回同一 aggregate root

常见信号：

- `XxxResultUseCase`
- `XxxCallbackUseCase`
- `XxxRetryUseCase`
- `XxxSyncUseCase`
- 同一主题方法散落到多个类型、文件或目录

## 审查维度 2：layer

至少不要漏看：

- 是否破坏 `service -> usecase -> repo -> data/kit/depend`
- 入口层组件是否开始承担业务编排、状态推进、补偿流程
- `biz / repo / data` 是否出现职责漂移

常见信号：

- `listener / service / consumer / server` 中出现业务分支或状态推进
- `repo / data` 中出现明显领域判断
- `biz` 直接感知协议、持久化或第三方 SDK 细节

## 审查维度 3：domain

至少不要漏看：

- 同一稳定主题是否只保留一个 `UseCase owner`
- 事务边界是否仍留在 `usecase`
- relation 是否仍统一收口在 `repo`
- 第三方依赖是否落在正确的 `repo / kit / data` 边界

常见信号：

- `usecase` 内部再拆 `Result / Callback / Retry / Sync`
- `repo` 或 `data` 中出现业务编排
- 事务从 `usecase` 漂移到其它层

## 审查维度 4：naming

至少不要漏看：

- 命名是否仍围绕稳定业务主题收敛
- 是否用技术实现名掩盖领域边界漂移
- 是否出现同一概念的多套近义命名

常见信号：

- `internal/data/<domain>` 中出现 `XxxClient / XxxProvider / XxxService`
- 同时出现 `Repo / Client / Provider / Adapter`
- 新命名明显围绕实现手段，而不是围绕主题

## 审查维度 5：service

至少不要漏看：

- `proto -> service -> server` 主线是否仍成立
- `service / server / gateway` 是否承担了业务编排
- 协议映射、codegen、provider/wire 装配是否保持闭环

常见信号：

- `service` 中出现明显状态推进或业务规则
- `proto` 与实际协议/文档不一致
- 接入层改动未补齐对应生成链路

## 审查维度 6：components

至少不要漏看：

- `listener / consumer / service / crontab / server` 是否吞错或伪装成功
- 运行时组件是否演化成隐藏的 `UseCase`
- event bus 时机变化是否改变事务或回滚语义
- `wire / provider / lifecycle` 变更是否影响装配与运行

常见信号：

- `_, _ = ...`
- 只记日志不返回
- 入口组件中出现大量业务分支、补偿或状态推进

## 审查维度 7：error-enum

至少不要漏看：

- 稳定值域是否已经收敛成 enum、typed const 或共享定义
- 第三方状态值是否直接进入 `biz` 参与分支或状态机
- 稳定协议字面量是否散落在实现中
- 错误语义是否真实表达，没有把失败伪装成成功

常见信号：

- `if x == "A"`、`switch status { ... }`
- 裸字符串或裸数字直接参与业务判断
- `not found`、状态错误被静默吞掉

## 审查维度 8：pkg

至少不要漏看：

- 是否把业务语义伪装成公共能力
- 是否存在应下沉却仍散落在局部实现中的稳定公共逻辑
- 是否为了“以后可能复用”过早下沉

常见信号：

- `internal/pkg` 中出现明确业务术语或业务错误
- 多个文件重复出现同类 helper 或转换逻辑
- 以“当前只在一处使用”为理由拒绝抽取稳定公共能力

## 审查维度 9：code-style

至少不要漏看：

- 失败路径是否真实表达
- 构造注入依赖是否被运行期判空兜底
- 是否为了局部方便破坏全局一致性

常见信号：

- 返回签名含 `error`，实现却吞错
- `nil, nil`
- `if u.repo == nil { ... }`

## 审查维度 10：shared-conventions

至少不要漏看：

- 共享表达是否仍保持统一口径
- 注释、日志与共享文案是否真实表达职责与边界
- 是否为局部实现重新发明共享风格

常见信号：

- 同类日志字段、注释风格明显不一致
- 用注释或日志掩盖失败路径和边界漂移
- 已有共享约定存在，但新实现另起一套
