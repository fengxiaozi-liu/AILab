---
description: 基于已就绪的 spec 生成 Kratos 微服务实施计划。
agent: plan
---

# /plan 命令

调用 **plan** agent 生成技术实施计划。

## 何时使用

- spec 已就绪（Status 为 Ready），需要进入技术设计阶段
- 需要从需求推导 Ent Schema、Proto 接口、分层实现方案

## 前置条件

已使用 `/specify` + `/clarify` 生成：
- `specs/<feature>/spec.md`（Status 为 Ready）

## 产出物

```
specs/<feature>/
├── spec.md              ← 特性规格（specify + clarify 产出）
└── plan.md              ← 实施计划（本命令产出）
```

## 工作流程

```
/plan [技术约束（可选）]
    ↓
前置检查 → 确认 spec Ready
    ↓
判定项目类型（BaseService / 业务 / 网关）
    ↓
Phase 0：领域调研（实体、状态、依赖）
    ↓
Phase 1：技术设计（Ent Schema / Proto / 枚举 / 异常）
    ↓
框架适配：用 kratos-patterns 扩充技术设计
    ↓
质量校验 → 汇报结果
```
