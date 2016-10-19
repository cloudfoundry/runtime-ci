FROM ubuntu:trusty

ENV HOME /root

RUN \
  apt-get update && \
  apt-get -y install \
    wget

RUN \
  wget https://s3.amazonaws.com/bosh-cli-artifacts/bosh-cli-0.0.85-linux-amd64 --output-document="/usr/bin/bosh"

RUN \
  chmod +x /usr/bin/bosh
