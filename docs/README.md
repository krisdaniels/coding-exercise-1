The service requires 2 env vars to be set: 
- SECRET: which contains the secret for the jwt token signing 
- ISSUER: the issuer which will be set on the iss claim of the jwt token and is also verified on authorized endpoints, /sum in this case

The scripts below will set these variables to a demo value automatically.

Running the service: 
- use start.sh in the root 

Building the docker container:
- in directory packages, run build-docker.sh

Running the docker container:
- in directory packages, run start-docker.sh

Provided examples in examples directory:
- auth.sh will call the auth endpoint and print out the token
- sum-external-token.sh, expects a token to be present in the TOKEN env var
- sum-fetch-token.sh will call the auth endpoint to fetch a token, but requires jq to be installed to extract the token from the response

Additional build script to run tests with coverage and print out failing tests:
- build/test.sh