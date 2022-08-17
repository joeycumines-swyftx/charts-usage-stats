#!/usr/bin/env bash

set -o pipefail

if [ -z "$REPLICA_SET_NAME" ]; then
  echo "Error: REPLICA_SET_NAME is not set" >&2
  exit 1
fi

if ! { SCRATCH_DIR="${SCRATCH_DIR:-scratch}" &&
  OUTPUT_NAME="${REPLICA_SET_NAME}-${START_TIME}-${END_TIME}" &&
  OUTPUT_DIR="${SCRATCH_DIR}/${OUTPUT_NAME}" &&
  OUTPUT_FILE="${OUTPUT_DIR}.tar.gz" &&
  mkdir -p "${OUTPUT_DIR}"; }; then
  echo "Error: failed to init output dir: ${SCRATCH_DIR}/${REPLICA_SET_NAME}-${START_TIME}-${END_TIME}" >&2
  exit 1
fi

if ! { hack/list-log-streams.sh -c "$AWS_CLI_OPTIONS" -g "$LOG_GROUP_NAME" -p "${REPLICA_SET_NAME}-" |
  while read -r log_stream_name; do
    echo "Dumping log stream: ${log_stream_name}" >&2
    hack/dump-log-stream.sh \
      -c "$AWS_CLI_OPTIONS" \
      -s "$START_TIME" \
      -e "$END_TIME" \
      -g "$LOG_GROUP_NAME" \
      -n "$log_stream_name" \
      >"${OUTPUT_DIR}/${log_stream_name}.jsonl" ||
      exit 1
  done &&
  tar -C "$SCRATCH_DIR" -czvf "$OUTPUT_FILE" "$OUTPUT_NAME"; }; then
  echo "Error: failed to archive log streams" >&2
  exit 1
fi

echo "Success: archived log streams" >&2
echo "$OUTPUT_FILE"
