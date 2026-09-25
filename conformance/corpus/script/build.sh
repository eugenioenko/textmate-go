#!/usr/bin/env bash
set -euo pipefail

project=${PROJECT_NAME:-textmate-go}
message=$(cat <<EOF
building ${project}
with stateful highlighting
EOF
)
printf '%s\n' "$message"
