bin/sts:
	go build -mod vendor -ldflags "-linkmode external -extldflags -static" -o bin/sts ./cmd/sts

.PHONY: dep
dep: 
	go mod tidy 
	go mod vendor

.PHONY: test
test:
	go test -v -mod vendor -coverprofile=cover.out ./... 

.PHONY: lint
lint:
	golangci-lint run

.PHONY: clean
clean:
	rm -rf bin
	go clean -mod vendor
