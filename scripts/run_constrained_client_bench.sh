#!/usr/bin/env bash
set -euo pipefail

repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
output_dir="$repo_root/experiments/results"
image_name="pvote-constrained-client:latest"
docker_bin="${DOCKER_BIN:-docker}"
network_name="pvote-c1m1-net"
receiver_name="pvote-c1m1-receiver"
raw_csv="$output_dir/constrained_client_rbpvss_results.csv"
transfer_samples="$output_dir/constrained_client_upload_samples_ms.txt"
summary_csv="$output_dir/constrained_client_c1m1.csv"

mkdir -p "$output_dir"
cleanup() {
  "$docker_bin" rm -f "$receiver_name" >/dev/null 2>&1 || true
  "$docker_bin" network rm "$network_name" >/dev/null 2>&1 || true
}
trap cleanup EXIT

"$docker_bin" build --platform linux/arm64 \
  --file "$repo_root/experiments/constrained-client/Dockerfile" \
  --tag "$image_name" \
  "$repo_root"
"$docker_bin" run --rm --platform linux/arm64 --cpus=1 --memory=1g --pids-limit=256 \
  --ulimit nofile=1024:1024 \
  --volume "$output_dir:/results" \
  "$image_name" \
  -out /results/constrained_client_rbpvss_results.csv -runs 20 -setup-runs 5 -only-n 100 -only-l 19

payload_bytes="$(awk -F, '$1 == 100 && $2 == 19 && $3 == 59 { printf "%d", $15 * 1024 }' "$raw_csv")"
share_mean_ms="$(awk -F, '$1 == 100 && $2 == 19 && $3 == 59 { print $11 }' "$raw_csv")"
share_std_ms="$(awk -F, '$1 == 100 && $2 == 19 && $3 == 59 { print $12 }' "$raw_csv")"
if [[ -z "$payload_bytes" || -z "$share_mean_ms" ]]; then
  echo "missing (100,19,59) row in $raw_csv" >&2
  exit 1
fi

"$docker_bin" network create "$network_name" >/dev/null
"$docker_bin" run -d --rm --name "$receiver_name" --network "$network_name" \
  --cap-add=NET_ADMIN \
  --entrypoint /bin/sh "$image_name" \
  -ec 'tc qdisc replace dev eth0 root netem delay 40ms 10ms distribution normal; exec constrained-receiver' >/dev/null
sleep 1
"$docker_bin" run --rm --platform linux/arm64 --network "$network_name" \
  --cap-add=NET_ADMIN --cpus=1 --memory=1g --pids-limit=256 --ulimit nofile=1024:1024 \
  --entrypoint /bin/sh "$image_name" -ec '
    tc qdisc replace dev eth0 root netem delay 40ms 10ms distribution normal rate 5mbit
    dd if=/dev/zero of=/tmp/ballot.bin bs=1 count="$1" status=none
    for _ in $(seq 1 20); do
      curl --silent --show-error --output /dev/null --write-out "%{time_total}\n" \
        --request POST --data-binary @/tmp/ballot.bin http://'"$receiver_name"':8080/ballot
    done
  ' sh "$payload_bytes" | awk '{ printf "%.6f\n", $1 * 1000 }' > "$transfer_samples"

upload_mean_ms="$(awk '{ sum += $1; n += 1 } END { if (n != 20) exit 1; printf "%.6f", sum/n }' "$transfer_samples")"
upload_std_ms="$(awk -v mean="$upload_mean_ms" '{ sum += ($1-mean)*($1-mean); n += 1 } END { printf "%.6f", sqrt(sum/(n-1)) }' "$transfer_samples")"
printf '%s\n' 'profile,n,l,t,runs,share_mean_ms,share_std_ms,payload_bytes,upload_mean_ms,upload_std_ms,cpu,memory,bandwidth,one_way_delay,jitter' > "$summary_csv"
printf 'C1M1,100,19,59,20,%s,%s,%s,%s,%s,1,1GiB,5Mbit/s,40ms,10ms\n' \
  "$share_mean_ms" "$share_std_ms" "$payload_bytes" "$upload_mean_ms" "$upload_std_ms" >> "$summary_csv"
