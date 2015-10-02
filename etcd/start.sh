#!/bin/bash

. ../shared.sh

docker run -d -v /usr/share/ca-certificates/:/etc/ssl/certs \
  --name etcd -p 4001:4001 -p 2380:2380 -p 2379:2379 \
  quay.io/coreos/etcd \
  -name etcd0 \
  -listen-client-urls http://0.0.0.0:2379,http://0.0.0.0:4001 \
  -advertise-client-urls http://172.17.42.1:2379,http://172.17.42.1:4001 \
  -initial-advertise-peer-urls http://172.17.42.1:2380 \
  -listen-peer-urls http://0.0.0.0:2380 \
  -initial-cluster etcd0=http://172.17.42.1:2380 \
  -initial-cluster-state new

