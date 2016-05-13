FROM ubuntu:trusty

RUN \
  apt-get update && \
  apt-get -y install \
    wget \
    python-software-properties \
    python-pip \
    jq

RUN \
  wget https://releases.hashicorp.com/vagrant/1.8.1/vagrant_1.8.1_x86_64.deb -P /tmp && \
  dpkg --install /tmp/vagrant_1.8.1_x86_64.deb && \
  vagrant plugin install vagrant-aws && \
  rm -rf /tmp/*

RUN \
  pip install awscli==1.6.6

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
  gem install bundler --no-rdoc --no-ri -v 1.11.2 && \
  # -v 1.11.2 is a workaround for: https://github.com/bundler/bundler/issues/4494
  gem install bosh_cli --version 1.3215.0 --no-rdoc --no-ri
