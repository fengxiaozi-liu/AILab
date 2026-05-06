import shutil
import subprocess
import sys
from pathlib import Path
import requests
from urllib.parse import quote, urlparse

from dataclasses import dataclass, field, asdict
import re

PROJECT_ROOT = Path(__file__).parent
ENTRY_POINT = PROJECT_ROOT / "command.py"
FILE_MAP = PROJECT_ROOT / "config" / "file_map.json"
AI_SUPPORT = PROJECT_ROOT.parent.parent.parent / "AISupport"
RELEASE_DIR = PROJECT_ROOT.parent / "release"
BUILD_DIR = PROJECT_ROOT.parent / "build"
EXE_NAME = "codepilot"

DEFAULT_GITLAB_URL = "http://10.0.0.191"
DEFAULT_GITLAB_PROJECT_PATH = "TeamA/AILab/Coding/Utils/CodePilot"
PACKAGE_NAME = EXE_NAME

type PLATFORMLabel = str
WindowsLabel = "Windows"
MacLabel = "MacOS"
LinuxLabel = "Linux"

# 各平台信息
_PLATFORM_INFO = {
    "win32": {"label": WindowsLabel, "bin": f"{EXE_NAME}.exe", "sep": ";"},
    "darwin": {"label": MacLabel, "bin": EXE_NAME, "sep": ":"},
    "linux": {"label": LinuxLabel, "bin": EXE_NAME, "sep": ":"},
}
_platform = _PLATFORM_INFO.get(sys.platform) or _PLATFORM_INFO["linux"]

_PLATFORM_RELEASE_BIN = {
    "win32": f"{EXE_NAME}.exe",
    "darwin": f"{EXE_NAME}-macos",
    "linux": f"{EXE_NAME}-linux",
}

type PlatformReleaseName = str
WinReleaseName = "codepilot.exe (Windows)"
MacReleaseName = "codepilot (macOS)"
LinuxReleaseName = "codepilot (Linux)"
ReleaseNameMap: dict[str, str] = {
    WinReleaseName: "codepilot.exe",
    MacReleaseName: "codepilot-macos",
    LinuxReleaseName: "codepilot-linux"
}


def find_gitlab_project_info() -> tuple[str, str]:
    """Resolve the GitLab host and project path from origin remote."""
    r = subprocess.run(
        ["git", "config", "--get", "remote.origin.url"],
        capture_output=True, text=True, cwd=PROJECT_ROOT.parent,
    )
    remote_url = r.stdout.strip()
    if not remote_url:
        return DEFAULT_GITLAB_URL, DEFAULT_GITLAB_PROJECT_PATH

    if remote_url.startswith(("http://", "https://")):
        parsed = urlparse(remote_url)
        project_path = parsed.path.lstrip("/")
        if project_path.endswith(".git"):
            project_path = project_path[:-4]
        if parsed.scheme and parsed.netloc and project_path:
            return f"{parsed.scheme}://{parsed.netloc}", project_path

    match = re.match(r"git@(?P<host>[^:]+):(?P<path>.+?)(?:\.git)?$", remote_url)
    if match:
        return f"http://{match.group('host')}", match.group("path")

    return DEFAULT_GITLAB_URL, DEFAULT_GITLAB_PROJECT_PATH


GITLAB_URL, GITLAB_PROJECT_PATH = find_gitlab_project_info()
GITLAB_PROJECT = quote(GITLAB_PROJECT_PATH, safe="")


