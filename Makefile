bin/sts:
	go build -ldflags "-linkmode external -extldflags -static" -o bin/sts ./cmd/sts

.PHONY: dep
dep: 
	go get -d github.com/go-sql-driver/mysql 
	go get -d github.com/gorilla/mux
	go get -d github.com/stretchr/testify/mock
	go get github.com/pressly/goose/cmd/goose

.PHONY: test
test:
	go test -v -coverprofile=cover.out ./...

.PHONY: lint
lint:
	golangci-lint run

.PHONY: clean
clean:
	rm -rf bin
	go clean
