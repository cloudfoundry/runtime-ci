FROM ubuntu:trusty

ENV HOME /root

RUN \
  apt-get update && \
  apt-get -y install \
    build-essential \
    software-properties-common \
    wget

RUN \
  add-apt-repository ppa:git-core/ppa -y && \
  apt-get update && \
  apt-get install -y git && \
  apt-get remove -y --purge software-properties-common


RUN \
  wget https://github.com/postmodern/ruby-install/archive/v0.4.1.tar.gz -P /tmp && \
  tar zxf /tmp/v0.4.1.tar.gz -C /tmp && \
  cd /tmp/ruby-install-0.4.1 && make install && \
  ruby-install ruby 2.3.1 && \
  rm -rf /usr/local/src/ruby-2.3.1 && \
  rm -rf /tmp/*

ENV PATH /opt/rubies/ruby-2.3.1/bin:$PATH
RUN \
  gem update --system && \
  gem install bundler --no-rdoc --no-ri && \
  gem install bosh_cli --version 1.3215.0 --no-rdoc --no-ri
