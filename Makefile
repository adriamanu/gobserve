test:
	go test -v ./...
build:
	go build
docker-build:
	docker build . -t goverwatch
ci:
	go test -v --covermode=count --coverprofile=./coverage.out ./...