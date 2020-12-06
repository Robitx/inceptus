# inceptus
[![Docs](https://img.shields.io/badge/docs-current-brightgreen.svg)](https://pkg.go.dev/github.com/robitx/inceptus)
[![Go Report Card](https://goreportcard.com/badge/github.com/robitx/inceptus)](https://goreportcard.com/report/github.com/robitx/inceptus)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/Robitx/inceptus/blob/main/LICENSE)

The inceptus is (for now) a personal framework for making (web) servers in go..

The goal of inceptus is to have all you need to start making a new full stack web app (PWA with golang backend).

Just call bootstrap.sh and have all the boring boilerplate ready, so you can invest time in what actually matters..

-------------------------------------
## Bootstraping

Clone this repo and follow the instructions in bootstrap.sh:
```
$./bootstrap.sh -h
Hey,
this script expects you already have existing git repository, into
which you want to bootstrap server based on inceptus server_template.

Usage:
    bootstrap.sh -h Display this help message.

    bootstrap.sh -n PROJECT_NAME -r REPOSITORY -d DIRECTORY
    where:
    - PROJECT_NAME will be used instead of "server_template"
    - REPOSITORY (like github.com/XXX/YYY) will be used instead of:
      "github.com/robitx/inceptus/server_template"
    - DIRECTORY is a folder where you want to boostrap the server:
      ├── cmd
      │   └── server
      │       └── main.go
      ├── conf
      │   ├── server_template.env
      │   └── server_template.yaml
      ├── internal
      │   ├── do
      │   │   └── do.go
      │   ├── env
      │   │   ├── config.go
      │   │   └── environment.go
      │   └── rest
      │       └── rest.go
      ├── main
      └── static
          ├── errors
          │   └── 404.html
          └── index.html
      ...
```

-------------------------------------

##  Local

Starting local instance of the server:
```
go build ./cmd/server/main.go; ./main -c conf/server_template.yaml
```

Trying the api:
```
url -D - -v -X GET localhost:9999/api/v1/echo -H "x-request-ids: requestID_XYZ" -H "accessToken: JWT_HERE" -d 'hey!'
```

For static site hello world go to http://localhost:9999/static/


-------------------------------------

##  Docker
Build:
```
docker build --tag server_template -f ./server_template.Dockerfile ./
```

Run:
```
docker run  --rm --name server_template -p 127.0.0.1:8080:8080 server_template:latest
```

Interactive shell:
```
docker exec -ti server_template sh
```
