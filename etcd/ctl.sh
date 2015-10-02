#!/bin/bash

. ../shared.sh

docker run --rm etcdctl -C http://172.17.42.1:2379 "$@"
