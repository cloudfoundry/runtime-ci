FROM golang:1.8

RUN \
  apt-get update && \
  apt-get -y install \
    curl

RUN \
  go get github.com/cloudfoundry-incubator/check-a-record
