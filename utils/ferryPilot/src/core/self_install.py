import os
import shutil
import subprocess
import sys
from pathlib import Path


# ── 平台常量 ────────────────────────────────────────────────────────────────

_WIN = sys.platform == "win32"
_INSTALL_DIR = Path.home() / ("bin" if _WIN else ".local/bin")
_EXE_NAME = "codepilot.exe" if _WIN else "codepilot"


# ── 内部工具 ────────────────────────────────────────────────────────────────

def _current_exe() -> Path:
    """返回当前运行的 exe 路径（仅 frozen 模式有效）。"""
    if not getattr(sys, 'frozen', False):
        raise RuntimeError("自安装仅支持打包后的可执行文件")
    return Path(sys.executable)


def _add_to_path_windows(install_dir: Path) -> bool:
    """将 install_dir 追加到 Windows 用户级 PATH（winreg），返回是否有更新。"""
    import winreg  # Windows 标准库，其他平台不会执行到此处

    reg_key = winreg.OpenKey(
        winreg.HKEY_CURRENT_USER,
        r"Environment",
        0,
        winreg.KEY_READ | winreg.KEY_WRITE,
    )
    try:
        current, _ = winreg.QueryValueEx(reg_key, "PATH")
    except FileNotFoundError:
        current = ""

    dirs = [p for p in current.split(";") if p]
    if str(install_dir) in dirs:
        return False  # 已存在，无需更新

    dirs.append(str(install_dir))
    new_path = ";".join(dirs)
    winreg.SetValueEx(reg_key, "PATH", 0, winreg.REG_EXPAND_SZ, new_path)
    winreg.CloseKey(reg_key)

    # 同步更新当前进程的 PATH，使派生的子进程能立即找到 codepilot
    os.environ["PATH"] = str(install_dir) + ";" + os.environ.get("PATH", "")

    # 广播环境变量变更，让当前会话感知（无需注销）
    import ctypes
    ctypes.windll.user32.SendMessageW(0xFFFF, 0x001A, 0, "Environment")
    return True


def _add_to_path_unix(install_dir: Path) -> tuple[Path, bool]:
    """将 export PATH 写入 ~/.bashrc，返回 (配置文件路径, 是否有更新)。"""
    rc = Path.home() / ".bashrc"

    export_line = f'export PATH="{install_dir}:$PATH"'
    content = rc.read_text(encoding="utf-8") if rc.exists() else ""
    if str(install_dir) in content:
        return rc, False  # 已存在

    with open(rc, "a", encoding="utf-8") as f:
        f.write(f"\n{export_line}\n")
    return rc, True


# ── 公开入口 ────────────────────────────────────────────────────────────────

def get_install_target() -> Path:
    """返回 codepilot 可执行文件的安装目标路径。"""
    return _INSTALL_DIR / _EXE_NAME


def run():
    """将 codepilot 安装到 PATH 并配置环境变量。"""
    current_exe = _current_exe()
    target = get_install_target()

    _INSTALL_DIR.mkdir(parents=True, exist_ok=True)
    shutil.copy2(current_exe, target)
    if not _WIN:
        target.chmod(0o755)
    print(f"已安装：{target}")

    if _WIN:
        updated = _add_to_path_windows(_INSTALL_DIR)
        if updated:
            print(f"已将 {_INSTALL_DIR} 添加到用户 PATH")

        # 直接开一个新的 PowerShell 窗口，继承当前进程已更新的 PATH
        # 用户可以在新窗口里立即执行 codepilot -g / codepilot -p
        subprocess.Popen(
            ["powershell", "-NoExit", "-Command",
             f'$env:PATH = "{_INSTALL_DIR};" + $env:PATH; '
             f'Write-Host "codepilot 安装成功，现在可以执行：codepilot -g 或 codepilot -p" -ForegroundColor Green'],
            creationflags=subprocess.CREATE_NEW_CONSOLE,
        )
        print("已打开新终端窗口，可直接使用 codepilot 命令。")
    else:
        rc, updated = _add_to_path_unix(_INSTALL_DIR)
        if updated:
            print(f"已将 {_INSTALL_DIR} 写入 {rc}")
            print(f"执行以下命令立即生效，无需重启终端：")
            print(f"  source {rc} && codepilot -g")
        else:
            print(f"{_INSTALL_DIR} 已在 PATH 中")
            print("现在可以执行：codepilot -g 或 codepilot -p")
