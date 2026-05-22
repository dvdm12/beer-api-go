#!/bin/bash
# sonar-coverage.sh
# Runs Go tests with coverage profiling for each microservice.
# Generates coverage.out required by sonar-analysis.sh.

set -e

SERVICES="create-service read-service update-service delete-service data-analysis"

for svc in $SERVICES; do
    echo ""
    echo "[coverage] Processing: ${svc}"

    cd "$svc"

    # Run tests and generate coverage profile
    if go test ./... -coverprofile=coverage.out -covermode=atomic 2>/dev/null; then
        echo "[coverage] OK - coverage.out generated for ${svc}"

        # Print total coverage summary
        go tool cover -func=coverage.out | grep "^total:" | awk '{print "[coverage] Total:", $3}'
    else
        # Tests failed but do not block the pipeline
        echo "[coverage] WARN - tests failed or no test files found in ${svc}"
        touch coverage.out
    fi

    cd ..
done