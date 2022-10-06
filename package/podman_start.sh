echo
echo " == PODMAN START"
echo
podman run --rm -v $(pwd)/..:/src                 \
           -ti                                    \
           --name debian                          \
           -d                                     \
           mvkvl/debian-golang-systemd:latest

# in docker run
#
# cd /src && ./app_package.sh {version}
#
