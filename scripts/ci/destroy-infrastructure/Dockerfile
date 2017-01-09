FROM ubuntu:trusty

RUN \
  apt-get update && \
  apt-get -y install \
    build-essential \
    curl \
    mercurial \
    software-properties-common \
    jq \
    wget

RUN \
  add-apt-repository ppa:git-core/ppa -y && \
  apt-get update && \
  apt-get install -y git unzip && \
  apt-get remove -y --purge software-properties-common

RUN \
  wget https://s3.amazonaws.com/bosh-init-artifacts/bosh-init-0.0.91-linux-amd64 -P /tmp && \
  mv /tmp/bosh-init-0.0.91-linux-amd64 /usr/local/bin/bosh-init && \
  chmod +x /usr/local/bin/bosh-init && \
  rm -rf /tmp/*

RUN \
  wget https://releases.hashicorp.com/terraform/0.8.2/terraform_0.8.2_linux_amd64.zip -P /tmp && \
  unzip /tmp/terraform_0.8.2_linux_amd64.zip -d /tmp && \
  mv /tmp/terraform /usr/local/bin/terraform && \
  chmod +x /usr/local/bin/terraform && \
  rm -rf /tmp/*

RUN wget "$(curl -s https://api.github.com/repos/cloudfoundry/bosh-bootloader/releases | jq -r .[0].assets[].browser_download_url | grep _linux)" && \
  mv bbl-* /usr/local/bin/bbl && \
  chmod +x /usr/local/bin/bbl

RUN \
  wget https://github.com/postmodern/ruby-install/archive/v0.4.1.tar.gz -P /tmp && \
  tar zxf /tmp/v0.4.1.tar.gz -C /tmp && \
  cd /tmp/ruby-install-0.4.1 && make install && \
  ruby-install ruby 2.1.8 && \
  rm -rf /usr/local/src/ruby-2.1.8 && \
  rm -rf /tmp/*

ENV PATH /opt/rubies/ruby-2.1.8/bin:$PATH
