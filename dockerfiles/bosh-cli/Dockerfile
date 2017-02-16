FROM ubuntu:trusty

RUN \
  apt-get update && \
  apt-get -y install \
    git \
    unzip \
    wget

# jq
RUN \
  wget https://github.com/stedolan/jq/releases/download/jq-1.5/jq-linux64 --output-document="/usr/bin/jq" && \
  cd /usr/bin && \
  echo "d8e36831c3c94bb58be34dd544f44a6c6cb88568 jq" | sha1sum -c - && \
  chmod +x jq

# bosh-cli
RUN \
  wget https://s3.amazonaws.com/bosh-cli-artifacts/bosh-cli-2.0.3-linux-amd64 --output-document="/usr/bin/bosh" && \
  cd /usr/bin && \
  echo "845b7a4e86208d492c55183a71961836f50f5c94 bosh" | sha1sum -c - && \
  chmod +x bosh

# bbl and dependencies
RUN \
  wget https://github.com/cloudfoundry/bosh-bootloader/releases/download/v2.1.1/bbl-v2.1.1_linux_x86-64 -P /tmp && \
  mv /tmp/bbl-* /usr/local/bin/bbl && \
  cd /usr/local/bin && \
  echo "020e37c53fd9f8370c7d3f3cd69b2607ef814b45 bbl" | sha1sum -c - && \
  chmod +x bbl

RUN \
  wget https://releases.hashicorp.com/terraform/0.8.4/terraform_0.8.4_linux_amd64.zip -P /tmp && \
  cd /tmp && \
  echo "190b2d404cc30227f7dec429d1f643d92fce8072 terraform_0.8.4_linux_amd64.zip" | sha1sum -c - && \
  unzip /tmp/terraform_0.8.4_linux_amd64.zip -d /tmp && \
  mv /tmp/terraform /usr/local/bin/terraform && \
  cd /usr/local/bin && \
  chmod +x terraform && \
  rm -rf /tmp/*

RUN \
  wget https://s3.amazonaws.com/bosh-init-artifacts/bosh-init-0.0.99-linux-amd64 -P /tmp && \
  mv /tmp/bosh-init-0.0.99-linux-amd64 /usr/local/bin/bosh-init && \
  cd /usr/local/bin && \
  echo "00ccaf07a11bd8206407f83f1e606e16f3475bf3 bosh-init" | sha1sum -c - && \
  chmod +x bosh-init
