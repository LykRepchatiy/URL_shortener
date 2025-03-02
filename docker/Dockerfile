# Используем официальный образ Go
FROM golang:1.24-alpine

# Создаем рабочую директорию
WORKDIR /app

# Копируем файлы модулей
COPY ../go.mod .
COPY ../go.sum .

# Скачиваем зависимости
RUN go mod download

# Копируем исходный код
COPY . .

# Собираем приложение
RUN go build -o main ./cmd/url-shortener/

# Команда для запуска приложения
CMD ["./main"]