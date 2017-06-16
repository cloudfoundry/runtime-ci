FROM ubuntu:trusty

RUN \
  apt-get update && \
  apt-get -y install \
    build-essential \
    wget \
    curl \
    openssh-client \
    software-properties-common \
    unzip

RUN \
  add-apt-repository ppa:git-core/ppa -y && \
  apt-get update && \
  apt-get install -y git && \
  apt-get remove -y --purge software-properties-common

ENV GOPATH /go
ENV PATH /go/bin:/usr/local/go/bin:$PATH
RUN \
  wget https://storage.googleapis.com/golang/go1.7.4.linux-amd64.tar.gz -P /tmp && \
  tar xzvf /tmp/go1.7.4.linux-amd64.tar.gz -C /usr/local && \
  mkdir $GOPATH && \
  rm -rf /tmp/*

# Install the cf CLI
RUN wget -q -O cf.deb "https://cli.run.pivotal.io/stable?release=debian64&version=6.22.2&source=github-rel" && \
  dpkg -i cf.deb

# Install the container networking CLI plugin
RUN wget -q -O /tmp/network-policy-plugin "https://github.com/cloudfoundry-incubator/netman-release/releases/download/v0.5.0/network-policy-plugin-linux64" && \
  chmod +x /tmp/network-policy-plugin && \
  cf install-plugin /tmp/network-policy-plugin -f && \
  rm -rf /tmp/*

RUN \
  wget https://github.com/cloudfoundry-incubator/routing-api-cli/releases/download/2.9.0/rtr-linux-amd64.tgz --output-document="/tmp/rtr-linux-amd64.tgz" && \
  echo "a702c7fc3316881387b2d333f7a4b0d01b6552c2 /tmp/rtr-linux-amd64.tgz" | sha1sum -c - && \
  tar xzf /tmp/rtr-linux-amd64.tgz -C /tmp && \
  mv /tmp/rtr-linux-amd64 /usr/bin/rtr && \
  chmod +x /usr/bin/rtr
