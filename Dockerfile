FROM 242617/go-builder:1.0.2 AS builder

WORKDIR /build
ADD . .
RUN make build

FROM alpine:3.12
RUN apk --no-cache add ca-certificates
COPY --from=builder /build/bin/app /opt/app

STOPSIGNAL 15
HEALTHCHECK --interval=5m --timeout=5s CMD curl -f http://127.0.0.1:8080/healthz || exit 1

EXPOSE 8080/tcp

ENTRYPOINT ["/opt/app"]
