FROM golang:latest

LABEL maintainer="Doh Halle <u_mhallen@hotmail.com>"

RUN mkdir /bank-operations

WORKDIR  /bank-operations

COPY go.mod .

COPY go.sum .

RUN go mod download

COPY . .

EXPOSE 8080

RUN go build && ./bank-operations
