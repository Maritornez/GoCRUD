# Первый этап: забилдить программу
FROM golang:1.22.5 AS builder
WORKDIR /app
COPY go.mod go.sum ./
# Загрузка зависимостей
RUN go mod download
# Копировать оставшийся код
COPY . .
# Сборка
WORKDIR /app/cmd
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /go/bin/gocrud .


# Второй этап: создание маленького контейнера без компилятора
FROM alpine:latest
# Установка необходимых пакетов
RUN apk --no-cache add ca-certificates
# Копирование бинарного файла
COPY --from=builder /go/bin/gocrud /usr/local/bin/gocrud
# Определение порта, на котором работает приложение
EXPOSE 8080
# Запустить бинарный файл
CMD ["gocrud"]
