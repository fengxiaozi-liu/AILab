# 发布流程

`ferryPilot` 的发布二进制由 GitHub Actions 构建。

## 触发方式

推送匹配 `v*` 的 tag：

```bash
git tag v0.1.0
git push origin v0.1.0
```

## 构建矩阵

发布 workflow 会构建：

| OS | Arch | Artifact |
|----|------|----------|
| windows | amd64 | `ferryPilot-windows-amd64.tar.gz` |
| linux | amd64 | `ferryPilot-linux-amd64.tar.gz` |
| darwin | amd64 | `ferryPilot-darwin-amd64.tar.gz` |
| darwin | arm64 | `ferryPilot-darwin-arm64.tar.gz` |

## Workflow 职责

发布 workflow 会：

1. checkout 仓库。
2. 根据 `utils/ferryPilotGo/go.mod` 设置 Go 环境。
3. 执行 `go test ./...`。
4. 交叉编译 `ferryPilot`。
5. 将可执行文件与 `config/file_map.json` 一起打包。
6. 将产物附加到 GitHub Release。

