#!/usr/bin/env bash

export GO111MODULE=on
export GOPROXY=https://goproxy.cn

rm -rf build
mkdir -p build/bin

cd src/github.com/sbot

go build -o sbot cmd/sbot/main.go

mv sbot ../../../build/bin

cp -rf conf ../../../build
cp -rf scripts ../../../build

cd ../../../
cp install.sh build

echo "build sbot into build dir is ok.................."
