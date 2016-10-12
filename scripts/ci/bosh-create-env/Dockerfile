FROM ubuntu:trusty

ENV HOME /root

RUN \
  apt-get update && \
  apt-get -y install \
    build-essential \
    wget \
    unzip

RUN \
  wget https://github.com/postmodern/ruby-install/archive/v0.4.1.tar.gz -P /tmp && \
  tar zxf /tmp/v0.4.1.tar.gz -C /tmp && \
  cd /tmp/ruby-install-0.4.1 && make install && \
  ruby-install ruby 2.1.8 && \
  rm -rf /usr/local/src/ruby-2.1.8 && \
  rm -rf /tmp/*

ENV PATH /opt/rubies/ruby-2.1.8/bin:$PATH

RUN \
  wget https://s3.amazonaws.com/bosh-cli-artifacts/bosh-cli-0.0.63-linux-amd64 --output-document="/usr/bin/bosh"

RUN \
  chmod +x /usr/bin/bosh

RUN \
  export CLOUD_SDK_REPO="cloud-sdk-$(lsb_release -c -s)" && \
  echo "deb http://packages.cloud.google.com/apt $CLOUD_SDK_REPO main" | sudo tee /etc/apt/sources.list.d/google-cloud-sdk.list && \
  curl https://packages.cloud.google.com/apt/doc/apt-key.gpg | sudo apt-key add - && \
  sudo apt-get update && sudo apt-get install -y google-cloud-sdk
