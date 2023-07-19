install:
	go get .

build:
	go build -o bin/main main.go

run:
	go run main.go

format:
	go fmt ./...
	
test:
	go test -v -cover