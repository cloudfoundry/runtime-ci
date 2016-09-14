FROM ubuntu:trusty

RUN \
  apt-get update && \
  apt-get -y install \
    build-essential \
    mercurial \
    software-properties-common \
    wget \
    curl

RUN \
  add-apt-repository ppa:git-core/ppa -y && \
  apt-get update && \
  apt-get install -y git && \
  apt-get remove -y --purge software-properties-common

ENV PATH /usr/local/go/bin:$PATH
RUN \
  wget https://storage.googleapis.com/golang/go1.6.3.linux-amd64.tar.gz -P /tmp && \
  tar xzvf /tmp/go1.6.3.linux-amd64.tar.gz -C /usr/local && \
  rm -rf /tmp/*

RUN wget -q -O cf.deb "https://cli.run.pivotal.io/stable?release=debian64&version=6.17.1&source=github-rel" && \
  dpkg -i cf.deb
