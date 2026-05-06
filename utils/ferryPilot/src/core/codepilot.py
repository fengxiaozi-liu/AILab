import json
import shutil
import sys
from dataclasses import dataclass, field
from pathlib import Path

import frontmatter
import tomli_w

from core import codepilot_enum as enum


def _get_base_dir() -> Path:
    if getattr(sys, "frozen", False):
        return Path(getattr(sys, "_MEIPASS", ""))
    return Path(__file__).parent


def _get_repo_root() -> Path:
    if getattr(sys, "frozen", False):
        return Path(getattr(sys, "_MEIPASS", ""))
    return Path(__file__).resolve().parents[4]


@dataclass(eq=False)
class TargetConfig:
    name: str = ""
    dirs: dict[str, str] = field(default_factory=dict)


def copy_dir(source_dir: Path, target_dir: Path) -> bool:
    try:
        if not source_dir.exists():
            print(f"Source path does not exist: {source_dir}")
            return False
        shutil.copytree(source_dir, target_dir, dirs_exist_ok=True)
        return True
    except Exception as e:
        print(f"Error during copy: {e}")
        return False


def _build_target_config(target_name: str, target_value_dict: dict) -> TargetConfig:
    config = TargetConfig(name=target_name)
    for key, value in target_value_dict.items():
        config.dirs[key] = value
    return config


@dataclass
class CodePilot:
    target: dict[str, dict[str, TargetConfig]] = field(default_factory=dict)
    branch: str = "main"
    home_dir: Path = field(default_factory=Path.home)
    base_dir: Path = field(default_factory=_get_base_dir)
    repo_root: Path = field(default_factory=_get_repo_root)

    def __post_init__(self):
        if getattr(sys, "frozen", False):
            config_path = self.base_dir / enum.FILE_MAP_FROZEN_PATH
        else:
            config_path = self.base_dir / enum.FILE_MAP_PATH

        with open(config_path, encoding="utf-8") as file:
            data = json.load(file)

        for config_type, config_value in data.items():
            if config_type not in enum.ConfigTypeList:
                print(f"{config_type} is not a valid config type")
                continue

            self.target[config_type] = {}
            target_section = config_value.get(enum.TargetSectionKey, {})
            for target_name, target_value_dict in target_section.items():
                if target_name not in enum.TargetAgentList:
                    print(f"{target_name} is not a valid target agent")
                    continue
                self.target[config_type][target_name] = _build_target_config(
                    target_name,
                    target_value_dict,
                )

    @property
    def ai_support_dir(self) -> Path:
        return self.repo_root / "AISupport"

    def has_source_config(self) -> bool:
        return self.ai_support_dir.is_dir()

    def get_support_projects(self) -> list[str]:
        if not self.ai_support_dir.is_dir():
            return []
        return sorted(
            item.name
            for item in self.ai_support_dir.iterdir()
            if item.is_dir() and not item.name.startswith(".")
        )

    def get_target_dir_config(self, config_type: enum.ConfigType, target: enum.TargetAgent) -> TargetConfig | None:
        return self.target.get(config_type, {}).get(target)

    def install_support_project(
        self,
        project_name: str,
        target: enum.TargetAgent,
        install_root: Path,
        config_type: enum.ConfigType,
    ) -> None:
        source_root = self.ai_support_dir / project_name
        if not source_root.is_dir():
            print(f"{project_name}: source project not found in AISupport")
            return

        target_config = self.get_target_dir_config(config_type, target)
        if target_config is None or not target_config.dirs:
            print(f"{target}: no target directories configured for {config_type}")
            return

        copied_any = False
        copied_any |= self._install_key(source_root, target_config, target, install_root, "skills")
        copied_any |= self._install_key(source_root, target_config, target, install_root, enum.SubAgentsKey, "sub-agents")

        if copied_any:
            print(f"{config_type} install completed: {project_name} -> {target}")
        else:
            print(f"{project_name}: no installable skills or sub-agents found")

    def install_global(self, project_name: str, target: enum.TargetAgent):
        self.install_support_project(project_name, target, self.home_dir, str(enum.ConfigGlobal))

    def install_project(self, project_name: str, target: enum.TargetAgent, install_root: Path):
        self.install_support_project(project_name, target, install_root, str(enum.ConfigProject))

    def _install_key(
        self,
        source_root: Path,
        target_config: TargetConfig,
        target: enum.TargetAgent,
        install_root: Path,
        target_key: str,
        source_dir_name: str | None = None,
    ) -> bool:
        source_dir = source_root / (source_dir_name or target_key)
        if not source_dir.exists():
            return False

        target_rel_path = target_config.dirs.get(target_key)
        if not target_rel_path:
            print(f"{source_root.name}: key '{target_key}' not found in target {target}, skip")
            return False

        full_target = install_root / target_rel_path
        if target == enum.CodexAgent and target_key == enum.SubAgentsKey:
            return copy_dir_as_toml(source_dir, full_target)
        return copy_dir(source_dir, full_target)


def copy_dir_as_toml(source_dir: Path, target_dir: Path) -> bool:
    try:
        if not source_dir.exists():
            print(f"Source path does not exist: {source_dir}")
            return False
        target_dir.mkdir(parents=True, exist_ok=True)
        for md_file in source_dir.glob("*.md"):
            toml_content = convert_md_to_toml(md_file)
            target_file = target_dir / f"{md_file.stem}.toml"
            target_file.write_text(toml_content, encoding="utf-8")
        return True
    except Exception as e:
        print(f"Error during md-to-toml conversion: {e}")
        return False


def convert_md_to_toml(md_file: Path) -> str:
    post = frontmatter.load(str(md_file))
    base = tomli_w.dumps(dict(post.metadata))
    body = post.content.strip()
    return f'{base}developer_instructions = """\n{body}\n"""\n'


if __name__ == "__main__":
    pilot = CodePilot()
    print(pilot.get_support_projects())
