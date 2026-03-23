---
name: kratos-entry
description: |
  用于 Kratos 接入层开发与生成验证，包括 gRPC/HTTP 服务注册、Gateway 代理与协议适配、proto/wire/ent 生成和构建校验。
  适用于新增或修改服务接口、路由、网关映射、代理逻辑、proto 定义、wire 注入、codegen 或构建验证的场景。
  触发关键词：grpc、http、server、gateway、route、proxy、proto、wire、codegen、make api、go build。
  DO NOT USE FOR：业务逻辑编排（→ kratos-domain）、依赖注入设计（→ kratos-components）。
---

# Kratos Entry

## 输入

- 必需：变更目标描述（接口 / 路由 / 生成类型）
- 可选：proto 文件路径、wire provider 路径、ent schema 路径
- 可选：`specs/<feature>/tasks.md`

## 工作流

1. 识别本次变更类型：`server` / `gateway` / `codegen`
2. 输出开始前结构化状态（见强制输出）
3. 按需加载对应参考文件（见参考文件清单）
4. 执行变更：服务注册 / 路由 / 代理 / 接口定义
5. 执行对应生成：`make api` / `wire gen` / `ent generate`
6. 执行最小构建验证：`go build ./...`
7. 输出完成后结构化状态（见强制输出）

## 约束

### MUST
- 接入层 MUST 只做协议适配、注册和代理，不承载业务编排
- Gateway MUST 只做代理与协议转换，不把业务逻辑下沉到网关
- 对外接口 MUST 做认证、授权、输入校验
- 上游调用和代理 MUST 设置超时
- proto 修改 MUST 先改 `.proto`，再生成，再实现（Contract First）
- proto side MUST 只使用 `admin`、`open`、`inner` 三类目录
  - `admin` 只允许引用 `admin` 和 `base/*`
  - `open` 只允许引用 `open` 和 `base/*`
  - `inner` 只允许引用 `inner` 和 `base/*`
  - `base/*` 只承载稳定公共结构（`Paging`、`Sort`、`TimeRange`、`TransField`）
- 聚合根主对象 MUST 使用顶层 message；仅服务于单个 RPC 的局部结构使用嵌套 message
- 路由注册与 provider MUST 按约定位置落位，不散落
- proto/wire/ent 变更 MUST 完成 codegen 后再提交
- codegen 后 MUST 执行 `go build` 验证无报错

### MUST NOT
- MUST NOT 手改 `*.pb.go` 或 `wire_gen.go`（生成物不可手改）
- MUST NOT 只改源文件不更新生成物
- MUST NOT 对外暴露底层错误细节
- MUST NOT 存在无认证保护的公开入口
- MUST NOT `admin` 与 `open` 相互引用
- MUST NOT `admin`、`open` 引用 `inner`
- MUST NOT 跨 side import 其他 public proto
- MUST NOT 把 `base/*` 当业务 side 扩展目录使用

### SHOULD
- 协议变化后 SHOULD 同步检查枚举、错误码和状态语义
- 新增依赖或服务 SHOULD 同步检查 wire provider 链路
- 聚合根主 request 中的主对象 ID SHOULD 优先使用 `id`，不追加聚合根前缀

## 强制输出

开始前输出：

```json
{
  "entryScope": "server | gateway | codegen",
  "changeTarget": "变更目标简述",
  "codegenPlan": ["make api", "wire gen"]
}
```

完成后输出：

```json
{
  "serverRegistrationCorrect": true,
  "gatewayClean": true,
  "codegenDone": true,
  "buildPassed": true
}
```

## 参考文件

| 文件 | 适用场景 | 加载时机 |
|------|----------|----------|
| `reference/proto-spec.md` | 新增或修改 proto 定义、枚举、错误码 | 按需 |
| `reference/server-spec.md` | 新增或修改 gRPC/HTTP 服务注册、路由 | 按需 |
| `reference/gateway-spec.md` | Gateway 代理、参数映射、协议适配变更 | 按需 |
| `reference/wire-spec.md` | 修改 wire provider 或依赖注入链路 | 按需 |
| `reference/codegen-spec.md` | 执行 proto/wire/ent 生成与构建验证 | 按需 |
