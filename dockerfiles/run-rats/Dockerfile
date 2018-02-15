FROM relintdockerhubpushbot/cats-concourse-task

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
