FROM ubuntu:trusty

ENV HOME /root

RUN \
  apt-get update && \
  apt-get -y install \
    build-essential \
    wget

RUN wget -q -O cf.deb "https://cli.run.pivotal.io/stable?release=debian64&version=6.17.1&source=github-rel" && \
  dpkg -i cf.deb

