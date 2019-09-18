FROM golang:1-alpine

RUN apk update && \
    apk upgrade && \
    apk add git make

ENV CGO_ENABLED=0
ENV GO111MODULE=off

WORKDIR /go/src/github.com/blend/go-sdk

ADD . .

RUN go get ./...

CMD [ "make", "ci" ]
