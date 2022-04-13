FROM golang:latest

LABEL maintainer="Doh Halle <u_mhallen@hotmail.com>"

RUN mkdir /app

WORKDIR  /app

COPY go.mod .

COPY go.sum .

RUN go mod download

COPY . .

ENV PORT 5000

RUN go build && ./bank-operations
