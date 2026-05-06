import argparse
import sys
from pathlib import Path

import questionary

from core import codepilot
from core import codepilot_enum as enum
from core import self_install
from version import VERSION


SKIP_CHOICE = "Skip"


def build_parser() -> argparse.ArgumentParser:
    parser = argparse.ArgumentParser(
        prog="codepilot",
        description="Install AI support skills and agents from AISupport.",
        formatter_class=argparse.RawDescriptionHelpFormatter,
        epilog=(
            "Examples:\n"
            "  codepilot -g              Install an AISupport package globally\n"
            "  codepilot -p              Install an AISupport package into the current project\n"
            "  codepilot -g -t codex     Install globally for Codex\n"
            "  codepilot -p -t cursor    Install into the current project for Cursor"
        ),
    )
    parser.add_argument(
        "-v", "--version",
        action="version",
        version=f"{VERSION}",
        help="Show the current version.",
    )
    parser.add_argument(
        "-g", "--global",
        dest="install_global",
        action="store_true",
        help="Install an AISupport package into the current user's home directory.",
    )
    parser.add_argument(
        "-p", "--project",
        dest="install_project",
        action="store_true",
        help="Install an AISupport package into the current project directory.",
    )
    parser.add_argument(
        "-t", "--target",
        default=None,
        choices=[enum.CodexAgent, enum.GeminiAgent, enum.CopilotAgent, enum.CursorAgent, enum.ClaudeAgent],
        help="Target AI agent. If omitted, an interactive selector is shown.",
    )
    return parser


def select_target() -> str | None:
    return questionary.select(
        "Select target agent:",
        choices=[
            enum.CodexAgent,
            enum.GeminiAgent,
            enum.CopilotAgent,
            enum.CursorAgent,
            enum.ClaudeAgent,
        ],
    ).ask()


def select_support_project(projects: list[str]) -> str | None:
    return questionary.select(
        "Select AISupport package:",
        choices=[*projects, questionary.Separator(), SKIP_CHOICE],
    ).ask()


def confirm_install(target: str, mode: str, package_name: str, install_root: Path) -> bool:
    print("\n--- Install confirmation ---")
    print(f"  Target agent : {target}")
    print(f"  Scope        : {mode}")
    print(f"  Package      : {package_name}")
    print(f"  Install root : {install_root}")
    print("----------------------------")
    return questionary.confirm("Confirm install?", default=True).ask() or False


def install_selected_package(
    pilot: codepilot.CodePilot,
    target: str,
    mode: str,
    install_root: Path,
    config_type: str,
) -> None:
    projects = pilot.get_support_projects()
    if not projects:
        print(f"No installable AISupport packages found in {pilot.ai_support_dir}")
        sys.exit(1)

    package_name = select_support_project(projects)
    if package_name is None or package_name == SKIP_CHOICE:
        print(f"Skipped {mode} install.")
        return

    if not confirm_install(target, mode, package_name, install_root):
        print("Install canceled.")
        return

    if config_type == enum.ConfigGlobal:
        pilot.install_global(package_name, target)
    else:
        pilot.install_project(package_name, target, install_root)


def main():
    parser = build_parser()
    args = parser.parse_args()

    if not args.install_global and not args.install_project:
        if getattr(sys, "frozen", False):
            current = Path(sys.executable).resolve()
            target_path = self_install.get_install_target().resolve()
            if current == target_path:
                parser.print_help()
            else:
                self_install.run()
        else:
            parser.print_help()
        return

    target = args.target or select_target()
    if target is None:
        print("Canceled.")
        sys.exit(0)

    pilot = codepilot.CodePilot()
    if not pilot.has_source_config():
        print(f"AISupport directory not found: {pilot.ai_support_dir}")
        sys.exit(1)

    if args.install_global:
        install_selected_package(
            pilot,
            target,
            "global",
            pilot.home_dir,
            enum.ConfigGlobal,
        )

    if args.install_project:
        install_selected_package(
            pilot,
            target,
            "project",
            Path.cwd(),
            enum.ConfigProject,
        )


if __name__ == "__main__":
    try:
        main()
    except Exception as e:
        print(f"\nError: {e}")
    finally:
        if getattr(sys, "frozen", False) and len(sys.argv) == 1:
            input("\nPress Enter to exit...")
