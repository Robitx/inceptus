# inceptus
The inceptus is (for now) a perosnal framework for making servers in go..

Clone this repo and follow instructions in bootstrap.sh =>
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
    - REPOSITORY will be used instead of:
      "github.com/robitx/inceptus/server_template"
    - DIRECTORY is a folder where you want to boostrap the server:
      ├── cmd
      │   └── server
      │       └── main.go
      ├── conf
      │   ├── PROJECT_NAME.env
      │   └── PROJECT_NAME.yaml
      └── internal
          ├── do
          │   └── do.go
          └── env
              ├── config.go
              └── environment.go
      ...

```