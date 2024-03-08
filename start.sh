#!/bin/bash

# SECRET=$(openssl rand -base64 48)
export SECRET="some-secret"
export ISSUER="some-issuer"

go run main.go
