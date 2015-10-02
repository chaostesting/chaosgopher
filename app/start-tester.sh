#!/bin/bash

. ../shared.sh

docker run --rm -v /var/run/docker.sock:/var/run/docker.sock \
  --name tester app ./run-tester.sh
