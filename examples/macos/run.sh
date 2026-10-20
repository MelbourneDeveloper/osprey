#!/usr/bin/env bash
# Implements [IOS-TARGET-MACOS]. Build and print the storage report:
#   run.sh [folder] [limit]    report on a folder (default: your home folder)
#   run.sh --smoke             assert the boundary, volume and walk answers
#   run.sh --build             build only
# OSPREY_MACOS_SKIP_BUILD=1 reuses the existing build.
set -euo pipefail

example_dir=$(cd "$(dirname "$0")" && pwd)
host="$example_dir/build/osprey-storage"
if [[ ${OSPREY_MACOS_SKIP_BUILD:-0} != 1 ]]; then
    bash "$example_dir/../../scripts/macos-host.sh" "$example_dir/storage" \
        "$example_dir/OspreyStorage" "$host" >&2
fi
[[ ${1:-} == --build ]] && exit 0
exec "$host" "$@"
