FROM golang:1.18.4-alpine3.16 AS builder
WORKDIR /app
COPY ./src .
RUN go mod download
RUN unset GOPATH && go build -o app ./cmd/main.go

###
FROM alpine:edge
COPY --from=builder /app/app ./
COPY ./docker/docker-entrypoint.sh /entrypoint.sh

ENTRYPOINT ["/entrypoint.sh"]
