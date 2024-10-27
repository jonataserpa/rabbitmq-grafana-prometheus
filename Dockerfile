# Use Go base image
FROM golang:1.19-alpine

# Set up working directory
WORKDIR /app

# Copy Go module files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the Go application
RUN go build -o main .

# Run the application
CMD ["./main"]
