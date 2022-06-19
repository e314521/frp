#!/bin/sh
curPath=$(readlink -f "$(dirname "$0")")
cd $curPath   
export LD_LIBRARY_PATH=$LD_LIBRARY_PATH:./client/libx64
go run ./cmd/frpc/main.go -c ./conf/frpc.ini
