#!/bin/sh

echo "haldo"

alias weave=/home/weave/weave

# setup networking
weave launch
eval $(weave env)

/opt/localkube/localkube start
