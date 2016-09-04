FROM ubuntu:trusty

ENV HOME /root

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

COPY Gemfile Gemfile.lock /

RUN \
  gem update --system && \
  gem install bundler --no-rdoc --no-ri && bundle install
