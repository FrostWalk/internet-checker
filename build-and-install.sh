#!/bin/bash

if [ "$(id -u)" -ne 0 ]; then
    echo "Error: This script must be executed as root."
    exit 1
fi

# Define the paths
SOURCE_PATH="./src"
BINARY_NAME="internet-checker"
DESTINATION_PATH="/usr/local/bin/"
SYSTEMD_PATH="./systemd/"
SYSTEMD_DESTINATION_PATH="/etc/systemd/system/"

export GOPATH=$(./src/gopath)

# Compile the Go program with maximum optimization
cd ${SOURCE_PATH} && go build -ldflags "-s -w" -o "/tmp/${BINARY_NAME}" "main.go" && cd ..

mv "/tmp/${BINARY_NAME}" ${DESTINATION_PATH}

cp "${SYSTEMD_PATH}internet-checker.service" "${SYSTEMD_DESTINATION_PATH}"

# Reload systemd services
systemctl daemon-reload

# Start and enable the service
systemctl start internet-checker
systemctl enable internet-checker