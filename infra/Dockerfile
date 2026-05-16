FROM jenkins/jenkins:lts
USER root

RUN apt-get update && \
    apt-get install -y wget curl unzip ca-certificates && \
    \
    # Install Go 1.25.8
    wget -O /tmp/go.tar.gz https://go.dev/dl/go1.25.8.linux-amd64.tar.gz && \
    tar -C /usr/local -xzf /tmp/go.tar.gz && \
    rm /tmp/go.tar.gz && \
    \
    # Install SonarScanner 7.1.0.4889
    curl -L -o /tmp/sonar-scanner.zip \
    https://binaries.sonarsource.com/Distribution/sonar-scanner-cli/sonar-scanner-cli-7.1.0.4889-linux-x64.zip && \
    unzip /tmp/sonar-scanner.zip -d /opt && \
    ln -sf /opt/sonar-scanner-7.1.0.4889-linux-x64/bin/sonar-scanner /usr/local/bin/sonar-scanner && \
    rm /tmp/sonar-scanner.zip && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

ENV PATH="/usr/local/go/bin:/opt/sonar-scanner-7.1.0.4889-linux-x64/bin:${PATH}"

USER jenkins