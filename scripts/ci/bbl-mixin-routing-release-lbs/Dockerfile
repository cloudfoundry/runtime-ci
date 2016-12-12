FROM golang:1.7

ENV HOME /root

RUN \
  apt-get update && \
  apt-get -y install \
    git \
    unzip

RUN \
  wget https://releases.hashicorp.com/terraform/0.7.5/terraform_0.7.5_linux_amd64.zip -P /tmp && \
  unzip /tmp/terraform_0.7.5_linux_amd64.zip -d /tmp && \
  mv /tmp/terraform /usr/local/bin/terraform && \
  chmod +x /usr/local/bin/terraform && \
  rm -rf /tmp/*

RUN wget "$(curl -s https://api.github.com/repos/cloudfoundry/bosh-bootloader/releases | jq -r .[0].assets[].browser_download_url | grep _linux)" && \
  mv bbl-* /usr/local/bin/bbl && \
  chmod +x /usr/local/bin/bbl
