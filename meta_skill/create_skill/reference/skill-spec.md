# Skill Spec

## SKILL.md 骨架

创建任何技能，SKILL.md 必须包含以下章节，顺序固定：

```
---
name: skill-name
description: [判别性描述] [USE FOR 场景] [触发关键词] [DO NOT USE FOR]
---

## 输入
## 工作流
## 约束
## 强制输出
## 参考文件
```

---

## description 写法

决策表：description 必须包含哪些元素。

| 元素 | 是否必须 | 示例 |
|------|----------|------|
| 领域 + 核心能力描述 | MUST | `用于 Kratos 接入层开发与生成验证` |
| USE FOR 场景 | MUST | `适用于新增或修改服务接口、路由...` |
| 触发关键词 | MUST | `触发关键词：grpc、http、proto、wire` |
| DO NOT USE FOR | MUST | `DO NOT USE FOR：业务逻辑编排（→ kratos-domain）` |

```yaml
# ✅
description: 用于 Kratos 接入层开发与生成验证，包括 gRPC/HTTP 服务注册、Gateway 代理与协议适配。
  适用于新增或修改服务接口、路由、网关映射、proto 定义、wire 注入的场景。
  触发关键词：grpc、http、server、gateway、proto、wire、codegen。
  DO NOT USE FOR：业务逻辑编排（→ kratos-domain）、依赖注入设计（→ kratos-components）。

# ❌ 只写了功能，没有触发词和排除场景
description: 用于 Kratos 项目的接入层开发。

# ❌ 过于宽泛，无判别力
description: 用于代码开发任务。
```

---

## 章节写法

### 输入

```markdown
## 输入

- 必需：变更目标描述（接口 / 路由 / 生成类型）
- 可选：proto 文件路径、wire provider 路径
- 可选：`specs/<feature>/tasks.md`

缺少必需输入时，MUST 先向用户提问，不得猜测继续。
```

### 工作流

工作流写有序步骤 + IF/THEN 分支，不写纯描述。

```markdown
# ✅ 有序步骤 + 分支
## 工作流
1. 识别变更类型：server / gateway / codegen
2. IF tasks.md 存在 → 按任务逐步推进
   IF tasks.md 不存在 → 先生成最小任务草案，用户确认后继续
3. 按需加载参考文件
4. 执行变更
5. 执行 codegen 与 build 验证
6. 输出完成状态

# ❌ 纯描述，无分支，无顺序语义
## 工作流
根据任务类型修改对应的文件，完成后执行构建验证。
```

### 降级路径

```markdown
## 降级路径

| 缺失条件 | 降级行为 |
|----------|----------|
| 无法读取 SERVER_NAME | 按目录结构判定，显式声明假设 |
| tasks.md 不存在 | 先生成最小任务草案，告知用户确认 |
| 变更类型无法识别 | 停止并向用户提问，不猜测继续 |
```

### 约束

```markdown
## 约束

### MUST
- 接入层 MUST 只做协议适配，不承载业务编排

### MUST NOT
- MUST NOT 手改 `*.pb.go` 或 `wire_gen.go`

### SHOULD
- 协议变化后 SHOULD 同步检查枚举和错误码
```

边界案例：何时用 MUST vs SHOULD。

| 条件 | 级别 |
|------|------|
| 违反会直接导致错误、安全问题、不可逆后果 | MUST / MUST NOT |
| 默认应该遵守，但有合理理由可绕过 | SHOULD |
| 推荐做法，不做也不影响正确性 | 可省略或放 reference |

### 强制输出

```markdown
## 强制输出

开始前输出：

\```json
{
  "scope": "server | gateway | codegen",
  "changeTarget": "变更目标简述",
  "codegenPlan": ["make api", "wire gen"]
}
\```

完成后输出：

\```json
{
  "codegenDone": true,
  "buildPassed": true
}
\```
```

```markdown
# ❌ 自然语言列表，不可被下游解析
## 强制输出
开始前说明：scope、变更目标、生成计划
完成后确认：codegen 是否完成、build 是否通过
```

### 参考文件

```markdown
## 参考文件

| 文件 | 适用场景 | 加载时机 |
|------|----------|----------|
| `reference/proto-spec.md` | 新增或修改 proto 定义、枚举、错误码 | 按需 |
| `reference/server-spec.md` | 新增或修改 gRPC/HTTP 服务注册、路由 | 按需 |
```

```markdown
# ❌ 散列引用，无场景说明
## 参考文件
- `reference/proto-spec.md`
- `reference/server-spec.md`
```

---

## 文件结构

```
skill-name/
├── SKILL.md           # 执行手册：工作流 + 约束 + 强制输出 + 参考文件清单
└── reference/         # 知识库：决策表 + 正反例（按需加载）
    └── xxx-spec.md
```

```
# ✅ 正确结构

skill-name/
├── SKILL.md
└── reference/
    ├── proto-spec.md
    └── server-spec.md

# ❌ 错误结构 —— 不应存在 rules/ 文件夹
skill-name/
├── SKILL.md
├── rules/              ← 删除，内容内联到 SKILL.md ## 约束
│   └── proto-rule.md
└── reference/
```
