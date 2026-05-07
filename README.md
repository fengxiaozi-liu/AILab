# AI Agent Skill

本仓库用于分发面向编码 Agent 的 AI 支持资产，并提供 `ferryPilot` 工具把这些资产安装到本地 Agent 运行时或当前项目目录。

## AISupport

`AISupport/` 是资产目录。`AISupport/<package>` 是安装单元，每个 package 可以包含 skills、sub-agents、templates、references 等内容。

## utils

`utils/` 存放仓库工具。当前主要工具是 `utils/ferryPilot/`。

### ferryPilot

`ferryPilot` 是 Go 编写的 CLI 安装器，用于把 `AISupport` 中的 package 安装到目标 Agent。工具介绍、构建命令、使用方式和配置说明请见 [`utils/ferryPilot/README.md`](utils/ferryPilot/README.md)。

#### 一键安装

Windows PowerShell:

```powershell
irm https://raw.githubusercontent.com/fengxiaozi-liu/AILab/main/utils/ferryPilot/scripts/install.ps1 | iex
```

Linux/macOS:

```bash
curl -fsSL https://raw.githubusercontent.com/fengxiaozi-liu/AILab/main/utils/ferryPilot/scripts/install.sh | sh
```

安装脚本会下载最新 Release，把 `ferryPilot` 放到用户目录并加入 PATH。重新打开终端后即可使用：
