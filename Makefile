VERSION = 0.9.4
COMMIT = ${shell git log --pretty='format:%h' -n 1}
BRANCH = ${shell git rev-parse --abbrev-ref HEAD}


init:
	go get -t -v ./...

test:
	go clean
	go test ./ ./should

cover:
	go clean
	go test -cover ./ ./should

test-watch:
	modd

# Make sure there's no debug code etc.
code-quality:
	@echo "Checking for debugging figments"
	@! grep --exclude Makefile --exclude-dir vendor -nIR 'y0ssar1an/q' *
	@! grep --exclude Makefile --exclude-dir vendor -nIR 'DEBUG' *

# package location for compiled-in values
INJECT_VARS_SITE = github.com/kindrid/gotest/gotest

# create executables for this plaform
build: init
	go build -v -ldflags "-X ${INJECT_VARS_SITE}.Version=${VERSION} -X ${INJECT_VARS_SITE}.Commit=${COMMIT}"  ./...

# create any distribution files
dist: build

release:
	github-release kindrid/gotest ${VERSION} ${BRANCH} copyTheChangeLogManually CHANGELOG.md

# Convention for our vendored builds on Semaphore
ci-build: build

# Semaphore preliminaries
ci-before: code-quality ci-build

# First semaphore job
ci-job1: code-quality test
