#!/bin/bash
# sonar-metrics.sh
# Fetches code quality metrics from the SonarQube API for each microservice.
# Appends all values to /tmp/sonar-metrics.env for use in sonar-report.sh.
# Note: file is initialized by sonar-coverage.sh — do not clear it here.

set -e

SONAR_URL="http://sonarqube:9000"
METRICS="coverage,bugs,vulnerabilities,code_smells,duplicated_lines_density,ncloc"
OUTPUT_FILE="/tmp/sonar-metrics.env"

declare -A PROJECT_KEY=(
    [create]="beer-api-go-create-service"
    [read]="beer-api-go-read-service"
    [update]="beer-api-go-update-service"
    [delete]="beer-api-go-delete-service"
    [analysis]="beer-api-go-data-analysis"
)

TOTAL_COV=0
TOTAL_BUGS=0
TOTAL_VULNS=0
TOTAL_SMELLS=0
SVC_COUNT=0

# Extracts a metric value from a JSON response by metric key name
parse_metric() {
    local response="$1"
    local metric="$2"
    echo "$response" \
        | grep -o "\"metric\":\"${metric}\",\"value\":\"[^\"]*\"" \
        | awk -F'"' '{print $8}'
}

# Extracts quality gate status from project_status response
parse_gate() {
    local response="$1"
    echo "$response" \
        | grep -o '"status":"[^"]*"' | head -1 \
        | awk -F'"' '{print $4}'
}

for svc in create read update delete analysis; do
    KEY="${PROJECT_KEY[$svc]}"
    SVC_UPPER=$(echo "$svc" | tr '[:lower:]' '[:upper:]')

    echo "[metrics] Fetching: ${KEY}"

    # Fetch all metrics in a single API call
    RESPONSE=$(curl -sf -u "${SONAR_TOKEN}:" \
        "${SONAR_URL}/api/measures/component?component=${KEY}&metricKeys=${METRICS}" \
        2>/dev/null || echo "{}")

    # Parse individual metric values
    COV=$(parse_metric   "$RESPONSE" "coverage")
    BUGS=$(parse_metric  "$RESPONSE" "bugs")
    VULNS=$(parse_metric "$RESPONSE" "vulnerabilities")
    SMELLS=$(parse_metric "$RESPONSE" "code_smells")
    DUPL=$(parse_metric  "$RESPONSE" "duplicated_lines_density")
    NCLOC=$(parse_metric "$RESPONSE" "ncloc")

    # Apply defaults for missing values
    COV=${COV:-0}
    BUGS=${BUGS:-0}
    VULNS=${VULNS:-0}
    SMELLS=${SMELLS:-0}
    DUPL=${DUPL:-0}
    NCLOC=${NCLOC:-0}

    # Fetch Quality Gate status
    GATE_RESPONSE=$(curl -sf -u "${SONAR_TOKEN}:" \
        "${SONAR_URL}/api/qualitygates/project_status?projectKey=${KEY}" \
        2>/dev/null || echo "{}")

    GATE=$(parse_gate "$GATE_RESPONSE")
    GATE=${GATE:-N/A}

    # Resolve gate CSS class and display label
    if [ "$GATE" = "OK" ]; then
        GATE_CLASS="gate-ok"
        GATE_LABEL="Passed"
    elif [ "$GATE" = "ERROR" ]; then
        GATE_CLASS="gate-err"
        GATE_LABEL="Failed"
    else
        GATE_CLASS="gate-na"
        GATE_LABEL="N/A"
    fi

    # Round decimal values to one place
    COV_R=$(printf "%.1f"  "$COV"  2>/dev/null || echo "0.0")
    DUPL_R=$(printf "%.1f" "$DUPL" 2>/dev/null || echo "0.0")

    # Append per-service variables to env file
    echo "SONAR_COV_${SVC_UPPER}=${COV_R}"            >> "$OUTPUT_FILE"
    echo "SONAR_BUGS_${SVC_UPPER}=${BUGS}"             >> "$OUTPUT_FILE"
    echo "SONAR_VULNS_${SVC_UPPER}=${VULNS}"           >> "$OUTPUT_FILE"
    echo "SONAR_SMELLS_${SVC_UPPER}=${SMELLS}"         >> "$OUTPUT_FILE"
    echo "SONAR_DUPL_${SVC_UPPER}=${DUPL_R}"           >> "$OUTPUT_FILE"
    echo "SONAR_NCLOC_${SVC_UPPER}=${NCLOC}"           >> "$OUTPUT_FILE"
    echo "SONAR_GATE_${SVC_UPPER}=${GATE_CLASS}"       >> "$OUTPUT_FILE"
    echo "SONAR_GATE_LABEL_${SVC_UPPER}=${GATE_LABEL}" >> "$OUTPUT_FILE"

    # Accumulate totals for global summary
    TOTAL_COV=$(awk "BEGIN {printf \"%.1f\", $TOTAL_COV + $COV}")
    TOTAL_BUGS=$((TOTAL_BUGS + BUGS))
    TOTAL_VULNS=$((TOTAL_VULNS + VULNS))
    TOTAL_SMELLS=$((TOTAL_SMELLS + SMELLS))
    SVC_COUNT=$((SVC_COUNT + 1))

    echo "[metrics] OK - ${svc}: cov=${COV_R}% bugs=${BUGS} vulns=${VULNS} smells=${SMELLS}"
done

# Append global summary to env file
AVG_COV=$(awk "BEGIN {printf \"%.1f\", $TOTAL_COV / $SVC_COUNT}")
echo "SONAR_TOTAL_COVERAGE=${AVG_COV}"    >> "$OUTPUT_FILE"
echo "SONAR_TOTAL_BUGS=${TOTAL_BUGS}"     >> "$OUTPUT_FILE"
echo "SONAR_TOTAL_VULNS=${TOTAL_VULNS}"   >> "$OUTPUT_FILE"
echo "SONAR_TOTAL_SMELLS=${TOTAL_SMELLS}" >> "$OUTPUT_FILE"

echo ""
echo "[metrics] All metrics written to ${OUTPUT_FILE}"