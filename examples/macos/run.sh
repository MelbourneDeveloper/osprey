#!/usr/bin/env bash
# Implements [IOS-TARGET-MACOS]. Build and print the storage report:
#   run.sh [folder] [limit]    report on a folder (default: your home folder)
#   run.sh --smoke             assert the boundary, volume and index answers
# OSPREY_MACOS_SKIP_BUILD=1 reuses the existing build.
set -euo pipefail

example_dir=$(cd "$(dirname "$0")" && pwd)
[[ ${OSPREY_MACOS_SKIP_BUILD:-0} == 1 ]] || bash "$example_dir/build.sh" >&2
host="$example_dir/build/osprey-storage"
if [[ ${1:-} != --smoke ]]; then
    exec "$host" "$@"
fi
# Index assertions are required wherever Spotlight indexes the boot volume.
# Hosted runners usually disable it, and that is reported rather than hidden.
if mdutil -s / 2>/dev/null | grep -q "Indexing enabled"; then
    exec "$host" --smoke --require-index
fi
echo "Spotlight indexing is disabled on /; index answers are not exercised." >&2
exec "$host" --smoke
