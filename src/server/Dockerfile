# Use the official Golang image for building
FROM golang:1.24-alpine AS builder

LABEL authors="richardsimonik"
# Set working directory
WORKDIR /app
# Install git (needed for go mod download)
RUN apk add --no-cache git
# Copy Go modules and dependencies
COPY go.mod go.sum ./
RUN go mod download && go mod verify
# Copy source code
COPY . .
# Build the application
RUN go build -o main .
# Use a minimal base image for final deployment
FROM alpine:latest
# Set working directory in the container
WORKDIR /root/
# Copy the built binary from the builder stage
COPY --from=builder /app/main .
# Expose the application port
EXPOSE 8080
# Run the application
CMD ["./main"]
