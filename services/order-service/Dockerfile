FROM golang:1.24-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

EXPOSE 8080

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o order-service .

CMD ["./order-service"]
