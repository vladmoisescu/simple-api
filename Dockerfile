FROM golang:1.16-alpine

WORKDIR /
COPY bin/simple-api /simple-api

EXPOSE 8080

ENTRYPOINT ["/simple-api"]