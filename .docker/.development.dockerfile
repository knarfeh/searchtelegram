FROM golang:1.8.3 as builder
WORKDIR /go/src/github.com/knarfeh/searchtelegram/
COPY . /go/src/github.com/knarfeh/searchtelegram/
