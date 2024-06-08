FROM golang:1.22.0-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o main server.go
RUN go mod tidy

CMD ["/app/main"]