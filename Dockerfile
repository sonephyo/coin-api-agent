FROM golang:1.22.5-alpine AS builder

WORKDIR /app

COPY . .

RUN go mod tidy

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -a -o output .

# Stage 2: Create a minimal final image with scratch
FROM scratch

# Copy the CA certificates for HTTPS support
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Set the working directory in the final container
WORKDIR /root/

# Copy the statically built Go binary from the builder stage
COPY --from=builder /app/output ./

# Define the command to run the application
CMD ["./output"]