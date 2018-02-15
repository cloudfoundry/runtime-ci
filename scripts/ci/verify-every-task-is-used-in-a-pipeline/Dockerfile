FROM ubuntu:trusty

RUN \
  echo "deb http://archive.ubuntu.com/ubuntu trusty-backports main restricted universe multiverse" > /etc/apt/sources.list.d/backports.list && \
  apt-get update && \
  apt-get -y install \
    build-essential \
    software-properties-common \
    ack-grep

RUN dpkg-divert --local --divert /usr/bin/ack --rename --add /usr/bin/ack-grep

RUN \
  add-apt-repository ppa:git-core/ppa -y && \
  apt-get update && \
  apt-get install -y git && \
  apt-get remove -y --purge software-properties-common
