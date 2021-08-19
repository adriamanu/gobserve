test:
	go test -v .
test-coverage:
	go test -v --covermode=count --coverprofile=coverage.out .
build:
	go build