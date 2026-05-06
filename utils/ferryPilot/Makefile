.PHONY: install package clean release upload upload-direct

# 安装依赖
install:
	pip install -r requirements.txt

# 仅打包为当前平台可执行文件 → release/
package:
	python src/package.py

# 发版（发版人执行）：打包 + 上传 + 创建 git tag + 创建 GitLab Release
# 用法：make release TAG=v1.0.0
release:
	@python -c "import sys; sys.exit(0) if '$(TAG)' else (print('请指定 TAG，例：make release TAG=v1.0.0') or sys.exit(1))"
	python src/package.py --release $(TAG)

# 补充平台产物（非发版人执行）：仅打包 + 上传当前平台二进制，不创建 tag 和 Release
# 用法：make upload TAG=v1.0.0
upload:
	@python -c "import sys; sys.exit(0) if '$(TAG)' else (print('请指定 TAG，例：make upload TAG=v1.0.0') or sys.exit(1))"
	python src/package.py --upload $(TAG)

# 直接上传已有产物到已有 Release，不打包不创建 tag
# 用法：make upload-direct TAG=v1.0.0 FILE=dist/codepilot-linux BIN=codepilot-linux
upload-direct:
	@python -c "import sys; sys.exit(0) if '$(TAG)' else (print('请指定 TAG') or sys.exit(1))"
	@python -c "import sys; sys.exit(0) if '$(BIN)' else (print('请指定 BIN，例：codepilot-linux') or sys.exit(1))"
	python src/package.py --upload-direct "$(TAG)" "$(FILE)" "$(BIN)"

# 清理打包产物
clean:
	python -c "import shutil; [shutil.rmtree(d, ignore_errors=True) for d in ['release', 'build']]"
