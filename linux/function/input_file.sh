#!/usr/bin/env bash

filename='/etc/hosts'

while read myline
do
  echo "$myline"
done < $filename
