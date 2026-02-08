---
description: 对 spec/plan/tasks 三件套做只读一致性与质量分析。
agent: analyze
---

# /analyze 命令

调用 **analyze** agent 在实现前做质量门检查。

## 何时使用

- tasks.md 生成后，实现前做一致性与覆盖率检查
- 怀疑 spec/plan/tasks 之间存在不一致时

## 前置条件

已生成完整三件套：
- `specs/<feature>/spec.md`
- `specs/<feature>/plan.md`
- `specs/<feature>/tasks.md`

## 产出物

只读分析报告（不写入文件），包含：
- 发现表（最多 50 条，按严重级别排序）
- RQ 覆盖汇总
- 统计指标

## 工作流程

```
/analyze [关注点（可选）]
    ↓
前置检查 → 确认三件套齐全
    ↓
加载文档关键段落
    ↓
构建 RQ→设计项→Task 追踪链
    ↓
7 类检测（覆盖缺口、不一致、重复、歧义、欠规范、宪法对齐、风险交叉）
    ↓
输出报告 + 下一步建议
```

## 示例

```
User: /analyze

Agent:
## 规格分析报告

已分析：`specs/manage-portfolio/`

### 统计

| 指标 | 数值 |
|------|------|
| RQ 总数 | 3 |
| 任务总数 | 19 |
| RQ→任务覆盖率 | 100% |
| CRITICAL | 0 |
| HIGH | 1 |

### 发现

| ID | 类别 | 严重级别 | 位置 | 摘要 | 建议 |
|----|------|---------|------|------|------|
| F-001 | 不一致 | HIGH | plan.md / tasks.md | Portfolio 实体在 plan 中有 6 个字段，tasks T005 只引用了 4 个 | 补齐 T005 描述或拆分任务 |
| F-002 | 歧义 | MEDIUM | spec.md RQ-003 | "高效查询"未量化 | 补充具体性能指标 |

无 CRITICAL 问题，可继续实施。建议先处理 F-001。

需要对高优先级问题给出具体修复建议吗？
```
