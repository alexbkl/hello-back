# Start from the latest golang base image
FROM golang:latest AS builder

# Add Maintainer Info
LABEL maintainer="Hello Decentralized <hi@hello.ws>"

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. 
# Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app
RUN go build -o bin/hello-back main.go

# Final stage
FROM golang:latest

WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/bin/hello-back /app/hello-back

# Expose port 6969 to the outside world
EXPOSE 6969

# Run the binary program produced by `go build`
CMD ["/app/hello-back"]