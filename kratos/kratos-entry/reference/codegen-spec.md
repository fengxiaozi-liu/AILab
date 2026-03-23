# Codegen Spec

## 生成链路决策

| 改动类型 | 执行命令 | 检查产物 |
|---------|---------|---------|
| 修改 `.proto` 文件 | `make generate` | `api/.../v1/*.pb.go` |
| 修改 `ent/schema/*.go` | `ent generate ./ent/schema` | `ent/*.go` |
| 修改 wire ProviderSet / 构造函数 | `cd cmd/server && wire` | `cmd/server/wire_gen.go` |
| 以上任意变动后 | `go build ./...` | 编译无错误 |

---

## proto 生成

```bash
# ✅ 修改 .proto 后执行
make generate

# ✅ 验证生成产物无误
go build ./...

# ❌ 手动修改 .pb.go 文件（会被下次生成覆盖）
# vim api/admin/account/v1/account.pb.go  ❌
```

---

## ent 生成

```bash
# ✅ 修改 ent/schema/*.go 后执行
ent generate ./ent/schema

# ✅ 也可通过 make 封装
make ent-generate

# ❌ 直接修改 ent/ 下的生成文件
# vim ent/account.go  ❌
```

---

## wire 生成

```bash
# ✅ 修改 ProviderSet 或构造函数后执行
cd cmd/server
wire

# ✅ wire 生成后校验编译
go build ./cmd/server/...

# ❌ 手动修改 wire_gen.go
# vim cmd/server/wire_gen.go  ❌
```

---

## ⚠️ 工具版本对齐

```bash
# ⚠️ 确认本地工具版本与 CI 一致，否则产物可能有格式差异
protoc --version
wire --version
ent version

# ⚠️ 生成产物需提交到仓库（CI 不重新生成），保证产物与源定义一致
git diff --stat  # 检查是否有未提交的生成文件
```

---

## 组合场景

```bash
# 完整：新增 Account 聚合根后的生成流程
# 1. 新增 ent schema
echo "新增 ent/schema/account.go"
ent generate ./ent/schema

# 2. 新增 proto 定义
echo "新增 api/admin/account/v1/account.proto"
make generate

# 3. 新增 wire provider
echo "更新 ProviderSet"
cd cmd/server && wire

# 4. 全量校验
cd ../.. && go build ./...
```

---

## 常见错误模式

```bash
# ❌ 修改 proto 后忘记执行 make generate
# 导致 .pb.go 与 .proto 不一致，运行时 panic

# ❌ 修改 wire ProviderSet 后忘记 wire 生成
# 报错：wire: no provider found for ...

# ❌ 手动修改生成文件
# 下次 make generate 覆盖手动改动

# ❌ 本地工具版本低于 CI
# 生成的 .pb.go 格式不同，PR 产生噪音 diff
```
