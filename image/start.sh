#!/bin/sh

alias weave=/home/weave/weave

# setup weave
weave launch-router
weave launch-proxy --without-dns --rewrite-inspect

# add localkube to network
weave expose -h "localkube.weave.local"

# setup SkyDNS to use docker networking IP
ip=$(ip -4 addr show dev docker0 | grep -m1 -o 'inet [.0-9]*' | sed 's/inet \([.0-9]*\)/\1/')
export DNS_SERVER=$ip:1970

localkube start
