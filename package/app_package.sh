#!/usr/bin/env bash

DIR="$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )"

VERSION=$1
if [ -z "$VERSION" ]; then
  VERSION=$(cat "${DIR}/../VERSION")
else
  shift
fi

./app_build.sh

echo
echo "PACKAGING mbridge:${VERSION}"
echo


DSTPATH="${DIR}/distr"
SRCDIR="$DSTPATH/bin/linux"

# https://www.internalpointers.com/post/build-binary-deb-package-practical-guide
function package_deb() {
  ARCH=$1
  DSTPATH=$2
  DSTDIR="${DSTPATH}/mbridge_${VERSION}_${ARCH}"
  mkdir -p "$DSTDIR/DEBIAN"
  install -d "${DSTDIR}/etc/mbridge"
  install -d "${DSTDIR}/usr/local/bin"
  install -d "${DSTDIR}/etc/systemd/system"
  install -m 0644 "${DIR}/debian/mbridge.service"  "${DSTDIR}/etc/systemd/system/mbridge.service"
  install -m 0644 "${DIR}/../conf/channels.json"     "${DSTDIR}/etc/mbridge/channels.json"
  install -m 0644 "${DIR}/../conf/logger.json"       "${DSTDIR}/etc/mbridge/logger.json"
  install -m 0755 "${SRCDIR}/${ARCH}/mbridge" "${DSTDIR}/usr/local/bin/mbridge"
  {
    echo "Package: mbridge"
    echo "Version: ${VERSION}"
    echo "Architecture: ${ARCH}"
    echo "Maintainer: Mikhail Kantur <mkantur@gmail.com>"
    echo "Description: modbus-to-http bridging service"
    echo "Depends: systemd"
    echo ""
  } > "$DSTDIR/DEBIAN/control"
  install -m 755 "${DIR}/debian/mbridge.postinst" "$DSTDIR/DEBIAN/postinst"
  install -m 755 "${DIR}/debian/mbridge.postrm"   "$DSTDIR/DEBIAN/postrm"
  dpkg-deb --build --root-owner-group "${DSTDIR}"
  rm -rf "${DSTDIR}"
}

package_deb amd64   "${DSTPATH}"
package_deb armhf   "${DSTPATH}"
package_deb armhf64 "${DSTPATH}"

