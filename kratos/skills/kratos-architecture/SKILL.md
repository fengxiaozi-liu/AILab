---
name: kratos-architecture
description: 用于 Go Kratos 项目的架构判断与知识路由。凡是需要判断聚合根、分层归属、命名边界、proto/service/server/gateway 落位、ent/wire/provider/listener 归属、变量与常量定义、组件注入、kit/data/repo 边界、internal/pkg 下沉与 shared conventions 时必须使用。
---

# Kratos 知识事实依据

## 职责

`kratos-architecture` 是当前 Kratos 项目结构性改动的架构判断与规范裁决来源。
它基于仓库事实和既有约束，为本次改动给出可执行判断。其职责包括：

- 判断当前改动围绕哪个 aggregate root 组织
- 判断代码应落在哪一层：`service / biz / data / listener / consumer / crontab / server / pkg`
- 判断 `proto / service / server / gateway / ent / wire / provider / listener` 等协议与组件应如何落位
- 判断 `repo / kit / data` 的职责边界与依赖方向
- 判断哪些能力应保留在当前业务域内，哪些能力可进入 `internal/pkg`
- 判断稳定语义应如何收敛，包括变量、常量、`error`、`enum`、`i18n`、`logging`、`comment`
- 判断实现是否符合仓库既有编码约束，是否破坏注入前提、偏离既有工程模式

## 证据优先级

所有判断都必须遵循以下优先级：

1. 当前仓库事实
2. 当前仓库已有实现模式
3. 本 skill references
4. 通用 Kratos 经验

禁止跳过仓库事实，直接套用外部经验。

## 使用方式

1. 先读取与当前任务直接相关的仓库文件。
2. 先识别本次要判断的问题类型，再选择对应的主要 reference。
3. 仅当主要 reference 不足以支持判断时，再补充读取相关 reference。

## References 路由说明

本 skill 的 references 按判断主题组织，用于把具体判断问题路由到对应的项目知识。

| 判断主题 | 主要 reference | 何时读取 | 产出目标 |
| --- | --- | --- | --- |
| 项目类型、目录职责、层级边界 | `references/layer.md` | 判断改动落层、目录职责、项目类型时 | 明确层级归属与目录边界 |
| 聚合根识别 | `references/aggregate-root.md` | 判断当前改动围绕哪个领域对象组织时 | 明确 aggregate root |
| 领域、usecase、repo、kit、data 边界 | `references/domain.md` | 判断 biz / repo / kit / data / usecase 职责与依赖方向时 | 明确领域边界与职责分工 |
| 命名与术语 | `references/naming.md` | 判断目录名、类型名、方法名、术语一致性时 | 明确命名归属 |
| proto、service、server、gateway | `references/service.md` | 判断协议定义、服务接口、服务暴露层时 | 明确协议与服务组件落位 |
| OpenAPI v3 文档注解规范 | `references/openapi-v3-spec.md` | 判断 OpenAPI 注解和文档输出时 | 明确注解与文档约束 |
| ent、listener、consumer、crontab、event、wire、provider | `references/components.md` | 判断运行时组件、数据组件、注入组件职责时 | 明确组件归属 |
| `internal/pkg` 下沉判断 | `references/pkg.md` | 判断能力是否可共享下沉时 | 明确是否允许进入 `internal/pkg` |
| error、enum、稳定值域、常量定义 | `references/error-enum.md` | 判断错误码、枚举、协议常量、稳定字面量时 | 明确共享语义落位 |
| i18n、logging、comment | `references/shared-conventions.md` | 判断共享约定时 | 明确 shared conventions |
| code style、反防御式编程、组件注入、是否需要判空 | `references/code-style.md` | 判断 nil-check、注入假设、编码约束时 | 明确实现约束 |

## 硬约束

- 不要把业务主流程塞进 `service / listener / server / consumer / crontab`
- 不要把带业务语义的能力伪装成 `internal/pkg` 公共能力
- 不要脱离 aggregate root 并行发明命名体系
- 不要在没有证据的情况下给出层级判断
- 不要只给结论，不给基于仓库事实的理由

## 补充约束

- 评审第三方 `client / repo / adapter` 代码时，先判断“这是 DTO 还是稳定协议字面量”，再判断落位
- 请求体、响应体、payload 结构体等 DTO 可以留在 `internal/data/...`
- header、path 模板、query key、content-type、第三方固定状态值、固定结果值等稳定协议字面量，必须收敛到 `internal/enum/<domain>`
- 禁止用“当前只在一个文件里使用”作为把稳定协议字面量留在实现文件中的理由
- 判断基础设施能力归属时，优先看仓库是否已存在 `kit` 聚合入口；已进入 `kit` 的能力默认通过注入复用，不在 `repo` 内重复构造
