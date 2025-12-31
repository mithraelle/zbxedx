#!/usr/bin/env bash
set -euo pipefail

ZBX_VERSION="${1:-7.4}"
ZBX_BRANCH="release/${ZBX_VERSION}"
ZBX_DIR="zabbix"

echo "==> Bootstrapping Zabbix ${ZBX_BRANCH}"

if [[ -d "${ZBX_DIR}" ]]; then
  echo "Zabbix directory already exists: ${ZBX_DIR}"
  echo "Remove it first if you want to re-bootstrap"
  exit 1
fi

git clone \
  --depth 1 \
  --branch "${ZBX_BRANCH}" \
  https://github.com/zabbix/zabbix.git \
  "${ZBX_DIR}"

echo "==> Updating go.mod"

go mod edit -require=golang.zabbix.com/agent2@v0.0.0
go mod edit -replace=golang.zabbix.com/agent2=./zabbix/src/go
go mod tidy

echo "==> Done"