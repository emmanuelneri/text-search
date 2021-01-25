#!/usr/bin/env bash
# jq required

USERNAME=api@gmail.com.br
PASSWORD=api

response=$(curl -s -X POST \
-d "username=${USERNAME}" \
-d "password=${PASSWORD}" \
http://localhost:8000/auth)

echo ${response}
#
ACCESS_TOKEN=$(echo ${response} | jq -r ".access_token")
#curl http://localhost:8000/search/v1/users
curl -H "Authorization: Bearer ${ACCESS_TOKEN}" http://localhost:8000/search/v1/users
