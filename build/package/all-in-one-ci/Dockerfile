FROM alpine:3.10 AS builder
RUN apk add --upgrade --no-cache ca-certificates tzdata
FROM scratch
COPY --from=builder /etc/ssl/certs/* /etc/ssl/certs/
COPY --from=builder /usr/share/zoneinfo/* /usr/share/zoneinfo/
COPY ./bin/myblog-server .
COPY ./web/dist/ ./dist/
ENTRYPOINT ["/myblog-server"]
