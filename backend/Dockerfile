# Builder stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Install migrate tool
RUN apk add --no-cache git
RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Copy dependency files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy all source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o /main ./cmd/api

# Final stage
FROM alpine:latest AS final

# Install necessary packages
RUN apk add --no-cache ca-certificates

# Create app directory
RUN mkdir /app

# Copy binary and migrate tool from builder stage
COPY --from=builder /main /app/main
COPY --from=builder /go/bin/migrate /usr/local/bin/migrate

# Copy migrations directory
COPY --from=builder /app/migrations /app/migrations
# migrate -database "postgres://user:password@db:5432/todo_db?sslmode=disable" -path /app/migrations up

# Expose port
EXPOSE 8080

# Run the application
CMD ["./app/main"]
