bin/sts:
	go build -o bin/sts github.com/illfate/social-tournaments-service/cmd/sts

.PHONY: dep
dep: 
	go get -d github.com/go-sql-driver/mysql 
	go get -d github.com/gorilla/mux

.PHONY: test
test:
	go test -v -coverprofile=cover.out ./...

.PHONY: lint
lint:
	golangci-lint run
