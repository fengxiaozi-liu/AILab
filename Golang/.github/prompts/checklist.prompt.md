---
description: 基于特性上下文生成定制化需求质量检查清单。
agent: checklist
---

# /checklist 命令

调用 **checklist** agent 生成需求质量检查清单。

## 何时使用

- 想检查 spec 的需求是否完整、清晰、一致
- 在 plan/tasks 之前或之后均可使用
- 需要针对特定主题（安全、边界、接口设计等）做需求审查

## 前置条件

- `specs/<feature>/spec.md` 存在

## 产出物

```
specs/<feature>/
├── spec.md
├── checklists/
│   └── <domain>.md      ← 需求质量清单（本命令产出）
```

## 工作流程

```
/checklist [主题/焦点（可选）]
    ↓
前置检查 → 确认 spec 存在
    ↓
动态澄清（最多 3 个问题）
    ↓
确定清单主题
    ↓
读取 spec/plan/tasks 相关片段
    ↓
生成清单 → 输出报告
```
