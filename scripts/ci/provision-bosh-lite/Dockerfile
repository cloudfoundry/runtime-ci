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
