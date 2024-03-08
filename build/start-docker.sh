#!/bin/bash

#SECRET=$(openssl rand -base64 32)
SECRET="some-secret"
ISSUER="some-issuer"
docker run -ti --rm -e SECRET=$SECRET -e ISSUER=$ISSUER -p 8080:8080 coding-exercise-1
