# ferryaide

`ferryaide` 是用于安装本仓库 `AISupport/` 资产的 Go CLI 工具。

## 构建

```bash
make test
make build
```

Windows:

```powershell
make test
make build
```

## 终端启动

从仓库根目录进入 `ferryaide` 工具目录：

```bash
cd utils/ferryaide
```

构建后可以直接从 `bin/` 目录启动本地可执行文件：

```bash
make build
./bin/ferryaide -p speckit
```

Windows PowerShell:

```powershell
cd utils\ferryaide
make build
.\bin\ferryaide -p
```

如果已经通过安装脚本或系统 PATH 安装了 `ferryaide`，可以在任意终端目录直接运行：

```bash
ferryaide -p speckit
```

## 使用

```bash
ferryaide -p speckit
ferryaide -g speckit
ferryaide -p -t codex speckit
ferryaide -p -t cursor speckit
ferryaide -p -t claude speckit
ferryaide -p -t copilot speckit
ferryaide -p -t gemini speckit
```

参数说明：

| 参数 | 说明 |
| --- | --- |
| `-p`, `--project` | 安装到当前项目目录 |
| `-g`, `--global` | 安装到当前用户的全局 Agent 目录 |
| `-t`, `--target` | 指定目标 Agent；省略时会在终端中选择 |
| `--config` | 使用外部 `file_map.json` 覆盖内置默认配置 |

如果省略 `-t/--target`，并且存在多个可选 Agent，`ferryaide` 会先在终端中提供上下键选择；随后如果省略 package 名称，并且存在多个可选 package，也会继续提供上下键选择。

## 配置

默认 `file_map.json` 已通过 Go `embed` 编译进可执行文件。默认配置使用 git 数据源，运行时会 clone 数据源仓库到系统临时目录，从 `AISupport/<package>` 安装文件，完成后清理临时目录。

需要覆盖默认配置时：

```bash
ferryaide -p --config path/to/file_map.json speckit
```
