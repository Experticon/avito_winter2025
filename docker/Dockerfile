# Используем базовый образ Go для сборки
FROM golang:1.23 AS builder

WORKDIR /app

# Копируем зависимости и устанавливаем их
COPY go.mod go.sum ./
RUN go mod download

# Копируем весь код в контейнер
COPY . ./

# Собираем сервер
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/main.go

# Финальный образ
FROM alpine:latest

WORKDIR /app

# Копируем собранный сервер
COPY --from=builder /app/server .

# Копируем Swagger UI и swagger.yaml в контейнер
COPY --from=builder /app/pkg/swagger-ui/dist /app/swagger-ui
COPY --from=builder /app/internal/go-server-server-generated/api/swagger.yaml /app/swagger.yaml

# Открываем порт
EXPOSE 8080

# Запускаем сервер
ENTRYPOINT ["./server"]
