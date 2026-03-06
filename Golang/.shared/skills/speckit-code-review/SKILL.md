---
name: speckit-code-review
description: |
  用于 Speckit 研发流程的代码质量审查，包括依据 spec、plan 和 tasks 检查已变更代码的验收完成度与架构合规性，甄别性能与安全坏味道。适用于代码落地实施 (implement) 完成后，准备合并、提交，或需要借助上级架构规约对底层实现进行拦截查杀的场景。触发关键词包括 code review、review、代码审查、架构合规、隐患检查。
---

# Spec Kit Code Review Skill（中文）

## 何时使用

- implement 执行完毕后，对代码做质量审查。
- 修复后再次审查，直到通过。

## 输入

- 必需：`specs/<feature>/tasks.md`（期望基线）
- 必需：`specs/<feature>/plan.md`（设计参照）
- 可选：`specs/<feature>/spec.md`（需求追溯）
- 实际代码文件（基于 tasks.md 中的文件路径）

## 目标

在代码提交前，识别完成度缺口、架构违规、性能隐患和安全问题。

## 约束

- **只读审查**：输出报告，不自动修改代码
- 报告写入 `specs/<feature>/review.md`
- 每次审查生成完整报告（覆盖上一次）

## 工作流

### §前置检查

1. 确认 `specs/<feature>/tasks.md` 存在。
2. 确认 `specs/<feature>/plan.md` 存在。
3. 若任一缺失，终止并提示先运行 `/tasks` 和 `/plan`。
4. 检查 tasks.md 中是否有 `[x]` 标记的已完成任务，若全部未完成则终止并提示先运行 `/implement`。

### §上下文加载

#### 从文档提取

| 文档 | 提取内容 |
|------|---------|
| tasks.md | 已完成任务 `[x]`、文件路径、RQ 引用、Phase 分组 |
| plan.md | 项目类型、Phase 1 设计表、框架适配扩充、测试策略 |
| spec.md | RQ 列表、约束与边界（若存在） |

#### 从代码提取

基于 tasks.md 中标记 `[x]` 的任务的文件路径，逐个读取实际代码文件。

优先级：
1. 先读 `internal/biz/` 下的代码（核心业务逻辑）
2. 再读 `internal/data/` 下的代码（数据访问）
3. 然后 `internal/service/` 下的代码（API 实现）
4. 最后 `internal/server/`、Ent Schema、Proto 等辅助文件

### §审查维度

按以下 4 个维度逐一审查，每个维度独立评分。

#### A. 完成度审查

对照 tasks.md 中标记 `[x]` 的任务：

| 检查项 | 规则 |
|--------|------|
| 文件存在 | 任务声明的文件路径必须存在 |
| 方法实现 | 任务描述的方法/函数必须有实际实现（非空桩、非 TODO） |
| RQ 覆盖 | 每个 RQ 至少有一个对应的已实现任务 |
| Wire 注册 | 声明的 Provider 必须已注册到 Wire ProviderSet |
| 服务注册 | gRPC/HTTP 服务必须已注册到 server |

对照 tasks.md 中标记 `[ ]` 的未完成任务：
- 列出未完成任务清单，标注影响范围

#### B. 架构合规审查

基于 kratos-patterns reference 规范检查：

| 检查项 | 规则 | 参考规范 |
|--------|------|---------|
| 层隔离 | biz 不 import data 包；service 不直接 import data 包 | layer-spec |
| 依赖方向 | service → biz → data（单向依赖） | layer-spec |
| 接口契约 | biz 层通过 Repo 接口访问数据，不直接操作 Ent Client | layer-spec |
| Depend 包装 | InnerRPC 调用必须通过 depend 包装，不直接调用 gRPC Client | depend-spec |
| 命名规范 | 文件名、方法名、变量名符合命名规范 | naming-spec |
| Proto 规范 | Proto 文件结构、字段命名、注解符合规范 | proto-spec |
| Ent Schema | Schema 定义、索引、边关系符合规范 | ent-spec |
| 枚举使用 | 枚举值定义和使用符合规范 | enum-spec |
| 错误处理 | 使用项目定义的错误类型，非裸 error | error-spec |
| 国际化 | 用户可见文本使用 i18n，非硬编码 | i18n-spec |
| Wire DI | Provider 函数签名正确，ProviderSet 组织合理 | codegen-spec |

#### C. 性能审查

