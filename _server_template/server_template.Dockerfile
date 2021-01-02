ARG BUILDER_IMAGE=golang:buster

# :debug has busybox shell
ARG RUN_IMAGE=gcr.io/distroless/base:debug

# to make WORK_DIR accessible under FROM
# calling `ARG WORK_DIR` is required
# see https://docs.docker.com/engine/reference/builder/#understand-how-arg-and-from-interact
ARG WORK_DIR=/www/server_template

FROM ${BUILDER_IMAGE} as builder
# making WORK_DIR usable inside FROM
ARG WORK_DIR

# up to date certs
RUN update-ca-certificates

# setup go environment
RUN go env -w \ 
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    GO111MODULE=on

WORKDIR ${WORK_DIR}

# prep dependencies
COPY go.mod .
RUN go mod download; go mod verify

# copy data into workdir
COPY . .

# build the app
RUN go build \
# use native 
    # force pure Go resolver on net and os/user
    # -tags='netgo osusergo' \
    # disable symbol table and DWARF debug info generation
    -ldflags='-s -w -extldflags "-static"' \ 
    # -ldflags='-s -w' \ 
    # force rebuild for all packages
    -a \
    # custom install suffix
    -installsuffix nocgo  \
    # executable name
    -o ./server_template \
    # path to code
    ./cmd/server/main.go

# helpers
RUN echo "\nChecking if binary is static:"; \
    readelf -d server_template; \
    ldd -v ./server_template || true; \
    echo "\nChecking files:"; \
    ls -Rla; \
    echo "\nChecking system and go setup:"; \
    uname -a; \
    go version; \
    go env

FROM ${RUN_IMAGE}
# making WORK_DIR usable inside FROM
ARG WORK_DIR

# hotfix for busybox issues..
COPY --from=amd64/busybox:1.32.0 /bin/busybox /busybox/busybox
# setup busybox shell
SHELL ["/busybox/sh", "-c"]

# copy the executable
COPY --from=builder ${WORK_DIR}/server_template ${WORK_DIR}/server_template
COPY --from=builder ${WORK_DIR}/conf/ ${WORK_DIR}/conf/
COPY --from=builder ${WORK_DIR}/static ${WORK_DIR}/static

# run stuff under nobody:nobody
RUN chown -R 65534:65534 ${WORK_DIR}

WORKDIR ${WORK_DIR}

# helper to check files
RUN ["/busybox/sh", "-c", "ls -Rla ./"]

# perform any further action as an unprivileged nobody:nobody user
USER 65534:65534

# lets run the app
# ENTRYPOINT ["./server_template", "-c", "conf/server_template.yaml"]
ENTRYPOINT ["./server_template", "-e", "server_template"]

# app should serve on 8080
EXPOSE 8080