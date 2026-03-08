"""
Manage AI runtime assets.
管理 AI 运行时的资产 (例如 prompts, workflows, skills 等)。
"""
from __future__ import annotations

import argparse
import json
import shutil
import subprocess
import sys
import tempfile
from dataclasses import dataclass
from datetime import datetime
from pathlib import Path
from typing import Optional

# --- Constants (常量定义) ---
# Git operations (Git 操作相关常量)
GITLAB_TREE_MARKER = "/-/tree/"
GIT_EXT = ".git"

# Directory names (常用目录与配置文件名)
SHARED_DIR = ".shared"
SHARED_AGENTS_DIR = ".shared/agents"
INSTALL_CONFIG_FILE = "install-config.json"

# Template markers (模板占位符)
RUNTIME_HOME_PLACEHOLDER = "{{runtime_home}}"
TIMESTAMP_PLACEHOLDER = "{{timestamp}}"
PRESERVED_SKILL_DIR_PREFIX = "."


# --- Configuration (配置及目标定义) ---
@dataclass(frozen=True)
class TargetConfig:
    name: str                                     # 目标名称（如 gemini, codex 等）
    source_dir: str                               # 源码所在的子目录
    user_dest_dir: Optional[str] = None           # 用户维度（Home目录）的目标安装路径
    project_dest_dir: Optional[str] = None        # 项目维度的目标安装路径
    generate_project_agents: bool = False         # 是否自动生成项目级的 agent 说明文件
    project_instructions_file: str = "AGENTS.md"  # 项目生成的说明文件名称
    managed_user_items: tuple[str, ...] = ()      # 用户目录下需要独立管理的特定资源集合
    managed_project_items: tuple[str, ...] = ()   # 项目目录下需要独立管理的特定资源集合
    shared_user_items: tuple[str, ...] = ()       # 用户目录下需要共享(从.shared中取)的资源集合
    shared_project_items: tuple[str, ...] = ()    # 项目目录下需要共享(从.shared中取)的资源集合

# 注册支持安装的 AI 环境标的目标
TARGETS = {
    "codex": TargetConfig(
        name="codex",
        source_dir=".codex",
        user_dest_dir=".codex",
        generate_project_agents=True,
        project_instructions_file="AGENTS.md",
        managed_user_items=("AGENTS.md", "prompts", "agents"),
        shared_user_items=("skills",),
    ),
    "github": TargetConfig(
        name="github",
        source_dir=".github",
        project_dest_dir=".github",
        managed_project_items=("agents", "prompts"),
        shared_project_items=("skills", "agents"),
    ),
    "copilot": TargetConfig(
        name="copilot",
        source_dir=".github",
        project_dest_dir=".github",
        generate_project_agents=True,
        project_instructions_file=".github/copilot-instructions.md",
        managed_project_items=("agents", "prompts"),
        shared_project_items=("skills",),
    ),
    "gemini": TargetConfig(
        name="gemini",
        source_dir=".gemini/.agents",
        user_dest_dir=".gemini/antigravity",
        project_dest_dir=".agent",
        generate_project_agents=True,
        project_instructions_file="GEMINI.md",
        managed_user_items=("workflows",),
        shared_user_items=("skills",),
    ),
}

# --- Utility functions (通用工具函数) ---
def get_repo_root() -> Path:
    """Returns the local source repository root. 
    (获取本地源码仓库的绝对根路径)"""
    return Path(__file__).resolve().parent.parent

def generate_timestamp() -> str:
    """Provides a unified timestamp string. 
    (生成统一格式的时间戳字符串，格式例：20231024-123000)"""
    return datetime.now().strftime("%Y%m%d-%H%M%S")

def get_runtime_home_text(runtime_home: Optional[Path], home_dir: Path) -> str:
    """Computes textual representation of the runtime home path. 
    (计算运行时根目录的字符串表示。如果属于 user home 下，则转换为 `~/...` 的简写形式)"""
    if not runtime_home:
        return "."
    try:
        return "~/" + runtime_home.relative_to(home_dir).as_posix()
    except (ValueError, RuntimeError):
        return runtime_home.as_posix()

