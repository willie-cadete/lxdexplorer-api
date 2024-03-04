# Stage 1: Build the Go application
FROM golang:1.22 AS builder

WORKDIR /app

# Copy the Go module files
COPY go.mod go.sum ./

# Download the Go module dependencies
RUN go mod download

# Copy the application source code
COPY . .

# Build the Go application
RUN go build

# Stage 2: Create a minimal runtime image
FROM golang:latest

# RUN apk --no-cache add ca-certificates

WORKDIR /app/

# Copy the built Go application from the previous stage
COPY --from=builder /app/lxdexplorer-api .

ENV GIN_MODE=release

# Expose the application port
EXPOSE 8080

# Run the Go application
CMD ["./lxdexplorer-api"]
