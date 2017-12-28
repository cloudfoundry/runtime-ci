FROM relintdockerhubpushbot/relint-base

# bosh-cli
ARG bosh_cli_version
RUN \
  wget https://s3.amazonaws.com/bosh-cli-artifacts/bosh-cli-${bosh_cli_version}-linux-amd64 --output-document="/usr/bin/bosh" && \
  cd /usr/bin && \
  chmod +x bosh

# bbl and dependencies
RUN \
  wget https://github.com/cloudfoundry/bosh-bootloader/releases/download/v5.10.1/bbl-v5.10.1_linux_x86-64 -P /tmp && \
  mv /tmp/bbl-* /usr/local/bin/bbl && \
  cd /usr/local/bin && \
  echo "cab73790adc8f2a508cd580fd2064837220562b1523e386600cef419f2220ffe bbl" | sha256sum -c - && \
  chmod +x bbl

RUN \
  wget https://releases.hashicorp.com/terraform/0.11.1/terraform_0.11.1_linux_amd64.zip -P /tmp && \
  cd /tmp && \
  echo "4e3d5e4c6a267e31e9f95d4c1b00f5a7be5a319698f0370825b459cb786e2f35 terraform_0.11.1_linux_amd64.zip" | sha256sum -c - && \
  unzip /tmp/terraform_0.11.1_linux_amd64.zip -d /tmp && \
  mv /tmp/terraform /usr/local/bin/terraform && \
  cd /usr/local/bin && \
  chmod +x terraform && \
  rm -rf /tmp/*

RUN \
  wget -q -O cf.deb \
  "https://cli.run.pivotal.io/stable?release=debian64&version=6.28.0&source=github-rel"

RUN dpkg -i cf.deb