def render_placeholders(content: str, runtime_text: str) -> str:
    """Replaces format placeholders with runtime values. 
    (替换文本内容中如 {{runtime_home}} 与 {{timestamp}} 等模板占位符)"""
    return content.replace(RUNTIME_HOME_PLACEHOLDER, runtime_text).replace(TIMESTAMP_PLACEHOLDER, generate_timestamp())

def render_text_file(path: Path, runtime_text: str) -> None:
    """Applies placeholder rendering to a text file if required. 
    (如果文件为 Markdown 且内部包含占位符，则进行实时文本渲染)"""
    if path.suffix.lower() != ".md":
        return
    try:
        content = path.read_text(encoding="utf-8")
    except UnicodeDecodeError:
        # 非 utf-8 纯文本将直接忽略报错
        return
    
    # 用简单的字符串检查来提升效率
    if RUNTIME_HOME_PLACEHOLDER not in content and TIMESTAMP_PLACEHOLDER not in content:
        return
    
    rendered = render_placeholders(content, runtime_text)
    if rendered != content:
        path.write_text(rendered, encoding="utf-8")

def robust_remove(path: Path) -> None:
    """Safely removes a file or directory branch. 
    (安全删除。如果是目录就级联删除整树；如果是文件则直接unlink)"""
    if not path.exists():
        return
    shutil.rmtree(path) if path.is_dir() else path.unlink()

# --- Git operations (Git 仓库操作) ---
def normalize_repo_url(repo_url: str, branch: Optional[str]) -> tuple[str, Optional[str]]:
    """(将可能包含 GitLab tree 路径的 URL 进行标准化解析，切分出真正的仓库 .git 地址与分支名)"""
    if GITLAB_TREE_MARKER not in repo_url:
        return repo_url, branch
    
    base, tree_branch = repo_url.split(GITLAB_TREE_MARKER, 1)
    url = base.rstrip("/")
    if not url.endswith(GIT_EXT):
        url += GIT_EXT
    return url, branch or tree_branch.strip("/")

def clone_repo(repo_url: str, branch: Optional[str]) -> Path:
    """(克隆远程 Git 仓库的某分区带(采用 depth=1 浅克隆以提速)，并返回临时目录路径)"""
    url, branch = normalize_repo_url(repo_url, branch)
    checkout_dir = Path(tempfile.mkdtemp(prefix="runtime-install-")) / "source"
    
    cmd = ["git", "clone", "--depth", "1"]
    if branch:
        cmd.extend(["--branch", branch])
    cmd.extend([url, str(checkout_dir)])
    
    subprocess.run(cmd, check=True)
    return checkout_dir

