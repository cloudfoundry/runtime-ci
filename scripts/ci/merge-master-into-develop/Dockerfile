FROM ubuntu:trusty

ENV HOME /root
RUN \
  apt-get update && \
  apt-get -y install \
    build-essential \
    software-properties-common

RUN \
  add-apt-repository ppa:git-core/ppa -y && \
  apt-get update && \
  apt-get install -y git && \
  apt-get remove -y --purge software-properties-common

