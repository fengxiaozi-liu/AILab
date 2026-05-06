# ferryPilot

`ferryPilot` 是一个 Go CLI 安装器，用于安装本仓库 `AISupport/` 目录中的 AI 支持资产。

当前实现独立于 `utils/ferryPilot` 下的旧 Python 实现；它不会导入、运行、复制或依赖这些旧文件。

## 构建

```bash
go build -o bin/ferryPilot ./cmd/ferryPilot
```

Windows 下可以使用：

```powershell
go build -o bin\ferryPilot.exe .\cmd\ferryPilot
```

## 测试

```bash
go test ./...
```

## 使用方式

```bash
ferryPilot -g speckit
ferryPilot -p speckit
ferryPilot -g -t codex speckit
ferryPilot -p -t cursor speckit
ferryPilot -p -t copilot speckit
ferryPilot -g -t claude speckit
ferryPilot -g -t gemini speckit
```

参数说明：

- `-g`, `--global`：安装到当前用户的 Agent 运行时目录。
- `-p`, `--project`：安装到当前项目目录。
- `-t`, `--target`：选择目标 Agent，默认值为 `codex`。

支持的 target 在 `config/file_map.json` 中配置。默认配置包含 `codex`、`cursor`、`claude`、`copilot` 和 `gemini`。

如果省略 package 名称，并且存在多个可选 package，`ferryPilot` 会提示用户选择。

## 数据源

`ferryPilot` 读取 `config/file_map.json` 来确定：

- 作为 AISupport 数据源的 git 仓库
- 默认 target
- 每个 target 的全局安装映射和项目安装映射

运行时会将配置中的 git 仓库 clone 到临时目录，然后从 `tmp/AISupport/<package>` 安装文件到所选目标目录。安装完成后，临时 checkout 会被清理。

## 发布

当推送匹配 `v*` 的 tag 时，GitHub Actions 会构建发布产物。支持的目标平台包括：

- `windows/amd64`
- `linux/amd64`
- `darwin/amd64`
- `darwin/arm64`

