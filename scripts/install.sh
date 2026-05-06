#!/usr/bin/env sh
set -eu

repo="fengxiaozi-liu/ai_agent_skill"
install_dir="$HOME/.local/bin"

os="$(uname -s)"
arch="$(uname -m)"

case "$os" in
  Linux)
    goos="linux"
    ;;
  Darwin)
    goos="darwin"
    ;;
  *)
    echo "Unsupported OS: $os" >&2
    exit 1
    ;;
esac

case "$arch" in
  x86_64|amd64)
    goarch="amd64"
    ;;
  arm64|aarch64)
    if [ "$goos" = "darwin" ]; then
      goarch="arm64"
    else
      echo "Unsupported Linux architecture: $arch" >&2
      exit 1
    fi
    ;;
  *)
    echo "Unsupported architecture: $arch" >&2
    exit 1
    ;;
esac

asset="ferryPilot-$goos-$goarch"
url="https://github.com/$repo/releases/latest/download/$asset"
target="$install_dir/ferryPilot"

mkdir -p "$install_dir"

if command -v curl >/dev/null 2>&1; then
  curl -fsSL "$url" -o "$target"
elif command -v wget >/dev/null 2>&1; then
  wget -q "$url" -O "$target"
else
  echo "curl or wget is required" >&2
  exit 1
fi

chmod +x "$target"

profile=""
if [ "$goos" = "darwin" ]; then
  profile="$HOME/.zshrc"
elif [ -n "${BASH_VERSION:-}" ]; then
  profile="$HOME/.bashrc"
else
  profile="$HOME/.profile"
fi

case ":$PATH:" in
  *":$install_dir:"*) ;;
  *)
    touch "$profile"
    if ! grep -F 'export PATH="$HOME/.local/bin:$PATH"' "$profile" >/dev/null 2>&1; then
      printf '\nexport PATH="$HOME/.local/bin:$PATH"\n' >> "$profile"
    fi
    ;;
esac

echo "ferryPilot installed to $target"
echo "Restart your terminal, then run: ferryPilot -p speckit"
