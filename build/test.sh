#!/bin/bash

rm -f coverage.out
rm -f coverage.html
rm -f testresults.json

go test ../... -coverprofile=coverage.out -json > testresults.json
go tool cover -html=coverage.out -o ./coverage.html

cat testresults.json | grep -i '"action":"fail"'

# python3 -m http.server 8090
