echo
echo " == PODMAN INIT"
echo
podman machine init --cpus 1 --disk-size 10 -m 4096 --volume /opt/work pod && \
podman machine start pod
