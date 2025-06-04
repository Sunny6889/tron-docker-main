#!/bin/bash

# Capture source and destination addresses for indicators
METRICS_PORT=""
METRICS_URL="http://localhost:${METRICS_PORT}/metrics"

# Config VictoriaMetrics
VICTORIA_METRICS_IP=""
VICTORIA_METRICS_PORT=""
VICTORIA_METRICS_URL="http://${VICTORIA_METRICS_IP}:${VICTORIA_METRICS_PORT}/insert/0/prometheus/api/v1/import/prometheus"

# Configurable label variables
GROUP=""
INSTANCE=""
JOB=""
SLEEP_SECONDS=1

# Build additional tag parameters
EXTRA_LABELS="extra_label=group=${GROUP}&extra_label=instance=${INSTANCE}&extra_label=job=${JOB}"

while true; do
  curl -s "$METRICS_URL" | \
  curl -X POST \
      --data-binary @- \
      -H "Content-Type: text/plain" \
      "${VICTORIA_METRICS_URL}?${EXTRA_LABELS}"
  sleep "$SLEEP_SECONDS"
done
