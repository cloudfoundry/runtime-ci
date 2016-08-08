FROM ubuntu:trusty

RUN \
  apt-get update && \
  apt-get -y install \
    build-essential \
    python-pip \
    dnsutils

RUN \
  pip install awscli==1.6.6

COPY resolv.conf /etc/resolv.conf
