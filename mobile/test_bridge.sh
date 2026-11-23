#!/bin/sh
set -eu
mobile_dir=$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)
"$mobile_dir/build_native.sh" macosx
mkdir -p "$mobile_dir/build/tests"
xcrun swiftc -parse-as-library -O -module-cache-path "${TMPDIR:-/tmp}/pvote-swift-cache" \
 -import-objc-header "$mobile_dir/ios/PVote/PVote-Bridging-Header.h" \
 -L "$mobile_dir/build/macosx" -lPVote -lresolv \
 "$mobile_dir/ios/PVote/NativeEngine.swift" "$mobile_dir/ios/PVote/ChainClient.swift" \
 "$mobile_dir/tests/BridgeSmoke.swift" -o "$mobile_dir/build/tests/bridge-smoke"
"$mobile_dir/build/tests/bridge-smoke" "$mobile_dir/fixtures/fixtures.json"
