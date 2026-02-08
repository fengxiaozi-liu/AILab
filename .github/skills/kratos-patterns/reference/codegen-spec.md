# 代码生成规范（Kratos 通用）

> 本文档描述 Kratos 项目中常见生成物的通用规则：proto / wire / ent 等。

## 核心原则

1. **生成物禁止手改**：生成目录应当视为 build artifact
2. **生成流程可重复**：本地与 CI 执行应一致
3. **变更要可追溯**：proto/ent 变更需在 MR/PR 中清晰呈现

## 典型生成物

- Proto：API 契约、DTO、validate
- Wire：依赖注入装配代码
- Ent：ORM schema -> 代码

## Go Demo（示意）

### 1) Wire

```bash
# 示例：在 cmd/server 下执行 wire
cd cmd/server
wire
```

### 2) Proto

```bash
# 示例：使用 buf 或 make generate（以项目为准）
make generate
```

## 常见陷阱

- 生成文件被手改导致 diff 难以维护
- 本地生成与 CI 生成不一致（版本/参数不同）
- 没有锁定工具版本导致生成结果漂移