# --- Asset Management (资产资源同步与管理) ---
class AssetSyncManager:
    """Handles file syncing and template formatting operations. 
    (处理文件系统拷贝同步以及模板格式化写入操作的管理类)"""
    
    def __init__(self, allow_overwrite: bool, runtime_text: Optional[str] = None):
        self.allow_overwrite = allow_overwrite
        self.runtime_text = runtime_text

    def copy_file(self, src: Path, dst: Path) -> None:
        """(拷贝单文件。目标如果有冲突根据策略覆盖。如果是带占位符的 md 文档顺手打入渲染内容)"""
        if dst.exists():
            if not self.allow_overwrite:
                raise FileExistsError(f"Destination exists: {dst}")
            robust_remove(dst)
            
        dst.parent.mkdir(parents=True, exist_ok=True)
        if self.runtime_text and src.suffix.lower() == ".md":
            try:
                content = src.read_text(encoding="utf-8")
                if RUNTIME_HOME_PLACEHOLDER in content or TIMESTAMP_PLACEHOLDER in content:
                    dst.write_text(render_placeholders(content, self.runtime_text), encoding="utf-8")
                    return
            except UnicodeDecodeError:
                pass
                
        shutil.copy2(src, dst)

    def copy_tree(self, src: Path, dst: Path) -> None:
        """(拷贝整棵目录树。对于树下带有占位符的 markdown 文件批量应用渲染替换)"""
        if dst.exists():
            if not self.allow_overwrite:
                raise FileExistsError(f"Destination exists: {dst}")
            robust_remove(dst)
            
        dst.parent.mkdir(parents=True, exist_ok=True)
        shutil.copytree(src, dst)
        
        # 递归遍历文件寻找需要渲染的配置文件
        if self.runtime_text:
            for item in dst.rglob("*.md"):
                render_text_file(item, self.runtime_text)

    def copy_skills_tree(self, src: Path, dst: Path) -> None:
        """(拷贝 skills 目录时保留目标目录下的系统技能目录，如 .system)"""
        preserved_items: list[tuple[Path, Path]] = []
        if dst.exists():
            for item in dst.iterdir():
                if not item.name.startswith(PRESERVED_SKILL_DIR_PREFIX):
                    continue

                backup_path = Path(tempfile.mkdtemp(prefix="skills-preserve-")) / item.name
                shutil.copytree(item, backup_path)
                preserved_items.append((backup_path, dst / item.name))

        self.copy_tree(src, dst)

        for backup_path, restore_path in preserved_items:
            if restore_path.exists():
                robust_remove(restore_path)
            shutil.copytree(backup_path, restore_path)
            robust_remove(backup_path.parent)

    def sync_items(self, src_dir: Path, dst_dir: Path, items: tuple[str, ...]) -> None:
        """(按照给定的清单按需挑选目录中的指定项进行拷贝同步)"""
        dst_dir.mkdir(parents=True, exist_ok=True)
        for name in items:
            src = src_dir / name
            if not src.exists():
                continue
                
            dst = dst_dir / name
            if src.resolve() == dst.resolve():
                continue
                
            if src.is_dir():
                if name == "skills":
                    self.copy_skills_tree(src, dst)
                else:
                    self.copy_tree(src, dst)
                if name == "skills":
                    # 对特别的 skills 目录执行后处理过滤逻辑
                    prune_unconfigured_skills(dst)
            else:
                self.copy_file(src, dst)

def prune_unconfigured_skills(skills_dir: Path) -> None:
    """(清理未配置或者被排除的 skills。取决于 install-config.json 中的白名单与黑名单规则)"""
    config_paths = (
        get_repo_root() / "scripts" / INSTALL_CONFIG_FILE,
        skills_dir / INSTALL_CONFIG_FILE,
    )
    
    config_file = next((p for p in config_paths if p.exists()), None)
    installable_set, excluded_set = set(), set()
    installable_skills = None
    
    if config_file:
        data = json.loads(config_file.read_text(encoding="utf-8"))
        installable_skills = data.get("installable_skills")
        if installable_skills is not None:
            installable_set = set(installable_skills)
        excluded_set = set(data.get("excluded_skills", []))
        
    for item in skills_dir.iterdir():
        if not item.is_dir():
            continue
        if item.name.startswith(PRESERVED_SKILL_DIR_PREFIX):
            continue
            
        # 不在 installable 白名单 或者 在 excluded 黑名单，则踢出目标项
        is_unconfigured = config_file and (installable_skills is not None and item.name not in installable_set)
        if is_unconfigured or item.name in excluded_set:
            robust_remove(item)

