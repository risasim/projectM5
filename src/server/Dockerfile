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
# Copy migration files to the container
COPY --from=builder /app/db/migrations ./db/migrations
# Copy soundEffects directory if it exists
COPY --from=builder /app/soundEffects ./soundEffects
# Create soundEffects directory if it doesn't exist (for runtime file uploads)
RUN mkdir -p ./soundEffects
# Copy already built frontend
COPY --from=builder /app/web ./web

# Expose the application port
EXPOSE 8080
# Run the application
CMD ["./main"]
