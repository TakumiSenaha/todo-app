FROM golang:1.23-alpine

WORKDIR /app

# Install tools
RUN apk add --no-cache git
RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
RUN go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

# Copy dependency files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy all source code
COPY . .

# migrate -database "postgresql://user:password@db:5432/todo_db?sslmode=disable" -path /app/migrations up
# sqlc generate

# Expose port
EXPOSE 8080

# Run with go run for code sync
CMD ["go", "run", "./cmd/api"]
