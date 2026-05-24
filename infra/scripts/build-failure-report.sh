#!/bin/bash
# build-failure-report.sh
# Generates the failure HTML report by reading the Jenkins build log
# and injecting stage statuses into the email-failure.html template.

set -e

TEMPLATE="infra/templates/email-failure.html"
LOG_FILE="/tmp/build.log"
OUTPUT="/tmp/email-failure-rendered.html"

if [ ! -f "$TEMPLATE" ]; then
    echo "[failure-report] ERROR - template not found: ${TEMPLATE}"
    exit 1
fi

if [ ! -f "$LOG_FILE" ]; then
    echo "[failure-report] WARN - build log not found, using empty trace"
    touch "$LOG_FILE"
fi

cp "$TEMPLATE" "$OUTPUT"

# Inject build log as error trace (last 60 lines)
BUILD_LOG=$(tail -60 "$LOG_FILE" | sed 's/&/\&amp;/g; s/</\&lt;/g; s/>/\&gt;/g')
sed -i "s|\${BUILD_LOG}|${BUILD_LOG}|g" "$OUTPUT"

# Detect failed stage from log
FAILED_STAGE="Unknown"
for stage in "Test Coverage" "SonarQube Analysis" "Quality Gate" "Fetch Metrics" "Black-box Tests"; do
    if grep -q "stage.*${stage}.*failed\|${stage}.*ERROR\|${stage}.*exit code" "$LOG_FILE" 2>/dev/null; then
        FAILED_STAGE="$stage"
        break
    fi
done

sed -i "s|\${FAILED_STAGE}|${FAILED_STAGE}|g" "$OUTPUT"
sed -i "s|\${FAILURE_CAUSE}|Check console log for details|g" "$OUTPUT"

# Detect stage results from log
detect_stage() {
    local stage_name="$1"
    local key="$2"

    if grep -q "\[${stage_name}\].*OK\|stage.*${stage_name}.*completed" "$LOG_FILE" 2>/dev/null; then
        sed -i "s|\${${key}_RESULT}|Passed|g" "$OUTPUT"
        sed -i "s|\${${key}_COLOR}|#3fb950|g" "$OUTPUT"
        sed -i "s|\${${key}_STATUS}|&#10003;|g" "$OUTPUT"
    elif grep -q "\[${stage_name}\].*ERROR\|stage.*${stage_name}.*failed" "$LOG_FILE" 2>/dev/null; then
        sed -i "s|\${${key}_RESULT}|Failed|g" "$OUTPUT"
        sed -i "s|\${${key}_COLOR}|#f85149|g" "$OUTPUT"
        sed -i "s|\${${key}_STATUS}|&#10007;|g" "$OUTPUT"
    else
        sed -i "s|\${${key}_RESULT}|Skipped|g" "$OUTPUT"
        sed -i "s|\${${key}_COLOR}|#8b949e|g" "$OUTPUT"
        sed -i "s|\${${key}_STATUS}|&mdash;|g" "$OUTPUT"
    fi
}

detect_stage "coverage" "STAGE_COVERAGE"
detect_stage "analysis" "STAGE_SONAR"
detect_stage "gate"     "STAGE_GATE"
detect_stage "metrics"  "STAGE_METRICS"
detect_stage "newman"   "STAGE_NEWMAN"

echo "[failure-report] Failure report ready at ${OUTPUT}"