#!/bin/bash

. ../shared.sh

println START FAILOVER MONITOR
docker run --rm -it --name mysql-failover \
  chaostesting/mysql-failover

