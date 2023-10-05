FROM cloudfoundry/relint-base

# bosh-cli
ARG bosh_cli_version
RUN set -eux; \
      apt-get -y install \
        build-essential \
        libreadline-dev \
        libsqlite3-dev \
        libssl-dev \
        libxml2-dev \
        libxslt1-dev \
        libyaml-dev \
        openssl \
        ruby \
        ruby-dev \
        sqlite3 \
        zlib1g-dev; \
      url="https://github.com/cloudfoundry/bosh-cli/releases/download/v${bosh_cli_version}/bosh-cli-${bosh_cli_version}-linux-amd64"; \
      wget "${url}" -O /usr/local/bin/bosh; \
      chmod +x /usr/local/bin/bosh; \
      bosh --version

# bbl and dependencies
ARG bbl_version
RUN \
  wget https://github.com/cloudfoundry/bosh-bootloader/releases/download/v${bbl_version}/bbl-v${bbl_version}_linux_amd64 -P /tmp && \
  mv /tmp/bbl-* /usr/local/bin/bbl && \
  cd /usr/local/bin && \
  chmod +x bbl

ARG terraform_version
RUN \
  wget https://releases.hashicorp.com/terraform/${terraform_version}/terraform_${terraform_version}_linux_amd64.zip -P /tmp && \
  cd /tmp && \
  unzip /tmp/terraform_${terraform_version}_linux_amd64.zip -d /tmp && \
  mv /tmp/terraform /usr/local/bin/terraform && \
  cd /usr/local/bin && \
  chmod +x terraform && \
  rm -rf /tmp/*

RUN \
  wget -q -O cf.deb \
  "https://cli.run.pivotal.io/stable?release=debian64&version=6.28.0&source=github-rel"

RUN dpkg -i cf.deb
