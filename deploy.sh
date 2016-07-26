#!/bin/bash
export GOPATH=$(pwd)
go build -ldflags "-s" -o ./gideon ./*.go
ssh 192.168.1.39 killall gideon
scp gideon chip@192.168.1.39:.
ssh 192.168.1.39 ./start_gideon.sh
