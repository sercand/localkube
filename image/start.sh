#!/bin/sh

alias weave=/home/weave/weave

# setup networking
weave launch-router
weave launch-proxy --without-dns --rewrite-inspect

localkube start
