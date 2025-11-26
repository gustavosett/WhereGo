FROM golang:1.24-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o api ./cmd/api

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/api .

COPY --from=builder /app/data ./data

EXPOSE 8080

ENV PREFORK=true

CMD ["./api"]
