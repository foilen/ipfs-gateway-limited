#!/bin/bash

set -e

RUN_PATH="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd $RUN_PATH

echo ----[ Create .deb ]----
DEB_FILE=ipfs-gateway-limited_${VERSION}_amd64.deb
DEB_PATH=$RUN_PATH/build/debian_out/ipfs-gateway-limited
rm -rf $DEB_PATH
mkdir -p $DEB_PATH $DEB_PATH/DEBIAN/ $DEB_PATH/usr/bin/

cat > $DEB_PATH/DEBIAN/control << _EOF
Package: ipfs-gateway-limited
Version: $VERSION
Maintainer: Foilen
Architecture: amd64
Description: front your local gateway with this light reverse proxy that translates your host to IPFS paths and only exposes those you need
_EOF

cat > $DEB_PATH/DEBIAN/postinst << _EOF
#!/bin/bash

set -e
_EOF
chmod +x $DEB_PATH/DEBIAN/postinst

cp -rv build/bin/* $DEB_PATH/usr/bin/

cd $DEB_PATH/..
dpkg-deb --no-uniform-compression --build ipfs-gateway-limited
mv ipfs-gateway-limited.deb $DEB_FILE
