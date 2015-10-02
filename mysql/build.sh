#!/bin/bash

. ../shared.sh

println BUILDING DOCKER IMAGES
docker build -t chaostesting/mysql-master -f master .
docker build -t chaostesting/mysql-slave -f slave .
docker build -t chaostesting/mysql-failover -f failover .

