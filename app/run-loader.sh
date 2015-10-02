#!/bin/bash

. ../shared.sh

export GOPATH=$PWD/vendor:$PWD
go get -v -d loader/...
go install loader
bin/loader
