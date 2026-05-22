#!/bin/bash
# sonar-gate.sh
# Checks the SonarQube Quality Gate status for each microservice.
# Exits with code 1 if any service fails the gate.

set -e

SONAR_URL="http://sonarqube:9000"
MAX_RETRIES=10
RETRY_DELAY=3

SERVICES="create-service read-service update-service delete-service data-analysis"

declare -A PROJECT_KEY=(
    [create-service]="beer-api-go-create-service"
    [read-service]="beer-api-go-read-service"
    [update-service]="beer-api-go-update-service"
    [delete-service]="beer-api-go-delete-service"
    [data-analysis]="beer-api-go-data-analysis"
)

FAILED_SERVICES=""

for svc in $SERVICES; do
    KEY="${PROJECT_KEY[$svc]}"
    echo ""
    echo "[gate] Checking: ${svc} (${KEY})"

    # Wait for analysis task to complete
    for i in $(seq 1 $MAX_RETRIES); do
        TASK_STATUS=$(curl -sf -u "${SONAR_TOKEN}:" \
            "${SONAR_URL}/api/ce/component?component=${KEY}" \
            | grep -o '"status":"[^"]*"' | head -1 \
            | awk -F'"' '{print $4}')

        if [ "$TASK_STATUS" = "SUCCESS" ] || [ "$TASK_STATUS" = "FAILED" ]; then
            break
        fi

        echo "[gate] Task status: ${TASK_STATUS} - retry ${i}/${MAX_RETRIES}"
        sleep $RETRY_DELAY
    done

    # Query Quality Gate result
    GATE_STATUS=$(curl -sf -u "${SONAR_TOKEN}:" \
        "${SONAR_URL}/api/qualitygates/project_status?projectKey=${KEY}" \
        | grep -o '"status":"[^"]*"' | head -1 \
        | awk -F'"' '{print $4}')

    if [ "$GATE_STATUS" = "OK" ]; then
        echo "[gate] PASSED - ${svc}"
    else
        echo "[gate] FAILED - ${svc} (status: ${GATE_STATUS})"
        FAILED_SERVICES="${FAILED_SERVICES} ${svc}"
    fi
done

# Fail pipeline if any service did not pass
if [ -n "$FAILED_SERVICES" ]; then
    echo ""
    echo "[gate] ERROR - Quality Gate failed for:${FAILED_SERVICES}"
    exit 1
fi

echo ""
echo "[gate] All services passed the Quality Gate"