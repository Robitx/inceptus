# server_template

## Prerequisities
Required for development:
 - [Go](https://golang.org/doc/install#install) 
 - [docker](https://docs.docker.com/get-docker/)
 - [docker-compose](https://docs.docker.com/compose/install/)

 Highly recommended:
 - [mkcert](https://github.com/FiloSottile/mkcert) for making locally-trusted development certificates,
  see section [HTTPS for local development](#https-for-local-development) 

-------------------------------------

## Dev stack
The current dev stack consists of three containers:
1) server_template_proxy based on [traefik](https://github.com/traefik/traefik/)
2) [server_template_database](#server_template_database) based on [postgres](https://www.postgresql.org/)
3) your [server_template](#server_template-itself) server bootstraped from [inceptus](https://github.com/Robitx/inceptus)

If you're new to docker-compose, you might want to check the [command reference](https://docs.docker.com/compose/reference/),
bellow are some commands for basic operations.

#### **Simple start of the stack**
attached to the terminal, prints logs into the terminl from all containers (stoppable by `Ctrl+C`, killable by `Ctrl+C Ctrl+C`)

    docker-compose -f docker-compose-dev.yaml up

#### **Detached start with rebuilding of changed images** 
(to apply new changes, don't stop the stack, just run this commnad again)

    docker-compose -f docker-compose-dev.yaml up --detach --build

To see logs for all containers:

    docker-compose -f docker-compose-dev.yaml logs -f

Logs just for server_template:

    docker-compose -f docker-compose-dev.yaml logs -f server_template
  
#### **Stopping/Kiling the stack**
    
    docker-compose -f docker-compose-dev.yaml stop
    docker-compose -f docker-compose-dev.yaml kill

  
#### **Full cleanup after the stack**

    docker-compose -f docker-compose-dev.yaml down -v --rmi all

-------------------------------------

## server_template_proxy
You can check the [dashboard](https://traefik.server_template.localhost) of the Traefik proxy.


## server_template_database
If you want to interact with the database directly, run:
```
docker exec -ti server_template_database bash
psql --host=localhost --username=_server_template_db_user --dbname=_server_template_db
```

Postgre image uses server_template_database volume to persist data => the sql scripts in _database/init/*.sql are applied only once => if you change them, shut down and remove the volume, to start cleanly next time:
```
docker-compose -f docker-compose-dev.yaml down -v
```

## server_template itself
The example server is still work in progress, but its already usable. Here are some already implemented features:
- multistage building that results in a docker image based on [distroless](https://github.com/GoogleContainerTools/distroless) that has only around 50MB in size (and could be sized down even more by eliminating busybox shell). You can attach to the container's shell by running:

      docker exec -ti server_template sh

- the server inside the docker image is run as nobody:nobody (just to be on safer side)
- no global variables, instead using dependency injections (might incorporate [wire](https://github.com/google/wire) in the future) 
- graceful shutdown (basic signal handlers (SIGTERM, SIGINT, ..) are implemented with inceptus/life package)  
  (If you're naughty and your app might produce zombies, you might want to integrate [tini](https://github.com/krallin/tini) or [dumb-init](https://github.com/Yelp/dumb-init).)
- config is specified either as yaml or preferably ENVs (with helper inceptus/conf package based on [viper](https://github.com/spf13/viper))
- json logging (with inceptus/log helper package based on [zerolog](https://github.com/rs/zerolog))
- static files only server (without exposing dir structure)
- work in progress router package (with some basic middleware based on [chi-go](https://github.com/go-chi/chi))
- ...


To attach to the container with interactive shell:
```
docker exec -ti server_template sh
```

The example app used for bootstraping

Trying the api:
```
curl -D - -v -X GET https://server_template.localhost/api/v1/echo -H "x-request-id: requestID_XYZ" -H "accessToken: JWT_HERE" -d 'hey!'
```

For static site hello world go to http://localhost:8080/static/

-------------------------------------

## HTTPS for local development
To get https for localhost:
- install [mkcert](https://github.com/FiloSottile/mkcert) and run the following
- create root [CA](https://en.wikipedia.org/wiki/Certificate_authority) and add it to your system and browsers (might need restart):

      mkcert -install
- generated CA can be found under:

      mkcert -CAROOT

- create certs used by the server_template:

      mkcert -cert-file _proxy/certs/local-cert.pem -key-file _proxy/certs/local-key.pem "*.server_template.localhost" "server_template.localhost"

And thats it, enjoy https for localhost.

*Technically mkcert setup is optional. If you don't do this step before running the docker-compose stack, traefik proxy will create default certs for you. But they won't be trusted by your browser and you have to deal with enoying insecure warnings on your own..*

-------------------------------------