---
description: 从自然语言需求生成特性规格文档，捕获需求并追踪澄清。
agent: specify
---

# /specify 命令

调用 **specify** agent 从需求描述生成特性规格。

## 何时使用

- 有新的功能需求，需要结构化记录
- 需求模糊，需要识别并追踪待澄清项
- 已有 spec 需要补充澄清结果

## 产出物

```
specs/<feature>/
├── spec.md                        ← 特性规格（需求描述 + 待澄清问题）
└── checklists/
    └── requirements.md            ← 质量校验清单
```
