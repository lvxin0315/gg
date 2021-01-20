FROM golang:latest as builder

EXPOSE 8080

ENV GOPROXY https://goproxy.io
ENV GO111MODULE on


RUN echo 'Asia/Shanghai' >/etc/timezone

WORKDIR $GOPATH/src/github.com/lvxin0315/gg

ADD . .

RUN go mod download

RUN go build -o /tmp/gg-server main.go

FROM ubuntu:20.04

WORKDIR /
COPY --from=builder  /tmp/gg-server /gg-server

ADD config.toml /config.toml


CMD ./gg-server