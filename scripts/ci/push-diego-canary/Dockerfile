FROM ubuntu:trusty

RUN \
  apt-get update && \
  apt-get -y install \
    build-essential \
    wget

RUN \
  wget \
    -P /tmp \
    http://go-cli.s3-website-us-east-1.amazonaws.com/releases/v6.16.1/cf-cli_6.16.1_linux_x86-64.tgz && \
  tar xvf /tmp/cf-cli_6.16.1_linux_x86-64.tgz -C /usr/local/bin && \
  rm -rf /tmp/*