| 检查项 | 模式 / 特征 |
|--------|------------|
| N+1 查询 | 循环内调用 Repo/数据库方法；应使用批量查询或 Eager Loading |
| 无分页查询 | 列表查询缺少 `.Limit()` / `.Offset()` 或分页参数 |
| 缺失索引 | Ent Schema 中频繁查询字段无 `.Indexes()` 定义 |
| 全表扫描 | 使用 `.All(ctx)` 无 `.Where()` 条件的大表查询 |
| Goroutine 泄漏 | 启动 goroutine 无退出机制（无 context/channel/done 信号） |
| 未关闭资源 | 打开文件/连接/Body 后无 `defer Close()` |
| 大对象拷贝 | 大 struct 值传递（应用指针） |
| 重复计算 | 循环内重复执行不变的计算或查询 |

#### D. 安全审查

| 检查项 | 模式 / 特征 |
|--------|------------|
| SQL 注入 | 使用字符串拼接构造查询（Ent 通常安全，但 raw SQL 需注意） |
| 权限检查 | 涉及数据变更的 biz 方法是否有权限/角色校验 |
| 敏感字段 | 密码、手机号、身份证等字段在响应中未脱敏 |
| 日志泄露 | 日志中输出了密码、token、密钥等敏感信息 |
| 硬编码密钥 | 代码中硬编码 secret/password/token（应从配置读取） |
| 输入校验 | Proto validate 标签是否覆盖用户输入字段 |
| 越权访问 | 数据查询是否包含租户/用户维度过滤 |

### §严重级别

| 级别 | 定义 | 处理方式 |
|------|------|---------|
| **BLOCKER** | 功能未实现（RQ 零覆盖）；安全漏洞（注入、越权、密钥泄露） | 必须修复后才能提交 |
| **MAJOR** | 架构违规（层穿透、依赖方向错误）；严重性能问题（N+1、全表扫描） | 强烈建议修复 |
| **MINOR** | 命名不规范；缺失 i18n；Wire 组织不规范 | 建议修复 |
| **INFO** | 代码风格优化；可选的性能微调 | 可选修复 |

### §报告生成

输出 Markdown 格式报告，写入 `specs/<feature>/review.md`。

#### 报告结构

```markdown
# Code Review Report

> Feature: <feature-short-name>
> 审查时间: <timestamp>
> 项目类型: <BaseService/业务项目/网关类项目>

## 审查摘要

| 维度 | BLOCKER | MAJOR | MINOR | INFO | 评分 |
|------|---------|-------|-------|------|------|
| 完成度 | | | | | /10 |
| 架构合规 | | | | | /10 |
| 性能 | | | | | /10 |
| 安全 | | | | | /10 |
| **总计** | | | | | **/40** |

## 完成度

### 已完成任务

| Task ID | RQ | 描述 | 文件 | 状态 |
|---------|-----|------|------|------|
| T-xxx | RQ-xxx | ... | ... | ✅ 已实现 / ⚠️ 部分实现 / ❌ 空桩 |

### 未完成任务

| Task ID | RQ | 描述 | 影响 |
|---------|-----|------|------|

### RQ 覆盖

| RQ | 任务覆盖 | 代码覆盖 | 状态 |
|----|---------|---------|------|

## 发现清单

| ID | 维度 | 严重级别 | 文件:行号 | 摘要 | 建议 |
|----|------|---------|-----------|------|------|
| CR-001 | ... | ... | ... | ... | ... |

## 修复建议优先级

按 BLOCKER → MAJOR → MINOR 排序，给出具体修复方案。
```

#### 评分规则

每个维度 10 分起，按发现项扣分：

| 级别 | 每项扣分 |
|------|---------|
| BLOCKER | -3 |
| MAJOR | -2 |
| MINOR | -1 |
| INFO | 0 |

最低 0 分，不出现负分。

### §下一步建议

根据审查结果：

- 有 BLOCKER → **必须修复**：建议运行 `/implement` 修复后再次 `/code-review`
- 仅 MAJOR → **建议修复**：列出具体修复点，建议修复后再次 `/code-review`
- 仅 MINOR/INFO → **可以提交**：列出改进建议，不阻塞提交
- 全部通过（总分 ≥ 36/40 且无 BLOCKER/MAJOR）→ **审查通过** ✅

## 输出

- 审查报告：`specs/<feature>/review.md`
