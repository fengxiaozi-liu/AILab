# 任务清单：[FEATURE NAME]

**Spec**: `specs/[feature]/spec.md`
**Plan**: `specs/[feature]/plan.md`
**Created**: [DATE]
**项目类型**: [BaseService / 业务 / 网关]

---

## 格式说明

- `[TaskID]`：`T001` 起顺序编号
- `[P]`：可并行执行（同 Phase 内、不同文件、无依赖冲突）
- `[RQ-xxx]`：该任务实现的 spec 需求编号（基建型任务可省略）
- 描述中必须包含精确文件路径

---

<!--
  ============================================================================
  重要说明：

  以下 Phase 结构根据项目类型选择对应模板。
  Phase 顺序严格按照 kratos-patterns「进行工作」步骤排列。
  每个任务必须标注 RQ 来源（基建型除外）。

  项目类型：
  - BaseService → 使用 BaseService Phase 模板
  - 业务项目   → 使用业务项目 Phase 模板
  - 网关类项目 → 使用网关类 Phase 模板

  根据判定结果保留对应章节，删除其余。
  ============================================================================
-->

## BaseService 项目

### Phase 1：Proto 接口设计与定义

**产出物**: `api/<domain>/<scope>/v1/*.proto`

- [ ] T001 [RQ-xxx] [描述] `[文件路径]`

### Phase 2：代码生成

**产出物**: Proto → Go 生成代码

- [ ] T00x 执行 Proto 代码生成

### Phase 3：枚举与异常定义

**产出物**: `internal/enum/<domain>/`、`internal/error/<domain>/`

- [ ] T00x [P] [RQ-xxx] [定义枚举] `[文件路径]`
- [ ] T00x [P] [RQ-xxx] [定义异常] `[文件路径]`

### Phase 4：InnerRPC 依赖包装

**产出物**: `internal/biz/depend/`、`internal/data/depend/`

- [ ] T00x [RQ-xxx] [包装依赖] `[文件路径]`

### Phase 5：国际化

**产出物**: `assets/i18n/active.*.toml`

- [ ] T00x [RQ-xxx] [编写 i18n] `[文件路径]`

### Phase 6：收尾

- [ ] T00x 文档更新
- [ ] T00x 代码审查

---

## 业务项目

### Phase 1：Ent Schema

**产出物**: `internal/data/ent/schema/*.go`

- [ ] T001 [P] [RQ-xxx] [编写 Schema] `[文件路径]`

### Phase 2：Ent 代码生成

**产出物**: `internal/data/ent/` 生成文件

- [ ] T00x 执行 Ent 代码生成

### Phase 3：Biz 层实现

**产出物**: `internal/biz/<domain>/`

- [ ] T00x [RQ-xxx] [实现业务逻辑] `[文件路径]`

### Phase 4：Data 层实现

**产出物**: `internal/data/<domain>/`

- [ ] T00x [RQ-xxx] [实现数据访问] `[文件路径]`

### Phase 5：Service 层实现

**产出物**: `internal/service/<domain>/`

- [ ] T00x [RQ-xxx] [实现 API] `[文件路径]`

### Phase 6：服务注册

**产出物**: `internal/server/` 路由注册

- [ ] T00x 注册服务路由 `[文件路径]`

### Phase 7：Wire 代码生成

**产出物**: `cmd/server/wire_gen.go`

- [ ] T00x 生成 Wire 依赖注入代码

### Phase 8：测试

- [ ] T00x [P] [RQ-xxx] 编写单元测试 `[文件路径]`
- [ ] T00x [P] [RQ-xxx] 编写集成测试 `[文件路径]`

### Phase 9：收尾

- [ ] T00x 文档更新
- [ ] T00x 代码审查

---

## 网关类项目

### Phase 1：Proxy 层实现

**产出物**: 代理路由实现

- [ ] T001 [RQ-xxx] [实现代理逻辑] `[文件路径]`

### Phase 2：服务注册

**产出物**: `internal/server/` 路由注册

- [ ] T00x 注册服务路由 `[文件路径]`

### Phase 3：Wire 代码生成

**产出物**: `cmd/server/wire_gen.go`

- [ ] T00x 生成 Wire 依赖注入代码

### Phase 4：收尾

- [ ] T00x 文档更新
- [ ] T00x 代码审查

---

## 依赖与执行顺序

### Phase 间依赖

<!--
  Phase 严格按序执行，前置 Phase 未完成不可开始后续 Phase。
  同一 Phase 内标记 [P] 的任务可并行。
-->

### 并行规则

- 同 Phase 内标记 `[P]` 的任务：操作不同文件且无依赖冲突，可并行
- 跨 Phase 严格串行：后置 Phase 依赖前置 Phase 的产出物
- 典型可并行场景：同层多个实体的 Schema/Biz/Data 编写

---

## 统计

| 指标 | 数值 |
|------|------|
| 总任务数 | [N] |
| 可并行任务数 | [N] |
| 覆盖 RQ 数 | [N] / [Total RQ] |
