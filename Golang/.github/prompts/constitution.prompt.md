---
description: 创建或更新项目宪法，沉淀治理原则并同步校验关键入口文档。
agent: constitution
---

# /constitution 命令

调用 **constitution** agent 维护项目宪法。

## 何时使用

- 首次建立项目治理原则
- 调整安全、质量、发布、流程等高层规则
- 规则文件较多，需要统一抽象并定版

## 调用模板（主 Agent 使用）

```markdown
## 宪法目标
- [初始化 / 修订]

## 变更动机
- [为什么要改]

## 已知约束
- [用户已确认的红线/偏好]

## 规则更新范围
- [需要新增/修订的 project rules 文件]

## 联动更新范围
- [需要同步更新的 agent/prompt/command 入口]
```
