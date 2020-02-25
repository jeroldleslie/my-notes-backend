#!/bin/bash


docker-machine start default
docker-machine env default
eval "$(docker-machine env default)"
host_ip=$(docker-machine ip default)
echo $host_ip
export POSTGRESQL_ADDRESS=postgres://postgres:postgres@$host_ip/notes?sslmode=disable

echo $POSTGRESQL_ADDRESS
docker-compose -f ./docker-compose_notes_api.yml up --build -d postgres_notes


while ! nc -z $host_ip 5432;
do
    echo Wait for database to start...;
    sleep 1;
done;
echo Connected to database;

docker-compose -f ./docker-compose_notes_api.yml up --build -d notes_api