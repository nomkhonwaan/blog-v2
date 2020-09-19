ARG NPM_AUTH_TOKEN
FROM golang:1.15-alpine AS builder
ENV CGO_ENABLED=0
WORKDIR /go/src/github.com/nomkhonwaan/myblog
RUN apk add --upgrade --no-cache ca-certificates curl git make tzdata
COPY go.mod go.sum Makefile ./
RUN make install
COPY . .
RUN make build
FROM node:12.18-alpine AS web-builder
ARG NPM_AUTH_TOKEN
WORKDIR /opt/app-root
RUN apk add --upgrade --no-cache curl gettext git make
COPY Makefile ./
COPY ./web/.npmrc ./web/package.json ./web/package-lock.json ./web/
RUN make install-web && \
    rm -f ./web/.npmrc
COPY . .
RUN make build-web
FROM scratch
COPY --from=builder /etc/ssl/certs/* /etc/ssl/certs/
COPY --from=builder /usr/share/zoneinfo/* /usr/share/zoneinfo/
COPY --from=builder /go/src/github.com/nomkhonwaan/myblog/bin/myblog-server .
COPY --from=web-builder /opt/app-root/web/dist/ ./dist/
ENTRYPOINT ["/myblog-server"]
