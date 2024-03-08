#!/bin/bash

curl localhost:8080/sum \
  -v \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -X POST -k \
  --data '{"a":{"b":4.2, "ww":[1.2,-3.4]},"c":-2.1,"z":0.5}'
