ARG PKG=github.com/nomkhonwaan/myblog
FROM golang:1.12-alpine AS builder
ARG PKG
ENV CGO_ENABLED=0
WORKDIR /go/src/${PKG}
COPY Gopkg.lock Gopkg.toml Makefile ./
RUN apk add -u --no-cache curl git make && \
    curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh && \
    dep ensure --vendor-only
COPY . .
RUN make build
FROM scratch
ARG PKG
COPY --from=builder /go/src/${PKG}/bin/myblog .
ENTRYPOINT ["/myblog"]