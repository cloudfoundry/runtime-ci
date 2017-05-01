FROM relintdockerhubpushbot/bosh-cli

RUN \
  curl -v -L -o ./spiff.zip https://github.com/cloudfoundry-incubator/spiff/releases/download/v1.0.8/spiff_linux_amd64.zip && \
  unzip ./spiff.zip -d /usr/local/bin && \
  chmod +x /usr/local/bin/spiff && \
  rm -f ./spiff.zip
