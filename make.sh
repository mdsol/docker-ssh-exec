#!/usr/bin/env bash
# runs goxc in each product directory
set -e

echo "Building docker image for docker-ssh-exec..."
docker build --no-cache=true --tag mdsol/docker-ssh-exec .
