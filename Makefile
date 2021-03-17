all: fmt lint vet test

fmt:
	go fmt $(shell find -name \*.go |grep -v _examples|xargs dirname|sort -u)

lint:
	golint $(shell find -name \*.go |xargs dirname|sort -u)

test:
	go test $(shell find -name \*.go |grep -v _examples|xargs dirname|sort -u)

vet:
	go vet $(shell find -name \*.go |grep -v _examples|xargs dirname|sort -u)
