#!/bin/sh

# Requirements
# - programs: wget, xcputranslate, go
# - argument $1: program name
# - environment variable ${TARGETPLATFORM}

# Output
# - Clear current directory
# - Binary program to /tmp/bin

echo "Marking ${TARGETPLATFORM} as unavailable"
echo "Clearing the current directory..."
rm -rf *
echo "Building unavailable program..."
wget -q https://raw.githubusercontent.com/qdm12/binpot/main/unavailable/main.go
GOARCH="$(xcputranslate translate -field arch -targetplatform ${TARGETPLATFORM})" \
GOARM="$(xcputranslate translate -field arm -targetplatform ${TARGETPLATFORM})" \
go build -trimpath \
  -ldflags="-s -w \
  -X 'main.name=${1}' \
  -X 'main.platform=${2}'" \
  -o /tmp/bin main.go
chmod 500 /tmp/bin
