#!/bin/bash

# If you change the default ports in .env, change them here too
docker build . --tag iperf-node && docker run \
    --env-file ./.env \
    --rm \
    -p 5001:5001 -p 8080:8080 \
    iperf-node
    