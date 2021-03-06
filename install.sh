#!/usr/bin/env bash
installDir=/opt/bot/sbot
cbotStoreDir=/opt/data/store/cbot
attackFileStoreDir=/opt/data/store/attack/tasks
cbotAttackSourceScriptStoreDir=/opt/data/store/source/script
cbotAttackScriptStoreDir=/opt/data/store/attack/script
confStoreDir=/opt/data/store/conf/
payloadsStoreDir=/opt/data/store/payloads

mkdir -p $installDir
mkdir -p $cbotStoreDir
mkdir -p $attackFileStoreDir
mkdir -p $cbotAttackScriptStoreDir
mkdir -p $cbotAttackSourceScriptStoreDir
mkdir -p $confStoreDir
mkdir -p $payloadsStoreDir

cp -fr bin $installDir
cp -fr conf $installDir
cp -fr conf/* $confStoreDir

cp -fr scripts/setup/* $cbotStoreDir
cp -fr scripts/jar/* $cbotStoreDir
cp -fr scripts/source/* $cbotAttackSourceScriptStoreDir
cp -fr scripts/attack/* $cbotAttackScriptStoreDir

cp -fr bin/cbot_linux $cbotStoreDir

chmod a+xrw $installDir/bin -R
chmod a+xrw $confStoreDir -R
chmod a+rw $cbotStoreDir -R
chmod a+rw $attackFileStoreDir -R
chmod a+rw $cbotAttackScriptStoreDir -R
chmod a+rw $cbotAttackSourceScriptStoreDir -R
