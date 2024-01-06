#!/bin/bash

# For the sake of convinience we do map the docker socket inside the devcontainer
DOCKER_SOCKET="/var/run/docker-host.sock"

# Go CI Lint
curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.54.0

# Docker Go SDK
go get github.com/docker/docker/client

# Go Cobra
go get github.com/spf13/cobra

# Go Viper
go get github.com/spf13/viper

# Avahi related stuff
go get github.com/godbus/dbus/v5
go get github.com/holoplot/go-avahi
go get github.com/miekg/dns

# Generate the cobra config file
cat << EOF > ~/.cobra.yaml
author: Christian Ege <ch@ege.io>
useViper: false
license: Apache-2.0
EOF

# Use the Docker host socket
cat << EOF > config.yaml
socket: ${DOCKER_SOCKET}
EOF

# Make the Docker Socket available
GROUP_ID=$(stat -c %g $DOCKER_SOCKET)
CURRENT_USER=$(whoami)
# Add a docker group with the right group id
sudo groupadd -g ${GROUP_ID} docker

# Add the dev user to this group
sudo adduser ${CURRENT_USER} docker
# Make the group available in the current shell
newgrp docker
