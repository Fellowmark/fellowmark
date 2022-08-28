# FellowMark
Education application for student to student cross reviewing of assignments.

# Requirements
1. docker
2. docker-compose
3. jq: https://stedolan.github.io/jq/

## Optional
1. go
2. go-tools (search gopls support for your editor)

# Setup
`.env` file contains the configuration variables.

## Local dev API setup
```shell
docker-compose --profile api-dev up
```

## Local dev front end setup
```shell
docker-compose --profile frontend-dev up
```

To start the fellowmark system, run
```shell
docker-compose --profile dev up
```

After api is running, populate with mock data by running the seed.sh script:
```
bash seed.sh
```

The seed.sh creates the following:
* 1 Staff 20 Student accounts
* 1 Module
* 1 Assignment and its marker-reviewee pairings
* 1 Question
* 2 Rubrics

Ping server for health
```sh
curl localhost:5000/health
```

Frontend can be accessed on `http://localhost:3000`
Backend can be accessed on `http://localhost:5000`
