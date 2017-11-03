FROM ubuntu:trusty

ENV HOME /root

RUN \
  apt-get update && \
  apt-get -y install \
    build-essential \
    wget \
    curl \
    unzip

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
  gem install bosh_cli --version 1.3215.0 --no-rdoc --no-ri

RUN \
  curl -v -L -o ./spiff.zip https://github.com/cloudfoundry-incubator/spiff/releases/download/v1.0.7/spiff_linux_amd64.zip && \
  unzip ./spiff.zip -d /usr/local/bin && \
  chmod +x /usr/local/bin/spiff && \
  rm -f ./spiff.zip
