FROM ubuntu:trusty

ENV HOME /root

RUN \
  apt-get update && \
  apt-get -y install \
    git \
    unzip \
    curl \
    wget \
    software-properties-common && \
  add-apt-repository ppa:brightbox/ruby-ng -y && \
  apt-get update && \
  apt-get -y install \
    ruby2.3 \
    ruby2.3-dev && \
  apt-get remove -y --purge software-properties-common


ENV PATH /opt/rubies/ruby-2.3.1/bin:$PATH
RUN \
  gem update --system && \
  gem install bundler --no-rdoc --no-ri && \
  gem install bosh_cli --version 1.3215.0 --no-rdoc --no-ri

RUN \
  curl -v -L -o ./spiff.zip https://github.com/cloudfoundry-incubator/spiff/releases/download/v1.0.7/spiff_linux_amd64.zip && \
  unzip ./spiff.zip -d /usr/local/bin && \
  chmod +x /usr/local/bin/spiff && \
  rm -f ./spiff.zip
