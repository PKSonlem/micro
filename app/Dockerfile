# Стадия 1: Сборка бинарного файла
FROM golang:1.23 AS builder

# Устанавливаем рабочую директорию внутри контейнера
WORKDIR /app

# Копируем файлы проекта
COPY go.mod go.sum ./
RUN go mod download

# Копируем исходный код
COPY . .

# Компиляция с оптимизацией размера бинарника
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app ./cmd/service/main.go

# Стадия 2: Минимальный контейнер для запуска
FROM alpine:latest

# Устанавливаем рабочую директорию
WORKDIR /root/

# Копируем скомпилированный бинарник из builder-стадии
COPY --from=builder /app/app .

# Делаем бинарник исполняемым
RUN chmod +x ./app

# Запускаем приложение
CMD ["./app"]
