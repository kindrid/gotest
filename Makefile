init:
	go get -t -v ./...

test:
	go test ./ ./should

test-watch:
	modd

build-ci: init
	go build -v ./...

ci-job1: test
