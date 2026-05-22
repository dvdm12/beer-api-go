#!/bin/bash
# sonar-report.sh
# Reads metrics from /tmp/sonar-metrics.env and injects all placeholder
# values into the HTML email template, producing a rendered output file.

set -e

TEMPLATE="infra/templates/email-report.html"
METRICS_FILE="/tmp/sonar-metrics.env"
OUTPUT="/tmp/email-report-rendered.html"

# Validate inputs
if [ ! -f "$TEMPLATE" ]; then
    echo "[report] ERROR - template not found: ${TEMPLATE}"
    exit 1
fi

if [ ! -f "$METRICS_FILE" ]; then
    echo "[report] ERROR - metrics file not found: ${METRICS_FILE}"
    echo "[report] Run sonar-metrics.sh before this script"
    exit 1
fi

# Copy template to output
cp "$TEMPLATE" "$OUTPUT"

# Load and inject all sonar metric variables
while IFS='=' read -r key value; do
    # Skip empty lines and comments
    [[ -z "$key" || "$key" == \#* ]] && continue
    sed -i "s|\${${key}}|${value}|g" "$OUTPUT"
done < "$METRICS_FILE"

echo "[report] Metrics injected from ${METRICS_FILE}"
echo "[report] Report written to ${OUTPUT}"