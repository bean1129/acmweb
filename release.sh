#!/bin/sh

run_container_id=`docker ps -q -f "name=acmsvr"`
docker cp ${run_container_id}:/acmsvr/run/acmsvr ./bin/