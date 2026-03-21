# I18n Rule

## Principles

- i18n key 稳定，翻译产物通过工具链生成。

## Specification

- 文案 key 按语义命名，不直接散落业务代码。
- 翻译文件通过 `goi18n extract` 和翻译工具生成。

## Prohibit

- 禁止手改 `active.*.toml` 产物。
