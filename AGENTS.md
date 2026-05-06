# AGENTS.md

本仓库是 AI 支持资产仓库，不是业务应用服务仓库。Agent 在这里工作时，应保持资产语义稳定、控制实现改动范围，并以 `specs/` 下的规格产物作为当前特性的意图来源。

## 当前特性上下文

当前 ferryPilot Go 重写相关文档位于：

- `specs/rewrite-ferrypilot-go/spec.md`
- `specs/rewrite-ferrypilot-go/plan.md`
- `specs/rewrite-ferrypilot-go/tasks.md`

实施时必须按 `tasks.md` 分 Phase 推进。只有相关文件已经真实落地并完成验证后，才能标记任务完成。

## 硬性边界

- Go 重写版本位于 `utils/ferryPilotGo`。
- 命令行工具名称统一为 `ferryPilot`。
- 不得依赖、导入、复制、执行或机械翻译 `utils/ferryPilot` 下的旧 Python 实现文件。
- 现有 `utils/ferryPilot/README.md` 和用户可见行为只能作为需求参考。
- 除非规格明确要求，不得改变 `AISupport/` 下资产的业务语义。

## 常用命令

```bash
cd utils/ferryPilotGo
go test ./...
go build -o bin/ferryPilot ./cmd/ferryPilot
```

## 文档维护要求

以下项目级文档应与实现事实保持同步：

- `README.md`
- `AGENTS.md`
- `docs/`
- `utils/ferryPilotGo/README.md`

