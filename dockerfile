# syntax=docker/dockerfile:1
FROM golang:1.16-alpine3.15

LABEL maintainer="Doh Halle <u_mhallen@hotmail.com>"

RUN mkdir /bank-operations

WORKDIR  /bank-operations

COPY go.mod .

COPY go.sum .

RUN go mod download

COPY . .

RUN go build -o /bo-api

#RUN ./bank-operations

EXPOSE 7000

CMD [ "/bo-api" ]
