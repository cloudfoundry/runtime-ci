FROM golang:1.20.7-bullseye

RUN \
  apt-get update && \
  apt-get -y install \
    ack-grep \
    build-essential \
    curl \
    git \
    jq \
    python3-pip \
    ruby \
    unzip \
    wget

RUN dpkg-divert --local --divert /usr/bin/ack --rename --add /usr/bin/ack-grep

# yq
RUN \
  pip3 install yq

# Add trusted relint ca certificate
ARG RELINT_CA_CERTIFICATE
ARG PREVIOUS_RELINT_CA_CERTIFICATE
RUN echo -n "$RELINT_CA_CERTIFICATE" > /usr/local/share/ca-certificates/relint.crt && \
    echo -n "$PREVIOUS_RELINT_CA_CERTIFICATE" > /usr/local/share/ca-certificates/previous_relint.crt && \
  /usr/sbin/update-ca-certificates -f
