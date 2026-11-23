#!/bin/sh
set -eu
mobile_dir=$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)
repo_dir=$(dirname "$mobile_dir")
platform=${1:-all}
go_bin=${PVOTE_GO:-}
if [ -z "$go_bin" ]; then
  go_bin=$(command -v go || true)
fi
if [ -z "$go_bin" ]; then
  for candidate in /opt/homebrew/bin/go /usr/local/opt/go/libexec/bin/go /usr/local/go/bin/go; do
    if [ -x "$candidate" ]; then go_bin=$candidate; break; fi
  done
fi
if [ -z "$go_bin" ]; then echo 'Go not found; set PVOTE_GO to its executable path.' >&2; exit 1; fi
if [ "$platform" = all ]; then
  "$0" iphoneos
  "$0" iphonesimulator
  exit 0
fi
case "$platform" in
 iphoneos) sdk=iphoneos; target=arm64-apple-ios17.0; goos=ios ;;
 iphonesimulator) sdk=iphonesimulator; target=arm64-apple-ios17.0-simulator; goos=ios ;;
 macosx) sdk=macosx; target=arm64-apple-macos14.0; goos=darwin ;;
 *) echo 'Usage: build_native.sh [all|iphoneos|iphonesimulator|macosx]' >&2; exit 2 ;;
esac
sdk_path=$(xcrun --sdk "$sdk" --show-sdk-path)
clang_path=$(xcrun --sdk "$sdk" --find clang)
revision=$(git -C "$repo_dir" rev-parse --short=12 HEAD)
if [ -n "$(git -C "$repo_dir" status --porcelain -- mobile baselines crypto chain compile go.mod go.sum)" ]; then revision="$revision-dirty"; fi
source_sha=$(python3 - "$repo_dir" <<'PY'
import hashlib, pathlib, sys
root = pathlib.Path(sys.argv[1])
files = [root / 'go.mod', root / 'go.sum']
for directory in ['crypto', 'chain', 'compile', 'mobile', 'baselines']:
    files.extend(p for p in (root / directory).rglob('*') if p.is_file()
                 and not {'build', 'results', 'xcuserdata', 'deriveddata'} & {s.lower() for s in p.relative_to(root).parts}
                 and (p.suffix in {'.go', '.swift', '.h', '.sol', '.sh', '.pbxproj', '.plist'}
                      or p.relative_to(root).as_posix() == 'mobile/fixtures/fixtures.json'))
h = hashlib.sha256()
for path in sorted(files):
    h.update(str(path.relative_to(root)).encode() + b'\0')
    h.update(path.read_bytes() + b'\0')
print(h.hexdigest())
PY
)
mkdir -p "$mobile_dir/build/$platform"
cd "$repo_dir"
GOCACHE="${TMPDIR:-/tmp}/pvote-go-cache" GOOS="$goos" GOARCH=arm64 CGO_ENABLED=1 \
 CC="$clang_path" CGO_CFLAGS="-target $target -isysroot $sdk_path -O2" \
 CGO_LDFLAGS="-target $target -isysroot $sdk_path" \
 "$go_bin" build -trimpath -buildmode=c-archive -ldflags "-s -w -X PVote/mobile/core.Revision=$revision -X PVote/mobile/core.SourceSHA=$source_sha" \
 -o "$mobile_dir/build/$platform/libPVote.a" ./mobile/bridge
echo "Built native RB-PVSS for $platform ($revision)"
