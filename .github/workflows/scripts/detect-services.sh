#!/bin/bash
# detect-services.sh
# Detects which services changed and outputs a JSON matrix for GitHub Actions.
# Usage: bash .github/scripts/detect-services.sh

ALL_SERVICES=(
  create-service
  read-service
  update-service
  delete-service
  data-analysis
)

CHANGED=()

if [[ "${GITHUB_REF}" != refs/tags/* ]]; then
  CHANGED_FILES=$(git diff --name-only HEAD~1 HEAD 2>/dev/null || git diff --name-only HEAD)
  for svc in "${ALL_SERVICES[@]}"; do
    if echo "$CHANGED_FILES" | grep -q "^${svc}/"; then
      CHANGED+=("\"$svc\"")
    fi
  done
fi

if [ ${#CHANGED[@]} -eq 0 ]; then
  for svc in "${ALL_SERVICES[@]}"; do
    CHANGED+=("\"$svc\"")
  done
fi

MATRIX="[$(IFS=,; echo "${CHANGED[*]}")]"
echo "services=$MATRIX" >> $GITHUB_OUTPUT