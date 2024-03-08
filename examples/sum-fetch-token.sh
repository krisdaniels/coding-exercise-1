#!/bin/bash

# sudo apt install jq
TOKEN=$(curl localhost:8080/auth -H "Content-Type: application/json" -X POST -k --data '{"username":"some-username","password":"some-password"}' | jq -r '(.token)')

curl localhost:8080/sum \
  -v \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -X POST -k \
  --data '{"a":{"b":4.2, "ww":[1.2,-3.4]},"c":-2.1,"z":0.5}'
