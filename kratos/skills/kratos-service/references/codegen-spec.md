# Codegen Spec

## 生成链路判定

| 条件 | 做法 |
|------|------|
| 修改 proto | 执行 proto 相关生成并检查产物差异 |
| 修改 wire/provider | 执行 wire 生成 |
| 修改 ent schema | 执行 ent 生成 |
| 修改注册链路或接入层 | 至少执行最小 build 验证 |

```text
// ✅ 先判定触发链路
proto change -> make generate
wire change  -> wire
ent change   -> go generate ./internal/data/ent
```

```text
// ❌ 不区分改动类型，提交前完全不跑生成
edit files -> git commit
```

---

## 生成物处理

| 条件 | 做法 |
|------|------|
| 生成结果超出预期 | 回看源定义和生成命令 |
| 生成物需要变更 | 通过修改源文件触发，不直接手改产物 |

```bash
# ✅ 修改 proto 后重新生成
make generate
```

```go
// ❌ 直接手改 *.pb.go / wire_gen.go
func (x *GetAccountRequest) GetId() uint32 { ... }
```

---

## 最小验证顺序

| 条件 | 做法 |
|------|------|
| 改动触发生成 | 生成后检查 diff，再跑 build |
| 只改 service/server 映射 | 至少跑 `go build ./...` 或定向构建 |

```text
// ✅ 推荐顺序
修改源定义
-> 执行生成
-> 检查 diff
-> go build ./...
-> 定向测试
```

```text
// ❌ 只跑生成，不做构建
修改 proto
-> make generate
-> commit
```

---

## 常用命令模板

| 条件 | 做法 |
|------|------|
| proto 生成 | `make generate` 或项目约定命令 |
| wire 生成 | `wire` 或 `wire ./cmd/server` |
| ent 生成 | `go generate ./internal/data/ent` |
| 全量构建 | `go build ./...` |

```bash
# ✅ proto
make generate
```

```bash
# ✅ wire
wire ./cmd/server
```

```bash
# ✅ ent
go generate ./internal/data/ent
```

```bash
# ✅ build
go build ./...
```

---

## 组合场景

```text
修改 api/open/v1/account.proto
-> make generate
-> 检查 *.pb.go / *_grpc.pb.go / validate 产物 diff
-> 如果 service/provider 有联动，再执行 wire ./cmd/server
-> go build ./...
```

这个组合场景同时满足：

- 改源文件即补生成
- 不直接修产物
- 生成与构建形成闭环

---

## 常见错误模式

```text
// ❌ 改 proto 不跑生成
api/account.proto
```

```text
// ❌ 手改 wire_gen.go
cmd/server/wire_gen.go
```

```text
// ❌ 只跑 codegen 不跑 build
make generate
```
