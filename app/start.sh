#!/bin/bash

. ../shared.sh

docker run --rm --name app -p 8080:8080 app ./run.sh
