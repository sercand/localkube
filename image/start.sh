#!/bin/sh

alias weave=/home/weave/weave

# setup weave
weave launch-router
weave launch-proxy --without-dns --rewrite-inspect

# add localkube to network
weave expose -h "localkube.weave.local"

localkube start
