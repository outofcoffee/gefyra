FROM golang:1.17 as build

ARG GOARCH=amd64

WORKDIR /go/src/app
COPY . .

RUN GO111MODULE=on go get
RUN GO111MODULE=on GOOS=linux GOARCH=${GOARCH} go build -ldflags="-w -s"
RUN chmod +x /go/src/app/gefyra

FROM debian:buster-slim

COPY --from=build /go/src/app/gefyra /usr/local/bin/gefyra
RUN mkdir -p /etc/gefyra/config
ENTRYPOINT [ "gefyra" ]
