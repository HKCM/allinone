#!/usr/bin/env bash
if read -s -t 5 -p "Please enter your name: " name
then
  echo "Hello $name, welcome to my script"
else
  echo 
  echo "Sorry too slow"
  exit 1
fi