#!/bin/bash

set -e

RUN_PATH="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd $RUN_PATH

echo ----[ Deploy in IPFS ]----
DEB_FILE=ipfs-gateway-limited_${VERSION}_amd64.deb
DEB_PATH=$RUN_PATH/build/debian_out/ipfs-gateway-limited
acilia ipfs publish add deploy.foilen.com /ipfs-gateway-limited/$DEB_FILE $DEB_PATH/../$DEB_FILE
