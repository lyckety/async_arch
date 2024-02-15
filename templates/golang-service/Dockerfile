FROM golang:1.21.4-alpine3.18 AS builder

ARG BIN_NAME=golang-base-image

WORKDIR /opt/app

COPY . .

FROM alpine:3.18

WORKDIR /srv

COPY migrations migrations

COPY --from=builder /opt/app/bin/golang-base-image /srv/golang-base-image

ENTRYPOINT ["/srv/golang-base-image"]
