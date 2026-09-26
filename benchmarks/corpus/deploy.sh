#!/usr/bin/env bash
# Build, tag, and deploy a release to the configured environment.
set -euo pipefail

readonly SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
readonly ENVIRONMENTS=(staging production)
LOG_FILE="${LOG_FILE:-/tmp/deploy-$(date +%Y%m%d).log}"
DRY_RUN=0

log() {
  local level="$1"; shift
  printf '%s [%s] %s\n' "$(date -Iseconds)" "$level" "$*" | tee -a "$LOG_FILE" >&2
}

usage() {
  cat <<USAGE
Usage: $(basename "$0") [-n] <environment> <version>

  -n   dry run; print commands instead of running them
USAGE
  exit 64
}

run() {
  if (( DRY_RUN )); then
    echo "+ $*"
  else
    "$@"
  fi
}

while getopts ":n" opt; do
  case "$opt" in
    n) DRY_RUN=1 ;;
    *) usage ;;
  esac
done
shift $((OPTIND - 1))

[[ $# -eq 2 ]] || usage
env="$1" version="$2"

if [[ ! " ${ENVIRONMENTS[*]} " =~ " ${env} " ]]; then
  log ERROR "unknown environment: $env"
  exit 1
fi

if ! [[ "$version" =~ ^v[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
  log ERROR "version must look like v1.2.3, got '$version'"
  exit 1
fi

log INFO "deploying $version to $env"
run docker build -t "app:${version}" "$SCRIPT_DIR/.."
run docker push "registry.example.com/app:${version}"
for host in $(awk -v e="$env" '$1 == e { print $2 }' "$SCRIPT_DIR/hosts.txt"); do
  run ssh "deploy@${host}" "systemctl restart app@${version#v}" || log WARN "restart failed on $host"
done
log INFO "done"
