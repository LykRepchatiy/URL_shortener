FROM golang:1.24.0

WORKDIR /app

COPY . .

RUN go build -o url-shortener .

EXPOSE 8080

CMD ["./url-shortener"]