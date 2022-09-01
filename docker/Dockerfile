FROM --platform=$BUILDPLATFORM golang:1.17-alpine AS build
WORKDIR /src
ARG TARGETOS TARGETARCH
RUN --mount=target=. \
    --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg \
    GOOS=$TARGETOS GOARCH=$TARGETARCH go build -o /out/myapp ./cmd/main.go

FROM alpine
COPY --from=build /out/myapp /bin