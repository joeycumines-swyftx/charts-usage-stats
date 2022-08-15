#!/usr/bin/env bash

set -o pipefail

aws_options=''
log_group_name=''
log_stream_prefix=''

print_usage() {
  cat <<EOF
Usage:
  $0 \\
    [-c aws_options] \\
    -g log_group_name -p log_stream_prefix

Description:
  Lists log streams by prefix, printing them one per line.
EOF
}
while getopts ":hc:g:p:" options; do
  case "${options}" in
  c)
    aws_options="${OPTARG}"
    ;;
  g)
    log_group_name="${OPTARG}"
    ;;
  p)
    log_stream_prefix="${OPTARG}"
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
if [ -z "$log_stream_prefix" ]; then
  echo "Error: -p log_stream_prefix is required" >&2
  exit 1
fi

aws ${aws_options} logs describe-log-streams \
  --output=json \
  --no-paginate \
  --log-group-name="${log_group_name}" \
  --log-stream-name-prefix="${log_stream_prefix}" |
  jq -r '.logStreams[].logStreamName'
