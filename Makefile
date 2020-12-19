GO111MODULE=on

.PHONY: vendor
vendor:
	GO111MODULE=${GO111MODULE} go get ./... && \
	GO111MODULE=${GO111MODULE} go mod tidy && \
	GO111MODULE=${GO111MODULE} go mod vendor

.PHONY: lint
lint:
ifeq (, $(shell which golangci-lint))
	$(error "No golangci-lint in $(PATH). Install it from https://github.com/golangci/golangci-lint")
endif
	golangci-lint run

.PHONY: test
test:
	GO111MODULE=on go test -mod vendor -bench=. -benchmem -coverpkg=./... -covermode=count -coverprofile=coverage.out -v $(go list ./... | grep -v -e benchmarks/ -e examples/ -e metrics/) && \
    GO111MODULE=${GO111MODULE} go tool cover -func=coverage.out && \
    GO111MODULE=on go tool cover -html=coverage.out -o=coverage.html
