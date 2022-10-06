echo
echo " == PODMAN SHUTDOWN"
echo
podman container stop $(podman container list -aq)
podman container rm $(podman container list -aq)
#podman image rm $(podman image list -q)
podman machine stop pod
podman machine rm -f pod
