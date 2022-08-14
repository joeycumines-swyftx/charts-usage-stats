#!/usr/bin/env bash

set -o pipefail

aws_options=''
start_time=''
end_time=''
log_group_name=''
log_stream_name=''
next_token=''
page_number=0

print_usage() {
  cat <<EOF
Usage:
  $0 \\
    [-c aws_options] [-s start_time] [-e end_time] \\
    -g log_group_name -n log_stream_name
EOF
}
while getopts ":hc:s:e:g:n:" options; do
  case "${options}" in
  c)
    aws_options="${OPTARG}"
    ;;
  s)
    start_time="${OPTARG}"
    ;;
  e)
    end_time="${OPTARG}"
    ;;
  g)
    log_group_name="${OPTARG}"
    ;;
  n)
    log_stream_name="${OPTARG}"
    ;;
  h)
    print_usage
    exit 0
    ;;
  :)
    echo "Error: -${OPTARG} requires an argument." >&2
    exit 1
    ;;
  *)
    echo "Error: Unknown option: -${OPTARG}" >&2
    exit 1
    ;;
  esac
done
if [ -z "$log_group_name" ]; then
  echo "Error: -g log_group_name is required" >&2
  exit 1
fi
if [ -z "$log_stream_name" ]; then
  echo "Error: -n log_stream_name is required" >&2
  exit 1
fi

while :; do
  page_number="$((page_number + 1))"
  echo "Dumping page ${page_number}" >&2
  command_options=(
    --start-from-head
    --output=json
    --log-group-name="$log_group_name"
    --log-stream-name="$log_stream_name"
  )
  if [ ! -z "$start_time" ]; then
    command_options+=(--start-time="$start_time")
  fi
  if [ ! -z "$end_time" ]; then
    command_options+=(--end-time="$end_time")
  fi
  if [ ! -z "$next_token" ]; then
    command_options+=(--next-token="$next_token")
  fi
  if ! data="$(aws ${aws_options} logs get-log-events "${command_options[@]}")"; then
    echo "Error: Failed to get log events" >&2
    exit 1
  fi
  if ! jq -c '.events[]' <<<"$data"; then
    echo "Error: Failed to output log events" >&2
    exit 1
  fi
  old_token="$next_token"
  if ! next_token="$(jq -r '.nextForwardToken' <<<"$data")"; then
    echo "Error: Failed to get next token" >&2
    exit 1
  fi
  if [ -z "$next_token" ] || [ "$next_token" = null ] || [ "$next_token" = "$old_token" ]; then
    break
  fi
done

echo "Dumped ${page_number} pages" >&2
