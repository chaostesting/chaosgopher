#!/bin/bash

. ../shared.sh

export GOPATH=$PWD/vendor:$PWD
go get -v -d tester/...
go install tester
bin/tester
