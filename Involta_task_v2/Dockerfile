# Вариант создания контейнера с компиляцией внутри контейнера

FROM golang:1.20
WORKDIR /app
COPY ./ ./
RUN go mod download
WORKDIR /app/cmd
RUN go build -o gocrud .

EXPOSE 8080

CMD [ "./gocrud" ]





