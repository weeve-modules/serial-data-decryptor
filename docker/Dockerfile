FROM golang:1.18.4-alpine3.16 AS builder
WORKDIR /app
COPY ./src .
RUN go mod download
RUN unset GOPATH && go build -o app ./cmd/main.go

###
FROM alpine:edge
# Oh yes, you need bash!
RUN apk add --no-cache bash
SHELL ["/bin/bash", "-c"]
COPY --from=builder /app/app ./
COPY ./docker/entrypoint.sh /entrypoint.sh

ENTRYPOINT ["/entrypoint.sh"]
