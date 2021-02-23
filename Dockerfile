FROM golang:1.15-alpine

WORKDIR /srv/app

COPY . .

RUN apk add -Uv \
  alpine-sdk
