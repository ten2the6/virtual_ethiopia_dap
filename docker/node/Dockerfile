# docker/node/Dockerfile
FROM golang:1.23-alpine

WORKDIR /app

# Copy go mod files first for better layer caching
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download && go mod verify

# Copy the rest of the application
COPY . .

# Build the application
RUN go build -o /bin/node cmd/node/main.go

# Expose the necessary ports
EXPOSE 3000
EXPOSE 30301

CMD ["/bin/node"]