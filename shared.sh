#!/bin/bash

set -euo pipefail

export MYSQL_ROOT_PASSWORD=chaostestingrootpassword
export MYSQL_DATABASE=chaostesting

function println() {
  echo -e "\033[36m ### $@\033[0m"
}

if [[ ! $(docker-machine inspect chaostesting 2> /dev/null) ]]; then
  println CREATING DOCKER MACHINE
  docker-machine create -d virtualbox --virtualbox-disk-size=3072 chaostesting
fi

eval "$(docker-machine env chaostesting)"
println USING DOCKER MACHINE: "$(docker-machine active)"

DOCKER_IP=`echo $DOCKER_HOST | cut -d '/' -f 3 | cut -d ':' -f 1`

