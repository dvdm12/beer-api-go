#!/bin/bash
# sonar-analysis.sh
# Runs SonarScanner for each microservice.
# Requires coverage.out to be present (run sonar-coverage.sh first).

set -e

export PATH="$PATH:/opt/sonar-scanner-7.1.0.4889-linux-x64/bin"

SERVICES="create-service read-service update-service delete-service data-analysis"

for svc in $SERVICES; do
    echo ""
    echo "[analysis] Processing: ${svc}"

    cd "$svc"

    # Validate required configuration file
    if [ ! -f sonar-project.properties ]; then
        echo "[analysis] ERROR - sonar-project.properties not found in ${svc}"
        exit 1
    fi

    # Warn if coverage report is missing
    if [ ! -f coverage.out ]; then
        echo "[analysis] WARN - coverage.out not found in ${svc}, coverage will show 0%"
    fi

    # Run static analysis
    sonar-scanner
    echo "[analysis] OK - scan completed for ${svc}"

    cd ..
done