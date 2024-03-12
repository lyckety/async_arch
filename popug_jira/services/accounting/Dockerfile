FROM golang:1.21.4-alpine3.18 AS builder

WORKDIR /opt/app

COPY . .

RUN apk update && \
    apk add make && \
    rm -rf /var/cache/apk/* && \
    make build

FROM alpine:3.18

WORKDIR /srv

COPY migrations migrations

COPY --from=builder /opt/app/bin/accounting /srv/accounting

ENTRYPOINT ["/srv/accounting"]
