#!/bin/bash
#create a docker image for each test, to have all containers play each other through a script
f=0
for i in test_programs/*; do
  FILE_NAME=`echo "$i" | cut -d'.' -f1`
  OBJECT=`echo "$FILE_NAME" | cut -d'/' -f2`
  docker run -dit -v "$PWD":/usr/src/client -w /usr/src/client --name $OBJECT golang:1.6 bash -c "go build -v && ./client $OBJECT python $i"
  ((f+=1))
done
