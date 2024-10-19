FROM golang:latest AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod tidy

COPY . .
RUN go build -o output .

CMD ["./output"]