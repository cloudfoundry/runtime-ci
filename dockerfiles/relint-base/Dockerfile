FROM golang:1.20.6-bullseye

ENV jq_version 1.6
ENV jq_checksum 056ba5d6bbc617c29d37158ce11fd5a443093949

RUN \
  apt-get update && \
  apt-get -y install \
    build-essential \
    git \
    ruby \
    unzip \
    wget \
    curl \
    ack-grep \
    python3-pip

RUN dpkg-divert --local --divert /usr/bin/ack --rename --add /usr/bin/ack-grep

# jq
RUN \
  wget https://github.com/stedolan/jq/releases/download/jq-${jq_version}/jq-linux64 --output-document="/usr/bin/jq" && \
  cd /usr/bin && \
  echo "${jq_checksum} jq" | sha1sum -c - && \
  chmod +x jq

# yq
RUN \
  pip3 install yq

# Add trusted relint ca certificate
ARG RELINT_CA_CERTIFICATE
ARG PREVIOUS_RELINT_CA_CERTIFICATE
RUN echo -n "$RELINT_CA_CERTIFICATE" > /usr/local/share/ca-certificates/relint.crt && \
    echo -n "$PREVIOUS_RELINT_CA_CERTIFICATE" > /usr/local/share/ca-certificates/previous_relint.crt && \
  /usr/sbin/update-ca-certificates -f
