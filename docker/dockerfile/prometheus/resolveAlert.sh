#!/bin/bash

START_TIME=$(TZ="Japan" date +"%Y-%m-%dT%H:%M:%S%z")

alertsData='{
    "comment": "fixed",
    "createdBy": "AWS",
    "status": "resolved",
    "startsAt": "'"${START_TIME}"'",
    "endsAt": "'"${START_TIME}"'",
    "matchers": [
      {
        "name": "alertname",
        "value": "Cpu TEST internal host message",
        "isRegex": false
      }
    ]
  }'

# API v1
# curl -XPOST -d"$alertsData" http://localhost:9093/api/v1/silences
curl -XPOST -d"$alertsData" -H "Content-Type: application/json" http://localhost:9093/api/v2/silences