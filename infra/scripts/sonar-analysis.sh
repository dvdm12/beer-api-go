#!/bin/bash
set -e

export PATH="$PATH:/var/jenkins_home/tools/hudson.plugins.sonar.SonarRunnerInstallation/SonarQube_Scanner/bin"

SERVICES="create-service read-service update-service delete-service data-analysis"

for svc in $SERVICES; do
    echo
    echo "Analyzing ${svc}..."

    cd "$svc"

    if [ ! -f sonar-project.properties ]; then
        echo "ERROR: sonar-project.properties not found in ${svc}"
        exit 1
    fi

    sonar-scanner

    cd ..
done