# ferryPilot

ferryPilot is a CLI installer for local AI support assets stored in this repository's `AISupport/` directory.

## Project Layout

```text
.
├── AISupport/
│   ├── kratos/
│   │   └── skills/
│   └── speckit/
│       ├── skills/
│       └── sub-agents/
└── utils/
    └── ferryPilot/
```

`AISupport/<package>` is the install unit. For example, selecting `speckit` installs everything under `AISupport/speckit/skills` and `AISupport/speckit/sub-agents`.

## Usage

```bash
# Install one AISupport package globally
codepilot -g

# Install one AISupport package into the current project
codepilot -p

# Install for a specific target agent
codepilot -g -t codex
codepilot -p -t cursor
```

## Behavior

- `-g / --global` installs into the current user's home directory.
- `-p / --project` installs into the current working directory.
- The installer scans only the first-level directories under `AISupport/`, such as `speckit` and `kratos`.
- After selecting a package, all installable content under that package is copied according to `src/config/file_map.json`.
- Codex `sub-agents/*.md` files are converted to `.toml` during installation, preserving the previous installer behavior.

## Build

```bash
make install
make package
```

The PyInstaller package includes both `src/config/file_map.json` and the repository-level `AISupport/` directory, so the released executable can install the bundled support assets directly.
