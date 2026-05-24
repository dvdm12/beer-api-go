#!/bin/bash
# build-failure-report.sh
# Generates the failure HTML email report by reading the Jenkins build log
# and injecting stage statuses into the email-failure.html template.

set -e

TEMPLATE="infra/templates/email-failure.html"
LOG_FILE="/tmp/build.log"
OUTPUT="/tmp/email-failure-rendered.html"

if [ ! -f "$TEMPLATE" ]; then
    echo "[failure-report] ERROR - template not found: ${TEMPLATE}"
    exit 1
fi

cp "$TEMPLATE" "$OUTPUT"

# Read log and escape HTML special characters
if [ -f "$LOG_FILE" ]; then
    BUILD_LOG=$(cat "$LOG_FILE" \
        | sed 's/&/\&amp;/g' \
        | sed 's/</\&lt;/g' \
        | sed 's/>/\&gt;/g' \
        | sed 's/"/\&quot;/g')
else
    BUILD_LOG="No log available"
fi

# Write log to temp file and use awk to inject
# Note: use 'content' instead of 'log' to avoid awk reserved keyword conflict
awk -v content="$BUILD_LOG" '{
    if (index($0, "${BUILD_LOG}") > 0) {
        gsub(/\$\{BUILD_LOG\}/, content)
    }
    print
}' "$OUTPUT" > /tmp/email-failure-tmp.html && mv /tmp/email-failure-tmp.html "$OUTPUT"

# Detect failed stage from log
FAILED_STAGE="Unknown"
if grep -q "\[coverage\].*ERROR" "$LOG_FILE" 2>/dev/null; then
    FAILED_STAGE="Test Coverage"
elif grep -q "\[analysis\].*ERROR\|EXECUTION FAILURE" "$LOG_FILE" 2>/dev/null; then
    FAILED_STAGE="SonarQube Analysis"
elif grep -q "\[gate\].*FAILED\|Quality Gate failed" "$LOG_FILE" 2>/dev/null; then
    FAILED_STAGE="Quality Gate"
elif grep -q "\[metrics\].*ERROR" "$LOG_FILE" 2>/dev/null; then
    FAILED_STAGE="Fetch Metrics"
elif grep -q "AssertionError\|newman.*failed" "$LOG_FILE" 2>/dev/null; then
    FAILED_STAGE="Black-box Tests"
fi

sed -i "s|\${FAILED_STAGE}|${FAILED_STAGE}|g"               "$OUTPUT"
sed -i "s|\${FAILURE_CAUSE}|Check console log for details|g" "$OUTPUT"

# Detect stage results from log
detect_stage() {
    local key="$1"
    local pass_pattern="$2"
    local fail_pattern="$3"

    if grep -q "$fail_pattern" "$LOG_FILE" 2>/dev/null; then
        sed -i "s|\${${key}_RESULT}|Failed|g"   "$OUTPUT"
        sed -i "s|\${${key}_COLOR}|#f85149|g"   "$OUTPUT"
        sed -i "s|\${${key}_STATUS}|&#10007;|g" "$OUTPUT"
    elif grep -q "$pass_pattern" "$LOG_FILE" 2>/dev/null; then
        sed -i "s|\${${key}_RESULT}|Passed|g"   "$OUTPUT"
        sed -i "s|\${${key}_COLOR}|#3fb950|g"   "$OUTPUT"
        sed -i "s|\${${key}_STATUS}|&#10003;|g" "$OUTPUT"
    else
        sed -i "s|\${${key}_RESULT}|Skipped|g"  "$OUTPUT"
        sed -i "s|\${${key}_COLOR}|#8b949e|g"   "$OUTPUT"
        sed -i "s|\${${key}_STATUS}|\&mdash;|g" "$OUTPUT"
    fi
}

detect_stage "STAGE_COVERAGE" "\[coverage\] Total:"          "\[coverage\] ERROR"
detect_stage "STAGE_SONAR"    "\[analysis\] OK"              "EXECUTION FAILURE"
detect_stage "STAGE_GATE"     "\[gate\] All services passed" "\[gate\] ERROR"
detect_stage "STAGE_METRICS"  "\[metrics\] All metrics"      "\[metrics\] ERROR"
detect_stage "STAGE_NEWMAN"   "newman.*passed"               "AssertionError"

echo "[failure-report] Failure report ready at ${OUTPUT}"