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
docker:
	docker build -t hello-back .
	docker run --env-file ./.env -p 6969:6969 hello-back
#will forcefully remove all unused Docker images
clean:
	docker image prune -f