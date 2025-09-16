FROM golang:1.25-bookworm

RUN set -eux; \
      apt-get update; \
      apt-get -y install \
        ack \
        jq

# Add trusted relint ca certificate
ARG RELINT_CA_CERTIFICATE
ARG PREVIOUS_RELINT_CA_CERTIFICATE
RUN set -e; \
      echo -n "${RELINT_CA_CERTIFICATE}" > /usr/local/share/ca-certificates/relint.crt; \
      echo -n "${PREVIOUS_RELINT_CA_CERTIFICATE}" > /usr/local/share/ca-certificates/previous_relint.crt; \
      /usr/sbin/update-ca-certificates -f
