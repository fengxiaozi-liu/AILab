# Change Control Rule

## Principles

- 证据优先，先确认再修改。
- Contract First，跨边界能力先定协议再落实现。

## Specification

- 修改前必须先验证仓库事实，不凭猜测实施。
- 仅实现与需求直接相关的变更，不顺手做无关重构。
- 变更 proto、wire、ent schema 时，后续必须进入生成与校验链路。

## Prohibit

- 禁止未确认需求边界时扩大改动范围。
- 禁止修改生成物代替修改源文件。
