FROM public.ecr.aws/docker/library/golang:latest AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy only the Go module files first to cache dependencies
COPY go.mod go.sum ./

# Download and cache Go dependencies
RUN go mod download

# Copy the rest of the project files
COPY . .

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