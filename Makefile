export GO111MODULE=on
bin/sts:
	go build -ldflags "-linkmode external -extldflags -static" -o bin/sts ./cmd/sts

.PHONY: dep
dep: 
	go tidy 

.PHONY: test
test:
	go test -v -coverprofile=cover.out ./...

.PHONY: lint
lint:
	golangci-lint run

.PHONY: clean
clean:
	rm -rf bin
	go clean github.com/illfate/social-tournaments-service

.PHONY: upgrade
upgrade:
	go get -u
