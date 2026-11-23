#!/bin/sh
set -eu
repo_dir=$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)
cd "$repo_dir"
out=${1:-"experiments/mac-generation-$(date +%Y%m%d_%H%M%S)"}
if [ -e "$out/summary.csv" ]; then
    echo 'Choose a new output directory to preserve the existing experiment.' >&2
    exit 1
fi
source_sha=$(python3 - <<'PY'
import hashlib,pathlib
root=pathlib.Path.cwd();files=[root/'go.mod',root/'go.sum']
for directory in ['mobile/core','baselines','crypto','chain','compile','cmd/score-generation-bench']:
    files.extend(p for p in (root/directory).rglob('*') if p.is_file() and p.suffix in ['.go','.sol'])
h=hashlib.sha256()
for p in sorted(files):
    h.update(str(p.relative_to(root)).encode()+b'\0');h.update(p.read_bytes()+b'\0')
print(h.hexdigest())
PY
)
revision=$(git rev-parse --short=12 HEAD)
if [ -n "$(git status --porcelain)" ]; then revision="$revision-dirty"; fi
build_dir=$(mktemp -d "${TMPDIR:-/tmp}/pvote-score-generation.XXXXXX")
trap 'rm -f "$build_dir/bench"; rmdir "$build_dir"' EXIT HUP INT TERM
GOCACHE="${TMPDIR:-/tmp}/pvote-go-cache" go build -trimpath \
    -ldflags "-X PVote/mobile/core.SourceSHA=$source_sha -X PVote/mobile/core.Revision=$revision" \
    -o "$build_dir/bench" ./cmd/score-generation-bench
"$build_dir/bench" --out "$out/raw"
python3 mobile/summarize_comparison.py "$out"/raw/benchmark-*.json --out "$out"
