FROM golang:1.13 as build

WORKDIR /go/src/app
COPY . .

RUN GO111MODULE=on go get
RUN GO111MODULE=on GOOS=linux GOARCH=amd64 go build -ldflags="-w -s"
RUN chmod +x /go/src/app/redis-queue-bridge

FROM debian:buster-slim

COPY --from=build /go/src/app/redis-queue-bridge /usr/local/bin/redis-queue-bridge
RUN mkdir -p /etc/bridge/config
ENTRYPOINT [ "redis-queue-bridge" ]
