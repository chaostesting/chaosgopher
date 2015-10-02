#!/bin/bash

. ../shared.sh

function start_container() {
  docker start $1 2>/dev/null || \
    docker run --name $1 -d \
      -m 2g \
      -e MYSQL_ROOT_PASSWORD=$MYSQL_ROOT_PASSWORD \
      -e MYSQL_DATABASE=$MYSQL_DATABASE \
      -p $2:$2 \
      $3
}

println STARTING CONTAINERS
start_container mysql-alpha 3306 chaostesting/mysql-alpha
start_container mysql-beta 3307 chaostesting/mysql-beta

