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
# armhf (hardware floatin point)  GOARM=6 / GOARM=7
function cleanup {
  rm -rf distr || true
}
function build_for_mac_m1() {
  mkdir -p distr/mac/m1
  GOOS=darwin GOARCH=arm64 go build -o distr/mac/m1/mbridge .
}
function build_for_mac_x86_64() {
  mkdir -p distr/mac/amd64
  GOOS=darwin GOARCH=amd64 go build -o distr/mac/amd64/mbridge .
}
function build_for_linux_amd64() {
  mkdir -p distr/linux/amd64
  GOOS=linux GOARCH=amd64 go build -o distr/linux/amd64/mbridge .
}
function build_for_linux_armhf() {
  mkdir -p distr/linux/armhf
  GOOS=linux GOARCH=arm GOARM=7 go build -o distr/linux/armhf/mbridge .
}
function build_for_linux_armhf_x64() {
  mkdir -p distr/linux/armhf_x64
  GOOS=linux GOARCH=arm64 GOARM=7 go build -o distr/linux/armhf_x64/mbridge .
}

cleanup                   && \
build_for_mac_m1          && \
build_for_mac_x86_64      && \
build_for_linux_amd64     && \
build_for_linux_armhf     && \
build_for_linux_armhf_x64 && \
echo "Done!"