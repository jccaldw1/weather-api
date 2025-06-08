# Start with an official Go image
FROM golang:1.23-alpine

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum and download dependencies
COPY go.mod ./
RUN go mod download

# Copy the source code
COPY . .

# Build the Go binary
RUN go build -o myapp .

# Expose port 8080
EXPOSE 8080

# Run the binary
CMD ["./myapp"]
