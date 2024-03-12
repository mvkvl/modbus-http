#!/usr/bin/env bash

PROGRAM=mbridge

DIR="$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )"
cd ${DIR}/..

VERSION=$1
if [ -z "$VERSION" ]; then
  VERSION=$(cat "build/VERSION")
else
  shift
fi

./build/build.sh

echo
echo "PACKAGING ${PROGRAM}:${VERSION}"
echo


DSTPATH="${DIR}/../distr"
SRCDIR="$DSTPATH/bin"

# https://www.internalpointers.com/post/build-binary-deb-package-practical-guide
#function package_deb() {
#  ARCH=$1
#  DSTPATH=$2
#  DSTDIR="${DSTPATH}/keito_${VERSION}_${ARCH}"
#  mkdir -p "$DSTDIR/DEBIAN"
#  install -d "${DSTDIR}/usr/local/bin"
#  install -m 0755 "${SRCDIR}/linux/${ARCH}/keito" "${DSTDIR}/usr/local/bin/keito"
#  {
#    echo "Package: keito"
#    echo "Version: ${VERSION}"
#    echo "Architecture: ${ARCH}"
#    echo "Maintainer: maintainer..."
#    echo "Description: key & token generator"
#    echo ""
#  } > "$DSTDIR/DEBIAN/control"
#  install -m 755 "${DIR}/debian/keito.postinst" "$DSTDIR/DEBIAN/postinst"
#  install -m 755 "${DIR}/debian/keito.postrm"   "$DSTDIR/DEBIAN/postrm"
#  dpkg-deb --build --root-owner-group "${DSTDIR}"
#  rm -rf "${DSTDIR}"
#}
function package_tgz() {
  OS=$1
  ARCH=$2
  DSTPATH=$3
  DSTFILE="${PROGRAM}_${VERSION}_${OS}_${ARCH}.tgz"
  cd "${SRCDIR}/${OS}/${ARCH}"                    && \
  tar cvfz "${DSTFILE}" ${PROGRAM}                && \
  mv "${DSTFILE}" "${DSTPATH}"
}
function package_zip() {
  OS=$1
  ARCH=$2
  DSTPATH=$3
  DSTFILE="${PROGRAM}_${VERSION}_${OS}_${ARCH}.zip"
  cd  "${SRCDIR}/${OS}/${ARCH}"                   && \
  zip "${DSTFILE}" ${PROGRAM}.exe                 && \
  mv  "${DSTFILE}" "${DSTPATH}"
}
#package_deb amd64         "${DSTPATH}"
#package_deb armhf         "${DSTPATH}"
#package_deb armhf64       "${DSTPATH}"
package_tgz linux amd64   "${DSTPATH}"
package_tgz linux armhf   "${DSTPATH}"
package_tgz linux armhf64 "${DSTPATH}"
package_tgz macos x86_64  "${DSTPATH}"
package_tgz macos aarch64 "${DSTPATH}"
package_zip win   x86_64  "${DSTPATH}"
package_zip win   x86     "${DSTPATH}"
