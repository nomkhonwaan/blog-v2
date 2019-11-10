FROM golang:1.13-alpine AS builder
ENV CGO_ENABLED=0
WORKDIR /go/src/github.com/nomkhonwaan/myblog
RUN apk add -u --no-cache ca-certificates curl git make tzdata
COPY go.mod go.sum Makefile ./
RUN make install
COPY . .
RUN make build
FROM scratch
COPY --from=builder /etc/ssl/certs/* /etc/ssl/certs/
COPY --from=builder /go/src/github.com/nomkhonwaan/myblog/bin/myblog .
COPY --from=builder /usr/share/zoneinfo/* /usr/share/zoneinfo/
ENTRYPOINT ["/myblog"]
