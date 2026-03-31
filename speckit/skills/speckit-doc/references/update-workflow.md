# Speckit Doc Update Workflow

## `init`

适用于仓库首次建档或已有架构文档质量不足，需要基于当前仓库事实重新整理时。

### 执行步骤

1. 扫描仓库目录、README、关键配置、现有系统文档
2. 识别系统边界、核心模块、外部依赖、关键流程、核心实体
3. 读取 `templates/architecture-template.md`
4. 按标准模板写入 `docs/architecture.md`

## `update`

适用于已有 `docs/architecture.md`，且当前 feature 的 speckit 产物已经形成时。

### 执行步骤

1. 读取现有 `docs/architecture.md`
2. 读取当前 feature 的 `spec.md`、`plan.md`、`tasks.md`
3. 标记以下变更点：
   - 新增或变更模块
   - 新增或变更核心流程
   - 新增或变更实体
   - 新增依赖、约束、边界
4. 回到已有文档中做定点更新
5. 若功能尚未完成，使用“规划中”或“待实现”标记
6. 将更新结果直接合并到对应章节，不追加过程记录

## 决策规则

- `spec.md` 主要提供业务目标、需求边界、核心实体线索
- `plan.md` 主要提供模块划分、接口/流程设计、关键决策
- `tasks.md` 主要提供实施范围和完成状态
- 若三者与代码事实冲突，以代码事实为“当前状态”，以 speckit 文档为“目标状态”

## 禁止事项

- 不得把未实现设计写成已上线事实
- 不得把实现细节直接堆成架构文档
- 不得用 feature 级临时术语破坏全局文档一致性
