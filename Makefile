install:
	go get .

build:
	go build -o bin/main cmd/main.go

run:
	go run cmd/main.go

format:
	go fmt ./...
	
test:
	go test -v -cover