def build_package():
    print(f">>> 当前平台：{_platform['label']}")

    # 清理上次产物
    for d in (RELEASE_DIR, BUILD_DIR):
        if d.exists():
            shutil.rmtree(d)
            print(f"已清理：{d}")

    print(">>> 开始打包...")
    add_data_args = ["--add-data", f"{FILE_MAP}{_platform['sep']}config"]
    if AI_SUPPORT.exists():
        add_data_args += ["--add-data", f"{AI_SUPPORT}{_platform['sep']}AISupport"]
    else:
        print(f">>> AISupport not found, package will not include local skills: {AI_SUPPORT}")

    result = subprocess.run(
        [
            sys.executable, "-m", "PyInstaller",
            "--onefile",
            "--name", EXE_NAME,
            *add_data_args,
            "--distpath", str(RELEASE_DIR),
            "--workpath", str(BUILD_DIR),
            "--noconfirm",
            str(ENTRY_POINT),
        ],
        cwd=PROJECT_ROOT,
    )

    # 清理 build 中间产物和 spec 文件
    if BUILD_DIR.exists():
        shutil.rmtree(BUILD_DIR)
    spec_file = PROJECT_ROOT / f"{EXE_NAME}.spec"
    if spec_file.exists():
        spec_file.unlink()

    if result.returncode != 0:
        print(">>> 打包失败！", file=sys.stderr)
        sys.exit(1)

    exe = RELEASE_DIR / _platform["bin"]
    print(f"\n打包成功：{exe}")


def find_gitlab_token() -> str | None:
    """从已有配置中查找 token，找不到返回 None。"""
    # 1. git config
    r = subprocess.run(
        ["git", "config", "--get", "gitlab.token"],
        capture_output=True, text=True,
    )
    if r.returncode == 0 and r.stdout.strip():
        print(">>> 已从 git config 读取 Token")
        return r.stdout.strip()
    return None


@dataclass(eq=False)
class AssetsBaseInfo:
    name: str
    url: str
    link_type: str = field(default="package")


@dataclass(eq=False)
class AssetsLinks:
    links: list[AssetsBaseInfo] = field(default_factory=list)


@dataclass(eq=False)
class GitlabReleaseInfo:
    name: str
    description: str
    tag_name: str
    assets: AssetsLinks = field(default_factory=AssetsLinks)


@dataclass(eq=False)
class GitLabConfig:
    base_url: str = field(default=f'{GITLAB_URL}/api/v4/projects/{GITLAB_PROJECT}/packages/generic/{PACKAGE_NAME}')
    tag_url: str = field(default="")  # 派生自 base_url，格式：base_url/%s
    link_url: str = field(default="")  # 派生自 release_url，格式：release_url/%s/assets/links
    get_token_url: str = field(default=f"{GITLAB_URL}/api/v4/user")
    release_url: str = field(default=f"{GITLAB_URL}/api/v4/projects/{GITLAB_PROJECT}/releases")
    token_header: dict[str, str] = field(default_factory=lambda: {"PRIVATE-TOKEN": find_gitlab_token() or ""})
    release_description_format: str = field(default=(
        "## 下载\n\n| 平台 | 文件 |\n|------|------|\n"
        "| Windows | [codepilot.exe](%s/codepilot.exe) |\n"
        "| Linux   | [codepilot-linux](%s/codepilot-linux) |\n"
        "| macOS   | [codepilot-macos](%s/codepilot-macos) |"
    ))

    def __post_init__(self):
        if not self.tag_url:
            self.tag_url = self.base_url + "/{tag}"
        if not self.link_url:
            self.link_url = self.release_url + "/{tag}/assets/links"


