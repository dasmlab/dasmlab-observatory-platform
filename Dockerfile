FROM docker.io/library/golang:1.22-bookworm AS builder
WORKDIR /src
COPY go.mod go.sum ./
COPY platform ./platform
RUN go mod download
COPY cmd ./cmd
COPY internal ./internal
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags "-s -w" -o /out/dpo-api ./cmd/dpo-api

FROM docker.io/library/alpine:3.20
RUN apk add --no-cache ca-certificates tzdata \
  && adduser -D -u 65532 -g 65532 dpo
WORKDIR /app
COPY --from=builder /out/dpo-api /app/dpo-api
COPY web /app/web
RUN mkdir -p /data && chown -R 65532:65532 /app /data
USER 65532
ENV DPO_LISTEN=:8080 \
    DPO_DATA_DIR=/data \
    DPO_STATIC_DIR=/app/web \
    DPO_TENANT=dasmlab.org
EXPOSE 8080
VOLUME ["/data"]
ENTRYPOINT ["/app/dpo-api"]
