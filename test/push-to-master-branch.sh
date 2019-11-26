#!/usr/bin/env bash
curl -X POST localhost:1323/github \
    -H "content-type: application/json" \
    -H "x-github-event: push" \
    -d @$(dirname $0)/data/webhook/webhook-push-to-master.json 