@dataclass(eq=False)
class PackageUtil:
    gitlabConfig: GitLabConfig = field(default_factory=GitLabConfig)

    def release(self, tag: str):
        """发版：打包 + 上传 + 创建 git tag + 创建 GitLab Release。发版人执行一次。"""
        self.refresh_token()
        self.update_version(tag)
        self.build_and_upload(tag)

        # 推送 tag
        result = subprocess.run(["git", "tag", tag], cwd=PROJECT_ROOT.parent, capture_output=True, text=True)
        if result.returncode != 0 and "already exists" not in result.stderr:
            print(result.stderr, file=sys.stderr)
            sys.exit(1)
        subprocess.run(["git", "push", "origin", tag], check=True, cwd=PROJECT_ROOT.parent)
        print(f">>> tag {tag} 已推送")

        # 创建 Release
        self.create_release(tag)
        print(">>> 完成！")

    def upload(self, tag: str):
        """补充平台产物：打包 + 上传当前平台二进制 + 更新 Release 中对应平台的下载链接，不创建 tag 和 Release。"""
        self.refresh_token()
        release_bin = self.build_and_upload(tag)
        find_name = {"win32": WinReleaseName, "darwin": MacReleaseName}.get(sys.platform, LinuxReleaseName)
        release_info = self.build_release_payload(tag)
        self.update_release_link(tag, release_bin, find_name, release_info.description)
        print(">>> 完成！")

    def upload_direct(self, tag: str, file_path: str, release_bin: str):
        """直接上传已有文件到 GitLab，更新 Release 对应平台下载链接与 Notes。不打包，不创建 tag 和 Release。"""
        self.refresh_token()
        source = Path(file_path).expanduser().resolve()
        if not source.exists() or not source.is_file():
            raise FileNotFoundError(f"文件不存在：{source}")
        find_name = {
            _PLATFORM_RELEASE_BIN["win32"]: WinReleaseName,
            _PLATFORM_RELEASE_BIN["darwin"]: MacReleaseName,
            _PLATFORM_RELEASE_BIN["linux"]: LinuxReleaseName,
        }.get(release_bin, release_bin)
        self.upload_bin(tag, source, release_bin)
        release_info = self.build_release_payload(tag)
        self.update_release_link(tag, release_bin, find_name, release_info.description)
        print(">>> 完成！")

    def refresh_token(self):
        token = self.gitlabConfig.token_header["PRIVATE-TOKEN"]
        if token and token != "":
            return
        token = input("请输入 Personal Access Token：").strip()
        if token == "":
            return
        self.gitlabConfig.token_header["PRIVATE-TOKEN"] = token
        try:
            r = requests.get(self.gitlabConfig.get_token_url, headers=self.gitlabConfig.token_header)
            if r.status_code == 401:
                print(">>> Token 无效，终止。", file=sys.stderr)
                return
        except Exception:
            print(">>> Token 无效，终止。", file=sys.stderr)
            return
        subprocess.run(["git", "config", "--global", "gitlab.token", token], check=True)
        print(">>> Token 验证通过，已保存到 git config --global gitlab.token")
        return

    def update_version(self, tag: str):
        """将 tag（如 v1.2.0）写入 version.py 的 VERSION 常量。"""
        version_file = PROJECT_ROOT / "version.py"
        content = version_file.read_text(encoding="utf-8")
        new_content = re.sub(
            r'^VERSION\s*=\s*"[^"]*"',
            f'VERSION = "{tag}"',
            content,
            flags=re.MULTILINE,
        )
        version_file.write_text(new_content, encoding="utf-8")

    def update_release_link(self, tag: str, release_bin: str, find_name: str, description: str = "") -> None:
        """更新 GitLab Release 中指定平台对应的 asset link，可同时更新 description。"""
        headers = self.gitlabConfig.token_header
        new_url = self.gitlabConfig.tag_url.format(tag=tag) + f"/{release_bin}"
        links_url = self.gitlabConfig.link_url.format(tag=tag)

        resp = requests.get(links_url, headers=headers)
        if not resp.ok:
            print(f">>> 无法获取 Release {tag} 的 links，跳过更新: {resp.text}", file=sys.stderr)
            return
        target = next((lk for lk in resp.json() if find_name in lk.get("name", "")), None)

        if target:
            asset_base_info = AssetsBaseInfo(name=find_name, url=new_url)
            r = requests.put(f"{links_url}/{target['id']}", headers=headers, json=asdict(asset_base_info))
            if r.ok:
                print(f">>> Release {tag} 的 {find_name} 下载链接已更新：{new_url}")
            else:
                print(f">>> 更新 Release link 失败: {r.text}", file=sys.stderr)
        else:
            print(f">>> 在 Release {tag} 中未找到 {find_name} 的下载链接，无法更新")

        # 同时更新 Release Notes（description）
        if description:
            r = requests.put(
                f"{self.gitlabConfig.release_url}/{tag}",
                headers=headers,
                json={"description": description},
            )
            if r.ok:
                print(f">>> Release {tag} 的 Notes 已更新")
            else:
                print(f">>> 更新 Release Notes 失败: {r.text}", file=sys.stderr)

    def build_and_upload(self, tag: str) -> str:
        """打包 + 重命名 + 上传当前平台二进制到 Package Registry，返回 release_bin 名。"""
        build_package()
        raw = RELEASE_DIR / _platform["bin"]
        release_bin = _PLATFORM_RELEASE_BIN[sys.platform]
        dst = RELEASE_DIR / release_bin
        if raw != dst:
            shutil.copy2(raw, dst)
        self.upload_bin(tag, dst, release_bin)
        return release_bin

    def upload_bin(self, tag: str, file_path: Path, release_bin: str) -> None:
        """上传指定文件到 GitLab Package Registry。"""
        bin_url = self.gitlabConfig.tag_url.format(tag=tag) + f"/{release_bin}"
        print(f"\n>>> 上传 {file_path.name} → {bin_url}")
        with open(file_path, "rb") as f:
            resp = requests.put(bin_url, headers=self.gitlabConfig.token_header, data=f)
        resp.raise_for_status()
        print(">>> 上传成功")

    def build_release_payload(self, tag: str) -> GitlabReleaseInfo:
        """构建 GitLab Release 的请求体（description + assets links）。"""
        tag_url = self.gitlabConfig.tag_url.format(tag=tag)
        release_name: str = f'Release {tag}'
        tag_name: str = tag
        description: str = self.gitlabConfig.release_description_format % (tag_url, tag_url, tag_url)
        links: list[AssetsBaseInfo] = []
        for link_name, value in ReleaseNameMap.items():
            asset_base = AssetsBaseInfo(name=link_name, url=f"{tag_url}/{value}")
            links.append(asset_base)
        asset_link = AssetsLinks(links)
        return GitlabReleaseInfo(release_name, description, tag_name, asset_link)

    def create_release(self, tag: str) -> None:
        """创建 GitLab Release，已存在则跳过。"""
        releases_url = self.gitlabConfig.release_url
        payload = self.build_release_payload(tag)
        exists = requests.get(f"{releases_url}/{tag}", headers=self.gitlabConfig.token_header).ok
        if exists:
            print(f">>> Release {tag} 已存在，跳过创建（可用 make upload 补充当前平台产物）")
            return
        requests.post(releases_url, headers=self.gitlabConfig.token_header, json=asdict(payload)).raise_for_status()
        print(f">>> Release {tag} 创建成功：{GITLAB_URL}/{GITLAB_PROJECT_PATH}/-/releases")


if __name__ == "__main__":
    if len(sys.argv) >= 2 and sys.argv[1] == "--release":
        if len(sys.argv) < 3:
            print("用法：python src/package.py --release <tag>")
            sys.exit(1)
        PackageUtil().release(tag=sys.argv[2])
    elif len(sys.argv) >= 2 and sys.argv[1] == "--upload":
        if len(sys.argv) < 3:
            print("用法：python src/package.py --upload <tag>")
            sys.exit(1)
        PackageUtil().upload(tag=sys.argv[2])
    elif len(sys.argv) >= 2 and sys.argv[1] == "--upload-direct":
        if len(sys.argv) < 5:
            print("用法：python src/package.py --upload-direct <tag> <file_path> <release_bin>")
            print("  release_bin 示例：codepilot.exe / codepilot-linux / codepilot-macos")
            sys.exit(1)
        PackageUtil().upload_direct(tag=sys.argv[2], file_path=sys.argv[3], release_bin=sys.argv[4])
    else:
        build_package()
