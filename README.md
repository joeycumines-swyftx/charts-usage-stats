# charts-usage-stats

Utils as part of manually extracting usage statistics from charts.

## Dependencies

### Commands

- `bash` note that you may need a more recent version than comes with MacOS
- `jq` command line JSON tool
- `aws` specifically AWS CLI v2

### Setup

It's recommended that you follow the
[k8s setup instructions](https://swyftx.atlassian.net/wiki/spaces/SREPS/pages/1809776803/Configuring+CLI+access+to+our+Kubernetes+infrastructure),
as part of setting up your AWS CLI access.

## Usage

### Dump all the logs for a replica set

Handy if you want all the logs, for all pods, between redeploy of a k8s deployment.

```bash
# you can find this multiple ways, you could use kubectl like:
#   kubectl --kubeconfig="$HOME"/.kube/dev.yaml -n swy-charts get replicasets
# or you could just go to cloudwatch, and infer it from the log stream names
REPLICA_SET_NAME='api-candle-server-5d8c6985f6'

# optional, can also accept strings like '2018-01-01T00:00:00Z'
# note: non-macos users use date, macos users must have coreutils installed
START_TIME="$(gdate -d '-1 day' +%s)000"
END_TIME="$(gdate +%s)000"

# https://github.com/Swyftx/Operations/blob/99291104af1028316b211b3332962d297ba4f8ee/application/modules/cluster/kubernetes/templates/fluent-bit/application-log.conf#L84
LOG_GROUP_NAME='/aws/containerinsights/development-cluster/application'

# change this for your target environment
AWS_CLI_OPTIONS='--profile=dev --region=ap-southeast-2'

# where you want to put the logs, which will be one file per log stream
OUTPUT_DIR='scratch'
mkdir -p "$OUTPUT_DIR"

# fail if the lhs of a pipe command fails, handy
set -o pipefail

hack/list-log-streams.sh -c "$AWS_CLI_OPTIONS" -g "$LOG_GROUP_NAME" -p "${REPLICA_SET_NAME}-" |
  while read -r log_stream_name; do
    hack/dump-log-stream.sh \
      -c "$AWS_CLI_OPTIONS" \
      -s "$START_TIME" \
      -e "$END_TIME" \
      -g "$LOG_GROUP_NAME" \
      -n "$log_stream_name" \
      >"${OUTPUT_DIR}/${log_stream_name}.jsonl" ||
      exit 1
  done
```
