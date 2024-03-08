#!/bin/bash
curl localhost:8080/auth \
  -v \
  -H "Content-Type: application/json" \
  -X POST -k \
  --data '{"username":"some-username","password":"some-password"}'
