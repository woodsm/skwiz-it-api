FROM golang:1.9.2 AS build
ADD . /src
WORKDIR /src
RUN go get -d -v -t
RUN go test --cover -v ./... --run UnitTest
RUN go build -v -o skwiz-it-api



FROM alpine:3.7
MAINTAINER 	Ben Kauffman <ben@krashidbuilt.com>

RUN apk apk update && apk upgrade
RUN apk add --no-cache ca-certificates

RUN mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2

WORKDIR /usr/local/bin
COPY --from=build /src/skwiz-it-api /usr/local/bin/skwiz-it-api
COPY --from=build /src/config.json /usr/local/bin/config.json
RUN chmod +x /usr/local/bin/skwiz-it-api

CMD ["skwiz-it-api"]

EXPOSE 3000