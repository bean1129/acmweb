#!/bin/sh

run_container_id=`docker ps -q -f "name=acmsvr"`
echo "Running Container Id: ${run_container_id}"
if [ "${run_container_id}" != "" ]; then
    docker kill ${run_container_id}
fi
container_id=`docker ps -qa -f "name=acmsvr"`
echo "Stopped Container Id: ${container_id}"
if [ "${container_id}" != "" ]; then
    docker rm ${container_id}
fi
images_id=`docker images -qa "acmsvr"`
echo "Images Id: ${images_id}"
if [ "${images_id}" != "" ]; then
    docker rmi ${images_id}
fi

docker build -f ./DockerFile -t acmsvr:v1 .
docker run -d -p 9088:9088 -p 9089:9758 --name acmsvr acmsvr:v1
run_container_id=`docker ps -q -f "name=acmsvr"`
docker cp ${run_container_id}:/acmsvr/run/acmsvr ./bin/