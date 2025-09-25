FROM golang:1.21-alpine AS builder

# Устанавливаем необходимые пакеты для CGO и SQLite
RUN apk add --no-cache gcc musl-dev sqlite-dev

WORKDIR /app
COPY . .
RUN go mod download
# Включаем CGO для SQLite
RUN CGO_ENABLED=1 go build -o rva_bot main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates sqlite
WORKDIR /root/

COPY --from=builder /app/rva_bot .
COPY --from=builder /app/env.production.example .env

# Создаем директорию для SQLite базы данных
RUN mkdir -p /data

CMD ["./rva_bot"]