# --- Primary installation flows (主要安装与卸载流程) ---
class Installer:
    """(负责核心资产资源的分发与清理)"""
    def __init__(self, source_root: Path, home_dir: Path, project_dir: Optional[Path], config: TargetConfig):
        self.source_root = source_root      # 源码根目录
        self.home_dir = home_dir            # 用户目录
        self.project_dir = project_dir      # 项目目录
        self.config = config                # 当前执行的目标配置

    def _setup_target_dir(
        self, 
        dest_rel_path: str,
        managed_items: tuple[str, ...],
        shared_items: tuple[str, ...],
        allow_overwrite: bool
    ) -> None:
        """Helper to sync or copy a generic directory. 
        (通用助手方法：向对应的目标路径做资源灌注(全量copy树 或 使用名单同步))"""
        # 判断安装去哪：用户目录 还是 项目级目录
        dest_dir = self.home_dir / dest_rel_path if dest_rel_path == self.config.user_dest_dir else self.project_dir / dest_rel_path  # type: ignore
        src_dir = self.source_root / self.config.source_dir
        
        runtime_text = get_runtime_home_text(dest_dir, self.home_dir)
        mgr = AssetSyncManager(allow_overwrite, runtime_text)

        # 依次处理需要管理同步的项目 和 全局共享项目
        if managed_items:
            mgr.sync_items(src_dir, dest_dir, managed_items)
        if shared_items:
            mgr.sync_items(self.source_root / SHARED_DIR, dest_dir, shared_items)
            
        # 没有特指管理清单时，走全量目录覆写模式
        if not managed_items and not shared_items:
            mgr.copy_tree(src_dir, dest_dir)
            

    def install(self, allow_overwrite: bool, agent: str) -> None:
        """(执行安装管道的主方法)"""
        # User assets setup (第一阶段：用户维度的资产配置)
        if self.config.user_dest_dir:
            self._setup_target_dir(
                self.config.user_dest_dir,
                self.config.managed_user_items,
                self.config.shared_user_items,
                allow_overwrite
            )

        # Project assets setup (第二阶段：项目维度的资产配置)
        if self.project_dir:
            if self.config.project_dest_dir:
                self._setup_target_dir(
                    self.config.project_dest_dir,
                    self.config.managed_project_items,
                    self.config.shared_project_items,
                    allow_overwrite
                )
                
            # 项目维度特需的 agent 总入口配置生成
            if self.config.generate_project_agents:
                self._generate_project_agents(allow_overwrite, agent)

    def _generate_project_agents(self, allow_overwrite: bool, agent: str) -> None:
        """(从 Agent 模板中生成当前项目的智能体说明文件指令)"""
        template_path = self.source_root / SHARED_AGENTS_DIR / f"{agent}.md"
        if not template_path.exists():
            raise FileNotFoundError(f"Agent template missing: {template_path}")
            
        runtime_home = self.home_dir / self.config.user_dest_dir if self.config.user_dest_dir else Path(self.config.project_dest_dir) if self.config.project_dest_dir else None
        runtime_text = get_runtime_home_text(runtime_home, self.home_dir)
        
        # 植入占位符
        content = render_placeholders(template_path.read_text(encoding="utf-8"), runtime_text)
        
        if self.project_dir:
             dst = self.project_dir / self.config.project_instructions_file
             if dst.exists() and not allow_overwrite:
                 raise FileExistsError(f"Exists: {dst}")
                 
             robust_remove(dst)
             dst.parent.mkdir(parents=True, exist_ok=True)
             dst.write_text(content, encoding="utf-8")

    def delete(self) -> None:
        """(执行卸载与清理流：把之前按照规则安装在各自目录的资源干干净净删掉)"""
        # Clean user assets (清理用户级别)
        if self.config.user_dest_dir:
            user_dest = self.home_dir / self.config.user_dest_dir
            if not self.config.managed_user_items and not self.config.shared_user_items:
                robust_remove(user_dest)  # 全部清理
            else:
                for item in self.config.managed_user_items + self.config.shared_user_items:
                    robust_remove(user_dest / item)  # 枚举清理

        # Clean project assets (清理项目级别)
        if self.project_dir:
            # 清理项目环境里的说明文件
            if self.config.generate_project_agents:
                robust_remove(self.project_dir / self.config.project_instructions_file)
                
            if self.config.project_dest_dir:
                project_dest = self.project_dir / self.config.project_dest_dir
                if not self.config.managed_project_items and not self.config.shared_project_items:
                    robust_remove(project_dest)  # 全部清理
                else:
                    for item in self.config.managed_project_items + self.config.shared_project_items:
                        robust_remove(project_dest / item)  # 枚举清理

