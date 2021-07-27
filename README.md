# NUS Peer Review
Education application for student to student cross reviewing of assignments.

# Requirements
1. docker
2. docker-compose

## Optional
1. go
2. go-tools (search gopls support for your editor)

# Setup
`.env` file contains the configuration variables.

## Local dev API setup
```sh
docker-compose --profile api-dev up
```
## Production API setup
```sh
docker-compose --profile api up
```
