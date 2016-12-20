
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
  wget https://s3.amazonaws.com/bosh-cli-artifacts/bosh-cli-0.0.124-linux-amd64 --output-document="/usr/bin/bosh" && \
  cd /usr/bin && \
  echo "5065d3c8291abab1ffef1bac6de6cea8d514e319 bosh" | sha1sum -c - && \
  chmod +x bosh
