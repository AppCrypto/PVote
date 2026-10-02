#!/usr/bin/env bash
set -euo pipefail

repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
output_dir="$repo_root/experiments/results"
docker_bin="${DOCKER_BIN:-docker}"
image_name="pvote-constrained-client:latest"

mkdir -p "$output_dir"
"$docker_bin" build --platform linux/arm64 \
  --file "$repo_root/experiments/constrained-client/Dockerfile" \
  --tag "$image_name" \
  "$repo_root"
"$docker_bin" run --rm --platform linux/arm64 --cpus=1 --memory=1g --pids-limit=256 \
  --ulimit nofile=1024:1024 --volume "$output_dir:/results" "$image_name" \
  -out /results/c1m1_rbpvss_share_sweep.csv -runs 20 -setup-runs 2 -share-only
"$docker_bin" run --rm --platform linux/arm64 --cpus=1 --memory=1g --pids-limit=256 \
  --ulimit nofile=1024:1024 --volume "$output_dir:/results" \
  --entrypoint /usr/local/bin/full-comparison-bench "$image_name" \
  -out /results/c1m1_score_domain_comparison.csv -runs 3 -score-domain-only
