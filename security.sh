#!/bin/bash

# Script for generate acl for redis and file with password for golang app

if [[ "$1" == "gen_pass" ]];
  then
    PASS=$(echo -n "$3" | sha256sum)
    echo -n -e "\nuser "$2" on ~* &* +@all #${PASS::-3}" >> ./redis/users.acl
    echo -n -e "$3" > ./app/RedisPass

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
