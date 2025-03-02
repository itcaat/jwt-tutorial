# Используем официальное изображение Go
FROM golang:1.23 as builder

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем файлы проекта
COPY go.mod go.sum ./
RUN go mod download

# Копируем код API
COPY . .

# Сборка бинарного файла
RUN CGO_ENABLED=0 GOOS=linux go build -o api cmd/main.go

# Создаём финальный минимальный образ
FROM alpine:latest

# Устанавливаем рабочую директорию
WORKDIR /app/

# Копируем бинарник из builder-образа
COPY --from=builder /app/api .

# Делаем бинарный файл исполняемым
RUN chmod +x ./api

# Открываем порт 8080
EXPOSE 8080

# Запуск приложения
CMD ["./api"]
