#!/bin/sh
git pull
cd cmd
go build -o chain_server main.go
cd ..
killall -9 chain_server
nohup ./cmd/chain_server &