# --- Argument Parsing & Bootstrap (参数解析与主例程引导) ---
def parse_args() -> argparse.Namespace:
    parser = argparse.ArgumentParser(description="Manage AI runtime assets. / 管理 AI 运行时相关资产。")
    subparsers = parser.add_subparsers(dest="command", required=True)

    # 注册 install 和 delete 指令
    for cmd in ("install", "delete"):
        sub = subparsers.add_parser(cmd)
        sub.add_argument("--target", choices=sorted(TARGETS.keys()), required=True, help="Target runtime. / 目标环境标识。")
        sub.add_argument("--source-root", help="Local source repository root. / 本地源码仓库路径覆盖。")
        sub.add_argument("--repo-url", help="Remote git repository URL. / 远端 Git 地址，有此项则意味着走在线克隆。")
        sub.add_argument("--branch", help="Remote git branch. / 分支指定。")
        sub.add_argument("--project-dir", help="Project override. / 项目路径覆盖指定。")
        sub.add_argument("--home", help="User home directory override. / 用户主目录覆盖指定。")
        sub.add_argument("--agent", default="kratos", help="Agent template. / Agent 基础模板选择方案配置。")
        
    return parser.parse_args()

def resolve_paths(args: argparse.Namespace, config: TargetConfig) -> tuple[Path, Optional[Path], Path, Optional[Path]]:
    """(解析命令行意图，计算出所有安装交互需要打交道的路径)"""
    temp_root = None
    # 远端模式下产生临时 checkout 源码区
    if args.repo_url:
        source_root = clone_repo(args.repo_url, args.branch)
        temp_root = source_root.parent
    else:
        source_root = Path(args.source_root).resolve() if args.source_root else get_repo_root()

    # 展开 ~ 为对应的真实用户目录
    home_dir = Path(args.home).expanduser().resolve() if args.home else Path.home().resolve()
    
    project_dir = None
    # 若有要求强改 project_dir 取命令行，如果明确在配置要求生成 project 相关的内容，则捕获 cwd
    if args.project_dir:
        project_dir = Path(args.project_dir).expanduser().resolve()
    elif config.generate_project_agents or config.project_dest_dir:
        project_dir = Path.cwd().resolve()
        
    return source_root, temp_root, home_dir, project_dir

def validate_paths(source_root: Path, project_dir: Optional[Path], config: TargetConfig) -> None:
    """(安全边界校验防呆：如果发现目标路径和源码是同一位置就阻断写入以防自相残杀销毁源码)"""
    if project_dir and config.project_dest_dir and project_dir == source_root:
        if (source_root / config.source_dir).resolve() == (project_dir / config.project_dest_dir).resolve():
            raise ValueError(f"Cannot overwrite source assets for {config.name}. / 不可将源文件配置原地自覆写！")

def main() -> int:
    args = parse_args()
    config = TARGETS[args.target]
    
    source_root, temp_root, home_dir, project_dir = resolve_paths(args, config)
    validate_paths(source_root, project_dir, config)
    
    installer = Installer(source_root, home_dir, project_dir, config)
    try:
        if args.command == "install":
            installer.install(allow_overwrite=True, agent=args.agent)
        elif args.command == "delete":
            installer.delete()
    finally:
        # 在线拉库的情况下，结束运行时一并清理临时暂存目录
        if temp_root and temp_root.exists():
            shutil.rmtree(temp_root, ignore_errors=True)

    return 0

if __name__ == "__main__":
    try:
        sys.exit(main())
    except subprocess.CalledProcessError as e:
        print(f"Command failed: {' '.join(map(str, e.cmd))}", file=sys.stderr)
        sys.exit(e.returncode)
    except Exception as e:
        print(f"Error: {e}", file=sys.stderr)
        sys.exit(1)
