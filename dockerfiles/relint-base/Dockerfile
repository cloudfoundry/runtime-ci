FROM ubuntu:trusty

RUN \
  apt-get update && \
  apt-get -y install \
    build-essential \
    git \
    ruby \
    unzip \
    wget \
    curl \
    ack-grep

RUN dpkg-divert --local --divert /usr/bin/ack --rename --add /usr/bin/ack-grep

ENV GOPATH /go
ENV PATH /go/bin:/usr/local/go/bin:$PATH
# https://golang.org/dl/ SHA256 Checksum: 1dfe664fa3d8ad714bbd15a36627992effd150ddabd7523931f077b3926d736d
RUN \
  wget https://dl.google.com/go/go1.11.2.linux-amd64.tar.gz -P /tmp && \
  tar xzvf /tmp/go1.11.2.linux-amd64.tar.gz -C /usr/local && \
  mkdir $GOPATH && \
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
