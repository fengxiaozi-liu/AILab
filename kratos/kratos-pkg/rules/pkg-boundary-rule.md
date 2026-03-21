# Pkg Boundary Rule

## Principles

- `internal/pkg` 只承载跨层、稳定、可复用的基础能力。
- 进入 `internal/pkg` 的能力必须脱离具体业务聚合和具体服务语义后仍然成立。

## Specification

- `context`、`metadata`、`middleware`、`proto helper`、`schema`、`seata`、`util` 这类基础能力可以放入 `internal/pkg`。
- 同一能力被多个模块复用，且边界明确、接口稳定时，才允许沉淀到 `internal/pkg`。
- `internal/pkg` 中的类型、函数和目录命名要围绕基础能力本身，不围绕具体业务动作命名。

## Prohibit

- 禁止把聚合对象、业务状态流转、业务校验、业务规则判断放入 `internal/pkg`。
- 禁止仅为单个调用点创建 `util`、`helper`、`common` 壳目录或壳函数。
- 禁止在 `internal/pkg` 中拼装业务 `reply`、聚合关系或下沉业务 `repo/usecase` 逻辑。

## Checklist

- 这段能力脱离当前业务服务后是否仍有复用价值。
- 这段能力是否可以用“基础设施能力”而不是“业务流程”来描述。
- 这段能力是否已经在其他 pkg 子目录存在近义实现。
