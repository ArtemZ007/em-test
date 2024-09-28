# Используем официальный образ Golang
FROM golang:1.23-alpine

# Устанавливаем рабочую директорию внутри контейнера
WORKDIR /app

# Копируем go.mod и go.sum и устанавливаем зависимости
COPY go.mod go.sum ./
RUN go mod download

# Копируем весь исходный код в контейнер
COPY . .

# Собираем приложение
RUN go build -o main .

# Открываем порт
EXPOSE 8080

# Команда для запуска приложения
CMD ["./main"]
