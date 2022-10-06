#!/usr/bin/env bash

VERSION=$1
if [ -z $VERSION ]; then
  VERSION=0.0.1
else
  shift
fi

./podman_init.sh            && \
./podman_build.sh           && \
./podman_start.sh           && \
podman exec -ti debian bash -c "cd /src/package && ./app_package.sh $VERSION" && \
./podman_shutdown.sh
