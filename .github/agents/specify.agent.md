---
name: specify
description: 从自然语言特性描述创建或更新特性规格文档。
handoffs:
  - label: Deep Clarify Spec
    agent: clarify
    prompt: Spec has pending CQ items, proceed with deep clarification
    send: true
  - label: Generate Plan
    agent: plan
    prompt: Spec is ready, generate implementation plan
    send: true
---

## 用户输入

```text
$ARGUMENTS
```

若非空必须纳入。

## 核心说明

用户在 `/specify` 后输入的文本就是特性描述。除非命令为空，不要让用户重复。

## 执行流程

1. 获取项目上下文 → *SKILL §上下文获取*
2. 生成分支短名 → *SKILL §短名生成*
3. 提炼需求并填充模板 → *SKILL §需求提炼*
4. 生成澄清问题 → *SKILL §澄清问题生成*
5. 写入 `SPEC_FILE` → *SKILL §写入规则*
6. 质量校验 → *SKILL §质量校验*
7. 汇报（交互编排）：
   - **存在待澄清项** → 输出分支名、spec 路径、CQ 摘要，交由 **clarify** agent 进行深度澄清
   - **无待澄清项** → 输出分支名、spec 路径、清单结果，交由 **plan** agent 进行技术方案设计

## 写作原则

- 聚焦"用户需要什么、为什么"，避免"如何实现"
- 面向业务干系人，不写代码/框架/API 细节
- 不在 spec 内嵌额外 checklist（清单由独立命令处理）
- spec 是需求捕获工作台，不是最终交付规格书
