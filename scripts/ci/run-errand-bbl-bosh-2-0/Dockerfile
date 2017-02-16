FROM ubuntu:trusty

RUN \
  apt-get update && \
  apt-get -y install \
    build-essential \
    wget

# bosh-cli
RUN \
  wget https://s3.amazonaws.com/bosh-cli-artifacts/bosh-cli-2.0.3-linux-amd64 --output-document="/usr/bin/bosh" && \
  cd /usr/bin && \
  echo "845b7a4e86208d492c55183a71961836f50f5c94 bosh" | sha1sum -c - && \
  chmod +x bosh
