FROM ubuntu:trusty

ENV HOME /root

RUN \
  apt-get update && \
  apt-get -y install \
    build-essential \
    wget

RUN wget -q -O cf.deb "https://cli.run.pivotal.io/stable?release=debian64&version=6.17.1&source=github-rel" && \
  dpkg -i cf.deb

RUN \
  wget https://s3.amazonaws.com/bosh-cli-artifacts/bosh-cli-0.0.124-linux-amd64 --output-document="/usr/bin/bosh" && \
  cd /usr/bin && \
  echo "5065d3c8291abab1ffef1bac6de6cea8d514e319 bosh" | sha1sum -c - && \
  chmod +x bosh

