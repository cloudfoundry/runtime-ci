FROM ubuntu:trusty

ENV HOME /root

RUN \
  apt-get update && \
  apt-get -y install \
    build-essential \
    software-properties-common \
    wget \
    unzip

RUN \
  add-apt-repository ppa:git-core/ppa -y && \
  apt-get update && \
  apt-get install -y git && \
  apt-get remove -y --purge software-properties-common

RUN \
  wget https://github.com/postmodern/ruby-install/archive/v0.4.1.tar.gz -P /tmp && \
  tar zxf /tmp/v0.4.1.tar.gz -C /tmp && \
  cd /tmp/ruby-install-0.4.1 && make install && \
  ruby-install ruby 2.1.8 && \
  rm -rf /usr/local/src/ruby-2.1.8 && \
  rm -rf /tmp/*

ENV PATH /opt/rubies/ruby-2.1.8/bin:$PATH

RUN \
  wget https://s3.amazonaws.com/bosh-cli-artifacts/bosh-cli-0.0.139-linux-amd64 --output-document="/usr/bin/bosh" && \
  cd /usr/bin && \
  echo "0c8363a379aa48f58aeb828b41eb9eb92fa57e6f bosh" | sha1sum -c - && \
  chmod +x bosh

RUN wget https://github.com/stedolan/jq/releases/download/jq-1.5/jq-linux64 --output-document /usr/local/bin/jq && \
  chmod +x /usr/local/bin/jq
