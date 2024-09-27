# Используем официальный образ Go в качестве базового
FROM golang:1.20-alpine

# Установим необходимые зависимости
RUN apk add --no-cache git

# Устанавливаем рабочую директорию внутри контейнера
WORKDIR /app

# Копируем файлы проекта
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Собираем приложение
RUN go build -o main ./cmd/api/main.go

# Определяем порт, который будет прослушивать приложение
EXPOSE 8080

# Команда для запуска приложения
CMD ["./main"]
