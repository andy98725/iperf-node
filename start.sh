#!/bin/bash

# If you change the default ports in .env, change them here too
docker build . --tag iperf-node && docker run \
    --env-file ./.env \
    --dns=8.8.8.8 \
    --rm \
    -p 5001:5001 -p 8080:8080 \
    iperf-node
    
    # -it --entrypoint sh \