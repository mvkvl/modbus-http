#!/usr/bin/env bash

DIR="$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )"
cd ${DIR}

VERSION=$1
if [ -z "$VERSION" ]; then
  VERSION=$(cat "${DIR}/../VERSION")
else
  shift
fi

./podman_init.sh            && \
./podman_build.sh           && \
./podman_start.sh           && \
podman exec -ti debian bash -c "cd /src/package && ./app_package.sh $VERSION" && \
./podman_shutdown.sh        && \
echo ""