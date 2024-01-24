#!/bin/bash
set -e
if [ $# -eq 0 ];then
  echo "No addition"
else
  echo "user say: $@ "
fi

if [ "$ENV" = 'DEV' ];then
  echo "Hi there, this is Environment variable ${EnVariable}"
  echo "Run in Development Server"
  nginx -g "daemon off;"
else
  echo "Hi there, this is Environment variable ${EnVariable}"
  echo "Run in Production Server"
  nginx -g "daemon off;"
fi
