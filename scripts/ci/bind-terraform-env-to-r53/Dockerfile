FROM ubuntu:trusty

RUN \
  apt-get update && \
  apt-get -y install \
    build-essential \
    mercurial \
    wget \
    python-pip \
    dnsutils \
    jq

RUN \
  pip install awscli==1.6.6
