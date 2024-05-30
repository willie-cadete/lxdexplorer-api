# Stage 1: Build the Go application
FROM golang:1.22 AS builder

ARG VERSION=dev

WORKDIR /app

# Copy the Go module files
COPY go.mod go.sum ./

# Download the Go module dependencies
RUN go mod download

# Copy the application source code
COPY . .

# Build the Go application
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags "-X main.version=$VERSION"

# Stage 2: Create a minimal runtime image
FROM debian:bookworm-slim

RUN apt-get update && apt-get install -y ca-certificates

WORKDIR /app/

# Copy the built Go application from the previous stage
COPY --from=builder /app/lxdexplorer-api .
COPY config.yaml.example /app/config.yaml

ENV GIN_MODE=release

# Expose the application port
EXPOSE 8080

# Run the Go application
CMD ["./lxdexplorer-api"]
