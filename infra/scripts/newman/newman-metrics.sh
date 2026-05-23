#!/bin/bash
# newman-metrics.sh
# Reads Newman metrics from /tmp/sonar-metrics.env and injects
# only Newman-related placeholders into the partial HTML report.
# Input is /tmp/sonar-report-partial.html from sonar-report.sh.
# Output is /tmp/newman-report-partial.html.

set -e

INPUT="/tmp/sonar-report-partial.html"
METRICS_FILE="/tmp/sonar-metrics.env"
OUTPUT="/tmp/newman-report-partial.html"

if [ ! -f "$INPUT" ]; then
    echo "[newman-metrics] ERROR - partial report not found: ${INPUT}"
    echo "[newman-metrics] Run sonar-report.sh before this script"
    exit 1
fi

if [ ! -f "$METRICS_FILE" ]; then
    echo "[newman-metrics] ERROR - metrics file not found: ${METRICS_FILE}"
    exit 1
fi

cp "$INPUT" "$OUTPUT"

# Inject only Newman variables
while IFS='=' read -r key value; do
    [[ -z "$key" || "$key" == \#* ]] && continue
    [[ "$key" != NEWMAN_* ]] && continue
    sed -i "s|\${${key}}|${value}|g" "$OUTPUT"
done < "$METRICS_FILE"

echo "[newman-metrics] Newman metrics injected into ${OUTPUT}"