FROM golang:latest AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod tidy

COPY . .
RUN go build -o ./out/go-sample-app .

CMD ["./out/go-sample-app"]