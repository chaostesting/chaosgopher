#!/bin/bash

. ../shared.sh

export GOPATH=$PWD/vendor:$PWD
go get -v -d app/...
go install app
bin/app
