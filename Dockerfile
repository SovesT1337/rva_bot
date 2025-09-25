FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o rva_bot main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY --from=builder /app/rva_bot .
COPY --from=builder /app/env.production.example .env

# Создаем директорию для SQLite базы данных
RUN mkdir -p /data

CMD ["./rva_bot"]
