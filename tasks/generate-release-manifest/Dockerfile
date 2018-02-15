FROM ubuntu:trusty

ENV HOME /root

RUN \
  apt-get update && \
  apt-get -y install \
    build-essential \
    curl \
    unzip

RUN \
  curl -v -L -o ./spiff.zip https://github.com/cloudfoundry-incubator/spiff/releases/download/v1.0.7/spiff_linux_amd64.zip && \
  unzip ./spiff.zip -d /usr/local/bin && \
  chmod +x /usr/local/bin/spiff && \
  rm -f ./spiff.zip
