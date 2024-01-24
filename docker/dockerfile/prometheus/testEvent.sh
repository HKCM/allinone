#!/bin/bash

DEFAULT_TIME=60
TARGET_TIME=DEFAULT_TIME
START_TIME=$(TZ="Japan" date +"%Y-%m-%dT%H:%M:%S%z")
# Setup 1 hour maintenance window
#END_TIME=$(TZ="Japan" date -d '+1 hour' +"%Y-%m-%dT%H:%M:%S%z")
END_TIME=$(TZ="Japan" date -d "+${TARGET_TIME} minutes" +"%Y-%m-%dT%H:%M:%S%z")

alertsData='[
  {
    "labels": {
      "alertname": "TEST alert from internal host",
      "severity": "critical"
    },
    "annotations": {
      "message": "This is just a test alert"
    },
    "startsAt": "'"${START_TIME}"'",
    "endsAt": "'"${END_TIME}"'""
  }
]'
# API v1
# curl -XPOST -d"$alertsData" http://localhost:9093/api/v1/alerts
curl -XPOST -d"$alertsData" -H "Content-Type: application/json" http://localhost:9093/api/v2/alerts