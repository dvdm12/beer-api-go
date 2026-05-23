#!/bin/bash
# sonar-report.sh
# Reads SonarQube metrics from /tmp/sonar-metrics.env and injects
# only SonarQube-related placeholders into the HTML template.
# Output is written to /tmp/sonar-report-partial.html.

set -e

TEMPLATE="infra/templates/email-report.html"
METRICS_FILE="/tmp/sonar-metrics.env"
OUTPUT="/tmp/sonar-report-partial.html"

if [ ! -f "$TEMPLATE" ]; then
    echo "[sonar-report] ERROR - template not found: ${TEMPLATE}"
    exit 1
fi

if [ ! -f "$METRICS_FILE" ]; then
    echo "[sonar-report] ERROR - metrics file not found: ${METRICS_FILE}"
    exit 1
fi

cp "$TEMPLATE" "$OUTPUT"

# Inject only SonarQube and coverage variables
while IFS='=' read -r key value; do
    [[ -z "$key" || "$key" == \#* ]] && continue
    [[ "$key" == NEWMAN_* ]] && continue
    sed -i "s|\${${key}}|${value}|g" "$OUTPUT"
done < "$METRICS_FILE"

echo "[sonar-report] SonarQube metrics injected into ${OUTPUT}"