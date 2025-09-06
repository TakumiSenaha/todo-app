# Builder stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy dependency files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy all source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o /main .

# Final stage
FROM alpine:latest AS final

# Create app directory
RUN mkdir /app

# Copy binary from builder stage
COPY --from=builder /main /app/main

# Expose port
EXPOSE 8080

# Run the application
CMD ["./app/main"]
