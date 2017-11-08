FROM relintdockerhubpushbot/relint-base

# bosh-cli
RUN \
  wget https://s3.amazonaws.com/bosh-cli-artifacts/bosh-cli-2.0.36-linux-amd64 --output-document="/usr/bin/bosh" && \
  cd /usr/bin && \
  echo "61033a0869547d687c09ea866fe4f6ce16e020fd bosh" | sha1sum -c - && \
  chmod +x bosh

# bbl and dependencies
RUN \
  wget https://github.com/cloudfoundry/bosh-bootloader/releases/download/v4.10.4/bbl-v4.10.4_linux_x86-64 -P /tmp && \
  mv /tmp/bbl-* /usr/local/bin/bbl && \
  cd /usr/local/bin && \
  echo "0393afe06f5ca1664553895f41c0f3d6186f191b bbl" | sha1sum -c - && \
  chmod +x bbl

RUN \
  wget https://releases.hashicorp.com/terraform/0.9.7/terraform_0.9.7_linux_amd64.zip -P /tmp && \
  cd /tmp && \
  echo "f7d9d2110acccb5b24fec86f2651f3f07af493c9 terraform_0.9.7_linux_amd64.zip" | sha1sum -c - && \
  unzip /tmp/terraform_0.9.7_linux_amd64.zip -d /tmp && \
  mv /tmp/terraform /usr/local/bin/terraform && \
  cd /usr/local/bin && \
  chmod +x terraform && \
  rm -rf /tmp/*

RUN \
  wget https://s3.amazonaws.com/bosh-init-artifacts/bosh-init-0.0.99-linux-amd64 -P /tmp && \
  mv /tmp/bosh-init-0.0.99-linux-amd64 /usr/local/bin/bosh-init && \
  cd /usr/local/bin && \
  echo "00ccaf07a11bd8206407f83f1e606e16f3475bf3 bosh-init" | sha1sum -c - && \
  chmod +x bosh-init

RUN \
  wget -q -O cf.deb \
  "https://cli.run.pivotal.io/stable?release=debian64&version=6.28.0&source=github-rel"

RUN dpkg -i cf.deb
