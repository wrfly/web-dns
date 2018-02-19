#!/bin/bash

SERVER="http://localhost:8080"
HOST="www.google.com
github.com
kfd.me"

TYPE="
A
AAAA
MX
ANY"

for host in `echo "$HOST"`; do
    for typ in `echo "$TYPE"`; do
        cmd="curl -s $SERVER/$host/$typ"
        echo `$cmd`
    done
done
