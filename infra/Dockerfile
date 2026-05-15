FROM jenkins/jenkins:lts

USER root

RUN apt-get update && \
    apt-get install -y wget && \
    wget https://go.dev/dl/go1.25.4.linux-amd64.tar.gz && \
    tar -C /usr/local -xzf go1.25.4.linux-amd64.tar.gz && \
    rm go1.25.4.linux-amd64.tar.gz && \
    apt-get clean

ENV PATH=$PATH:/usr/local/go/bin

USER jenkins