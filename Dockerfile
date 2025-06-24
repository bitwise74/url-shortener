FROM golang:1.24.4

WORKDIR /app
COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN GOOS=linux go build -o url-shortener

EXPOSE 8080

CMD ["./url-shortener"]