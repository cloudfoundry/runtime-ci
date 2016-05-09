FROM ubuntu:trusty

RUN \
  echo "deb http://archive.ubuntu.com/ubuntu trusty-backports main restricted universe multiverse" > /etc/apt/sources.list.d/backports.list && \
  apt-get update && \
  apt-get -y install \
    build-essential \
    wget \
    software-properties-common \
    jq/trusty-backports \
    curl \
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
  ruby-install ruby 2.3.1 && \
  rm -rf /usr/local/src/ruby-2.3.1 && \
  rm -rf /tmp/*

ENV PATH /opt/rubies/ruby-2.3.1/bin:$PATH
RUN \
  gem update --system && \
  gem install bundler --no-rdoc --no-ri

RUN \
  curl -v -L -o ./spiff.zip https://github.com/cloudfoundry-incubator/spiff/releases/download/v1.0.7/spiff_linux_amd64.zip && \
  unzip ./spiff.zip -d /usr/local/bin && \
  chmod +x /usr/local/bin/spiff && \
  rm -f ./spiff.zip
