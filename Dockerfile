# Start from the alpine3.17 golang base image
FROM golang:alpine3.17 as builder

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
COPY . /app

# Build the Go app
RUN go build -o build/hello-back cmd/main.go

# Final stage
FROM golang:alpine3.17

# Copy the binary from the builder stage
COPY --from=builder /app/build/hello-back /

# Expose port 8080 to the outside world
EXPOSE 8080

# Run the binary program produced by `go build`
CMD ["/hello-back"]