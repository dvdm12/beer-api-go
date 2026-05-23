#!/bin/bash
# build-report.sh
# Orchestrates the final HTML report generation.
# Runs sonar-report.sh and newman-metrics.sh in sequence,
# then copies the final rendered file for Jenkins to read.

set -e

FINAL="/tmp/email-report-rendered.html"

echo "[build-report] Starting report generation..."

# Step 1: Inject SonarQube metrics
chmod +x infra/scripts/sonar/sonar-report.sh
./infra/scripts/sonar/sonar-report.sh

# Step 2: Inject Newman metrics
chmod +x infra/scripts/newman/newman-metrics.sh
./infra/scripts/newman/newman-metrics.sh

# Step 3: Promote final output
cp /tmp/newman-report-partial.html "$FINAL"

echo "[build-report] Final report ready at ${FINAL}"