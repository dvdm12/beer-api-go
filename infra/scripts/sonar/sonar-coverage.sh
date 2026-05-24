#!/bin/bash
# sonar-coverage.sh
# Runs Go tests with coverage profiling for each microservice.
# Generates coverage.out required by sonar-analysis.sh.
# Counts passing tests per service and writes to /tmp/sonar-metrics.env.
# Only covers business logic packages: controllers, services, repository.
# Exits with code 1 if any unit test fails, aborting the pipeline.

set -e

SERVICES="create-service read-service update-service delete-service data-analysis"
METRICS_FILE="/tmp/sonar-metrics.env"

# Packages containing business logic to measure coverage against
TEST_PACKAGES="./internal/controllers ./internal/services ./internal/repository"

# Initialize metrics file
> "$METRICS_FILE"

TOTAL_TESTS=0

for svc in $SERVICES; do
    echo ""
    echo "[coverage] Processing: ${svc}"

    cd "$svc"

    # Derive SVC_KEY: remove -service suffix, map data-analysis to analysis, uppercase
    SVC_KEY=$(echo "$svc" \
        | sed 's/-service$//' \
        | sed 's/data-analysis/analysis/' \
        | tr '[:lower:]' '[:upper:]')

    # Run tests with JSON output and coverage profile
    # Any test failure exits immediately and aborts the pipeline
    if ! go test $TEST_PACKAGES \
        -coverprofile=coverage.out \
        -covermode=atomic \
        -json > test-report.json 2>&1; then

        echo "[coverage] ERROR - unit tests failed in ${svc}"
        grep '"Action":"fail"' test-report.json | grep '"Test"' || true
        exit 1
    fi

    echo "[coverage] OK - coverage.out generated for ${svc}"

    # Print total coverage summary
    go tool cover -func=coverage.out | grep "^total:" | awk '{print "[coverage] Total:", $3}'

    # Count passing tests from JSON report
    TEST_COUNT=$(grep '"Action":"pass"' test-report.json \
        | grep '"Test"' \
        | wc -l \
        | awk '{print $1}')

    echo "[coverage] Tests passed: ${TEST_COUNT} for ${svc}"

    # Write test count to metrics file
    echo "SONAR_TESTS_${SVC_KEY}=${TEST_COUNT}" >> "$METRICS_FILE"

    TOTAL_TESTS=$((TOTAL_TESTS + TEST_COUNT))

    cd ..
done

# Write global test total
echo "SONAR_TOTAL_TESTS=${TOTAL_TESTS}" >> "$METRICS_FILE"

echo ""
echo "[coverage] Total tests across all services: ${TOTAL_TESTS}"
echo "[coverage] Test counts written to ${METRICS_FILE}"