#!/usr/bin/env bash
# runs goxc in each product directory
set -e

goxc
echo "Building static linux binary for docker-ssh-exec..."
mkdir -p pkg
buildcmd='CGO_ENABLED=0 go build -a --installsuffix cgo --ldflags="-s" -o pkg/docker-ssh-exec'
docker run --rm -it -v "$GOPATH":/gopath -v "$(pwd)":/app -e "GOPATH=/gopath" \
  -w /app golang:1.6 sh -c "$buildcmd"

echo "Building docker image for docker-ssh-exec..."
docker build --no-cache=true --tag mdsol/docker-ssh-exec .
rm -f pkg/docker-ssh-exec

exit 0
