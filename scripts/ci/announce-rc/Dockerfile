FROM ubuntu:trusty

ENV HOME /root

RUN \
  echo "deb http://archive.ubuntu.com/ubuntu trusty-backports main restricted universe multiverse" > /etc/apt/sources.list.d/backports.list && \
  apt-get update && \
  apt-get -y install \
    build-essential \
    curl \
    jq/trusty-backports \
    software-properties-common

RUN \
  add-apt-repository ppa:git-core/ppa -y && \
  apt-get update && \
  apt-get install -y git && \
  apt-get remove -y --purge software-properties-common

