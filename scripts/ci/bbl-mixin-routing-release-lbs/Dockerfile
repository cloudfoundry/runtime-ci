FROM golang:1.7

ENV HOME /root

RUN \
  apt-get update && \
  apt-get -y install \
    git \
    awscli \
    unzip

RUN \
  wget https://releases.hashicorp.com/terraform/0.7.13/terraform_0.7.13_linux_amd64.zip -P /tmp && \
  unzip /tmp/terraform_0.7.13_linux_amd64.zip -d /tmp && \
  mv /tmp/terraform /usr/local/bin/terraform && \
  chmod +x /usr/local/bin/terraform && \
  rm -rf /tmp/*

RUN wget https://github.com/stedolan/jq/releases/download/jq-1.5/jq-linux64 --output-document /usr/local/bin/jq && \
  chmod +x /usr/local/bin/jq

RUN wget "$(curl -s https://api.github.com/repos/cloudfoundry/bosh-bootloader/releases | jq -r .[0].assets[].browser_download_url | grep _linux)" && \
  mv bbl-* /usr/local/bin/bbl && \
  chmod +x /usr/local/bin/bbl

RUN \
  wget https://s3.amazonaws.com/bosh-cli-artifacts/bosh-cli-0.0.124-linux-amd64 --output-document="/usr/bin/bosh" && \
  cd /usr/bin && \
  echo "5065d3c8291abab1ffef1bac6de6cea8d514e319 bosh" | sha1sum -c - && \
  chmod +x bosh
