#!/bin/bash

set -e

RUN_PATH="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd $RUN_PATH

DOMAIN_NAME=deploy.foilen.com
IPFS_ROOT=/com.foilen.deploy

echo ----[ Fetch IPFS root if missing ]----
IPFS_CURRENT_PATH=$(ipfs resolve /ipns/$DOMAIN_NAME)
IPFS_ROOT_DIR_PATH=/ipfs/$(ipfs files stat $IPFS_ROOT 2> /dev/null | head -n 1)
if [ "$IPFS_ROOT_DIR_PATH" == "/ipfs/" ]; then
    echo ----[ IPFS root $IPFS_ROOT is missing. Getting it ]----
    ipfs files cp $IPFS_CURRENT_PATH $IPFS_ROOT
else
    echo ----[ IPFS root $IPFS_ROOT is present ]----
    echo Current $DOMAIN_NAME resolves to
    echo $IPFS_CURRENT_PATH
    echo Current root $IPFS_ROOT resulves to $IPFS_ROOT_DIR_PATH
    echo $IPFS_ROOT_DIR_PATH
fi

echo ----[ Add to IPFS ]----
DEB_FILE=ipfs-gateway-limited_${VERSION}_amd64.deb
DEB_PATH=$RUN_PATH/build/debian_out/ipfs-gateway-limited
IPFS_FILE_ID=$(ipfs add -q $DEB_PATH/../$DEB_FILE | tail -n1)
echo IPFS_FILE_ID: $IPFS_FILE_ID

echo ----[ Put to IPFS under $IPFS_ROOT/ipfs-gateway-limited/$DEB_FILE ]----
ipfs files cp /ipfs/$IPFS_FILE_ID $IPFS_ROOT/ipfs-gateway-limited/$DEB_FILE

echo; echo
echo ----[ New dnslink ]----
IPFS_ROOT_DIR_ID=$(ipfs files stat $IPFS_ROOT | head -n 1)
echo You can update your DNS Link for $DOMAIN_NAME to 
echo dnslink=/ipfs/$IPFS_ROOT_DIR_ID
echo
