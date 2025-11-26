FROM golang:1.24-bookworm AS builder

WORKDIR /app
COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN apt-get update && apt-get install -y --no-install-recommends dumb-init
RUN CGO_ENABLED=0 GOOS=linux GOAMD64=v3 go build -ldflags="-s -w" -o api ./cmd/api

FROM gcr.io/distroless/base-debian12

WORKDIR /

COPY --from=builder /usr/bin/dumb-init /usr/bin/dumb-init
COPY --from=builder /app/api /api
COPY --from=builder /app/data ./data

EXPOSE 8080

ENV PREFORK=true

USER nonroot:nonroot

ENTRYPOINT ["/usr/bin/dumb-init", "--"]
CMD ["/api"]
