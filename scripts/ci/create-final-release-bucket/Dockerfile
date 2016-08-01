FROM ubuntu:trusty

ENV HOME /root

RUN \
  apt-get update && \
  apt-get -y install \
    build-essential \
    wget \
    python-software-properties \
    python-pip \
    jq

RUN \
  pip install awscli==1.6.6
