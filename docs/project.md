# 项目概览

本项目用于打包 AI 研发支持资产，并提供 CLI 安装器，把这些资产安装到受支持的 Agent 运行时或项目目录中。

## 仓库结构

```text
AISupport/
  kratos/
  speckit/
utils/
  ferryPilotGo/
specs/
docs/
```

## AISupport Packages

`AISupport/<package>` 是安装单元。`AISupport/` 下的一级目录会被 `ferryPilot` 展示为可选择的 package。

当前 package 包括：

- `speckit`：结构化规格、澄清、计划、任务、实施、评审和总结工作流。
- `kratos`：面向 Kratos 项目的 skills 和 references。

## 安装器

`ferryPilot` 是本仓库的安装器。当前实现使用 Go 编写，位于 `utils/ferryPilotGo`。

安装器支持：

- 使用 `-g` 进行全局安装
- 使用 `-p` 安装到当前项目
- 使用 `-t` 选择 target agent
- 通过参数或交互提示选择 package
- 基于 `config/file_map.json` 下载 git 数据源到临时目录
- 将需要转换的 Codex sub-agent markdown 转换为 TOML

## 维护规则

Go 安装器是独立实现。它不得导入、复制、运行或依赖 `utils/ferryPilot` 下的旧 Python 实现。

