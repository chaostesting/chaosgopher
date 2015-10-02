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
start_container mysql-master 3307 chaostesting/mysql-master
start_container mysql-slave 3308 chaostesting/mysql-slave

