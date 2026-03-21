---
name: kratos-patterns
description: |
  用于 Kratos 项目的路由判定与子技能选择。适用于任何 Kratos 代码修改、重构、代码评审、生成物更新，或不确定应加载哪个 skill 的场景。触发关键词包括 SERVER_NAME、项目类型判定、BaseService、业务服务、网关、路由、技能选择、实现前审计。
---

# Kratos Patterns

## 必读规则

- `./rules/routing-rule.md`
- `./rules/change-control-rule.md`
- `./rules/audit-output-rule.md`

## 触发时机

涉及代码更新、规范重组、生成物更新或代码评审时加载本技能。

## 识别工作域

通过 `.env.*` 中的 `SERVER_NAME` 判定项目类型：

| 项目类型 | SERVER_NAME 特征 | 工作流 |
|----------|-----------------|--------|
| BaseService | 以 `BaseService` 结尾 | 抽象定义 |
| 网关类 | 包含 `GatewayService` 或 `OpenapiService` | 接口代理 |
| 业务项目 | 不属于以上两类 | 业务实现 |

无法读取 `SERVER_NAME` 时按以下顺序兜底：

1. 若存在明显的 `gateway` / `proxy` 目录或 `GatewayService/OpenapiService` 字样，判定为网关类。
2. 若仓库主要包含 proto、枚举、错误码、依赖封装等抽象定义，判定为 BaseService。
3. 其他情况判定为业务项目。

## 子技能选择

路由完成后，按本次“职责变化”优先选择聚合 skill：

- 触发 `kratos-domain`：聚合边界、layer 分层、命名、UseCase、Repo、业务测试
- 触发 `kratos-components`：Ent、EventBus、Crontab、Depend、Config
- 触发 `kratos-pkg`：internal/pkg 下的 context、metadata、middleware、proto helper、schema、seata、util
- 触发 `kratos-conventions`：Error、Enum、i18n
- 触发 `kratos-entry`：Server、Gateway、Codegen、Proto/Wire/Build

参考文档说明：

- 以上 5 个聚合 skill 已自带各自的 `reference/` 目录。
- 查阅规范时，优先进入聚合 skill 自身目录，不再回退到旧 skill。

## 路由原则

1. 路径只能作为证据提示，不是唯一触发条件；优先看本次改动的职责和语义。
2. 每次编码必须按“路由 -> 聚合 skill -> 实施”的顺序执行。
3. 一个任务可同时加载多个聚合 skill，但应保持最小集合。

## 强制输出

开始前输出一段“路由与加载清单”：

- `WorkDomain:` BaseService / 业务 / 网关
- `Evidence:` 用于判定的证据
- `SubSkills:` 本次将加载的聚合 skill 列表
- `Rules:` 本次默认需要加载的规则文件列表

提交前输出“自检摘要”：

- 每个已加载聚合 skill 的检查点是否通过（Yes/No）
- 运行了哪些 codegen / build 验证
