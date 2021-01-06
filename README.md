# inceptus
[![Docs](https://img.shields.io/badge/docs-current-brightgreen.svg)](https://pkg.go.dev/github.com/robitx/inceptus)
[![Go Report Card](https://goreportcard.com/badge/github.com/robitx/inceptus)](https://goreportcard.com/report/github.com/robitx/inceptus)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/Robitx/inceptus/blob/main/LICENSE)

The inceptus is (for now a personal) framework for making (web) servers in go..

*The goal of inceptus is to have all the basic boilerplate that go servers need so you can focus on the application/business/domain logic instead.*

Just call bootstrap.sh and have all the boring boilerplate ready, so you can invest time into what actually matters..

-------------------------------------
## Bootstraping

Clone this repo, follow the simple instructions in [bootstrap.sh](https://github.com/Robitx/inceptus/blob/main/bootstrap.sh) and you're ready to Go!

```
$./bootstrap.sh 
Hello there,
this script allows you to quickly bootstrap new golang server based on inceptus/_server_template.

Example:
./bootstrap.sh -n shiny_new_project -r github.com/XXXX/shiny_new_project -d /tmp/shiny_new_project

cd /tmp/shiny_new_project; docker-compose -f docker-compose-dev.yaml up;

Usage:
    bootstrap.sh -h Display this help message.

    bootstrap.sh -n PROJECT_NAME -r REPOSITORY -d DIRECTORY
    where:
      - PROJECT_NAME will be used instead of "server_template"
      - REPOSITORY (like github.com/XXX/YYY) will be used instead of:
        "github.com/robitx/inceptus/_server_template"
      - DIRECTORY is a folder where you want to boostrap the server

```

-------------------------------------
## _server_template
For more details head over to the [server_template readme](https://github.com/Robitx/inceptus/blob/main/_server_template/README.md).
