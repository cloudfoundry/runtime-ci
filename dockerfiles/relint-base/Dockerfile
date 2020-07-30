FROM golang:1.14

ENV jq_version 1.6
ENV jq_checksum 056ba5d6bbc617c29d37158ce11fd5a443093949

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
    go-dep \
    python3-pip

RUN dpkg-divert --local --divert /usr/bin/ack --rename --add /usr/bin/ack-grep

RUN go get -u github.com/onsi/ginkgo/...

# jq
RUN \
  wget https://github.com/stedolan/jq/releases/download/jq-${jq_version}/jq-linux64 --output-document="/usr/bin/jq" && \
  cd /usr/bin && \
  echo "${jq_checksum} jq" | sha1sum -c - && \
  chmod +x jq

# yq
RUN \
  pip3 install yq
