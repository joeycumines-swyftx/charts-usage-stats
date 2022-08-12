#!/usr/bin/env bash

set -o pipefail

aws_options=(
  --profile=dev
  --region=ap-southeast-2
)
log_group_name=/aws/containerinsights/development-cluster/application
log_stream_name=api-candle-server-76cb8dc75-95z87_swy-charts_api-candle-server-57d9d49ced39d685ed70161b422c8dea62d944606fc0445caef478b4cd9e6af8
start_time="$(gdate -d '-1 day' +%s)000"
end_time="$(gdate -d '+1 day' +%s)000"
next_token=''
page_number=0

while :; do
  page_number="$((page_number + 1))"
  echo "Dumping page ${page_number}" >&2
  command_options=(
    --start-from-head
    --output=json
    --log-group-name="$log_group_name"
    --log-stream-name="$log_stream_name"
    --start-time="$start_time"
    --end-time="$end_time"
  )
  if [ ! -z "$next_token" ]; then
    command_options+=(--next-token="$next_token")
  fi
  if ! data="$(aws "${aws_options[@]}" logs get-log-events "${command_options[@]}")"; then
    echo "Failed to get log events" >&2
    exit 1
  fi
  if ! jq -c '.events[]' <<<"$data"; then
    echo "Failed to output log events" >&2
    exit 1
  fi
  old_token="$next_token"
  if ! next_token="$(jq -r '.nextForwardToken' <<<"$data")"; then
    echo "Failed to get next token" >&2
    exit 1
  fi
  if [ -z "$next_token" ] || [ "$next_token" = null ] || [ "$next_token" = "$old_token" ]; then
    break
  fi
done

echo "Dumped ${page_number} pages" >&2
