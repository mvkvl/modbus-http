#!/usr/bin/env bash

# go tool dist list

# GOOS
#Linux 	        linux
#MacOS X 	      darwin
#Windows 	      windows
#FreeBSD 	      freebsd
#NetBSD 	      netbsd
#OpenBSD 	      openbsd
#DragonFly BSD 	dragonfly
#Plan 9 	      plan9
#Native Client 	nacl
#Android 	      android

# GOARCH
#x386 	                  386
#AMD64 	                  amd64
#AMD64 с 32-указателями 	amd64p32
#ARM 	                    arm
#ARM 	                    arm64

# GOARM
# armel (softfloat)               GOARM=5
# armhf (hardware floating point) GOARM=6 / GOARM=7

DIR="$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )"

DSTPATH=$1
if [ -z $DSTPATH ]; then
  DSTPATH="${DIR}/distr"
fi

function cleanup {
  rm -rf ${DSTPATH}/bin || true
}

function build() {
  mkdir -p ${DSTPATH}/bin/$4/$5/
  cd ${DIR}/../src
  GOOS=$1 GOARCH=$2 GOARM=$3 go build -ldflags "-s -w" -o ${DSTPATH}/bin/$4/$5/mbridge .
}

echo
echo "BUILDING mbridge"
echo

cleanup                                     && \
build linux   amd64   ""    linux   amd64   && \
build darwin  arm64   ""    macos   amd64   && \
build darwin  amd64   ""    macos   m1      && \
build linux   arm     7     linux   armhf   && \
build linux   arm64   7     linux   armhf64 && \
echo "Build complete"

