#!/bin/bash

# Script for generate acl for redis and file with password for golang app

if [[ "$1" == "gen_pass" ]];
  then
    PASS=$(echo -n "$2" | sha256sum)
    echo -n "$2" > golangpass
    echo -n "user admin +@all on #${PASS::-3}" > redisacl.acl
elif [[ "$1" == "gen_certs" ]];
    then
      echo "Not work yet!"
elif [[ "$1" == "help" ]];
    then
      echo -e "For generate certs use: \n \
      ./security.sh gen_certs \n \
      For generate password use: \n \
      ./security.sh gen_pass your_secret_pass"
else
  echo -e "Not correct pamameters! Use: \n \
  ./security.sh help"
fi
