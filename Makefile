bin/sts:
	go build -o bin/sts ./cmd/sts

.PHONY: dep
dep: 
	go get -d github.com/go-sql-driver/mysql 
	go get -d github.com/gorilla/mux
	go get -d github.com/stretchr/testify/mock

.PHONY: test
test:
	go test -v -coverprofile=cover.out ./...

.PHONY: lint
lint:
	golangci-lint run
