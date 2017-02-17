init:
	go get -t -v ./...

test:
	go test ./ ./should

test-watch:
	modd

build-ci: init
	go build -v ./...

# Make sure there's no debug code etc.
ci-code-quality:
	@echo "Checking for debugging figments"
	@! grep --exclude Makefile --exclude-dir vendor -nIR 'y0ssar1an/q' *
	@! grep --exclude Makefile --exclude-dir vendor -nIR 'DEBUG' *

# CI First Parallel Job
ci-job1: ci-code-quality test
