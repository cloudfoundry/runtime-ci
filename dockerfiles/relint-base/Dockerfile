FROM ubuntu:bionic

ENV go_version 1.12.12

RUN \
  apt-get update && \
  apt-get -y install \
    build-essential \
    git \
    ruby \
    unzip \
    wget \
    curl \
    ack-grep \
    go-dep

RUN dpkg-divert --local --divert /usr/bin/ack --rename --add /usr/bin/ack-grep

ENV GOPATH /go
ENV PATH /go/bin:/usr/local/go/bin:$PATH
# https://golang.org/dl/ SHA256 Checksum: 2a3fdabf665496a0db5f41ec6af7a9b15a49fbe71a85a50ca38b1f13a103aeec
RUN \
  wget https://golang.org/dl/go${go_version}.linux-amd64.tar.gz -P /tmp && \
  tar xzvf /tmp/go${go_version}.linux-amd64.tar.gz -C /usr/local && \
  mkdir ${GOPATH} && \
  go get -u github.com/onsi/ginkgo/... && \
  rm -rf /tmp/*

# jq
RUN \
  wget https://github.com/stedolan/jq/releases/download/jq-1.5/jq-linux64 --output-document="/usr/bin/jq" && \
  cd /usr/bin && \
  echo "d8e36831c3c94bb58be34dd544f44a6c6cb88568 jq" | sha1sum -c - && \
  chmod +x jq

# yq
RUN \
  wget https://github.com/mikefarah/yq/releases/download/1.14.0/yq_linux_amd64 --output-document="/usr/bin/yq" && \
  chmod +x /usr/bin/yq
