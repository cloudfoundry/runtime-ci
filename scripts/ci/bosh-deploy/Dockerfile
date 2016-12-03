FROM ubuntu:trusty

ENV HOME /root

RUN \
  apt-get update && \
  apt-get -y install \
    git \
    build-essential \
    wget

RUN \
  wget https://s3.amazonaws.com/bosh-cli-artifacts/bosh-cli-0.0.124-linux-amd64 --output-document="/usr/bin/bosh" && \
  cd /usr/bin && \
  echo "cf78a36dd3eaa9fd206ca3647b628b69c7be9d90 bosh" | sha1sum -c - && \
  chmod +x bosh
