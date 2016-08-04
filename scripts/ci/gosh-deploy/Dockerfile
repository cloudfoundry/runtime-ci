FROM ubuntu:trusty

ENV HOME /root

RUN \
  apt-get update && \
  apt-get -y install \
    build-essential \
    mercurial \
    software-properties-common \
    wget

RUN \
  add-apt-repository ppa:git-core/ppa -y && \
  apt-get update && \
  apt-get install -y git && \
  apt-get remove -y --purge software-properties-common

RUN \
  wget https://s3.amazonaws.com/bosh-cli-artifacts/bosh-cli-0.0.18-linux-amd64 --output-document="/usr/bin/bosh"

RUN \
  chmod +x /usr/bin/bosh

ENV GOPATH $HOME/go
ENV PATH $HOME/go/bin:/usr/local/go/bin:$PATH

RUN \
  wget https://storage.googleapis.com/golang/go1.6.2.linux-amd64.tar.gz -P /tmp && \
  tar xzvf /tmp/go1.6.2.linux-amd64.tar.gz -C /usr/local && \
  mkdir $GOPATH && \
  rm -rf /tmp/*

RUN \
  wget https://s3.amazonaws.com/bosh-init-artifacts/bosh-init-0.0.91-linux-amd64 -P /tmp && \
  mv /tmp/bosh-init-0.0.91-linux-amd64 /usr/local/bin/bosh-init && \
  chmod +x /usr/local/bin/bosh-init && \
  rm -rf /tmp/*

RUN \
  go get github.com/pivotal-cf-experimental/bosh-bootloader/bbl && \
  cd $GOPATH/src/github.com/pivotal-cf-experimental/bosh-bootloader/bbl && \
  git checkout develop && \
  go install github.com/pivotal-cf-experimental/bosh-bootloader/bbl
