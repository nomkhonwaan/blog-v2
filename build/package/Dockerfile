FROM golang:1.12-alpine AS builder
ENV CGO_ENABLED=0
ENV GO111MODULE=on
WORKDIR /opt/app-root
COPY go.mod go.sum Makefile ./
RUN apk add -u --no-cache ca-certificates curl git make && \
    make install
COPY . .
RUN make build
FROM scratch
COPY --from=builder /etc/ssl/certs/* /etc/ssl/certs/
COPY --from=builder /opt/app-root/bin/myblog .
ENTRYPOINT ["/myblog"]