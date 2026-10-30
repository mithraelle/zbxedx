#!/usr/bin/env bash
set -euo pipefail

ZBX_REPO="https://github.com/zabbix/zabbix.git"
ZBX_DIR="zabbix"
ZBX_DEFAULT_VERSION="7.4"

usage() {
  cat <<EOF
Usage: $(basename "$0") [--force] [VERSION]

Fetches the Zabbix source tree into ./${ZBX_DIR} and points go.mod at it.

VERSION:
  X.Y           latest released X.Y.Z tag (default: ${ZBX_DEFAULT_VERSION})
  X.Y.Z         exact tag, including pre-releases such as 7.4.17rc1
  release/X.Y   tip of the release branch (may contain unreleased changes)

Options:
  -f, --force   replace an existing ${ZBX_DIR} directory
  -h, --help    show this help
EOF
}

FORCE=0
ZBX_VERSION=""

while [[ $# -gt 0 ]]; do
  case "$1" in
    -f|--force) FORCE=1 ;;
    -h|--help) usage; exit 0 ;;
    -*) echo "Unknown option: $1" >&2; usage >&2; exit 1 ;;
    *)
      if [[ -n "${ZBX_VERSION}" ]]; then
        echo "Only one VERSION may be given" >&2
        exit 1
      fi
      ZBX_VERSION="$1"
      ;;
  esac
  shift
done

ZBX_VERSION="${ZBX_VERSION:-${ZBX_DEFAULT_VERSION}}"

cd "$(dirname "$0")"

# Resolve VERSION to a branch or tag before touching the existing tree
if [[ "${ZBX_VERSION}" =~ ^release/[0-9]+\.[0-9]+$ ]]; then
  ZBX_REF="${ZBX_VERSION}"
elif [[ "${ZBX_VERSION}" =~ ^[0-9]+\.[0-9]+\.[0-9]+([a-z]+[0-9]+)?$ ]]; then
  ZBX_REF="${ZBX_VERSION}"
elif [[ "${ZBX_VERSION}" =~ ^[0-9]+\.[0-9]+$ ]]; then
  echo "==> Looking up the latest ${ZBX_VERSION}.x release"
  ZBX_REF="$(git ls-remote --tags --refs "${ZBX_REPO}" "refs/tags/${ZBX_VERSION}.*" \
    | sed 's#.*refs/tags/##' \
    | grep -E "^${ZBX_VERSION//./\\.}\.[0-9]+$" \
    | sort -t. -k3,3n \
    | tail -n 1 || true)"
  if [[ -z "${ZBX_REF}" ]]; then
    echo "No ${ZBX_VERSION}.x release found; use release/${ZBX_VERSION} or an exact tag" >&2
    exit 1
  fi
  echo "==> Resolved ${ZBX_VERSION} to ${ZBX_REF}"
else
  echo "Unrecognized VERSION: ${ZBX_VERSION}" >&2
  usage >&2
  exit 1
fi

if [[ -d "${ZBX_DIR}" ]]; then
  if [[ "${FORCE}" -eq 1 ]]; then
    echo "==> Removing existing ${ZBX_DIR}"
    rm -rf "${ZBX_DIR}"
  else
    echo "Zabbix directory already exists: ${ZBX_DIR}"
    echo "Use --force to replace it"
    exit 1
  fi
fi

echo "==> Bootstrapping Zabbix ${ZBX_REF}"

git -c advice.detachedHead=false clone \
  --depth 1 \
  --branch "${ZBX_REF}" \
  "${ZBX_REPO}" \
  "${ZBX_DIR}"

echo "==> Updating go.mod"

go mod edit -require=golang.zabbix.com/agent2@v0.0.0
go mod edit -replace=golang.zabbix.com/agent2=./zabbix/src/go
go mod tidy

echo "==> Done: Zabbix ${ZBX_REF} ($(git -C "${ZBX_DIR}" rev-parse --short HEAD))"
