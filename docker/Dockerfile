FROM golang:1.18.4-alpine3.16 AS builder

# Set destination for COPY
WORKDIR /app

# Download Go modules
COPY go.mod .
COPY go.sum .
RUN go mod download

# Build
RUN unset GOPATH && go build -o app ./cmd/main.go
RUN echo Built the Golang binary to ./app

###
FROM alpine:edge
RUN apk add --no-cache bash
SHELL ["/bin/bash", "-c"]
COPY --from=builder /app/app ./
COPY ./docker/entrypoint.sh /entrypoint.sh

RUN ["chmod", "+x", "/entrypoint.sh"]

ENTRYPOINT ["/entrypoint.sh"]
###