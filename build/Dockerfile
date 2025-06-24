# syntax=docker/dockerfile:1
FROM golang:1.24 AS builder

ARG VERSION
ARG BUILD
ARG BUILD_DATE
RUN apt install ca-certificates
ADD . /go-server
WORKDIR /go-server/server
RUN sed -i -e "s|@{build}|${BUILD}|g" -e "s|@{version}|${VERSION}|g" -e "s|@{buildDate}|${BUILD_DATE}|g" internal/main.go
RUN go build -o ../main internal/main.go

FROM ubuntu:latest
WORKDIR /go-server
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /go-server/main .
EXPOSE 8080
USER root:root
CMD ["./main"]
