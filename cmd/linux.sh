#!/bin/bash

# shellcheck disable=SC2034
IP_ADDR="45.90.35.231"

cd ../
export GOOS=linux
go build .
export GOOS=windows

#ssh "root@$IP_ADDR" "service linksparser stop"
#scp -r linksparser "root@$IP_ADDR:/var/www/html"
#ssh -i "~/.ssh/id_rsa" "root@$IP_ADDR" "service linksparser restart"

sleep 15