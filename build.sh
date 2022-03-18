#!/usr/bin/env bash

export GO111MODULE=on
export GOPROXY=https://goproxy.cn

rm -rf build
mkdir -p build/bin

cd GBWServerBot/src/github.com/sbot

go build -o sbot cmd/sbot/main.go
go build -o SbotClient cmd/client/main.go
go build -o SbotFileClient cmd/fileclient/main.go
go build -o SbotAttackPayloadClient cmd/attackpayloadclient/main.go
go build -o SbotAttackQueryClient   cmd/query/main.go
go build -o SbotAttackTargetsClient cmd/attacktargets/main.go
go build -o SbotAttackScriptsClient cmd/attackscripts/main.go


mv sbot ../../../../build/bin
mv SbotClient ../../../../build/bin
mv SbotFileClient ../../../../build/bin
mv SbotAttackPayloadClient ../../../../build/bin
mv SbotAttackQueryClient ../../../../build/bin
mv SbotAttackTargetsClient ../../../../build/bin
mv SbotAttackScriptsClient ../../../../build/bin

cp -rf conf ../../../../build
cp -rf scripts ../../../../build

cd ../../../../

cd GBWClientBot/src/github.com/cbot

go build -o cbot_linux cmd/cbotlinux/main.go
go build -o AttackDump cmd/attack/main.go
go build -o AttackHadoopIPC cmd/hadoop/attack.go

mv cbot_linux ../../../../build/bin
mv AttackDump ../../../../build/bin
mv AttackHadoopIPC ../../../../build/bin

cd ../../../../

cp install.sh build

echo "build sbot into build dir is ok.................."
