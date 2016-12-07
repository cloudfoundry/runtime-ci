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
  echo "5065d3c8291abab1ffef1bac6de6cea8d514e319 bosh" | sha1sum -c - && \
  chmod +x bosh

RUN \
  wget https://github.com/cloudfoundry-incubator/routing-api-cli/releases/download/2.7.0/rtr-linux-amd64.tgz --output-document="/tmp/rtr-linux-amd64.tgz" && \
  echo "24ea7aa66b2829d9b4ba0315c88cadbc9bc47e55 /tmp/rtr-linux-amd64.tgz" | sha1sum -c - && \
  tar xzf /tmp/rtr-linux-amd64.tgz -C /tmp && \
  mv /tmp/rtr-linux-amd64 /usr/bin/rtr && \
  chmod +x /usr/bin/rtr
