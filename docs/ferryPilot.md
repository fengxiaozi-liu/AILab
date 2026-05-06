# ferryPilot Go 安装器

`ferryPilot` 用于把 `AISupport/` 中的 package 安装到受支持的本地 Agent target。

## 源码位置

```text
utils/ferryPilotGo/
```

## 构建

```bash
cd utils/ferryPilotGo
go build -o bin/ferryPilot ./cmd/ferryPilot
```

Windows 本地构建可以使用：

```powershell
cd utils/ferryPilotGo
go build -o bin\ferryPilot.exe .\cmd\ferryPilot
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

如果没有提供 package 名称，并且存在多个可选 package，`ferryPilot` 会提示用户选择。

## 资产处理

`ferryPilot` 读取 `config/file_map.json` 来获取 git 数据源和安装映射。运行时会将配置中的 git 仓库 clone 到临时目录，然后从 `tmp/AISupport/<package>` 安装到目标目录。安装完成后，临时 checkout 会被清理。

## Target 行为

默认配置支持：

- `codex`
- `cursor`
- `claude`
- `copilot`
- `gemini`

其中映射为需要转换的 `sub-agents/*.md` 内容会输出为 `.toml` 文件。

