init:
	go get -t -v ./...

test:
	go test ./ ./cuke ./should

test-watch:
	modd

build-ci: init
	go build -v ./...
