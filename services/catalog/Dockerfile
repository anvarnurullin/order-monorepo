# syntax=docker/dockerfile:1
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Копируем go.mod и go.sum и ставим зависимости
COPY go.mod go.sum ./
RUN go mod download

# Копируем весь код
COPY . .

# Сборка бинарника
RUN go build -o catalog ./cmd/catalog

# Финальный образ
FROM alpine:3.18

WORKDIR /app

COPY --from=builder /app/catalog .

# Порт для сервиса
EXPOSE 8082

# Запуск
CMD ["./catalog"]
