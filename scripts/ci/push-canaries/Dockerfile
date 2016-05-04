FROM ubuntu:trusty

RUN \
  apt-get update && \
  apt-get -y install \
    build-essential \
    wget

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
  gem install bundler --no-rdoc --no-ri

RUN \
  wget \
    -P /tmp \
    http://go-cli.s3-website-us-east-1.amazonaws.com/releases/v6.16.1/cf-cli_6.16.1_linux_x86-64.tgz && \
  tar xvf /tmp/cf-cli_6.16.1_linux_x86-64.tgz -C /usr/local/bin && \
  rm -rf /tmp/*
