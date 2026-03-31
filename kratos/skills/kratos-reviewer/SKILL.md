---
name: kratos-reviewer
description: 当任务需要检查 Go Kratos 项目的代码、改动或变更是否符合分层、命名、聚合根、Wire、internal/pkg、runtime 组件、data-access 规则以及 kratos-architecture 规范时使用。仅用于代码审查、结构检查和风险识别；不要用于直接生成完整代码或纯语法解释。
---

# Kratos 代码审查

## 作用

用于审查 Kratos 项目中的代码、改动与变更，判断其是否违反 `kratos-architecture` 定义的架构规范与仓库既有实现约束。

它负责检查：

- 项目改动与变更是否符合 `kratos-architecture` 的判断依据与规范要求
- 是否放错层
- 是否命名错误
- 是否违反聚合根驱动规则
- 是否误放 `internal/pkg`
- 是否对 Wire 注入依赖做多余 nil 判断
- 是否将 ent 细节泄漏到 biz
- 是否把复杂业务逻辑放进 `consumer / listener / service`
- 是否存在不合理的 `error / enum` 定义位置

## 什么时候使用

- 新代码生成后做检查
- PR / diff 审查
- 架构一致性检查
- 风险识别与问题定位
- 需要结构化审查结果时

## 什么时候不要使用

- 仅要求直接生成代码
- 仅做概念解释
- 与 Kratos 无关的普通代码问题

## 事实来源

审查必须基于：

1. 仓库中的实际代码与改动
2. references下的检查清单与常见违规规则
3. `kratos-architecture` 及其 references 作为知识与事实依据

## 决策优先级

1. 明确的仓库事实
2. `kratos-reviewer` checklist
3. `common-violations`
4. `kratos-architecture` references知识文档内容
5. 一般性 Go 最佳实践

若仓库本地模式与 checklist 冲突，必须明确指出该模式是否是历史遗留，而不是默认合理。

## reference 路由

- 统一审查项 → `references/review-checklist.md`
- 常见错误模式 → `references/common-violations.md`
- 若需要具体规则依据，应回看 `kratos-architecture` 中对应 references

## 工作流

1. 先查看相关代码或改动。
2. 判断属于哪类检查：
   - 分层
   - 聚合根 / 命名
   - Wire / nil-check
   - `internal/pkg`
   - runtime 组件
   - ent / data-access
   - `error / enum`
   - etc. 
3. 按 checklist 检查。
4. 对照 `common-violations` 识别反模式。
5. 输出结构化审查结果。

## 必须遵守的高优先级规则

- 不要只说“有问题”，必须指出问题类型、位置和原因
- 不要把风格偏好当成架构违规
- 不要忽略仓库已有上下文直接套通用 Go 习惯
- 若无证据，不要臆测 bug
- 审查结果必须区分：架构问题 / 可选优化 / 风险提示

## 推荐输出骨架

```md
Review Result
- [ARCH] ...
  Evidence: ...
  Reason: ...
- [RISK] ...
  Evidence: ...
  Reason: ...
- [OPT] ...
  Evidence: ...
  Reason: ...
```
