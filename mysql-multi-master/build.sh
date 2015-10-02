#!/bin/bash

. ../shared.sh

println BUILDING DOCKER IMAGES
docker build -t chaostesting/mysql-alpha -f alpha .
docker build -t chaostesting/mysql-beta -f beta .

