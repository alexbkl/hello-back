# Start from the alpine3.17 golang base image
FROM golang:alpine as builder

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
RUN CGO_ENABLED=0 GOOS=linux go build -o build/backend cmd/main.go

# Final stage
FROM golang:alpine

# Copy the binary from the builder stage
COPY --from=builder /app/build/backend /app/backend

# Copy the env file
# COPY --from=builder .env /app

# Expose port 8080 to the outside world
EXPOSE 8080

# Run the binary program produced by `go build`
CMD ["/app/backend"]