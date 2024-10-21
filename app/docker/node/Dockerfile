# Use the official Go image to build the application
FROM golang:latest

# Set the current working directory inside the container
WORKDIR /app

# Copy the Go module files
COPY ../../go.mod ../../go.sum ./

# Download the Go module dependencies
RUN go mod download

# Copy the source code into the container
COPY cmd/node/ ./cmd/node/
COPY pkg/ ./pkg/

# Build the Go application
RUN (CGO_ENABLED=0 GOOS=linux cd ./cmd/node && go build -o ./node)

# Expose the port your application runs on
EXPOSE 8080

# Command to run the application
CMD ["./cmd/node/node"]
