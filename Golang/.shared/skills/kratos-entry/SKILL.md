---
name: kratos-entry
description: |
  用于 Kratos 接入层与生成验证，包括 gRPC/HTTP 服务注册、Gateway 代理与协议适配、proto/wire/ent 生成和构建校验。适用于新增或修改服务接口、路由、网关映射、代理逻辑、proto 定义、wire 注入、codegen 或构建验证的场景。触发关键词包括 grpc、http、server、gateway、route、proxy、proto、wire、codegen、make api、go build。
---

# Kratos Entry

## 必读规则

- `./rules/proto-rule.md`
- `./rules/server-rule.md`
- `./rules/codegen-rule.md`

## 按需参考

- Proto：`./reference/proto-spec.md`
- Server：`./reference/server-spec.md`
- Gateway：`./reference/gateway-spec.md`
- Wire：`./reference/wire-spec.md`
- Codegen：`./reference/codegen-spec.md`

## 读取顺序

先读 `./rules/*.md` 明确接入边界和生成约束，再按当前任务读取必要的 `./reference/*.md`。

按需读取以下参考文档：

- Server：`./reference/server-spec.md`
- Gateway：`./reference/gateway-spec.md`
- Codegen：`./reference/codegen-spec.md`
- Proto：`./reference/proto-spec.md`
- Wire：`./reference/wire-spec.md`

## 何时使用

- 新增或修改 gRPC/HTTP 服务实现、注册、路由
- 新增或修改 Gateway/Openapi 代理逻辑、参数映射、协议适配
- 修改 proto、wire provider、ent schema 后需要生成物校验
- 提交前需要确认生成物和 `go build` 可通过

## 核心约束

1. 接入层只做协议适配、注册和代理，不承载业务编排。
2. Gateway 只做代理与协议转换，不把业务逻辑下沉到网关。
3. 任何接口定义、依赖注入或 schema 变化，都要同步完成 codegen 和 build 验证。
4. 路由注册与 provider 更新按规范位置落位，避免散落式注册。

## 实施流程

1. 判断本次变更属于 server、gateway 还是 codegen 校验。
2. 修改服务注册、路由、代理或接口定义。
3. 执行对应的 proto/wire/ent 生成。
4. 执行 `go build` 或仓库规定的最小编译验证。

## 强制输出

开始前输出：

- `EntryScope:` server / gateway / codegen
- `RouteChange:` 路由、代理或接口变更摘要
- `CodegenPlan:` 计划执行的生成与构建命令

提交前输出：

- 服务注册是否按规范落位（Yes/No）
- Gateway 是否未承载业务逻辑（Yes/No）
- 是否完成 codegen + build 验证（Yes/No）
