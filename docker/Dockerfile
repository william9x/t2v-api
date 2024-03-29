FROM golang:1.21-alpine AS builder
ARG BUILD_MODULE
ARG BUILD_VERSION
ARG BUILD_COMMIT_HASH
ARG BUILD_TIME
ARG BS_PKG=github.com/Braly-Ltd/t2v-api-$BUILD_MODULE/bootstrap
ENV GO111MODULE=on
RUN apk add --no-cache git
RUN apk add --no-cache ca-certificates && update-ca-certificates
COPY ./src /go/src
WORKDIR /go/src/$BUILD_MODULE
# export APP_PROFILES=dev GIN_MODE=release && ./binary
RUN env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o ./binary \
    -ldflags="-X '$BS_PKG.Version=$BUILD_VERSION' -X '$BS_PKG.CommitHash=$BUILD_COMMIT_HASH' -X '$BS_PKG.BuildTime=$BUILD_TIME'"

FROM scratch
ARG BUILD_MODULE
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /go/src/$BUILD_MODULE/config /app/config
COPY --from=builder /go/src/$BUILD_MODULE/binary /app
EXPOSE 8080
WORKDIR /app
ENTRYPOINT ["./binary"]
