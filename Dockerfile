FROM golang:1.13 as build

WORKDIR /go/src/app
COPY . .

RUN GO111MODULE=on go get
RUN GO111MODULE=on GOOS=linux GOARCH=amd64 go build -ldflags="-w -s"
RUN chmod +x /go/src/app/gefyra

FROM debian:buster-slim

COPY --from=build /go/src/app/gefyra /usr/local/bin/gefyra
RUN mkdir -p /etc/gefyra/config
ENTRYPOINT [ "gefyra" ]
